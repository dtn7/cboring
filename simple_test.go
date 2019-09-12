package cboring

import (
	"bytes"
	"reflect"
	"testing"
)

func TestBoolean(t *testing.T) {
	tests := []struct {
		data []byte
		b    bool
	}{
		{[]byte{0xF4}, false},
		{[]byte{0xF5}, true},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if b, err := ReadBoolean(buff); err != nil {
			t.Fatal(err)
		} else if b != test.b {
			t.Fatalf("Resulting bool %t is not %t", b, test.b)
		}

		// Write
		buff.Reset()
		if err := WriteBoolean(test.b, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestFloat32(t *testing.T) {
	tests := []struct {
		data []byte
		f    float32
	}{
		{[]byte{0xfa, 0x47, 0xc3, 0x50, 0x00}, 100000.0},
		{[]byte{0xfa, 0x7f, 0x7f, 0xff, 0xff}, 3.4028234663852886e+38},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if f, err := ReadFloat32(buff); err != nil {
			t.Fatal(err)
		} else if f != test.f {
			t.Fatalf("Resulting float %f is not %f", f, test.f)
		}

		// Write
		buff.Reset()
		if err := WriteFloat32(test.f, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		data []byte
		f    float64
	}{
		{[]byte{0xfb, 0x3f, 0xf1, 0x99, 0x99, 0x99, 0x99, 0x99, 0x9a}, 1.1},
		{[]byte{0xfb, 0x7e, 0x37, 0xe4, 0x3c, 0x88, 0x00, 0x75, 0x9c}, 1.0e+300},
		{[]byte{0xfb, 0xc0, 0x10, 0x66, 0x66, 0x66, 0x66, 0x66, 0x66}, -4.1},
	}

	for _, test := range tests {
		// Read
		buff := bytes.NewBuffer(test.data)
		if f, err := ReadFloat64(buff); err != nil {
			t.Fatal(err)
		} else if f != test.f {
			t.Fatalf("Resulting float %f is not %f", f, test.f)
		}

		// Write
		buff.Reset()
		if err := WriteFloat64(test.f, buff); err != nil {
			t.Fatal(err)
		}

		if bb := buff.Bytes(); !reflect.DeepEqual(bb, test.data) {
			t.Fatalf("Serialized data mismatches: %x != %x", bb, test.data)
		}
	}
}
