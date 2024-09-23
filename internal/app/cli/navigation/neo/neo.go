// Package neo предоставляет подкоманду neo
package neo

import (
	"asvsoft/internal/app/cli/common"
	neom8t "asvsoft/internal/app/sensors/neo-m8t"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dstCfg    *serial_port.SerialPortConfig
	srcCfg    *serial_port.SerialPortConfig
	neoConfig neom8t.Config
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neo",
		Short: "Режим чтения данных с последовательного порта",
		RunE:  Handler,
	}
	dstCfg = common.AddSerialDestinationFlags(cmd)
	srcCfg = common.AddSerialSourceFlags(cmd, "")

	cmd.Flags().IntVar(
		&neoConfig.Rate, "rate",
		1, "navigation solution rate in second",
	)

	return cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	srcPort, err := serial_port.New(srcCfg)
	if err != nil {
		return fmt.Errorf("cannot open serial port '%s': %v", srcCfg.Port, err)
	}
	defer srcPort.Close()

	var dstPort *serial_port.SerialPort

	if !dstCfg.TransmittingDisabled {
		dstPort, err = serial_port.New(dstCfg)
		if err != nil {
			return fmt.Errorf("cannot open serial port '%s': %v", dstCfg.Port, err)
		}
	}
	defer dstPort.Close()

	neo, err := neom8t.New(&neoConfig, srcPort)
	if err != nil {
		return fmt.Errorf("cannot create neo instance: %v", err)
	}

	packer := proto.NewPacker()

	for {
		measure, err := neo.Measure()
		if err != nil {
			log.Errorf("cannot decode message: %v", err)
			continue
		}

		log.Printf("read the message: %#v", measure)

		b, err := packer.Pack(measure, proto.GNSSModuleAddr, proto.WritingModeA)
		if err != nil {
			log.Errorf("cannot pack measure: %v", err)
			continue
		}

		if dstCfg.TransmittingDisabled {
			continue
		}

		_, err = dstPort.Write(b)
		if err != nil {
			log.Errorf("cannot write measures: %v", err)
			continue
		}
	}
}
