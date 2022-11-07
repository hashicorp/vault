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

// Config should be used to supply configuration when creating a new Vault logger
type Config struct {
	Name        string
	LogLevel    log.Level
	LogFormat   LogFormat
	LogFilePath string
}

func (c Config) IsFormatJson() bool {
	return c.LogFormat == JSONFormat
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

// NewVaultLogger creates a new logger with the specified level and a Vault
// formatter
func NewVaultLogger(config Config) (log.InterceptLogger, error) {
	return NewVaultLoggerWithWriter(config, log.DefaultOutput)
}

// NewVaultLoggerWithWriter creates a new logger with the specified
// configuration and writer
func NewVaultLoggerWithWriter(config Config, w io.Writer) (log.InterceptLogger, error) {
	// If out is os.Stdout and Vault is being run as a Windows Service, writes will
	// fail silently, which may inadvertently prevent writes to other writers.
	// noErrorWriter is used as a wrapper to suppress any errors when writing to out.
	writers := []io.Writer{noErrorWriter{w: w}}

	if config.LogFilePath != "" {
		dir, fileName := filepath.Split(config.LogFilePath)
		if fileName == "" {
			fileName = "vault-agent.log"
		}
		logFile := &LogFile{
			fileName: fileName,
			logPath:  dir,
		}
		if err := logFile.openNew(); err != nil {
			return nil, fmt.Errorf("failed to set up file logging: %w", err)
		}
		writers = append(writers, logFile)
	}

	logger := log.NewInterceptLogger(&log.LoggerOptions{
		Name:       config.Name,
		Level:      config.LogLevel,
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
