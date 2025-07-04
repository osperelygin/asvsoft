// Package lidar предоставляет подкоманду lidar
package lidar

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lidar",
		Short: "Модуль обработки данных лидара",
		RunE:  common.Handler(&cfg, common.LidarMode),
	}
	cfg.SrcSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
