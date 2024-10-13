// Package cli предоставляет корневую команду asvsoft
package cli

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/cli/controller"
	depthmeter "asvsoft/internal/app/cli/depth-meter"
	"asvsoft/internal/app/cli/lidar"
	"asvsoft/internal/app/cli/navigation"
	"asvsoft/internal/app/ctxutils"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func RootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:               "asvsoft",
		Short:             "ПО БКУ безэкипажным надводным аппаратом",
		PersistentPreRunE: persistentPreRunHandler,
	}

	rootCmd.PersistentFlags().StringVar(
		&common.LogLevel, "loglevel",
		"info", "",
	)

	rootCmd.AddCommand(
		controller.Cmd(),
		depthmeter.Cmd(),
		navigation.Cmd(),
		lidar.Cmd(),
	)

	return &rootCmd
}

func persistentPreRunHandler(cmd *cobra.Command, args []string) error { // nolint: revive
	lvl, err := log.ParseLevel(common.LogLevel)
	if err != nil {
		return err
	}

	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "Jan _2 15:04:05.000",
		FullTimestamp:   true,
	})

	appinfo := ctxutils.GetAppInfo(cmd.Context())
	log.Infof(
		"BuildTime: %s, BuildCommit: %s, BuildBranch: %s",
		appinfo.BuildTime, appinfo.BuildCommit, appinfo.BuildBranch,
	)

	return nil
}
