package logging

import (
	"io"
	"os"
	"strings"

	log "github.com/hashicorp/go-hclog"
)

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
		JSONFormat: useJson(),
	}
	return log.New(opts)
}

func useJson() bool {
	logFormat := os.Getenv("VAULT_LOG_FORMAT")
	if logFormat == "" {
		logFormat = os.Getenv("LOGXI_FORMAT")
	}
	switch strings.ToLower(logFormat) {
	case "json", "vault_json", "vault-json", "vaultjson":
		return true
	default:
		return false
	}
}
