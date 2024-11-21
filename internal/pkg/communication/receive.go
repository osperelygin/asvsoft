package communication

import (
	"asvsoft/internal/pkg/proto"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

func Recieve(r io.Reader) (any, error) {
	rawData, err := proto.Read(r)
	if err != nil {
		return nil, fmt.Errorf("read failed: %v", err)
	}

	log.Debugf("raw received data: %+v", rawData)

	data, err := proto.Unpack(rawData)
	if err != nil {
		return nil, fmt.Errorf("unpack failed: %v", err)
	}

	log.Infof("received data: %+v", data)

	return data, nil
}
