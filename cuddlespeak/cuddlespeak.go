package main

import (
	"bytes"
	"encoding"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"
	"time"

	"../cuddlemaster"
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

	portname := flag.String("port", "/dev/ttyUSB0", "the serial port name")

	// parse flags
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		fatalUsage()
	} else if *help {
		flag.Usage()
		os.Exit(0)
	}

	// open serial port
	port, err := cuddlemaster.OpenPort(*portname)
	if err != nil {
		log.Fatalln(err)
	}
	defer port.Close()
	log.Println("Connected to", *portname)

	// run command
	switch true {
	case *ribs:
		runcmd(port, msgtype.RibsAddress)
	case *purr:
		runcmd(port, msgtype.PurrAddress)
	case *spine:
		runcmd(port, msgtype.SpineAddress)
	case *headx:
		runcmd(port, msgtype.HeadXAddress)
	case *heady:
		runcmd(port, msgtype.HeadYAddress)
	}
}

func runcmd(conn net.Conn, addr msgtype.RemoteAddress) {
	// run command
	switch flag.Arg(0) {
	case "setpid":
		if flag.NArg() < 4 {
			fatalUsage()
		}

		var kp, ki, kd float32
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(1)), "%f", &kp)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(2)), "%f", &ki)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(3)), "%f", &kd)

		if *debug {
			log.Printf("parsed pid kp=%f ki=%f kd=%f", kp, ki, kd)
		}

		sendcmd(conn, &msgtype.SetPID{addr, kp, ki, kd})

	case "setpoint":
		if flag.NArg() < 5 {
			fatalUsage()
		}

		if flag.NArg()%2 != 1 {
			log.Fatal(os.Stderr, "Error: duration and setpoint must be given in pairs")
		}

		var delayInt, loopInt int

		fmt.Fscanf(bytes.NewBufferString(flag.Arg(1)), "%d", &delayInt)
		fmt.Fscanf(bytes.NewBufferString(flag.Arg(2)), "%d", &loopInt)

		if delayInt < 0 || loopInt < 0 {
			log.Fatal(os.Stderr, "Error: delay and loop must be positive")
		}

		delay := uint16(delayInt)
		loop := uint16(loopInt)

		setpoints := make([]msgtype.SetpointValue, (flag.NArg()-3)/2)
		for i := 3; i < flag.NArg(); i += 2 {
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

		sendcmd(conn, &msgtype.Setpoint{addr, delay, loop, setpoints})

	case "ping":
		sendcmd(conn, &msgtype.Ping{addr})
		conn.SetReadDeadline(time.Now().Add(time.Second))
		io.Copy(os.Stdout, conn)

	case "test":
		sendcmd(conn, &msgtype.Test{addr})
		conn.SetReadDeadline(time.Now().Add(time.Minute * 5))
		io.Copy(os.Stdout, conn)

	case "value":
		sendcmd(conn, &msgtype.Value{addr})
		conn.SetReadDeadline(time.Now().Add(time.Second))
		io.Copy(os.Stdout, conn)

	default:
		fatalUsage()
	}

	if *debug {
		log.Printf("sent %s message to address %d", flag.Arg(1), addr)
	}
}

func sendcmd(conn io.Writer, m encoding.BinaryMarshaler) {
	if bs, err := m.MarshalBinary(); err != nil {
		log.Fatalln(err)
	} else if _, err := conn.Write(bs); err != nil {
		log.Fatalln(err)
	}
}

var header = `Cuddlespeak is a tool for testing the Cuddlebot actuators.

Usage:

    %s [flags] command [arguments]

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

    $ %s -ribs setpid 40.4 1.0 -1.0

    $ %s -ribs setpoint 0 forever 1000 26075 1000 0

    $ %s -ribs ping

    $ %s -ribs test
    ... test results ...

    $ %s -ribs value
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
