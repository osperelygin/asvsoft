// Package depthmeter предоставляет подкоманду depth-meter
package depthmeter

import (
	"asvsoft/internal/app/cli/common"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
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
	srcCfg = common.AddSerialSourceFlags(cmd)

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
		// TODO: убрать ResetInputBuffer и time.Sleep
		err = srcPort.ResetInputBuffer()
		if err != nil {
			log.Warnf("port: reset input buffer failed: %v", err)
		}

		time.Sleep(50 * time.Millisecond)

		measure, err := dm.ReadMeasure()
		if err != nil {
			log.Errorf("cannot read measure: %v", err)

			if pErr := new(serial.PortError); errors.As(err, &pErr) && pErr.Code() == serial.PortClosed {
				srcPort, err = serial_port.New(srcPort.Cfg)
				if err != nil {
					return fmt.Errorf("port closed and failed to reopen: %w", err)
				}

				dm = depthmeter.New(srcPort)

				log.Warn("port successfully reopened")

				continue
			}

			continue
		}

		b, err := packer.Pack(measure, proto.DepthMeterModuleAddr, proto.WritingModeA)
		if err != nil {
			log.Errorf("cannot pack measure: %v", err)
			continue
		}

		log.Printf("transmitted: %#v", measure)

		if dstCfg.TransmittingDisabled {
			continue
		}

		_, err = dstPort.Write(b)
		if err != nil {
			log.Errorf("cannot write measures: %v", err)
		}
	}
}
