package common

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/measurer"

	"github.com/spf13/cobra"
)

func Handler(cmd *cobra.Command, args []string, cfg *config.Config, mode RunMode) error { // nolint: revive
	ctx := config.WrapContext(cmd.Context(), cfg)

	m, t, err := Init(ctx, mode)
	if err != nil {
		return err
	}

	err = measurer.Run(ctx, m, t)
	if err != nil {
		return err
	}

	return nil
}
