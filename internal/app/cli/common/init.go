package common

import (
	"asvsoft/internal/app/config"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	"asvsoft/internal/pkg/measurer"
	"asvsoft/internal/pkg/proto"
	serialport "asvsoft/internal/pkg/serial-port"
	"asvsoft/internal/pkg/transmitter"
	"context"

	log "github.com/sirupsen/logrus"
)

type RunMode int

const (
	DepthMeterMode RunMode = iota
	NeoM8tMode
	ImuMode
	NavMode
)

func Init(ctx context.Context, mode RunMode) (m measurer.Measurer, t transmitter.Transmitter, err error) {
	cfg := config.FromContext(ctx)

	var srcPort *serialport.Wrapper

	if mode != ImuMode {
		srcPort, err = serialport.New(cfg.SrcSerialPort)
		if err != nil {
			return nil, nil, err
		}

		srcPort.SetLogger(log.StandardLogger())
	}

	var dstPort *serialport.Wrapper

	if !cfg.DstSerialPort.TransmittingDisabled {
		dstPort, err = serialport.New(cfg.DstSerialPort)
		if err != nil {
			return nil, nil, err
		}

		dstPort.SetLogger(log.StandardLogger())
	}

	var ct *transmitter.CommonTransmitter

	switch mode {
	case DepthMeterMode:
		m = depthmeter.New(srcPort)

		ct = transmitter.New(proto.DepthMeterModuleAddr, proto.WritingModeA)
		if !cfg.DstSerialPort.TransmittingDisabled {
			ct.SetWritter(dstPort)
		}
	case NeoM8tMode:
		m, err = neom8t.New(cfg.NeoM8t, srcPort)
		if err != nil {
			return nil, nil, err
		}

		ct = transmitter.New(proto.GNSSModuleAddr, proto.WritingModeA)
		if !cfg.DstSerialPort.TransmittingDisabled {
			ct.SetWritter(dstPort)
		}
	case ImuMode:
		m, err = sensehat.NewIMU(cfg.Imu)
		if err != nil {
			return nil, nil, err
		}

		ct = transmitter.New(proto.IMUModuleAddr, proto.WritingModeA)
		if !cfg.DstSerialPort.TransmittingDisabled {
			ct.SetWritter(dstPort)
		}
	case NavMode:
		panic("implement me")
	}

	t = ct

	return m, t, nil
}
