package payloadblock

import (
	"fmt"
	"io"

	"github.com/dtn7/cboring"
)

type payloadBlock struct {
	BlockType         uint64
	BlockNumber       uint64
	BlockControlFlags uint64
	CRCType           uint64
	Data              []byte
}

func newPayloadBlock(data []byte) payloadBlock {
	return payloadBlock{
		BlockType:         1,
		BlockNumber:       1,
		BlockControlFlags: 0x02,
		CRCType:           0,
		Data:              data}
}

func (pb *payloadBlock) MarshalCbor(w io.Writer) error {
	// Start an array with five elements
	if err := cboring.WriteArrayLength(5, w); err != nil {
		return err
	}

	// Write the four fields
	fields := []uint64{pb.BlockType, pb.BlockNumber, pb.BlockControlFlags, pb.CRCType}
	for _, f := range fields {
		if err := cboring.WriteUInt(f, w); err != nil {
			return err
		}
	}

	// Write the data blob
	if err := cboring.WriteByteString(pb.Data, w); err != nil {
		return err
	}

	return nil
}

func (pb *payloadBlock) UnmarshalCbor(r io.Reader) error {
	// Start of an array with five elements
	if l, err := cboring.ReadArrayLength(r); err != nil {
		return err
	} else if l != 5 {
		return fmt.Errorf("Expected array with length 5, got %d", l)
	}

	fields := []*uint64{&pb.BlockType, &pb.BlockNumber, &pb.BlockControlFlags, &pb.CRCType}
	for _, f := range fields {
		if n, err := cboring.ReadUInt(r); err != nil {
			return err
		} else {
			*f = n
		}
	}

	if data, err := cboring.ReadByteString(r); err != nil {
		return err
	} else {
		pb.Data = data
	}

	return nil
}
