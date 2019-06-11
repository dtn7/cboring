package cboring

import (
	"fmt"
	"io"
)

// MajorType defines a Major Type, as specified in RFC7049, section 2.1
type MajorType = byte

const (
	UInt       MajorType = 0
	NegInt     MajorType = 1
	ByteString MajorType = 2
	TextString MajorType = 3
	Array      MajorType = 4
	Map        MajorType = 5
	Tagging    MajorType = 6
	Etc        MajorType = 7
)

type flag string

func (f flag) Error() string {
	return string(f)
}

const (
	IndefiniteLength = flag("Indefinite Length Array")
	BreakCode        = flag("Break Stop Code")
)

func readMajorType(b byte) (major MajorType, adds byte) {
	major = b >> 5
	adds = b & 0x1F
	return
}

// ReadMajors parses a (major) type definition from the Reader.
func ReadMajors(r io.Reader) (m MajorType, n uint64, err error) {
	var buff [8]byte
	tmpBuff := buff[:1]

	if _, rerr := r.Read(tmpBuff); rerr != nil {
		err = rerr
		return
	}

	m, adds := readMajorType(tmpBuff[0])
	if adds <= 23 {
		n = uint64(adds)
	} else if 24 <= adds && adds <= 27 {
		l := 1 << (adds - 24)
		tmpBuff = buff[:l]

		if rn, rerr := r.Read(tmpBuff); rerr != nil {
			err = rerr
			return
		} else if rn != l {
			err = fmt.Errorf("ReadMajors: Read %d bytes instead of %d", rn, l)
			return
		}

		for i := 0; i < l; i++ {
			n = n<<8 | uint64(tmpBuff[i])
		}
	} else if adds == 31 && m == Array {
		err = IndefiniteLength
	} else if adds == 31 && m == Etc {
		err = BreakCode
	} else {
		err = fmt.Errorf("ReadMajors: Other additional information %d", adds)
	}

	return
}

func writeMajorType(major MajorType, adds byte) byte {
	return (major << 5) | (adds & 0x1F)
}

// WriteMajors composes a (major) type definition into the Writer.
func WriteMajors(m MajorType, n uint64, w io.Writer) (err error) {
	var buff [9]byte
	var bc = 0

	if n < 24 {
		buff[0] = writeMajorType(m, byte(n))
	} else {
		var mt uint8
		if n < 1<<8 {
			bc = 1
			mt = 24
		} else if n < 1<<16 {
			bc = 2
			mt = 25
		} else if n < 1<<32 {
			bc = 4
			mt = 26
		} else {
			bc = 8
			mt = 27
		}

		buff[0] = writeMajorType(m, mt)
		for i := bc; i > 0; i-- {
			buff[i] = byte(n & 0xFF)
			n = n >> 8
		}
	}

	if wn, werr := w.Write(buff[:bc+1]); werr != nil {
		err = werr
	} else if wn != bc+1 {
		err = fmt.Errorf("WriteMajors: Wrote %d instead of %d bytes", wn, bc+1)
	}
	return
}
