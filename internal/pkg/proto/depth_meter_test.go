package proto

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackDepthMeterDataSuccess(t *testing.T) {
	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		sentMsg := NewMessage(DepthMeterModuleID, WritingModeA, _depthMeterData)

		msgBytes, err := sentMsg.Marshal()
		require.NoError(t, err)

		receivedMsg := new(Message)

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
