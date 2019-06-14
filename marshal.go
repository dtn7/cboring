package cboring

import "io"

// CborMarshaler is the interface implemented by an object that can both marshal
// itself into a CBOR form and unmarshal a CBOR representation of itself.
type CborMarshaler interface {
	MarshalCbor(w io.Writer) error
	UnmarshalCbor(r io.Reader) error
}

// Marshal writes a CBOR representation of a CborMarshaler into the Writer.
func Marshal(data CborMarshaler, w io.Writer) error {
	return data.MarshalCbor(w)
}

// Unmarshal reads a CBOR representation from a Reader into a CborMarshaler.
func Unmarshal(data CborMarshaler, r io.Reader) error {
	return data.UnmarshalCbor(r)
}
