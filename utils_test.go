package gostun

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestNonces(T *testing.T) {

	//use max processors to ensure that cross-core generation doesn't create collisions
	runtime.GOMAXPROCS(runtime.NumCPU())

	sst := time.Now().UnixNano()
	timeout := 40 //seconds

	//generate 10,000 nonces. They should not be equal to each other
	nonces := make([][]byte, 10000)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func(n [][]byte, i int) {
			n[i] = createNonce(sst)
			wg.Done()
		}(nonces, i)
	}
	wg.Wait()

	//check to see if all nonces are valid
	for _, v := range nonces {
		if !nonceValid(v, sst, timeout) {
			T.Errorf("Encountered invalid Nonce!")
			break
		}
	}

	for i := 0; i < 10000; i++ {
		for j := i + 1; j < 10000; j++ {
			if string(nonces[i]) == string(nonces[j]) {
				T.Errorf("Encountered collision on Nonce #%d with Nonce #%d", i, j)
				return
			}
		}
	}

}
