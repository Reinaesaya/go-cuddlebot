// Securely exposes a local serial port on the network
package tserialproxy

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type readWriter struct {
	io.Reader
	io.Writer
}

// Encrypt the pipe using AES128
func encrypt(pipe io.ReadWriter, secret []byte, iv []byte) (io.ReadWriter, error) {
	// create AES-128 cipher
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	// initialize stream
	s := cipher.NewOFB(block, iv)
	r := &cipher.StreamReader{S: s, R: pipe}
	w := &cipher.StreamWriter{S: s, W: pipe}

	// wrap in reader-writer
	epipe := &readWriter{Reader: r, Writer: w}

	return epipe, nil
}

// Encrypt the pipe for a server connection
func encryptServer(pipe io.ReadWriter, secret []byte) (io.ReadWriter, []byte, error) {
	// generate IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, nil, err
	}

	// encrypt
	epipe, err := encrypt(pipe, secret, iv)
	if err != nil {
		return nil, nil, err
	}

	return epipe, iv, err
}

// Encrypt the client for a server connection
func encryptClient(pipe io.ReadWriter, secret []byte) (io.ReadWriter, []byte, error) {
	// receive IV from remote
	iv := make([]byte, aes.BlockSize)
	if _, err := pipe.Read(iv); err != nil {
		return nil, nil, err
	}

	// encrypt
	epipe, err := encrypt(pipe, secret, iv)
	if err != nil {
		return nil, nil, err
	}

	return epipe, iv, err
}
