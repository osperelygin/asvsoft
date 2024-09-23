package neom8t

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/daedaleanai/ublox"
	"github.com/daedaleanai/ublox/ubx"
	log "github.com/sirupsen/logrus"
)

const (
	FullMode       = "ALL"
	NavPosslshMode = "NAV-POSLLH"
	NavVelnedMode  = "NAV-VELNED"
)

type Config struct {
	Mode string
}

type NeoM8t struct {
	d   *ublox.Decoder
	cfg *Config
}

func New(cfg *Config, r io.Reader) *NeoM8t {
	return &NeoM8t{
		d:   ublox.NewDecoder(r),
		cfg: cfg,
	}
}

func (n *NeoM8t) Measure() (*proto.GNSSData, error) {
	var (
		data                               proto.GNSSData
		navPosllhMsgRead, navVelnedMsgRead bool
	)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout exceeded, failed to read measures")
		default:
			msg, err := n.d.Decode()
			if err != nil {
				log.Errorf("cannot decode msg: %v", err)
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
			}

			if !navPosllhMsgRead || !navVelnedMsgRead {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			return &data, nil
		}
	}
}
