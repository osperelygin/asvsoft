// Package proto предоставляет коммуникацию определенными в унифицированном протоколе
// сообщениями между модулями БКУ
package proto

import (
	"asvsoft/internal/pkg/encoder"
	"asvsoft/pkg/crc8"
	"bytes"
	"fmt"
	"io"
	"time"
)

type ModuleID uint8

const (
	CheckModuleID   ModuleID = 0xF0
	ControlModuleID ModuleID = 0x01
)

const (
	RadioTelemetryModuleID ModuleID = 0x21 + iota*0x10
	CommunicationModule
	IMUModuleID
	GNSSModuleID
	NavigationModuleID
	DepthMeterModuleID
	LidarModuleID
)

type MessageID uint8

const (
	ReadingModeA = 0x11 + iota
	ReadingModeB
	ReadingModeC
	WritingModeA
	WritingModeB
	WritingModeC
)

const (
	headerSize       = 2
	sytemByteSize    = 1
	moduleIDSize     = 1
	msgIDSize        = 1
	timestampSize    = 4
	payloadBytesSize = 1
	checkSumSize     = 1
)

const serviceBytesSize = headerSize +
	sytemByteSize +
	moduleIDSize +
	msgIDSize +
	timestampSize +
	payloadBytesSize +
	checkSumSize

const payloadFirstByte = serviceBytesSize - checkSumSize

const (
	defaultBuffSize    = 512
	defaultReadRetries = 1024
)

var header = []byte{0xFA, 0xFA}

const (
	dummySystemByte byte = 0xFF
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

// Pack ...
func Pack(data any, moduleID ModuleID, msgID MessageID) ([]byte, error) {
	var (
		err     error
		payload []byte
	)

	switch moduleID {
	case DepthMeterModuleID:
		payload, err = packDepthMeterData(data.(*DepthMeterData), msgID)
	case LidarModuleID:
		payload, err = packLidarData(data.(*LidarData), msgID)
	case IMUModuleID:
		payload, err = packIMUData(data.(*IMUData), msgID)
	case GNSSModuleID:
		payload, err = packGNSSData(data.(*GNSSData), msgID)
	case CheckModuleID:
		payload, err = packCheckData(data.(*CheckData), msgID)
	default:
		panic(fmt.Sprintf("Pack is not implemented for this addr (%x)", moduleID))
	}

	if err != nil {
		return nil, err
	}

	payloadSize := len(payload)
	ts := uint32(time.Now().Unix())

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, serviceBytesSize+payloadSize)))

	err = enc.Encode(header, dummySystemByte, uint8(moduleID), uint8(msgID), ts, uint8(payloadSize), payload)
	if err != nil {
		return nil, err
	}

	checkSum := crc8.ChecksumSMBus(enc.Bytes()[headerSize:])

	err = enc.Encode(checkSum)
	if err != nil {
		return nil, err
	}

	return enc.Bytes(), nil
}

// Unpack ...
func Unpack(data []byte) (out any, err error) {
	dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(data)))
	defer dec.Close()

	// Пропускаем байты синхронизации
	_, err = dec.Discard(headerSize)
	if err != nil {
		return nil, err
	}

	var (
		rawModuleID, rawMsgID, systemByte uint8
		ts                                uint32
		payloadSize                       uint8
	)

	err = dec.Decode(&systemByte, &rawModuleID, &rawMsgID, &ts, &payloadSize)
	if err != nil {
		return nil, err
	}

	_ = systemByte

	payload, err := dec.Slice(int(payloadSize))
	if err != nil {
		return nil, err
	}

	var checkSum uint8

	err = dec.Decode(&checkSum)
	if err != nil {
		return nil, err
	}

	if checkSum != crc8.ChecksumSMBus(data[headerSize:len(data)-checkSumSize]) {
		return nil, fmt.Errorf("check sum missmatch")
	}

	moduleID := ModuleID(rawModuleID)
	msgID := MessageID(rawMsgID)

	switch moduleID {
	case DepthMeterModuleID:
		out, err = unpackDepthMeterData(payload, msgID)
	case LidarModuleID:
		out, err = unpackLidarData(payload, msgID)
	case IMUModuleID:
		out, err = unpackIMUData(payload, msgID)
	case GNSSModuleID:
		out, err = unpackGNSSData(payload, msgID)
	case CheckModuleID:
		out, err = unpackCheckData(payload, msgID)
	default:
		panic(fmt.Sprintf("Unpack is not implemented for this addr (%x)", rawModuleID))
	}

	return out, err
}
