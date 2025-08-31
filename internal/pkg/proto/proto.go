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
	CheckModuleID       ModuleID = 0xF0
	ControlModuleID     ModuleID = 0x01
	RegistratorModuleID ModuleID = 0xFF
)

const (
	RadioTelemetryModuleID ModuleID = 0x21 + iota*0x10
	CommunicationModule
	IMUModuleID
	GNSSModuleID
	NavigationModuleID
	DepthMeterModuleID
	LidarModuleID
	CameraModuleID
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
	SyncRequest MessageID = 0xFA + iota
	SyncResponse
	ResponseOK
	ResponseFail
)

const (
	headerSize       = 2
	sytemByteSize    = 1
	moduleIDSize     = 1
	msgIDSize        = 1
	systemTimeSize   = 4
	payloadBytesSize = 1
	checkSumSize     = 1
)

const serviceBytesSize = headerSize +
	sytemByteSize +
	moduleIDSize +
	msgIDSize +
	systemTimeSize +
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
	return ReadWithLimitV2(r, defaultReadRetries)
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

		d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(svcBuff[payloadFirstByte-payloadBytesSize : payloadFirstByte])))
		defer d.Close()

		var payloadSize uint8

		err := d.Decode(&payloadSize)
		if err != nil {
			return nil, fmt.Errorf("decode payload size failed: %w", err)
		}

		rawData = make([]byte, serviceBytesSize+int(payloadSize))
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

// ReadWithLimitV2 аналогично Read, но с возможностью указать лимит не по умолчанию.
func ReadWithLimitV2(r io.Reader, limit int) ([]byte, error) {
	var (
		rawData []byte
		svcBuff = make([]byte, payloadFirstByte)
	)

	_, err := r.Read(svcBuff)
	if err != nil {
		return nil, fmt.Errorf("proto.Read failed: %w", err)
	}

	for retries := limit; retries > 0 && len(rawData) == 0; retries-- {
		start := bytes.Index(svcBuff, header)
		if start == -1 {
			svcBuff[0] = svcBuff[payloadFirstByte-1]

			_, err = r.Read(svcBuff[1:])
			if err != nil {
				return nil, fmt.Errorf("proto.Read failed: %w", err)
			}

			continue
		}

		cursor := len(svcBuff) - start
		_ = copy(svcBuff[0:cursor], svcBuff[start:])

		_, err = r.Read(svcBuff[cursor:])
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}

		d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(svcBuff[payloadFirstByte-payloadBytesSize : payloadFirstByte])))
		defer d.Close()

		var payloadSize uint8

		err := d.Decode(&payloadSize)
		if err != nil {
			return nil, fmt.Errorf("decode payload size failed: %w", err)
		}

		rawData = make([]byte, serviceBytesSize+int(payloadSize))

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

// Message contain msg meta information and payload
type Message struct {
	ModuleID    ModuleID
	MsgID       MessageID
	SystemTime  uint32
	PayloadSize uint8
	Payload     Packer
	CheckSum    uint8
}

func NewMessage(moduleID ModuleID, msgID MessageID, payload Packer) *Message {
	return &Message{
		ModuleID: moduleID,
		MsgID:    msgID,
		Payload:  payload,
	}
}

func (m Message) String() string {
	return fmt.Sprintf(
		"{moduleID:%#X,msgID:%#X,ts:%d,payloadSize:%d,payload:%s,checksum: %#X}",
		m.ModuleID, m.MsgID, m.SystemTime, m.PayloadSize, m.Payload, m.CheckSum,
	)
}

type Packer interface {
	Pack(msgID MessageID) ([]byte, error)
	Unpack(b []byte, msgID MessageID) error
}

// Marshal ...
func (m *Message) Marshal() ([]byte, error) {
	var (
		err        error
		rawPayload []byte
	)

	switch m.MsgID {
	case SyncRequest, ResponseOK, ResponseFail:
	default:
		rawPayload, err = m.Payload.Pack(m.MsgID)
	}

	if err != nil {
		return nil, err
	}

	m.PayloadSize = uint8(len(rawPayload))
	m.SystemTime = systemTime()

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, serviceBytesSize+int(m.PayloadSize))))

	err = enc.Encode(header, dummySystemByte, uint8(m.ModuleID), uint8(m.MsgID), m.SystemTime, m.PayloadSize, rawPayload)
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

	err = dec.Decode(&systemByte, &moduleID, &msgID, &m.SystemTime, &m.PayloadSize)
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

	m.ModuleID = ModuleID(moduleID)
	m.MsgID = MessageID(msgID)

	checkSum := crc8.ChecksumSMBus(data[headerSize : len(data)-checkSumSize])

	if m.CheckSum != checkSum {
		return fmt.Errorf(
			"check sum missmatch: recieved cs: %#X, calculated cs: %#X, message: %s",
			m.CheckSum, checkSum, m,
		)
	}

	switch m.MsgID {
	case SyncRequest, ResponseOK, ResponseFail:
	default:
		err = m.unpack(rawPayload)
	}

	return err
}

func (m *Message) unpack(rawPayload []byte) error {
	switch m.ModuleID {
	case DepthMeterModuleID:
		m.Payload = &DepthMeterData{}
	case LidarModuleID:
		m.Payload = &LidarData{}
	case IMUModuleID:
		m.Payload = &IMUData{}
	case GNSSModuleID:
		m.Payload = &GNSSData{}
	case CameraModuleID:
		m.Payload = &CameraData{}
	case CheckModuleID:
		m.Payload = &CheckData{}
	default:
		panic(fmt.Sprintf("Unpack is not implemented for this addr (%x)", m.ModuleID))
	}

	return m.Payload.Unpack(rawPayload, m.MsgID)
}

var startStamp = time.Now().UnixMilli()

func SetStartStamp(stamp uint32) {
	startStamp = int64(stamp) * 1000
}

func GetStartStamp() uint32 {
	return uint32(startStamp / 1000)
}

func systemTime() uint32 {
	return uint32(time.Now().UnixMilli() - startStamp)
}

// lock-free implementation for updating startStamp when now-startStamp > math.Uint32
//
//
// func init() {
// 	startStamp.Store(time.Now().UnixMilli())
// }

// var startStamp atomic.Int64

// func systemTime() uint32 {
// 	now := time.Now().UnixMilli()
// 	start := startStamp.Load()

// 	systemTime := now - start
// 	if systemTime <= math.MaxUint32 {
// 		return uint32(systemTime)
// 	}

// 	if startStamp.CompareAndSwap(start, now) {
// 		return 0
// 	}

// 	start = startStamp.Load()

// 	systemTime = now - start
// 	if systemTime < 0 {
// 		return 0
// 	}

// 	return uint32(systemTime)
// }
