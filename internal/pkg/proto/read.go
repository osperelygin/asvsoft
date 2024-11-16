package proto

import (
	"bytes"
	"fmt"
	"io"
)

// Read ищет в потоке принимаемых байтов синхронизовачный заголовок
// и затем вычитает фрейм протокола. Возвращает полученный фрейм и ошибку.
func Read(r io.Reader) ([]byte, error) {
	return ReadWithLimit(r, defaultReadRetries)
}

// ReadWithLimit аналогично Read, но с возможностью указать лимит не по умолчанию.
func ReadWithLimit(r io.Reader, limit int) ([]byte, error) {
	var (
		rawData []byte
		svcBuff = make([]byte, serviceBytesSize)
	)

	_, err := r.Read(svcBuff[:headerSize])
	if err != nil {
		return nil, fmt.Errorf("proto.Read failed: %w", err)
	}

	for retries := limit; len(rawData) == 0 && retries > 0; retries-- {
		if !bytes.Equal(header, svcBuff[:headerSize]) {
			copy(svcBuff[0:headerSize-1], svcBuff[1:headerSize])

			_, err = r.Read(svcBuff[headerSize-1 : headerSize])
			if err != nil {
				return nil, fmt.Errorf("proto.Read failed: %w", err)
			}

			continue
		}

		_, err = r.Read(svcBuff[headerSize:payloadFirstByte])
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}

		payloadSize := svcBuff[payloadFirstByte-1]

		rawData = make([]byte, serviceBytesSize+payloadSize)
		copy(rawData, svcBuff[:payloadFirstByte])

		_, err = r.Read(rawData[payloadFirstByte:])
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("frame not found after %d bytes reading", limit)
	}

	return rawData, nil
}
