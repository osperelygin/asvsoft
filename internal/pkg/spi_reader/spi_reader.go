// Package spireader ...
package spireader

import (
	"bytes"

	"github.com/stianeikeland/go-rpio/v4"
)

type Reader interface {
	Read(n int) []byte
}

type SPIReader struct {
}

func (r *SPIReader) Read(n int) []byte {
	return rpio.SpiReceive(n)
}

func NewBytesReaderAdapter(b *bytes.Reader) *BytesReaderAdapter {
	return &BytesReaderAdapter{b: b}
}

type BytesReaderAdapter struct {
	b *bytes.Reader
}

func (r *BytesReaderAdapter) Read(n int) []byte {
	b := make([]byte, n)
	_, _ = r.b.Read(b)

	return b
}

// func Read(r io.Reader) ([]byte, error) {
// 	headerBuff := rpio.SpiReceive(3)

// 	if !isEqual(syncFramePart, headerBuff) {
// 		headerBuff[0] = headerBuff[1]
// 		headerBuff[1] = headerBuff[2]
// 		headerBuff[2] = rpio.SpiReceive(1)[1]
// 	}

// 	id := rpio.SpiReceive(1)
// 	ts := rpio.SpiReceive(4)
// 	payloadSize := int(rpio.SpiReceive(1)[0])
// 	payload := rpio.SpiReceive(payloadSize)
// 	checkSum := rpio.SpiReceive(checkSumSize)

// 	rawData := make([]byte, 0, payloadSize+serviceFramePartSize)
// 	rawData = append(rawData, id...)
// 	rawData = append(rawData, ts...)
// 	rawData = append(rawData, byte(payloadSize))
// 	rawData = append(rawData, payload...)
// 	rawData = append(rawData, checkSum...)

// 	return rawData, nil
// }
