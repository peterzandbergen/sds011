package sds011

import (
	"context"
	"errors"
	"fmt"

	"go.bug.st/serial.v1"
)

type Client struct {
	p    serial.Port
	mode *serial.Mode
	port string
}

type option func(*Client) error

var (
	ErrNotImplemented = errors.New("not implemented")
)

func WithPort(port string) option {
	return func(c *Client) error {
		c.port = port
		return nil
	}
}

func WithSerialSettings(mode *serial.Mode) option {
	return func(c *Client) error {
		c.mode = mode
		return nil
	}
}

func Open(options ...option) (*Client, error) {
	c := &Client{}
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil, fmt.Errorf("open error: %w", err)
		}
	}
	// Open the port
	var err error
	p, err := serial.Open(c.port, c.mode)
	if err != nil {
		return nil, err
	}
	c.p = p
	return c, nil
}

func (c *Client) Close() {
	if c.p != nil {
		c.p.Close()
		c.p = nil
	}
}

type ReportingMode uint8
type DeviceId uint16
type WorkingMode uint8

const (
	ActiveMode = ReportingMode(iota)
	QueryMode

	DefaultDeviceId = DeviceId(0xFFFF)
)

const (
	Work = WorkingMode(iota)
	Sleep
)

func (c *Client) SetDataReportingMode(mode ReportingMode, id DeviceId) (DeviceId, error) {
	return 0, ErrNotImplemented
}

func (c *Client) GetDataReportingMode(id DeviceId) (ReportingMode, DeviceId, error) {
	return 0, 0, ErrNotImplemented
}

type Data struct {
	PM25 int
	PM10 int
}

func (c *Client) QueryData(id DeviceId) (Data, DeviceId, error) {
	return Data{}, 0, ErrNotImplemented
}

func (c *Client) SetDevideId(newId, currentId DeviceId) (DeviceId, error) {
	return 0, ErrNotImplemented
}

func (c *Client) SetWorkingMode(m WorkingMode, id DeviceId) (WorkingMode, error) {
	return 0, ErrNotImplemented
}

func (c *Client) GetWorkingMode(id DeviceId) (WorkingMode, DeviceId, error) {
	return 0, 0, ErrNotImplemented
}

func (c *Client) SetWorkingPeriod(sleepPeriod int, id DeviceId) (int, DeviceId, error) {
	return 0, 0, ErrNotImplemented
}

func (c *Client) GetWorkingPeriod(id DeviceId) (int, DeviceId, error) {
	return 0, 0, ErrNotImplemented
}

type FirmwareVersion struct {
	year, month, day int
}

func (c *Client) GetFirmwareVersion(id DeviceId) (FirmwareVersion, DeviceId, error) {
	return FirmwareVersion{}, 0, ErrNotImplemented
}

func (c *Client) readData() (Data, error) {
	return Data{}, ErrNotImplemented
}

func (c *Client) ReceiveData(ctx context.Context, sleepTime int, id DeviceId) (<-chan Data, error) {
	// Set working period and start working mode.
	var err error
	_, err = c.SetDataReportingMode(ActiveMode, DefaultDeviceId)
	if err != nil {
		return nil, err
	}
	_, _, err = c.SetWorkingPeriod(sleepTime, DefaultDeviceId)
	if err != nil {
		return nil, err
	}

	cout := make(chan Data, 10)
	go func() {
		for {
			// Listen for data messages and send to channel.
			select {
			case <-ctx.Done():
				close(cout)
				return
			default:
				d, err := c.readData()
				if err != nil {
					close(cout)
					return
				}
				cout <- d
			}
		}
	}()
	return cout, nil
}
