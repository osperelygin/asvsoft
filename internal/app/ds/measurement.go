package ds

func NewMeasurement(data any, err error) *Measurement {
	return &Measurement{
		data: data,
		err:  err,
	}
}

type Measurement struct {
	data any
	err  error
}

func (m *Measurement) Data() any {
	return m.data
}

func (m *Measurement) Error() error {
	return m.err
}
