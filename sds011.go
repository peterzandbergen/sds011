package sds011

import (
	"errors"

	"go.bug.st/serial.v1"
)

type SDS011 struct {
	p serial.Port
}

type option func(*SDS011) 

func WithPort(port string) option {
	return func(s *SDS011) {

	}
}

func Open(options ...option) (*SDS011, error) {
	return nil, errors.New("not implemented yet")
}
