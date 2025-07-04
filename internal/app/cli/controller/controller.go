// Package controller предоставляет подкоманду controller
package controller

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/communication"
	"asvsoft/internal/pkg/proto"
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
		Short: "Контроллер управления",
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

	log.Debugln("waiting sync request...")

	err = communication.NewSyncer(proto.ControlModuleID).WithReadWriter(srcPort).Serve()
	if err != nil {
		return fmt.Errorf("cannot sync: %v", err)
	}

	r := communication.NewReceiver(srcPort)
	defer func() {
		err = r.Close()
		if err != nil {
			log.Errorf("cannot close receiver: %v", err)
		}
	}()

	log.Debugln("receiving messages...")

	for {
		msg, err := r.Recieve()
		if err != nil {
			log.Errorf("receive failed: %v", err)
		}

		// TODO: обработка полученных данных
		_ = msg
	}
}
