package cboring

import (
	"bytes"
	"reflect"
	"testing"
)

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
		buff := bytes.NewBuffer(test.data)
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
