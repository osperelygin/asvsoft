package proto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackLidarDataSuccess(t *testing.T) {
	data := &LidarData{
		Speed:      0xe05,
		StartAngle: 0x8bf,
		Points: [12]Point{
			{Distance: 0x218, Intensity: 0xff},
			{Distance: 0x218, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x219, Intensity: 0xff},
			{Distance: 0x21a, Intensity: 0xff},
			{Distance: 0x21a, Intensity: 0xff},
			{Distance: 0x21a, Intensity: 0xff},
			{Distance: 0x21a, Intensity: 0xff}},
		EndAngle:  0x976,
		Timestamp: 0x513c,
	}

	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		sentMsg := NewMessage(LidarModuleID, WritingModeA, data)

		msgBytes, err := sentMsg.Marshal()
		require.NoError(t, err)

		receivedMsg := new(Message)

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
