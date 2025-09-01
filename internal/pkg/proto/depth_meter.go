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

func (dmd DepthMeterData) String() string {
	type _DepthMeterData DepthMeterData
	return fmt.Sprintf("%+v", _DepthMeterData(dmd))
}

const (
	depthMeterPaylodSizeModeA = 12
)

func (dmd *DepthMeterData) Pack(msgID MessageID) ([]byte, error) {
	var buf *bytes.Buffer

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, depthMeterPaylodSizeModeA))

		err := encoder.NewEncoder(buf).Encode(dmd.ID, dmd.SystemTime, dmd.Distance, dmd.Status, dmd.Strength, dmd.Precision)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packDepthMeterData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), nil
}

func (dmd *DepthMeterData) Unpack(b []byte, msgID MessageID) error {
	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(b)))

		err := enc.Decode(&dmd.ID, &dmd.SystemTime, &dmd.Distance, &dmd.Status, &dmd.Strength, &dmd.Precision)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("packDepthMeterData is not implemented for this message ID: %x", msgID))
	}

	return nil
}
