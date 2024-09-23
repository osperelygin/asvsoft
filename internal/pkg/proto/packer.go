// package proto предоставляет коммуникацию определенными в унифицированном протоколе
// сообщениями между модулями БКУ
package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
	"time"

	"github.com/howeyc/crc16"
)

type Addr uint8

const (
	ControlModuleAddr Addr = iota * 0x10
	DepthMeterModuleAddr
	LidarModuleAddr
	CommunicationModuleAddr
	NavigationModuleAddr
	GNSSModuleAddr
	IMUModuleAddr
)

type MessageID uint8

const (
	ReadingModeA MessageID = iota
	ReadingModeB
	ReadingModeC
	WritingModeA
	WritingModeB
	WritingModeC
)

type Bitmask uint8

const (
	ModuleAddrBitmask Bitmask = 0xF0
	MessageIDBitmask  Bitmask = 0x0F
)

const (
	checkSumSize         = 2
	syncFramePartSize    = 3
	serviceFramePartSize = 11
)

var syncFramePart = []byte{0x57, 0x10, 0xFF}

type Packer struct{}

func NewPacker() Packer {
	return Packer{}
}

// Pack ...
func (p *Packer) Pack(data any, addr Addr, msgID MessageID) ([]byte, error) {
	var err error
	var payload []byte

	switch addr {
	case DepthMeterModuleAddr:
		payload, err = p.packDepthMeterData(data.(*DepthMeterData), msgID)
	case IMUModuleAddr:
		payload, err = p.packIMUData(data.(*IMUData), msgID)
	case GNSSModuleAddr:
		payload, err = p.packGNSSData(data.(*GNSSData), msgID)
	default:
		panic(fmt.Sprintf("Pack is not implemented for this addr (%x)", addr))
	}

	if err != nil {
		return nil, err
	}

	payloadSize := uint8(len(payload))
	id := uint8(addr) | uint8(msgID)
	ts := uint32(time.Now().Unix())

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, serviceFramePartSize+payloadSize)))

	err = enc.Encode(syncFramePart, id, ts, payloadSize, payload)
	if err != nil {
		return nil, err
	}

	checkSum := crc16.ChecksumCCITT(enc.Bytes()[syncFramePartSize:])
	if err = enc.Encode(checkSum); err != nil {
		return nil, err
	}

	return enc.Bytes(), nil
}

// Unpack ...
func (p *Packer) Unpack(data []byte) (out any, err error) {
	dec := encoder.NewDecoder(io.NopCloser(bytes.NewReader(data)))
	defer dec.Close()

	// Пропускаем байты синхронизации
	_, err = dec.Discard(syncFramePartSize)
	if err != nil {
		return nil, err
	}

	var id uint8
	var ts uint32
	var payloadSize uint8

	err = dec.Decode(&id, &ts, &payloadSize)
	if err != nil {
		return nil, err
	}

	payload, err := dec.Slice(int(payloadSize))
	if err != nil {
		return nil, err
	}

	var checkSum uint16

	err = dec.Decode(&checkSum)
	if err != nil {
		return nil, err
	}

	if checkSum != crc16.ChecksumCCITT(data[syncFramePartSize:len(data)-checkSumSize]) {
		return nil, fmt.Errorf("check sum missmatch")
	}

	addr := Addr(id) & Addr(ModuleAddrBitmask)
	msgID := MessageID(id) & MessageID(MessageIDBitmask)

	switch addr {
	case DepthMeterModuleAddr:
		out, err = p.unpackDepthMeterData(payload, msgID)
	case IMUModuleAddr:
		out, err = p.unpackIMUData(payload, msgID)
	case GNSSModuleAddr:
		out, err = p.unpackGNSSData(payload, msgID)
	default:
		panic(fmt.Sprintf("Unpack is not implemented for this addr (%x)", addr))
	}

	return out, err
}

var (
	tmpBuff        = make([]byte, 1)
	idBuff         = make([]byte, 1)
	tsBuff         = make([]byte, 4)
	paylodSizeBuff = make([]byte, 1)
	checkSumBuff   = make([]byte, 2)
)

// Read ищет в потоке принимаемых байтов синхронизовачный заголовок
// и затем вычитает фрейм протокола. Возвращает полученный фрейм и ошибку.
func Read(r io.Reader, limit int) ([]byte, error) {
	var rawData []byte

	headerBuff := make([]byte, 3)
	r.Read(headerBuff)

	for idx := 0; idx < limit && rawData == nil; idx += 1 {
		if !isEqual(syncFramePart, headerBuff) {
			headerBuff[0] = headerBuff[1]
			headerBuff[1] = headerBuff[2]
			r.Read(tmpBuff)
			headerBuff[2] = tmpBuff[0]
			continue
		}

		r.Read(idBuff)
		r.Read(tsBuff)
		r.Read(paylodSizeBuff)

		payloadSize := int(paylodSizeBuff[0])
		// TODO: получать из пула
		payloadBuff := make([]byte, payloadSize)

		r.Read(payloadBuff)
		r.Read(checkSumBuff)

		rawData = make([]byte, 0, payloadSize+serviceFramePartSize)
		rawData = append(rawData, syncFramePart...)
		rawData = append(rawData, idBuff...)
		rawData = append(rawData, tsBuff...)
		rawData = append(rawData, byte(payloadSize))
		rawData = append(rawData, payloadBuff...)
		rawData = append(rawData, checkSumBuff...)
	}

	if rawData == nil {
		return nil, fmt.Errorf("frame not found after %d bytes reading", limit)
	}

	return rawData, nil
}

func isEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx += 1 {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}
