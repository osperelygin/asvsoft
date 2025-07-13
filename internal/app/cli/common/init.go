package common

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/app/sensors/camera"
	"asvsoft/internal/app/sensors/check"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	"asvsoft/internal/app/sensors/lidar"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	"asvsoft/internal/pkg/communication"
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	serialport "asvsoft/internal/pkg/serial-port"
	"context"
	"fmt"
	"slices"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type RunMode int

const (
	DepthMeterMode RunMode = iota + 1
	LidarMode
	NeoM8tMode
	ImuMode
	MockCameraMode
	CameraMode
	NavMode
	RegistratorMode
	CheckMode
)

var (
	requiredSrcSerialPortRunMode = []RunMode{
		DepthMeterMode, LidarMode, NeoM8tMode,
	}
)

// Init общая функция инициализации модуля камеры, лидара, ИНС и ГНСС, измерителя глубины ,
// модуля навигации и модуля проверки. Требуемые для работы модуля порты-источники и
// порты-назначения обернуты в объекте sender'a.
func Init(ctx context.Context, mode RunMode, opts ...ModuleOptions) (*communication.Sender, *communication.Syncer, error) {
	cfg := config.FromContext(ctx)

	var (
		srcPort *serialport.Wrapper
		err     error
	)

	if slices.Contains(requiredSrcSerialPortRunMode, mode) {
		srcPort, err = serialport.New(cfg.SensorSerialPort)
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
		m    communication.MeasureCloser
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
	case CameraMode:
		cam, err := camera.New()
		if err != nil {
			return nil, nil, err
		}

		m = cam.WithLogger(
			logger.Wrap(logrus.StandardLogger(), "[camera]"),
		)

		addr = proto.CameraModuleID
	case MockCameraMode:
		// hardcode
		m, err = camera.NewMockCamera([][3]int16{
			{1090, -12711, -16544},
			{1087, -12768, -16545},
			{1085, -12669, -16545},
		})
		if err != nil {
			return nil, nil, err
		}

		addr = proto.CameraModuleID
	default:
		panic(fmt.Sprintf("unknown run mode: %q", addr))
	}

	sendMode := proto.WritingModeA
	if len(opts) > 0 && opts[0].SendMode != 0 {
		sendMode = opts[0].SendMode
	}

	sndr := communication.NewSender(m, addr, sendMode)
	sncr := communication.NewSyncer(addr)

	if !cfg.ControllerSerialPort.TransmittingDisabled {
		dstPort, err := serialport.New(cfg.ControllerSerialPort)
		if err != nil {
			return nil, nil, err
		}

		err = dstPort.ResetOutputBuffer()
		if err != nil {
			log.Errorf("cannot reset output buffer: %v", err)
		}

		err = dstPort.ResetInputBuffer()
		if err != nil {
			log.Errorf("cannot reset input buffer: %v", err)
		}

		dstPort.SetLogger(log.StandardLogger())
		sndr.WithReadWriteCloser(dstPort).
			WithSleep(cfg.ControllerSerialPort.Sleep).
			WithSync(cfg.ControllerSerialPort.Sync)
		sncr.WithReadWriter(dstPort)
	}

	return sndr, sncr, nil
}
