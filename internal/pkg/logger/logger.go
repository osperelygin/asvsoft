// Package logger ...
package logger

var _ Logger = (*DummyLogger)(nil)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type DummyLogger struct{}

func (dm DummyLogger) Infof(format string, args ...interface{})  {} // nolint: revive
func (dm DummyLogger) Errorf(format string, args ...interface{}) {} // nolint: revive
func (dm DummyLogger) Warnf(format string, args ...interface{})  {} // nolint: revive
func (dm DummyLogger) Debugf(format string, args ...interface{}) {} // nolint: revive
