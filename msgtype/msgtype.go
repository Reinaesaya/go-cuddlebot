package msgtype

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/mikepb/go-crc16"
)

const (

	// Board addresses.
	InvalidAddress RemoteAddress = 0   // invalid address
	RibsAddress                  = 'r' // ribs actuator
	PurrAddress                  = 'p' // purr motor
	SpineAddress                 = 's' // spine actuator
	HeadXAddress                 = 'x' // head yaw actuator
	HeadYAddress                 = 'y' // head pitch actuator

	// Board addresses as strings.
	RibsAddressString  = "ribs"  // ribs actuator
	PurrAddressString  = "purr"  // purr motor
	SpineAddressString = "spine" // spine actuator
	HeadXAddressString = "headx" // head yaw actuator
	HeadYAddressString = "heady" // head pitch actuator

	// Message types.
	kInvalidType uint8 = 0   // invalid message
	kPing              = '?' // ping an actuator
	kPong              = '.' // respond to ping
	kSetPID            = 'c' // send PID coefficients
	kSetpoint          = 'g' // send setpoints
	kSleep             = 'z' // deactivate motor output
	kTest              = 't' // run internal tests
	kValue             = 'v' // get position value

)

// Loop setpoints forever.
const LOOP_INFINITE uint16 = 0xffff

type RemoteAddress uint8

// Data for simple message types.
type simpleType struct {
	Addr RemoteAddress `json:"addr"`
}

// Ping message type.
type Ping simpleType

// SetPID message type.
type SetPID struct {
	Addr RemoteAddress `json:"addr"`
	Kp   float32       `json:"kp"`
	Ki   float32       `json:"ki"`
	Kd   float32       `json:"kd"`
}

// Setpoint message type.
type Setpoint struct {
	Addr      RemoteAddress   `json:"addr"`
	Delay     uint16          `json:"delay"`
	Loop      uint16          `json:"loop"`
	Setpoints []SetpointValue `json:"setpoints"`
}

// Setpoint value.
type SetpointValue struct {
	Duration uint16 `json:"duration"` // offset 0x00, duration in ms
	Setpoint uint16 `json:"setpoint"` // offset 0x02, setpoint
}

// Sleep message type.
type Sleep simpleType

// Test message type.
type Test simpleType

// Value message type.
type Value simpleType

// Invalid message error.
var InvalidAddressError = errors.New("Invalid address")

// Invalid message error.
var InvalidMessageError = errors.New("Invalid message")

// Serialize address to text.
func (a *RemoteAddress) MarshalText() ([]byte, error) {
	var s string

	switch *a {
	case RibsAddress:
		s = RibsAddressString
	case PurrAddress:
		s = PurrAddressString
	case SpineAddress:
		s = SpineAddressString
	case HeadXAddress:
		s = HeadXAddressString
	case HeadYAddress:
		s = HeadYAddressString
	default:
		return nil, InvalidAddressError
	}

	return bytes.NewBufferString(s).Bytes(), nil
}

// Deserialize address from text.
func (a *RemoteAddress) UnmarshalText(text []byte) error {
	s, err := bytes.NewBuffer(text).ReadString(0)
	if err != io.EOF {
		return err
	}

	switch s {
	case RibsAddressString:
		*a = RibsAddress
	case PurrAddressString:
		*a = PurrAddress
	case SpineAddressString:
		*a = SpineAddress
	case HeadXAddressString:
		*a = HeadXAddress
	case HeadYAddressString:
		*a = HeadYAddress
	default:
		return InvalidAddressError
	}

	return nil
}

// Encode ping message.
func (m *Ping) MarshalBinary() (data []byte, err error) {
	return marshalBinary(m.Addr, kPing)
}

// Encode PID message.
func (m *SetPID) MarshalBinary() (data []byte, err error) {
	var b bytes.Buffer
	h := crc16.NewANSI()
	ww := io.MultiWriter(&b, h)
	// write header
	if _, err = ww.Write([]byte{uint8(m.Addr), kSetPID, 12, 0}); err != nil {
		return
	}
	// write data
	d := []float32{m.Kp, m.Ki, m.Kd}
	if err = binary.Write(ww, binary.LittleEndian, d); err != nil {
		return
	}
	// write checksum
	sum := h.Sum16()
	if err = binary.Write(ww, binary.LittleEndian, sum); err != nil {
		return
	}
	// return data
	return b.Bytes(), nil
}

// Write set Setpoint message.
func (m *Setpoint) MarshalBinary() (data []byte, err error) {
	nsetpoints := len(m.Setpoints)

	// (1024-6)/4 = 254 max setpoints for 1024 byte max data
	if nsetpoints <= 0 || nsetpoints > 254 {
		return nil, InvalidMessageError
	}

	var b bytes.Buffer
	h := crc16.NewANSI()
	ww := io.MultiWriter(&b, h)

	// write header
	if _, err = ww.Write([]uint8{uint8(m.Addr), kSetpoint}); err != nil {
		return
	}
	// write size
	size := uint16(6 + 4*nsetpoints)
	if err = binary.Write(ww, binary.LittleEndian, size); err != nil {
		return
	}
	// write data
	d := []uint16{m.Delay, m.Loop, uint16(nsetpoints)}
	if err = binary.Write(ww, binary.LittleEndian, d); err != nil {
		return
	}
	// write setpoint data
	if err = binary.Write(ww, binary.LittleEndian, m.Setpoints); err != nil {
		return
	}
	sum := h.Sum16()
	if err = binary.Write(ww, binary.LittleEndian, sum); err != nil {
		return
	}
	// return data
	return b.Bytes(), nil
}

// Write sleep test message.
func (m *Sleep) MarshalBinary() (data []byte, err error) {
	return marshalBinary(m.Addr, kSleep)
}

// Write set test message.
func (m *Test) MarshalBinary() (data []byte, err error) {
	return marshalBinary(m.Addr, kTest)
}

// Write request value message.
func (m *Value) MarshalBinary() (data []byte, err error) {
	return marshalBinary(m.Addr, kValue)
}

// Marshal simple message.
func marshalBinary(addr RemoteAddress, msgtype uint8) (data []byte, err error) {
	var b bytes.Buffer
	h := crc16.NewANSI()
	ww := io.MultiWriter(&b, h)
	// write header
	if _, err = ww.Write([]uint8{uint8(addr), msgtype, 0, 0}); err != nil {
		return
	}
	sum := h.Sum16()
	if err = binary.Write(ww, binary.LittleEndian, sum); err != nil {
		return
	}
	// return contents
	return b.Bytes(), nil
}
