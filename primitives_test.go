package cboring

import (
	"bytes"
	"reflect"
	"testing"
)

/*** UInt ***/

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
		buff := &bytes.Buffer{}
		_, _ = buff.Write(test.data)
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

func TestReadUIntError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		{0xFF},
		// Wrong additionals for major type 0
		{0x1F},
		// Empty stream
		{},
		// Incomplete streams
		{0x18}, {0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadUInt(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

/*** NInt ***/

func TestNInt(t *testing.T) {
	tests := []struct {
		data []byte
		numb int64
	}{
		{[]byte{0x20}, -1},
		{[]byte{0x29}, -10},
		{[]byte{0x36}, -23},
		{[]byte{0x37}, -24},
		{[]byte{0x38, 0x18}, -25},
		{[]byte{0x38, 0x63}, -100},
		{[]byte{0x39, 0x03, 0xe7}, -1000},
		{[]byte{0x3a, 0x00, 0x0f, 0x42, 0x3F}, -1000000},
		{[]byte{0x3b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x0F, 0xFF}, -1000000000000},
		{[]byte{0x3b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, 0}, /* can not represent -2^64
		but underflow would lead to 0, so this is the best representation */
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if n, err := ReadNInt(buff); err != nil {
			t.Fatal(err)
		} else if int64(^n) != test.numb {
			t.Fatalf("Resulting negative int %d is not %d", int64(^n), test.numb)
		}

		// Write
		buff.Reset()
		if err := WriteNInt(uint64(^test.numb), buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestReadNIntError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		{0xFF},
		// Wrong additionals for major type 1
		{NInt | 0x1F},
		// Empty stream
		{},
		// Incomplete streams
		{NInt | 0x18}, {0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadNInt(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

/*** Int ***/

func TestInt(t *testing.T) {
	tests := []struct {
		data []byte
		numb int64
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
		{[]byte{0x20}, -1},
		{[]byte{0x29}, -10},
		{[]byte{0x36}, -23},
		{[]byte{0x37}, -24},
		{[]byte{0x38, 0x18}, -25},
		{[]byte{0x38, 0x63}, -100},
		{[]byte{0x39, 0x03, 0xe7}, -1000},
		{[]byte{0x3a, 0x00, 0x0f, 0x42, 0x3F}, -1000000},
		{[]byte{0x3b, 0x00, 0x00, 0x00, 0xe8, 0xd4, 0xa5, 0x0F, 0xFF}, -1000000000000},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if n, err := ReadInt(buff); err != nil {
			t.Fatal(err)
		} else if n != test.numb {
			t.Fatalf("Resulting int %d is not %d", n, test.numb)
		}

		// Write
		buff.Reset()
		if err := WriteInt(test.numb, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestReadIntError(t *testing.T) {
	tests := [][]byte{
		// Wrong major type
		{0xFF},
		// Wrong additionals for major type 0
		{0x1F},
		// Wrong additionals for major type 1
		{NInt | 0x1F},
		// Empty stream
		{},
		// Incomplete streams
		{0x18}, {0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		// Number to big
		{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		// Number to small
		{NInt | 0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, err := ReadInt(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

/*** ByteString ***/

func TestByteStringLen(t *testing.T) {
	tests := []struct {
		data []byte
		len  uint64
	}{
		{[]byte{0x40}, 0},
		{[]byte{0x44}, 4},
		{[]byte{0x58, 0x37}, 55},
		{[]byte{0x59, 0x0A, 0x50}, 2640},
	}

	for _, test := range tests {
		// Read
		buff := &bytes.Buffer{}
		_, _ = buff.Write(test.data)
		if n, err := ReadByteStringLen(buff); err != nil {
			t.Fatal(err)
		} else if n != test.len {
			t.Fatalf("Resulting length %d is not %d", n, test.len)
		}

		// Write
		buff.Reset()
		if err := WriteByteStringLen(test.len, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

/*** TextString ***/

func TestTextStringLen(t *testing.T) {
	tests := []struct {
		data []byte
		len  uint64
	}{
		{[]byte{0x60}, 0},
		{[]byte{0x61}, 1},
		{[]byte{0x78, 0x1A}, 26},
		{[]byte{0x79, 0x07, 0xD0}, 2000},
	}

	for _, test := range tests {
		// Read
		buff := &bytes.Buffer{}
		_, _ = buff.Write(test.data)
		if n, err := ReadTextStringLen(buff); err != nil {
			t.Fatal(err)
		} else if n != test.len {
			t.Fatalf("Resulting length %d is not %d", n, test.len)
		}

		// Write
		buff.Reset()
		if err := WriteTextStringLen(test.len, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

/*** Array ***/

func TestArrayLen(t *testing.T) {
	tests := []struct {
		data []byte
		len  uint64
	}{
		{[]byte{0x80}, 0},
		{[]byte{0x81}, 1},
		{[]byte{0x98, 0x19}, 25},
	}

	for _, test := range tests {
		// Read
		buff := &bytes.Buffer{}
		_, _ = buff.Write(test.data)
		if n, err := ReadArrayLength(buff); err != nil {
			t.Fatal(err)
		} else if n != test.len {
			t.Fatalf("Resulting length %d is not %d", n, test.len)
		}

		// Write
		buff.Reset()
		if err := WriteArrayLength(test.len, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

/*** Map ***/

func TestMapPairLen(t *testing.T) {
	tests := []struct {
		data []byte
		len  uint64
	}{
		{[]byte{0xA0}, 0},
		{[]byte{0xA1}, 1},
		{[]byte{0xB8, 0x19}, 25},
	}

	for _, test := range tests {
		// Read
		buff := &bytes.Buffer{}
		_, _ = buff.Write(test.data)
		if n, err := ReadMapPairLength(buff); err != nil {
			t.Fatal(err)
		} else if n != test.len {
			t.Fatalf("Resulting length %d is not %d", n, test.len)
		}

		// Write
		buff.Reset()
		if err := WriteMapPairLength(test.len, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}
