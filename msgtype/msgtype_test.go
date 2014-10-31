package msgtype

import (
	"bytes"
	"testing"
)

func TestPing(t *testing.T) {
	writeExpect(t, func(w RPC) error {
		return w.Ping(3)
	}, []byte{3, '?', 0, 0, 95, 210})
}

func TestSetPID(t *testing.T) {
	writeExpect(t, func(w RPC) error {
		return w.SetPID(2, 1.0, 2.0, 3.0)
	}, []byte{2, 'c', 12, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 95, 128})
}

func TestSetpoint(t *testing.T) {
	// no setpoints
	{
		var b bytes.Buffer
		w := RPC{Writer: &b}
		if err := w.Setpoint(4, 13, 0xffff, []Setpoint{}); err == nil {
			t.Fatal("Setpoint did not return an error for empty set")
		}
	}
	// one setpoint
	writeExpect(t, func(w RPC) error {
		return w.Setpoint(4, 13, 0xffff, []Setpoint{
			Setpoint{Duration: 16, Setpoint: 8},
		})
	}, []byte{
		4, 'g', 10, 0,
		13, 0, 255, 255, 1, 0,
		16, 0, 8, 0,
		239, 11,
	})
	// three setpoint
	writeExpect(t, func(w RPC) error {
		return w.Setpoint(4, 13, 0xffff, []Setpoint{
			Setpoint{Duration: 16, Setpoint: 8},
			Setpoint{Duration: 17, Setpoint: 95},
			Setpoint{Duration: 1000, Setpoint: 256},
		})
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
	writeExpect(t, func(w RPC) error {
		return w.RunTests(7)
	}, []byte{7, 't', 0, 0, 41, 39})
}

func TestRequestPosition(t *testing.T) {
	writeExpect(t, func(w RPC) error {
		return w.RequestPosition(1)
	}, []byte{1, 'v', 0, 0, 155, 172})
}

func writeExpect(t *testing.T, f func(RPC) error, expect []byte) {
	var b bytes.Buffer
	w := RPC{Writer: &b}
	if err := f(w); err != nil {
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
