// Package transmitter ...
package transmitter

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

var _ Transmitter = (*CommonTransmitter)(nil)

type Transmitter interface {
	Transmit(ctx context.Context, data any) error
}

func New(addr proto.Addr, mode proto.MessageID) *CommonTransmitter {
	return &CommonTransmitter{
		addr: addr,
		mode: mode,
	}
}

type CommonTransmitter struct {
	w    io.Writer
	addr proto.Addr
	mode proto.MessageID
}

func (ct *CommonTransmitter) SetWritter(w io.Writer) *CommonTransmitter {
	ct.w = w
	return ct
}

func (ct *CommonTransmitter) Transmit(_ context.Context, data any) error {
	if ct.w == nil {
		return nil
	}

	b, err := proto.Pack(data, ct.addr, ct.mode)
	if err != nil {
		return fmt.Errorf("cannot pack measure: %w", err)
	}

	log.Debugf("packed measure: %v", b)

	_, err = ct.w.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	return nil
}
