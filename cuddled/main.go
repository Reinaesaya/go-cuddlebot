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
	"os/signal"
	"runtime"
	"syscall"

	"../cuddlemaster"
	"github.com/mikepb/go-serial"
)

func main() {

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
		cmd := exec.Command("/bin/setserial", *portname, "low_latency")
		if err := cmd.Run(); err != nil {
			log.Println("exec:", err)
		} else if output, err := cmd.CombinedOutput(); err != nil {
			log.Println("exec:", err)
		} else {
			log.Println(output)
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
		log.Fatalln("serial:", err)
	}
	log.Println("Opened serial port:", *portname)
	defer port.Close()

	// set debug
	cuddlemaster.Debug = *debug
	// set up listener for web server
	cuddlemaster.ListenAndServe(*listenaddr, port)

	// set up interrupt signal channel
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// handle connections or interrupts
	killSignal := <-interrupt
	log.Println("Got signal:", killSignal)
	log.Println("Shutting down...")
	os.Exit(0)
}
