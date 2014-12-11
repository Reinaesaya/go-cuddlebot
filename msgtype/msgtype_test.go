package msgtype

import (
	"bytes"
	"encoding"
	"io"
	"testing"
)

func TestRemoteAddressMarshalText(t *testing.T) {
	addr := InvalidAddress
	if _, err := addr.MarshalText(); err == nil {
		t.Fatal("RemoteAddress.MarshalText() did not error on invalid address")
	}
	testRemoteAddress(t, RibsAddress, "ribs", "RibsAddress")
	testRemoteAddress(t, PurrAddress, "purr", "PurrAddress")
	testRemoteAddress(t, SpineAddress, "spine", "SpineAddress")
	testRemoteAddress(t, HeadXAddress, "headx", "HeadXAddress")
	testRemoteAddress(t, HeadYAddress, "heady", "HeadYAddress")
}

func testRemoteAddress(t *testing.T, addr RemoteAddress, str, name string) {
	if bs, err := addr.MarshalText(); err != nil {
		t.Fatal(err)
	} else if s, err := bytes.NewBuffer(bs).ReadString(0); err != io.EOF && err != nil {
		t.Fatal(err)
	} else if s != str {
		t.Fatalf("%s != '%s'", name, str)
	}
}

func TestPing(t *testing.T) {
	testMarshalExpect(t, &Ping{'r'}, []byte{'r', '?', 0, 0, 6, 51})
}

func TestSetPID(t *testing.T) {
	testMarshalExpect(t, &SetPID{2, 1.0, 2.0, 3.0},
		[]byte{2, 'c', 12, 0, 0, 0, 128, 63, 0, 0, 0, 64, 0, 0, 64, 64, 93, 23})
}

func TestSetpoint(t *testing.T) {
	// no setpoints
	{
		setpoint := &Setpoint{4, 13, 0xffff, []SetpointValue{}}
		if _, err := setpoint.MarshalBinary(); err == nil {
			t.Fatal("Setpoint did not return an error for empty set")
		}
	}
	// one setpoint
	testMarshalExpect(t, &Setpoint{4, 13, 0xffff, []SetpointValue{
		SetpointValue{Duration: 16, Setpoint: 8},
	}}, []byte{
		4, 'g', 10, 0,
		13, 0, 255, 255, 1, 0,
		16, 0, 8, 0,
		150, 136,
	})
	// three setpoint
	testMarshalExpect(t, &Setpoint{4, 13, 0xffff, []SetpointValue{
		SetpointValue{Duration: 16, Setpoint: 8},
		SetpointValue{Duration: 17, Setpoint: 95},
		SetpointValue{Duration: 1000, Setpoint: 256},
	}}, []byte{
		4, 'g', 18, 0,
		13, 0, 255, 255, 3, 0,
		16, 0, 8, 0,
		17, 0, 95, 0,
		232, 3, 0, 1,
		144, 36,
	})
}

func TestSleep(t *testing.T) {
	testMarshalExpect(t, &Sleep{9},
		[]byte{9, 'z', 0, 0, 181, 68})
}

func TestRunTests(t *testing.T) {
	testMarshalExpect(t, &Test{7}, []byte{7, 't', 0, 0, 241, 215})
}

func TestRequestPosition(t *testing.T) {
	testMarshalExpect(t, &Value{1}, []byte{1, 'v', 0, 0, 70, 158})
}

func testMarshalExpect(t *testing.T, m encoding.BinaryMarshaler, expect []byte) {
	subject, err := m.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Bytes: %v", subject)
	if len(subject) != len(expect) {
		t.Fatalf("Expected %d bytes, got %d", len(expect), len(subject))
	}
	for i, v := range subject {
		if v != expect[i] {
			t.Fatal("Bytes did not match expected")
		}
	}
}
