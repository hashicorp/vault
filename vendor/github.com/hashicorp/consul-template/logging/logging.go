package logging

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"time"

	gsyslog "github.com/hashicorp/go-syslog"
	"github.com/hashicorp/logutils"
)

// Log time format
const timeFmt = "2006-01-02T15:04:05.000Z0700"

// Levels are the log levels we respond to=o.
var Levels = []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERR"}

type logWriter struct {
	out io.Writer
}

// To let me replace in tests
var now = func() string {
	return time.Now().Format(timeFmt)
}

// writer to output date / time in a standard format
func (writer logWriter) Write(bytes []byte) (int, error) {
	if len(bytes) == 0 {
		return 0, nil
	}
	return fmt.Fprintf(writer.out, "%s %s", now(), bytes)
}

// Config is the configuration for this log setup.
type Config struct {
	// Level is the log level to use.
	Level string `json:"level"`

	// Syslog and SyslogFacility are the syslog configuration options.
	Syslog         bool   `json:"syslog"`
	SyslogFacility string `json:"syslog_facility"`
	// SyslogName is the progname as it will appear in syslog output (if enabled).
	SyslogName string `json:"name"`

	// Writer is the output where logs should go. If syslog is enabled, data will
	// be written to writer in addition to syslog.
	Writer io.Writer `json:"-"`
}

func Setup(config *Config) error {
	logOutput, err := newWriter(config)
	if err != nil {
		return err
	}

	log.SetFlags(0)
	log.SetOutput(logOutput)

	return nil
}

// Creates a log writer w/ filtering
func newWriter(config *Config) (io.Writer, error) {
	var logOutput io.Writer = logWriter{out: config.Writer}
	logLevel := logutils.LogLevel(strings.ToUpper(config.Level))

	logOutput, err := newLogFilter(logOutput, logLevel)
	if err != nil {
		return nil, err
	}

	if config.Syslog {
		log.Printf("[DEBUG] (logging) enabling syslog on %s", config.SyslogFacility)

		l, err := gsyslog.NewLogger(gsyslog.LOG_NOTICE, config.SyslogFacility, config.SyslogName)
		if err != nil {
			return nil, fmt.Errorf("error setting up syslog logger: %s", err)
		}
		syslog := &SyslogWrapper{l, logOutput.(*logutils.LevelFilter)}
		logOutput = io.MultiWriter(logOutput, syslog)
	}

	return logOutput, nil
}

// NewLogFilter returns a LevelFilter that is configured with the log levels that
// we use.
func newLogFilter(out io.Writer, logLevel logutils.LogLevel) (*logutils.LevelFilter, error) {
	if out == nil {
		out = ioutil.Discard
	}

	logFilter := &logutils.LevelFilter{
		Levels:   Levels,
		MinLevel: logLevel,
		Writer:   out,
	}

	if !validateLevelFilter(logLevel, logFilter) {
		levels := make([]string, 0, len(logFilter.Levels))
		for _, level := range logFilter.Levels {
			levels = append(levels, string(level))
		}
		return nil, fmt.Errorf("invalid log level %q, valid log levels are %s",
			logLevel, strings.Join(levels, ", "))
	}
	return logFilter, nil
}

// validateLevelFilter verifies that the log levels within the filter are valid.
func validateLevelFilter(min logutils.LogLevel, filter *logutils.LevelFilter) bool {
	for _, level := range filter.Levels {
		if level == min {
			return true
		}
	}
	return false
}
