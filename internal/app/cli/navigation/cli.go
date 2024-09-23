// Package navigation предоставляет подкоманду nav
package navigation

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/cli/navigation/neo"
	sensehat "asvsoft/internal/app/cli/navigation/sense-hat"
	"asvsoft/internal/pkg/proto"
	serialport "asvsoft/internal/pkg/serial-port"
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
	srcCfg *serialport.Config
)

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "nav",
		Short: "блок навигации",
		RunE:  Handler,
	}
	srcCfg = common.AddSerialSourceFlags(&cmd)

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
		port, err = serialport.New(srcCfg)
		if err != nil {
			return fmt.Errorf("cannot init gnss port: %v", err)
		}
	default:
		panic(fmt.Sprintf("unknown mode: '%s'", mode))
	}

	defer port.Close()

	packer := proto.Packer{}

	for {
		rawData, err := proto.Read(port)
		if err != nil {
			log.Errorf("read data from source port failed: %v", err)
		}

		data, err := packer.Unpack(rawData)
		if err != nil {
			log.Errorf("unpack failed: %v", err)
		}

		log.Printf("received: %+v", data)
	}
}
