package log

import (
	"io"
	"log"

	hcllog "github.com/hashicorp/go-hclog"
)

// The idea behind this struct is that we're going to be adding to each log method with
// openTelemetry calls, so that we don't have to add a diagnose-specific struct into random
// places in the codebase (eg serviceRegistration).

// For now, I'm going to prototype this as a struct that contains a bunch of log-lines in mem.
type Logger struct {
	hcllog.Logger
	logLines []string
}

func (l Logger) Log(level hcllog.Level, msg string, args ...interface{}) {
	l.logLines = append(l.logLines, level.String()+":"+msg)
	l.Logger.Log(level, msg, args)
}

// Emit a message and key/value pairs at the TRACE level
func (l Logger) Trace(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.Logger.Trace(msg, args)
}

// Emit a message and key/value pairs at the DEBUG level
func (l Logger) Debug(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.Logger.Debug(msg, args)
}

// Emit a message and key/value pairs at the INFO level
func (l Logger) Info(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.Logger.Info(msg, args)
}

// Emit a message and key/value pairs at the WARN level
func (l Logger) Warn(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.Logger.Warn(msg, args)
}

// Emit a message and key/value pairs at the ERROR level
func (l Logger) Error(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.Logger.Error(msg, args)
}

func (l Logger) IsTrace() bool {
	return l.Logger.IsTrace()
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (l Logger) IsDebug() bool {
	return l.Logger.IsDebug()
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (l Logger) IsInfo() bool {
	return l.Logger.IsInfo()
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (l Logger) IsWarn() bool {
	return l.Logger.IsWarn()
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (l Logger) IsError() bool {
	return l.Logger.IsError()
}

// ImpliedArgs returns With key/value pairs
func (l Logger) ImpliedArgs() []interface{} {
	return l.Logger.ImpliedArgs()
}

// Creates a sublogger that will always have the given key/value pairs
func (l Logger) With(args ...interface{}) hcllog.Logger {
	return l.Logger.With(args)
}

// Returns the Name of the logger
func (l Logger) Name() string {
	return l.Logger.Name()
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l Logger) Named(name string) hcllog.Logger {
	return l.Logger.Named(name)
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l Logger) ResetNamed(name string) hcllog.Logger {
	return l.Logger.ResetNamed(name)
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l Logger) SetLevel(level hcllog.Level) {
	l.Logger.SetLevel(level)
}

// Return a value that conforms to the stdlib log.Logger interface
func (l Logger) StandardLogger(opts *hcllog.StandardLoggerOptions) *log.Logger {
	return l.Logger.StandardLogger(opts)
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (l Logger) StandardWriter(opts *hcllog.StandardLoggerOptions) io.Writer {
	return l.Logger.StandardWriter(opts)
}

type InterceptLogger struct {
	hcllog.InterceptLogger
	logLines []string
}

func (l InterceptLogger) Log(level hcllog.Level, msg string, args ...interface{}) {
	l.logLines = append(l.logLines, level.String()+":"+msg)
	l.InterceptLogger.Log(level, msg, args)
}

// Emit a message and key/value pairs at the TRACE level
func (l InterceptLogger) Trace(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.InterceptLogger.Trace(msg, args)
}

// Emit a message and key/value pairs at the DEBUG level
func (l InterceptLogger) Debug(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.InterceptLogger.Debug(msg, args)
}

// Emit a message and key/value pairs at the INFO level
func (l InterceptLogger) Info(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.InterceptLogger.Info(msg, args)
}

// Emit a message and key/value pairs at the WARN level
func (l InterceptLogger) Warn(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.InterceptLogger.Warn(msg, args)
}

// Emit a message and key/value pairs at the ERROR level
func (l InterceptLogger) Error(msg string, args ...interface{}) {
	l.logLines = append(l.logLines, msg)
	l.InterceptLogger.Error(msg, args)
}

func (l InterceptLogger) IsTrace() bool {
	return l.InterceptLogger.IsTrace()
}

// Indicate if DEBUG logs would be emitted. This and the other Is* guards
func (l InterceptLogger) IsDebug() bool {
	return l.InterceptLogger.IsDebug()
}

// Indicate if INFO logs would be emitted. This and the other Is* guards
func (l InterceptLogger) IsInfo() bool {
	return l.InterceptLogger.IsInfo()
}

// Indicate if WARN logs would be emitted. This and the other Is* guards
func (l InterceptLogger) IsWarn() bool {
	return l.InterceptLogger.IsWarn()
}

// Indicate if ERROR logs would be emitted. This and the other Is* guards
func (l InterceptLogger) IsError() bool {
	return l.InterceptLogger.IsError()
}

// ImpliedArgs returns With key/value pairs
func (l InterceptLogger) ImpliedArgs() []interface{} {
	return l.InterceptLogger.ImpliedArgs()
}

// Creates a sublogger that will always have the given key/value pairs
func (l InterceptLogger) With(args ...interface{}) hcllog.Logger {
	return l.InterceptLogger.With(args)
}

// Returns the Name of the logger
func (l InterceptLogger) Name() string {
	return l.InterceptLogger.Name()
}

// Create a logger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l InterceptLogger) Named(name string) hcllog.Logger {
	return l.InterceptLogger.Named(name)
}

// Create a logger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l InterceptLogger) ResetNamed(name string) hcllog.Logger {
	return l.InterceptLogger.ResetNamed(name)
}

// Updates the level. This should affect all related loggers as well,
// unless they were created with IndependentLevels. If an
// implementation cannot update the level on the fly, it should no-op.
func (l InterceptLogger) SetLevel(level hcllog.Level) {
	l.InterceptLogger.SetLevel(level)
}

// Return a value that conforms to the stdlib log.Logger interface
func (l InterceptLogger) StandardLogger(opts *hcllog.StandardLoggerOptions) *log.Logger {
	return l.InterceptLogger.StandardLogger(opts)
}

// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
func (l InterceptLogger) StandardWriter(opts *hcllog.StandardLoggerOptions) io.Writer {
	return l.InterceptLogger.StandardWriter(opts)
}

// RegisterSink adds a SinkAdapter to the InterceptLogger
func (l InterceptLogger) RegisterSink(sink hcllog.SinkAdapter) {
	l.InterceptLogger.RegisterSink(sink)
}

// DeregisterSink removes a SinkAdapter from the InterceptLogger
func (l InterceptLogger) DeregisterSink(sink hcllog.SinkAdapter) {
	l.InterceptLogger.DeregisterSink(sink)
}

// Create a interceptlogger that will prepend the name string on the front of all messages.
// If the logger already has a name, the new value will be appended to the current
// name. That way, a major subsystem can use this to decorate all it's own logs
// without losing context.
func (l InterceptLogger) NamedIntercept(name string) InterceptLogger {
	return InterceptLogger{InterceptLogger: l.InterceptLogger.NamedIntercept(name)}
}

// Create a interceptlogger that will prepend the name string on the front of all messages.
// This sets the name of the logger to the value directly, unlike Named which honor
// the current name as well.
func (l InterceptLogger) ResetNamedIntercept(name string) InterceptLogger {
	return InterceptLogger{InterceptLogger: l.InterceptLogger.ResetNamedIntercept(name)}
}

// Deprecated: use StandardLogger
func (l InterceptLogger) StandardLoggerIntercept(opts *hcllog.StandardLoggerOptions) *log.Logger {
	return l.InterceptLogger.StandardLoggerIntercept(opts)
}

// Deprecated: use StandardWriter
func (l InterceptLogger) StandardWriterIntercept(opts *hcllog.StandardLoggerOptions) io.Writer {
	return l.InterceptLogger.StandardWriterIntercept(opts)
}
