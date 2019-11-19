package jerasures

import (
	"log"
	"reflect"
	"testing"
)

func TestReedSolVand(t *testing.T) {
	var origData []byte
	for i := 0; i < 124; i++ {
		origData = append(origData, []byte("aa")...)
	}
	rsv := NewReedSolVand(16, 4)
	encodedData, encodedParity, blockSize, _ := rsv.Encode(origData)
	log.Printf("blockSize = %v, k = %v \n", blockSize, rsv.k)
	missingIDs := []int{0, rsv.k - 1, -1}
	encodedData[0] = make([]byte, blockSize)
	encodedData[rsv.k - 1] = make([]byte, blockSize)

	// FIXME: wrong testing here
	recoveredData := rsv.Decode(encodedData, encodedParity, blockSize, missingIDs)

	if !reflect.DeepEqual(origData, recoveredData[:len(origData)]) {
		log.Print(recoveredData[:len(origData)])
		log.Print(origData)
		t.Fatalf("failed to decode data")
	}
}
