package payloadblock

import (
	"bytes"
	"fmt"
	"math/rand"
	"reflect"
	"testing"

	"github.com/dtn7/cboring"
)

func TestPayload(t *testing.T) {
	var pb = newPayloadBlock([]byte("hello"))
	var pbData = []byte{
		0x85, 0x01, 0x01, 0x02, 0x00, 0x45, 0x68, 0x65, 0x6C, 0x6C, 0x6F}

	t.Run("marshal", func(t *testing.T) {
		buff := new(bytes.Buffer)
		if err := cboring.Marshal(&pb, buff); err != nil {
			t.Fatalf("Marshaling failed: %v", err)
		}

		if data := buff.Bytes(); !reflect.DeepEqual(data, pbData) {
			t.Fatalf("CBOR differs: %x != %x", data, pbData)
		}
	})

	t.Run("unmarshal", func(t *testing.T) {
		pbTmp := payloadBlock{}
		buff := bytes.NewBuffer(pbData)

		if err := cboring.Unmarshal(&pbTmp, buff); err != nil {
			t.Fatalf("Unmarshaling failed: %v", err)
		}

		if !reflect.DeepEqual(pbTmp, pb) {
			t.Fatalf("PayloadBlock differs: %v != %v", pbTmp, pb)
		}
	})
}

func BenchmarkPayload(b *testing.B) {
	sizes := []int{
		// Ridiculously small
		0, 1, 128, 256,
		// Kibibytes
		1024, 10240, 102400,
		// Mebibytes
		1048576, 10485760, 104857600,
	}

	for _, size := range sizes {
		rndData := make([]byte, size)
		rand.Seed(0)
		rand.Read(rndData)

		b.Run(fmt.Sprintf("marshal-%d", size), func(b *testing.B) {
			// Setup a buffer, like in the unmarshaling test
			pbTmp := newPayloadBlock(rndData)

			buff := new(bytes.Buffer)
			if err := cboring.Marshal(&pbTmp, buff); err != nil {
				b.Fatalf("Marshaling failed: %v", err)
			}

			// Benchmark starts here
			b.ResetTimer()

			pb := payloadBlock{}
			if err := cboring.Unmarshal(&pb, buff); err != nil {
				b.Fatalf("Unmarshaling failed: %v", err)
			}
		})

		b.Run(fmt.Sprintf("unmarshal-%d", size), func(b *testing.B) {
			pb := newPayloadBlock(rndData)

			buff := new(bytes.Buffer)
			if err := cboring.Marshal(&pb, buff); err != nil {
				b.Fatalf("Marshaling failed: %v", err)
			}
		})
	}
}
