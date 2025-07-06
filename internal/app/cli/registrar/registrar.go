// Package registrar предоставляет подкоманду registrar
package registrar

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/proto"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registrar",
		Short: "Бортовой регистратор",
		RunE:  common.ControllerHandler(proto.RegistratorModuleID),
	}
	cmd.Flags().StringVarP(
		&common.CtrlCfgPath, "config", "c",
		"/etc/asvsoft/config.yaml",
		"Path to config",
	)

	return cmd
}
