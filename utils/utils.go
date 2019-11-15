package utils

import "C"
import "unsafe"

func ceill(f float64) int {
	return int(f + 0.9)
}

func getAlignedDataSize(k, w int, dataLen int) int {
	wordSize := w / 8
	alignmentMultiple := k * wordSize
	return ceill(float64(dataLen)/float64(alignmentMultiple)) * alignmentMultiple
}

func PrepareDataForEncode(k, m, w int, data []byte) ([][]byte, [][]byte, int) {
	dataLen := len(data)
	alignedDataLen := getAlignedDataSize(k, w, dataLen)

	blockSize := alignedDataLen / k
	payloadSize := blockSize

	// prepare encode data
	encodedData := make([][]byte, k)
	cursor := 0
	for i := 0; i < k; i++ {
		copySize := payloadSize
		if dataLen < payloadSize {
			copySize = dataLen
		}

		if dataLen > 0 {
			encodedData[i] = data[cursor : cursor+copySize]
		}

		cursor += copySize
		dataLen -= copySize
	}

	// encode parity
	encodedParity := make([][]byte, m)
	for k := 0; k < m; k++ {
		encodedParity[k] = make([]byte, blockSize)
	}

	return encodedData, encodedParity, blockSize
}

// ConvertResultData
// converts returned result data to byte
func ConvertResultData(ed [][]byte, blockSize int) []byte {
	data := []byte{}

	for _, d := range ed {
		data = append(data, d...)
	}
	return data
}

func ConvertResultDatas(ed [](*C.char), blockSize int) []byte {
	data := []byte{}

	for _, d := range ed {
		data = append(data, C.GoBytes(unsafe.Pointer(d), C.int(blockSize))...)
	}
	return data
}

func PrepareDataForDecode(k, m int, encodedData, encodedParity [][]byte) ([]*C.char, []*C.char) {
	ed := make([](*C.char), k)
	for i, v := range encodedData {
		ed[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}
	ep := make([](*C.char), m)
	for i, v := range encodedParity {
		ep[i] = (*C.char)(unsafe.Pointer(&v[0]))
	}
	return ed, ep
}
