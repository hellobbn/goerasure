package jerasures

/*
#include <stdio.h>
char* jerasure_bytes_at(char **data, int i) {
  return data[i];
}
int jerasure_int_at(int **schedule, int i, int j) {
	return schedule[i][j];
}
*/
import "C"
import "unsafe"

// blockToC receives a Go byte matrix as input and outputs a C char matrix.
func BlockToC(data [][]byte) **C.char {
	if len(data) < 1 {
		panic("no data given")
	}

	var b *C.char
	ptrSize := unsafe.Sizeof(b)

	//Allocate the char** list
	ptr := C.malloc(C.size_t(len(data)) * C.size_t(ptrSize))

	//Assign each byte slice to its appropriate offset
	for i := 0; i < len(data); i++ {
		element := (**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*ptrSize))
		*element = (*C.char)(unsafe.Pointer(&data[i][0]))
	}

	return ((**C.char)(ptr))
}

// CToBlock received a C char matrix as input and outputs a Go byte matrix.
func CToBlock(dataC **C.char, data [][]byte) {
	for i := range data {
		data[i] = C.GoBytes(unsafe.Pointer(C.jerasure_bytes_at(dataC, C.int(i))), C.int(len(data[i])))
	}
}

//intSliceToC converts a Go slice into a C int pointer array
func IntSliceToC(slice []int) *C.int {
	return (*C.int)(unsafe.Pointer(&slice[0]))
}
