package proto

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

var (
	buffPool = sync.Pool{
		New: func() any {
			buff := make([]byte, defaultBuffSize)
			return &buff
		},
	}

	svcBuffPool = sync.Pool{
		New: func() any {
			buff := make([]byte, servicePartSize)
			return &buff
		},
	}
)

// Read ищет в потоке принимаемых байтов синхронизовачный заголовок
// и затем вычитает фрейм протокола. Возвращает полученный фрейм и ошибку.
func Read(r io.Reader) ([]byte, error) {
	return ReadWithLimit(r, defaultReadRetries)
}

// ReadWithLimit аналогично Read, но с возможностью указать лимит не по умолчанию.
func ReadWithLimit(r io.Reader, limit int) ([]byte, error) {
	var frame []byte

	svcBuff := *svcBuffPool.Get().(*[]byte)
	defer svcBuffPool.Put(&svcBuff)

	_, err := r.Read(svcBuff[:3])
	if err != nil {
		return nil, fmt.Errorf("proto.Read failed: %w", err)
	}

	for retries := limit; len(frame) == 0 && retries > 0; retries-- {
		if !bytes.Equal(syncFramePart, svcBuff[:3]) {
			copy(svcBuff[0:2], svcBuff[1:3])

			_, err = r.Read(svcBuff[2:3])
			if err != nil {
				return nil, fmt.Errorf("proto.Read failed: %w", err)
			}

			continue
		}

		_, err = r.Read(svcBuff[3:9])
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}

		payloadBuff := *buffPool.Get().(*[]byte)
		defer buffPool.Put(&payloadBuff)

		payloadSize := svcBuff[8]
		payload := payloadBuff[:payloadSize]

		for _, part := range [][]byte{payload, svcBuff[9:]} {
			_, err = r.Read(part)
			if err != nil {
				return nil, fmt.Errorf("proto.Read failed: %w", err)
			}
		}

		buff := *buffPool.Get().(*[]byte)
		defer buffPool.Put(&buff)

		buff = buff[:0]
		frameBuff := bytes.NewBuffer(buff)

		for _, part := range [][]byte{svcBuff[:9], payload, svcBuff[9:]} {
			_, err := frameBuff.Write(part)
			if err != nil {
				return nil, fmt.Errorf("cannot write to destination: %w", err)
			}
		}

		frame = frameBuff.Bytes()
	}

	if len(frame) == 0 {
		return nil, fmt.Errorf("frame not found after %d bytes reading", limit)
	}

	return frame, nil
}
