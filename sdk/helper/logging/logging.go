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
	opts := &log.LoggerOptions{
		Level:      level,
		Output:     w,
		JSONFormat: ParseEnvLogFormat() == JSONFormat,
	}
	return log.New(opts)
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
		return UnspecifiedFormat, fmt.Errorf("Unknown log format: %s", format)
	}
}

// ParseEnvLogFormat parses the log format from an environment variable.
func ParseEnvLogFormat() LogFormat {
	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "json", "vault_json", "vault-json", "vaultjson":
		return JSONFormat
	case "standard":
		return StandardFormat
	default:
		return UnspecifiedFormat
	}
}
