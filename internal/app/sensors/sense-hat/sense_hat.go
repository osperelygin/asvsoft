package sensehat

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/encoder"
	"asvsoft/internal/pkg/proto"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/d2r2/go-i2c"
)

const (
	offsetCalculatingTries = 1024
	offsetCalculatingSleep = 16 * time.Millisecond
	wrongMeasureThreshold  = 3 << 13
)

type configCmd struct {
	reg byte
	val byte
}

type SenseHAT struct {
	buf         []byte
	config      internalConfig
	inertialBus *i2c.I2C
	magnBus     *i2c.I2C
	gxOffset    int16
	gyOffset    int16
	gzOffset    int16
}

func New(cmnCfg *config.SenseHATConfig) (*SenseHAT, error) {
	cfg := getInternalConfig(*cmnCfg)

	err := cfg.validate()
	if err != nil {
		return nil, err
	}

	s := &SenseHAT{config: cfg}

	switch cfg.Mode {
	case FullMode:
		s.buf = make([]byte, 24)
	case IntertialMode:
		s.buf = make([]byte, 16)
	default:
		return nil, fmt.Errorf("cannot create imu: unknown mode: '%s'", cfg.Mode)
	}

	defer func() {
		if err != nil {
			err = s.Close()
			if err != nil {
				log.Errorf("failed to close imu: %v", err)
			}
		}
	}()

	s.inertialBus, err = initInertialSensors(cfg)
	if err != nil {
		return nil, err
	}

	if cfg.Gyr.RemoveOffset {
		var gx, gy, gz int

		for i := 0; i < offsetCalculatingTries; i++ {
			m, err := s.measure()
			if err != nil {
				return nil, fmt.Errorf("cannot remove offset: %w", err)
			}

			gx += int(m.Gx)
			gy += int(m.Gy)
			gz += int(m.Gz)

			time.Sleep(offsetCalculatingSleep)
		}

		s.gxOffset = int16(gx / offsetCalculatingTries)
		s.gyOffset = int16(gy / offsetCalculatingTries)
		s.gzOffset = int16(gz / offsetCalculatingTries)

		log.Infof("gyro offset: x=%d, y=%d, z=%d", s.gxOffset, s.gyOffset, s.gzOffset)
	}

	s.magnBus, err = initMagnSensor(cfg)
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *SenseHAT) Measure(_ context.Context) (proto.Packer, error) {
	time.Sleep(s.config.Period)
	return s.measure()
}

// ErrWrongMeasure ...
var ErrWrongMeasure = errors.New("wrong measure")

func (s *SenseHAT) measure() (*proto.IMUData, error) {
	b, err := s.RawMeasure()
	if err != nil {
		return nil, err
	}

	m := &proto.IMUData{}
	decoder := encoder.NewDecoder(io.NopCloser(bytes.NewBuffer(b)))

	switch s.config.Mode {
	case FullMode:
		m.AccFactor = int16(s.config.Acc.rangeSensitivity())
		m.GyrFactor = int16(s.config.Gyr.rangeSensitivity())
		err = decoder.Decode(
			&m.Ax, &m.Ay, &m.Az,
			&m.Gx, &m.Gy, &m.Gz,
			&m.Mx, &m.My, &m.Mz,
		)
	case IntertialMode:
		m.AccFactor = int16(s.config.Acc.rangeSensitivity())
		m.GyrFactor = int16(s.config.Gyr.rangeSensitivity())
		err = decoder.Decode(
			&m.Ax, &m.Ay, &m.Az,
			&m.Gx, &m.Gy, &m.Gz,
		)
	}

	if err != nil {
		return nil, err
	}

	// handle "wrong" measures
	if m.Gx > wrongMeasureThreshold ||
		m.Gy > wrongMeasureThreshold ||
		m.Gz > wrongMeasureThreshold {
		return m, fmt.Errorf("%w: (%d, %d, %d)", ErrWrongMeasure, m.Gx, m.Gy, m.Gz)
	}

	s.removeGyroOffset(m)

	return m, nil
}

func (s *SenseHAT) removeGyroOffset(data *proto.IMUData) {
	data.Gx -= s.gxOffset
	data.Gy -= s.gyOffset
	data.Gz -= s.gzOffset
}

