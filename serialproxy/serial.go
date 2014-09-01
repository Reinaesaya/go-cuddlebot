package serialproxy

import (
	"io"

	"github.com/mikepb/go-serial"
)

// Open a serial port
func openSerialPort(name string) (port io.ReadWriteCloser, err error) {
	return serial.Open(name, serial.Options{
		Baudrate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   serial.PARITY_NONE,
	})
}
