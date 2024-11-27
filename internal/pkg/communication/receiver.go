package communication

import (
	"asvsoft/internal/pkg/proto"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type Receiver struct {
	rc io.ReadCloser
}

func NewReceiver(rc io.ReadCloser) *Receiver {
	return &Receiver{
		rc: rc,
	}
}

func (r *Receiver) Recieve() (*proto.Message, error) {
	rawData, err := proto.Read(r.rc)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	log.Debugf("raw received data: %+v", rawData)

	msg, err := proto.Unpack(rawData)
	if err != nil {
		return nil, fmt.Errorf("unpack failed: %v", err)
	}

	log.Infof("received data: %v", msg)

	return msg, nil
}

func (r *Receiver) Close() error {
	return r.rc.Close()
}
