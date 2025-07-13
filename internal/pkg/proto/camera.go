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
	Yaw, Pitch, Roll             int16
	CurrentChunck, TotalChunckes uint8
	// RawImagePart сырое кодированное изображение
	RawImagePart []byte
}

func (cd CameraData) String() string {
	return fmt.Sprintf(
		"{Yaw:%d,Pitch:%d,Roll:%d,RawImageLen:%d,CurrentChunck:%d,TotalChunckes:%d}",
		cd.Yaw, cd.Pitch, cd.Roll, len(cd.RawImagePart), cd.CurrentChunck, cd.TotalChunckes,
	)
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

		err = encoder.NewEncoder(buf).Encode(in.Yaw, in.Pitch, in.Roll)
		if err != nil {
			return nil, err
		}
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, len(in.RawImagePart)+2))

		err = encoder.NewEncoder(buf).Encode(in.CurrentChunck, in.TotalChunckes, in.RawImagePart)
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
	case WritingModeB:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		out.RawImagePart = make([]byte, len(in)-2)

		err = dec.Decode(&out.CurrentChunck, &out.TotalChunckes, &out.RawImagePart)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
