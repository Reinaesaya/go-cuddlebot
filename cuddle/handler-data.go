package cuddle

import (
	"io"
	"net/http"
)

type dataMessage struct {
}

func dataHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	if req.Method != "GET" {
		return MethodNotAllowed
	}
	// http://wiki.analog.com/software/linux/docs/iio/iio
	// https://archive.fosdem.org/2012/schedule/event/693/127_iio-a-new-subsystem.pdf
	return NotImplementedError
}
