package jerasures

// #include <jerasure.h>
// #include <jerasure/reed_sol.h>
// #include <stdlib.h>
// #cgo CFLAGS: -I/usr/local/include/jerasure/
// #cgo LDFLAGS: -lJerasure
import "C"

import (
	"unsafe"

	"github.com/hellobbn/goerasure/utils"
)

//structure for ReedSol Code with Vandermonde Matrix
// ReedSolVand implements Coder
type ReedSolVand struct {
	matrix (*C.int)
	k      int
	m      int
	w      int
	blkSize int	// saves blkint
}

// NewReedSolVand:
// Creates a new ReedSolVand Object
// Using this function, w is by default set to 8
func NewReedSolVand(k, m int) ReedSolVand {
	rscode := ReedSolVand{
		k: k,
		m: m,
		w: 8, // by default, TODO: try to be more specific
	}
	rscode.matrix = C.reed_sol_vandermonde_coding_matrix(C.int(k), C.int(m), C.int(rscode.w))
	return rscode
}

// Encode:
// Encodes dara using RS and Vandermonde matrix
func (rsCode ReedSolVand) Encode(data []byte) ([][]byte, [][]byte, int, error) {
	edBytes, epBytes, blockSize := utils.PrepareDataForEncode(rsCode.k, rsCode.m, rsCode.w, data)

	edC := BlockToC(edBytes)
	epC := BlockToC(epBytes)

	C.jerasure_matrix_encode(C.int(rsCode.k), C.int(rsCode.m), C.int(rsCode.w),
		rsCode.matrix,
		edC,
		epC,
		C.int(blockSize))

	CToBlock(edC, edBytes)
	CToBlock(epC, epBytes)

	C.free(unsafe.Pointer(edC))
	C.free(unsafe.Pointer(epC))

	// TODO: save blockSize
	rsCode.blkSize = blockSize
	return edBytes, epBytes, blockSize, nil
}

// Decode:
// decodes data
func (rsCode ReedSolVand) Decode(encodedData, encodedParity [][]byte, blockSize int, missingIDs []int) [][]byte {
	edC := BlockToC(encodedData)
	epC := BlockToC(encodedParity)

	missingIDsC := IntSliceToC(missingIDs)

	ret := C.jerasure_matrix_decode(C.int(rsCode.k), C.int(rsCode.m), C.int(rsCode.w),
		rsCode.matrix, 1,
		missingIDsC,
		edC,
		epC,
		C.int(blockSize))

	if ret == -1 {
		panic("ERROR decoding.")
	}

	CToBlock(edC, encodedData)
	CToBlock(epC, encodedParity)
	C.free(unsafe.Pointer(edC))
	C.free(unsafe.Pointer(epC))

	return encodedData
}

func (rsv ReedSolVand) DatBlks() int {
	return rsv.k;
}

func (rsv ReedSolVand) PariBlks() int {
	return rsv.m;
}

func (rsv ReedSolVand) BlkSize() int {
	return rsv.blkSize;
}
