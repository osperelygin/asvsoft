package communication

import (
	"context"
)

var _ Measurement = (*CommonMeasurement)(nil)

type Measurement interface {
	Data() any
	Error() error
}

type Measurer interface {
	Measure(ctx context.Context) Measurement
	Close() error
}

func NewCommonMeasurement(data any, err error) *CommonMeasurement {
	return &CommonMeasurement{
		data: data,
		err:  err,
	}
}

type CommonMeasurement struct {
	data any
	err  error
}

func (m *CommonMeasurement) Data() any {
	return m.data
}

func (m *CommonMeasurement) Error() error {
	return m.err
}
