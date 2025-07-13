package camera

import (
	"asvsoft/internal/pkg/logger"
	"asvsoft/internal/pkg/proto"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

const (
	defaultSocketPath    = "/tmp/camera.sock"
	defaultMsgBufferSize = 1 << 16
)

var (
	msgBuffer = make([]byte, defaultMsgBufferSize)
)

type Camera struct {
	log      logger.Logger
	listener net.Listener
	conn     net.Conn
}

func New() (*Camera, error) {
	err := os.Remove(defaultSocketPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("fail to remove old socket: %w", err)
	}

	listener, err := net.Listen("unix", defaultSocketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create listenter %w:", err)
	}

	return &Camera{listener: listener, log: logger.DummyLogger{}}, nil
}

func (c *Camera) WithLogger(log logger.Logger) *Camera {
	c.log = log
	return c
}

func (c *Camera) Close() error {
	errConnClose := c.conn.Close()
	errListenerClose := c.listener.Close()
	errSocketRemove := os.Remove(defaultSocketPath)

	if errConnClose != nil || errListenerClose != nil || errSocketRemove != nil {
		return fmt.Errorf(
			"connection close error: %w, listener close error: %w, socket remove error: %w",
			errConnClose, errListenerClose, errSocketRemove,
		)
	}

	return nil
}

func (c *Camera) Measure(ctx context.Context) (any, error) {
	return c.measure()
}

func (c *Camera) measure() (*proto.CameraData, error) {
	var err error

	if c.conn == nil {
		c.log.Debugf("connection is nil, waiting connection...")

		c.conn, err = c.listener.Accept()
		if err != nil {
			return nil, fmt.Errorf("failed to accept connection: %w", err)
		}

		c.log.Debugf("connection is estabilished, start reading...")
	}

	n, err := c.conn.Read(msgBuffer)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = c.conn.Close()
			c.conn = nil
		}

		return nil, fmt.Errorf("failed to read data from connection: %w", err)
	}

	c.log.Debugf("successfully read %d bytes", n)

	return &proto.CameraData{RawImagePart: msgBuffer[:n]}, nil
}
