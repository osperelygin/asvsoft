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
	// RawImage сырое кодированное изображение
	RawImage []byte
}

func (cd CameraData) String() string {
	return fmt.Sprintf(
		"{Yaw:%d,Pitch:%d,Roll:%d,RawImageLen:%d}",
		cd.Yaw, cd.Pitch, cd.Roll, len(cd.RawImage),
	)
}

const (
	cameraDataSizeModeA = 6
)

func packCameraData(in *CameraData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		res []byte
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

		res = buf.Bytes()
	case WritingModeB:
		res = in.RawImage
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return res, err
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
	case WritingModeB:
		out.RawImage = in
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
