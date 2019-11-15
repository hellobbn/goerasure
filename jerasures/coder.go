// This file defines the coder interface for jerasure

package jerasures

// Coder interface:
// The interface to encode and decode data, and also contains other
// useful methods
type Coder interface {
	// Encode method
	Encode(data []byte) ([][]byte, [][]byte, int, error)
	// Decode method
	Decode(encodedData, encodedParity [][]byte, blockSize int, missingIDs []int) [][]byte
}
