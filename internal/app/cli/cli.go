// Package cli предоставляет корневую команду asvsoft
package cli

import (
	"asvsoft/internal/app/cli/camera"
	"asvsoft/internal/app/cli/check"
	"asvsoft/internal/app/cli/controller"
	depthmeter "asvsoft/internal/app/cli/depth-meter"
	"asvsoft/internal/app/cli/lidar"
	"asvsoft/internal/app/cli/navigation"
	neom8t "asvsoft/internal/app/cli/neo-m8t"
	"asvsoft/internal/app/cli/registrar"
	sensehat "asvsoft/internal/app/cli/sense-hat"
	"asvsoft/internal/app/ctxutils"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var logLevel string

func RootCmd() *cobra.Command {
	rootCmd := cobra.Command{
		Use:               "asvsoft",
		Short:             "ПО БКУ безэкипажным надводным аппаратом",
		PersistentPreRunE: persistentPreRunHandler,
	}

	rootCmd.PersistentFlags().StringVar(
		&logLevel, "loglevel",
		"info", "",
	)

	rootCmd.AddCommand(
		controller.Cmd(),
		depthmeter.Cmd(),
		navigation.Cmd(),
		lidar.Cmd(),
		neom8t.Cmd(),
		sensehat.Cmd(),
		check.Cmd(),
		camera.Cmd(),
		registrar.Cmd(),
	)

	return &rootCmd
}

func persistentPreRunHandler(cmd *cobra.Command, args []string) error { // nolint: revive
	lvl, err := log.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	log.SetLevel(lvl)
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: time.StampMilli,
		FullTimestamp:   true,
		ForceColors:     true,
	})

	appinfo := ctxutils.GetAppInfo(cmd.Context())
	log.Infof(
		"BuildTime: %s, BuildCommit: %s, BuildBranch: %s",
		appinfo.BuildTime, appinfo.BuildCommit, appinfo.BuildBranch,
	)

	return nil
}
