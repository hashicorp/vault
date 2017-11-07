package logbridge

import (
	"log"

	hclog "github.com/hashicorp/go-hclog"
)

type Logger struct {
	hclogger hclog.Logger
}

func NewLogger(hclogger hclog.Logger) *Logger {
	return &Logger{hclogger: hclogger}
}
func (l *Logger) Trace(msg string, args ...interface{}) {
	l.hclogger.Trace(msg, args...)
}
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.hclogger.Debug(msg, args...)
}
func (l *Logger) Info(msg string, args ...interface{}) {
	l.hclogger.Info(msg, args...)
}
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.hclogger.Warn(msg, args...)
}
func (l *Logger) Error(msg string, args ...interface{}) {
	l.hclogger.Error(msg, args...)
}
func (l *Logger) IsTrace() bool {
	return l.hclogger.IsTrace()
}
func (l *Logger) IsDebug() bool {
	return l.hclogger.IsDebug()
}
func (l *Logger) IsInfo() bool {
	return l.hclogger.IsInfo()
}
func (l *Logger) IsWarn() bool {
	return l.hclogger.IsWarn()
}
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		hclogger: l.hclogger.With(args...),
	}
}
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		hclogger: l.hclogger.Named(name),
	}
}
func (l *Logger) ResetNamed(name string) *Logger {
	return &Logger{
		hclogger: l.hclogger.ResetNamed(name),
	}
}
func (l *Logger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return l.hclogger.StandardLogger(opts)
}
func (l *Logger) LogxiLogger() *LogxiLogger {
	return &LogxiLogger{
		l: l,
	}
}

// This is only for compatibility with whatever the fuck is up with the errors
// coming back from warn/error in Logxi's API. Don't use this directly.
type LogxiLogger struct {
	l *Logger
}

func (l *LogxiLogger) Trace(msg string, args ...interface{}) {
	l.l.Trace(msg, args...)
}
func (l *LogxiLogger) Debug(msg string, args ...interface{}) {
	l.l.Debug(msg, args...)
}
func (l *LogxiLogger) Info(msg string, args ...interface{}) {
	l.l.Info(msg, args...)
}
func (l *LogxiLogger) Warn(msg string, args ...interface{}) error {
	l.l.Warn(msg, args...)
	return nil
}
func (l *LogxiLogger) Error(msg string, args ...interface{}) error {
	l.l.Error(msg, args...)
	return nil
}
func (l *LogxiLogger) Fatal(msg string, args ...interface{}) {
	panic(msg)
}
func (l *LogxiLogger) Log(level int, msg string, args []interface{}) {
	panic(msg)
}
func (l *LogxiLogger) IsTrace() bool {
	return l.l.IsTrace()
}
func (l *LogxiLogger) IsDebug() bool {
	return l.l.IsDebug()
}
func (l *LogxiLogger) IsInfo() bool {
	return l.l.IsInfo()
}
func (l *LogxiLogger) IsWarn() bool {
	return l.l.IsWarn()
}
func (l *LogxiLogger) SetLevel(level int) {
	panic("set level")
}
func (l *LogxiLogger) With(args ...interface{}) *LogxiLogger {
	return l.l.With(args...).LogxiLogger()
}
func (l *LogxiLogger) Named(name string) *LogxiLogger {
	return l.l.Named(name).LogxiLogger()
}
func (l *LogxiLogger) ResetNamed(name string) *LogxiLogger {
	return l.l.ResetNamed(name).LogxiLogger()
}
func (l *LogxiLogger) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return l.l.StandardLogger(opts)
}
