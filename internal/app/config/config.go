// Package config ...
package config

import (
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	serialport "asvsoft/internal/pkg/serial-port"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type ModuleConfig struct {
	SensorSerialPort      *serialport.Config
	ControllerSerialPort  *serialport.Config
	RegestratorSerialPort *serialport.Config
	NeoM8t                *neom8t.Config
	Imu                   *sensehat.ImuConfig
}

type ControllerConfig struct {
	Modules map[string]ModuleConnectionConfig `yaml:"modules" mapstructure:"modules"`
}

type ModuleConnectionConfig struct {
	Listener *serialport.Config `yaml:"listener" mapstructure:"listener"`
	Enabled  bool               `yaml:"enabled" mapstructure:"enabled"`
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
	v.WatchConfig()

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read in config: %w", err)
	}

	var cfg ControllerConfig

	err = v.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
