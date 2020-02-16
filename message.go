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

func (r *Request) CalculateChecksum() int {
	return checksum(r[2 : len(r)-2])
}

func (r *Response) CalculateChecksum() int {
	return checksum(r[2 : len(r)-2])
}

func (r DataResponse) PM25() int {
	d := r.Response[2:4]
	h := int(d[1])
	l := int(d[0])
	return h*256+l
}

func (r DataResponse) PM10() int {
	d := r.Response[4:6]
	h := int(d[1])
	l := int(d[0])
	return h*256+l
}

func isResponseID(b byte) bool {
	return b == cIDData || b == cIDStatus
}

func isDataID(b byte) bool {
	return b == cIDData
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
	if len(b) < responseLen {
		return 0, ErrBufferTooShort
	}
	// Read until head found or error occured.
	for {
		var p1 [1]byte
		for {
			n, err := r.Read(p1[:])
			if err != nil || n == 0 {
				return n, err
			}
			if p1[0] == head {
				b[0] = head
				break
			}
		}
		// Read next byte
		n, err := r.Read(p1[:])
		if err != nil || n == 0 {
			return n, err
		}
		if isResponseID(p1[0]) {
			b[1] = p1[0]
			break
		}
	}
	// head and commandID read.
	i := 2
	for i < 10 {
		n, err := r.Read(b[i:responseLen])
		if err != nil {
			return i + n, err
		}
		i += n
	}
	return i, nil
}
