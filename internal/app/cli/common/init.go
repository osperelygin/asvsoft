package common

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/app/sensors/check"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	"asvsoft/internal/app/sensors/lidar"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	"asvsoft/internal/pkg/communication"
	"asvsoft/internal/pkg/proto"
	serialport "asvsoft/internal/pkg/serial-port"
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type RunMode int

const (
	DepthMeterMode RunMode = iota
	LidarMode
	NeoM8tMode
	ImuMode
	NavMode
	CheckMode
)

func Init(ctx context.Context, mode RunMode) (communication.Measurer, communication.Transmitter, error) {
	cfg := config.FromContext(ctx)

	var (
		srcPort *serialport.Wrapper
		err     error
	)

	if mode != ImuMode && mode != CheckMode {
		srcPort, err = serialport.New(cfg.SrcSerialPort)
		if err != nil {
			return nil, nil, err
		}

		srcPort.SetLogger(log.StandardLogger())

		err = srcPort.ResetInputBuffer()
		if err != nil {
			log.Errorf("cannot reset input buffer: %v", err)
		}
	}

	var (
		m    communication.Measurer
		addr proto.ModuleID
	)

	switch mode {
	case DepthMeterMode:
		m = depthmeter.New(srcPort)
		addr = proto.DepthMeterModuleID
	case LidarMode:
		m = lidar.New(srcPort)
		addr = proto.LidarModuleID
	case NeoM8tMode:
		m, err = neom8t.New(cfg.NeoM8t, srcPort)
		if err != nil {
			return nil, nil, err
		}

		addr = proto.GNSSModuleID
	case ImuMode:
		m, err = sensehat.NewIMU(cfg.Imu)
		if err != nil {
			return nil, nil, err
		}

		addr = proto.IMUModuleID
	case NavMode:
		panic("implement me")
	case CheckMode:
		addr = proto.CheckModuleID
		m = check.New()
	default:
		panic(fmt.Sprintf("unknown run mode: %q", addr))
	}

	t := communication.NewCommonTransmitter(addr, proto.WritingModeA)

	if !cfg.DstSerialPort.TransmittingDisabled {
		dstPort, err := serialport.New(cfg.DstSerialPort)
		if err != nil {
			return nil, nil, err
		}

		err = dstPort.ResetOutputBuffer()
		if err != nil {
			log.Errorf("cannot reset output buffer: %v", err)
		}

		dstPort.SetLogger(log.StandardLogger())
		t.WithWritter(dstPort)
	}

	return m, t, nil
}
