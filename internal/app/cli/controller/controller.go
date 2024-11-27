// Package controller предоставляет подкоманду controller
package controller

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/communication"
	serialport "asvsoft/internal/pkg/serial-port"
	"fmt"

	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

var (
	controllerConfig *serialport.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Режим чтения данных с последовательного порта",
		RunE:  Handler,
	}
	controllerConfig = common.AddSerialSourceFlags(cmd)

	return cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	srcPort, err := serialport.New(controllerConfig)
	if err != nil {
		return fmt.Errorf("cannot open serial port %q: %v", controllerConfig.Port, err)
	}
	defer srcPort.Close()

	for {
		msg, err := communication.Recieve(srcPort)
		if err != nil {
			log.Errorf("receive failed: %v", err)
		}

		// TODO: обработка полученных данных
		_ = msg
	}
}
