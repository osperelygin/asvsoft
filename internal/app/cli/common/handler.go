package common

import (
	"asvsoft/internal/app/config"
	"asvsoft/internal/pkg/communication"

	"github.com/spf13/cobra"
)

func Handler(cfg *config.Config, mode RunMode) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := config.WrapContext(cmd.Context(), cfg)

		m, t, err := Init(ctx, mode)
		if err != nil {
			return err
		}

		err = communication.Entrypoint(ctx, m, t)
		if err != nil {
			return err
		}

		return nil
	}
}
