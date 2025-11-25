package cboring

import (
	"fmt"
	"io"
	"math"
)

const (
	simpleFalse byte = 20
	simpleTrue  byte = 21
	simpleNull  byte = 22
)

// ReadBoolean reads a bool value from the Reader.
func ReadBoolean(r io.Reader) (b bool, err error) {
	var buff [1]byte

	if _, dataErr := r.Read(buff[:1]); dataErr != nil {
		err = dataErr
		return
	}

	major, adds := readMajorType(buff[0])
	if major != SimpleData {
		err = fmt.Errorf("ReadBoolean: Expected major 0x%x, got 0x%x", SimpleData, major)
		return
	}

	switch adds {
	case simpleFalse:
		b = false
	case simpleTrue:
		b = true
	case simpleNull:
		err = FlagNull
	default:
		err = fmt.Errorf("ReadBoolean: Unknown additional 0x%x", adds)
	}

	return
}

// WriteBoolean writes a bool into the Writer.
func WriteBoolean(b bool, w io.Writer) (err error) {
	var adds byte
	if b {
		adds = simpleTrue
	} else {
		adds = simpleFalse
	}

	payload := writeMajorType(SimpleData, adds)
	_, err = w.Write([]byte{payload})

	return
}

// ReadFloat32 reads a float32 value from the Reader.
func ReadFloat32(r io.Reader) (f float32, err error) {
	if fbits, fbitsErr := ReadExpectMajors(SimpleData, r); fbitsErr != nil {
		err = fbitsErr
	} else {
		f = math.Float32frombits(uint32(fbits))
	}

	return
}

// WriteFloat32 writes a float32 into the Writer.
func WriteFloat32(f float32, w io.Writer) (err error) {
	fbits := math.Float32bits(f)
	return WriteMajors(SimpleData, uint64(fbits), w)
}

// ReadFloat64 reads a float64 value from the Reader.
func ReadFloat64(r io.Reader) (f float64, err error) {
	if fbits, fbitsErr := ReadExpectMajors(SimpleData, r); fbitsErr != nil {
		err = fbitsErr
	} else {
		f = math.Float64frombits(fbits)
	}

	return
}

// WriteFloat64 writes a float64 into the Writer.
func WriteFloat64(f float64, w io.Writer) (err error) {
	fbits := math.Float64bits(f)
	return WriteMajors(SimpleData, fbits, w)
}

// WriteNull writes a null into the Writer.
func WriteNull(w io.Writer) (err error) {
	return WriteMajors(SimpleData, uint64(simpleNull), w)
}
