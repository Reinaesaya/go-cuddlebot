package main

import (
	"flag"

	. "cs.ubc.ca/spin/cuddlebot/serialproxy"
)

var addr = flag.String("addr", ":3641", "tcp service address")
var secret = flag.String("secret", "", "shared secret")
var serialPort = flag.String("port", "", "path to serial port device")

func main() {
	flag.Parse()

	proxy := &Server{
		Addr:       *addr,
		Secret:     *secret,
		SerialPort: *serialPort,
	}

	proxy.Listen()
}
