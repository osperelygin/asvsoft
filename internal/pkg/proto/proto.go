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

// Message contain msg meta information and payload
type Message struct {
	ModuleID    ModuleID
	MsgID       MessageID
	SystemTime  uint32
	PayloadSize uint8
	Payload     any
	CheckSum    uint8
}

func (m Message) String() string {
	return fmt.Sprintf(
		"{moduleID:%#X,msgID:%#X,ts:%d,payloadSize:%d,payload:%+v,checksum: %#X}",
		m.ModuleID, m.MsgID, m.SystemTime, m.PayloadSize, m.Payload, m.CheckSum,
	)
}

// Marshal ...
func (m *Message) Marshal(data any, moduleID ModuleID, msgID MessageID) ([]byte, error) {
	var (
		err        error
		rawPayload []byte
	)

	m.ModuleID = moduleID
	m.MsgID = msgID
	m.Payload = data

	switch msgID {
	case SyncRequest:
		// just chill
	case SyncResponse:
		enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, 4)))

		err = enc.Encode(data.(uint32))
		if err != nil {
			return nil, err
		}

		rawPayload = enc.Bytes()
	default:
		rawPayload, err = m.pack(data)
	}

	if err != nil {
		return nil, err
	}

	m.PayloadSize = uint8(len(rawPayload))
	m.SystemTime = systemTime()

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, serviceBytesSize+int(m.PayloadSize))))

	err = enc.Encode(header, dummySystemByte, uint8(moduleID), uint8(msgID), m.SystemTime, m.PayloadSize, rawPayload)
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

func (m *Message) pack(data any) ([]byte, error) {
	var (
		rawPayload []byte
		err        error
	)

	switch m.ModuleID {
	case DepthMeterModuleID:
		rawPayload, err = packDepthMeterData(data.(*DepthMeterData), m.MsgID)
	case LidarModuleID:
		rawPayload, err = packLidarData(data.(*LidarData), m.MsgID)
	case IMUModuleID:
		rawPayload, err = packIMUData(data.(*IMUData), m.MsgID)
	case GNSSModuleID:
		rawPayload, err = packGNSSData(data.(*GNSSData), m.MsgID)
	case CameraModuleID:
		rawPayload, err = packCameraData(data.(*CameraData), m.MsgID)
	case CheckModuleID:
		rawPayload, err = packCheckData(data.(*CheckData), m.MsgID)
	default:
		panic(fmt.Sprintf("Pack is not implemented for this addr (%x)", m.ModuleID))
	}

	return rawPayload, err
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

	if m.CheckSum != crc8.ChecksumSMBus(data[headerSize:len(data)-checkSumSize]) {
		return fmt.Errorf("check sum missmatch")
	}

	m.ModuleID = ModuleID(moduleID)
	m.MsgID = MessageID(msgID)

	switch m.MsgID {
	case SyncRequest:
	case SyncResponse:
		d := encoder.NewDecoder(io.NopCloser(bytes.NewReader(rawPayload)))
		defer d.Close()
		m.Payload, err = d.U32()
	default:
		err = m.unpack(rawPayload)
	}

	return err
}

func (m *Message) unpack(rawPayload []byte) error {
	var err error

	switch m.ModuleID {
	case DepthMeterModuleID:
		m.Payload, err = unpackDepthMeterData(rawPayload, m.MsgID)
	case LidarModuleID:
		m.Payload, err = unpackLidarData(rawPayload, m.MsgID)
	case IMUModuleID:
		m.Payload, err = unpackIMUData(rawPayload, m.MsgID)
	case GNSSModuleID:
		m.Payload, err = unpackGNSSData(rawPayload, m.MsgID)
	case CameraModuleID:
		m.Payload, err = unpackCameraData(rawPayload, m.MsgID)
	case CheckModuleID:
		m.Payload, err = unpackCheckData(rawPayload, m.MsgID)
	default:
		panic(fmt.Sprintf("Unpack is not implemented for this addr (%x)", m.ModuleID))
	}

	return err
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
