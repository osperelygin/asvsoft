package proto

type Addr uint8

const (
	ControlModuleAddr Addr = iota * 0x10
	DepthMeterModuleAddr
	LidarModuleAddr
	CommunicationModuleAddr
	NavigationModuleAddr
	GNSSModuleAddr
	IMUModuleAddr
	CheckModuleAddr
)

type MessageID uint8

const (
	ReadingModeA MessageID = iota
	ReadingModeB
	ReadingModeC
	WritingModeA
	WritingModeB
	WritingModeC
)

type Bitmask uint8

const (
	ModuleAddrBitmask Bitmask = 0xF0
	MessageIDBitmask  Bitmask = 0x0F
)
