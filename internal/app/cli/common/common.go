package common

import (
	"asvsoft/internal/pkg/serial_wrapper"
	"fmt"
	"io"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

const (
	defaultBaudrate = 4800
	defaultTimeout  = 5 * time.Second
)

type SerialConfig struct {
	Port     string
	Timeout  time.Duration
	BaudRate int
}

// AddSerialDestinationFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса назначения и возвращает его конфиг. По умолчанию используется порт
// /dev/ttySOFT0 со скоростью 4800 bit/sec
func AddSerialDestinationFlags(cmd *cobra.Command) *SerialConfig {
	var config SerialConfig

	cmd.Flags().StringVar(
		&config.Port, "dst",
		"/dev/ttySOFT0", "target port to sending measures",
	)

	cmd.Flags().IntVar(
		&config.BaudRate, "dst-baudrate",
		defaultBaudrate, "serial port baud rate",
	)

	return &config
}

// AddSerialSourceFlags добавляем команде флаги последовательного конфигурации
// последовательного интерфейса источника и возвращает его конфиг. По умолчанию используется порт
// /dev/ttyAMA0 со скоростью 4800 bit/sec и таймаутом 5 секунд .
func AddSerialSourceFlags(cmd *cobra.Command, prefix string) *SerialConfig {
	var config SerialConfig

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

func InitSerialPort(cfg *SerialConfig) (serial.Port, error) {
	srcPort, err := serial.Open(cfg.Port, &serial.Mode{BaudRate: cfg.BaudRate})
	if err != nil {
		return nil, fmt.Errorf("cannot open serial port '%s': %v", cfg.Port, err)
	}

	if cfg.Timeout != 0 {
		err = srcPort.SetReadTimeout(cfg.Timeout)
		if err != nil {
			return nil, fmt.Errorf("cannot set read timeout: %v", err)
		}
	}

	return serial_wrapper.New(srcPort), nil
}

func CloseSerialPort(port io.Closer) {
	err := port.Close()
	if err != nil {
		log.Errorf("cannot close serial port: %v", err)
	}
}
