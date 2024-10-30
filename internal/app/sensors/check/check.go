// Package check provides dummy check sensor
package check

import (
	"asvsoft/internal/app/ds"
	"asvsoft/internal/pkg/measurer"
	"asvsoft/internal/pkg/proto"
	"context"
	"math/rand"
	"time"
)

// Measurer test measurer
type Measurer struct {
}

// New creates new CheckMeasurer instance
func New() *Measurer {
	return &Measurer{}
}

// Measure sleep 500 ms and returns measurement with random uint32
func (dm *Measurer) Measure(_ context.Context) measurer.Measurement {
	time.Sleep(500 * time.Millisecond)
	return ds.NewMeasurement(&proto.CheckData{Value: rand.Uint32()}, nil) // nolint: gosec
}

// Close ...
func (dm *Measurer) Close() error {
	return nil
}
