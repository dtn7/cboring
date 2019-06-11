package cboring

import "io"

/*** Uint ***/

// ReadUInt expects an unsigned integer at the Reader's position and returns it.
func ReadUInt(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(UInt, r)
}

// WriteUInt serializes an unsigned integer into the Writer.
func WriteUInt(n uint64, w io.Writer) error {
	return WriteMajors(UInt, n, w)
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
