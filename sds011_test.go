package sds011

import "testing"


func TestOpenPort(t *testing.T) {
	s := &SDS011{}
	s, err := Open(WithPort("usb"))
	if err != nil {
		t.Errorf("error opening SDS011: %s", err.Error())
	}
	_ = s
}

