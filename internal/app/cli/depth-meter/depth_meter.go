package depthmeter

import (
	"asvsoft/internal/app/cli/depth-meter/gpio"
	"asvsoft/internal/app/cli/depth-meter/serial"

	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "depth-meter",
		Short: "Обработки и передача данных измерителя глубины",
	}

	cmd.AddCommand(
		serial.Cmd(),
		gpio.Cmd(),
	)

	return &cmd
}
