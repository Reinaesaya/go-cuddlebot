package tserialproxy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"io"
	"log"
	"net"
)

type Client struct {
	Addr   string
	Secret string

	secretHash []byte
	encrypt    bool
}

func (c *Client) Connect() {
	conn, err := net.Dial("tcp", c.Addr)
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

	c.handleConnection(conn)
}

func (c Client) handleConnection(conn net.Conn) {
	var reader io.Reader
	var writer io.Writer

	// set up encryption
	if c.encrypt {

		// receive IV from server
		iv := make([]byte, aes.BlockSize)
		if _, err := conn.Read(iv); err != nil {
			panic(err)
		}

		// debug
		log.Printf("iv: %x", iv)

		// create AES-128 cipher
		block, err := aes.NewCipher(c.secretHash)
		if err != nil {
			panic(err)
		}

		// initialize stream
		stream := cipher.NewOFB(block, iv)
		reader = &cipher.StreamReader{S: stream, R: conn}
		writer = &cipher.StreamWriter{S: stream, W: conn}
	} else {
		reader = conn
		writer = conn
	}

	// do something
	if c, err := writer.Write([]byte{0}); err != nil {
		log.Print("write: ", err)
	} else {
		log.Print("write: ", c)
	}

	buf := make([]byte, 1)
	if c, err := reader.Read(buf); err != nil {
		log.Print("read: ", err)
	} else {
		log.Print("read: ", c, " ", buf)
	}
}
