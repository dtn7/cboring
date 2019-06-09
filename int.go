package cboring

import (
	"fmt"
	"io"
)

func ReadUint(r io.Reader) (n uint64, err error) {
	var buff [8]byte
	tmpBuff := buff[:1]

	if rn, rerr := r.Read(tmpBuff); rerr != nil {
		err = rerr
		return
	} else if rn != 1 {
		err = fmt.Errorf("ReadUint: Read %d bytes instead of %d", rn, 1)
		return
	}

	major, adds := ParseMajor(tmpBuff[0])
	if major != UInt {
		err = fmt.Errorf("ReadUint: Wrong Major Type: %d instead of %d", major, UInt)
		return
	}

	if adds <= 23 {
		n = uint64(adds)
	} else if 24 <= adds && adds <= 27 {
		l := 1 << (adds - 24)
		tmpBuff = buff[:l]

		if rn, rerr := r.Read(tmpBuff); rerr != nil {
			err = rerr
			return
		} else if rn != l {
			err = fmt.Errorf("ReadUint: Read %d bytes instead of %d", rn, l)
			return
		}

		for i := 0; i < l; i++ {
			n = n<<8 | uint64(tmpBuff[i])
		}
	} else {
		err = fmt.Errorf("ReadUint: Unknown additional information %d", adds)
	}

	return
}
