// Package neom8t предоставляет подкоманду neo
package neom8t

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"

	"github.com/spf13/cobra"
)

var (
	cfg config.ModuleConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neo-m8t",
		Short: "Модуль обработки данных ГНСС",
		RunE:  common.ModuleHandler(&cfg, common.NeoM8tMode),
	}
	cfg.ControllerSerialPort = common.AddSerialDestinationFlags(cmd)
	cfg.SensorSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.NeoM8t = new(neom8t.Config)

	cmd.Flags().IntVar(
		&cfg.NeoM8t.Rate, "rate",
		1, "navigation solution rate in second",
	)

	return cmd
}
