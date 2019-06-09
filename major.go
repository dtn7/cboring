package cboring

import (
	"fmt"
	"io"
)

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

func ReadMajorType(b byte) (major MajorType, adds byte) {
	major = b >> 5
	adds = b & 0x1F
	return
}

func ReadMajorFields(r io.Reader) (m MajorType, n uint64, err error) {
	var buff [8]byte
	tmpBuff := buff[:1]

	if _, rerr := r.Read(tmpBuff); rerr != nil {
		err = rerr
		return
	}

	m, adds := ReadMajorType(tmpBuff[0])
	if adds <= 23 {
		n = uint64(adds)
	} else if 24 <= adds && adds <= 27 {
		l := 1 << (adds - 24)
		tmpBuff = buff[:l]

		if rn, rerr := r.Read(tmpBuff); rerr != nil {
			err = rerr
			return
		} else if rn != l {
			err = fmt.Errorf("ReadMajorFields: Read %d bytes instead of %d", rn, l)
			return
		}

		for i := 0; i < l; i++ {
			n = n<<8 | uint64(tmpBuff[i])
		}
	} else {
		err = fmt.Errorf("ReadMajorFields: Other additional information %d", adds)
	}

	return
}

func WriteMajorType(major MajorType, adds byte) byte {
	return (major << 5) | (adds & 0x1F)
}

func WriteMajorFields(m MajorType, n uint64, w io.Writer) (err error) {
	var buff [9]byte
	var bc = 0

	if n < 24 {
		buff[0] = WriteMajorType(m, byte(n))
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

		buff[0] = WriteMajorType(m, mt)
		for i := bc; i > 0; i-- {
			buff[i] = byte(n & 0xFF)
			n = n >> 8
		}
	}

	if wn, werr := w.Write(buff[:bc+1]); werr != nil {
		err = werr
	} else if wn != bc+1 {
		err = fmt.Errorf("WriteMajorFields: Wrote %d instead of %d bytes", wn, bc+1)
	}
	return
}
