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

	err = enc.Encode(header, uint8(moduleID), uint8(msgID), dummySystemByte, ts, uint8(payloadSize), payload)
	if err != nil {
		return nil, err
	}

	checkSum := crc16.ChecksumCCITT(enc.Bytes()[headerSize:])
	if err = enc.Encode(checkSum); err != nil {
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

	err = dec.Decode(&rawModuleID, &rawMsgID, &systemByte, &ts, &payloadSize)
	if err != nil {
		return nil, err
	}

	_ = systemByte

	payload, err := dec.Slice(int(payloadSize))
	if err != nil {
		return nil, err
	}

	var checkSum uint16

	err = dec.Decode(&checkSum)
	if err != nil {
		return nil, err
	}

	if checkSum != crc16.ChecksumCCITT(data[headerSize:len(data)-checkSumSize]) {
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
