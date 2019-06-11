package cboring

import (
	"fmt"
	"io"
)

func readString(len int, r io.Reader) (data []byte, err error) {
	data = make([]byte, len)
	if rn, rerr := r.Read(data); err != nil {
		err = rerr
	} else if rn != len {
		err = fmt.Errorf("readString: read length mismatches: %d != %d", rn, len)
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

	return readString(int(n), r)
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

	if rdata, rerr := readString(int(n), r); rerr != nil {
		err = rerr
	} else {
		data = string(rdata)
	}
	return
}

// WriteTextString writes a byte string into the Writer.
func WriteTextString(data string, w io.Writer) error {
	if err := WriteTextStringLen(uint64(len(data)), w); err != nil {
		return err
	}

	// WriteString instead of w.Write to save a cast
	if n, err := io.WriteString(w, data); err != nil {
		return err
	} else if n != len(data) {
		return fmt.Errorf("WriteTextString: Wrote %d instead of %d bytes",
			n, len(data))
	}
	return nil
}
