package endpoint

import (
	"fmt"
	"io"

	"github.com/dtn7/cboring"
)

func (eid *endpointID) MarshalCbor(w io.Writer) error {
	// Start an array with two elements
	if err := cboring.WriteArrayLength(2, w); err != nil {
		return err
	}

	// URI code: scheme name
	if err := cboring.WriteUInt(eid.SchemeName, w); err != nil {
		return err
	}

	// SSP
	switch eid.SchemeSpecificPart.(type) {
	case uint64:
		// dtn:none
		if err := cboring.WriteUInt(0, w); err != nil {
			return err
		}

	case string:
		// dtn:whatsoever
		if err := cboring.WriteTextString(eid.SchemeSpecificPart.(string), w); err != nil {
			return err
		}

	case [2]uint64:
		// ipn:23.42
		var ssps [2]uint64 = eid.SchemeSpecificPart.([2]uint64)
		if err := cboring.WriteArrayLength(2, w); err != nil {
			return err
		}

		for _, ssp := range ssps {
			if err := cboring.WriteUInt(ssp, w); err != nil {
				return err
			}
		}
	}

	return nil
}

func (eid *endpointID) UnmarshalCbor(r io.Reader) error {
	// Start of an array with two elements
	if l, err := cboring.ReadArrayLength(r); err != nil {
		return err
	} else if l != 2 {
		return fmt.Errorf("Expected array with length 2, got %d", l)
	}

	// URI code: scheme name
	if sn, err := cboring.ReadUInt(r); err != nil {
		return err
	} else {
		eid.SchemeName = sn
	}

	// SSP
	if m, n, err := cboring.ReadMajors(r); err != nil {
		return err
	} else {
		switch m {
		case cboring.UInt:
			// dtn:none
			eid.SchemeSpecificPart = n

		case cboring.TextString:
			// dtn:whatsoever
			if tmp, err := cboring.ReadRawBytes(n, r); err != nil {
				return err
			} else {
				eid.SchemeSpecificPart = string(tmp)
			}

		case cboring.Array:
			// ipn:23.42
			if n != 2 {
				return fmt.Errorf("Expected array with length 2, got %d", n)
			}

			var ssps [2]uint64
			for i := 0; i < 2; i++ {
				if n, err := cboring.ReadUInt(r); err != nil {
					return err
				} else {
					ssps[i] = n
				}
			}

			eid.SchemeSpecificPart = ssps
		}
	}

	return nil
}
