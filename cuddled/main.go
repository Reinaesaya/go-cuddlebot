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
	"os/exec"
	"runtime"
	"time"

	"../cuddlemaster"
	"github.com/mikepb/go-serial"
	"github.com/stretchr/graceful"
)

func main() {
	l := log.New(os.Stdout, "[cuddled] ", 0)

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

	// set low-latency option on serial port
	if runtime.GOOS == "linux" {
		l := log.New(os.Stdout, "[setserial] ", 0)
		cmd := exec.Command("/bin/setserial", *portname, "low_latency")
		if err := cmd.Run(); err != nil {
			l.Println(err)
		} else if output, err := cmd.CombinedOutput(); err != nil {
			l.Println(err)
		} else {
			l.Println(output)
		}
	}

	// connect serial port
	port, err := serial.Open(*portname, serial.Options{
		Baudrate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   serial.PARITY_NONE,
	})
	if err != nil {
		l := log.New(os.Stdout, "[serial] ", 0)
		l.Fatalln(err)
	}
	l.Println("Connected to", *portname)
	defer port.Close()

	// update setpoints in background
	go cuddlemaster.UpdateSetpoints(port)

	// set debug
	cuddlemaster.Debug = *debug
	// create server instance
	mux := cuddlemaster.New()

	// run with graceful shutdown
	graceful.Run(*listenaddr, time.Second, mux)
}
