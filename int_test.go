package cboring

import (
	"bytes"
	"reflect"
	"testing"
)

func TestUInt(t *testing.T) {
	tests := []struct {
		data []byte
		numb uint64
	}{
		{[]byte{0x00}, 0},
		{[]byte{0x01}, 1},
		{[]byte{0x0a}, 10},
		{[]byte{0x17}, 23},
		{[]byte{0x18, 0x18}, 24},
		{[]byte{0x18, 0x19}, 25},
		{[]byte{0x18, 0x64}, 100},
		{[]byte{0x19, 0x03, 0xe8}, 1000},
		{[]byte{0x1a, 0x00, 0x0f, 0x42, 0x40}, 1000000},
		{[]byte{0x1b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x10, 0x00}, 1000000000000},
		{[]byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 18446744073709551615},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if n, err := ReadUInt(buff); err != nil {
			t.Fatal(err)
		} else if n != test.numb {
			t.Fatalf("Resulting uint %d is not %d", n, test.numb)
		}

		// Write
		buff.Reset()
		if err := WriteUInt(test.numb, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestNegInt(t *testing.T) {
	tests := []struct {
		data []byte
		numb int64
	}{
		{[]byte{0x20}, -1},
		{[]byte{0x29}, -10},
		{[]byte{0x38, 0x63}, -100},
		{[]byte{0x39, 0x03, 0xe7}, -1000},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if n, err := ReadNegInt(buff); err != nil {
			t.Fatal(err)
		} else if n != test.numb {
			t.Fatalf("Resulting int %d is not %d", n, test.numb)
		}

		// Write
		buff.Reset()
		if err := WriteNegInt(test.numb, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestReadUIntError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		[]byte{0xFF},
		// Wrong additionals for major type 0
		[]byte{0x1F},
		// Empty stream
		[]byte{},
		// Incomplete streams
		[]byte{0x18}, []byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadUInt(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

func TestReadNegIntError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		[]byte{0xFF},
		// Wrong additionals for major type 0
		[]byte{0x3F},
		// Empty stream
		[]byte{},
		// Incomplete streams
		[]byte{0x38}, []byte{0x39, 0x03},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadNegInt(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}
