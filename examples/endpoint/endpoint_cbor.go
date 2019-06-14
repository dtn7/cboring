package main

import (
	"io"

	"github.com/dtn7/cboring"
)

func (eid *EndpointID) MarshalCbor(w io.Writer) error {
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

func (eid *EndpointID) UnmarshalCbor(r io.Reader) error {
	// TODO
	return nil
}
