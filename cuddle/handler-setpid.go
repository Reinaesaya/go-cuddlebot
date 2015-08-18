package cuddle

import (
	"encoding/json"
	"io"
	"net/http"

	"../msgtype"
)

type setpidMessage struct {
	Addr      *msgtype.RemoteAddress `json:"addr"`
	Kp	   *float32            	`json:"kp"`
	Ki	   *float32            	`json:"ki"`
	Kd	   *float32            	`json:"kd"`
}

func (s *setpidMessage) bind(m *msgtype.SetPID) error {
	if s.Addr == nil {
		return InvalidMessageError
	}

	m.Addr = *s.Addr
	m.Kp = *s.Kp
	m.Ki = *s.Ki
	m.Kd = *s.Kd

	return nil
}

func setpidHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	if req.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return MethodNotAllowed
	}

	var data setpidMessage
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return &Error{Message: err.Error()}
	}

	var message msgtype.SetPID
	if err := data.bind(&message); err != nil {
		return err
	}

	QueueMessage(&message)

	io.WriteString(w, `{"ok":true}`)

	return nil
}
