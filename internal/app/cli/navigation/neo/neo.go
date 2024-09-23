package neo

import (
	"asvsoft/internal/app/cli/common"
	neom8t "asvsoft/internal/pkg/neo-m8t"
	"asvsoft/internal/pkg/proto"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	dstCfg *common.SerialConfig
	srcCfg *common.SerialConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "neo",
		Short: "Режим чтения данных с последовательного порта",
		Run:   Handler,
	}
	dstCfg = common.AddSerialDestinationFlags(cmd)
	srcCfg = common.AddSerialSourceFlags(cmd, "")

	return cmd
}

func Handler(cmd *cobra.Command, args []string) {
	srcPort, err := common.InitSerialPort(srcCfg)
	if err != nil {
		log.Fatalf("cannot open serial port '%s': %v", srcCfg.Port, err)
	}

	defer common.CloseSerialPort(srcPort)

	dstPort, err := common.InitSerialPort(dstCfg)
	if err != nil {
		log.Fatalf("cannot open serial port '%s': %v", dstCfg.Port, err)
	}

	defer common.CloseSerialPort(dstPort)

	packer := proto.NewPacker()
	neo := neom8t.New(&neom8t.Config{}, srcPort)

	for {
		measure, err := neo.Measure()
		if err != nil {
			log.Errorf("cannot decode message: %v", err)
			continue
		}

		b, err := packer.Pack(measure, proto.GNSSModuleAddr, proto.WritingModeA)
		if err != nil {
			log.Errorf("cannot pack measure: %v", err)
			continue
		}

		_, err = dstPort.Write(b)
		if err != nil {
			log.Errorf("cannot write measures: %v", err)
			continue
		}

		log.Printf("read the message: %#v", measure)
	}
}
