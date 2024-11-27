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
		var sentMsg Message

		msgBytes, err := sentMsg.Marshal(data, LidarModuleID, WritingModeA)
		require.NoError(t, err)

		var receivedMsg Message

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
