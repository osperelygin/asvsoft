// Package sensehat предоставляет подкоманду sense-hat
package sensehat

import (
	"asvsoft/internal/app/cli/common"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	period time.Duration
	dstCfg *serial_port.SerialPortConfig
	imuCfg *sensehat.ImuConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sense-hat",
		Short: "Обработка и передача данных SENSE HAT (C)",
		RunE:  Handler,
	}

	dstCfg = common.AddSerialDestinationFlags(cmd)
	imuCfg = sensehat.NewImuConfig()

	// TODO: реализовать различные режимы
	cmd.Flags().DurationVar(
		&period, "period",
		10*time.Millisecond, "период чтения данных",
	)

	cmd.Flags().StringVar(
		&imuCfg.Mode, "mode",
		sensehat.IntertialMode, "режим чтения данных",
	)

	cmd.Flags().Float32Var(
		&imuCfg.Acc.Order, "acc-order",
		125, "частота обновления данных на регистрах АСС в Гц",
	)

	cmd.Flags().Float32Var(
		&imuCfg.Gyr.Order, "gyr-order",
		125, "частота обновления данных на регистрах гироскопа в Гц",
	)

	cmd.Flags().Float32Var(
		&imuCfg.Mag.Order, "mag-order",
		20, "частота обновления данных на регистрах магнитометра в Гц",
	)

	cmd.Flags().IntVar(
		&imuCfg.Acc.Range, "acc-range",
		2, "диапазон измерений АСС в g",
	)

	cmd.Flags().IntVar(
		&imuCfg.Gyr.Range, "gyr-range",
		128, "диапазон измерений гироскопа в град/с",
	)

	return cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	imu, err := sensehat.NewIMU(imuCfg)
	if err != nil {
		return fmt.Errorf("cannot init imu: %v", err)
	}

	defer func() {
		err = imu.Close()
		if err != nil {
			log.Errorf("failed to close imu: %v", err)
		}
	}()

	dstPort, err := serial_port.New(dstCfg)
	if err != nil {
		return fmt.Errorf("cannot open serial port '%s': %v", dstCfg.Port, err)
	}

	defer dstPort.Close()

	var packer proto.Packer

	for {
		// TODO: оптимизировать аллокацию памяти в ReadRegBytes
		measure, err := imu.Measure()
		if err != nil {
			log.Errorf("imu measure failed : %v", err)
			continue
		}

		log.Infof("decode measure: %#v", measure)

		b, err := packer.Pack(measure, proto.IMUModuleAddr, proto.WritingModeA)
		if err != nil {
			log.Errorf("cannot pack data: %v", err)
			continue
		}

		_, err = dstPort.Write(b)
		if err != nil {
			log.Errorf("cannot write data to target: %v", err)
		}

		time.Sleep(period)
	}
}
