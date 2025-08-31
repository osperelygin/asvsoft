package communication

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"io"
)

const (
	DefaultChunkSize    = 250
	DefaultRetriesLimit = 10
)

type MeasureCloser interface {
	io.Closer
	Measure(ctx context.Context) (proto.Packer, error)
}
