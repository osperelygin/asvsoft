package proto

const (
	checkSumSize         = 2
	syncFramePartSize    = 3
	serviceFramePartSize = 11
)

var syncFramePart = []byte{0x57, 0x10, 0xFF}

const (
	DefaultReadRetries = 1 << 10
)
