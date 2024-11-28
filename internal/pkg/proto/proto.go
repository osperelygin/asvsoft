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
	ReadingModeA MessageID = 0x11 + iota
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
	timestampSize    = 8
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

const (
	dummySystemByte byte = 0xFF
)

var header = []byte{0xFA, 0xFA}

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

type Message struct {
	ModuleID    ModuleID
	MsgID       MessageID
	Timestamp   uint64
	PayloadSize uint8
	Payload     any
	CheckSum    uint8
}

func (m Message) String() string {
	return fmt.Sprintf(
		"{moduleID:%#X,msgID:%#X,ts:%d,payloadSize:%d,payload:%+v,checksum: %#X}",
		m.ModuleID, m.MsgID, m.Timestamp, m.PayloadSize, m.Payload, m.CheckSum,
	)
}

// Marshal ...
func (m *Message) Marshal(data any, moduleID ModuleID, msgID MessageID) ([]byte, error) {
	var (
		err     error
		payload []byte
	)

	m.ModuleID = moduleID
	m.MsgID = msgID
	m.Payload = data

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

	m.PayloadSize = uint8(len(payload))
	m.Timestamp = uint64(time.Now().UnixMilli())

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, serviceBytesSize+int(m.PayloadSize))))

	err = enc.Encode(header, dummySystemByte, uint8(moduleID), uint8(msgID), m.Timestamp, m.PayloadSize, payload)
	if err != nil {
		return nil, err
	}

	m.CheckSum = crc8.ChecksumSMBus(enc.Bytes()[headerSize:])

	err = enc.Encode(m.CheckSum)
	if err != nil {
		return nil, err
	}

	return enc.Bytes(), nil
}

// Unmarshal ...
func (m *Message) Unmarshal(data []byte) error {
	dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(data)))
	defer dec.Close()

	// Пропускаем байты синхронизации
	_, err := dec.Discard(headerSize)
	if err != nil {
		return err
	}

	var (
		systemByte      uint8
		moduleID, msgID uint8
	)

	err = dec.Decode(&systemByte, &moduleID, &msgID, &m.Timestamp, &m.PayloadSize)
	if err != nil {
		return err
	}

	_ = systemByte

	rawPayload, err := dec.Slice(int(m.PayloadSize))
	if err != nil {
		return err
	}

	err = dec.Decode(&m.CheckSum)
	if err != nil {
		return err
	}

	if m.CheckSum != crc8.ChecksumSMBus(data[headerSize:len(data)-checkSumSize]) {
		return fmt.Errorf("check sum missmatch")
	}

	m.ModuleID = ModuleID(moduleID)
	m.MsgID = MessageID(msgID)

	switch m.ModuleID {
	case DepthMeterModuleID:
		m.Payload, err = unpackDepthMeterData(rawPayload, m.MsgID)
	case LidarModuleID:
		m.Payload, err = unpackLidarData(rawPayload, m.MsgID)
	case IMUModuleID:
		m.Payload, err = unpackIMUData(rawPayload, m.MsgID)
	case GNSSModuleID:
		m.Payload, err = unpackGNSSData(rawPayload, m.MsgID)
	case CheckModuleID:
		m.Payload, err = unpackCheckData(rawPayload, m.MsgID)
	default:
		panic(fmt.Sprintf("Unpack is not implemented for this addr (%x)", moduleID))
	}

	return err
}
