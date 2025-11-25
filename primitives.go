package cboring

import (
	"fmt"
	"io"
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

/*** Nint ***/

// ReadNInt expects a negative integer at the Reader's position and returns -N - 1, with N the actual number.
// when read into int64 N = int64(^n) check if value is negative
func ReadNInt(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(NInt, r)
}

// WriteNInt serializes a negative integer into the Writer. n has to be -N - 1, with N the actual number.
// when used on int64 use n = uint64(^N)
func WriteNInt(n uint64, w io.Writer) error {
	return WriteMajors(NInt, n, w)
}

/*** int ***/

// ReadInt expects either an unsigned or negative integer at the Reader's position and returns it if the value fits int64.
func ReadInt(r io.Reader) (n int64, err error) {
	major, num, err := ReadMajors(r)
	n = int64(num)
	if n < 0 {
		if major == UInt {
			err = fmt.Errorf("ReadInt: Returned Integer %d to big for int64", num)
		} else if major == NInt {
			err = fmt.Errorf("ReadInt: Returned Integer -%d to small for int64",
				num+1) // might overflow but highly unlikely
		} else {
			err = fmt.Errorf("ReadInt: Wrong Major Type: 0x%x instead of 0x00 or 0x20",
				major)
		}
	} else if major == NInt {
		n = ^n
	} else if major != UInt {
		err = fmt.Errorf("ReadInt: Wrong Major Type: 0x%x instead of 0x00 or 0x20",
			major)
	}
	return
}

// WriteInt serializes an integer into the Writer, either as UInt or NInt.
func WriteInt(n int64, w io.Writer) error {
	if n < 0 {
		return WriteNInt(uint64(^n), w)
	}
	return WriteUInt(uint64(n), w)
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

/*** Map ***/

// ReadMapPairLength expects a map at the Reader's position and returns the
// amount of pairs stored.
func ReadMapPairLength(r io.Reader) (n uint64, err error) {
	return ReadExpectMajors(Map, r)
}

// WriteMapPairLength writes the type definition for a map with the given
// amount of pairs into the Writer.
func WriteMapPairLength(n uint64, w io.Writer) error {
	return WriteMajors(Map, n, w)
}
