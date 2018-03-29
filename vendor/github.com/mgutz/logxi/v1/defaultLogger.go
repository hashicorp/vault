package log

import (
	"fmt"
	"io"
)

// DefaultLogger is the default logger for this package.
type DefaultLogger struct {
	writer    io.Writer
	name      string
	level     int
	formatter Formatter
}

// NewLogger creates a new default logger. If writer is not concurrent
// safe, wrap it with NewConcurrentWriter.
func NewLogger(writer io.Writer, name string) Logger {
	formatter, err := createFormatter(name, logxiFormat)
	if err != nil {
		panic("Could not create formatter")
	}
	return NewLogger3(writer, name, formatter)
}

// NewLogger3 creates a new logger with a writer, name and formatter. If writer is not concurrent
// safe, wrap it with NewConcurrentWriter.
func NewLogger3(writer io.Writer, name string, formatter Formatter) Logger {
	var level int
	if name != "__logxi" {
		// if err is returned, then it means the log is disabled
		level = getLogLevel(name)
		if level == LevelOff {
			return NullLog
		}
	}

	log := &DefaultLogger{
		formatter: formatter,
		writer:    writer,
		name:      name,
		level:     level,
	}

	// TODO loggers will be used when watching changes to configuration such
	// as in consul, etcd
	loggers.Lock()
	loggers.loggers[name] = log
	loggers.Unlock()
	return log
}

// New creates a colorable default logger.
func New(name string) Logger {
	return NewLogger(colorableStdout, name)
}

// Trace logs a debug entry.
func (l *DefaultLogger) Trace(msg string, args ...interface{}) {
	l.Log(LevelTrace, msg, args)
}

// Debug logs a debug entry.
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.Log(LevelDebug, msg, args)
}

// Info logs an info entry.
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.Log(LevelInfo, msg, args)
}

// Warn logs a warn entry.
func (l *DefaultLogger) Warn(msg string, args ...interface{}) error {
	if l.IsWarn() {
		defer l.Log(LevelWarn, msg, args)

		for _, arg := range args {
			if err, ok := arg.(error); ok {
				return err
			}
		}

		return nil
	}
	return nil
}

func (l *DefaultLogger) extractLogError(level int, msg string, args []interface{}) error {
	defer l.Log(level, msg, args)

	for _, arg := range args {
		if err, ok := arg.(error); ok {
			return err
		}
	}
	return fmt.Errorf(msg)
}

// Error logs an error entry.
func (l *DefaultLogger) Error(msg string, args ...interface{}) error {
	return l.extractLogError(LevelError, msg, args)
}

// Fatal logs a fatal entry then panics.
func (l *DefaultLogger) Fatal(msg string, args ...interface{}) {
	l.extractLogError(LevelFatal, msg, args)
	defer panic("Exit due to fatal error: ")
}

// Log logs a leveled entry.
func (l *DefaultLogger) Log(level int, msg string, args []interface{}) {
	// log if the log level (warn=4) >= level of message (err=3)
	if l.level < level || silent {
		return
	}
	l.formatter.Format(l.writer, level, msg, args)
}

// IsTrace determines if this logger logs a debug statement.
func (l *DefaultLogger) IsTrace() bool {
	// DEBUG(7) >= TRACE(10)
	return l.level >= LevelTrace
}

// IsDebug determines if this logger logs a debug statement.
func (l *DefaultLogger) IsDebug() bool {
	return l.level >= LevelDebug
}

// IsInfo determines if this logger logs an info statement.
func (l *DefaultLogger) IsInfo() bool {
	return l.level >= LevelInfo
}

// IsWarn determines if this logger logs a warning statement.
func (l *DefaultLogger) IsWarn() bool {
	return l.level >= LevelWarn
}

// SetLevel sets the level of this logger.
func (l *DefaultLogger) SetLevel(level int) {
	l.level = level
}

// SetFormatter set the formatter for this logger.
func (l *DefaultLogger) SetFormatter(formatter Formatter) {
	l.formatter = formatter
}
