package encoder

import (
	"asvsoft/internal/pkg/common"
	"bufio"
	"fmt"
	"io"
	"sync"
)

var readBufferPool = sync.Pool{
	New: func() any {
		return bufio.NewReaderSize(nil, 32*common.KB)
	},
}

// A Decoder reads and decodes binary values from an input stream.
type Decoder struct {
	r            *bufio.Reader
	c            io.Closer
	numBytesRead int
}

// NewDecoder returns a new decoder that reads from r.
func NewDecoder(r io.ReadCloser) *Decoder {
	rd := readBufferPool.Get().(*bufio.Reader)
	rd.Reset(r)

	return &Decoder{
		r: rd,
		c: r,
	}
}

func (dec *Decoder) Read(p []byte) (int, error) {
	n, err := dec.r.Read(p)
	dec.numBytesRead += n
	return n, err
}

// Discard skips the next n bytes, returning the number of bytes discarded.
func (dec *Decoder) Discard(n int) (discarded int, err error) {
	discarded, err = dec.r.Discard(n)
	dec.numBytesRead += discarded
	return discarded, err
}

// Close closes reader
func (dec *Decoder) Close() error {
	dec.r.Reset(nil)
	readBufferPool.Put(dec.r)

	if dec.c == nil {
		return nil
	}

	return dec.c.Close()
}

// Slice reads exactly n and returns slice of byte
func (dec *Decoder) Slice(n int) ([]byte, error) {
	buf := make([]byte, n)
	n, err := io.ReadFull(dec.r, buf)
	dec.numBytesRead += n
	return buf, err
}

// NumBytesRead number of bytes read
func (dec *Decoder) NumBytesRead() int {
	return dec.numBytesRead
}

// Decode reads the next binary value from its
// input and stores it in the value pointed to by v.
func (dec *Decoder) Decode(values ...any) error {
	var err error
	for _, untyped := range values {
		switch v := untyped.(type) {
		case *uint8:
			*v, err = dec.U8()
		case *uint16:
			*v, err = dec.U16()
		case *common.Uint24:
			*v, err = dec.U24()
		case *uint32:
			*v, err = dec.U32()
		case *int32:
			*v, err = dec.I32()
		case *int16:
			*v, err = dec.I16()
		case *[]byte:
			var n int
			n, err = io.ReadFull(dec.r, *v)
			dec.numBytesRead += n
		default:
			panic(fmt.Sprintf("Decode is not implemented for this type (%T)", v))
		}

		if err != nil {
			return err
		}
	}

	return nil
}

// U8 reads and returns a single byte
func (dec *Decoder) U8() (uint8, error) {
	b, err := dec.r.ReadByte()
	if err == nil {
		dec.numBytesRead++
	}

	return b, err
}

// U16 reads and returns two bytes
// func (dec *Decoder) U16() (uint16, error) {
// 	c1, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c2, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return uint16(c1) | uint16(c2)<<8, nil
// }

// U16 reads and returns two bytes
func (dec *Decoder) U16() (uint16, error) {
	return decodeBytes[uint16](dec)
}

// U24 reads and returns three bytes
// func (dec *Decoder) U24() (uint24, error) {
// 	c1, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c2, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c3, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return uint24(c1) | uint24(c2)<<8 | uint24(c3)<<16, nil
// }

// U24 reads and returns three bytes
func (dec *Decoder) U24() (common.Uint24, error) {
	return decodeBytes[common.Uint24](dec)
}

// U32 reads and returns four bytes
// func (dec *Decoder) U32() (uint32, error) {
// 	c1, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c2, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c3, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	c4, err := dec.U8()
// 	if err != nil {
// 		return 0, err
// 	}

// 	return uint32(c1) | uint32(c2)<<8 | uint32(c3)<<16 | uint32(c4)<<24, nil
// }

// U32 reads and returns four bytes
func (dec *Decoder) U32() (uint32, error) {
	return decodeBytes[uint32](dec)
}

func (dec *Decoder) I16() (int16, error) {
	return decodeSignedBytes[int16](dec)
}

func (dec *Decoder) I32() (int32, error) {
	return decodeSignedBytes[int32](dec)
}

func decodeBytes[T ubytes](d *Decoder) (res T, err error) {
	n := bytesOf(res)

	for i := 0; i < n; i += 1 {
		c, err := d.U8()
		if err != nil {
			return res, err
		}

		res |= T(c) << (i * 8)
	}

	return res, nil
}

func decodeSignedBytes[T sbytes](d *Decoder) (res T, err error) {
	n := bytesOf(res)

	for i := 0; i < n-1; i += 1 {
		c, err := d.U8()
		if err != nil {
			return res, err
		}

		res |= T(c) << (i * 8)
	}

	c, err := d.U8()
	if err != nil {
		return res, err
	}

	res |= T(c&127) << ((n - 1) * 8)

	if c&128 == 128 {
		res *= -1
	}

	return res, nil
}
