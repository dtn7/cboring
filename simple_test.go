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
