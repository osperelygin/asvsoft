// Package depthmeter предоставляет подкоманду depth-meter
package depthmeter

import (
	"asvsoft/internal/app/cli/common"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	dstCfg *serial_port.SerialPortConfig
	srcCfg *serial_port.SerialPortConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "depthmeter",
		Short: "Режим чтения данных с последовательного порта",
		RunE:  Handler,
	}
	dstCfg = common.AddSerialDestinationFlags(cmd)
	srcCfg = common.AddSerialSourceFlags(cmd, "")

	return cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	srcPort, err := serial_port.New(srcCfg)
	if err != nil {
		return fmt.Errorf("cannot init depth meter port '%s': %v", srcCfg.Port, err)
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

	dm := depthmeter.New(srcPort)
	packer := proto.NewPacker()

	for {
		_ = srcPort.ResetInputBuffer()

		measure, err := dm.ReadMeasure()
		if err != nil {
			log.Errorf("cannot read measure: %v", err)
			continue
		}

		b, err := packer.Pack(measure, proto.DepthMeterModuleAddr, proto.WritingModeA)
		if err != nil {
			log.Errorf("cannot pack measure: %v", err)
			continue
		}

		log.Printf("%+v", measure)

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
