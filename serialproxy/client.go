package serialproxy

import (
	"crypto/md5"
	"io"
	"log"
	"net"
	"strings"
)

type Client struct {
	Addr   string
	Secret string

	secretHash []byte
	encrypt    bool
}

func (c *Client) Connect() {
	var proto string
	if strings.HasPrefix(c.Addr, "/") {
		proto = "unix"
	} else {
		proto = "tcp"
	}

	conn, err := net.Dial(proto, c.Addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	c.encrypt = len(c.Secret) > 0

	if c.encrypt {
		h := md5.New()
		io.WriteString(h, c.Secret)
		c.secretHash = h.Sum(nil)
	}

	log.Printf("Connecting to %s %s", proto, c.Addr)

	var p io.ReadWriter

	// set up encryption
	if c.encrypt {
		ec, iv, err := encryptClient(conn, c.secretHash)
		if err != nil {
			log.Fatal(err)
		}
		p = ec

		// debug
		log.Printf("debug: iv: %x", iv)
	} else {
		p = conn
	}

	// do something
	count, err := p.Write([]byte{0})
	log.Printf("write (%d): %v", count, err)

	buf := make([]byte, 1)
	count, err = p.Read(buf)
	log.Printf("write (%d): %v", count, err)
}
