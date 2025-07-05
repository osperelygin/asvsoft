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

// Receive читает данный из r.rc и распаковывает пакет в сообщение.
func (r *Receiver) Receive() (proto.Message, error) {
	var msg proto.Message

	rawData, err := proto.Read(r.rc)
	if err != nil {
		return msg, fmt.Errorf("read msg failed: %v", err)
	}

	log.Debugf("raw received msg: %+v", rawData)

	err = msg.Unmarshal(rawData)
	if err != nil {
		return msg, fmt.Errorf("unmarshal msg failed: %v", err)
	}

	return msg, nil
}

func (r *Receiver) Close() error {
	return r.rc.Close()
}
