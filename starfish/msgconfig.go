package starfish

import "io"

type MsgConfigMethod uint8
type MsgFeedbackType uint8

const (
	CONFIG_METHOD_NULL MsgConfigMethod = iota
	CONFIG_METHOD_QUERY
	CONFIG_METHOD_UPDATE
)

const (
	FEEDBACK_DISABLED MsgFeedbackType = iota
	FEEDBACK_OPEN_LOOP
	FEEDBACK_POSITION
	FEEDBACK_VELOCITY
	FEEDBACK_FORCE
	FEEDBACK_CURRENT
)

// PacketType = 1
type MsgConfigReq struct {
	MsgHeader
	ConfigMethod MsgConfigMethod
	FeedbackType MsgFeedbackType
	Offset       float32 // radians
	GainP        float32
	GainI        float32
	GainD        float32
	Reserved__   [20]uint8
	Checksum     int32
}

// PacketType = 3
type MsgConfigResp struct {
	MsgHeader
	ConfigMethod MsgConfigMethod
	FeedbackType MsgFeedbackType
	Offset       float32 // radians
	GainP        float32
	GainI        float32
	GainD        float32
	Version      [4]byte
	Reserved__   [16]uint8
	Checksum     int32
}

func (m *MsgConfigReq) Validate() error {
	if m.Type != MSG_CONFIG_REQ {
		return InvalidMessageError
	}
	return nil
}

func (m *MsgConfigReq) Update() error {
	m.Type = MSG_CONFIG_REQ
	return nil
}

func (m *MsgConfigReq) ReadFrom(r io.Reader) (int64, error) {
	return ReadFrom(m, r)
}

func (m *MsgConfigReq) WriteTo(w io.Writer) (int64, error) {
	return WriteTo(m, w)
}

func (m *MsgConfigReq) UnmarshalBinary(data []byte) error {
	return unmartialBinary(m, data)
}

func (m *MsgConfigReq) MarshalBinary() ([]byte, error) {
	return martialBinary(m)
}

func (m *MsgConfigResp) Validate() error {
	if m.Type != MSG_CONFIG_RESP {
		return InvalidMessageError
	}
	return nil
}

func (m *MsgConfigResp) Update() error {
	m.Type = MSG_CONFIG_RESP
	return nil
}

func (m *MsgConfigResp) ReadFrom(r io.Reader) (int64, error) {
	return ReadFrom(m, r)
}

func (m *MsgConfigResp) WriteTo(w io.Writer) (int64, error) {
	return WriteTo(m, w)
}

func (m *MsgConfigResp) UnmarshalBinary(data []byte) error {
	return unmartialBinary(m, data)
}

func (m *MsgConfigResp) MarshalBinary() ([]byte, error) {
	m.Type = MSG_CONFIG_RESP
	return martialBinary(m)
}
