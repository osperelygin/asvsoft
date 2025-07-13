// Package communication ...
package communication

import (
	"asvsoft/internal/pkg/proto"
	"asvsoft/internal/pkg/utils"
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

type MeasureCloser interface {
	io.Closer
	Measure(ctx context.Context) (any, error)
}

func NewSender(m MeasureCloser, addr proto.ModuleID, mode proto.MessageID) *Sender {
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
	m     MeasureCloser
	rwc   io.ReadWriteCloser
	addr  proto.ModuleID
	mode  proto.MessageID
	sleep time.Duration
	sync  bool
}

func (s *Sender) WithReadWriteCloser(rw io.ReadWriteCloser) *Sender {
	s.rwc = rw
	return s
}

func (s *Sender) WithSync(sync bool) *Sender {
	s.sync = sync
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
			err := s.Send(measure)
			if err != nil {
				log.Errorf("cannot transmit measure: %v", err)
			}
		}
	}

	return nil
}

var (
	chunkedRequestModules = map[proto.ModuleID]bool{
		proto.CameraModuleID: true,
	}
)

// Send упаковывает измерения согласно унифицированному протоколу и отправляет пакет в s.rw.
func (s *Sender) Send(data any) error {
	if chunkedRequestModules[s.addr] {
		return s.chunkedSend(data)
	}

	return s.send(data)
}

func (s *Sender) send(data any) error {
	var msg proto.Message

	b, err := msg.Marshal(data, s.addr, s.mode)
	if err != nil {
		return fmt.Errorf("cannot marshal msg: %w", err)
	}

	if s.rwc == nil {
		log.Debugf("s.wc == nil: mock sending msg: %+v", msg)
		return nil
	}

	err = utils.RunWithRetries(func() error {
		_, err := s.rwc.Write(b)
		if err != nil {
			return fmt.Errorf("cannot write measures: %w", err)
		}

		err = s.waitOK()
		if err != nil {
			return fmt.Errorf("failed to wait ok message: %w", err)
		}

		return nil
	}, logrus.StandardLogger(), 2, 0)

	if err != nil {
		return err
	}

	log.Debugf("raw sent msg: %+v", b)
	log.Infof("sent msg: %s", msg)

	time.Sleep(s.sleep)

	return nil
}

const chunkSize = 256

func (s *Sender) chunkedSend(data any) error {
	cameraData, ok := data.(*proto.CameraData)
	if !ok {
		return fmt.Errorf("unexpected data")
	}

	rawImage := cameraData.RawImagePart

	chunkes := len(rawImage) / chunkSize
	if len(rawImage)%chunkSize != 0 {
		chunkes++
	}

	for i := 1; i <= chunkes; i++ {
		start := (i - 1) * chunkSize
		end := min(i*chunkSize, len(rawImage))

		err := s.send(&proto.CameraData{
			RawImagePart:  rawImage[start:end],
			CurrentChunck: uint8(i),
			TotalChunckes: uint8(chunkes),
		})
		if err != nil {
			return fmt.Errorf("failed to send #%d chunk, drop package: %w", i, err)
		}
	}

	return nil
}

func (s *Sender) waitOK() error {
	if !s.sync {
		return nil
	}

	rawResp, err := proto.Read(s.rwc)
	if err != nil {
		return fmt.Errorf("failed to read ok message: %w", err)
	}

	var msg proto.Message

	err = msg.Unmarshal(rawResp)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	if msg.MsgID != proto.ResponseOK {
		return fmt.Errorf("response is not ok: %d", msg.MsgID)
	}

	log.Debugf("successfully got ok msg: %s", msg)

	return nil
}

func (s *Sender) Close() error {
	if s.rwc == nil {
		return nil
	}

	return s.rwc.Close()
}
