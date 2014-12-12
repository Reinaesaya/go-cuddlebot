/*

CuddleD runs Cuddlemaster.

Interrupt handling based on example at:
https://github.com/takama/daemon

*/
package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/stretchr/graceful"

	"../cuddle"
)

func main() {
	l := log.New(os.Stdout, "[cuddled] ", 0)
	e := log.New(os.Stderr, "[cuddled] ", 0)

	// define flags
	debug := flag.Bool("debug", false, "print debug messages")
	help := flag.Bool("help", false, "print help")
	portname := flag.String("port", "/dev/ttyUSB0", "the serial port name")
	listenaddr := flag.String("listen", ":http", "the address on which to listen")

	// parse flags
	flag.Parse()

	// print help
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// do not accept arguments
	if flag.NArg() > 0 {
		flag.Usage()
		os.Exit(1)
	}

	// connect serial port
	port, err := cuddle.OpenPort(*portname)
	if err != nil {
		e.Fatalln(err)
	}
	defer port.Close()
	l.Println("Connected to", *portname)

	// update setpoints in background
	go cuddle.SendQueuedMessagesTo(port)

	// set debug
	cuddle.Debug = *debug
	// create server instance
	mux := cuddle.New()

	// run with graceful shutdown
	graceful.Run(*listenaddr, time.Second, mux)
}
