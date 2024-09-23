package proto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		packer := Packer{}

		packedData, err := packer.Pack(depthMeterData, DepthMeterModuleAddr, WritingModeA)
		assert.NoError(t, err)

		out, err := packer.Unpack(packedData)
		assert.NoError(t, err)

		assert.Equal(t, depthMeterData, out.(*DepthMeterData))
	})
}
