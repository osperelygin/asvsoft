package serial_wrapper

import (
	"go.bug.st/serial"
)

type SerialPortWrapper struct {
	serial.Port
}

func New(port serial.Port) serial.Port {
	return &SerialPortWrapper{Port: port}
}

func (r *SerialPortWrapper) Read(p []byte) (n int, err error) {
	for n < len(p) {
		c, err := r.Port.Read(p[n:])
		if err != nil {
			return n, err
		}

		n += c
	}

	return n, nil
}
