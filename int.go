package cboring

import (
	"fmt"
	"io"
)

// ReadUInt expects an unsigned integer at the Reader's position and returns it.
func ReadUInt(r io.Reader) (n uint64, err error) {
	m, n, err := ReadMajors(r)
	if err == nil && m != UInt {
		err = fmt.Errorf("ReadUInt: Wrong Major Type: %d instead of %d", m, UInt)
	}
	return
}

// WriteUInt serializes an unsigned integer into the Writer.
func WriteUInt(n uint64, w io.Writer) error {
	return WriteMajors(UInt, n, w)
}

// ReadNegInt expects a negative integer at the Reader's position and returns it.
func ReadNegInt(r io.Reader) (n int64, err error) {
	m, tmp, err := ReadMajors(r)
	if err == nil && m != NegInt {
		err = fmt.Errorf("ReadNegInt: Wrong Major Type: %d instead of %d", m, NegInt)
	}
	n = -1 - int64(tmp)
	return
}

// WriteNegInt serializes a negative integer into the Writer.
func WriteNegInt(n int64, w io.Writer) error {
	if n >= 0 {
		return fmt.Errorf("WriteNegInt: Expected negative integer")
	}

	return WriteMajors(NegInt, uint64((n+1)*(-1)), w)
}
