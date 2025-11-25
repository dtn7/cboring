package cboring

import (
	"bufio"
	"bytes"
	"io"
	"reflect"
	"testing"

	"pgregory.net/rapid"
)

func TestReadMajorsSmall(t *testing.T) {
	tests := []MajorType{UInt, ByteString, TextString, Array}

	for _, test := range tests {
		for i := uint64(0); i <= 23; i++ {
			r := bytes.NewBuffer([]byte{test | byte(i)})
			if m, n, err := ReadMajors(r); err != nil {
				t.Fatal(err)
			} else if m != test {
				t.Fatalf("Resulting type %d mismatches %d", m, test)
			} else if n != i {
				t.Fatalf("Resulting uint %d is not %d", n, i)
			}
		}
	}
}

func TestReadMajorsBig(t *testing.T) {
	tests := []struct {
		data  []byte
		major MajorType
		numb  uint64
	}{
		{[]byte{0x18, 0x64}, UInt, 100},
		{[]byte{0x58, 0x40}, ByteString, 64},
		{[]byte{0x78, 0x20}, TextString, 32},
		{[]byte{0x98, 0x19}, Array, 25},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test.data)
		if m, n, err := ReadMajors(r); err != nil {
			t.Fatal(err)
		} else if m != test.major {
			t.Fatalf("Resulting type %d mismatches %d", m, test.major)
		} else if n != test.numb {
			t.Fatalf("Resulting uint %d is not %d", n, test.numb)
		}
	}
}

func TestReadMajorsError(t *testing.T) {
	tests := [][]byte{
		// Empty stream
		{},
		// Incomplete streams
		{0x18}, {0x19, 0x03},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, _, err := ReadMajors(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

func TestWriteMajorsSmall(t *testing.T) {
	tests := []MajorType{UInt, ByteString, TextString, Array}

	for _, test := range tests {
		for i := uint64(0); i <= 23; i++ {
			var buff bytes.Buffer

			if err := WriteMajors(test, i, &buff); err != nil {
				t.Fatal(err)
			}

			if m, n, err := ReadMajors(&buff); err != nil {
				t.Fatal(err)
			} else if m != test {
				t.Fatalf("Resulting type %d mismatches %d", m, test)
			} else if n != i {
				t.Fatalf("Resulting uint %d is not %d", n, i)
			}
		}
	}
}

func TestWriteMajorsBig(t *testing.T) {
	tests := []struct {
		data  []byte
		major MajorType
		numb  uint64
	}{
		{[]byte{0x18, 0x64}, UInt, 100},
		{[]byte{0x58, 0x40}, ByteString, 64},
		{[]byte{0x78, 0x20}, TextString, 32},
		{[]byte{0x98, 0x19}, Array, 25},
	}

	for _, test := range tests {
		var buff bytes.Buffer

		if err := WriteMajors(test.major, test.numb, &buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestReadExpect(t *testing.T) {
	var buff bytes.Buffer
	for i := byte(0); i < 255; i++ {
		// Read correct input
		buff.Reset()
		buff.WriteByte(i)

		if err := ReadExpect(i, &buff); err != nil {
			t.Fatalf("ReadExpect errored for %d: %v", i, err)
		}

		// Read invalid input
		buff.Reset()
		buff.WriteByte(i + 1)

		if err := ReadExpect(i, &buff); err == nil {
			t.Fatalf("ReadExpect did not errored for %d", i)
		}
	}
}

func TestReadExampleArray(t *testing.T) {
	var data = []byte{
		0x9f, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b,
		0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17,
		0x18, 0x18, 0x18, 0x19, 0xff}

	var buff = bytes.NewBuffer(data)

	// Should start with an indefinite length array
	if err := ReadExpect(IndefiniteArray, buff); err != nil {
		t.Fatalf("Data does not start with an indefinite length array: %v", err)
	}

	// Read numbers until break stop code (should be 1..25)
	for c := uint64(1); ; c++ {
		n, err := ReadUInt(buff)

		if err != nil && err != FlagBreakCode {
			t.Fatal(FlagBreakCode)
		} else if err == FlagBreakCode {
			if c != 26 {
				t.Fatalf("Break stop code appeared at %d, not at %d", c, 26)
			}

			break
		} else if c > 25 {
			t.Fatalf("Reached %d, which is greater than %d", c, 25)
		} else if c != n {
			t.Fatalf("Read %d, not %d", n, c)
		}
	}
}

const (
	dataMinSize  = 0
	dataMaxSize  = 1024
	dataMinElems = 0
	dataMaxElems = 50000
)

func TestReadBigData(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		payload := rapid.SliceOfN(rapid.Byte(), dataMinSize, dataMaxSize).Draw(t, "Payload")
		elem := rapid.IntRange(dataMinElems, dataMaxElems).Draw(t, "Elements")
		buff := new(bytes.Buffer)
		bw := bufio.NewWriter(buff)

		for i := 0; i < elem; i++ {
			if err := WriteByteString(payload, bw); err != nil {
				t.Fatalf("Writing no %d errored: %v", i, err)
			}
		}

		if err := bw.Flush(); err != nil {
			t.Fatalf("Flushing errored: %v", err)
		}

		br := bufio.NewReader(buff)
		for i := 0; i < elem; i++ {
			tmp, err := ReadByteString(br)
			if err != nil {
				t.Fatalf("Reading no %d errored: %v", i, err)
			} else if len(tmp) != len(payload) {
				t.Fatalf("Length no %d mismatches: %d != %d", i, len(tmp), len(payload))
			}
			for j := 0; j < len(tmp); j++ {
				if tmp[j] != payload[j] {
					t.Fatalf("Wrong value at %d: %d != %d", i, j, len(tmp))
				}
			}
		}
	})
}

func asUntyped[T any](read func(reader io.Reader) (T, error)) func(reader io.Reader) (any, error) {
	return func(reader io.Reader) (any, error) {
		return read(reader)
	}
}

func TestNull(t *testing.T) {
	tests := []struct {
		data  []byte
		value any
		read  func(reader io.Reader) (any, error)
	}{
		{[]byte{Null, 0xfb, 0xc0, 0x10, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66}, -4.1, asUntyped(ReadFloat64)},
		{[]byte{Null, 0x64, 0x74, 0x65, 0x73, 0x74}, "test", asUntyped(ReadTextString)},
		{[]byte{Null, SimpleData | simpleTrue}, true, asUntyped(ReadBoolean)},
		{[]byte{Null, 0x0a}, uint64(10), asUntyped(ReadUInt)},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if v, err := test.read(buff); err != FlagNull {
			t.Fatalf("found Value %s", v)
		}
		if v, err := test.read(buff); err != nil {
			t.Fatal(err)
		} else if v != test.value {
			t.Fatalf("Resulting value %s not expected %s", v, test.value)
		}
	}

	buff := &bytes.Buffer{}
	err := WriteNull(buff)
	if err != nil {
		t.Fatal(err)
	} else if buff.Len() != 1 || buff.Bytes()[0] != Null {
		t.Fatal("written value not null")
	}
}
