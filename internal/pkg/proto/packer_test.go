package proto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	t.Run("успешное чтение фрейма протокола из потока байтов", func(t *testing.T) {
		packer := Packer{}

		noiseBytes := []byte{0, 0, 0, syncFramePart[0], syncFramePart[1], 0, 0}

		packedData, err := packer.Pack(depthMeterData, DepthMeterModuleAddr, WritingModeA)
		assert.NoError(t, err)

		rawData := make([]byte, 0, len(packedData)+len(noiseBytes))
		rawData = append(rawData, noiseBytes...)
		rawData = append(rawData, packedData...)

		b, err := Read(bytes.NewReader(rawData), 1<<10)
		assert.NoError(t, err)
		assert.Equal(t, packedData, b)
	})

	t.Run("отсутствие фрейма в потоке байтов", func(t *testing.T) {
		emptyFlow := make([]byte, 1<<11)

		b, err := Read(bytes.NewReader(emptyFlow), 1<<10)
		assert.Nil(t, b)
		assert.Error(t, err)
	})
}

// func TestSPIReader(t *testing.T) {
// 	t.Run("успешное чтение фрейма протокола из потока байтов", func(t *testing.T) {
// 		packer := Packer{}

// 		noiseBytes := []byte{0, 0, 0, syncFramePart[0], syncFramePart[1], 0, 0}

// 		packedData, err := packer.Pack(depthMeterData, DepthMeterModuleAddr, WritingModeA)
// 		assert.NoError(t, err)

// 		rawData := make([]byte, 0, len(packedData)+len(noiseBytes))
// 		rawData = append(rawData, noiseBytes...)
// 		rawData = append(rawData, packedData...)

// 		r := spireader.NewBytesReaderAdapter(bytes.NewReader(rawData))

// 		b, err := Read(r, 1<<10)
// 		assert.NoError(t, err)
// 		assert.Equal(t, packedData, b)
// 	})

// 	t.Run("отсутствие фрейма в потоке байтов", func(t *testing.T) {
// 		emptyFlow := make([]byte, 1 << 11)

// 		r := spireader.NewBytesReaderAdapter(bytes.NewReader(emptyFlow))

// 		b, err := Read(r, 1<<10)
// 		assert.Nil(t, b)
// 		assert.Error(t, err)
// 	})
// }
