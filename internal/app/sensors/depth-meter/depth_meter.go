// Package depthmeter предоставляет подкоманду depthmeter
package depthmeter

import (
	"asvsoft/internal/pkg/encoder"
	"asvsoft/internal/pkg/proto"
	"bytes"
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

const (
	headerSize = 3
	frameSize  = 16
)

var (
	// frameHeader последовательность байт протокола для синхронизации
	frameHeader = []byte{0x57, 0x00, 0xff}
	// rawData преаллоцированный массив байт для чтения данных протокола
	rawData = make([]byte, 2*frameSize)
)

type DepthMeter struct {
	r io.ReadCloser
}

func New(r io.ReadCloser) *DepthMeter {
	return &DepthMeter{
		r: r,
	}
}

func (dm *DepthMeter) Measure(_ context.Context) (any, error) {
	return dm.measure()
}

func (dm *DepthMeter) Close() error {
	return dm.r.Close()
}

func (dm *DepthMeter) measure() (*proto.DepthMeterData, error) {
	_, err := dm.r.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("cannot read bytes from port: %w", err)
	}

	start := bytes.Index(rawData, frameHeader)
	if start == -1 {
		return nil, fmt.Errorf("cannot find frame header: %v", rawData)
	}

	if start+frameSize > 2*frameSize {
		return nil, fmt.Errorf("cannot parse binary data: %v", rawData)
	}

	frame := rawData[start : start+frameSize]

	log.Debugf("raw read measure: %v", frame)

	sum := 0
	for idx := 0; idx < frameSize-1; idx++ {
		sum += int(frame[idx])
	}

	if sum%256 != int(frame[frameSize-1]) {
		return nil, fmt.Errorf("check sum missmatch")
	}

	d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(frame)))

	defer func() {
		err = d.Close()
		if err != nil {
			log.Errorf("cannot close decoder: %v", err)
		}
	}()

	_, err = d.Discard(headerSize)
	if err != nil {
		return nil, fmt.Errorf("cannot discard frame header: %w", err)
	}

	var measure proto.DepthMeterData

	err = d.Decode(&measure.ID, &measure.SystemTime, &measure.Distance, &measure.Status, &measure.Strength, &measure.Precision)
	if err != nil {
		return nil, fmt.Errorf("cannot decode measure: %w", err)
	}

	return &measure, nil
}
