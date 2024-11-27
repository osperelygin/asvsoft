package proto

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCheckDataSuccess(t *testing.T) {
	data := &CheckData{Value: rand.Uint32()}

	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		packedData, err := Pack(data, CheckModuleID, WritingModeA)
		require.NoError(t, err)

		out, err := Unpack(packedData)
		require.NoError(t, err)

		require.Equal(t, &Message{
			ModuleID:  CheckModuleID,
			MsgID:     WritingModeA,
			Payload:   data,
			Timestamp: out.Timestamp,
			CheckSum:  out.CheckSum,
		}, out)
	})
}
