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
		"{len(RawImagePart):%d,CurrentChunck:%d,TotalChunckes:%d}",
		len(cd.RawImagePart), cd.CurrentChunck, cd.TotalChunckes,
	)
}

const (
	cameraDataSizeModeA = 6
)

func (cd *CameraData) Pack(msgID MessageID) ([]byte, error) {
	var buf *bytes.Buffer

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, cameraDataSizeModeA))

		err := encoder.NewEncoder(buf).Encode(cd.Yaw, cd.Pitch, cd.Roll)
		if err != nil {
			return nil, err
		}
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, len(cd.RawImagePart)+2))

		err := encoder.NewEncoder(buf).Encode(cd.CurrentChunck, cd.TotalChunckes, cd.RawImagePart)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), nil
}

func (cd *CameraData) Unpack(in []byte, msgID MessageID) error {
	switch msgID {
	case WritingModeA:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		err := dec.Decode(&cd.Yaw, &cd.Pitch, &cd.Roll)
		if err != nil {
			return err
		}
	case WritingModeB:
		dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		cd.RawImagePart = make([]byte, len(in)-2)

		err := dec.Decode(&cd.CurrentChunck, &cd.TotalChunckes, &cd.RawImagePart)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("packLidarData is not implemented for this message ID: %x", msgID))
	}

	return nil
}
