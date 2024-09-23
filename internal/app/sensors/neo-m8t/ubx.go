// Package neom8t осуществляет чтение и конфигурацию измерений neo m8t
package neom8t

import (
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/serial_port"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/daedaleanai/ublox"
	"github.com/daedaleanai/ublox/ubx"
	log "github.com/sirupsen/logrus"
	"go.bug.st/serial"
)

const (
	FullMode       = "ALL"
	NavPosslshMode = "NAV-POSLLH"
	NavVelnedMode  = "NAV-VELNED"
)

type Config struct {
	Rate int // Период получения навигационного решения в секундах
	Mode string
}

type NeoM8t struct {
	port *serial_port.SerialPort
	cfg  *Config
	d    *ublox.Decoder
}

func New(cfg *Config, port *serial_port.SerialPort) (*NeoM8t, error) {
	n := &NeoM8t{
		cfg:  cfg,
		port: port,
	}

	n.d = ublox.NewDecoder(n.port)

	// configurate NAV-POSLLH, NAV-VELNED rate
	err := n.configurate(0x02, 0x12)
	if err != nil {
		return nil, fmt.Errorf("cannot configurate message: %w", err)
	}

	log.Infoln("waitting 3 seconds for the configuration to be applied")
	time.Sleep(3 * time.Second)

	return n, nil
}

func (n *NeoM8t) configurate(msgIDList ...byte) error {
	rateBytes := [6]byte{}
	rateBytes[1] = byte(n.cfg.Rate)

	cfgMsg := ubx.CfgMsg2{
		MsgClass: 0x01,
		Rate:     rateBytes,
	}

	for _, msgID := range msgIDList {
		cfgMsg.MsgID = msgID

		b, err := ubx.Encode(cfgMsg)
		if err != nil {
			return fmt.Errorf("cannot encode cfg msg: %w", err)
		}

		log.Infof("writting cfg msg: %v", b)

		_, err = n.port.Write(b)
		if err != nil {
			return fmt.Errorf("cannot write cfg msg: %w", err)
		}

		time.Sleep(500 * time.Millisecond)
	}

	return nil
}

func (n *NeoM8t) Measure() (*proto.GNSSData, error) {
	var (
		data                               proto.GNSSData
		navPosllhMsgRead, navVelnedMsgRead bool
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*n.cfg.Rate)*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("failed to read measures: %w", ctx.Err())
		default:
			msg, err := n.d.Decode()
			if err != nil {
				log.Errorf("cannot decode msg: %v", err)

				portError := &serial.PortError{}
				if errors.As(err, &portError) && portError.Code() == serial.PortClosed {
					// Пересоздаем порт
					port, err := serial_port.New(n.port.Cfg)
					if err != nil {
						return nil, fmt.Errorf("port closed and failed to reopen: %w", err)
					}

					n.port = port
					n.d = ublox.NewDecoder(n.port)

					log.Warn("port successfully reopened")
				}
			}

			if navPosllhMsg, ok := msg.(*ubx.NavPosllh); ok {
				navPosllhMsgRead = true
				data.ITowNAVPOSLLH = navPosllhMsg.ITOW_ms
				data.Lon = navPosllhMsg.Lon_dege7
				data.Lat = navPosllhMsg.Lat_dege7
				data.Height = navPosllhMsg.Height_mm
				data.HMSL = navPosllhMsg.HMSL_mm
				data.HAcc = navPosllhMsg.HAcc_mm
				data.VAcc = navPosllhMsg.VAcc_mm
				log.Printf("read NAV-POSLLH msg: %#v", navPosllhMsg)
			} else if navVelnedMsg, ok := msg.(*ubx.NavVelned); ok {
				navVelnedMsgRead = true
				data.ITowNAVVELNED = navVelnedMsg.ITOW_ms
				data.VelN = navVelnedMsg.VelN_cm_s
				data.VelE = navVelnedMsg.VelE_cm_s
				data.VelD = navVelnedMsg.VelD_cm_s
				data.Speed = navVelnedMsg.Speed_cm_s
				data.GSppeed = navVelnedMsg.GSpeed_cm_s
				data.Heading = navVelnedMsg.Heading_dege5
				data.SAcc = navVelnedMsg.SAcc_cm_s
				data.CAcc = navVelnedMsg.CAcc_dege5
				log.Printf("read NAV-VELNED msg: %#v", navVelnedMsg)
			}

			if !navPosllhMsgRead || !navVelnedMsgRead {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			return &data, nil
		}
	}
}
