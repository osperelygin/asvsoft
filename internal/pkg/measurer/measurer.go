// Package measurer ...
package measurer

import (
	"asvsoft/internal/pkg/transmitter"
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
)

type Measurement interface {
	Data() any
	Error() error
}

type Measurer interface {
	Measure(ctx context.Context) Measurement
	Close() error
}

func Run(ctx context.Context, m Measurer, t transmitter.Transmitter) error {
	ctx, cancel := context.WithCancel(ctx)

	quit := make(chan os.Signal, 2)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	measurementChan := make(chan Measurement)

	go func() {
		for {
			select {
			case <-ctx.Done():
				err := m.Close()
				if err != nil {
					log.Errorf("failed to close measurer: %v", err)
				}

				close(measurementChan)

				return
			default:
				measurementChan <- m.Measure(ctx)
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
		case measurement := <-measurementChan:
			measure, err := measurement.Data(), measurement.Error()
			if err != nil {
				log.Errorf("cannot read measure: %v", err)

				continue
			}

			log.Infof("read measure: %+v", measure)

			err = t.Transmit(ctx, measure)
			if err != nil {
				log.Errorf("cannot transmit measure: %v", err)
			}
		}
	}

	return nil
}
