package starfish

import (
	"bytes"
	"encoding"
	"encoding/binary"
	"errors"
	"io"
)

type UpdatableMessage interface {
	Update() error
}

type ValidatedMessage interface {
	Validate() error
}

var InsufficientBytesError = errors.New("Insufficient bytes")
var InvalidMessageError = errors.New("Invalid message")
var InvalidAddressError = errors.New("Invalid address")

// Read into struct using BinaryUnmarshaler interface to invoke
// validation.
func ReadFrom(m interface{}, r io.Reader) (int64, error) {
	size := binary.Size(m)
	switch m.(type) {
	case encoding.BinaryUnmarshaler:
		// read into buffer
		buf := make([]byte, 0, size)
		if _, err := io.ReadFull(r, buf); err != nil {
			return 0, err
		}
		// use marshaler to validate checksum
		return int64(size), m.(encoding.BinaryUnmarshaler).UnmarshalBinary(buf)
	default:
		return int64(size), binary.Read(r, binary.LittleEndian, m)
	}
}

// Write struct using BinaryMarshaler interface to calculate checksum.
func WriteTo(m interface{}, w io.Writer) (int64, error) {
	size := binary.Size(m)
	switch m.(type) {
	case encoding.BinaryMarshaler:
		// use marshaler to calculate checksum
		buf, err := m.(encoding.BinaryMarshaler).MarshalBinary()
		if err != nil {
			return 0, err
		}
		// write buffer data
		n, err := w.Write(buf)
		return int64(n), err
	default:
		// support other types by default
		return int64(size), binary.Write(w, binary.LittleEndian, m)
	}
}

func unmartialBinary(m interface{}, data []byte) error {
	size := binary.Size(m)

	// check for expected number of bytes
	if len(data) < size {
		return InsufficientBytesError
	}
	// validate checksum
	if err := validateChecksum(data, size); err != nil {
		return err
	}
	// unmarshal data
	if err := binary.Read(bytes.NewBuffer(data), binary.LittleEndian, m); err != nil {
		return err
	}
	// validate if validated message
	return validate(m)
}

func martialBinary(m interface{}) ([]byte, error) {
	size := binary.Size(m)

	// update if updatable message
	switch m.(type) {
	case UpdatableMessage:
		if err := m.(UpdatableMessage).Update(); err != nil {
			return nil, err
		}
	}

	// convert struct to byte array
	w := bytes.NewBuffer(make([]byte, 0, size))
	if err := binary.Write(w, binary.LittleEndian, m); err != nil {
		return nil, err
	}

	// remove checksum
	w.Truncate(size - 4)
	// calculate and write checksum
	if err := writeChecksum(w); err != nil {
		return nil, err
	}

	return w.Bytes(), nil
}

func validate(m interface{}) error {
	switch m.(type) {
	case ValidatedMessage:
		return m.(ValidatedMessage).Validate()
	}
	return nil
}

func validateChecksum(b []byte, size int) error {
	var csum int32
	wc := bytes.NewBuffer(b)
	nints := int(size/4) - 1
	var i32 int32
	for i := 0; i < nints; i++ {
		if err := binary.Read(wc, binary.LittleEndian, &i32); err != nil {
			return err
		}
		csum += i32
	}
	if err := binary.Read(wc, binary.LittleEndian, &i32); err != nil {
		return err
	}
	if csum != i32 {
		return InvalidMessageError
	}
	return nil
}

func writeChecksum(b *bytes.Buffer) error {
	var csum int32
	wc := bytes.NewBuffer(b.Bytes())
	nints := int(b.Len() / 4)
	for i := 0; i < nints; i++ {
		var i32 int32
		if err := binary.Read(wc, binary.LittleEndian, &i32); err != nil {
			return err
		}
		csum += i32
	}
	return binary.Write(b, binary.LittleEndian, csum)
}
