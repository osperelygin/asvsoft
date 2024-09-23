package sensehat

import (
	"asvsoft/internal/app/ds"
	"asvsoft/internal/pkg/encoder"
	"asvsoft/internal/pkg/measurer"
	"asvsoft/internal/pkg/proto"
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/d2r2/go-i2c"
)

type configCmd struct {
	reg byte
	val byte
}

type IMU struct {
	buf         []byte
	config      *ImuConfig
	inertialBus *i2c.I2C
	magnBus     *i2c.I2C
}

func NewIMU(config *ImuConfig) (*IMU, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	imu := &IMU{config: config}

	switch config.Mode {
	case FullMode:
		imu.buf = make([]byte, 24)
	case IntertialMode:
		imu.buf = make([]byte, 16)
	default:
		return nil, fmt.Errorf("cannot create imu: unknown mode: '%s'", config.Mode)
	}

	defer func() {
		if err != nil {
			err = imu.Close()
			if err != nil {
				log.Errorf("failed to close imu: %v", err)
			}
		}
	}()

	imu.inertialBus, err = initInertialSensors(config)
	if err != nil {
		return nil, err
	}

	imu.magnBus, err = initMagnSensor(config)
	if err != nil {
		return nil, err
	}

	return imu, nil
}

func (imu *IMU) Measure(_ context.Context) measurer.Measurement {
	time.Sleep(imu.config.Period)
	return ds.NewMeasurement(imu.measure())
}

func (imu *IMU) measure() (*proto.IMUData, error) {
	b, err := imu.RawMeasure()
	if err != nil {
		return nil, err
	}

	measure := &proto.IMUData{}
	decoder := encoder.NewDecoder(io.NopCloser(bytes.NewBuffer(b)))

	switch imu.config.Mode {
	case FullMode:
		measure.AccFactor = int16(imu.config.Acc.rangeSensitivity())
		measure.GyrFactor = int16(imu.config.Gyr.rangeSensitivity())
		err = decoder.Decode(
			&measure.Ax, &measure.Ay, &measure.Az,
			&measure.Gx, &measure.Gy, &measure.Gz,
			&measure.Mx, &measure.My, &measure.Mz,
		)
	case IntertialMode:
		measure.AccFactor = int16(imu.config.Acc.rangeSensitivity())
		measure.GyrFactor = int16(imu.config.Gyr.rangeSensitivity())
		err = decoder.Decode(
			&measure.Ax, &measure.Ay, &measure.Az,
			&measure.Gx, &measure.Gy, &measure.Gz,
		)
	}

	if err != nil {
		return nil, err
	}

	return measure, nil
}

func (imu *IMU) RawMeasure() ([]byte, error) {
	var err error

	switch imu.config.Mode {
	case FullMode:
		b, _, err := imu.inertialBus.ReadRegBytes(QMI8658RegisterAxL, 12)
		if err != nil {
			return nil, err
		}

		copy(imu.buf[:12], b)

		b, _, err = imu.readMagMeasure()
		if err != nil {
			return nil, err
		}

		copy(imu.buf[12:], b)
	case IntertialMode:
		imu.buf, _, err = imu.inertialBus.ReadRegBytes(QMI8658RegisterAxL, 12)
	}

	return imu.buf, err
}

func (imu *IMU) readMagMeasure() ([]byte, int, error) {
	tries := 20

	for ; tries > 0; tries-- {
		// TODO: читает нули, разобраться почему
		b, err := imu.magnBus.ReadRegU8(AK09918_ST1)
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

	return imu.magnBus.ReadRegBytes(AK09918_HXL, 6)
}

func (imu *IMU) Close() error {
	if imu.inertialBus != nil {
		err := imu.inertialBus.Close()
		if err != nil {
			return fmt.Errorf("cannot close inertial bus: %v", err)
		}
	}

	if imu.magnBus != nil {
		err := imu.magnBus.Close()
		if err != nil {
			return fmt.Errorf("cannot close magn bus: %v", err)
		}
	}

	return nil
}

func initInertialSensors(config *ImuConfig) (*i2c.I2C, error) {
	if !(config.Acc.Enable || config.Gyr.Enable) {
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

	if config.Acc.Enable {
		configCmds = append(configCmds, configCmd{QMI8658RegisterCtrl2, config.Acc.order() | config.Acc.rangeValue()})
		sensorEnabled |= QMI8658_CTRL7_ACC_ENABLE
	}

	if config.Gyr.Enable {
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

func initMagnSensor(config *ImuConfig) (*i2c.I2C, error) {
	if !config.Mag.Enable {
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
