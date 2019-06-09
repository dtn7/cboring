package cboring

import (
	"fmt"
	"io"
)

func ReadUInt(r io.Reader) (n uint64, err error) {
	m, n, err := ReadMajorFields(r)
	if err == nil && m != UInt {
		err = fmt.Errorf("ReadUInt: Wrong Major Type: %d instead of %d", m, UInt)
	}
	return
}

func WriteUInt(n uint64, w io.Writer) error {
	return WriteMajorFields(UInt, n, w)
}

func ReadNegInt(r io.Reader) (n int64, err error) {
	m, tmp, err := ReadMajorFields(r)
	if err == nil && m != NegInt {
		err = fmt.Errorf("ReadNegInt: Wrong Major Type: %d instead of %d", m, NegInt)
	}
	n = -1 - int64(tmp)
	return
}

func WriteNegInt(n int64, w io.Writer) error {
	if n >= 0 {
		return fmt.Errorf("WriteNegInt: Expected negative integer")
	}

	return WriteMajorFields(NegInt, uint64((n+1)*(-1)), w)
}
