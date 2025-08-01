// Package controller предоставляет подкоманду controller
package controller

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/proto"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Контроллер управления",
		RunE:  common.ControllerHandler(proto.ControlModuleID),
	}
	cmd.Flags().StringVarP(
		&common.CtrlCfgPath, "config", "c",
		"/etc/asvsoft/config.yaml",
		"Path to config",
	)

	return cmd
}
