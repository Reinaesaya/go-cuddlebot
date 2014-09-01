// Securely exposes a local serial port on the network
package serialproxy

import (
	"crypto/md5"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
)

type Server struct {
	Addr       string
	Secret     string
	SerialPort string
	Logger     *log.Logger

	encrypt    bool
	secretHash []byte
}

func (s *Server) Listen() {
	if s.Logger == nil {
		s.Logger = log.New(os.Stderr, "[serialproxy] ", log.LstdFlags)
	}

	var proto string
	if strings.HasPrefix(s.Addr, "/") {
		proto = "unix"
	} else {
		proto = "tcp"
	}

	listener, err := net.Listen(proto, s.Addr)
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

	s.Logger.Printf("Listening on %s %s", proto, s.Addr)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				s.Logger.Fatal(err)
			}
			defer conn.Close()
			// no goroutines; only accept one connection at a time
			s.handleConnection(conn)
		}
	}()

	// handle signals
	signotify := make(chan os.Signal, 1)
	signal.Notify(signotify, os.Interrupt, os.Kill)
	sig := <-signotify
	s.Logger.Printf("Got %v signal, shutting down", sig)
}

func (s Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	var c io.ReadWriter

	// set up encryption
	if s.encrypt {
		ec, iv, err := encryptServer(conn, s.secretHash)
		if err != nil {
			log.Fatal(err)
		}
		c = ec

		// debug
		log.Printf("debug: iv: %x", iv)

		// send IV to client
		conn.Write(iv)
	} else {
		c = conn
	}

	// open serial port
	port, err := openSerialPort(s.SerialPort)
	if err != nil {
		log.Fatal(err)
	}
	defer port.Close()

	// pipe data to/from conn and port
	s.pipe(port, c)
}

// Pipe I/O to/from remote and local
func (s *Server) pipe(local io.ReadWriter, remote io.ReadWriter) {
	var wg sync.WaitGroup

	// remote -> local
	go func() {
		wg.Add(1)
		c, err := io.Copy(local, remote)
		s.Logger.Printf("serialproxy: conn: remote -> local (%d): %v", c, err)
		wg.Done()
	}()

	// local -> remote
	go func() {
		wg.Add(1)
		c, err := io.Copy(remote, local)
		s.Logger.Printf("serialproxy: conn: local -> remote (%d): %v", c, err)
		wg.Done()
	}()

	wg.Wait()
}
