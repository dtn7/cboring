package cboring

import "io"

type CborMarshaler interface {
	MarshalCbor(w io.Writer) error
	UnmarshalCbor(r io.Reader) error
}

/*
func Marshal(data CborMarshaler, w io.Writer) error {

}

func Unmarshal(data CborMarshaler, r io.Reader) error {
}
*/
