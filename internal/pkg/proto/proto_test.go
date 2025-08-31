package proto

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var _depthMeterData = &DepthMeterData{
	ID:         1,
	SystemTime: uint32(time.Now().Unix()),
	Distance:   1000,
	Status:     1,
	Strength:   1 << 10,
	Precision:  1 << 5,
}

func BenchmarkRead(b *testing.B) {
	sentMsg := NewMessage(DepthMeterModuleID, WritingModeA, _depthMeterData)
	msgBytes, _ := sentMsg.Marshal()

	noiseBytes := []byte{0x01, 0x00, 0xFF, header[0], header[1], 0x05, 0x06}

	rawData := make([]byte, 0, len(msgBytes)+len(noiseBytes))
	rawData = append(rawData, noiseBytes...)
	rawData = append(rawData, msgBytes...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Read(bytes.NewReader(rawData))
		if err != nil {
			b.Fatalf("Read return error: %v", err)
		}
	}
}

func TestRead(t *testing.T) {
	t.Run("успешное чтение фрейма протокола из потока байтов", func(t *testing.T) {
		sentMsg := NewMessage(DepthMeterModuleID, WritingModeA, _depthMeterData)

		msgBytes, err := sentMsg.Marshal()
		require.NoError(t, err)

		noiseBytes := []byte{0x01, 0x00, 0xFF, header[0], 0x05, 0x06}

		rawData := make([]byte, 0, len(msgBytes)+len(noiseBytes))
		rawData = append(rawData, noiseBytes...)
		rawData = append(rawData, msgBytes...)

		b, err := Read(bytes.NewReader(rawData))
		require.NoError(t, err)
		require.Equal(t, msgBytes, b, "неожиданное упакованное сообщение")
	})

	t.Run("отсутствие фрейма в потоке байтов", func(t *testing.T) {
		emptyFlow := make([]byte, 1<<11)

		b, err := Read(bytes.NewReader(emptyFlow))
		require.Nil(t, b)
		require.Error(t, err)
	})
}
