/*

Cuddlemaster implements a web server that communicates with the
Cuddlebot actuators.

*/
package cuddle

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/phyber/negroni-gzip/gzip"
)

type customHandler func(w http.ResponseWriter, req *http.Request, body io.Reader) error

var Debug = false

func New() http.Handler {
	// set up handlers
	http.HandleFunc("/1/setpoint.json", makeHandler(setpointHandler))
	http.HandleFunc("/1/sleep.json", makeHandler(sleepHandler))
	http.Handle("/1/data.json", negroni.New(
		gzip.Gzip(gzip.DefaultCompression),
		negroni.Wrap(makeHandler(dataHandler)),
	))

	// use negroni
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(http.DefaultServeMux)

	return http.Handler(n)
}

func makeHandler(fn customHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if err := fn(w, req, req.Body); err != nil {
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, `{"ok":false,"error":"InternalServerError"}`)
			}
		}
	}
}
