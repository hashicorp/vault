package logging

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	log "github.com/hashicorp/go-hclog"
)

const (
	UnspecifiedFormat LogFormat = iota
	StandardFormat
	JSONFormat
)

type LogFormat int

// LogConfig should be used to supply configuration when creating a new Vault logger
type LogConfig struct {
	name        string
	logLevel    log.Level
	logFormat   LogFormat
	logFilePath string
}

func NewLogConfig(name string, logLevel log.Level, logFormat LogFormat, logFilePath string) LogConfig {
	return LogConfig{
		name:        name,
		logLevel:    logLevel,
		logFormat:   logFormat,
		logFilePath: strings.TrimSpace(logFilePath),
	}
}

func (c LogConfig) IsFormatJson() bool {
	return c.logFormat == JSONFormat
}

// Stringer implementation
func (lf LogFormat) String() string {
	switch lf {
	case UnspecifiedFormat:
		return "unspecified"
	case StandardFormat:
		return "standard"
	case JSONFormat:
		return "json"
	}

	// unreachable
	return "unknown"
}

// noErrorWriter is a wrapper to suppress errors when writing to w.
type noErrorWriter struct {
	w io.Writer
}

func (w noErrorWriter) Write(p []byte) (n int, err error) {
	_, _ = w.w.Write(p)
	// We purposely return n == len(p) as if write was successful
	return len(p), nil
}

// Setup creates a new logger with the specified configuration and writer
func Setup(config LogConfig, w io.Writer) (log.InterceptLogger, error) {
	// Validate the log level
	if config.logLevel.String() == "unknown" {
		return nil, fmt.Errorf("invalid log level: %v", config.logLevel)
	}

	// If out is os.Stdout and Vault is being run as a Windows Service, writes will
	// fail silently, which may inadvertently prevent writes to other writers.
	// noErrorWriter is used as a wrapper to suppress any errors when writing to out.
	writers := []io.Writer{noErrorWriter{w: w}}

	if config.logFilePath != "" {
		dir, fileName := filepath.Split(config.logFilePath)
		if fileName == "" {
			fileName = "vault-agent.log"
		}
		logFile := NewLogFile(dir, fileName)
		if err := logFile.openNew(); err != nil {
			return nil, fmt.Errorf("failed to set up file logging: %w", err)
		}
		writers = append(writers, logFile)
	}

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Name:       config.name,
		Level:      config.logLevel,
		Output:     io.MultiWriter(writers...),
		JSONFormat: config.IsFormatJson(),
	})
	return logger, nil
}

// ParseLogFormat parses the log format from the provided string.
func ParseLogFormat(format string) (LogFormat, error) {
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "":
		return UnspecifiedFormat, nil
	case "standard":
		return StandardFormat, nil
	case "json":
		return JSONFormat, nil
	default:
		return UnspecifiedFormat, fmt.Errorf("unknown log format: %s", format)
	}
}

// ParseLogLevel returns the hclog.Level that corresponds with the provided level string.
// This differs hclog.LevelFromString in that it supports additional level strings.
func ParseLogLevel(logLevel string) (log.Level, error) {
	var result log.Level
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))

	switch logLevel {
	case "trace":
		result = log.Trace
	case "debug":
		result = log.Debug
	case "notice", "info", "":
		result = log.Info
	case "warn", "warning":
		result = log.Warn
	case "err", "error":
		result = log.Error
	default:
		return -1, errors.New(fmt.Sprintf("unknown log level: %s", logLevel))
	}

	return result, nil
}

// TranslateLoggerLevel returns the string that corresponds with logging level of the hclog.Logger.
func TranslateLoggerLevel(logger log.Logger) (string, error) {
	var result string

	if logger.IsTrace() {
		result = "trace"
	} else if logger.IsDebug() {
		result = "debug"
	} else if logger.IsInfo() {
		result = "info"
	} else if logger.IsWarn() {
		result = "warn"
	} else if logger.IsError() {
		result = "error"
	} else {
		return "", fmt.Errorf("unknown log level")
	}

	return result, nil
}
