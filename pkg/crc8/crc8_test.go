package crc8

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChecksum(t *testing.T) {
	t.Run("check crc8 algorithm", func(t *testing.T) {
		checkString := []byte("123456789")

		tables := []*Table{
			smBusTable,
			cdma2000Table,
			darcTable,
			dvb2S2Table,
			ebuTable,
			iCODETable,
			ituTable,
		}

		for _, table := range tables {
			checksum := table.Checksum(checkString)
			require.Equal(t, table.params.Check, checksum, "expecting calculated checksum")
		}
	})
}
