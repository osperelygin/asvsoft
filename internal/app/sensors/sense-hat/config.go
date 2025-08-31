// Package sensehat предоставляет функционал для чтения и конфигурации измерения sense hat (c)
package sensehat

import (
	"asvsoft/internal/app/config"
	"fmt"
)

const (
	IntertialMode = "inertial"
	FullMode      = "full"
)

type internalConfig struct {
	config.SenseHATConfig
	Acc sensorConfig
	Gyr sensorConfig
	Mag sensorConfig
}

type sensorConfig struct {
	c        config.SenseHATSensorConfig
	enable   bool
	orderMap map[float32]byte
	rangeMap map[int]rangeConfig
}

type rangeConfig struct {
	value       byte
	sensitivity int
}

func getInternalConfig(cmnCfg config.SenseHATConfig) internalConfig {
	var cfg internalConfig

	cfg.Period = cmnCfg.Period
	cfg.Mode = cmnCfg.Mode

	cfg.Acc = sensorConfig{
		c: cmnCfg.Acc,
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
	}

	cfg.Gyr = sensorConfig{
		c: cmnCfg.Gyr,
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
	}

	cfg.Mag = sensorConfig{
		c: cmnCfg.Mag,
		orderMap: map[float32]byte{
			10:  AK09918_CONTINUOUS_10HZ,
			20:  AK09918_CONTINUOUS_20HZ,
			50:  AK09918_CONTINUOUS_50HZ,
			100: AK09918_CONTINUOUS_100HZ,
		},
		rangeMap: map[int]rangeConfig{
			0: rangeConfig{},
		},
	}

	return cfg
}

func (c *internalConfig) validate() (err error) {
	switch c.Mode {
	case IntertialMode:
		c.Acc.enable = true
		c.Gyr.enable = true
	case FullMode:
		c.Acc.enable = true
		c.Gyr.enable = true
		c.Mag.enable = true
	default:
		return fmt.Errorf("unknown mode: '%s'", c.Mode)
	}

	if !(c.Acc.enable || c.Gyr.enable || c.Mag.enable) {
		return fmt.Errorf("acc, gyro and magn disable, nothing to do")
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

func (c *sensorConfig) validate() error {
	if !c.enable {
		return nil
	}

	if _, ok := c.orderMap[c.c.Order]; !ok {
		return fmt.Errorf("unknown sensor order: %f", c.c.Order)
	}

	if _, ok := c.rangeMap[c.c.Range]; !ok {
		return fmt.Errorf("unknown sensor range: %d", c.c.Range)
	}

	return nil
}

func (c *sensorConfig) rangeValue() byte {
	return c.rangeMap[c.c.Range].value
}

func (c *sensorConfig) rangeSensitivity() int {
	return c.rangeMap[c.c.Range].sensitivity
}

func (c *sensorConfig) order() byte {
	return c.orderMap[c.c.Order]
}
