package communication

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

func NewSender(addr proto.ModuleID, mode proto.MessageID) *Sender {
	return &Sender{
		addr: addr,
		mode: mode,
	}
}

type Sender struct {
	w    io.Writer
	addr proto.ModuleID
	mode proto.MessageID
}

func (s *Sender) WithWritter(w io.Writer) *Sender {
	s.w = w
	return s
}

func (s *Sender) Send(_ context.Context, data any) error {
	if s.w == nil {
		return nil
	}

	b, err := proto.Pack(data, s.addr, s.mode)
	if err != nil {
		return fmt.Errorf("cannot pack measure: %w", err)
	}

	_, err = s.w.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	log.Debugf("raw transmitted data: %+v", b)
	log.Infof("transmitted data: %+v", data)

	return nil
}
