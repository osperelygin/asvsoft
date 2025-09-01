package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

const checkDataPayloadSize = 8

type CheckData struct {
	Value uint32
}

func (cd CheckData) String() string {
	type _CheckData CheckData
	return fmt.Sprintf("%+v", _CheckData(cd))
}

func (cd *CheckData) Pack(msgID MessageID) ([]byte, error) {
	var buf *bytes.Buffer

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, checkDataPayloadSize))

		err := encoder.NewEncoder(buf).Encode(cd.Value)
		if err != nil {
			return nil, err
		}
	default:
		panic(fmt.Sprintf("packGNSSData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), nil
}

func (cd *CheckData) Unpack(in []byte, msgID MessageID) error {
	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))

		err := enc.Decode(&cd.Value)
		if err != nil {
			return err
		}
	default:
		panic(fmt.Sprintf("unpackCheckData is not implemented for this message ID: %x", msgID))
	}

	return nil
}
