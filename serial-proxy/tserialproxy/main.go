package main

import (
	tsp ".."
	"flag"
)

var addr = flag.String("addr", ":3641", "tcp service address")
var secret = flag.String("secret", "", "shared secret")
var serialPort = flag.String("port", "", "path to serial port device")

func main() {
	flag.Parse()

	proxy := &tsp.Proxy{
		Addr:       *addr,
		Secret:     *secret,
		SerialPort: *serialPort,
	}

	proxy.Listen()
}
