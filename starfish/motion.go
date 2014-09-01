package starfish

import (
	"math"
	"math/rand"
)

type ChangePoint struct {
	Duration float64 // seconds
	Angle    float64 // radians
	Omega    float64 // radians/sec
}

type MotorState struct {
	Period  float64 // seconds
	Torque  float64 // Neuton-metres
	Current float64 // amperes
	ChangePoint
}

type CycleState struct {
	MotorState
	Amplitude float64 // cm
}

type PeriodState struct {
	MotorState
	currpos  float64 // cm
	prevpos  float64 // cm
	Velocity float64 // cm/sec
}

func (s *CycleState) NextPoint() *ChangePoint {
	m := &ChangePoint{Duration: s.Duration}
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

	return m
}

func (s *PeriodState) NextPoint() *ChangePoint {
	m := &ChangePoint{Duration: s.Duration}
	delta := s.currpos - s.prevpos

	var omega float64
	if delta != 0 {
		omega = s.Velocity * math.Pi / math.Abs(delta)
		if s.Angle < omega*s.Velocity || s.Angle <= math.Pi {
			a := math.Min(math.Pi, s.Angle)
			s.Angle += omega * delta
			m.Angle = 0.5*delta*(1-math.Cos(a)) + s.prevpos
			m.Omega = 0.5 * s.Angle * math.Sin(a)
			return m
		}
	}

	m.Angle = s.Angle
	m.Omega = 0
	s.Angle = 0
	s.currpos = math.Pi * (2*rand.Float64() - 1) / 4
	s.prevpos = s.Angle

	return m
}
