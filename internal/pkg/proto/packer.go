// Package proto предоставляет коммуникацию определенными в унифицированном протоколе
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

var syncFramePart = []byte{0x57, 0x10, 0xFF}

type Packer struct{}

func NewPacker() *Packer {
	return &Packer{}
}

// Pack ...
func (p *Packer) Pack(data any, addr Addr, msgID MessageID) ([]byte, error) {
	var (
		err     error
		payload []byte
	)

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

	enc := encoder.NewEncoder(bytes.NewBuffer(make([]byte, 0, servicePartSize+payloadSize)))

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

	var (
		id          uint8
		ts          uint32
		payloadSize uint8
	)

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
