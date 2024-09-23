// Package serial_port предоставляет структур для работы с последовательным портом
package serial_port

import (
	"fmt"
	"time"

	"go.bug.st/serial"
)

type SerialPort struct {
	serial.Port
	Cfg *SerialPortConfig
}

type SerialPortConfig struct {
	Port                 string
	Timeout              time.Duration
	BaudRate             int
	TransmittingDisabled bool
}

func New(cfg *SerialPortConfig) (*SerialPort, error) {
	port, err := serial.Open(cfg.Port, &serial.Mode{BaudRate: cfg.BaudRate})
	if err != nil {
		return nil, fmt.Errorf("cannot open serial port '%s': %v", cfg.Port, err)
	}

	if cfg.Timeout != 0 {
		err = port.SetReadTimeout(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("cannot set read timeout: %v", err)
		}
	}

	return &SerialPort{
		Port: port,
		Cfg:  cfg,
	}, nil
}

func (sp *SerialPort) Read(p []byte) (n int, err error) {
	for n < len(p) {
		c, err := sp.Port.Read(p[n:])
		if err != nil {
			return n, err
		}

		n += c
	}

	return n, nil
}

func (sp *SerialPort) Close() error {
	if sp.Port == nil {
		return nil
	}

	err := sp.Port.Close()
	if err != nil {
		return err
	}

	return nil
}
