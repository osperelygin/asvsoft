package serial

import (
	"asvsoft/internal/app/cli/common"
	depthmeter "asvsoft/internal/app/sensors/depth-meter"
	"asvsoft/internal/pkg/proto"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	dstCfg *common.SerialConfig
	srcCfg *common.SerialConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serial",
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
		log.Fatalf("cannot init depth meter port '%s': %v", srcCfg.Port, err)
	}

	defer common.CloseSerialPort(srcPort)

	dstPort, err := common.InitSerialPort(dstCfg)
	if err != nil {
		log.Fatalf("cannot open serial port '%s': %v", dstCfg.Port, err)
	}

	defer common.CloseSerialPort(dstPort)

	dm := depthmeter.New(srcPort)
	packer := proto.NewPacker()

	for {
		srcPort.ResetInputBuffer()
		time.Sleep(50 * time.Millisecond)

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

		_, err = dstPort.Write(b)
		if err != nil {
			log.Errorf("cannot write measures: %v", err)
			continue
		}

		log.Println(measure.SystemTime, measure.Distance, measure.Strength, measure.Precision)
	}
}
