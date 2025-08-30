// Package controller предоставляет подкоманду controller
package controller

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/proto"

	"github.com/spf13/cobra"
)

var (
	ctrlCfgPath *string
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Контроллер управления",
		RunE:  common.ControllerHandler(proto.ControlModuleID, ctrlCfgPath),
	}
	cmd.Flags().StringVarP(
		ctrlCfgPath, "config", "c",
		"/etc/asvsoft/config.yaml",
		"Path to config",
	)

	return cmd
}
