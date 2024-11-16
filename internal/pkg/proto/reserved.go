package proto

type ModuleID uint8

const (
	CheckModuleID   ModuleID = 0xF0
	ControlModuleID ModuleID = 0x01
)

const (
	RadioTelemetryModuleID ModuleID = 0x21 + iota*0x10
	CommunicationModule
	IMUModuleID
	GNSSModuleID
	NavigationModuleID
	DepthMeterModuleID
	LidarModuleID
)

type MessageID uint8

const (
	ReadingModeA = 0x11 + iota
	ReadingModeB
	ReadingModeC
	WritingModeA
	WritingModeB
	WritingModeC
)

const (
	headerSize       = 2
	sytemByteSize    = 1
	moduleIDSize     = 1
	msgIDSize        = 1
	timestampSize    = 4
	payloadBytesSize = 1
	checkSumSize     = 2
)

const serviceBytesSize = headerSize +
	sytemByteSize +
	moduleIDSize +
	msgIDSize +
	timestampSize +
	payloadBytesSize +
	checkSumSize

const payloadFirstByte = serviceBytesSize - checkSumSize

const (
	defaultBuffSize    = 512
	defaultReadRetries = 1024
)

var header = []byte{0xFA, 0xFA}

const (
	dummySystemByte byte = 0xFF
)
