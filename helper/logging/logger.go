// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package logging

import (
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-multierror"
)

const (
	UnspecifiedFormat LogFormat = iota
	StandardFormat
	JSONFormat
)

// defaultRotateDuration is the default time taken by the agent to rotate logs
const defaultRotateDuration = 24 * time.Hour

type LogFormat int

// LogConfig should be used to supply configuration when creating a new Vault logger
type LogConfig struct {
	// Name is the name the returned logger will use to prefix log lines.
	Name string

	// LogLevel is the minimum level to be logged.
	LogLevel hclog.Level

	// LogFormat is the log format to use, supported formats are 'standard' and 'json'.
	LogFormat LogFormat

	// LogFilePath is the path to write the logs to the user specified file.
	LogFilePath string

	// LogRotateDuration is the user specified time to rotate logs
	LogRotateDuration time.Duration

	// LogRotateBytes is the user specified byte limit to rotate logs
	LogRotateBytes int

	// LogRotateMaxFiles is the maximum number of past archived log files to keep
	LogRotateMaxFiles int

	// DefaultFileName should be set to the value to be used if the LogFilePath
	// ends in a path separator such as '/var/log/'
	// Examples of the default name are as follows: 'vault', 'agent' or 'proxy.
	// The creator of this struct *must* ensure that it is assigned before doing
	// anything with LogConfig!
	DefaultFileName string
}

// NewLogConfig should be used to initialize the LogConfig struct.
func NewLogConfig(defaultFileName string) (*LogConfig, error) {
	defaultFileName = strings.TrimSpace(defaultFileName)
	if defaultFileName == "" {
		return nil, errors.New("default file name is required")
	}

	return &LogConfig{DefaultFileName: defaultFileName}, nil
}

func (c *LogConfig) isLevelInvalid() bool {
	return c.LogLevel == hclog.NoLevel || c.LogLevel == hclog.Off || c.LogLevel.String() == "unknown"
}

func (c *LogConfig) isFormatJson() bool {
	return c.LogFormat == JSONFormat
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

// parseFullPath takes a full path intended to be the location for log files and
// breaks it down into a directory and a file name. It checks both of these for
// the common globbing character '*' and returns an error if it is present.
func parseFullPath(fullPath string) (directory, fileName string, err error) {
	directory, fileName = filepath.Split(fullPath)

	globChars := "*?["
	if strings.ContainsAny(directory, globChars) {
		err = multierror.Append(err, fmt.Errorf("directory contains glob character"))
	}
	if fileName == "" {
		fileName = "vault.log"
	} else if strings.ContainsAny(fileName, globChars) {
		err = multierror.Append(err, fmt.Errorf("file name contains globbing character"))
	}

	return directory, fileName, err
}

// Setup creates a new logger with the specified configuration and writer
func Setup(config *LogConfig, w io.Writer) (hclog.InterceptLogger, error) {
	// Validate the log level
	if config.isLevelInvalid() {
		return nil, fmt.Errorf("invalid log level: %v", config.LogLevel)
	}

	// If out is os.Stdout and Vault is being run as a Windows Service, writes will
	// fail silently, which may inadvertently prevent writes to other writers.
	// noErrorWriter is used as a wrapper to suppress any errors when writing to out.
	writers := []io.Writer{noErrorWriter{w: w}}

	// Create a file logger if the user has specified the path to the log file
	if config.LogFilePath != "" {
		dir, fileName, err := parseFullPath(config.LogFilePath)
		if err != nil {
			return nil, err
		}
		if fileName == "" {
			fileName = fmt.Sprintf("%s.log", config.DefaultFileName)
		}
		if config.LogRotateDuration == 0 {
			config.LogRotateDuration = defaultRotateDuration
		}

		logFile := &LogFile{
			fileName:         fileName,
			logPath:          dir,
			duration:         config.LogRotateDuration,
			maxBytes:         config.LogRotateBytes,
			maxArchivedFiles: config.LogRotateMaxFiles,
		}
		if err := logFile.pruneFiles(); err != nil {
			return nil, fmt.Errorf("failed to prune log files: %w", err)
		}
		if err := logFile.openNew(); err != nil {
			return nil, fmt.Errorf("failed to setup logging: %w", err)
		}
		writers = append(writers, logFile)
	}

	logger := hclog.NewInterceptLogger(&hclog.LoggerOptions{
		Name:              config.Name,
		Level:             config.LogLevel,
		IndependentLevels: true,
		Output:            io.MultiWriter(writers...),
		JSONFormat:        config.isFormatJson(),
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
func ParseLogLevel(logLevel string) (hclog.Level, error) {
	var result hclog.Level
	logLevel = strings.ToLower(strings.TrimSpace(logLevel))

	switch logLevel {
	case "trace":
		result = hclog.Trace
	case "debug":
		result = hclog.Debug
	case "notice", "info", "":
		result = hclog.Info
	case "warn", "warning":
		result = hclog.Warn
	case "err", "error":
		result = hclog.Error
	default:
		return -1, errors.New(fmt.Sprintf("unknown log level: %s", logLevel))
	}

	return result, nil
}

// TranslateLoggerLevel returns the string that corresponds with logging level of the hclog.Logger.
func TranslateLoggerLevel(logger hclog.Logger) (string, error) {
	logLevel := logger.GetLevel()

	switch logLevel {
	case hclog.Trace, hclog.Debug, hclog.Info, hclog.Warn, hclog.Error:
		return logLevel.String(), nil
	default:
		return "", fmt.Errorf("unknown log level")
	}
}
