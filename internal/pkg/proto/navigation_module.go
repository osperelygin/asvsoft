package proto

import (
	"asvsoft/internal/pkg/encoder"
	"bytes"
	"fmt"
	"io"
)

const (
	imuDataPayloadSizeModeA  = 18
	imuDataPayloadSizeModeB  = 12
	imuDataPayloadSizeModeC  = 6
	gnssDataPayloadSizeModeA = 128
	gnssDataPayloadSizeModeB = 56
	gnssDataPayloadSizeModeC = 72
)

// IMUData - данные АСС и гироскопов
type IMUData struct {
	AccFactor  int16
	GyrFactor  int16
	Gx, Gy, Gz int16
	Ax, Ay, Az int16
	Mx, My, Mz int16
}

func (d IMUData) String() string {
	type _IMUData IMUData
	return fmt.Sprintf("%+v", _IMUData(d))
}

// GNSSData - данные ГНСС
type GNSSData struct {
	// UBX-NAVPOSLLH
	ITowNAVPOSLLH uint32
	Lon, Lat      int32
	Height, HMSL  int32
	HAcc, VAcc    uint32
	// UBX-NAVVELNED
	ITowNAVVELNED    uint32
	VelN, VelE, VelD int32
	Speed, GSppeed   uint32
	Heading          int32
	SAcc, CAcc       uint32
}

func (d GNSSData) String() string {
	type _GNSSData GNSSData
	return fmt.Sprintf("%+v", _GNSSData(d))
}

func (d *IMUData) Pack(msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeA))
		err = encoder.NewEncoder(buf).Encode(
			d.AccFactor,
			d.Ax, d.Ay, d.Az,
			d.GyrFactor,
			d.Gx, d.Gy, d.Gz,
		)
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeB))
		err = encoder.NewEncoder(buf).Encode(
			d.AccFactor,
			d.Ax, d.Ay, d.Az,
			d.GyrFactor,
			d.Gx, d.Gy, d.Gz,
			d.Mx, d.My, d.Mz,
		)
	case WritingModeC:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeC))
		err = encoder.NewEncoder(buf).Encode(
			d.Mx, d.My, d.Mz,
		)
	default:
		panic(fmt.Sprintf("packIMUData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func (d *IMUData) Unpack(in []byte, msgID MessageID) error {
	var err error

	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.AccFactor,
			&d.Ax, &d.Ay, &d.Az,
			&d.GyrFactor,
			&d.Gx, &d.Gy, &d.Gz,
		)
	case WritingModeB:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.AccFactor,
			&d.Ax, &d.Ay, &d.Az,
			&d.GyrFactor,
			&d.Gx, &d.Gy, &d.Gz,
			&d.Mx, &d.My, &d.Mz,
		)
	case WritingModeC:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.Mx, &d.My, &d.Mz,
		)
	default:
		panic(fmt.Sprintf("unpackIMUData is not implemented for this message ID: %x", msgID))
	}

	return err
}

func (d *GNSSData) Pack(msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeA))
		err = encoder.NewEncoder(buf).Encode(
			d.ITowNAVPOSLLH,
			d.Lon, d.Lat,
			d.Height, d.HMSL,
			d.HAcc, d.VAcc,
			d.ITowNAVVELNED,
			d.VelN, d.VelE, d.VelD,
			d.Speed, d.GSppeed,
			d.Heading,
			d.SAcc, d.CAcc,
		)
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeB))
		err = encoder.NewEncoder(buf).Encode(
			d.ITowNAVPOSLLH,
			d.Lon, d.Lat,
			d.Height, d.HMSL,
			d.HAcc, d.VAcc,
		)
	case WritingModeC:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeC))
		err = encoder.NewEncoder(buf).Encode(
			d.ITowNAVVELNED,
			d.VelN, d.VelE, d.VelD,
			d.Speed, d.GSppeed,
			d.Heading,
			d.SAcc, d.CAcc,
		)
	default:
		panic(fmt.Sprintf("packGNSSData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func (d *GNSSData) Unpack(in []byte, msgID MessageID) (err error) {
	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.ITowNAVPOSLLH,
			&d.Lon, &d.Lat,
			&d.Height, &d.HMSL,
			&d.HAcc, &d.VAcc,
			&d.ITowNAVVELNED,
			&d.VelN, &d.VelE, &d.VelD,
			&d.Speed, &d.GSppeed,
			&d.Heading,
			&d.SAcc, &d.CAcc,
		)
	case WritingModeB:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.ITowNAVPOSLLH,
			&d.Lon, &d.Lat,
			&d.Height, &d.HMSL,
			&d.HAcc, &d.VAcc,
		)
	case WritingModeC:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&d.ITowNAVVELNED,
			&d.VelN, &d.VelE, &d.VelD,
			&d.Speed, &d.GSppeed,
			&d.Heading,
			&d.SAcc, &d.CAcc,
		)
	default:
		panic(fmt.Sprintf("unpackGNSSData is not implemented for this message ID: %x", msgID))
	}

	return err
}
