package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/peterzandbergen/sds011"

	"go.bug.st/serial.v1"
)

// Blocking till port is open
func openPort(ctx context.Context, name string, mode *serial.Mode) (serial.Port, error) {
	for {
		// Open the port
		port, err := serial.Open("/dev/ttyUSB0", mode)
		if err == nil {
			return port, nil
		}
		log.Printf("failed to open port %s: %s", name, err.Error())
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
			break
		}
	}
}

var (
	port = flag.String("serport", "/dev/ttyUSB0", "serial port name")
)

func logMessages(ctx context.Context, port serial.Port) error {
	p := make([]byte, 10)
	for {
		_, err := sds011.ReadResponseBytes(ctx, port, p)
		if err != nil {
			log.Printf("error reading bytes: %s", err.Error())
			return err
		}
		di, err := sds011.NewResponse(p)
		if err != nil {
			log.Printf("error creating response: %s", err.Error())
			return err
		}
		d, ok := di.(*sds011.DataResponse)
		if !ok {
			log.Printf("error casting to *sds011.Data")
		}
		log.Printf("PM25: %d, PM10: %d", d.PM25(), d.PM10())
	}
}

func main() {
	mode := &serial.Mode{
		BaudRate: 9600,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			port, err := openPort(ctx, *port, mode)
			if errors.Is(err, context.Canceled) {
				break
			}
			err = logMessages(ctx, port)
			if errors.Is(err, context.Canceled) {
				break
			}
			port.Close()
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c)
	s := <-c
	log.Printf("signal caught: %s", s.String())
}
