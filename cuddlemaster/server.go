/*

Cuddlemaster implements a web server that communicates with the
Cuddlebot actuators.

*/
package cuddlemaster

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
	http.HandleFunc("/setpoint", makeHandler(setpointHandler))
	http.Handle("/data", negroni.New(
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
		if err := fn(w, req, req.Body); err != nil {
			if err := json.NewEncoder(w).Encode(err); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, `{"ok":false,"error":"InternalServerError"}`)
			}
		}
	}
}
