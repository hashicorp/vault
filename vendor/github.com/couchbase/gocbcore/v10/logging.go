package gocbcore

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// LogLevel specifies the severity of a log message.
type LogLevel int

// Various logging levels (or subsystems) which can categorize the message.
// Currently these are ordered in decreasing severity.
const (
	LogError LogLevel = iota
	LogWarn
	LogInfo
	LogDebug
	LogTrace
	LogSched
	LogMaxVerbosity
)

func redactUserData(v interface{}) string {
	return fmt.Sprintf("<ud>%v</ud>", v)
}

func redactMetaData(v interface{}) string {
	return fmt.Sprintf("<md>%v</md>", v)
}

func redactSystemData(v interface{}) string {
	return fmt.Sprintf("<sd>%v</sd>", v)
}

// LogRedactLevel specifies the degree with which to redact the logs.
type LogRedactLevel int

const (
	// RedactNone indicates to perform no redactions
	RedactNone LogRedactLevel = iota

	// RedactPartial indicates to redact all possible user-identifying information from logs.
	RedactPartial

	// RedactFull indicates to fully redact all possible identifying information from logs.
	RedactFull
)

// SetLogRedactionLevel specifies the level with which logs should be redacted.
func SetLogRedactionLevel(level LogRedactLevel) {
	globalLogRedactionLevel = level
}

func isLogRedactionLevelNone() bool {
	return globalLogRedactionLevel == RedactNone
}

func isLogRedactionLevelPartial() bool {
	return globalLogRedactionLevel == RedactPartial
}

func isLogRedactionLevelFull() bool {
	return globalLogRedactionLevel == RedactFull
}

func logLevelToString(level LogLevel) string {
	switch level {
	case LogError:
		return "error"
	case LogWarn:
		return "warn"
	case LogInfo:
		return "info"
	case LogDebug:
		return "debug"
	case LogTrace:
		return "trace"
	case LogSched:
		return "sched"
	}

	return fmt.Sprintf("unknown (%d)", level)
}

// Logger defines a logging interface. You can either use one of the default loggers
// (DefaultStdioLogger(), VerboseStdioLogger()) or implement your own.
type Logger interface {
	// Outputs logging information:
	// level is the verbosity level
	// offset is the position within the calling stack from which the message
	// originated. This is useful for contextual loggers which retrieve file/line
	// information.
	Log(level LogLevel, offset int, format string, v ...interface{}) error
}

type defaultLogger struct {
	Level    LogLevel
	GoLogger *log.Logger
}

func (l *defaultLogger) Log(level LogLevel, offset int, format string, v ...interface{}) error {
	if level > l.Level {
		return nil
	}
	s := fmt.Sprintf(format, v...)
	return l.GoLogger.Output(offset+2, s)
}

var (
	globalDefaultLogger = defaultLogger{
		GoLogger: log.New(os.Stderr, "GOCB ", log.Lmicroseconds|log.Lshortfile), Level: LogDebug,
	}

	globalVerboseLogger = defaultLogger{
		GoLogger: globalDefaultLogger.GoLogger, Level: LogMaxVerbosity,
	}

	globalLogger            Logger
	globalLogRedactionLevel LogRedactLevel
)

// DefaultStdioLogger gets the default standard I/O logger.
//  gocbcore.SetLogger(gocbcore.DefaultStdioLogger())
func DefaultStdioLogger() Logger {
	return &globalDefaultLogger
}

// VerboseStdioLogger is a more verbose level of DefaultStdioLogger(). Messages
// pertaining to the scheduling of ordinary commands (and their responses) will
// also be emitted.
//  gocbcore.SetLogger(gocbcore.VerboseStdioLogger())
func VerboseStdioLogger() Logger {
	return &globalVerboseLogger
}

// SetLogger sets a logger to be used by the library. A logger can be obtained via
// the DefaultStdioLogger() or VerboseStdioLogger() functions. You can also implement
// your own logger using the Logger interface.
func SetLogger(logger Logger) {
	globalLogger = logger
}

type redactableLogValue interface {
	redacted() interface{}
}

func logExf(level LogLevel, offset int, format string, v ...interface{}) {
	if globalLogger != nil {
		if level <= LogInfo && !isLogRedactionLevelNone() {
			// We only redact at info level or below.
			for i, iv := range v {
				if redactable, ok := iv.(redactableLogValue); ok {
					v[i] = redactable.redacted()
				}
			}
		}

		err := globalLogger.Log(level, offset+1, format, v...)
		if err != nil {
			log.Printf("Logger error occurred (%s)\n", err)
		}
	}
}

func logDebugf(format string, v ...interface{}) {
	logExf(LogDebug, 1, format, v...)
}

func logSchedf(format string, v ...interface{}) {
	logExf(LogSched, 1, format, v...)
}

func logWarnf(format string, v ...interface{}) {
	logExf(LogWarn, 1, format, v...)
}

func logErrorf(format string, v ...interface{}) {
	logExf(LogError, 1, format, v...)
}

func logInfof(format string, v ...interface{}) {
	logExf(LogInfo, 1, format, v...)
}

func reindentLog(indent, message string) string {
	reindentedMessage := strings.Replace(message, "\n", "\n"+indent, -1)
	return fmt.Sprintf("%s%s", indent, reindentedMessage)
}
