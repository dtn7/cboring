package cboring

import (
	"fmt"
	"io"
)

// ReadArrayLength expects an array at the Reader's position and returns its
// length.
func ReadArrayLength(r io.Reader) (n uint64, err error) {
	m, n, err := ReadMajors(r)
	if err == nil && m != Array {
		err = fmt.Errorf("ReadArrayLength: Wrong Major Type: %d instead of %d",
			m, Array)
	}
	return
}

// WriteArrayLength writes the type definition for an array with the given
// length into the Writer.
func WriteArrayLength(n uint64, w io.Writer) error {
	return WriteMajors(Array, n, w)
}
