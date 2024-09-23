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

var (
	headerBuff     = make([]byte, 3)
	tmpBuff        = make([]byte, 1)
	idBuff         = make([]byte, 1)
	tsBuff         = make([]byte, 4)
	paylodSizeBuff = make([]byte, 1)
	checkSumBuff   = make([]byte, 2)
)

var buffPool = make(map[int][]byte)

// buffFromPool возвращает из пула слайс размера size, что позволяет оптимизирует
// аллокацию памяти при каждом вызове метода Read. Длина сообщений протокола  постоянна,
// поэтому нет необходимости следить за размером пула.
func buffFromPool(size int) []byte {
	buff, ok := buffPool[size]
	if !ok {
		buffPool[size] = make([]byte, size)
		buff = buffPool[size]
	}

	return buff
}

// Read ищет в потоке принимаемых байтов синхронизовачный заголовок
// и затем вычитает фрейм протокола. Возвращает полученный фрейм и ошибку.
func Read(r io.Reader) ([]byte, error) {
	return ReadWithLimit(r, DefaultReadRetries)
}

// ReadWithLimit аналогично Read, но с возможностью указать лимит не по умолчанию.
func ReadWithLimit(r io.Reader, limit int) ([]byte, error) {
	rawData := []byte{}

	_, err := r.Read(headerBuff)
	if err != nil {
		return nil, fmt.Errorf("proto.Read failed: %w", err)
	}

	for retries := limit; len(rawData) == 0 && retries > 0; retries-- {
		if !isEqual(syncFramePart, headerBuff) {
			headerBuff[0] = headerBuff[1]
			headerBuff[1] = headerBuff[2]

			_, err = r.Read(tmpBuff)
			if err != nil {
				return nil, fmt.Errorf("proto.Read failed: %w", err)
			}

			headerBuff[2] = tmpBuff[0]

			continue
		}

		err := readFrameParts(r, idBuff, tsBuff, paylodSizeBuff)
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}

		payloadSize := int(paylodSizeBuff[0])

		payloadBuff := buffFromPool(payloadSize)

		err = readFrameParts(r, payloadBuff, checkSumBuff)
		if err != nil {
			return nil, fmt.Errorf("proto.Read failed: %w", err)
		}

		rawData = buffFromPool(payloadSize + serviceFramePartSize)

		mergeSlices(
			rawData,
			syncFramePart,
			idBuff,
			tsBuff,
			paylodSizeBuff,
			payloadBuff,
			checkSumBuff,
		)
	}

	if len(rawData) == 0 {
		return nil, fmt.Errorf("frame not found after %d bytes reading", limit)
	}

	return rawData, nil
}

func isEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for idx := 0; idx < len(a); idx++ {
		if a[idx] != b[idx] {
			return false
		}
	}

	return true
}

func readFrameParts(r io.Reader, parts ...[]byte) error {
	for _, part := range parts {
		_, err := r.Read(part)
		if err != nil {
			return err
		}
	}

	return nil
}

func mergeSlices(dst []byte, slices ...[]byte) {
	idx := 0

	for _, sl := range slices {
		for _, b := range sl {
			dst[idx] = b
			idx++
		}
	}
}
