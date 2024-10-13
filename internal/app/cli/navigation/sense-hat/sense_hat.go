// Package sensehat предоставляет подкоманду sense-hat
package sensehat

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"
	sensehat "asvsoft/internal/app/sensors/sense-hat"
	"time"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sense-hat",
		Short: "Обработка и передача данных SENSE HAT (C)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.Handler(cmd, args, &cfg, common.ImuMode)
		},
	}

	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)
	cfg.Imu = sensehat.NewImuConfig()

	// TODO: реализовать различные режимы
	cmd.Flags().DurationVar(
		&cfg.Imu.Period, "period",
		10*time.Millisecond, "период чтения данных",
	)

	cmd.Flags().StringVar(
		&cfg.Imu.Mode, "mode",
		sensehat.IntertialMode, "режим чтения данных",
	)

	cmd.Flags().Float32Var(
		&cfg.Imu.Acc.Order, "acc-order",
		125, "частота обновления данных на регистрах АСС в Гц",
	)

	cmd.Flags().Float32Var(
		&cfg.Imu.Gyr.Order, "gyr-order",
		125, "частота обновления данных на регистрах гироскопа в Гц",
	)

	cmd.Flags().Float32Var(
		&cfg.Imu.Mag.Order, "mag-order",
		20, "частота обновления данных на регистрах магнитометра в Гц",
	)

	cmd.Flags().IntVar(
		&cfg.Imu.Acc.Range, "acc-range",
		2, "диапазон измерений АСС в g",
	)

	cmd.Flags().IntVar(
		&cfg.Imu.Gyr.Range, "gyr-range",
		128, "диапазон измерений гироскопа в град/с",
	)

	return cmd
}
