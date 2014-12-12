package main

import (
	"bufio"
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

	"../cuddle"
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
	port, err := cuddle.OpenPort(*portname)
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

		kpS := flag.Arg(1)
		kiS := flag.Arg(2)
		kdS := flag.Arg(3)

		var kp, ki, kd float32
		fmt.Fscanf(bytes.NewBufferString(kpS), "%f", &kp)
		fmt.Fscanf(bytes.NewBufferString(kiS), "%f", &ki)
		fmt.Fscanf(bytes.NewBufferString(kdS), "%f", &kd)

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

		delayS := flag.Arg(1)
		loopS := flag.Arg(2)

		var delay, loop int
		fmt.Fscanf(bytes.NewBufferString(delayS), "%d", &delay)
		if loopS == "forever" {
			loop = 0xffff
		} else {
			fmt.Fscanf(bytes.NewBufferString(loopS), "%d", &loop)
		}

		if delay < 0 || loop < 0 {
			log.Fatal(os.Stderr, "Error: delay and loop must be positive")
		}

		setpoints := make([]msgtype.SetpointValue, (flag.NArg()-3)/2)
		for i := 3; i < flag.NArg(); i += 2 {
			durationS := flag.Arg(i)
			setpointS := flag.Arg(i + 1)

			var duration, setpoint int
			if durationS == "forever" {
				duration = 0xffff
			} else {
				fmt.Fscanf(bytes.NewBufferString(durationS), "%d", &duration)
			}
			fmt.Fscanf(bytes.NewBufferString(setpointS), "%d", &setpoint)

			if duration < 0 || setpoint < 0 {
				log.Fatal(os.Stderr, "Error: duration and setpoint must be positive")
			}

			j := (i - 3) / 2

			setpoints[j].Duration = uint16(duration)
			setpoints[j].Setpoint = uint16(setpoint)
		}

		sendcmd(conn, &msgtype.Setpoint{addr,
			uint16(delay), uint16(loop), setpoints})

	case "ping":
		sendcmd(conn, &msgtype.Ping{addr})
		conn.SetReadDeadline(time.Now().Add(time.Second))
		buf := make([]byte, 1)
		conn.Read(buf)
		os.Stdout.Write(buf)
		os.Stdout.WriteString("\n")

	case "test":
		sendcmd(conn, &msgtype.Test{addr})

	case "value":
		sendcmd(conn, &msgtype.Value{addr})
		conn.SetReadDeadline(time.Now().Add(time.Second))
		if line, _, err := bufio.NewReader(conn).ReadLine(); err != nil {
			log.Fatalln(err)
		} else {
			os.Stdout.Write(line)
			os.Stdout.WriteString("\n")
		}

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
