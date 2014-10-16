package msgtype

import (
	"encoding/binary"
	"hash/crc32"
	"io"
)

/* Board addresses. */
const (
	ADDR_INVALID    = iota // invalid address
	ADDR_RIBS              // ribs actuator
	ADDR_PURR              // purr motor
	ADDR_SPINE             // spine actuator
	ADDR_HEAD_YAW          // head yaw actuator
	ADDR_HEAD_PITCH        // head pitch actuator
)

/* Message types. */
const (
	kInvalidType = iota // invalid message
	kPing               // ping an actuator
	kPong               // respond to ping
	kSetPID             // send PID coefficients
	kSetpoint           // send setpoints
	kTest               // run internal tests
	kValue              // get position value
)

/* Setpoint value. */
type Setpoint struct {
	Duration uint16 // offset 0x00, duration in ms
	Setpoint uint16 // offset 0x02, setpoint
}

/* RPC wrapper. */
type RPC struct {
	io.Writer
}

/* Write ping message. */
func (w *RPC) Ping(addr uint8) error {
	return binary.Write(w, binary.LittleEndian,
		[]byte{addr, kPing, 0, 0, 0, 0})
}

/* Write pong message. */
func (w *RPC) Pong(addr uint8) error {
	return binary.Write(w, binary.LittleEndian,
		[]byte{addr, kPong, 0, 0, 0, 0})
}

/* Write set PID message. */
func (w *RPC) SetPID(addr uint8, kp, ki, kd float32) error {
	// write header
	header := []uint8{addr, kSetPID, 12, 0}
	if err := binary.Write(w, binary.LittleEndian, header); err != nil {
		return err
	}
	// write data
	h := crc32.NewIEEE()
	mw := io.MultiWriter(w, h)
	data := []float32{kp, ki, kd}
	if err := binary.Write(mw, binary.LittleEndian, data); err != nil {
		return err
	}
	// write checksum
	var checksum uint16 = uint16(h.Sum32())
	return binary.Write(w, binary.LittleEndian, checksum)
}

/* Write set Setpoint message. */
func (w *RPC) Setpoint(
	addr uint8, delay, loop uint16, setpoints []Setpoint) error {

	// write header
	header := []uint8{addr, kSetpoint}
	if err := binary.Write(w, binary.LittleEndian, header); err != nil {
		return err
	}
	// write size
	var size uint16 = uint16(6 + 4*len(setpoints))
	if err := binary.Write(w, binary.LittleEndian, size); err != nil {
		return err
	}
	// write data
	h := crc32.NewIEEE()
	mw := io.MultiWriter(w, h)
	data := []uint16{delay, loop, uint16(len(setpoints))}
	if err := binary.Write(mw, binary.LittleEndian, data); err != nil {
		return err
	}
	if err := binary.Write(mw, binary.LittleEndian, setpoints); err != nil {
		return err
	}
	// write checksum
	var checksum uint16 = uint16(h.Sum32())
	return binary.Write(w, binary.LittleEndian, checksum)
}

/* Write set test message. */
func (w *RPC) RunTests(addr uint8) error {
	return binary.Write(w, binary.LittleEndian,
		[]byte{addr, kTest, 0, 0, 0, 0})
}

/* Write set value message. */
func (w *RPC) RequestPosition(addr uint8) error {
	return binary.Write(w, binary.LittleEndian,
		[]byte{addr, kValue, 0, 0, 0, 0})
}
