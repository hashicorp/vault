package log

// Trace logs a trace statement. On terminals file and line number are logged.
func Trace(msg string, args ...interface{}) {
	DefaultLog.Trace(msg, args...)
}

// Debug logs a debug statement.
func Debug(msg string, args ...interface{}) {
	DefaultLog.Debug(msg, args...)
}

// Info logs an info statement.
func Info(msg string, args ...interface{}) {
	DefaultLog.Info(msg, args...)
}

// Warn logs a warning statement. On terminals it logs file and line number.
func Warn(msg string, args ...interface{}) {
	DefaultLog.Warn(msg, args...)
}

// Error logs an error statement with callstack.
func Error(msg string, args ...interface{}) {
	DefaultLog.Error(msg, args...)
}

// Fatal logs a fatal statement.
func Fatal(msg string, args ...interface{}) {
	DefaultLog.Fatal(msg, args...)
}

// IsTrace determines if this logger logs a trace statement.
func IsTrace() bool {
	return DefaultLog.IsTrace()
}

// IsDebug determines if this logger logs a debug statement.
func IsDebug() bool {
	return DefaultLog.IsDebug()
}

// IsInfo determines if this logger logs an info statement.
func IsInfo() bool {
	return DefaultLog.IsInfo()
}

// IsWarn determines if this logger logs a warning statement.
func IsWarn() bool {
	return DefaultLog.IsWarn()
}
