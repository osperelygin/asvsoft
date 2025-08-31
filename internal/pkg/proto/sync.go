package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"io"
)

type SyncData uint32

func (sd *SyncData) Pack(msgID MessageID) ([]byte, error) {
	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, 4)))

	err := enc.Encode(sd)
	if err != nil {
		return nil, err
	}

	return enc.Bytes(), nil
}

func (sd *SyncData) Unpack(b []byte, msgID MessageID) error {
	d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(b)))
	v, err := d.U32()
	if err != nil {
		return err
	}

	*sd = SyncData(v)

	return nil
}
