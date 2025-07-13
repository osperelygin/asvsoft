package utils

import (
	"asvsoft/internal/pkg/logger"
	"fmt"
	"time"
)

// RunWithRetries - синхронно запускает функцию с ретраями.
func RunWithRetries(f func() error, log logger.Logger, retriesLimit int, retriesDelay time.Duration) error {
	if f == nil {
		return fmt.Errorf("nothing to do: f is nil")
	}

	var err error

	for i := range retriesLimit {
		err = f()
		if err != nil {
			log.Debugf("[retry %d]: failed to call f: %v", i, err)
			time.Sleep(retriesDelay)

			continue
		}

		break
	}

	if err != nil {
		return fmt.Errorf("failed to run f after %d retries: %w", retriesLimit, err)
	}

	return nil
}
