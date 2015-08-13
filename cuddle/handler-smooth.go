package cuddle

import (
	"encoding/json"
	"io"
	"net/http"

	"../msgtype"
)

type smoothMessage struct {
	Addr      *msgtype.RemoteAddress `json:"addr"`
	Time 	   *uint16             	`json:"time"`
	Setpoint *[]uint16              `json:"setpoint"`
}

func (s *smoothMessage) bind(m *msgtype.Smooth) error {
	if s.Addr == nil || s.Time == nil || s.Setpoint == nil {
		return InvalidMessageError
	}

	spvalue := *s.Setpoint
	setpoint := make([]msgtype.SetpointValue, 1)
	setpoint[0] = msgtype.SetpointValue{
		Duration: spvalue[0],
		Setpoint: spvalue[1],
	}

	m.Addr = *s.Addr
	m.Time = *s.Time
	m.Setpoint = setpoint

	return nil
}

func smoothHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	if req.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return MethodNotAllowed
	}

	var data smoothMessage
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return &Error{Message: err.Error()}
	}

	var message msgtype.Smooth
	if err := data.bind(&message); err != nil {
		return err
	}

	QueueMessage(&message)

	io.WriteString(w, `{"ok":true}`)

	return nil
}
