// Package lidar предоставляет подкоманду lidar
package lidar

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
		Use:   "lidar",
		Short: "Модуль обработки данных лидара",
		RunE:  common.ModuleHandler(&cfg, common.LidarMode),
	}
	cfg.SensorSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.ControllerSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
