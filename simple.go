package cboring

import (
	"fmt"
	"io"
)

const (
	simpleFalse byte = 20
	simpleTrue  byte = 21
)

// ReadBoolean reads a bool value from the Reader.
func ReadBoolean(r io.Reader) (b bool, err error) {
	var buff [1]byte

	if _, dataErr := r.Read(buff[:1]); dataErr != nil {
		err = dataErr
		return
	}

	major, adds := readMajorType(buff[0])
	if major != SimpleData {
		err = fmt.Errorf("ReadBoolean: Expected major 0x%x, got 0x%x", SimpleData, major)
		return
	}

	switch adds {
	case simpleFalse:
		b = false
	case simpleTrue:
		b = true
	default:
		err = fmt.Errorf("ReadBoolean: Unknown additional 0x%x", adds)
	}

	return
}

// WriteBoolean writes a bool into the Writer.
func WriteBoolean(b bool, w io.Writer) (err error) {
	var adds byte
	if b {
		adds = simpleTrue
	} else {
		adds = simpleFalse
	}

	payload := writeMajorType(SimpleData, adds)
	_, err = w.Write([]byte{payload})

	return
}
