package starfish

import (
	"errors"
	"io"
	"math"
	"math/rand"
)

type ChangePoint struct {
	Duration float64 // seconds
	Angle    float64 // radians
	Omega    float64 // radians/sec
}

type MotorState struct {
	Address MsgAddress

	// configuration
	FeedbackType MsgFeedbackType
	Offset       float32 // radians
	GainP        float32
	GainI        float32
	GainD        float32

	// state
	Status MsgStatus
	ChangePoint
	Torque  float64 // Neuton-metres
	Current float64 // amperes

	// parameters
	Period float64 // seconds
	// starfish cycle parameters
	Amplitude float64 // cm
	// starfish period+velocity parameters
	currpos  float64 // cm
	prevpos  float64 // cm
	Velocity float64 // cm/sec

	// ADC0 float32 // volts
	// ADC1 float32 // volts
	// ADC2 float32 // volts
	// ADC3 float32 // volts
	// ADC4 float32 // volts
	// ADC5 float32 // volts
	// ADC6 float32 // volts
	// ADC7 float32 // volts
}

var ConfigurationUpdateError = errors.New("Unable to update configuration")

// Send a message to configure the motor and use the response to
// validate the configuration.
func (s *MotorState) Configure(port io.ReadWriter) error {
	// send config request
	req := MsgConfigReq{
		MsgHeader:    MsgHeader{Address: s.Address},
		ConfigMethod: CONFIG_METHOD_UPDATE,
		FeedbackType: s.FeedbackType,
		Offset:       s.Offset,
		GainP:        s.GainP,
		GainI:        s.GainI,
		GainD:        s.GainD,
	}
	if _, err := req.WriteTo(port); err != nil {
		return err
	}
	// receive config response
	var resp MsgConfigResp
	if _, err := resp.ReadFrom(port); err != nil {
		return err
	}
	// check config
	if resp.FeedbackType != s.FeedbackType || resp.Offset != s.Offset ||
		resp.GainP != s.GainP || resp.GainI != s.GainI || resp.GainD != s.GainD {
		return ConfigurationUpdateError
	}

	return nil
}

// Send a message to update the motor state and use the response to
// update the stored state.
func (s *MotorState) UpdateState(port io.ReadWriter) error {
	// send update request
	req := MsgMotionReq{
		MsgHeader:    MsgHeader{Address: s.Address},
		MotionMethod: MOTION_QUEUE_SETPOINT,
		Angle:        float32(s.Angle),
		Omega:        float32(s.Omega),
		Duration:     float32(s.Duration),
	}
	if _, err := req.WriteTo(port); err != nil {
		return err
	}
	// receive update response
	var resp MsgMotionResp
	if _, err := resp.ReadFrom(port); err != nil {
		return err
	}
	// update state
	s.Status = resp.Status
	s.Angle = float64(resp.Angle)
	s.Omega = float64(resp.Omega)
	s.Duration = float64(resp.Duration)
	s.Torque = float64(resp.Torque)

	return nil
}

// Advance the motor state to the next period in the StarFish cycle
// model.
func (s *MotorState) NextCycle() {
	// m := &ChangePoint{Duration: s.Duration}
	m := s
	f := 1.0 / s.Period

	if s.Angle > math.Pi {
		m.Angle = s.Amplitude * math.Cos(s.Angle)
		m.Omega = -4.0 * math.Pi * s.Amplitude * f * math.Sin(s.Angle)
		s.Angle += 4.0 * math.Pi * s.Duration * f
		if s.Angle > 2.0*math.Pi {
			s.Angle -= 2.0 * math.Pi
		}
	} else {
		m.Angle = s.Amplitude * math.Cos(s.Angle)
		m.Omega = -math.Pi * s.Amplitude * f * math.Sin(s.Angle)
		s.Angle += math.Pi * s.Duration * f
	}
}

// Advance the motor state to the next period in the StarFish
// period+velocity model.
func (s *MotorState) NextPeriod() {
	// m := &ChangePoint{Duration: s.Duration}
	m := s
	delta := s.currpos - s.prevpos

	var omega float64
	if delta != 0 {
		omega = s.Velocity * math.Pi / math.Abs(delta)
		if s.Angle < omega*s.Velocity || s.Angle <= math.Pi {
			a := math.Min(math.Pi, s.Angle)
			s.Angle += omega * delta
			m.Angle = 0.5*delta*(1-math.Cos(a)) + s.prevpos
			m.Omega = 0.5 * s.Angle * math.Sin(a)
		}
	}

	m.Angle = s.Angle
	m.Omega = 0
	s.Angle = 0
	s.currpos = math.Pi * (2*rand.Float64() - 1) / 4
	s.prevpos = s.Angle
}
