// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package observability

import (
	"os"
	"strconv"

	instana "github.com/instana/go-sensor"
)

const (
	// defaultInstanaAgentPort is the default port for the Instana agent
	defaultInstanaAgentPort = 42699
)

// init initializes the Instana sensor for test tracing if the required
// environment variables are set. This allows opt-in tracing for tests
// without requiring changes to individual test files.
//
//	import _ "github.com/hashicorp/vault/helper/testhelpers/observability"
func init() {
	// Only initialize if the agent key is provided
	agentKey := os.Getenv("INSTANA_AGENT_KEY")
	if agentKey == "" {
		return
	}

	// Get agent host, default to localhost
	agentHost := os.Getenv("INSTANA_AGENT_HOST")
	if agentHost == "" {
		agentHost = "localhost"
	}

	// Get agent port, default to defaultInstanaAgentPort
	agentPort := defaultInstanaAgentPort
	if portStr := os.Getenv("INSTANA_AGENT_PORT"); portStr != "" {
		if port, err := strconv.Atoi(portStr); err == nil && port > 0 {
			agentPort = port
		}
	}

	// Initialize Instana sensor with configuration
	// Note: INSTANA_ENDPOINT_URL is used by the agent, not the Go sensor directly
	opts := &instana.Options{
		Service:   "vault-tests",
		AgentHost: agentHost,
		AgentPort: agentPort,
	}

	instana.StartMetrics(opts)
}
