package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/mikepb/go-serial"

	"../msgtype"
)

var debug = flag.Bool("debug", false, "print debug messages")

func main() {
	// define actuator flags
	help := flag.Bool("help", false, "print help")
	ribs := flag.Bool("ribs", false, "send command to ribs actuator")
	purr := flag.Bool("purr", false, "send command to purr actuator")
	spine := flag.Bool("spine", false, "send command to spine actuator")
	headx := flag.Bool("headx", false, "send command to head yaw actuator")
	heady := flag.Bool("heady", false, "send command to head pitch actuator")

	// parse flags
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 2 {
		os.Exit(1)
	} else if *help {
		flag.Usage()
		os.Exit(0)
	}

	// open serial port
	portname := flag.Arg(0)
	port, err := serial.Open(portname, serial.Options{
		Baudrate: 115200,
		DataBits: 8,
		StopBits: 1,
		Parity:   serial.PARITY_NONE,
	})
	if err != nil {
		log.Fatal(err)
	} else if *debug {
		log.Printf("opened %s", portname)
	}
	defer port.Close()

	// buffer output
	bw := bufio.NewWriter(port)

	// rpc wrapper
	rpc := msgtype.RPC{Writer: bw}

	// run command
	switch true {
	case *ribs:
		runcmd(rpc, msgtype.ADDR_RIBS)
	case *purr:
		runcmd(rpc, msgtype.ADDR_PURR)
	case *spine:
		runcmd(rpc, msgtype.ADDR_SPINE)
	case *headx:
		runcmd(rpc, msgtype.ADDR_HEAD_YAW)
	case *heady:
		runcmd(rpc, msgtype.ADDR_HEAD_PITCH)
	}

	// flush output
	bw.Flush()
}

func runcmd(rpc msgtype.RPC, addr uint8) {
	// run command
	switch flag.Arg(1) {
	case "setpid":
		if flag.NArg() < 5 {
			fatalUsage()
		}

		var kp, ki, kd float32
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(2)), "%f", &kp)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(3)), "%f", &ki)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(4)), "%f", &kd)

		if *debug {
			log.Printf("parsed pid kp=%f ki=%f kd=%f", kp, ki, kd)
		}

		rpc.SetPID(addr, kp, ki, kd)

	case "setpoint":
		if flag.NArg() < 6 {
			fatalUsage()
		}

		if flag.NArg()%2 != 0 {
			log.Fatal(os.Stderr, "Error: duration and setpoint must be given in pairs")
		}

		var delayInt, loopInt int

		fmt.Fscanf(bytes.NewBufferString(flag.Arg(2)), "%d", &delayInt)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(3)), "%d", &loopInt)

		if delayInt < 0 || loopInt < 0 {
			log.Fatal(os.Stderr, "Error: delay and loop must be positive")
		}

		delay := uint16(delayInt)
		loop := uint16(loopInt)

		setpoints := make([]msgtype.Setpoint, (flag.NArg()-4)/2)
		for i := 4; i < flag.NArg(); i += 2 {
			var duration, setpoint int

			fmt.Fscanf(bytes.NewBufferString(flag.Arg(i)), "%d", &duration)
			fmt.Fscanf(bytes.NewBufferString(flag.Arg(i+1)), "%d", &setpoint)

			if delayInt < 0 || loopInt < 0 {
				log.Fatal(os.Stderr, "Error: duration and setpoint must be positive")
			}

			j := (i - 4) / 2

			setpoints[j].Duration = uint16(duration)
			setpoints[j].Setpoint = uint16(setpoint)
		}

		rpc.Setpoint(addr, delay, loop, setpoints)

	case "ping":
		rpc.Ping(addr)
	case "test":
		rpc.RunTests(addr)
	case "value":
		rpc.RequestPosition(addr)

	default:
		fatalUsage()
	}

	if *debug {
		log.Printf("sent %s message to address %d", flag.Arg(1), addr)
	}
}

var header = `Cuddlespeak is a tool for testing the Cuddlebot actuators.

Usage:

    %s [flags] port command [arguments]

The flags are:

`

var footer = `

The commands are:

    setpid      set the PID coefficients
    setpoint    send setpoints
    ping        send a ping
    test        send test command
    value       read motor position

The setpid command accepts these arguments:

    kp          float: the P coefficient
    ki          float: the I coefficient
    kd          float: the D coefficient

The setpoint command accepts these arguments:

    delay       uint: the P coefficient
    loop        uint: the number of times to repeat this group of
                setpoints or "forever" to loop indefinitely
    [duration setpoint]+
                one or more setpoints consisting of groups of two
                uints in order: duration setpoint; with duration in
                milliseconds and setpoint in (1 / 2^16) increments of
                a circle

Examples:

    $ %s -ribs /dev/ttyUSB0 setpid 40.4 1.0 -1.0

    $ %s -ribs /dev/ttyUSB0 setpoint 0 forever 1000 26075 1000 0

    $ %s -ribs /dev/ttyUSB0 ping

    $ %s -ribs /dev/ttyUSB0 test
    ... test results ...

    $ %s -ribs /dev/ttyUSB0 value
    0.1

`

func usage() {
	name := path.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, header, name)

	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(os.Stderr, "    -%-10s %s\n", f.Name, f.Usage)
	})

	fmt.Fprintf(os.Stderr, footer, name, name, name, name, name)
}

func fatalUsage() {
	usage()
	os.Exit(1)
}
