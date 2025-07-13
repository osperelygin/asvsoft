package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

const pointNums = 12

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

func (d LidarData) String() string {
	type _LidarData LidarData
	return fmt.Sprintf("%+v", _LidarData(d))
}

type Point struct {
	Distance  uint16
	Intensity uint8
}

const (
	lidarPaylodSizeModeA = 2 + 2 + pointNums*3 + 2 + 2
)

func packLidarData(in *LidarData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, lidarPaylodSizeModeA))
		enc := encoder.NewEncoder(buf)

		err = enc.Encode(in.Speed, in.StartAngle)
		if err != nil {
			return nil, err
		}

		for _, point := range in.Points {
			err = enc.Encode(point.Distance, point.Intensity)
			if err != nil {
				return nil, err
			}
		}

		err = enc.Encode(in.EndAngle, in.Timestamp)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func unpackLidarData(in []byte, msgID MessageID) (out *LidarData, err error) {
	out = new(LidarData)

	switch msgID {
	case WritingModeA:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		err = dec.Decode(&out.Speed, &out.StartAngle)
		if err != nil {
			return nil, err
		}

		for i := range out.Points {
			err = dec.Decode(&out.Points[i].Distance, &out.Points[i].Intensity)
			if err != nil {
				return nil, err
			}
		}

		err = dec.Decode(&out.EndAngle, &out.Timestamp)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
