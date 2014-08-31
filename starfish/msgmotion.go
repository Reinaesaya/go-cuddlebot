package starfish

type MsgMotionMethod uint8
type MsgStatus uint8

const (
	MOTION_STOP MsgMotionMethod = iota
	MOTION_STATUS_REQ
	MOTION_QUEUE_SETPOINT
	MOTION_QUEUE_NEXT
)

const (
	STATUS_STOP MsgStatus = iota
	STATUS_IDLE
	STATUS_MOVING
	STATUS_ERROR
)

// PacketType = 0
type MsgMotionReq struct {
	MsgHeader
	MotionMethod MsgMotionMethod
	Reserved0__  uint8
	Position     float32 // radians
	Velocity     float32 // radians/sec
	Duration     float32 // seconds
	Reserved1__  uint8
	SkinDVal     uint8
	Reserved2__  uint8
	SkinMode     uint8
	Reserved3__  [8]uint8
	Checksum     int32
}

// PacketType = 2
type MsgMotionResp struct {
	MsgHeader
	Status      MsgStatus
	Reserved0__ uint8
	Position    float32 // radians
	Velocity    float32 // radians/sec
	Duration    float32 // seconds
	LoopIndex   uint8
	Reserved1__ [3]uint8
	Torque      float32 // Neuton-metres
	Current     float32 // amperes
	ADC0        float32 // volts
	ADC1        float32 // volts
	ADC2        float32 // volts
	ADC3        float32 // volts
	ADC4        float32 // volts
	ADC5        float32 // volts
	ADC6        float32 // volts
	ADC7        float32 // volts
	Checksum    int32
}

func (m *MsgMotionReq) Validate() error {
	if m.Type != MSG_MOTION_REQ {
		return InvalidMessageError
	}
	return nil
}

func (m *MsgMotionReq) Update() error {
	m.Type = MSG_MOTION_REQ
	return nil
}

func (m *MsgMotionReq) UnmarshalBinary(data []byte) error {
	return unmartialBinary(m, data)
}

func (m *MsgMotionReq) MarshalBinary() ([]byte, error) {
	return martialBinary(m)
}

func (m *MsgMotionResp) Validate() error {
	if m.Type != MSG_MOTION_RESP {
		return InvalidMessageError
	}
	return nil
}

func (m *MsgMotionResp) Update() error {
	m.Type = MSG_MOTION_RESP
	return nil
}

func (m *MsgMotionResp) UnmarshalBinary(data []byte) error {
	return unmartialBinary(m, data)
}

func (m *MsgMotionResp) MarshalBinary() ([]byte, error) {
	return martialBinary(m)
}
