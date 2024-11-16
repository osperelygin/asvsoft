package proto

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckDataSuccess(t *testing.T) {
	data := &CheckData{Value: rand.Uint32()}

	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		packedData, err := Pack(data, CheckModuleID, WritingModeA)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)

		assert.Equal(t, data, out.(*CheckData))
	})
}
