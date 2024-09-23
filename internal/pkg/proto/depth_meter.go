package proto

import (
	"asvsoft/internal/pkg/common"
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

type DepthMeterData struct {
	// Идентификатор сенсора
	ID uint8
	// Системное время
	SystemTime uint32
	// Измеренное расстояние
	Distance common.Uint24
	// Статус измерения
	Status uint8
	// Сила измеренного сигнала
	Strength uint16
	// Точность измеренного сигнала
	Precision uint8
}

const (
	depthMeterPaylodSizeModeA = 12
)

func (p *Packer) packDepthMeterData(in *DepthMeterData, msgID MessageID) ([]byte, error) {
	var buf *bytes.Buffer
	var err error

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, depthMeterPaylodSizeModeA))
		err = encoder.NewEncoder(buf).Encode(in.ID, in.SystemTime, in.Distance, in.Status, in.Strength, in.Precision)
	default:
		panic(fmt.Sprintf("packDepthMeterData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func (p *Packer) unpackDepthMeterData(in []byte, msgID MessageID) (*DepthMeterData, error) {
	var err error
	out := new(DepthMeterData)

	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(&out.ID, &out.SystemTime, &out.Distance, &out.Status, &out.Strength, &out.Precision)
	default:
		panic(fmt.Sprintf("packDepthMeterData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
