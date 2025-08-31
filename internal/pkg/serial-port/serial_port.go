// Package serialport предоставляет структур для работы с последовательным портом
package serialport

import (
	"asvsoft/internal/pkg/logger"
	"errors"
	"fmt"
	"time"

	"go.bug.st/serial"
)

var (
	ErrReadTimeout = errors.New("read timeout")
)

type Wrapper struct {
	serial.Port
	logger logger.Logger
	Cfg    Config
}

type Config struct {
	Port     string        `yaml:"port" mapstructure:"port"`
	BaudRate int           `yaml:"baudrate" mapstructure:"baudrate"`
	Timeout  time.Duration `yaml:"timeout" mapstructure:"timeout"`
}

func New(cfg Config) (*Wrapper, error) {
	port, err := newSerialPort(cfg)
	if err != nil {
		return nil, err
	}

	return &Wrapper{
		Port: port,
		Cfg:  cfg,
	}, nil
}

func newSerialPort(cfg Config) (serial.Port, error) {
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

	return port, nil
}

func (w *Wrapper) SetLogger(logger logger.Logger) *Wrapper {
	w.logger = logger
	return w
}

func (w *Wrapper) Logger() logger.Logger {
	if w.logger == nil {
		return logger.DummyLogger{}
	}

	return w.logger
}

func (w *Wrapper) Read(p []byte) (n int, err error) {
	for n < len(p) {
		c, err := w.Port.Read(p[n:])
		if err != nil {
			err = w.portClosedFallback(err)
			if err != nil {
				return n, err
			}

			n = 0

			continue
		}

		if c == 0 {
			return 0, ErrReadTimeout
		}

		n += c
	}

	return n, nil
}

func (w *Wrapper) portClosedFallback(err error) error {
	pErr := new(serial.PortError)
	if errors.As(err, &pErr) && pErr.Code() == serial.PortClosed {
		w.Port, err = newSerialPort(w.Cfg)
		if err != nil {
			return fmt.Errorf("port closed and failed to reopen: %w", err)
		}

		w.Logger().Warnf("serail port was reopened")

		return nil
	}

	return err
}

func (w *Wrapper) Close() error {
	if w.Port == nil {
		return nil
	}

	err := w.Port.Close()
	if err != nil {
		return err
	}

	return nil
}
