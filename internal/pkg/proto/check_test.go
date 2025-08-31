package proto

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckDataSuccess(t *testing.T) {
	data := &CheckData{Value: rand.Uint32()}

	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		sentMsg := NewMessage(CheckModuleID, WritingModeA, data)

		msgBytes, err := sentMsg.Marshal()
		require.NoError(t, err)

		receivedMsg := new(Message)

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
