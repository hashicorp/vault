package gosnowflake

import (
	"context"
	"fmt"
	rlog "github.com/sirupsen/logrus"
	"io"
	"path"
	"runtime"
	"time"
)

//SFSessionIDKey is context key of session id
const SFSessionIDKey contextKey = "LOG_SESSION_ID"

//SFSessionUserKey is context key of  user id of a session
const SFSessionUserKey contextKey = "LOG_USER"

//LogKeys these keys in context should be included in logging messages when using logger.WithContext
var LogKeys = [...]contextKey{SFSessionIDKey, SFSessionUserKey}

//SFLogger Snowflake logger interface to expose FieldLogger defined in logrus
type SFLogger interface {
	rlog.Ext1FieldLogger
	SetLogLevel(level string) error
	WithContext(ctx context.Context) *rlog.Entry
	SetOutput(output io.Writer)
}

//SFCallerPrettyfier to provide base file name and function name from calling frame used in SFLogger
func SFCallerPrettyfier(frame *runtime.Frame) (string, string) {
	return path.Base(frame.Function), fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
}

type defaultLogger struct {
	inner *rlog.Logger
}

//SetLogLevel set logging level for calling defaultLogger
func (log *defaultLogger) SetLogLevel(level string) error {
	actualLevel, err := rlog.ParseLevel(level)
	if err != nil {
		return err
	}
	log.inner.SetLevel(actualLevel)
	return nil
}

//WithContext return Entry to include fields in context
func (log *defaultLogger) WithContext(ctx context.Context) *rlog.Entry {
	fields := context2Fields(ctx)
	return log.inner.WithFields(*fields)
}

//CreateDefaultLogger return a new instance of SFLogger with default config
func CreateDefaultLogger() SFLogger {
	var rLogger = rlog.New()
	var formatter = rlog.TextFormatter{CallerPrettyfier: SFCallerPrettyfier}
	rLogger.SetReportCaller(true)
	rLogger.SetFormatter(&formatter)
	var ret = defaultLogger{inner: rLogger}
	return &ret //(&ret).(*SFLogger)
}

