// Package depthmeter предоставляет подкоманду depth-meter
package depthmeter

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/measurer"

	"github.com/spf13/cobra"
)

var (
	cfg config.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depthmeter",
		Short: "Режим чтения данных с последовательного порта",
		RunE:  Handler,
	}
	cfg.SrcSerialPort = common.AddSerialSourceFlags(cmd)
	cfg.DstSerialPort = common.AddSerialDestinationFlags(cmd)

	return cmd
}

func Handler(cmd *cobra.Command, args []string) error { // nolint: revive
	ctx := config.WrapContext(cmd.Context(), &cfg)

	m, t, err := common.Init(ctx, common.DepthMeterMode)
	if err != nil {
		return err
	}

	err = measurer.Run(ctx, m, t)
	if err != nil {
		return err
	}

	return nil
}
