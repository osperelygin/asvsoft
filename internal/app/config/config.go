// Package config ...
package config

import (
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	serialport "asvsoft/internal/pkg/serial-port"
)

type Config struct {
	SrcSerialPort *serialport.Config
	DstSerialPort *serialport.Config
	NeoM8t        *neom8t.Config
	Imu           *sensehat.ImuConfig
}
