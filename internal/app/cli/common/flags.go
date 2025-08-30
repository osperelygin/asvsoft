// Package common предоставляет методы для добавления коммандам общих cli флагов
package common

import (
	"asvsoft/internal/pkg/communication"
	serialport "asvsoft/internal/pkg/serial-port"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

const (
	DefaultSerialPort         = "/dev/ttySC0"
	DefaultSerialPortBaudrate = 921600
	DefaultSerialPortTimeout  = 500 * time.Millisecond
)

// AddSerialDestinationFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса назначения и возвращает его конфиг. По умолчанию используется порт
// /dev/ttySOFT0 со скоростью 4800 bit/sec
func AddSerialDestinationFlags(cmd *cobra.Command) *serialport.Config {
	config := addSerialSourceFlagsWithPrefix(cmd, "dst")

	cmd.Flags().DurationVar(
		&config.Sleep, "dst-sleep",
		0, "sleep after transmitting data via serail port",
	)

	cmd.Flags().BoolVar(
		&config.Sync, "dst-sync",
		true, "wait ok message after sending own message",
	)

	cmd.Flags().IntVar(
		&config.ChunkSize, "dst-chunk-size",
		communication.DefaultChunkSize, "wait ok message after sending own message",
	)

	cmd.Flags().IntVar(
		&config.ChunkSize, "dst-retries-limit",
		communication.DefaultRetriesLimit, "wait ok message after sending own message",
	)

	cmd.Flags().BoolVar(
		&config.TransmittingDisabled, "transmitting-disabled",
		false, "disble transmitting to destination port",
	)

	return config
}

// AddSerialSourceFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса источника и возвращает его конфиг. По умолчанию используется порт
// /dev/ttyAMA0 со скоростью 4800 bit/sec и таймаутом 5 секунд .
func AddSerialSourceFlags(cmd *cobra.Command) *serialport.Config {
	return addSerialSourceFlagsWithPrefix(cmd, "")
}

func addSerialSourceFlagsWithPrefix(cmd *cobra.Command, prefix string) *serialport.Config {
	// addSerialSourceFlagsWithPrefix аналогично AddSerialSourceFlags, но с воможностью добавить флагам префикс.
	var config serialport.Config

	cmd.Flags().StringVar(
		&config.Port, strings.Trim(prefix+"-"+"port", "-"),
		DefaultSerialPort, "target port to sending measures",
	)

	cmd.Flags().IntVar(
		&config.BaudRate, strings.Trim(prefix+"-"+"baudrate", "-"),
		DefaultSerialPortBaudrate, "serial port baud rate",
	)

	cmd.Flags().DurationVar(
		&config.Timeout, strings.Trim(prefix+"-"+"timeout", "-"),
		DefaultSerialPortTimeout, "serial port timeout",
	)

	return &config
}
