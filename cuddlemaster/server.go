/*

Cuddlemaster implements a web server that communicates with the
Cuddlebot actuators.

*/
package cuddlemaster

import (
	"log"
	"net"
	"net/http"
	"time"
)

var Debug = false
var port net.Conn

func ListenAndServe(addr string, p net.Conn) {
	port = p

	// set up handlers
	http.HandleFunc("/pid", makeHandler(setPIDHandler))
	http.HandleFunc("/setpoint", makeHandler(setpointHandler))
	http.HandleFunc("/data", makeHandler(dataHandler))

	// listen and serve
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			log.Fatalln("http:", err)
		}
	}()
	time.Sleep(time.Millisecond * 500)
	log.Println("Listening on:", addr)
}
