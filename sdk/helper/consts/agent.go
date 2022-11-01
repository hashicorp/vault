package consts

// AgentPathCacheClear is the path that the agent will use as its cache-clear
// endpoint.
const AgentPathCacheClear = "/agent/v1/cache-clear"

// AgentPathMetrics is the path the the agent will use to expose its internal
// metrics.
const AgentPathMetrics = "/agent/v1/metrics"

// AgentPathQuit is the path that the agent will use to trigger stopping it.
const AgentPathQuit = "/agent/v1/quit"
