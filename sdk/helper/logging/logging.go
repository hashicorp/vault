package logging

import (
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/hashicorp/go-hclog"
)

type LogFormat int

const (
	UnspecifiedFormat LogFormat = iota
	StandardFormat
	JSONFormat
)

// Stringer implementation
func (l LogFormat) String() string {
	switch l {
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
func NewVaultLogger(level log.Level) log.Logger {
	return NewVaultLoggerWithWriter(log.DefaultOutput, level)
}

// NewVaultLoggerWithWriter creates a new logger with the specified level and
// writer and a Vault formatter
func NewVaultLoggerWithWriter(w io.Writer, level log.Level) log.Logger {
	// If out is os.Stdout and Vault is being run as a Windows Service, writes will
	// fail silently, which may inadvertently prevent writes to other writers.
	// noErrorWriter is used as a wrapper to suppress any errors when writing to out.
	w = noErrorWriter{w: w}

	opts := &log.LoggerOptions{
		Level:             level,
		IndependentLevels: true,
		Output:            w,
		JSONFormat:        ParseEnvLogFormat() == JSONFormat,
	}
	return log.New(opts)
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

// ParseEnvLogFormat parses the log format from an environment variable.
func ParseEnvLogFormat() LogFormat {
	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	switch strings.ToLower(logFormat) {
	case "json", "vault_json", "vault-json", "vaultjson":
		return JSONFormat
	case "standard":
		return StandardFormat
	default:
		return UnspecifiedFormat
	}
}
