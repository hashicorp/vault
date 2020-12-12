package gocbcore

import (
	"fmt"
	"sync"
	"time"
)

type clusterAgent struct {
	tlsConfig            *dynTLSConfig
	defaultRetryStrategy RetryStrategy

	httpMux     *httpMux
	tracer      *tracerComponent
	http        *httpComponent
	diagnostics *diagnosticsComponent
	n1ql        *n1qlQueryComponent
	analytics   *analyticsQueryComponent
	search      *searchQueryComponent
	views       *viewQueryComponent

	revLock sync.Mutex
	revID   int64

	configWatchLock sync.Mutex
	configWatchers  []routeConfigWatcher
}

func createClusterAgent(config *clusterAgentConfig) *clusterAgent {
	var tlsConfig *dynTLSConfig
	if config.UseTLS {
		tlsConfig = createTLSConfig(config.Auth, config.TLSRootCAProvider)
	}

	httpCli := createHTTPClient(config.HTTPMaxIdleConns, config.HTTPMaxIdleConnsPerHost,
		config.HTTPIdleConnectionTimeout, tlsConfig)

	tracer := config.Tracer
	if tracer == nil {
		tracer = noopTracer{}
	}
	tracerCmpt := newTracerComponent(tracer, "", config.NoRootTraceSpans)

	c := &clusterAgent{
		tlsConfig: tlsConfig,
		tracer:    tracerCmpt,

		defaultRetryStrategy: config.DefaultRetryStrategy,
	}
	if c.defaultRetryStrategy == nil {
		c.defaultRetryStrategy = newFailFastRetryStrategy()
	}

	circuitBreakerConfig := config.CircuitBreakerConfig
	auth := config.Auth
	userAgent := config.UserAgent

	var httpEpList []string
	for _, hostPort := range config.HTTPAddrs {
		if c.tlsConfig == nil {
			httpEpList = append(httpEpList, fmt.Sprintf("http://%s", hostPort))
		} else {
			httpEpList = append(httpEpList, fmt.Sprintf("https://%s", hostPort))
		}
	}

	c.httpMux = newHTTPMux(circuitBreakerConfig, c)
	c.http = newHTTPComponent(
		httpComponentProps{
			UserAgent:            userAgent,
			DefaultRetryStrategy: c.defaultRetryStrategy,
		},
		httpCli,
		c.httpMux,
		auth,
		c.tracer,
	)
	c.n1ql = newN1QLQueryComponent(c.http, c, c.tracer)
	c.analytics = newAnalyticsQueryComponent(c.http, c.tracer)
	c.search = newSearchQueryComponent(c.http, c.tracer)
	c.views = newViewQueryComponent(c.http, c.tracer)
	// diagnostics at this level will never need to hook KV. There are no persistent connections
	// so Diagnostics calls should be blocked. Ping and WaitUntilReady will only try HTTP services.
	c.diagnostics = newDiagnosticsComponent(nil, c.httpMux, c.http, "", c.defaultRetryStrategy, nil)

	// Kick everything off.
	cfg := &routeConfig{
		mgmtEpList: httpEpList,
		revID:      -1,
	}

	c.httpMux.OnNewRouteConfig(cfg)

	return c
}

func (agent *clusterAgent) RegisterWith(cfgMgr configManager) {
	cfgMgr.AddConfigWatcher(agent)
}

func (agent *clusterAgent) UnregisterWith(cfgMgr configManager) {
	cfgMgr.RemoveConfigWatcher(agent)
}

func (agent *clusterAgent) AddConfigWatcher(watcher routeConfigWatcher) {
	agent.configWatchLock.Lock()
	agent.configWatchers = append(agent.configWatchers, watcher)
	agent.configWatchLock.Unlock()
}

func (agent *clusterAgent) RemoveConfigWatcher(watcher routeConfigWatcher) {
	var idx int
	agent.configWatchLock.Lock()
	for i, w := range agent.configWatchers {
		if w == watcher {
			idx = i
		}
	}

	if idx == len(agent.configWatchers) {
		agent.configWatchers = agent.configWatchers[:idx]
	} else {
		agent.configWatchers = append(agent.configWatchers[:idx], agent.configWatchers[idx+1:]...)
	}
	agent.configWatchLock.Unlock()
}

