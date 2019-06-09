package cboring

import (
	"bytes"
	"testing"
)

func TestReadUintSmall(t *testing.T) {
	for i := uint64(0); i <= 23; i++ {
		r := bytes.NewBuffer([]byte{byte(i)})
		if n, err := ReadUint(r); err != nil {
			t.Fatal(err)
		} else if n != i {
			t.Fatalf("Resulted uint %d is not %d", n, i)
		}
	}
}

func TestReadUintBig(t *testing.T) {
	tests := []struct {
		data []byte
		numb uint64
	}{
		{[]byte{0x18, 0x18}, 24},
		{[]byte{0x18, 0x19}, 25},
		{[]byte{0x18, 0x64}, 100},
		{[]byte{0x19, 0x03, 0xe8}, 1000},
		{[]byte{0x1a, 0x00, 0x0f, 0x42, 0x40}, 1000000},
		{[]byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00}, 1000000000000},
		{[]byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 18446744073709551615},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test.data)
		if n, err := ReadUint(r); err != nil {
			t.Fatal(err)
		} else if n != test.numb {
			t.Fatalf("Resulted uint %d is not %d", n, test.numb)
		}
	}
}

func TestReadUintMultiple(t *testing.T) {
	numbs := []uint64{0, 1000, 24, 25, 1000000000000}
	data := []byte{
		0x00,
		0x19, 0x03, 0xe8,
		0x18, 0x18,
		0x18, 0x19,
		0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00,
	}

	r := bytes.NewBuffer(data)
	for _, numb := range numbs {
		if n, err := ReadUint(r); err != nil {
			t.Fatal(err)
		} else if n != numb {
			t.Fatalf("Resulted uint %d is not %d", n, numb)
		}
	}
}

func TestReadUintError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		[]byte{0xFF},
		// Wrong additionals for major type 0
		[]byte{0x1F},
		// Empty stream
		[]byte{},
		// Incomplete streams
		[]byte{0x18},
		[]byte{0x19, 0x03},
		[]byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadUint(r); err == nil {
			t.Fatalf("Wrong type %x did not errored", test)
		}
	}
}
