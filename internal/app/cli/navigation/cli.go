package navigation

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/app/cli/navigation/neo"
	sensehat "asvsoft/internal/app/cli/navigation/sense-hat"
	"asvsoft/internal/pkg/proto"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.bug.st/serial"
)

var (
	mode           string
	senseHatConfig *common.SerialConfig
	neoConfig      *common.SerialConfig
)

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "nav",
		Short: "блок навигации",
		Run:   Handler,
	}
	senseHatConfig = common.AddSerialSourceFlags(&cmd, "sensehat")
	neoConfig = common.AddSerialSourceFlags(&cmd, "neo")

	cmd.Flags().StringVar(
		&mode, "mode",
		"gnss", "режим чтения данных: gnss/imu",
	)
	cmd.AddCommand(
		sensehat.Cmd(),
		neo.Cmd(),
	)

	return &cmd
}

func Handler(cmd *cobra.Command, args []string) {
	var port serial.Port
	var err error

	switch mode {
	case "gnss":
		port, err = common.InitSerialPort(neoConfig)
		if err != nil {
			log.Fatalf("cannot init neo port: %v", err)
		}
	case "imu":
		port, err = common.InitSerialPort(neoConfig)
		if err != nil {
			log.Fatalf("cannot init neo port: %v", err)
		}
	default:
		panic("unknown mode")
	}

	defer common.CloseSerialPort(port)

	// TODO: добавить обработку sigterm

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
