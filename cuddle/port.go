package cuddle

import (
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func OpenPort(name string) (io.ReadWriteCloser, error) {

	// http://projectgus.com/2011/10/notes-on-ftdi-latency-with-arduino/
	// http://faumarz.blogspot.ca/2014/06/change-ftdi-usb-serial-latency-in-linux.html
	// https://forum.openwrt.org/viewtopic.php?id=47367
	/*

	   Linux:

	   # sed -i 's/=\d+/=1/' /etc/modules.d/usb-serial-ftdi
	   # stty -F /dev/ttyUSB0 115200 raw
	   # setserial /dev/ttyUSB0 low_latency
	   # cat /sys/bus/usb-serial/devices/ttyUSB0/latency_timer
	   16
	   # echo 1 > /sys/bus/usb-serial/devices/ttyUSB0/latency_timer
	   # cat /sys/bus/usb-serial/devices/ttyUSB0/latency_timer
	   1

	   In code:

	   struct serial_struct ser_info;
	   ioctl(serial, TIOCGSERIAL, &ser_info);
	   ser_info.flags |= ASYNC_LOW_LATENCY;
	   ioctl(serial, TIOCSSERIAL, &ser_info);

	*/

	if runtime.GOOS == "linux" {
		// execWithLogging("setserial", "/bin/setserial", name, "low_latency")
		execWithLogging("stty", "/bin/stty", "-F", name, "115200", "raw")
	} else if runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" {
		execWithLogging("setserial", "/bin/stty", "-f", name, "115200", "raw")
	}

	return os.OpenFile(name, os.O_RDWR, 0600)
}

func execWithLogging(name string, args ...string) {
	l := log.New(os.Stdout, "["+name+"] ", 0)
	l.Println(strings.Join(args, " "))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		l := log.New(os.Stderr, "["+name+"] ", 0)
		l.Println(err)
	}
}
