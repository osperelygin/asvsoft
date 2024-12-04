package common

import (
	"asvsoft/internal/app/config"
	"fmt"

	"github.com/spf13/cobra"
)

func Handler(cfg *config.Config, mode RunMode) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := config.WrapContext(cmd.Context(), cfg)

		sndr, sncr, err := Init(ctx, mode)
		if err != nil {
			return err
		}

		err = sncr.Sync()
		if err != nil {
			return fmt.Errorf("cannot sync: %v", err)
		}

		err = sndr.Start(ctx)
		if err != nil {
			return err
		}

		return nil
	}
}
