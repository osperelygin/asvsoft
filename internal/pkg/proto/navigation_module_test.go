package proto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPackIMUDataSuccess(t *testing.T) {
	data := &IMUData{
		Gx: -100, Gy: 101, Gz: 110,
		Ax: 1000, Ay: -1001, Az: 1010,
		Mx: 10000, My: 10001, Mz: -10010,
	}

	t.Run("успешная упаковка и распаковка данных сообщения A", func(t *testing.T) {
		packedData, err := Pack(data, IMUModuleID, WritingModeA)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, data, out.(*IMUData))
	})

	t.Run("успешная упаковка и распаковка данных сообщения B", func(t *testing.T) {
		packedData, err := Pack(data, IMUModuleID, WritingModeB)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, &IMUData{
			Gx: data.Gx, Gy: data.Gy, Gz: data.Gz,
			Ax: data.Ax, Ay: data.Ay, Az: data.Az,
		}, out.(*IMUData))
	})

	t.Run("успешная упаковка и распаковка данных сообщения C", func(t *testing.T) {
		packedData, err := Pack(data, IMUModuleID, WritingModeC)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, &IMUData{
			Mx: data.Mx, My: data.My, Mz: data.Mz,
		}, out.(*IMUData))
	})
}

func TestPackGNSSSDataSuccess(t *testing.T) {
	data := &GNSSData{
		ITowNAVPOSLLH: uint32(time.Now().Unix()) - 3,
		Lon:           37,
		Lat:           50,
		Height:        15000,
		HMSL:          20000,
		HAcc:          14000,
		VAcc:          5000,
		ITowNAVVELNED: uint32(time.Now().Unix()) - 5,
		VelN:          -10,
		VelE:          5,
		VelD:          0,
		Speed:         7,
		GSppeed:       21,
		Heading:       87,
		SAcc:          3,
		CAcc:          40,
	}

	t.Run("успешная упаковка и распаковка данных сообщения A", func(t *testing.T) {
		packedData, err := Pack(data, GNSSModuleID, WritingModeA)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, data, out.(*GNSSData))
	})

	t.Run("успешная упаковка и распаковка данных сообщения B", func(t *testing.T) {
		packedData, err := Pack(data, GNSSModuleID, WritingModeB)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, &GNSSData{
			ITowNAVPOSLLH: data.ITowNAVPOSLLH,
			Lon:           data.Lon,
			Lat:           data.Lat,
			Height:        data.Height,
			HMSL:          data.HMSL,
			HAcc:          data.HAcc,
			VAcc:          data.VAcc,
		}, out.(*GNSSData))
	})

	t.Run("успешная упаковка и распаковка данных сообщения C", func(t *testing.T) {
		packedData, err := Pack(data, GNSSModuleID, WritingModeC)
		assert.NoError(t, err)

		out, err := Unpack(packedData)
		assert.NoError(t, err)
		assert.Equal(t, &GNSSData{
			ITowNAVVELNED: data.ITowNAVVELNED,
			VelN:          data.VelN,
			VelE:          data.VelE,
			VelD:          data.VelD,
			Speed:         data.Speed,
			GSppeed:       data.GSppeed,
			Heading:       data.Heading,
			SAcc:          data.SAcc,
			CAcc:          data.CAcc,
		}, out.(*GNSSData))
	})
}
