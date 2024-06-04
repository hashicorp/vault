// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package consts

import "time"

// AgentPathCacheClear is the path that the agent will use as its cache-clear
// endpoint.
const AgentPathCacheClear = "/agent/v1/cache-clear"

// AgentPathMetrics is the path the agent will use to expose its internal
// metrics.
const AgentPathMetrics = "/agent/v1/metrics"

// AgentPathQuit is the path that the agent will use to trigger stopping it.
const AgentPathQuit = "/agent/v1/quit"

// DefaultMinBackoff is the default minimum backoff time for agent and proxy
const DefaultMinBackoff = 1 * time.Second

// DefaultMaxBackoff is the default max backoff time for agent and proxy
const DefaultMaxBackoff = 5 * time.Minute
