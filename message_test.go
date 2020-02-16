package sds011

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestBufferSize(t *testing.T) {
	var p [9]byte
	_, err := ReadResponseBytes(context.Background(), nil, p[:])
	if errors.Is(err, nil) {
		t.Errorf("expected error, got nil")
		return
	}
	if !errors.Is(err, ErrBufferTooShort) {
		t.Errorf("expected %s, received %s", ErrBufferTooShort, err)
		return
	}
}

var respValidSize = []byte{
	0xAA,
	0xC5,
	9: 0xAB,
}

func TestValidSize(t *testing.T) {
	r := bytes.NewBuffer(respValidSize)
	p := make([]byte, 13)
	n, err := ReadResponseBytes(context.Background(), r, p)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	p = p[:n]
	if n != responseLen {
		t.Fatalf("error, exptected %d chars, received %d", responseLen, n)
	}
}