func (s *SenseHAT) RawMeasure() ([]byte, error) {
	var err error

	switch s.config.Mode {
	case FullMode:
		b, _, err := s.inertialBus.ReadRegBytes(QMI8658RegisterAxL, 12)
		if err != nil {
			return nil, err
		}

		copy(s.buf[:12], b)

		b, _, err = s.readMagMeasure()
		if err != nil {
			return nil, err
		}

		copy(s.buf[12:], b)
	case IntertialMode:
		s.buf, _, err = s.inertialBus.ReadRegBytes(QMI8658RegisterAxL, 12)
	}

	log.Debugf("raw read measure: %v", s.buf)

	return s.buf, err
}

func (s *SenseHAT) readMagMeasure() ([]byte, int, error) {
	tries := 20

	for ; tries > 0; tries-- {
		// TODO: читает нули, разобраться почему
		b, err := s.magnBus.ReadRegU8(AK09918_ST1)
		if err != nil {
			return nil, 0, err
		}

		// Можно читать измерения с регистров вывода
		if (b & 0x01) != 0 {
			break
		}

		time.Sleep(10 * time.Millisecond)
	}

	if tries == 0 {
		return nil, 0, fmt.Errorf("failed to read magn measure: all tries failed")
	}

	return s.magnBus.ReadRegBytes(AK09918_HXL, 6)
}

func (s *SenseHAT) Close() error {
	if s.inertialBus != nil {
		err := s.inertialBus.Close()
		if err != nil {
			return fmt.Errorf("cannot close inertial bus: %v", err)
		}
	}

	if s.magnBus != nil {
		err := s.magnBus.Close()
		if err != nil {
			return fmt.Errorf("cannot close magn bus: %v", err)
		}
	}

	return nil
}

func initInertialSensors(config internalConfig) (*i2c.I2C, error) {
	if !(config.Acc.enable || config.Gyr.enable) {
		return nil, nil
	}

	bus, err := i2c.NewI2C(I2CAddImuQMI8658, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot create new i2c: %v", err)
	}

	b, _, err := bus.ReadRegBytes(0x00, 1)
	if err != nil {
		return nil, fmt.Errorf("read from register failed: %v", err)
	}

	if b[0] != 0x05 {
		return nil, fmt.Errorf("unexpected byte was read: %d", b[0])
	}

	var sensorEnabled byte = 0x80

	configCmds := make([]configCmd, 0, 5)
	configCmds = append(configCmds, configCmd{QMI8658RegisterCtrl1, 0x60})

	if config.Acc.enable {
		configCmds = append(configCmds, configCmd{QMI8658RegisterCtrl2, config.Acc.order() | config.Acc.rangeValue()})
		sensorEnabled |= QMI8658_CTRL7_ACC_ENABLE
	}

	if config.Gyr.enable {
		configCmds = append(configCmds, configCmd{QMI8658RegisterCtrl3, config.Gyr.order() | config.Gyr.rangeValue()})
		sensorEnabled |= QMI8658_CTRL7_GYR_ENABLE
	}

	configCmds = append(configCmds,
		configCmd{QMI8658RegisterCtrl5, 0x00},
		configCmd{QMI8658RegisterCtrl7, sensorEnabled},
	)

	for _, cmd := range configCmds {
		err = bus.WriteRegU8(cmd.reg, cmd.val)
		if err != nil {
			return nil, fmt.Errorf("cannot write byte to reg: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	return bus, nil
}

func initMagnSensor(config internalConfig) (*i2c.I2C, error) {
	if !config.Mag.enable {
		return nil, nil
	}

	bus, err := i2c.NewI2C(I2CAddImuAK09918, 1)
	if err != nil {
		return nil, fmt.Errorf("cannot create new i2c: %v", err)
	}

	b, err := bus.ReadRegU8(AK09918_WIA2)
	if err != nil {
		return nil, fmt.Errorf("read from register failed: %v", err)
	}

	if b != 0x0C {
		return nil, fmt.Errorf("unexpected byte was read: %x", b)
	}

	configCmds := []configCmd{
		{AK09918_CNTL3, AK09918_SRST_BIT},
		{AK09918_CNTL2, config.Mag.order()},
	}

	for _, cmd := range configCmds {
		err = bus.WriteRegU8(cmd.reg, cmd.val)
		if err != nil {
			return nil, fmt.Errorf("cannot write byte to reg: %v", err)
		}

		time.Sleep(10 * time.Millisecond)
	}

	return bus, nil
}
