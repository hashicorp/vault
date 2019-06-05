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
		JSONFormat: EnvLogFormat() == JSONFormat,
	}
	return log.New(opts)
}

// SpecifyLogFormat returns a LogFormat by checking each of the provided formats, and
// returning the first one that is not unspecified.  If all of the paramters
// are unspecified, return UnspecifiedFormat.
func SpecifyLogFormat(formats ...LogFormat) LogFormat {

	for _, format := range formats {
		if format != UnspecifiedFormat {
			return format
		}
	}
	return UnspecifiedFormat
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

// EnvLogFormat returns whether there is an environment variable that specifies
// the log format.  Currently the only log format that can be specified is JSON.
func EnvLogFormat() LogFormat {
	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "json", "vault_json", "vault-json", "vaultjson":
		return JSONFormat
	default:
		return UnspecifiedFormat
	}
}
