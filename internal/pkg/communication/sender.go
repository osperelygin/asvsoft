// Package communication ...
package communication

import (
	"asvsoft/internal/pkg/proto"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

type Measurer interface {
	Measure(ctx context.Context) (any, error)
	Close() error
}

func NewSender(m Measurer, addr proto.ModuleID, mode proto.MessageID) *Sender {
	return &Sender{
		m:    m,
		addr: addr,
		mode: mode,
	}
}

func (s *Sender) WithSleep(sleep time.Duration) *Sender {
	s.sleep = sleep
	return s
}

type Sender struct {
	m     Measurer
	wc    io.WriteCloser
	addr  proto.ModuleID
	mode  proto.MessageID
	sleep time.Duration
}

func (s *Sender) WithWritter(rw io.WriteCloser) *Sender {
	s.wc = rw
	return s
}

// Start асинхронно получает измерения от измерителя s.m и отправляет их в s.wc.
func (s *Sender) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	measureChan := make(chan any)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := s.m.Close()
				if err != nil {
					log.Errorf("failed to close measurer: %v", err)
				}

				close(measureChan)

				return
			default:
				measure, err := s.m.Measure(ctx)
				if err != nil {
					log.Errorf("cannot read measure: %v", err)

					continue
				}

				log.Infof("read measure: %s", measure)

				measureChan <- measure
			}
		}
	}()

LOOP:
	for {
		select {
		case <-quit:
			log.Infoln("signal called, cancel operations")
			cancel()
			break LOOP
		case measure := <-measureChan:
			err := s.Send(ctx, measure)
			if err != nil {
				log.Errorf("cannot transmit measure: %v", err)
			}
		}
	}

	return nil
}

// Send упаковывает измерения согласно унифицированному протоколу и отправляет пакет в s.rw.
func (s *Sender) Send(_ context.Context, data any) error {
	var msg proto.Message

	b, err := msg.Marshal(data, s.addr, s.mode)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	if s.wc == nil {
		log.Debugf("s.wc == nil: mock sending msg: %+v", msg)
		return nil
	}

	_, err = s.wc.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write measures: %w", err)
	}

	log.Debugf("raw sent msg: %+v", b)
	log.Infof("sent msg: %s", msg)

	time.Sleep(s.sleep)

	return nil
}

func (s *Sender) Close() error {
	if s.wc == nil {
		return nil
	}

	return s.wc.Close()
}
