package endpoint

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/dtn7/cboring"
)

var endpointTests = []struct {
	eid  string
	cbor []byte
}{
	{"dtn:none", []byte{0x82, 0x01, 0x00}},
	{"dtn:foo", []byte{0x82, 0x01, 0x63, 0x66, 0x6F, 0x6F}},
	{"dtn:foo/bar", []byte{0x82, 0x01, 0x67, 0x66, 0x6F, 0x6F, 0x2F, 0x62, 0x61, 0x72}},
	{"ipn:0.0", []byte{0x82, 0x02, 0x82, 0x00, 0x00}},
	{"ipn:23.42", []byte{0x82, 0x02, 0x82, 0x17, 0x18, 0x2A}},
}

func TestEndpoint(t *testing.T) {
	for _, test := range endpointTests {
		t.Run(fmt.Sprintf("marshal-%s", test.eid), func(t *testing.T) {
			e := newEndpointID(test.eid)

			buff := new(bytes.Buffer)
			if err := cboring.Marshal(&e, buff); err != nil {
				t.Fatalf("Marshaling %s failed: %v", test.eid, err)
			}

			if data := buff.Bytes(); !reflect.DeepEqual(data, test.cbor) {
				t.Fatalf("CBOR differs: %x != %x", data, test.cbor)
			}
		})

		t.Run(fmt.Sprintf("unmarshal-%s", test.eid), func(t *testing.T) {
			e := endpointID{}

			buff := bytes.NewBuffer(test.cbor)
			if err := cboring.Unmarshal(&e, buff); err != nil {
				t.Fatalf("Unmarshaling %s failed: %v", test.eid, err)
			}

			if e.String() != test.eid {
				t.Fatalf("EID differs: %s != %s", e.String(), test.eid)
			}
		})
	}
}

func BenchmarkEndpoint(b *testing.B) {
	for _, test := range endpointTests {
		b.Run(fmt.Sprintf("marshal-%s", test.eid), func(b *testing.B) {
			e := newEndpointID(test.eid)

			buff := new(bytes.Buffer)
			if err := cboring.Marshal(&e, buff); err != nil {
				b.Fatalf("Marshaling %s failed: %v", test.eid, err)
			}
		})

		b.Run(fmt.Sprintf("unmarshal-%s", test.eid), func(b *testing.B) {
			e := endpointID{}

			buff := bytes.NewBuffer(test.cbor)
			if err := cboring.Unmarshal(&e, buff); err != nil {
				b.Fatalf("Unmarshaling %s failed: %v", test.eid, err)
			}
		})
	}
}
