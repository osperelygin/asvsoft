// Package registrar предоставляет подкоманду registrar
package registrar

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
		Use:   "registrar",
		Short: "Бортовой регистратор",
		RunE:  common.ControllerHandler(proto.RegistratorModuleID, ctrlCfgPath),
	}
	cmd.Flags().StringVarP(
		ctrlCfgPath, "config", "c",
		"/etc/asvsoft/config.yaml",
		"Path to config",
	)

	return cmd
}
