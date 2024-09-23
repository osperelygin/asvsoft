package proto

const (
	checkSumSize      = 2
	syncFramePartSize = 3
	servicePartSize   = 11 // header + id + ts + payloadSize + checkSumSize = 3 + 1 + 4 + 1 + 2
)

const (
	defaultBuffSize    = 512
	defaultReadRetries = 1024
)
