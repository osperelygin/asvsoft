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
		Short: "Режим чтения данных с последовательного порта",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.Handler(cmd, args, &cfg, common.CheckMode)
		},
	}
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
