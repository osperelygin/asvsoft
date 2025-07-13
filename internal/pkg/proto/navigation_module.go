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

type NavigationModuleData struct {
	IMUData
	GNSSData
}

func packIMUData(in *IMUData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeA))
		err = encoder.NewEncoder(buf).Encode(
			in.AccFactor,
			in.Ax, in.Ay, in.Az,
			in.GyrFactor,
			in.Gx, in.Gy, in.Gz,
		)
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeB))
		err = encoder.NewEncoder(buf).Encode(
			in.AccFactor,
			in.Ax, in.Ay, in.Az,
			in.GyrFactor,
			in.Gx, in.Gy, in.Gz,
			in.Mx, in.My, in.Mz,
		)
	case WritingModeC:
		buf = bytes.NewBuffer(make([]byte, 0, imuDataPayloadSizeModeC))
		err = encoder.NewEncoder(buf).Encode(
			in.Mx, in.My, in.Mz,
		)
	default:
		panic(fmt.Sprintf("packIMUData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func unpackIMUData(in []byte, msgID MessageID) (out *IMUData, err error) {
	out = new(IMUData)

	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.AccFactor,
			&out.Ax, &out.Ay, &out.Az,
			&out.GyrFactor,
			&out.Gx, &out.Gy, &out.Gz,
		)
	case WritingModeB:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.AccFactor,
			&out.Ax, &out.Ay, &out.Az,
			&out.GyrFactor,
			&out.Gx, &out.Gy, &out.Gz,
			&out.Mx, &out.My, &out.Mz,
		)
	case WritingModeC:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.Mx, &out.My, &out.Mz,
		)
	default:
		panic(fmt.Sprintf("unpackIMUData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}

func packGNSSData(in *GNSSData, msgID MessageID) ([]byte, error) {
	var (
		buf *bytes.Buffer
		err error
	)

	switch msgID {
	case WritingModeA:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeA))
		err = encoder.NewEncoder(buf).Encode(
			in.ITowNAVPOSLLH,
			in.Lon, in.Lat,
			in.Height, in.HMSL,
			in.HAcc, in.VAcc,
			in.ITowNAVVELNED,
			in.VelN, in.VelE, in.VelD,
			in.Speed, in.GSppeed,
			in.Heading,
			in.SAcc, in.CAcc,
		)
	case WritingModeB:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeB))
		err = encoder.NewEncoder(buf).Encode(
			in.ITowNAVPOSLLH,
			in.Lon, in.Lat,
			in.Height, in.HMSL,
			in.HAcc, in.VAcc,
		)
	case WritingModeC:
		buf = bytes.NewBuffer(make([]byte, 0, gnssDataPayloadSizeModeC))
		err = encoder.NewEncoder(buf).Encode(
			in.ITowNAVVELNED,
			in.VelN, in.VelE, in.VelD,
			in.Speed, in.GSppeed,
			in.Heading,
			in.SAcc, in.CAcc,
		)
	default:
		panic(fmt.Sprintf("packGNSSData is not implemented for this message ID: %x", msgID))
	}

	return buf.Bytes(), err
}

func unpackGNSSData(in []byte, msgID MessageID) (out *GNSSData, err error) {
	out = new(GNSSData)

	switch msgID {
	case WritingModeA:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.ITowNAVPOSLLH,
			&out.Lon, &out.Lat,
			&out.Height, &out.HMSL,
			&out.HAcc, &out.VAcc,
			&out.ITowNAVVELNED,
			&out.VelN, &out.VelE, &out.VelD,
			&out.Speed, &out.GSppeed,
			&out.Heading,
			&out.SAcc, &out.CAcc,
		)
	case WritingModeB:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.ITowNAVPOSLLH,
			&out.Lon, &out.Lat,
			&out.Height, &out.HMSL,
			&out.HAcc, &out.VAcc,
		)
	case WritingModeC:
		enc := encoder.NewDecoder(io.NopCloser(bytes.NewReader(in)))
		err = enc.Decode(
			&out.ITowNAVVELNED,
			&out.VelN, &out.VelE, &out.VelD,
			&out.Speed, &out.GSppeed,
			&out.Heading,
			&out.SAcc, &out.CAcc,
		)
	default:
		panic(fmt.Sprintf("unpackGNSSData is not implemented for this message ID: %x", msgID))
	}

	return out, err
}
