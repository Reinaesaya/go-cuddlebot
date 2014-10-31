package msgtype

import (
	"bufio"
	"encoding/binary"
	"errors"
	"hash/crc32"
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

/* Setpoint value. */
type Setpoint struct {
	Duration uint16 // offset 0x00, duration in ms
	Setpoint uint16 // offset 0x02, setpoint
}

/* RPC wrapper. */
type RPC struct {
	io.Reader
	io.Writer
}

/* Invalid message error. */
var InvalidMessageError = errors.New("Invalid message")

/* Write ping message. */
func (w *RPC) Ping(addr uint8) error {
	return w.writeSimpleAddressedMessage(addr, kPing)
}

/* Write set PID message. */
func (w *RPC) SetPID(addr uint8, kp, ki, kd float32) error {
	h := crc32.NewIEEE()
	bw := bufio.NewWriterSize(w, 18)
	mw := io.MultiWriter(bw, h)
	// write header
	header := []uint8{addr, kSetPID, 12, 0}
	if err := binary.Write(mw, binary.LittleEndian, header); err != nil {
		return err
	}
	// write data
	data := []float32{kp, ki, kd}
	if err := binary.Write(mw, binary.LittleEndian, data); err != nil {
		return err
	}
	// write checksum
	var checksum uint16 = uint16(h.Sum32())
	if err := binary.Write(bw, binary.LittleEndian, checksum); err != nil {
		return err
	}
	// flush buffer
	return bw.Flush()
}

/* Write set Setpoint message. */
func (w *RPC) Setpoint(
	addr uint8, delay, loop uint16, setpoints []Setpoint) error {
	nsetpoints := len(setpoints)

	if nsetpoints <= 0 || nsetpoints > 0xffff {
		return InvalidMessageError
	}

	h := crc32.NewIEEE()
	bw := bufio.NewWriter(w)
	mw := io.MultiWriter(bw, h)

	// write header
	header := []uint8{addr, kSetpoint}
	if err := binary.Write(mw, binary.LittleEndian, header); err != nil {
		return err
	}
	// write size
	size := uint16(6 + 4*nsetpoints)
	if err := binary.Write(mw, binary.LittleEndian, size); err != nil {
		return err
	}
	// write data
	data := []uint16{delay, loop, uint16(nsetpoints)}
	if err := binary.Write(mw, binary.LittleEndian, data); err != nil {
		return err
	}
	if err := binary.Write(mw, binary.LittleEndian, setpoints); err != nil {
		return err
	}
	// write checksum
	var checksum uint16 = uint16(h.Sum32())
	if err := binary.Write(bw, binary.LittleEndian, checksum); err != nil {
		return err
	}
	// flush buffer
	return bw.Flush()
}

/* Write set test message. */
func (w *RPC) RunTests(addr uint8) error {
	return w.writeSimpleAddressedMessage(addr, kTest)
}

/* Write set value message. */
func (w *RPC) RequestPosition(addr uint8) error {
	return w.writeSimpleAddressedMessage(addr, kValue)
}

/* Write simple message. */
func (w *RPC) writeSimpleAddressedMessage(addr, msgtype uint8) error {
	h := crc32.NewIEEE()
	bw := bufio.NewWriterSize(w, 6)
	mw := io.MultiWriter(bw, h)
	// write header
	header := []uint8{addr, msgtype, 0, 0}
	if err := binary.Write(mw, binary.LittleEndian, header); err != nil {
		return err
	}
	// write checksum
	var checksum uint16 = uint16(h.Sum32())
	if err := binary.Write(bw, binary.LittleEndian, checksum); err != nil {
		return err
	}
	// flush buffer
	return bw.Flush()
}