func (agent *clusterAgent) OnNewRouteConfig(cfg *routeConfig) {
	agent.revLock.Lock()
	// This could be coming from multiple agents so we need to make sure that it's up to date with what we've seen.
	if cfg.revID <= agent.revID {
		agent.revLock.Unlock()
		return
	}

	logDebugf("Cluster agent applying config rev id: %d\n", cfg.revID)

	agent.revID = cfg.revID
	agent.revLock.Unlock()
	agent.configWatchLock.Lock()
	watchers := agent.configWatchers
	agent.configWatchLock.Unlock()

	for _, watcher := range watchers {
		watcher.OnNewRouteConfig(cfg)
	}
}

// N1QLQuery executes a N1QL query against a random connected agent.
func (agent *clusterAgent) N1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return agent.n1ql.N1QLQuery(opts, cb)
}

// PreparedN1QLQuery executes a prepared N1QL query against a random connected agent.
func (agent *clusterAgent) PreparedN1QLQuery(opts N1QLQueryOptions, cb N1QLQueryCallback) (PendingOp, error) {
	return agent.n1ql.PreparedN1QLQuery(opts, cb)
}

// AnalyticsQuery executes an analytics query against a random connected agent.
func (agent *clusterAgent) AnalyticsQuery(opts AnalyticsQueryOptions, cb AnalyticsQueryCallback) (PendingOp, error) {
	return agent.analytics.AnalyticsQuery(opts, cb)
}

// SearchQuery executes a Search query against a random connected agent.
func (agent *clusterAgent) SearchQuery(opts SearchQueryOptions, cb SearchQueryCallback) (PendingOp, error) {
	return agent.search.SearchQuery(opts, cb)
}

// ViewQuery executes a view query against a random connected agent.
func (agent *clusterAgent) ViewQuery(opts ViewQueryOptions, cb ViewQueryCallback) (PendingOp, error) {
	return agent.views.ViewQuery(opts, cb)
}

// DoHTTPRequest will perform an HTTP request against one of the HTTP
// services which are available within the SDK, using a random connected agent.
func (agent *clusterAgent) DoHTTPRequest(req *HTTPRequest, cb DoHTTPRequestCallback) (PendingOp, error) {
	return agent.http.DoHTTPRequest(req, cb)
}

// Ping pings all of the servers we are connected to and returns
// a report regarding the pings that were performed.
func (agent *clusterAgent) Ping(opts PingOptions, cb PingCallback) (PendingOp, error) {
	for _, srv := range opts.ServiceTypes {
		if srv == MemdService {
			return nil, wrapError(errInvalidArgument, "memd service is not valid for use with clusterAgent")
		} else if srv == CapiService {
			return nil, wrapError(errInvalidArgument, "capi service is not valid for use with clusterAgent")
		}
	}

	if len(opts.ServiceTypes) == 0 {
		opts.ServiceTypes = []ServiceType{CbasService, FtsService, N1qlService, MgmtService}
		opts.ignoreMissingServices = true
	}

	return agent.diagnostics.Ping(opts, cb)
}

// WaitUntilReady returns whether or not the Agent has seen a valid cluster config.
func (agent *clusterAgent) WaitUntilReady(deadline time.Time, opts WaitUntilReadyOptions, cb WaitUntilReadyCallback) (PendingOp, error) {
	for _, srv := range opts.ServiceTypes {
		if srv == MemdService {
			return nil, wrapError(errInvalidArgument, "memd service is not valid for use with clusterAgent")
		} else if srv == CapiService {
			return nil, wrapError(errInvalidArgument, "capi service is not valid for use with clusterAgent")
		}
	}

	if len(opts.ServiceTypes) == 0 {
		opts.ServiceTypes = []ServiceType{CbasService, FtsService, N1qlService, MgmtService}
	}

	return agent.diagnostics.WaitUntilReady(deadline, opts, cb)
}

// Close shuts down the agent, closing the underlying http client. This does not cause the agent
// to unregister itself with any configuration providers so be sure to do that first.
func (agent *clusterAgent) Close() error {
	// Close the transports so that they don't hold open goroutines.
	agent.http.Close()

	return nil
}
