package cboring

import (
	"fmt"
	"io"
)

// MajorType defines a Major Type, as specified in RFC7049, section 2.1
type MajorType = byte

const (
	UInt       MajorType = 0x00
	NegInt     MajorType = 0x20
	ByteString MajorType = 0x40
	TextString MajorType = 0x60
	Array      MajorType = 0x80
)

const (
	IndefiniteArray byte = 0x9F
	BreakCode       byte = 0xFF
)

type flag byte

func (f flag) Error() string {
	return string(f)
}

const (
	flagIndefiniteArray = flag(iota)
	flagBreakCode       = flag(iota)
)

func readMajorType(b byte) (major MajorType, adds byte) {
	major = b & 0xE0
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

	switch b := tmpBuff[0]; b {
	case IndefiniteArray:
		err = flagIndefiniteArray

	case BreakCode:
		err = flagBreakCode

	default:
		var adds byte
		m, adds = readMajorType(b)

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
		} else {
			err = fmt.Errorf("ReadMajors: Other additional information %d", adds)
		}
	}

	return
}

// ReadExpectMajors parses the next (major) type, which must equal the requested
// one. This function wraps ReadMajors.
func ReadExpectMajors(m MajorType, r io.Reader) (n uint64, err error) {
	mTmp, n, err := ReadMajors(r)
	if err == nil && m != mTmp {
		err = fmt.Errorf("ReadExpectMajors: Wrong Major Type: 0x%x instead of 0x%x",
			m, mTmp)
	}
	return
}

func writeMajorType(major MajorType, adds byte) byte {
	return major | adds
}

// WriteMajors composes a (major) type definition into the Writer.
func WriteMajors(m MajorType, n uint64, w io.Writer) (err error) {
	var buff [9]byte
	var bc = 0

	if n < 24 {
		buff[0] = writeMajorType(m, byte(n))
	} else {
		var mt byte
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
