package main

import (
	"flag"
	"log"

	. "cs.ubc.ca/spin/cuddlebot/serialproxy"
)

var addr = flag.String("addr", "127.0.0.1:3641", "tcp service address")
var secret = flag.String("secret", "", "shared secret")

func main() {
	flag.Parse()

	client := &Client{
		Addr:   *addr,
		Secret: *secret,
	}

	log.Printf("Connecting to %s", client.Addr)

	p, err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer p.Close()

	// do something
	count, err := p.Write([]byte{0})
	log.Printf("write (%d): %v", count, err)

	buf := make([]byte, 1)
	count, err = p.Read(buf)
	log.Printf("write (%d): %v", count, err)
}
