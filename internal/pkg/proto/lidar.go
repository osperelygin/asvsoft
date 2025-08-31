package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

const (
	pointNums            = 12
	lidarPaylodSizeModeA = 2 + 2 + pointNums*3 + 2 + 2
)

type Point struct {
	Distance  uint16
	Intensity uint8
}

type LidarData struct {
	// Speed скорость вращения лидара в град/с
	Speed uint16
	// StartAngle начальный угол точек пакета в 0.01 град/с
	StartAngle uint16
	// Points массив точек измерения
	Points [pointNums]Point
	// EndAngle конечный угол точек пакета в 0.01 град/с
	EndAngle uint16
	// Timestamp
	Timestamp uint16
}

func (ld LidarData) String() string {
	type _LidarData LidarData
	return fmt.Sprintf("%+v", _LidarData(ld))
}

func (ld *LidarData) Pack(msgID MessageID) ([]byte, error) {
	var buf *bytes.Buffer

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, lidarPaylodSizeModeA))
		enc := encoder.NewEncoder(buf)

		err := enc.Encode(ld.Speed, ld.StartAngle)
		if err != nil {
			return nil, err
		}

		for _, point := range ld.Points {
			err = enc.Encode(point.Distance, point.Intensity)
			if err != nil {
				return nil, err
			}
		}

		err = enc.Encode(ld.EndAngle, ld.Timestamp)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), nil
}

func (ld *LidarData) Unpack(b []byte, msgID MessageID) error {
	switch msgID {
	case WritingModeA:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(b)))

		err := dec.Decode(&ld.Speed, &ld.StartAngle)
		if err != nil {
			return err
		}

		for i := range ld.Points {
			err = dec.Decode(&ld.Points[i].Distance, &ld.Points[i].Intensity)
			if err != nil {
				return err
			}
		}

		err = dec.Decode(&ld.EndAngle, &ld.Timestamp)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return nil
}
