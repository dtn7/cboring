package cboring

import (
	"fmt"
	"io"
)

func ReadUint(r io.Reader) (n uint64, err error) {
	m, n, err := ParseMajorFields(r)
	if err == nil && m != UInt {
		err = fmt.Errorf("ReadUint: Wrong Major Type: %d instead of %d", m, UInt)
	}
	return
}
