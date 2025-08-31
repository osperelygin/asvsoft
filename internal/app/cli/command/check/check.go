// Package check предоставляет подкоманду depth-meter
package check

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
		Use:   "check",
		Short: "Тестовый модуль",
		RunE:  common.ModuleHandler(&cfg, common.CheckMode),
	}
	cfg.ControllerSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
