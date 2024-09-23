// Package cli предоставляет корневую команду asvsoft
package cli

import (
	"asvsoft/internal/app/cli/controller"
	depthmeter "asvsoft/internal/app/cli/depth-meter"
	"asvsoft/internal/app/cli/navigation"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "Jan _2 15:04:05.000",
		FullTimestamp:   true,
	})

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
