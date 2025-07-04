// Package camera предоставляет подкоманду camera
package camera

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
		Use:   "camera",
		Short: "Модуль обработки данных камеры",
		RunE:  common.Handler(&cfg, common.CameraMode),
	}

	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
