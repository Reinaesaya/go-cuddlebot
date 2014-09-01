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

func Read(m interface{}, r io.Reader) error {
	// read into struct
	if err := binary.Read(r, binary.LittleEndian, m); err != nil {
		return err
	}
	// validate if validated message
	return validate(m)
}

func Write(m interface{}, w io.Writer) error {
	// use marshaler to calculate checksum
	switch m.(type) {
	case encoding.BinaryMarshaler:
		if buf, err := m.(encoding.BinaryMarshaler).MarshalBinary(); err != nil {
			return err
		} else {
			_, err = w.Write(buf)
			return err
		}
	}
	// should not reach here
	return InvalidMessageError
}

func validate(m interface{}) error {
	switch m.(type) {
	case ValidatedMessage:
		return m.(ValidatedMessage).Validate()
	}
	return nil
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
