package starfish

import "io"

type MsgType uint8
type MsgAddress uint8

const (
	MSG_CONFIG_REQ  MsgType = 1
	MSG_CONFIG_RESP MsgType = 3
	MSG_MOTION_REQ  MsgType = 0
	MSG_MOTION_RESP MsgType = 2
)

const (
	ADDR_NONE MsgAddress = iota
	ADDR_RIBS
	ADDR_HEAD_PITCH
	ADDR_HEAD_YAW
	ADDR_SPINE
	ADDR_PURR
)

type MsgHeader struct {
	Type    MsgType
	Address MsgAddress
}

func (m *MsgHeader) Validate() error {
	switch m.Type {
	case MSG_CONFIG_REQ:
	case MSG_CONFIG_RESP:
	case MSG_MOTION_REQ:
	case MSG_MOTION_RESP:
	default:
		return InvalidMessageError
	}
	switch m.Address {
	case ADDR_NONE:
	case ADDR_RIBS:
	case ADDR_HEAD_PITCH:
	case ADDR_HEAD_YAW:
	case ADDR_SPINE:
	case ADDR_PURR:
	default:
		return InvalidAddressError
	}
	return nil
}

func (m *MsgHeader) ReadFrom(r io.Reader) (int64, error) {
	return ReadFrom(m, r)
}

func (m *MsgHeader) WriteTo(w io.Writer) (int64, error) {
	return WriteTo(m, w)
}

func (m *MsgHeader) UnmartialBinary(data []byte) error {
	return unmartialBinary(m, data)
}

func (m *MsgHeader) MartialBinary() ([]byte, error) {
	return martialBinary(m)
}
