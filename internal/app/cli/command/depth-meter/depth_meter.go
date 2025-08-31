// Package depthmeter предоставляет подкоманду depth-meter
package depthmeter

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
		Use:   "depthmeter",
		Short: "Модуль обработки данных измерителя глубины",
		RunE:  common.ModuleHandler(&cfg, common.DepthMeterMode),
	}
	cfg.SensorSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.ControllerSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
