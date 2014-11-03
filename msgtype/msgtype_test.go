package msgtype

import (
	"bytes"
	"io"
	"testing"
)

func TestPing(t *testing.T) {
	writeExpect(t, func(w io.Writer) error {
		_, err := (&Ping{3}).WriteTo(w)
		return err
	}, []byte{3, '?', 0, 0, 95, 210})
}

func TestSetPID(t *testing.T) {
	writeExpect(t, func(w io.Writer) error {
		_, err := (&SetPID{2, 1.0, 2.0, 3.0}).WriteTo(w)
		return err
	}, []byte{2, 'c', 12, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 95, 128})
}

func TestSetpoint(t *testing.T) {
	// no setpoints
	{
		var b bytes.Buffer
		setpoint := &Setpoint{4, 13, 0xffff, []SetpointValue{}}
		if _, err := setpoint.WriteTo(&b); err == nil {
			t.Fatal("Setpoint did not return an error for empty set")
		}
	}
	// one setpoint
	writeExpect(t, func(w io.Writer) error {
		_, err := (&Setpoint{4, 13, 0xffff, []SetpointValue{
			SetpointValue{Duration: 16, Setpoint: 8},
		}}).WriteTo(w)
		return err
	}, []byte{
		4, 'g', 10, 0,
		13, 0, 255, 255, 1, 0,
		16, 0, 8, 0,
		239, 11,
	})
	// three setpoint
	writeExpect(t, func(w io.Writer) error {
		_, err := (&Setpoint{4, 13, 0xffff, []SetpointValue{
			SetpointValue{Duration: 16, Setpoint: 8},
			SetpointValue{Duration: 17, Setpoint: 95},
			SetpointValue{Duration: 1000, Setpoint: 256},
		}}).WriteTo(w)
		return err
	}, []byte{
		4, 'g', 18, 0,
		13, 0, 255, 255, 3, 0,
		16, 0, 8, 0,
		17, 0, 95, 0,
		232, 3, 0, 1,
		210, 234,
	})
}

func TestRunTests(t *testing.T) {
	writeExpect(t, func(w io.Writer) error {
		_, err := (&Test{7}).WriteTo(w)
		return err
	}, []byte{7, 't', 0, 0, 41, 39})
}

func TestRequestPosition(t *testing.T) {
	writeExpect(t, func(w io.Writer) error {
		_, err := (&Value{1}).WriteTo(w)
		return err
	}, []byte{1, 'v', 0, 0, 155, 172})
}

func writeExpect(t *testing.T, f func(io.Writer) error, expect []byte) {
	var b bytes.Buffer
	if err := f(&b); err != nil {
		t.Fatal(err)
	}
	t.Logf("Bytes: %v", b.Bytes())
	if b.Len() != len(expect) {
		t.Fatalf("Expected %d bytes, got %d", len(expect), b.Len())
	}
	subject := b.Bytes()
	for i, v := range subject {
		if v != expect[i] {
			t.Fatal("Bytes did not match expected")
		}
	}
}
