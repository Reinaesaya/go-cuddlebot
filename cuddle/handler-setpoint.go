package cuddle

import (
	"encoding/json"
	"io"
	"net/http"

	"../msgtype"
)

type setpointMessage struct {
	Addr      *msgtype.RemoteAddress `json:"addr"`
	Delay     uint16                 `json:"delay"`
	Loop      *uint16                `json:"loop"`
	Setpoints *[]uint16              `json:"setpoints"`
}

func (s *setpointMessage) bind(m *msgtype.Setpoint) error {
	if s.Addr == nil || s.Loop == nil || s.Setpoints == nil {
		return InvalidMessageError
	}

	spvalues := *s.Setpoints
	nsetpoints := len(spvalues)
	setpoints := make([]msgtype.SetpointValue, nsetpoints/2)

	for i := 0; i < nsetpoints; i += 2 {
		setpoints[i/2] = msgtype.SetpointValue{
			Duration: spvalues[i],
			Setpoint: spvalues[i+1],
		}
	}

	m.Addr = *s.Addr
	m.Delay = s.Delay
	m.Loop = *s.Loop
	m.Setpoints = setpoints

	return nil
}

func setpointHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	if req.Method != "PUT" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return MethodNotAllowed
	}

	var data setpointMessage
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return &Error{Message: err.Error()}
	}

	var message msgtype.Setpoint
	if err := data.bind(&message); err != nil {
		return err
	}

	QueueMessage(&message)

	io.WriteString(w, `{"ok":true}`)

	return nil
}
