package cuddle

import (
	"encoding"
	"log"
	"net"
	"os"
	"time"
)

var messageQueue = make(chan encoding.BinaryMarshaler, 10)
var messageQueueOut = log.New(os.Stdout, "[message] ", 0)
var messageQueueErr = log.New(os.Stderr, "[message] ", 0)

func SendQueuedMessagesTo(p net.Conn) {
	for {
		message := <-messageQueue
		p.SetWriteDeadline(time.Now().Add(10 * time.Millisecond))
		if buf, err := message.MarshalBinary(); err != nil {
			messageQueueErr.Printf("Failed marshal message %x %x",
				err.Error(), message)
		} else if _, err := p.Write(buf); err != nil {
			messageQueueErr.Printf("Failed to send message %x %x",
				err.Error(), message)
		} else {
			messageQueueOut.Printf("Completed message send %x", message)
		}
	}
}

func QueueMessage(message encoding.BinaryMarshaler) {
	messageQueue <- message
}
