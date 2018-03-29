package pluginutil

import (
	"bytes"
	"fmt"
	stdlog "log"
	"strings"

	log "github.com/hashicorp/go-hclog"
)

// pluginLogFaker is a wrapper on logxi.Logger that
// implements hclog.Logger
type hclogFaker struct {
	logger log.Logger

	name    string
	implied []interface{}
}

func (f *hclogFaker) buildLog(msg string, args ...interface{}) (string, []interface{}) {
	if f.name != "" {
		msg = fmt.Sprintf("%s: %s", f.name, msg)
	}
	args = append(f.implied, args...)

	return msg, args
}

func (f *hclogFaker) Trace(msg string, args ...interface{}) {
	msg, args = f.buildLog(msg, args...)
	f.logger.Trace(msg, args...)
}

func (f *hclogFaker) Debug(msg string, args ...interface{}) {
	msg, args = f.buildLog(msg, args...)
	f.logger.Debug(msg, args...)
}

func (f *hclogFaker) Info(msg string, args ...interface{}) {
	msg, args = f.buildLog(msg, args...)
	f.logger.Info(msg, args...)
}

func (f *hclogFaker) Warn(msg string, args ...interface{}) {
	msg, args = f.buildLog(msg, args...)
	f.logger.Warn(msg, args...)
}

func (f *hclogFaker) Error(msg string, args ...interface{}) {
	msg, args = f.buildLog(msg, args...)
	f.logger.Error(msg, args...)
}

func (f *hclogFaker) IsTrace() bool {
	return f.logger.IsTrace()
}

func (f *hclogFaker) IsDebug() bool {
	return f.logger.IsDebug()
}

func (f *hclogFaker) IsInfo() bool {
	return f.logger.IsInfo()
}

func (f *hclogFaker) IsWarn() bool {
	return f.logger.IsWarn()
}

func (f *hclogFaker) IsError() bool {
	return !f.logger.IsTrace() && !f.logger.IsDebug() && !f.logger.IsInfo() && !f.IsWarn()
}

func (f *hclogFaker) With(args ...interface{}) log.Logger {
	var nf = *f
	nf.implied = append(nf.implied, args...)
	return f
}

func (f *hclogFaker) Named(name string) log.Logger {
	var nf = *f
	if nf.name != "" {
		nf.name = nf.name + "." + name
	}
	return &nf
}

func (f *hclogFaker) ResetNamed(name string) log.Logger {
	var nf = *f
	nf.name = name
	return &nf
}

func (f *hclogFaker) StandardLogger(opts *log.StandardLoggerOptions) *stdlog.Logger {
	if opts == nil {
		opts = &log.StandardLoggerOptions{}
	}

	return stdlog.New(&stdlogAdapter{f, opts.InferLevels}, "", 0)
}

// Provides a io.Writer to shim the data out of *log.Logger
// and back into our Logger. This is basically the only way to
// build upon *log.Logger.
type stdlogAdapter struct {
	hl          log.Logger
	inferLevels bool
}

// Take the data, infer the levels if configured, and send it through
// a regular Logger
func (s *stdlogAdapter) Write(data []byte) (int, error) {
	str := string(bytes.TrimRight(data, " \t\n"))

	if s.inferLevels {
		level, str := s.pickLevel(str)
		switch level {
		case log.Trace:
			s.hl.Trace(str)
		case log.Debug:
			s.hl.Debug(str)
		case log.Info:
			s.hl.Info(str)
		case log.Warn:
			s.hl.Warn(str)
		case log.Error:
			s.hl.Error(str)
		default:
			s.hl.Info(str)
		}
	} else {
		s.hl.Info(str)
	}

	return len(data), nil
}

// Detect, based on conventions, what log level this is
func (s *stdlogAdapter) pickLevel(str string) (log.Level, string) {
	switch {
	case strings.HasPrefix(str, "[DEBUG]"):
		return log.Debug, strings.TrimSpace(str[7:])
	case strings.HasPrefix(str, "[TRACE]"):
		return log.Trace, strings.TrimSpace(str[7:])
	case strings.HasPrefix(str, "[INFO]"):
		return log.Info, strings.TrimSpace(str[6:])
	case strings.HasPrefix(str, "[WARN]"):
		return log.Warn, strings.TrimSpace(str[7:])
	case strings.HasPrefix(str, "[ERROR]"):
		return log.Error, strings.TrimSpace(str[7:])
	case strings.HasPrefix(str, "[ERR]"):
		return log.Error, strings.TrimSpace(str[5:])
	default:
		return log.Info, str
	}
}
