package cboring

import (
	"fmt"
	"io"
)

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
