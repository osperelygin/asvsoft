// Package navigation предоставляет подкоманду nav
package navigation

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/cli/navigation/neo"
	sensehat "asvsoft/internal/app/cli/navigation/sense-hat"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

const (
	gnssMode = "gnss"
	imuMode  = "imu"
)

var (
	mode   string
	srcCfg *serial_port.SerialPortConfig
)

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "nav",
		Short: "блок навигации",
		RunE:  Handler,
	}
	srcCfg = common.AddSerialSourceFlags(&cmd, "")

	cmd.Flags().StringVar(
		&mode, "mode",
		gnssMode, "режим чтения данных: gnss/imu",
	)
	cmd.AddCommand(
		sensehat.Cmd(),
		neo.Cmd(),
	)

	return &cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	var (
		port serial.Port
		err  error
	)

	switch mode {
	case gnssMode, imuMode:
		port, err = serial_port.New(srcCfg)
		if err != nil {
			return fmt.Errorf("cannot init gnss port: %v", err)
		}
	default:
		panic(fmt.Sprintf("unknown mode: '%s'", mode))
	}

	defer port.Close()

	packer := proto.Packer{}

	for {
		rawData, err := proto.Read(port, 1<<10)
		if err != nil {
			log.Errorf("read data from source port failed: %v", err)
		}

		data, err := packer.Unpack(rawData)
		if err != nil {
			log.Errorf("unpack failed: %v", err)
		}

		log.Printf("received measures: %#v", data)
	}
}
