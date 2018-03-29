package logging

import (
	"log"

	"github.com/hashicorp/go-hclog"
)

func NewNullLogger() hclog.Logger {
	return &nullLogger{}
}

type nullLogger struct{}

func (l *nullLogger) Trace(msg string, args ...interface{}) {}

func (l *nullLogger) Debug(msg string, args ...interface{}) {}

func (l *nullLogger) Info(msg string, args ...interface{}) {}

func (l *nullLogger) Warn(msg string, args ...interface{}) {}

func (l *nullLogger) Error(msg string, args ...interface{}) {}

func (l *nullLogger) IsTrace() bool { return false }

func (l *nullLogger) IsDebug() bool { return false }

func (l *nullLogger) IsInfo() bool { return false }

func (l *nullLogger) IsWarn() bool { return false }

func (l *nullLogger) IsError() bool { return false }

func (l *nullLogger) With(args ...interface{}) hclog.Logger { return l }

func (l *nullLogger) Named(name string) hclog.Logger { return l }

func (l *nullLogger) ResetNamed(name string) hclog.Logger { return l }

func (l *nullLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return &log.Logger{}
}
