package cuddlemaster

import (
	"log"
	"net"
	"os"

	"../msgtype"
)

var setpointQueue = make(chan msgtype.Setpoint, 10)
var setpointQueueOut = log.New(os.Stdout, "[setpoint] ", 0)
var setpointQueueErr = log.New(os.Stderr, "[setpoint] ", 0)

func ListenForSetpointUpdates(p net.Conn) {
	for {
		setpoint := <-setpointQueue
		if buf, err := setpoint.MarshalBinary(); err != nil {
			setpointQueueErr.Printf(
				"Failed marshal message for address %c: %x",
				setpoint.Addr, err.Error())
		} else if _, err := p.Write(buf); err != nil {
			setpointQueueErr.Printf(
				"Failed setpoint update for address %c: %x",
				setpoint.Addr, err.Error())
		} else {
			setpointQueueOut.Printf(
				"Completed setpoint update for address %c", setpoint.Addr)
		}
	}
}

func QueueSetpoint(setpoint *msgtype.Setpoint) {
	setpointQueue <- *setpoint
	setpointQueueOut.Printf("Queued setpoint update for address %c", setpoint.Addr)
}
