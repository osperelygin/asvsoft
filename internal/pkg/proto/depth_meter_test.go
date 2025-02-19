package proto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var depthMeterData = &DepthMeterData{
	ID:         1,
	SystemTime: uint32(time.Now().Unix()),
	Distance:   1000,
	Status:     1,
	Strength:   1 << 10,
	Precision:  1 << 5,
}

func TestPackDepthMeterDataSuccess(t *testing.T) {
	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		var sentMsg Message

		msgBytes, err := sentMsg.Marshal(depthMeterData, DepthMeterModuleID, WritingModeA)
		require.NoError(t, err)

		var receivedMsg Message

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
