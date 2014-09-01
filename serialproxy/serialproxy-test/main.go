package main

import (
	tsp ".."
	"flag"
	"log"
)

var addr = flag.String("addr", "127.0.0.1:3641", "tcp service address")
var secret = flag.String("secret", "", "shared secret")

func main() {
	flag.Parse()

	client := &tsp.Client{
		Addr:   *addr,
		Secret: *secret,
	}

	client.Connect()
}
