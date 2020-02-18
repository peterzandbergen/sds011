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

var dataResponse = []byte{
	0xAA,
	0xC0,
	0xD4, // PM2.5 low
	0x04, // PM2.5 high
	0x3A, // PM10 low
	0x0A, // PM10 hig
	0xA1, // Device ID 1
	0x60, // Device ID 2
	0x1D, // Checksum
	0xAB,
}

func TestDataResponse(t *testing.T) {
	r := bytes.NewBuffer(dataResponse)
	p := make([]byte, 10)
	n, err := ReadResponseBytes(context.Background(), r, p)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	p = p[:n]
	if n != responseLen {
		t.Fatalf("error, exptected %d chars, received %d", responseLen, n)
	}
	di, err := NewResponse(p)
	d := di.(*DataResponse)
	if ! d.ChecksumValid() {
		t.Fatal("checksum invald")
	}
	t.Logf("PM25 %f", d.PM25())
	t.Logf("PM10 %f", d.PM10())
}


func TestDataResponseIsValid(t *testing.T) {
	r := bytes.NewBuffer(dataResponse)
	p := make([]byte, 10)
	n, err := ReadResponseBytes(context.Background(), r, p)
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}
	p = p[:n]
	if n != responseLen {
		t.Fatalf("error, exptected %d chars, received %d", responseLen, n)
	}
	di, err := NewResponse(p)
	d := di.(*DataResponse)
	if ! d.IsValid() {
		t.Fatal("response invald")
	}
	t.Logf("PM25 %f", d.PM25())
	t.Logf("PM10 %f", d.PM10())
}
