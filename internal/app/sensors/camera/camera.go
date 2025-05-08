// Package camera предоставляет функционал для чтения данных камеры
package camera

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"errors"
	"time"
)

type Camera struct {
	i    int
	data [][3]int16
}

// ErrNoMoreData ошибка отсутсвия данных
var ErrNoMoreData = errors.New("no more data")

func NewCamera(data [][3]int16) (*Camera, error) {
	return &Camera{data: data}, nil
}

func (c *Camera) Measure(_ context.Context) (any, error) {
	return c.measure()
}

func (c *Camera) Close() error {
	return nil
}

func (c *Camera) measure() (*proto.CameraData, error) {
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
