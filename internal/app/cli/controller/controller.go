// Package controller предоставляет подкоманду controller
package controller

import (
	"asvsoft/internal/app/cli/common"
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var (
	depthMeterConfig *serial_port.SerialPortConfig
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "controller",
		Short: "Режим чтения данных с последовательного порта",
		RunE:  Handler,
	}
	depthMeterConfig = common.AddSerialSourceFlags(cmd, "depth-meter")

	return cmd
}

func Handler(_ *cobra.Command, _ []string) error {
	srcPort, err := serial_port.New(depthMeterConfig)
	if err != nil {
		return fmt.Errorf("cannot open serial port '%s': %v", depthMeterConfig.Port, err)
	}

	defer srcPort.Close()

	packer := proto.Packer{}

	for {
		rawData, err := proto.Read(srcPort, 1<<10)
		if err != nil {
			log.Errorf("read failed: %v", err)
		}

		data, err := packer.Unpack(rawData)
		if err != nil {
			log.Errorf("unpack failed: %v", err)
		}

		measures, ok := data.(*proto.DepthMeterData)
		if !ok {
			log.Errorf("cast data to *proto.DepthMeterData failed: %v", err)
		}

		log.Printf("received measures: %v", measures)
	}
}

// func Handler(_ *cobra.Command, _ []string) {
// 	initSPI()

// 	packer := proto.Packer{}
// 	spiReader := spireader.SPIReader{}

// 	for {
// 		rawData, err := proto.Read(&spiReader, 1<<10)
// 		if err != nil {
// 			log.Printf("spi read failed: %v", err)
// 		}

// 		data, err := packer.Unpack(rawData)
// 		if err != nil {
// 			log.Printf("unpack failed: %v", err)
// 		}

// 		measures, ok := data.(*proto.DepthMeterData)
// 		if !ok {
// 			log.Printf("cast data to *proto.DepthMeterData failed: %v", err)
// 		}

// 		log.Printf("received measures: %#v", measures)
// 	}
// }

// func initSPI() {
// 	err := rpio.Open()
// 	if err != nil {
// 		log.Fatalf("cannot open rpio: %v", err)
// 	}

// 	err = rpio.SpiBegin(rpio.Spi0)
// 	if err != nil {
// 		log.Fatalf("cannot spi begin: %v", err)
// 	}

// 	rpio.SpiChipSelect(1)

// 	log.Println("init SPI successful")
// }
