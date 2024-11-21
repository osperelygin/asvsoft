// Package lidar предоставляет структуру для получения измерений лидара
package lidar

import (
	"asvsoft/internal/app/ds"
	"asvsoft/internal/pkg/encoder"
	"asvsoft/internal/pkg/measurer"
	"asvsoft/internal/pkg/proto"
	"bytes"
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	verLenByte  = 0x2c
	frameHeader = 0x54
	frameSize   = 1 + 1 + 2 + 2 + 12*3 + 2 + 2 + 1
)

type Lidar struct {
	r io.ReadCloser
	// frameBuff буффер для чтения данных протокола
	frameBuff []byte
	// lastIndex индекс последнего необработанного байта буффера чтения
	lastIndex int
}

func New(r io.ReadCloser) *Lidar {
	return &Lidar{
		r:         r,
		frameBuff: make([]byte, 2*frameSize),
		lastIndex: 0,
	}
}

func (l *Lidar) Measure(_ context.Context) measurer.Measurement {
	return ds.NewMeasurement(l.read())
}

func (l *Lidar) Close() error {
	return l.r.Close()
}

func (l *Lidar) read() (*proto.LidarData, error) {
	var end int
	defer func() {
		if end == 0 {
			l.lastIndex = 0
		}

		l.lastIndex = 2*frameSize - end
		copy(l.frameBuff[:l.lastIndex], l.frameBuff[end:])
	}()

	_, err := l.r.Read(l.frameBuff[l.lastIndex:])
	if err != nil {
		return nil, fmt.Errorf("cannot read bytes from port: %w", err)
	}

	start := bytes.IndexByte(l.frameBuff, frameHeader)
	if start == -1 {
		return nil, fmt.Errorf("cannot find header byte: %v", l.frameBuff)
	}

	if start+frameSize > 2*frameSize {
		return nil, fmt.Errorf("cannot parse binary data: %v", l.frameBuff)
	}

	end = start + frameSize

	frame := l.frameBuff[start:end]

	log.Debugf("[lidar] read frame: %v", frame)

	if frame[1] != verLenByte {
		return nil, fmt.Errorf("unexpected ver len byte: %v", verLenByte)
	}

	checkSum := calcCheckSum(frame[:frameSize-1])
	if checkSum != frame[frameSize-1] {
		return nil, fmt.Errorf("check sum missmatch")
	}

	d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(frame[2:])))
	defer func() {
		err = d.Close()
		if err != nil {
			log.Errorf("cannot close decoder: %v", err)
		}
	}()

	var m proto.LidarData

	err = d.Decode(&m.Speed, &m.StartAngle)
	if err != nil {
		return nil, fmt.Errorf("cannot decode measure: %w", err)
	}

	for i := range m.Points {
		err = d.Decode(&m.Points[i].Distance, &m.Points[i].Intensity)
		if err != nil {
			return nil, fmt.Errorf("cannot decode measure: %w", err)
		}
	}

	err = d.Decode(&m.EndAngle, &m.Timestamp)
	if err != nil {
		return nil, fmt.Errorf("cannot decode measure: %w", err)
	}

	return &m, nil
}
