// Package config ...
package config

import (
	"asvsoft/internal/pkg/communication"
	serialport "asvsoft/internal/pkg/serial-port"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type ModuleConfig struct {
	SensorSerialPort      *SerialPortConfig
	ControllerSerialPort  *SerialPortConfig
	RegistratorSerialPort *SerialPortConfig
	NeoM8t                *NeoM8tConfig
	SenseHAT              *SenseHATConfig
}

type ControllerConfig struct {
	Modules map[string]*ModuleConnectionConfig `yaml:"modules" mapstructure:"modules"`
}

type ModuleConnectionConfig struct {
	Listener *SerialPortConfig `yaml:"listener" mapstructure:"listener"`
	Enabled  bool              `yaml:"enabled" mapstructure:"enabled"`
}

type SerialPortConfig struct {
	serialport.Config
	// Sync флаг включения функционала гарантированной доставки сообщений. В случае конфига
	// сервера - будут отправляться ok-сообщения, в случае конфига клиента - будет ожидание
	// ok-сообщения от сервера.
	Sync                 bool `yaml:"sync" mapstructure:"sync"`
	ChunkSize            int  `yaml:"chunk_size" mapstructure:"chunk_size"`
	RetriesLimit         int  `yaml:"retries_limit" mapstructure:"retries_limit"`
	TransmittingDisabled bool
	Sleep                time.Duration
}

func (c *SerialPortConfig) SetDefaults() {
	if c.ChunkSize == 0 {
		c.ChunkSize = communication.DefaultChunkSize
	}

	if c.RetriesLimit == 0 {
		c.RetriesLimit = communication.DefaultRetriesLimit
	}
}

func (c SerialPortConfig) String() string {
	return fmt.Sprintf(
		"port: %q, baudrate: %d, timeout: %v, sync: %v, sleep: %v, transmitting_disabled: %v",
		c.Port, c.BaudRate, c.Timeout, c.Sync, c.Sleep, c.TransmittingDisabled,
	)
}

func (c SerialPortConfig) Short() serialport.Config {
	return serialport.Config{
		Port:     c.Port,
		BaudRate: c.BaudRate,
		Timeout:  c.Timeout,
	}
}

type NeoM8tConfig struct {
	// Rate период получения навигационного решения в секундах
	Rate int
}

type SenseHATConfig struct {
	Period time.Duration
	Mode   string
	Acc    SenseHATSensorConfig
	Gyr    SenseHATSensorConfig
	Mag    SenseHATSensorConfig
}

type SenseHATSensorConfig struct {
	Order        float32
	Range        int
	RemoveOffset bool
}

func NewControllerConfig(cfgPath string) (*ControllerConfig, error) {
	v := viper.New()

	if cfgPath == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if !strings.HasSuffix(cfgPath, ".yaml") {
		return nil, fmt.Errorf("config file type must be yaml")
	}

	v.SetConfigType("yaml")
	v.SetConfigFile(cfgPath)

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read in config: %w", err)
	}

	var cfg ControllerConfig

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	for _, c := range cfg.Modules {
		if c.Listener == nil {
			continue
		}

		c.Listener.SetDefaults()
	}

	return &cfg, nil
}
