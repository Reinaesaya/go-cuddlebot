package serialproxy

import (
	"crypto/md5"
	"io"
	"net"
	"strings"
)

type Client struct {
	Addr   string
	Secret string
}

func (c *Client) Connect() (io.ReadWriteCloser, error) {
	var proto string
	if strings.HasPrefix(c.Addr, "/") {
		proto = "unix"
	} else {
		proto = "tcp"
	}

	conn, err := net.Dial(proto, c.Addr)
	if err != nil {
		return nil, err
	}

	// set up encryption
	if len(c.Secret) > 0 {
		h := md5.New()
		io.WriteString(h, c.Secret)
		hash := h.Sum(nil)

		ec, _, err := encryptClient(conn, hash)
		if err != nil {
			return nil, err
		}

		p := &readWriteCloser{
			ReadWriter: ec,
			Closer:     conn,
		}

		return p, nil
	}

	return conn, nil
}
