package sds011

import (
	"context"
	"errors"
	"io"
)

const (
	head       = 0xAA
	cIDRequest = 0xB4
	cIDData    = 0xC0
	cIDStatus  = 0xC5
	tail       = 0xAB

	requestLen  = 19
	responseLen = 10
)

const (
	ReportingModeRespType   = 2
	DeviceIdRespType        = 5
	WorkingModeRespType     = 6
	FirmwareVersionRespType = 7
	WorkingPeriodRespType   = 8
)

var (
	ResponseOK        = errors.New("response found")
	ErrBufferTooShort = errors.New("buffer must be at least 10 long")
	ErrBadType        = errors.New("bad type")
	ErrPortReadError  = errors.New("port read error")
)

type Request [requestLen]byte

type Response [responseLen]byte

type reportingModeRequest struct {
}

type ReportingModeResponse struct {
	Response
}

type WorkingModeResponse struct {
	Response
}

type DeviceIdResponse struct {
	Response
}

type DataResponse struct {
	Response
}

func checksum(d []byte) int {
	var sum int
	for _, n := range d {
		i := int(n)
		sum += i
	}
	return sum % 256
}

func (r *Response) ChecksumValid() bool {
	return r.Checksum() == r.CalculateChecksum()
}

func (r *Response) Checksum() int {
	return int(r[8])
}

func (r *Response) CalculateChecksum() int {
	return checksum(r[2 : len(r)-2])
}

// IsValid checks head, command ID, tail and checksum.
func (r *Response) IsValid() bool {
	return r[0] == head &&
		r[len(r)-1] == tail &&
		isResponseID(r[1]) &&
		r.ChecksumValid()
}

func (r *Request) CalculateChecksum() int {
	return checksum(r[2 : len(r)-2])
}

func (r DataResponse) PM25() float64 {
	d := r.Response[2:4]
	h := int(d[1])
	l := int(d[0])
	return (float64(h)*256.0 + float64(l)) / 10.0
}

func (r DataResponse) PM10() float64 {
	d := r.Response[4:6]
	h := int(d[1])
	l := int(d[0])
	return (float64(h)*256.0 + float64(l)) / 10.0
}

func (r *DataResponse) IsValid() bool {
	return r.Response[0] == head &&
		r.Response[len(r.Response)-1] == tail &&
		r.Response[1] == cIDData &&
		r.Response.ChecksumValid()
}

func isResponseID(b byte) bool {
	return b == cIDData || b == cIDStatus
}


func newResponse(b []byte) (Response, error) {
	if len(b) < responseLen {
		return Response{}, ErrBufferTooShort
	}
	r := Response{}
	copy(r[:], b)
	return r, nil
}

func NewResponse(p []byte) (interface{}, error) {
	if len(p) < responseLen {
		return nil, ErrBufferTooShort
	}
	var r Response
	r, err := newResponse(p)
	if err != nil {
		return nil, err
	}

	switch r[1] {
	case cIDData:
		return &DataResponse{
			Response: r,
		}, nil
	case cIDStatus:
		switch r[2] {
		case ReportingModeRespType:
			return &ReportingModeResponse{
				Response: r,
			}, nil
		case DeviceIdRespType:
			return nil, ErrNotImplemented
		case WorkingModeRespType:
			return nil, ErrNotImplemented
		case FirmwareVersionRespType:
			return nil, ErrNotImplemented
		case WorkingPeriodRespType:
			return nil, ErrNotImplemented
		}
	}
	return nil, ErrBadType
}

// ReadResponseBytes reads a response from r.
// It scans for a start and command byte and reads the following 8 bytes.
// All 10 bytes are returned in b.
// b needs to be at lease 10 bytes long or an error will be returned.
// n will always be 10 for now.
func ReadResponseBytes(ctx context.Context, r io.Reader, b []byte) (n int, err error) {
	// read treats a zero bytes read as an error
	read := func(n int, err error) (int, error) {
		if err != nil {
			return 0, err
		}
		if n == 0 {
			return 0, ErrPortReadError
		}
		return n, err
	}

	if len(b) < responseLen {
		return 0, ErrBufferTooShort
	}
	// Read until head found or error occured.
	for {
		var p1 [1]byte
		for {
			_, err := read(r.Read(p1[:]))
			if err != nil {
				return 0, err
			}
			if p1[0] == head {
				b[0] = head
				break
			}
		}
		// Read next byte
		_, err := read(r.Read(p1[:]))
		if err != nil {
			return 0, err
		}
		if isResponseID(p1[0]) {
			b[1] = p1[0]
			break
		}
	}
	// head and commandID read.
	i := 2
	for i < 10 {
		n, err := read(r.Read(b[i:responseLen]))
		if err != nil {
			return i + n, err
		}
		i += n
	}
	return i, nil
}
