package cuddlespeak

import (
	"."
	"bytes"
	"testing"
)

func TestWritePing(t *testing.T) {
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WritePing(3)
	}, []byte{3, 1, 0, 0, 0, 0})
}

func TestWritePong(t *testing.T) {
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WritePong(5)
	}, []byte{5, 2, 0, 0, 0, 0})
}

func TestWriteSetPID(t *testing.T) {
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteSetPID(2, 1.0, 2.0, 3.0)
	}, []byte{2, 3, 12, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 177, 150})
}

func TestWriteSetpoint(t *testing.T) {
	// no setpoints
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteSetpoint(4, 13, 0xffff, []cuddlespeak.Setpoint{})
	}, []byte{
		4, 4, 6, 0,
		13, 0, 255, 255, 0, 0,
		97, 128,
	})
	// one setpoint
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteSetpoint(4, 13, 0xffff, []cuddlespeak.Setpoint{
			cuddlespeak.Setpoint{Duration: 16, Setpoint: 8},
		})
	}, []byte{
		4, 4, 10, 0,
		13, 0, 255, 255, 1, 0,
		16, 0, 8, 0,
		34, 161,
	})
	// three setpoint
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteSetpoint(4, 13, 0xffff, []cuddlespeak.Setpoint{
			cuddlespeak.Setpoint{Duration: 16, Setpoint: 8},
			cuddlespeak.Setpoint{Duration: 17, Setpoint: 95},
			cuddlespeak.Setpoint{Duration: 1000, Setpoint: 256},
		})
	}, []byte{
		4, 4, 18, 0,
		13, 0, 255, 255, 3, 0,
		16, 0, 8, 0,
		17, 0, 95, 0,
		232, 3, 0, 1,
		167, 52,
	})
}

func TestWriteTest(t *testing.T) {
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteTest(7)
	}, []byte{7, 5, 0, 0, 0, 0})
}

func TestWriteValue(t *testing.T) {
	writeExpect(t, func(w cuddlespeak.Writer) error {
		return w.WriteValue(1)
	}, []byte{1, 6, 0, 0, 0, 0})
}

func writeExpect(t *testing.T, f func(cuddlespeak.Writer) error, expect []byte) {
	var b bytes.Buffer
	w := cuddlespeak.Writer{Writer: &b}
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
