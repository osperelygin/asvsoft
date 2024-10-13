// Package neo предоставляет подкоманду neo
package neo

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neo",
		Short: "Режим чтения данных с последовательного порта",
		RunE: func(cmd *cobra.Command, args []string) error {
			return common.Handler(cmd, args, &cfg, common.NeoM8tMode)
		},
	}
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)
	cfg.SrcSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.NeoM8t = new(neom8t.Config)

	cmd.Flags().IntVar(
		&cfg.NeoM8t.Rate, "rate",
		1, "navigation solution rate in second",
	)

	return cmd
}
