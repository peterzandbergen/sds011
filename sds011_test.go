package sds011

import "testing"


func TestOpenPort(t *testing.T) {
	s := &Client{}
	s, err := Open(WithPort("usb"))
	// s, err := Open(WithPort("usb"), WithSerialSettings(nil))
	if err != nil {
		t.Errorf("error opening Client on port: %s", err.Error())
	}
	_ = s
}

