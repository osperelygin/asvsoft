// Package depthmeter предоставляет подкоманду depthmeter
package depthmeter

import (
	"asvsoft/internal/pkg/encoder"
	"asvsoft/internal/pkg/proto"
	"bytes"
	"fmt"
	"io"
)

const (
	frameHeaderSize = 3
	totalFrameSize  = 16
)

var (
	// frameHeader последовательность байт протокола для синхронизации
	frameHeader = []byte{0x57, 0x00, 0xff}
	// rawData преаллоцированный массив байт для чтения данных протокола
	rawData = make([]byte, 2*totalFrameSize)
)

type DepthMeter struct {
	r io.Reader
}

func New(r io.Reader) *DepthMeter {
	return &DepthMeter{
		r: r,
	}
}

func (dm *DepthMeter) ReadMeasure() (*proto.DepthMeterData, error) {
	_, err := dm.r.Read(rawData)
	if err != nil {
		return nil, fmt.Errorf("cannot read bytes from port: %v", err)
	}

	start := bytes.Index(rawData, frameHeader)
	if start == -1 {
		return nil, fmt.Errorf("cannot find frame header: %v", rawData)
	}

	if start+totalFrameSize > 2*totalFrameSize {
		return nil, fmt.Errorf("cannot parse binary data: %v", rawData)
	}

	sum := 0
	payload := rawData[start : start+totalFrameSize]

	// log.Println(payload)

	for idx := 0; idx < totalFrameSize-1; idx++ {
		sum += int(payload[idx])
	}

	if sum%256 != int(payload[totalFrameSize-1]) {
		return nil, fmt.Errorf("check sum missmatch")
	}

	measure := proto.DepthMeterData{}
	d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(payload)))

	_, err = d.Discard(frameHeaderSize)
	if err != nil {
		return nil, fmt.Errorf("cannot discard frame header: %v", err)
	}

	err = d.Decode(&measure.ID, &measure.SystemTime, &measure.Distance, &measure.Status, &measure.Strength, &measure.Precision)
	if err != nil {
		return nil, fmt.Errorf("cannot decode measure: %v", err)
	}

	d.Close()

	return &measure, nil
}
