package gpio

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/stianeikeland/go-rpio/v4"
)

var ()

func Cmd() *cobra.Command {
	cmd := cobra.Command{
		Use:   "gpio",
		Short: "Режим чтения данных c GPIO пинов",
		Run:   Handler,
	}

	return &cmd
}

const (
	byteSize   = 8
	bufferSize = 32
)

func Handler(cmd *cobra.Command, args []string) {
	err := rpio.Open()
	if err != nil {
		log.Fatalf("cannot open rpio: %v", err)
	}

	defer func() {
		err := rpio.Close()
		log.Printf("cannot close rpio: %v", err)
	}()

	tx := rpio.Pin(14)
	rx := rpio.Pin(15)

	tx.Output()
	rx.Input()

	dataBuffer := [bufferSize]byte{}
	byteBuffer := [byteSize]byte{}

	for {
		for b := 0; b < bufferSize; b += 1 {
			for idx := 0; idx < byteSize; idx += 1 {
				byteBuffer[idx] = byte(rx.Read())
			}

			dataBuffer[b] = convertToByte(byteBuffer)

			log.Println("readed buffer:", dataBuffer)
		}
	}
}

func convertToByte(bits [byteSize]byte) byte {
	var res byte
	
	for idx := 0; idx < byteSize; idx += 1 {
		res += bits[idx] << idx
	}

	return res
}
