package proto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	t.Run("успешное чтение фрейма протокола из потока байтов", func(t *testing.T) {
		packedData, err := Pack(depthMeterData, DepthMeterModuleAddr, WritingModeA)
		assert.NoError(t, err)

		noiseBytes := []byte{0x01, 0x00, 0xFF, syncFramePart[0], syncFramePart[1], 0x05, 0x06}

		rawData := make([]byte, 0, len(packedData)+len(noiseBytes))
		rawData = append(rawData, noiseBytes...)
		rawData = append(rawData, packedData...)

		b, err := Read(bytes.NewReader(rawData))
		assert.NoError(t, err)
		assert.Equal(t, packedData, b)
	})

	t.Run("отсутствие фрейма в потоке байтов", func(t *testing.T) {
		emptyFlow := make([]byte, 1<<11)

		b, err := Read(bytes.NewReader(emptyFlow))
		assert.Nil(t, b)
		assert.Error(t, err)
	})
}

func BenchmarkRead(b *testing.B) {
	packedData, _ := Pack(depthMeterData, DepthMeterModuleAddr, WritingModeA)

	noiseBytes := []byte{0x01, 0x00, 0xFF, syncFramePart[0], syncFramePart[1], 0x05, 0x06}

	rawData := make([]byte, 0, len(packedData)+len(noiseBytes))
	rawData = append(rawData, noiseBytes...)
	rawData = append(rawData, packedData...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Read(bytes.NewReader(rawData))
		if err != nil {
			b.Fatalf("Read return error: %v", err)
		}
	}
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
