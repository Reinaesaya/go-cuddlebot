// Securely exposes a local serial port on the network
package tserialproxy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"io"
	"log"
	"net"

	"github.com/ziutek/serial"
)

type Proxy struct {
	Addr       string
	Secret     string
	SerialPort string

	secretHash []byte
	encrypt    bool
}

func (s *Proxy) Listen() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	s.encrypt = len(s.Secret) > 0

	if s.encrypt {
		h := md5.New()
		io.WriteString(h, s.Secret)
		s.secretHash = h.Sum(nil)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		// no goroutines; only accept one connection at a time
		s.handleConnection(conn)
	}

	return
}

func (s Proxy) handleConnection(conn net.Conn) {
	defer conn.Close()

	var reader io.Reader
	var writer io.Writer

	if s.encrypt {
		block, err := aes.NewCipher(s.secretHash)
		if err != nil {
			panic(err)
		}

		iv := make([]byte, aes.BlockSize)
		if _, err := rand.Read(iv); err != nil {
			panic(err)
		}

		log.Printf("%x", iv)

		stream := cipher.NewOFB(block, iv)
		reader = &cipher.StreamReader{S: stream, R: conn}
		writer = &cipher.StreamWriter{S: stream, W: conn}

		conn.Write(iv)
	} else {
		reader = conn
		writer = conn
	}

	_ = reader
	_ = writer
}
