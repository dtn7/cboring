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

func ReadNegInt(r io.Reader) (n int64, err error) {
	m, tmp, err := ReadMajorFields(r)
	if err == nil && m != NegInt {
		err = fmt.Errorf("ReadNegInt: Wrong Major Type: %d instead of %d", m, NegInt)
	}
	n = -1 - int64(tmp)
	return
}
