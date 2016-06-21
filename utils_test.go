package gostun

import (
	"encoding/binary"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var timeout = 40 //seconds

//used to create nonce that is beyond timeout time
func createLaggedNonce(sst int64) []byte {

	timestampBytes := make([]byte, 4)
	timestamp := time.Now().Unix() - 41
	binary.BigEndian.PutUint32(timestampBytes, uint32(timestamp))

	//make the first 8 bytes of the nonce
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, uint64(sst^timestamp))

	cnt := make([]byte, 4)
	binary.BigEndian.PutUint32(cnt, atomic.AddUint32(&counter, 1))

	//append the timestamp to the end and return
	return append(nonce, append(timestampBytes, cnt...)...)
}

func TestNonces(T *testing.T) {

	//use max processors to ensure that cross-core generation doesn't create collisions
	runtime.GOMAXPROCS(runtime.NumCPU())

	sst := time.Now().UnixNano()

	invalidNonces := [][]byte{
		{1, 2, 4, 6, 32, 12},
		{2, 4, 6, 32, 12, 223, 12, 32, 12, 12, 233, 41, 53, 11, 100, 23},
		createNonce(12333212),
		append(createNonce(sst)[:10], []byte{255, 255, 0, 0, 0, 0}...),
		{20, 90, 7, 235, 3, 229, 48, 247, 87, 104, 233, 151, 0, 0, 0, 1},
		createLaggedNonce(sst),
	}

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

	//check to see if all invalid nonces are indeed invalid
	for _, v := range invalidNonces {
		if nonceValid(v, sst, timeout) {
			T.Errorf("Encountered valid Nonce when expecting invalid nonce!")
			fmt.Println(v)
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
