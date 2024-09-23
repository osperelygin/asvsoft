// Package common предоставляет методы для добавления коммандам общих cli флагов
package common

import (
	"asvsoft/internal/pkg/serial_port"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// LogLevel уровень логирования
var LogLevel string

const (
	defaultBaudrate = 4800
	defaultTimeout  = 5 * time.Second
)

// AddSerialDestinationFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса назначения и возвращает его конфиг. По умолчанию используется порт
// /dev/ttySOFT0 со скоростью 4800 bit/sec
func AddSerialDestinationFlags(cmd *cobra.Command) *serial_port.SerialPortConfig {
	var config serial_port.SerialPortConfig

	cmd.Flags().StringVar(
		&config.Port, "dst",
		"/dev/ttySOFT0", "target port to sending measures",
	)

	cmd.Flags().IntVar(
		&config.BaudRate, "dst-baudrate",
		defaultBaudrate, "serial port baud rate",
	)

	cmd.Flags().BoolVar(
		&config.TransmittingDisabled, "transmitting-disabled",
		false, "disble transmitting to destination port",
	)

	return &config
}

// AddSerialSourceFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса источника и возвращает его конфиг. По умолчанию используется порт
// /dev/ttyAMA0 со скоростью 4800 bit/sec и таймаутом 5 секунд .
func AddSerialSourceFlags(cmd *cobra.Command) *serial_port.SerialPortConfig {
	return AddSerialSourceFlagsWithPrefix(cmd, "")
}

// AddSerialSourceFlagsWithPrefix аналогично AddSerialSourceFlags, но с воможностью добавить флагам префикс.
func AddSerialSourceFlagsWithPrefix(cmd *cobra.Command, prefix string) *serial_port.SerialPortConfig {
	var config serial_port.SerialPortConfig

	cmd.Flags().StringVar(
		&config.Port, strings.Trim(prefix+"-"+"port", "-"),
		"/dev/ttyAMA0", "target port to sending measures",
	)

	cmd.Flags().IntVar(
		&config.BaudRate, strings.Trim(prefix+"-"+"baudrate", "-"),
		defaultBaudrate, "serial port baud rate",
	)

	cmd.Flags().DurationVar(
		&config.Timeout, strings.Trim(prefix+"-"+"timeout", "-"),
		defaultTimeout, "serial port timeout",
	)

	return &config
}
