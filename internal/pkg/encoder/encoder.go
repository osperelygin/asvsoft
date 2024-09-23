// Package encoder занимается кодирование/декодирование бинарных данных
package encoder

import (
	"asvsoft/internal/pkg/common"
	"bytes"
	"fmt"
)

// An Encoder writes binary values to an output stream.
type Encoder struct {
	w *bytes.Buffer // where to send the data
}

// NewEncoder returns a new encoder that writes to w.
func NewEncoder(w *bytes.Buffer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// Encode writes the binary encoding of values to the stream
func (enc *Encoder) Encode(values ...any) error {
	for _, untyped := range values {
		var err error
		switch v := untyped.(type) {
		case uint8:
			err = enc.U8(v)
		case uint16:
			err = enc.U16(v)
		case common.Uint24:
			err = enc.U24(v)
		case uint32:
			err = enc.U32(v)
		case int16:
			err = enc.I16(v)
		case int32:
			err = enc.I32(v)
		case []byte:
			err = enc.Slice(v)
		default:
			panic(fmt.Sprintf("Encode is not implemented for this type (%T)", v))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// Bytes returns encoded bytes
func (enc *Encoder) Bytes() []byte {
	return enc.w.Bytes()
}

// U8 writes a single byte to the stream
func (enc *Encoder) U8(v uint8) error {
	return enc.w.WriteByte(v)
}

// U16 writes two bytes to the stream
// func (enc *Encoder) U16(v uint16) error {
// 	b := make([]byte, 2)
// 	b[0] = byte(v)
// 	b[1] = byte(v >> 8)

// 	_, err := enc.w.Write(b)
// 	return err
// }

// U16 writes two bytes to the stream
func (enc *Encoder) U16(v uint16) error {
	return encodeBytes(enc, v)
}

// U16 writes two bytes to the stream
// func (enc *Encoder) U24(v uint24) error {
// 	b := make([]byte, 3)
// 	b[0] = byte(v)
// 	b[1] = byte(v >> 8)
// 	b[2] = byte(v >> 16)

// 	_, err := enc.w.Write(b)
// 	return err
// }

// U24 writes three bytes to the stream
func (enc *Encoder) U24(v common.Uint24) error {
	return encodeBytes(enc, v)
}

// U32 writes two bytes to the stream
// func (enc *Encoder) U32(v uint32) error {
// 	b := make([]byte, 4)
// 	b[0] = byte(v)
// 	b[1] = byte(v >> 8)
// 	b[2] = byte(v >> 16)
// 	b[3] = byte(v >> 24)

// 	_, err := enc.w.Write(b)
// 	return err
// }

// U32 writes two bytes to the stream
func (enc *Encoder) U32(v uint32) error {
	return encodeBytes(enc, v)
}

// I16 writes two bytes to the stream
func (enc *Encoder) I16(v int16) error {
	return encodeSignedBytes(enc, v)
}

// I32 writes two bytes to the stream
func (enc *Encoder) I32(v int32) error {
	return encodeSignedBytes(enc, v)
}

func encodeBytes[T ubytes](enc *Encoder, v T) error {
	n := bytesOf(v)
	b := make([]byte, 0, n)

	for i := 0; i < n; i++ {
		b = append(b, byte(v>>(i*8)))
	}

	_, err := enc.w.Write(b)

	return err
}

func encodeSignedBytes[T sbytes](enc *Encoder, v T) error {
	n := bytesOf(v)
	b := make([]byte, 0, n)
	isNegative := false

	if v < 0 {
		v *= -1
		isNegative = true
	}

	for i := 0; i < n; i++ {
		b = append(b, byte(v>>(i*8)))
	}

	if isNegative {
		b[n-1] |= 128
	}

	_, err := enc.w.Write(b)

	return err
}

// Slice writes slice to the stream
func (enc *Encoder) Slice(v []byte) error {
	_, err := enc.w.Write(v)
	return err
}
