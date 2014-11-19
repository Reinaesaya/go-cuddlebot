package cuddlemaster

import (
	"encoding/json"
	"io"
	"net/http"

	"../msgtype"
)

type setpointMessage struct {
	Addr      *string   `json:"addr"`
	Delay     *uint16   `json:"delay"`
	Loop      *uint16   `json:"loop"`
	Setpoints *[]uint16 `json:"setpoints"`
}

func (s *setpointMessage) bind(m *msgtype.Setpoint) error {
	var addr uint8

	if s.Addr == nil || s.Delay == nil || s.Loop == nil || s.Setpoints == nil {
		return InvalidMessageError
	}

	switch *s.Addr {
	case "ribs":
		addr = msgtype.ADDR_RIBS
	case "purr":
		addr = msgtype.ADDR_PURR
	case "spine":
		addr = msgtype.ADDR_SPINE
	case "headx":
		addr = msgtype.ADDR_HEAD_YAW
	case "heady":
		addr = msgtype.ADDR_HEAD_PITCH
	default:
		return InvalidAddressError
	}

	spvalues := *s.Setpoints
	nsetpoints := len(spvalues)
	setpoints := make([]msgtype.SetpointValue, nsetpoints/2)

	for i := 0; i < nsetpoints; i += 2 {
		spvalue := setpoints[i/2]
		spvalue.Duration = spvalues[i]
		spvalue.Setpoint = spvalues[i+1]
	}

	m.Addr = addr
	m.Delay = *s.Delay
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

	setpointQueue <- message
	io.WriteString(w, `{"ok":true}`)

	return nil
}
