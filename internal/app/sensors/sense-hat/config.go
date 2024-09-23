// Package sensehat предоставляет функционал для чтения и конфигурации измерения sense hat (c)
package sensehat

import (
	"fmt"
	"time"
)

const (
	IntertialMode = "inertial"
	FullMode      = "full"
)

type ConfigError struct {
	reason string
}

func (err *ConfigError) Error() string {
	return fmt.Sprintf("IMU config error was occured: %s", err.reason)
}

type ImuConfig struct {
	Period time.Duration
	Mode   string
	Acc    SensorConfig
	Gyr    SensorConfig
	Mag    SensorConfig
}

type SensorConfig struct {
	Enable   bool
	Order    float32
	Range    int
	orderMap map[float32]byte
	rangeMap map[int]rangeConfig
}

type rangeConfig struct {
	value       byte
	sensitivity int
}

func (c *SensorConfig) validate() error {
	if !c.Enable {
		return nil
	}

	if _, ok := c.orderMap[c.Order]; !ok {
		return &ConfigError{fmt.Sprintf("unknown sensor order: %f", c.Order)}
	}

	if _, ok := c.rangeMap[c.Range]; !ok {
		return &ConfigError{fmt.Sprintf("unknown sensor range: %d", c.Range)}
	}

	return nil
}

func (c *SensorConfig) rangeValue() byte {
	return c.rangeMap[c.Range].value
}

func (c *SensorConfig) rangeSensitivity() int {
	return c.rangeMap[c.Range].sensitivity
}

func (c *SensorConfig) order() byte {
	return c.orderMap[c.Order]
}

func NewImuConfig() *ImuConfig {
	return &ImuConfig{
		Acc: SensorConfig{
			orderMap: map[float32]byte{
				8000:  QMI8658AccOdr_8000Hz,
				4000:  QMI8658AccOdr_4000Hz,
				2000:  QMI8658AccOdr_2000Hz,
				1000:  QMI8658AccOdr_1000Hz,
				500:   QMI8658AccOdr_500Hz,
				250:   QMI8658AccOdr_250Hz,
				125:   QMI8658AccOdr_125Hz,
				62.5:  QMI8658AccOdr_62_5Hz,
				31.25: QMI8658AccOdr_31_25Hz,
				128:   QMI8658AccOdr_LowPower_128Hz,
				21:    QMI8658AccOdr_LowPower_21Hz,
				11:    QMI8658AccOdr_LowPower_11Hz,
				3:     QMI8658AccOdr_LowPower_3Hz,
			},
			rangeMap: map[int]rangeConfig{
				2:  {QMI8658AccRange2g, 1 << 14},
				4:  {QMI8658AccRange4g, 1 << 13},
				8:  {QMI8658AccRange8g, 1 << 12},
				16: {QMI8658AccRange16g, 1 << 11},
			},
		},
		Gyr: SensorConfig{
			orderMap: map[float32]byte{
				8000:  QMI8658GyrOdr8000Hz,
				4000:  QMI8658GyrOdr4000Hz,
				2000:  QMI8658GyrOdr2000Hz,
				1000:  QMI8658GyrOdr1000Hz,
				500:   QMI8658GyrOdr500Hz,
				250:   QMI8658GyrOdr250Hz,
				125:   QMI8658GyrOdr125Hz,
				62.5:  QMI8658GyrOdr62_5Hz,
				31.25: QMI8658GyrOdr31_25Hz,
			},
			rangeMap: map[int]rangeConfig{
				16:   {QMI8658GyrRange16dps, 1 << 11},
				32:   {QMI8658GyrRange32dps, 1 << 10},
				64:   {QMI8658GyrRange64dps, 1 << 9},
				128:  {QMI8658GyrRange128dps, 1 << 8},
				256:  {QMI8658GyrRange256dps, 1 << 7},
				512:  {QMI8658GyrRange512dps, 1 << 6},
				1024: {QMI8658GyrRange1024dps, 1 << 5},
				2048: {QMI8658GyrRange2048dps, 1 << 4},
			},
		},
		Mag: SensorConfig{
			orderMap: map[float32]byte{
				10:  AK09918_CONTINUOUS_10HZ,
				20:  AK09918_CONTINUOUS_20HZ,
				50:  AK09918_CONTINUOUS_50HZ,
				100: AK09918_CONTINUOUS_100HZ,
			},
			rangeMap: map[int]rangeConfig{
				0: rangeConfig{},
			},
		},
	}
}

func (c *ImuConfig) validate() (err error) {
	switch c.Mode {
	case IntertialMode:
		c.Acc.Enable = true
		c.Gyr.Enable = true
	case FullMode:
		c.Acc.Enable = true
		c.Gyr.Enable = true
		c.Mag.Enable = true
	default:
		return &ConfigError{reason: fmt.Sprintf("unknown mode: '%s'", c.Mode)}
	}

	if !(c.Acc.Enable || c.Gyr.Enable || c.Mag.Enable) {
		return &ConfigError{"acc, gyro and magn disable, nothing to do"}
	}

	err = c.Acc.validate()
	if err != nil {
		return err
	}

	err = c.Gyr.validate()
	if err != nil {
		return err
	}

	err = c.Mag.validate()
	if err != nil {
		return err
	}

	return nil
}
