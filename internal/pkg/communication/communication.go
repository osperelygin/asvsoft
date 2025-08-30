package communication

import (
	"context"
	"io"
)

const (
	DefaultChunkSize    = 250
	DefaultRetriesLimit = 10
)

type MeasureCloser interface {
	io.Closer
	Measure(ctx context.Context) (any, error)
}
