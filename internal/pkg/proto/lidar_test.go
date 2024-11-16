package proto

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
		packedData, err := Pack(data, LidarModuleID, WritingModeA)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)

		assert.Equal(t, data, out.(*LidarData))
	})
}
