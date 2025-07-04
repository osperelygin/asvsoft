// Package depthmeter предоставляет подкоманду depth-meter
package depthmeter

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
		Use:   "depthmeter",
		Short: "Модуль обработки данных измерителя глубины",
		RunE:  common.Handler(&cfg, common.DepthMeterMode),
	}
	cfg.SrcSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}
