package msgtype

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
)

/* Board addresses. */
const (
	ADDR_INVALID    uint8 = 0   // invalid address
	ADDR_RIBS             = 'r' // ribs actuator
	ADDR_PURR             = 'p' // purr motor
	ADDR_SPINE            = 's' // spine actuator
	ADDR_HEAD_YAW         = 'x' // head yaw actuator
	ADDR_HEAD_PITCH       = 'y' // head pitch actuator
)

/* Message types. */
const (
	kInvalidType uint8 = 0   // invalid message
	kPing              = '?' // ping an actuator
	kPong              = '.' // respond to ping
	kSetPID            = 'c' // send PID coefficients
	kSetpoint          = 'g' // send setpoints
	kTest              = 't' // run internal tests
	kValue             = 'v' // get position value
)

/* Loop setpoints forever. */
const LOOP_INFINITE uint16 = 0xffff

/* Data for simple message types. */
type simpleType struct {
	Addr uint8
}

/* Ping message type. */
type Ping simpleType

/* SetPID message type. */
type SetPID struct {
	Addr uint8
	Kp   float32
	Ki   float32
	Kd   float32
}

/* Setpoint message type. */
type Setpoint struct {
	Addr      uint8
	Delay     uint16
	Loop      uint16
	Setpoints []SetpointValue
}

/* Setpoint value. */
type SetpointValue struct {
	Duration uint16 // offset 0x00, duration in ms
	Setpoint uint16 // offset 0x02, setpoint
}

/* Test message type. */
type Test simpleType

/* Value message type. */
type Value simpleType

/* Invalid message error. */
var InvalidMessageError = errors.New("Invalid message")

/* Write ping message. */
func (m *Ping) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, m.Addr, kPing)
}

/* Write set PID message. */
func (m *SetPID) WriteTo(w io.Writer) (int64, error) {
	h := NewModbus()
	b := bufio.NewWriter(w)
	ww := io.MultiWriter(b, h)
	// write header
	if _, err := ww.Write([]byte{m.Addr, kSetPID, 12, 0}); err != nil {
		return 0, err
	}
	// write data
	data := []float32{m.Kp, m.Ki, m.Kd}
	if err := binary.Write(ww, binary.LittleEndian, data); err != nil {
		return 0, err
	}
	// write checksum
	if _, err := ww.Write([]uint8{h.lo, h.hi}); err != nil {
		return 0, err
	}
	// flush buffer
	n := b.Buffered()
	if err := b.Flush(); err != nil {
		return 0, err
	}
	return int64(n), nil
}

/* Write set Setpoint message. */
func (m *Setpoint) WriteTo(w io.Writer) (int64, error) {
	nsetpoints := len(m.Setpoints)

	// (1024-6)/4 = 254 max setpoints for 1024 byte max data
	if nsetpoints <= 0 || nsetpoints > 254 {
		return 0, InvalidMessageError
	}

	h := NewModbus()
	b := bufio.NewWriter(w)
	ww := io.MultiWriter(b, h)

	// write header
	if _, err := ww.Write([]uint8{m.Addr, kSetpoint}); err != nil {
		return 0, err
	}
	// write size
	size := uint16(6 + 4*nsetpoints)
	if err := binary.Write(ww, binary.LittleEndian, size); err != nil {
		return 0, err
	}
	// write data
	data := []uint16{m.Delay, m.Loop, uint16(nsetpoints)}
	if err := binary.Write(ww, binary.LittleEndian, data); err != nil {
		return 0, err
	}
	// write setpoint data
	if err := binary.Write(ww, binary.LittleEndian, m.Setpoints); err != nil {
		return 0, err
	}
	// write checksum
	if _, err := ww.Write([]uint8{h.lo, h.hi}); err != nil {
		return 0, err
	}
	// flush buffer
	n := b.Buffered()
	if err := b.Flush(); err != nil {
		return 0, err
	}
	return int64(n), nil
}

/* Write set test message. */
func (m *Test) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, m.Addr, kTest)
}

/* Write request value message. */
func (m *Value) WriteTo(w io.Writer) (int64, error) {
	return writeTo(w, m.Addr, kValue)
}

/* Write simple message. */
func writeTo(w io.Writer, addr, msgtype uint8) (int64, error) {
	h := NewModbus()
	b := bufio.NewWriter(w)
	ww := io.MultiWriter(b, h)
	// write header
	if _, err := ww.Write([]uint8{addr, msgtype, 0, 0}); err != nil {
		return 0, err
	}
	// write checksum
	if _, err := ww.Write([]uint8{h.lo, h.hi}); err != nil {
		return 0, err
	}
	// flush buffer
	n := b.Buffered()
	if err := b.Flush(); err != nil {
		return 0, err
	}
	return int64(n), nil
}
