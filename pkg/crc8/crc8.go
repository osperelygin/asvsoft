// Package crc8 implements the 8-bit cyclic redundancy check.
package crc8

import "math/bits"

// Params represents parameters of a CRC-8 algorithm including polynomial and initial value.
// Read more about parameters: http://www.zlib.net/crc_v3.txt
type Params struct {
	Poly   uint8
	Init   uint8
	RefIn  bool
	RefOut bool
	XorOut uint8
	Check  uint8
	Name   string
}

// Predifined CRC-8 algorithm params. Source: https://reveng.sourceforge.io/crc-catalogue/1-15.htm#crc.cat-bits.8
var (
	SMBus    = Params{0x07, 0x00, false, false, 0x00, 0xF4, "CRC-8"}
	CDMA2000 = Params{0x9B, 0xFF, false, false, 0x00, 0xDA, "CRC-8/CDMA2000"}
	DARC     = Params{0x39, 0x00, true, true, 0x00, 0x15, "CRC-8/DARC"}
	DVB_S2   = Params{0xD5, 0x00, false, false, 0x00, 0xBC, "CRC-8/DVB-S2"}
	EBU      = Params{0x1D, 0xFF, true, true, 0x00, 0x97, "CRC-8/EBU"}
	I_CODE   = Params{0x1D, 0xFD, false, false, 0x00, 0x7E, "CRC-8/I-CODE"}
	ITU      = Params{0x07, 0x00, false, false, 0x55, 0xA1, "CRC-8/ITU"}
)

var (
	smBusTable    = MakeTable(SMBus)
	cdma2000Table = MakeTable(CDMA2000)
	darcTable     = MakeTable(DARC)
	dvb2S2Table   = MakeTable(DVB_S2)
	ebuTable      = MakeTable(EBU)
	iCODETable    = MakeTable(I_CODE)
	ituTable      = MakeTable(ITU)
)

func ChecksumSMBus(data []byte) uint8 {
	return smBusTable.Checksum(data)
}

func ChecksumCDMA2000(data []byte) uint8 {
	return cdma2000Table.Checksum(data)
}
func ChecksumDARC(data []byte) uint8 {
	return darcTable.Checksum(data)
}
func ChecksumDVBS2(data []byte) uint8 {
	return dvb2S2Table.Checksum(data)
}
func ChecksumEBU(data []byte) uint8 {
	return ebuTable.Checksum(data)
}
func ChecksumICODE(data []byte) uint8 {
	return iCODETable.Checksum(data)
}
func ChecksumITU(data []byte) uint8 {
	return ituTable.Checksum(data)
}

// Table is a 256-byte table representing polynomial and algorithm settings for efficient processing.
type Table struct {
	params Params
	data   [256]uint8
}

// MakeTable returns the Table constructed from the specified algorithm.
func MakeTable(p Params) *Table {
	t := new(Table)
	t.params = p

	for n := 0; n < 256; n++ {
		crc := uint8(n)

		for i := 0; i < 8; i++ {
			bit := (crc & 0x80) != 0
			crc <<= 1

			if bit {
				crc ^= p.Poly
			}
		}

		t.data[n] = crc
	}

	return t
}

// Checksum returns CRC checksum of data using specified algorithm represented by the Table.
func (table *Table) Checksum(data []byte) uint8 {
	crc := table.params.Init
	crc = table.update(crc, data)

	return table.complete(crc)
}

// update returns the result of adding the bytes in data to the crc.
func (table *Table) update(crc uint8, data []byte) uint8 {
	if table.params.RefIn {
		for _, d := range data {
			d = bits.Reverse8(d)
			crc = table.data[crc^d]
		}
	} else {
		for _, d := range data {
			crc = table.data[crc^d]
		}
	}

	return crc
}

// complete returns the result of CRC calculation and post-calculation processing of the crc.
func (table *Table) complete(crc uint8) uint8 {
	if table.params.RefOut {
		crc = bits.Reverse8(crc)
	}

	return crc ^ table.params.XorOut
}
