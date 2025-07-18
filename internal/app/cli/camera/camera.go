// Package camera предоставляет подкоманду camera
package camera

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/proto"

	"github.com/spf13/cobra"
)

var (
	cfg config.ModuleConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "camera",
		Short: "Модуль обработки данных камеры",
		RunE: common.ModuleHandler(
			&cfg,
			common.CameraMode,
			common.ModuleOptions{SendMode: proto.WritingModeB},
		),
	}

	cfg.ControllerSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
