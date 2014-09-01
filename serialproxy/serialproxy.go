/*

Package serialproxy implements a simple encrypting proxy to expose a
local serial port over the network.

*/
package serialproxy

import (
	"io"
)

type readWriter struct {
	io.Reader
	io.Writer
}

type readWriteCloser struct {
	io.ReadWriter
	io.Closer
}
