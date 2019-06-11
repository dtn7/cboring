package cboring

import (
	"bytes"
	"reflect"
	"testing"
)

func TestReadMajorsSmall(t *testing.T) {
	tests := []MajorType{UInt, NegInt, ByteString, TextString, Array}

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
		{[]byte{0x38, 0x63}, NegInt, 99}, // 99 = abs(-1 - 100)
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
		[]byte{},
		// Incomplete streams
		[]byte{0x18}, []byte{0x19, 0x03},
	}

	for _, test := range tests {
		r := bytes.NewBuffer(test)
		if _, _, err := ReadMajors(r); err == nil {
			t.Fatalf("Illegal input %x did not errored", test)
		}
	}
}

func TestWriteMajorsSmall(t *testing.T) {
	tests := []MajorType{UInt, NegInt, ByteString, TextString, Array}

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
		{[]byte{0x38, 0x63}, NegInt, 99}, // 99 = abs(-1 - 100)
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
