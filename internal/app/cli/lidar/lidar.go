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
		Short: "Режим чтения данных с последовательного порта",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.Handler(cmd, args, &cfg, common.LidarMode)
		},
	}
	cfg.SrcSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
