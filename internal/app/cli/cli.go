package cli

import (
	"asvsoft/internal/app/cli/controller"
	depthmeter "asvsoft/internal/app/cli/depth-meter"
	"asvsoft/internal/app/cli/navigation"

	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:   "asvsoft",
		Short: "ПО БКУ безэкипажным надводным аппаратом",
	}

	rootCmd.AddCommand(
		controller.Cmd(),
		depthmeter.Cmd(),
		navigation.Cmd(),
	)

	return &rootCmd
}
