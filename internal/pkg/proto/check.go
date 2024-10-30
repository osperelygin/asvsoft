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

func packCheckData(in *CheckData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, checkDataPayloadSize))
		err = encoder.NewEncoder(buf).Encode(in.Value)
	default:
		panic(fmt.Sprintf("packGNSSData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func unpackCheckData(in []byte, msgID MessageID) (out *CheckData, err error) {
	out = new(CheckData)

	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(&out.Value)
	default:
		panic(fmt.Sprintf("unpackCheckData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
