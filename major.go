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

func ParseMajorType(b byte) (major MajorType, adds byte) {
	major = b >> 5
	adds = b & 0x1F
	return
}

func ParseMajorFields(r io.Reader) (m MajorType, n uint64, err error) {
	var buff [8]byte
	tmpBuff := buff[:1]

	if _, rerr := r.Read(tmpBuff); rerr != nil {
		err = rerr
		return
	}

	m, adds := ParseMajorType(tmpBuff[0])
	if adds <= 23 {
		n = uint64(adds)
	} else if 24 <= adds && adds <= 27 {
		l := 1 << (adds - 24)
		tmpBuff = buff[:l]

		if rn, rerr := r.Read(tmpBuff); rerr != nil {
			err = rerr
			return
		} else if rn != l {
			err = fmt.Errorf("ParseMajorFields: Read %d bytes instead of %d", rn, l)
			return
		}

		for i := 0; i < l; i++ {
			n = n<<8 | uint64(tmpBuff[i])
		}
	} else {
		err = fmt.Errorf("ParseMajorFields: Other additional information %d", adds)
	}

	return
}
