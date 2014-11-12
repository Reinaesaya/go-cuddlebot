package cuddlemaster

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"../msgtype"
)

var responseWrittenError = errors.New("Response written placeholder error")

type customHandler func(w http.ResponseWriter, req *http.Request, body io.Reader) error

type setPIDMessage struct {
	Addr *string  `json:"addr"`
	Kp   *float32 `json:"kp"`
	Ki   *float32 `json:"ki"`
	Kd   *float32 `json:"kd"`
}

type setpointMessage struct {
	Addr      *string   `json:"addr"`
	Delay     uint16    `json:"delay"`
	Loop      uint16    `json:"loop"`
	Setpoints *[]uint16 `json:"setpoints"`
}

type dataMessage struct {
}

func makeHandler(fn customHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// encoding := req.Header.Get("Content-Encoding")
		body := req.Body
		if err := fn(w, req, body); err == responseWrittenError {
			// no-op
		} else if err != nil {
			json.NewEncoder(w).Encode(&jsonError{OK: false, Error: err.Error()})
		} else {
			io.WriteString(w, `{"ok":true}`)
		}
	}
}

func addrToInt(v string) (uint8, error) {
	switch v {
	case "ribs":
		return msgtype.ADDR_RIBS, nil
	case "purr":
		return msgtype.ADDR_PURR, nil
	case "spine":
		return msgtype.ADDR_SPINE, nil
	case "headyaw":
		return msgtype.ADDR_HEAD_YAW, nil
	case "headpitch":
		return msgtype.ADDR_HEAD_PITCH, nil
	default:
		return msgtype.ADDR_INVALID, AddressError
	}
}

func setPIDHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	var data setPIDMessage

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return err
	}

	if data.Addr == nil || data.Kp == nil || data.Ki == nil || data.Kd == nil {
		return MissingFieldsError
	}

	addr, err := addrToInt(*data.Addr)
	if err != nil {
		return err
	}

	message := &msgtype.SetPID{
		Addr: addr,
		Kp:   *data.Kp,
		Ki:   *data.Ki,
		Kd:   *data.Kd,
	}

	if _, err := message.WriteTo(port); err != nil {
		return err
	}

	return nil
}

func setpointHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	var data setpointMessage

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return err
	}

	if data.Addr == nil || data.Setpoints == nil {
		return MissingFieldsError
	}

	addr, err := addrToInt(*data.Addr)
	if err != nil {
		return err
	}

	spdata := *data.Setpoints
	nsetpoints := len(spdata)
	if nsetpoints == 0 || nsetpoints%2 != 0 {
		return InvalidSetpointError
	}

	setpoints := make([]msgtype.SetpointValue, nsetpoints/2)
	for i := 0; i < nsetpoints; i += 2 {
		spvalue := setpoints[i/2]
		spvalue.Duration = spdata[i]
		spvalue.Setpoint = spdata[i+1]
	}

	message := &msgtype.Setpoint{
		Addr:      addr,
		Delay:     data.Delay,
		Loop:      data.Loop,
		Setpoints: setpoints,
	}

	if _, err := message.WriteTo(port); err != nil {
		return err
	}

	return nil
}

func dataHandler(w http.ResponseWriter, req *http.Request, body io.Reader) error {
	// http://wiki.analog.com/software/linux/docs/iio/iio
	// https://archive.fosdem.org/2012/schedule/event/693/127_iio-a-new-subsystem.pdf
	return nil
}
