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

	"github.com/chimera/rs232"
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
	log.SetPrefix("conn: ")
	defer conn.Close()

	var reader io.Reader
	var writer io.Writer

	// set up encryption
	if s.encrypt {

		// generate IV
		iv := make([]byte, aes.BlockSize)
		if _, err := rand.Read(iv); err != nil {
			panic(err)
		}

		// send IV to client
		conn.Write(iv)

		// debug
		log.Printf("iv: %x", iv)

		// create AES-128 cipher
		block, err := aes.NewCipher(s.secretHash)
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

	// open serial port
	p, err := rs232.Open(s.SerialPort, rs232.Options{
		BitRate:  230400,
		DataBits: 8,
		StopBits: 1,
		Parity:   rs232.PARITY_NONE,
		Timeout:  1,
	})
	if err != nil {
		panic(err)
	}
	defer p.Close()

	done := make(chan bool, 1)
	closed := false

	// serial -> conn
	go func() {
		if c, err := io.Copy(writer, p); err != nil {
			if closed {
				log.Print("download: ", c, " (closed) ", err)
			} else {
				log.Print("download: ", c, " ", err)
			}
		} else {
			log.Print("download: ", c)
		}
		done <- true
	}()

	// conn -> serial
	go func() {
		if c, err := io.Copy(p, reader); err != nil {
			if closed {
				log.Print("upload: ", c, " (closed) ", err)
			} else {
				log.Print("upload: ", c, " ", err)
			}
		} else {
			log.Print("upload: ", c)
		}
		done <- true
	}()

	<-done
	closed = true
}
