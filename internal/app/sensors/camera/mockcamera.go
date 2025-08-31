// Package camera предоставляет функционал для чтения данных камеры
package camera

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"errors"
	"time"
)

type MockCamera struct {
	i    int
	data [][3]int16
}

// ErrNoMoreData ошибка отсутсвия данных
var ErrNoMoreData = errors.New("no more data")

func NewMockCamera(data [][3]int16) (*MockCamera, error) {
	return &MockCamera{data: data}, nil
}

func (c *MockCamera) Measure(_ context.Context) (proto.Packer, error) {
	return c.measure()
}

func (c *MockCamera) Close() error {
	return nil
}

func (c *MockCamera) measure() (*proto.CameraData, error) {
	time.Sleep(40 * time.Millisecond)

	if c.i >= len(c.data) {
		return nil, ErrNoMoreData
	}

	data := &proto.CameraData{
		Yaw:   c.data[c.i][0],
		Pitch: c.data[c.i][1],
		Roll:  c.data[c.i][2],
	}

	c.i++

	return data, nil
}
