package gocbcore

import (
	"errors"
	"sync"
	"time"
)

// AgentGroup represents a collection of agents that can be used for performing operations
// against a cluster. It holds an internal special agent type which does not create its own
// memcached connections but registers itself for cluster config updates on all agents that
// are created through it.
type AgentGroup struct {
	agentsLock  sync.Mutex
	boundAgents map[string]*Agent
	// clusterAgent holds no memcached connections but can be used for cluster level (i.e. http) operations.
	// It sets its own internal state by listening to cluster config updates on underlying agents.
	clusterAgent *clusterAgent

	config *AgentGroupConfig
}

// CreateAgentGroup will return a new AgentGroup with a base config of the config provided.
// Volatile: AgentGroup is subject to change or removal.
func CreateAgentGroup(config *AgentGroupConfig) (*AgentGroup, error) {
	logInfof("SDK Version: gocbcore/%s", goCbCoreVersionStr)
	logInfof("Creating new agent group: %+v", config)

	c := config.toAgentConfig()
	agent, err := CreateAgent(c)
	if err != nil {
		return nil, err
	}

	ag := &AgentGroup{
		config:      config,
		boundAgents: make(map[string]*Agent),
	}

	ag.clusterAgent = createClusterAgent(&clusterAgentConfig{
		HTTPAddrs:                 config.HTTPAddrs,
		UserAgent:                 config.UserAgent,
		UseTLS:                    config.UseTLS,
		Auth:                      config.Auth,
		TLSRootCAProvider:         config.TLSRootCAProvider,
		HTTPMaxIdleConns:          config.HTTPMaxIdleConns,
		HTTPMaxIdleConnsPerHost:   config.HTTPMaxIdleConnsPerHost,
		HTTPIdleConnectionTimeout: config.HTTPIdleConnectionTimeout,
		Tracer:                    config.Tracer,
		NoRootTraceSpans:          config.NoRootTraceSpans,
		DefaultRetryStrategy:      config.DefaultRetryStrategy,
		CircuitBreakerConfig:      config.CircuitBreakerConfig,
	})
	ag.clusterAgent.RegisterWith(agent.cfgManager)

	ag.boundAgents[config.BucketName] = agent

	return ag, nil
}

// OpenBucket will attempt to open a new bucket against the cluster.
// If an agent using the specified bucket name already exists then this will not open a new connection.
func (ag *AgentGroup) OpenBucket(bucketName string) error {
	if bucketName == "" {
		return wrapError(errInvalidArgument, "bucket name cannot be empty")
	}

	existing := ag.GetAgent(bucketName)
	if existing != nil {
		return nil
	}

	config := ag.config.toAgentConfig()
	config.BucketName = bucketName

	agent, err := CreateAgent(config)
	if err != nil {
		return err
	}

	ag.clusterAgent.RegisterWith(agent.cfgManager)

	ag.agentsLock.Lock()
	ag.boundAgents[bucketName] = agent
	ag.agentsLock.Unlock()
	ag.maybeCloseGlobalAgent()

	return nil
}

// GetAgent will return the agent, if any, corresponding to the bucket name specified.
func (ag *AgentGroup) GetAgent(bucketName string) *Agent {
	if bucketName == "" {
		// We don't allow access to the global level agent. We close that agent on OpenBucket so we don't want
		// to return an agent that we then later close. Doing so would only lead to pain.
		return nil
	}

	ag.agentsLock.Lock()
	existingAgent := ag.boundAgents[bucketName]
	ag.agentsLock.Unlock()
	if existingAgent != nil {
		return existingAgent
	}

	return nil
}

// Close will close all underlying agents.
func (ag *AgentGroup) Close() error {
	var firstError error
	ag.agentsLock.Lock()
	for _, agent := range ag.boundAgents {
		ag.clusterAgent.UnregisterWith(agent.cfgManager)
		if err := agent.Close(); err != nil && firstError == nil {
			firstError = err
		}
	}
	ag.agentsLock.Unlock()
	if err := ag.clusterAgent.Close(); err != nil && firstError == nil {
		firstError = err
	}

	return firstError
}

