package encoder

import (
	"asvsoft/internal/pkg/common"
	"fmt"
)

type ubytes interface {
	uint8 | uint16 | common.Uint24 | uint32 | int32
}

type sbytes interface {
	int8 | int16 | int32
}

func bytesOf(untyped any) int {
	switch v := untyped.(type) {
	case int16, uint16:
		return 2
	case common.Uint24:
		return 3
	case int32, uint32:
		return 4
	default:
		panic(fmt.Sprintf("bytesOf is not implemented for this type (%T)", v))
	}
}

// func bytesOfV2(untyped any) int {
// 	return int(unsafe.Sizeof(untyped)) / 8
// }
