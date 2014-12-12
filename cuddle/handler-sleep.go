package cuddle

import (
	"encoding/json"
	"io"
	"net/http"

	"../msgtype"
)

type sleepMessage struct {
	Addr *[]msgtype.RemoteAddress `json:"addr"`
}

func sleepHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	if req.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return MethodNotAllowed
	}

	var data sleepMessage
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return &Error{Message: err.Error()}
	}

	if data.Addr == nil || len(*data.Addr) == 0 {
		return InvalidMessageError
	}

	for _, addr := range *data.Addr {
		QueueMessage(&msgtype.Sleep{addr})
	}

	io.WriteString(w, `{"ok":true}`)

	return nil
}