// N1QLQuery executes a N1QL query against a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) N1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return ag.clusterAgent.N1QLQuery(opts, cb)
}

// PreparedN1QLQuery executes a prepared N1QL query against a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) PreparedN1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return ag.clusterAgent.PreparedN1QLQuery(opts, cb)
}

// AnalyticsQuery executes an analytics query against a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) AnalyticsQuery(opts AnalyticsQueryOptions, cb AnalyticsQueryCallback) (PendingOp, error) {
	return ag.clusterAgent.AnalyticsQuery(opts, cb)
}

// SearchQuery executes a Search query against a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) SearchQuery(opts SearchQueryOptions, cb SearchQueryCallback) (PendingOp, error) {
	return ag.clusterAgent.SearchQuery(opts, cb)
}

// ViewQuery executes a view query against a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) ViewQuery(opts ViewQueryOptions, cb ViewQueryCallback) (PendingOp, error) {
	return ag.clusterAgent.ViewQuery(opts, cb)
}

// DoHTTPRequest will perform an HTTP request against one of the HTTP
// services which are available within the SDK, using a random connected agent.
// If no agent is connected then this will block until one is available or the deadline is reached.
func (ag *AgentGroup) DoHTTPRequest(req *HTTPRequest, cb DoHTTPRequestCallback) (PendingOp, error) {
	return ag.clusterAgent.DoHTTPRequest(req, cb)
}

// WaitUntilReady returns whether or not the AgentGroup can ping the requested services.
func (ag *AgentGroup) WaitUntilReady(deadline time.Time, opts WaitUntilReadyOptions,
	cb WaitUntilReadyCallback) (PendingOp, error) {
	return ag.clusterAgent.WaitUntilReady(deadline, opts, cb)
}

// Ping pings all of the servers we are connected to and returns
// a report regarding the pings that were performed.
func (ag *AgentGroup) Ping(opts PingOptions, cb PingCallback) (PendingOp, error) {
	return ag.clusterAgent.Ping(opts, cb)
}

// Diagnostics returns diagnostics information about the client.
// Mainly containing a list of open connections and their current
// states.
func (ag *AgentGroup) Diagnostics(opts DiagnosticsOptions) (*DiagnosticInfo, error) {
	var agents []*Agent
	ag.agentsLock.Lock()
	// There's no point in trying to get diagnostics from clusterAgent as it has no kv connections.
	// In fact it doesn't even expose a Diagnostics function.
	for _, agent := range ag.boundAgents {
		agents = append(agents, agent)
	}
	ag.agentsLock.Unlock()

	if len(agents) == 0 {
		return nil, errors.New("no agents available")
	}

	var firstError error
	var diags []*DiagnosticInfo
	for _, agent := range agents {
		report, err := agent.diagnostics.Diagnostics(opts)
		if err != nil && firstError == nil {
			firstError = err
			continue
		}

		diags = append(diags, report)
	}

	if len(diags) == 0 {
		return nil, firstError
	}

	var overallReport DiagnosticInfo
	var connected int
	var expected int
	for _, report := range diags {
		expected++
		overallReport.MemdConns = append(overallReport.MemdConns, report.MemdConns...)
		if report.State == ClusterStateOnline {
			connected++
		}
		if report.ConfigRev > overallReport.ConfigRev {
			overallReport.ConfigRev = report.ConfigRev
		}
	}

	if connected == expected {
		overallReport.State = ClusterStateOnline
	} else if connected > 0 {
		overallReport.State = ClusterStateDegraded
	} else {
		overallReport.State = ClusterStateOffline
	}

	return &overallReport, nil
}

func (ag *AgentGroup) maybeCloseGlobalAgent() {
	ag.agentsLock.Lock()
	// Close and delete the global level agent that we created on Connect.
	agent := ag.boundAgents[""]
	if agent == nil {
		ag.agentsLock.Unlock()
		return
	}
	logDebugf("Shutting down global level agent")
	delete(ag.boundAgents, "")
	ag.agentsLock.Unlock()

	ag.clusterAgent.UnregisterWith(agent.cfgManager)
	if err := agent.Close(); err != nil {
		logDebugf("Failed to close agent: %s", err)
	}
}
