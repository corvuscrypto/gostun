package gostun

import (
	"encoding/binary"
	"sync/atomic"
	"time"
)

var counter uint32

func createNonce(sst int64) []byte {

	timestampBytes := make([]byte, 4)
	timestamp := time.Now().Unix()
	binary.BigEndian.PutUint32(timestampBytes, uint32(timestamp))

	//make the first 8 bytes of the nonce
	nonce := make([]byte, 8)
	binary.BigEndian.PutUint64(nonce, uint64(sst^timestamp))

	cnt := make([]byte, 4)
	binary.BigEndian.PutUint32(cnt, atomic.AddUint32(&counter, 1))

	//append the timestamp to the end and return
	return append(nonce, append(timestampBytes, cnt...)...)
}

func nonceValid(nonce []byte, sst int64, timeout int) bool {
	if len(nonce) != 16 {
		return false
	}
	timestamp := uint64(binary.BigEndian.Uint32(nonce[8:12]))
	if uint64(time.Now().Unix())-timestamp > uint64(timeout) {
		return false
	}
	fix := nonce[:8]
	if binary.BigEndian.Uint64(fix)^timestamp != uint64(sst) {
		return false
	}
	return true
}
