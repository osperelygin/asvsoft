// Package logger ...
package logger

import (
	"fmt"
)

var _ Logger = (*DummyLogger)(nil)

type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
	Warnf(format string, args ...any)
	Debugf(format string, args ...any)
	Tracef(format string, args ...any)
}

type DummyLogger struct{}

func (dm DummyLogger) Infof(format string, args ...any)  {} // nolint: revive
func (dm DummyLogger) Errorf(format string, args ...any) {} // nolint: revive
func (dm DummyLogger) Warnf(format string, args ...any)  {} // nolint: revive
func (dm DummyLogger) Debugf(format string, args ...any) {} // nolint: revive
func (dm DummyLogger) Tracef(format string, args ...any) {} // nolint: revive

type Wrapper struct {
	prefix string
	logger Logger
}

func Wrap(logger Logger, prefix string) *Wrapper {
	return &Wrapper{logger: logger, prefix: prefix}
}

func (l Wrapper) Infof(format string, args ...any) {
	l.logger.Infof("%s %s", l.prefix, fmt.Sprintf(format, args...))
}

func (l Wrapper) Errorf(format string, args ...any) {
	l.logger.Errorf("%s %s", l.prefix, fmt.Sprintf(format, args...))
}

func (l Wrapper) Warnf(format string, args ...any) {
	l.logger.Warnf("%s %s", l.prefix, fmt.Sprintf(format, args...))
}

func (l Wrapper) Debugf(format string, args ...any) {
	l.logger.Debugf("%s %s", l.prefix, fmt.Sprintf(format, args...))
}

func (l Wrapper) Tracef(format string, args ...any) {
	l.logger.Tracef("%s %s", l.prefix, fmt.Sprintf(format, args...))
}
