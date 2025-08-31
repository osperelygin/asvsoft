// Package neom8t предоставляет подкоманду neo
package neom8t

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"

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
	cfg.NeoM8t = new(config.NeoM8tConfig)

	cmd.Flags().IntVar(
		&cfg.NeoM8t.Rate, "rate",
		1, "navigation solution rate in second",
	)

	return cmd
}
