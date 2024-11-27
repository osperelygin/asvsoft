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

	var msg proto.Message

	b, err := msg.Marshal(data, s.addr, s.mode)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	_, err = s.w.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	log.Debugf("raw sent msg: %+v", b)
	log.Infof("sent msg: %+v", msg)

	return nil
}
