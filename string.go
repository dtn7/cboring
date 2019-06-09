package cboring

import (
	"fmt"
	"io"
)

// ReadByteStringLen expects a byte string at the Reader's position and returns
// the byte string's length.
func ReadByteStringLen(r io.Reader) (n uint64, err error) {
	m, n, err := ReadMajors(r)
	if err == nil && m != ByteString {
		err = fmt.Errorf("ReadByteStringLen: Wrong Major Type: %d instead of %d",
			m, ByteString)
	}
	return
}

// ReadByteStringLen expects a byte string at the Reader's position and returns
// the byte string.
func ReadByteString(r io.Reader) (data []byte, err error) {
	n, err := ReadByteStringLen(r)
	if err != nil {
		return
	}

	data = make([]byte, n)
	if rn, rerr := r.Read(data); err != nil {
		err = rerr
	} else if rn != int(n) {
		err = fmt.Errorf("ReadByteString: read length mismatches: %d != %d", rn, n)
	}
	return
}

// ReadTextStringLen expects a text string at the Reader's position and returns
// the text string's length.
func ReadTextStringLen(r io.Reader) (n uint64, err error) {
	m, n, err := ReadMajors(r)
	if err == nil && m != TextString {
		err = fmt.Errorf("TextByteStringLen: Wrong Major Type: %d instead of %d",
			m, TextString)
	}
	return
}

// ReadTextStringLen expects a text string at the Reader's position and returns
// the text string.
func ReadTextString(r io.Reader) (data string, err error) {
	n, err := ReadTextStringLen(r)
	if err != nil {
		return
	}

	tmpData := make([]byte, n)
	if rn, rerr := r.Read(tmpData); err != nil {
		err = rerr
	} else if rn != int(n) {
		err = fmt.Errorf("ReadTextString: read length mismatches: %d != %d", rn, n)
	}

	data = string(tmpData)
	return
}
