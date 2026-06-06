// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pkiexternalca

import (
	"io"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/agent/config"
)

// ServerConfig holds configuration for the PKI External CA server.
type ServerConfig struct {
	Logger      hclog.Logger
	AgentConfig *config.Config
	LogLevel    hclog.Level
	LogWriter   io.Writer
}
