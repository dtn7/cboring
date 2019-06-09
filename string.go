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

// WriteByteStringLen writes the type definition for a byte string with the
// given length into the Writer.
func WriteByteStringLen(n uint64, w io.Writer) error {
	return WriteMajors(ByteString, n, w)
}

// WriteByteString writes a byte string into the Writer.
func WriteByteString(data []byte, w io.Writer) error {
	if err := WriteByteStringLen(uint64(len(data)), w); err != nil {
		return err
	}

	if n, err := w.Write(data); err != nil {
		return err
	} else if n != len(data) {
		return fmt.Errorf("WriteByteString: Wrote %d instead of %d bytes",
			n, len(data))
	}
	return nil
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

// WriteTextStringLen writes the type definition for a text string with the
// given length into the Writer.
func WriteTextStringLen(n uint64, w io.Writer) error {
	return WriteMajors(TextString, n, w)
}

// WriteTextString writes a byte string into the Writer.
func WriteTextString(data string, w io.Writer) error {
	if err := WriteTextStringLen(uint64(len(data)), w); err != nil {
		return err
	}

	if n, err := w.Write([]byte(data)); err != nil {
		return err
	} else if n != len(data) {
		return fmt.Errorf("WriteTextString: Wrote %d instead of %d bytes",
			n, len(data))
	}
	return nil
}
