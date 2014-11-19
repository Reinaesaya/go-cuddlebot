package cuddlemaster

import (
	"log"
	"net"
	"os"

	"../msgtype"
)

var setpointQueue = make(chan msgtype.Setpoint, 10)

func UpdateSetpoints(p net.Conn) {
	l := log.New(os.Stdout, "[setpoint] ", 0)
	for {
		setpoint := <-setpointQueue
		if _, err := setpoint.WriteTo(p); err != nil {
			l := log.New(os.Stderr, "[setpoint] ", 0)
			l.Println(err)
		} else {
			l.Printf("Updated setpoints for address %c", setpoint.Addr)
		}
	}
}
