// Package check предоставляет подкоманду depth-meter
package check

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
		Use:   "check",
		Short: "Тестовый модуль",
		RunE:  common.Handler(&cfg, common.CheckMode),
	}
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