// WithField allocates a new entry and adds a field to it.
// Debug, Print, Info, Warn, Error, Fatal or Panic must be then applied to
// this new returned entry.
// If you want multiple fields, use `WithFields`.
func (log *defaultLogger) WithField(key string, value interface{}) *rlog.Entry {
	return log.inner.WithField(key, value)

}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (log *defaultLogger) WithFields(fields rlog.Fields) *rlog.Entry {
	return log.inner.WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (log *defaultLogger) WithError(err error) *rlog.Entry {
	return log.inner.WithError(err)
}

// Overrides the time of the log entry.
func (log *defaultLogger) WithTime(t time.Time) *rlog.Entry {
	return log.inner.WithTime(t)
}

func (log *defaultLogger) Logf(level rlog.Level, format string, args ...interface{}) {
	log.inner.Logf(level, format, args...)
}

func (log *defaultLogger) Tracef(format string, args ...interface{}) {
	log.inner.Tracef(format, args...)
}

func (log *defaultLogger) Debugf(format string, args ...interface{}) {
	log.inner.Debugf(format, args...)
}

func (log *defaultLogger) Infof(format string, args ...interface{}) {
	log.inner.Infof(format, args...)
}

func (log *defaultLogger) Printf(format string, args ...interface{}) {
	log.inner.Printf(format, args...)
}

func (log *defaultLogger) Warnf(format string, args ...interface{}) {
	log.inner.Warnf(format, args...)
}

func (log *defaultLogger) Warningf(format string, args ...interface{}) {
	log.inner.Warningf(format, args...)
}

func (log *defaultLogger) Errorf(format string, args ...interface{}) {
	log.inner.Errorf(format, args...)
}

func (log *defaultLogger) Fatalf(format string, args ...interface{}) {
	log.inner.Fatalf(format, args...)
}

func (log *defaultLogger) Panicf(format string, args ...interface{}) {
	log.inner.Panicf(format, args...)
}

func (log *defaultLogger) Log(level rlog.Level, args ...interface{}) {
	log.inner.Log(level, args...)
}

func (log *defaultLogger) LogFn(level rlog.Level, fn rlog.LogFunction) {
	log.inner.LogFn(level, fn)
}

func (log *defaultLogger) Trace(args ...interface{}) {
	log.inner.Trace(args...)
}

func (log *defaultLogger) Debug(args ...interface{}) {
	log.inner.Debug(args...)
}

func (log *defaultLogger) Info(args ...interface{}) {
	log.inner.Info(args...)
}

func (log *defaultLogger) Print(args ...interface{}) {
	log.inner.Print(args...)
}

func (log *defaultLogger) Warn(args ...interface{}) {
	log.inner.Warn(args...)
}

func (log *defaultLogger) Warning(args ...interface{}) {
	log.inner.Warning(args...)
}

func (log *defaultLogger) Error(args ...interface{}) {
	log.inner.Error(args...)
}

func (log *defaultLogger) Fatal(args ...interface{}) {
	log.inner.Fatal(args...)
}

func (log *defaultLogger) Panic(args ...interface{}) {
	log.inner.Panic(args...)
}

func (log *defaultLogger) TraceFn(fn rlog.LogFunction) {
	log.inner.TraceFn(fn)
}

func (log *defaultLogger) DebugFn(fn rlog.LogFunction) {
	log.inner.DebugFn(fn)
}

func (log *defaultLogger) InfoFn(fn rlog.LogFunction) {
	log.inner.InfoFn(fn)
}

func (log *defaultLogger) PrintFn(fn rlog.LogFunction) {
	log.inner.PrintFn(fn)
}

func (log *defaultLogger) WarnFn(fn rlog.LogFunction) {
	log.inner.PrintFn(fn)
}

func (log *defaultLogger) WarningFn(fn rlog.LogFunction) {
	log.inner.WarningFn(fn)
}

func (log *defaultLogger) ErrorFn(fn rlog.LogFunction) {
	log.inner.ErrorFn(fn)
}

func (log *defaultLogger) FatalFn(fn rlog.LogFunction) {
	log.inner.FatalFn(fn)
}

func (log *defaultLogger) PanicFn(fn rlog.LogFunction) {
	log.inner.PanicFn(fn)
}

func (log *defaultLogger) Logln(level rlog.Level, args ...interface{}) {
	log.inner.Logln(level, args...)
}

func (log *defaultLogger) Traceln(args ...interface{}) {
	log.inner.Traceln(args...)
}

func (log *defaultLogger) Debugln(args ...interface{}) {
	log.inner.Debugln(args...)
}

func (log *defaultLogger) Infoln(args ...interface{}) {
	log.inner.Infoln(args...)
}

func (log *defaultLogger) Println(args ...interface{}) {
	log.inner.Println(args...)
}

func (log *defaultLogger) Warnln(args ...interface{}) {
	log.inner.Warnln(args...)
}

func (log *defaultLogger) Warningln(args ...interface{}) {
	log.inner.Warningln(args...)
}

func (log *defaultLogger) Errorln(args ...interface{}) {
	log.inner.Errorln(args...)
}

func (log *defaultLogger) Fatalln(args ...interface{}) {
	log.inner.Fatalln(args...)
}

func (log *defaultLogger) Panicln(args ...interface{}) {
	log.inner.Panicln(args...)
}

func (log *defaultLogger) Exit(code int) {
	log.inner.Exit(code)
}

// SetLevel sets the logger level.
func (log *defaultLogger) SetLevel(level rlog.Level) {
	log.inner.SetLevel(level)
}

// GetLevel returns the logger level.
func (log *defaultLogger) GetLevel() rlog.Level {
	return log.inner.GetLevel()
}

// AddHook adds a hook to the logger hooks.
func (log *defaultLogger) AddHook(hook rlog.Hook) {
	log.inner.AddHook(hook)

}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (log *defaultLogger) IsLevelEnabled(level rlog.Level) bool {
	return log.inner.IsLevelEnabled(level)
}

// SetFormatter sets the logger formatter.
func (log *defaultLogger) SetFormatter(formatter rlog.Formatter) {
	log.inner.SetFormatter(formatter)
}

// SetOutput sets the logger output.
func (log *defaultLogger) SetOutput(output io.Writer) {
	log.inner.SetOutput(output)
}

func (log *defaultLogger) SetReportCaller(reportCaller bool) {
	log.inner.SetReportCaller(reportCaller)
}

// SetLogger set a new logger of SFLogger interface for gosnowflake
func SetLogger(inLogger *SFLogger) {
	logger = *inLogger //.(*defaultLogger)
}

// GetLogger return logger that is not public
func GetLogger() SFLogger {
	return logger
}

func context2Fields(ctx context.Context) *rlog.Fields {
	var fields = rlog.Fields{}
	if ctx == nil {
		return &fields
	}

	for i := 0; i < len(LogKeys); i++ {
		if ctx.Value(LogKeys[i]) != nil {
			fields[string(LogKeys[i])] = ctx.Value(LogKeys[i])
		}
	}
	return &fields
}
