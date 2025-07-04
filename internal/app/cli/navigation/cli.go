// Package navigation предоставляет подкоманду nav
package navigation

import (
	"asvsoft/internal/app/cli/common"
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
	dstCfg *serialport.Config // nolint: unused
)

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "navigation",
		Short: "Модуль навигации",
		RunE:  Handler,
	}
	srcCfg = common.AddSerialSourceFlags(&cmd)
	dstCfg = common.AddSerialDestinationFlags(&cmd)

	cmd.Flags().StringVar(
		&mode, "mode",
		gnssMode, "режим чтения данных: gnss/imu",
	)

	return &cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	var (
		port serial.Port
		err  error
	)

	// TODO: перейти к общему методу инициализации

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

	for {
		rawData, err := proto.Read(port)
		if err != nil {
			log.Errorf("read data from source port failed: %v", err)
		}

		var msg proto.Message

		err = msg.Unmarshal(rawData)
		if err != nil {
			log.Errorf("unpack failed: %v", err)
		}

		log.Printf("received data: %+v", msg)
	}
}
