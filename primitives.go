package cboring

import (
	"fmt"
	"io"
	"math"
)

/*** Uint ***/

// ReadUInt expects an unsigned integer at the Reader's position and returns it.
func ReadUInt(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(UInt, r)
}

// WriteUInt serializes an unsigned integer into the Writer.
func WriteUInt(n uint64, w io.Writer) error {
	return WriteMajors(UInt, n, w)
}

/*** NegInt ***/

// ReadNegInt expects a negative integer at the Reader's position and returns it.
func ReadNegInt(r io.Reader) (n int64, err error) {
	un, err := ReadExpectMajors(NegInt, r)
	if err == nil {
		if un > math.MaxInt64 {
			err = fmt.Errorf("ReadNegInt: Received number is too small for int64")
		} else {
			n = -1 - int64(un)
		}
	}
	return
}

// WriteNegInt serializes a negative integer into the Writer.
func WriteNegInt(n int64, w io.Writer) error {
	if n >= 0 {
		return fmt.Errorf("WriteNegInt: Expected negative integer")
	}

	return WriteMajors(NegInt, uint64((n+1)*(-1)), w)
}

/*** ByteString ***/

// ReadByteStringLen expects a byte string at the Reader's position and returns
// the byte string's length.
func ReadByteStringLen(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(ByteString, r)
}

// WriteByteStringLen writes the type definition for a byte string with the
// given length into the Writer.
func WriteByteStringLen(n uint64, w io.Writer) error {
	return WriteMajors(ByteString, n, w)
}

/*** TextString ***/

// ReadTextStringLen expects a text string at the Reader's position and returns
// the text string's length.
func ReadTextStringLen(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(TextString, r)
}

// WriteTextStringLen writes the type definition for a text string with the
// given length into the Writer.
func WriteTextStringLen(n uint64, w io.Writer) error {
	return WriteMajors(TextString, n, w)
}

/*** Array ***/

// ReadArrayLength expects an array at the Reader's position and returns its
// length.
func ReadArrayLength(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(Array, r)
}

// WriteArrayLength writes the type definition for an array with the given
// length into the Writer.
func WriteArrayLength(n uint64, w io.Writer) error {
	return WriteMajors(Array, n, w)
}
