package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

// CameraData - данные, полученные после обработки модуля камеры
type CameraData struct {
	// Углы ориентации в 0.0001 град
	Yaw, Pitch, Roll int16
}

const (
	cameraDataSizeModeA = 6
)

func packCameraData(in *CameraData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, cameraDataSizeModeA))
		enc := encoder.NewEncoder(buf)

		err = enc.Encode(in.Yaw, in.Pitch, in.Roll)
		if err != nil {
			return nil, err
		}

	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func unpackCameraData(in []byte, msgID MessageID) (out *CameraData, err error) {
	out = new(CameraData)

	switch msgID {
	case WritingModeA:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		err = dec.Decode(&out.Yaw, &out.Pitch, &out.Roll)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
