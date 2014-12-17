package main

import (
	"bufio"
	"bytes"
	"encoding"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path"

	"../cuddle"
	"../msgtype"
)

var debug = flag.Bool("debug", false, "print debug messages")
var n = flag.Bool("n", false, "parse arguments, but don't send command")

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
	args := flag.Args()

	if flag.NArg() < 1 {
		fatalUsage()
	} else if *help {
		flag.Usage()
		os.Exit(0)
	}

	// open serial port
	var port io.ReadWriteCloser
	if !*n {
		port, err := cuddle.OpenPort(*portname)
		if err != nil {
			log.Fatalln(err)
		}
		defer port.Close()
		log.Println("Connected to", *portname)
	}

	// run command
	switch true {
	case *ribs:
		runcmd(port, msgtype.RibsAddress, args)
	case *purr:
		runcmd(port, msgtype.PurrAddress, args)
	case *spine:
		runcmd(port, msgtype.SpineAddress, args)
	case *headx:
		runcmd(port, msgtype.HeadXAddress, args)
	case *heady:
		runcmd(port, msgtype.HeadYAddress, args)
	}
}

func runcmd(conn io.ReadWriter, addr msgtype.RemoteAddress, args []string) {
	cmd, args := args[0], args[1:]

	// run command
	switch cmd {
	case "setpid":
		if len(args) != 3 {
			fatalUsage()
		}

		kpS, kiS, kdS := args[0], args[1], args[2]

		var kp, ki, kd float32
		fmt.Fscanf(bytes.NewBufferString(kpS), "%f", &kp)
		fmt.Fscanf(bytes.NewBufferString(kiS), "%f", &ki)
		fmt.Fscanf(bytes.NewBufferString(kdS), "%f", &kd)

		if *debug {
			log.Printf("parsed pid kp=%f ki=%f kd=%f", kp, ki, kd)
		}

		sendcmd(conn, &msgtype.SetPID{addr, kp, ki, kd})

	case "setpoint":
		if len(args) < 3 {
			fatalUsage()
		}

		delayS, loopS, args := args[0], args[1], args[2:]

		var delay, loop int
		fmt.Fscanf(bytes.NewBufferString(delayS), "%d", &delay)
		if loopS == "forever" {
			loop = 0xffff
		} else {
			fmt.Fscanf(bytes.NewBufferString(loopS), "%d", &loop)
		}

		if delay < 0 || loop < 0 {
			log.Fatalln("Error: delay and loop must be positive")
		}

		if len(args)%2 != 0 {
			log.Fatalln("Error: duration and setpoint must be given in pairs")
		}

		setpoints := make([]msgtype.SetpointValue, len(args)/2)
		for i := 0; i < len(args); i += 2 {
			durationS, setpointS := args[i], args[i+1]

			var duration, setpoint int
			if durationS == "forever" {
				duration = 0xffff
			} else {
				fmt.Fscanf(bytes.NewBufferString(durationS), "%d", &duration)
			}
			fmt.Fscanf(bytes.NewBufferString(setpointS), "%d", &setpoint)

			if duration < 0 || setpoint < 0 {
				log.Fatalln("Error: duration and setpoint must be positive")
			}

			j := i / 2

			setpoints[j].Duration = uint16(duration)
			setpoints[j].Setpoint = uint16(setpoint)
		}

		sendcmd(conn, &msgtype.Setpoint{addr,
			uint16(delay), uint16(loop), setpoints})

	case "ping":
		if len(args) != 0 {
			fatalUsage()
		}

		sendcmd(conn, &msgtype.Ping{addr})

		if !*n {
			buf := make([]byte, 1)
			conn.Read(buf)
			os.Stdout.Write(buf)
			os.Stdout.WriteString("\n")
		}

	case "test":
		if len(args) != 0 {
			fatalUsage()
		}

		sendcmd(conn, &msgtype.Test{addr})

	case "value":
		if len(args) != 0 {
			fatalUsage()
		}

		sendcmd(conn, &msgtype.Value{addr})

		if !*n {
			if line, _, err := bufio.NewReader(conn).ReadLine(); err != nil {
				log.Fatalln(err)
			} else {
				os.Stdout.Write(line)
				os.Stdout.WriteString("\n")
			}
		}

	default:
		fatalUsage()
	}

	if *debug {
		log.Printf("sent %s message to address %d", cmd, addr)
	}
}

func sendcmd(conn io.Writer, m encoding.BinaryMarshaler) {
	if bs, err := m.MarshalBinary(); err != nil {
		log.Fatalln(err)
	} else if !*n {
		if _, err := conn.Write(bs); err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("ok", m)
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
