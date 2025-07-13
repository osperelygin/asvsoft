package proto

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPackCameraDataSuccess(t *testing.T) {
	t.Run("успешная упакова и распаковка данных", func(t *testing.T) {
		var sentMsg Message

		rawImage, err := os.ReadFile("testdata/1752339024.jpeg")
		require.NoError(t, err)

		cameraData := &CameraData{RawImagePart: rawImage}

		msgBytes, err := sentMsg.Marshal(cameraData, CameraModuleID, WritingModeB)
		require.NoError(t, err)

		var receivedMsg Message

		err = receivedMsg.Unmarshal(msgBytes)
		require.NoError(t, err)

		require.Equal(t, sentMsg, receivedMsg)
	})
}
