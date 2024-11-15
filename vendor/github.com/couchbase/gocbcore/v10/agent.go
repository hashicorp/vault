// Package gocbcore implements methods for low-level communication
// with a Couchbase Server cluster.
package gocbcore

import (
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Agent represents the base client handling connections to a Couchbase Server.
// This is used internally by the higher level classes for communicating with the cluster,
// it can also be used to perform more advanced operations with a cluster.
type Agent struct {
	clientID             string
	bucketName           string
	defaultRetryStrategy RetryStrategy

	pollerController configPollerController
	kvMux            *kvMux
	httpMux          *httpMux
	dialer           *memdClientDialerComponent

	cfgManager   *configManagementComponent
	errMap       *errMapComponent
	collections  *collectionsComponent
	tracer       *tracerComponent
	http         *httpComponent
	diagnostics  *diagnosticsComponent
	crud         *crudComponent
	observe      *observeComponent
	stats        *statsComponent
	n1ql         *n1qlQueryComponent
	analytics    *analyticsQueryComponent
	search       *searchQueryComponent
	views        *viewQueryComponent
	zombieLogger *zombieLoggerComponent

	// These connection settings are only ever changed when ForceReconnect or ReconfigureSecurity are called.
	connectionSettingsLock sync.Mutex
	auth                   AuthProvider
	authMechanisms         []AuthMechanism
	tlsConfig              *dynTLSConfig

	srvDetails  *srvDetails
	shutdownSig chan struct{}
}

// HTTPClient returns a pre-configured HTTP Client for communicating with
// Couchbase Server.  You must still specify authentication information
// for any dispatched requests.
func (agent *Agent) HTTPClient() *http.Client {
	return agent.http.cli
}

type srvDetails struct {
	Addrs  routeEndpoints
	Record SRVRecord
}

// CreateAgent creates an agent for performing normal operations.
func CreateAgent(config *AgentConfig) (*Agent, error) {
	return createAgent(config)
}

func createAgent(config *AgentConfig) (*Agent, error) {
	logInfof("SDK Version: gocbcore/%s", goCbCoreVersionStr)
	logInfof("Creating new agent: %+v", config)

	tracer := config.TracerConfig.Tracer
	if tracer == nil {
		tracer = noopTracer{}
	}

	tracerCmpt := newTracerComponent(tracer, config.BucketName, config.TracerConfig.NoRootTraceSpans, config.MeterConfig.Meter)

	c := &Agent{
		clientID:   formatCbUID(randomCbUID()),
		bucketName: config.BucketName,
		tracer:     tracerCmpt,

		defaultRetryStrategy: config.DefaultRetryStrategy,

		errMap: newErrMapManager(config.BucketName),
		auth:   config.SecurityConfig.Auth,

		shutdownSig: make(chan struct{}),
	}

	tlsConfig, err := setupTLSConfig(config.SeedConfig.MemdAddrs, config.SecurityConfig)
	if err != nil {
		return nil, err
	}
	c.tlsConfig = tlsConfig

	httpIdleConnTimeout := 1000 * time.Millisecond
	if config.HTTPConfig.IdleConnectionTimeout > 0 {
		httpIdleConnTimeout = config.HTTPConfig.IdleConnectionTimeout
	}
	httpConnectTimeout := 30 * time.Second
	if config.HTTPConfig.ConnectTimeout > 0 {
		httpConnectTimeout = config.HTTPConfig.ConnectTimeout
	}

	circuitBreakerConfig := config.CircuitBreakerConfig
	userAgent := config.UserAgent
	useMutationTokens := config.IoConfig.UseMutationTokens
	disableDecompression := config.CompressionConfig.DisableDecompression
	useCompression := config.CompressionConfig.Enabled
	useCollections := config.IoConfig.UseCollections
	useJSONHello := !config.IoConfig.DisableJSONHello
	usePITRHello := config.IoConfig.EnablePITRHello
	useXErrorHello := !config.IoConfig.DisableXErrorHello
	useSyncReplicationHello := !config.IoConfig.DisableSyncReplicationHello
	useResourceUnits := config.InternalConfig.EnableResourceUnitsTrackingHello
	compressionMinSize := 32
	compressionMinRatio := 0.83
	useDurations := config.IoConfig.UseDurations
	useOutOfOrder := config.IoConfig.UseOutOfOrderResponses
	UseClusterMapNotifications := config.IoConfig.UseClusterMapNotifications

	kvConnectTimeout := 7000 * time.Millisecond
	if config.KVConfig.ConnectTimeout > 0 {
		kvConnectTimeout = config.KVConfig.ConnectTimeout
	}

	serverWaitTimeout := 5 * time.Second
	if config.KVConfig.ServerWaitBackoff > 0 {
		serverWaitTimeout = config.KVConfig.ServerWaitBackoff
	}

	kvPoolSize := 1
	if config.KVConfig.PoolSize > 0 {
		kvPoolSize = config.KVConfig.PoolSize
	}

	maxQueueSize := 2048
	if config.KVConfig.MaxQueueSize > 0 {
		maxQueueSize = config.KVConfig.MaxQueueSize
	}

	kvBufferSize := uint(0)
	if config.KVConfig.ConnectionBufferSize > 0 {
		kvBufferSize = config.KVConfig.ConnectionBufferSize
	}

	confHTTPRetryDelay := 10 * time.Second
	if config.ConfigPollerConfig.HTTPRetryDelay > 0 {
		confHTTPRetryDelay = config.ConfigPollerConfig.HTTPRetryDelay
	}

	confHTTPRedialPeriod := 10 * time.Second
	if config.ConfigPollerConfig.HTTPRedialPeriod > 0 {
		confHTTPRedialPeriod = config.ConfigPollerConfig.HTTPRedialPeriod
	}

	confHTTPMaxWait := 5 * time.Second
	if config.ConfigPollerConfig.HTTPMaxWait > 0 {
		confHTTPMaxWait = config.ConfigPollerConfig.HTTPMaxWait
	}

	confCccpMaxWait := 3 * time.Second
	if config.ConfigPollerConfig.CccpMaxWait > 0 {
		confCccpMaxWait = config.ConfigPollerConfig.CccpMaxWait
	}

	confCccpPollPeriod := 2500 * time.Millisecond
	if config.ConfigPollerConfig.CccpPollPeriod > 0 {
		confCccpPollPeriod = config.ConfigPollerConfig.CccpPollPeriod
	}

	if config.CompressionConfig.MinSize > 0 {
		compressionMinSize = config.CompressionConfig.MinSize
	}
	if config.CompressionConfig.MinRatio > 0 {
		compressionMinRatio = config.CompressionConfig.MinRatio
		if compressionMinRatio >= 1.0 {
			compressionMinRatio = 1.0
		}
	}
	if c.defaultRetryStrategy == nil {
		c.defaultRetryStrategy = newFailFastRetryStrategy()
	}

	c.authMechanisms = authMechanismsFromConfig(config.SecurityConfig.AuthMechanisms, tlsConfig != nil)

	httpEpList := routeEndpoints{}
	var srcHTTPAddrs []routeEndpoint
	for _, hostPort := range config.SeedConfig.HTTPAddrs {
		if config.SecurityConfig.UseTLS && !config.SecurityConfig.NoTLSSeedNode {
			ep := routeEndpoint{
				Address:    fmt.Sprintf("https://%s", hostPort),
				IsSeedNode: true,
			}
			httpEpList.SSLEndpoints = append(httpEpList.SSLEndpoints, ep)
			srcHTTPAddrs = append(srcHTTPAddrs, ep)
		} else {
			ep := routeEndpoint{
				Address:    fmt.Sprintf("http://%s", hostPort),
				IsSeedNode: true,
			}
			httpEpList.NonSSLEndpoints = append(httpEpList.NonSSLEndpoints, ep)
			srcHTTPAddrs = append(srcHTTPAddrs, ep)
		}
	}

	if config.OrphanReporterConfig.Enabled {
		zombieLoggerInterval := 10 * time.Second
		zombieLoggerSampleSize := 10
		if config.OrphanReporterConfig.ReportInterval > 0 {
			zombieLoggerInterval = config.OrphanReporterConfig.ReportInterval
		}
		if config.OrphanReporterConfig.SampleSize > 0 {
			zombieLoggerSampleSize = config.OrphanReporterConfig.SampleSize
		}

		c.zombieLogger = newZombieLoggerComponent(zombieLoggerInterval, zombieLoggerSampleSize)
		go c.zombieLogger.Start()
	}

	kvServerList := routeEndpoints{}
	var srcMemdAddrs []routeEndpoint
	for _, seed := range config.SeedConfig.MemdAddrs {
		if config.SecurityConfig.UseTLS && !config.SecurityConfig.NoTLSSeedNode {
			kvServerList.SSLEndpoints = append(kvServerList.SSLEndpoints, routeEndpoint{
				Address:    seed,
				IsSeedNode: true,
			})
			srcMemdAddrs = kvServerList.SSLEndpoints
		} else {
			kvServerList.NonSSLEndpoints = append(kvServerList.NonSSLEndpoints, routeEndpoint{
				Address:    seed,
				IsSeedNode: true,
			})
			srcMemdAddrs = kvServerList.NonSSLEndpoints
		}
	}
	if config.SeedConfig.SRVRecord != nil {
		c.srvDetails = &srvDetails{
			Addrs:  kvServerList,
			Record: *config.SeedConfig.SRVRecord,
		}
	}

	var seedNodeAddr string
	if config.SecurityConfig.NoTLSSeedNode {
		host, err := parseSeedNode(config.SeedConfig.HTTPAddrs)
		if err != nil {
			return nil, err
		}

		seedNodeAddr = host
	}

	c.cfgManager = newConfigManager(
		configManagerProperties{
			NetworkType:  config.IoConfig.NetworkType,
			SrcMemdAddrs: srcMemdAddrs,
			SrcHTTPAddrs: srcHTTPAddrs,
			UseTLS:       tlsConfig != nil,
			SeedNodeAddr: seedNodeAddr,
		},
	)

	c.dialer = newMemdClientDialerComponent(
		memdClientDialerProps{
			ServerWaitTimeout:    serverWaitTimeout,
			KVConnectTimeout:     kvConnectTimeout,
			ClientID:             c.clientID,
			CompressionMinSize:   compressionMinSize,
			CompressionMinRatio:  compressionMinRatio,
			DisableDecompression: disableDecompression,
			NoTLSSeedNode:        config.SecurityConfig.NoTLSSeedNode,
			ConnBufSize:          kvBufferSize,
		},
		bootstrapProps{
			HelloProps: helloProps{
				CollectionsEnabled:             useCollections,
				MutationTokensEnabled:          useMutationTokens,
				CompressionEnabled:             useCompression,
				DurationsEnabled:               useDurations,
				OutOfOrderEnabled:              useOutOfOrder,
				JSONFeatureEnabled:             useJSONHello,
				XErrorFeatureEnabled:           useXErrorHello,
				SyncReplicationEnabled:         useSyncReplicationHello,
				PITRFeatureEnabled:             usePITRHello,
				ResourceUnitsEnabled:           useResourceUnits,
				ClusterMapNotificationsEnabled: UseClusterMapNotifications,
			},
			Bucket:        c.bucketName,
			UserAgent:     userAgent,
			ErrMapManager: c.errMap,
		},
		circuitBreakerConfig,
		c.zombieLogger,
		c.tracer,
		c.cfgManager,
	)
	c.kvMux = newKVMux(
		kvMuxProps{
			QueueSize:          maxQueueSize,
			PoolSize:           kvPoolSize,
			CollectionsEnabled: useCollections,
			NoTLSSeedNode:      config.SecurityConfig.NoTLSSeedNode,
		},
		c.cfgManager,
		c.errMap,
		c.tracer,
		c.dialer,
		&kvMuxState{
			tlsConfig:          tlsConfig,
			authMechanisms:     c.authMechanisms,
			auth:               config.SecurityConfig.Auth,
			expectedBucketName: c.bucketName,
		},
	)
	c.collections = newCollectionIDManager(
		collectionIDProps{
			MaxQueueSize:         config.KVConfig.MaxQueueSize,
			DefaultRetryStrategy: c.defaultRetryStrategy,
		},
		c.kvMux,
		c.tracer,
		c.cfgManager,
	)
	c.httpMux = newHTTPMux(
		circuitBreakerConfig,
		c.cfgManager,
		&httpClientMux{tlsConfig: tlsConfig, auth: config.SecurityConfig.Auth},
		config.SecurityConfig.NoTLSSeedNode,
	)
	c.http = newHTTPComponent(
		httpComponentProps{
			UserAgent:            userAgent,
			DefaultRetryStrategy: c.defaultRetryStrategy,
		},
		httpClientProps{
			maxIdleConns:        config.HTTPConfig.MaxIdleConns,
			maxIdleConnsPerHost: config.HTTPConfig.MaxIdleConnsPerHost,
			idleTimeout:         httpIdleConnTimeout,
			connectTimeout:      httpConnectTimeout,
		},
		c.httpMux,
		c.tracer,
	)

	var poller configPollerController
	if len(config.SeedConfig.MemdAddrs) == 0 && config.BucketName == "" {
		// The http poller can't run without a bucket. We don't trigger an error for this case
		// because AgentGroup users who use memcached buckets on non-default ports will end up here.
		logDebugf("No bucket name specified and only http addresses specified, not running config poller")
		c.diagnostics = newDiagnosticsComponent(c.kvMux, c.httpMux, c.http, c.bucketName, c.defaultRetryStrategy, nil)
	} else {
		if config.SecurityConfig.NoTLSSeedNode {
			poller = newSeedConfigController(srcHTTPAddrs[0].Address, c.bucketName,
				httpPollerProperties{
					httpComponent:        c.http,
					confHTTPRetryDelay:   confHTTPRetryDelay,
					confHTTPRedialPeriod: confHTTPRedialPeriod,
					confHTTPMaxWait:      confHTTPMaxWait,
				}, c.cfgManager)
		} else {
			var httpPoller *httpConfigController
			if c.bucketName != "" {
				httpPoller = newHTTPConfigController(
					c.bucketName,
					httpPollerProperties{
						httpComponent:        c.http,
						confHTTPRetryDelay:   confHTTPRetryDelay,
						confHTTPRedialPeriod: confHTTPRedialPeriod,
						confHTTPMaxWait:      confHTTPMaxWait,
					},
					c.httpMux,
					c.cfgManager,
				)
			}
			cccpFetcher := newCCCPConfigFetcher(confCccpMaxWait)
			poller = newPollerController(
				newCCCPConfigController(
					cccpPollerProperties{
						confCccpPollPeriod: confCccpPollPeriod,
						cccpConfigFetcher:  cccpFetcher,
					},
					c.kvMux,
					c.cfgManager,
					c.isPollingFallbackError,
					c.onCCCPNoConfigFromAnyNode,
				),
				httpPoller,
				c.cfgManager,
				c.isPollingFallbackError,
			)
			c.cfgManager.SetConfigFetcher(cccpFetcher)
		}
		c.pollerController = poller
		c.diagnostics = newDiagnosticsComponent(c.kvMux, c.httpMux, c.http, c.bucketName, c.defaultRetryStrategy, c.pollerController)
	}
	c.dialer.AddBootstrapFailHandler(c.diagnostics)
	c.dialer.AddCCCPUnsupportedHandler(c)
	c.cfgManager.AddConfigWatcher(c.dialer)

	c.observe = newObserveComponent(c.collections, c.defaultRetryStrategy, c.tracer, c.kvMux)
	c.crud = newCRUDComponent(c.collections, c.defaultRetryStrategy, c.tracer, c.errMap, c.kvMux, c.kvMux, disableDecompression)
	c.stats = newStatsComponent(c.kvMux, c.defaultRetryStrategy, c.tracer)
	c.n1ql = newN1QLQueryComponent(c.http, c.cfgManager, c.tracer)
	c.analytics = newAnalyticsQueryComponent(c.http, c.tracer)
	c.search = newSearchQueryComponent(c.http, c.cfgManager, c.tracer)
	c.views = newViewQueryComponent(c.http, c.tracer)

	// Kick everything off.
	cfg := &routeConfig{
		kvServerList: kvServerList,
		mgmtEpList:   httpEpList,
		revID:        -1,
	}

	c.httpMux.OnNewRouteConfig(cfg)
	c.kvMux.OnNewRouteConfig(cfg)

	if c.pollerController != nil {
		go c.pollerController.Run()
	}

	return c, nil
}

// Close shuts down the agent, disconnecting from all servers and failing
// any outstanding operations with ErrShutdown.
func (agent *Agent) Close() error {
	logInfof("Agent closing")
	poller := agent.pollerController
	if poller != nil {
		poller.Stop()
	}
	routeCloseErr := agent.kvMux.Close()
	agent.cfgManager.Close()

	if agent.zombieLogger != nil {
		agent.zombieLogger.Stop()
	}

	// Close the transports so that they don't hold open goroutines.
	agent.http.Close()
	close(agent.shutdownSig)

	logInfof("Agent close complete")

	return routeCloseErr
}

// ClientID returns the unique id for this agent
func (agent *Agent) ClientID() string {
	return agent.clientID
}

// MemdEps returns all the available endpoints for performing KV/DCP operations (using the memcached binary protocol).
// As apposed to other endpoints, these will have the 'couchbase(s)://' scheme prefix.
func (agent *Agent) MemdEps() []string {
	snapshot, err := agent.kvMux.PipelineSnapshot()
	if err != nil {
		return []string{}
	}
	return snapshot.state.KVEps()
}

// CapiEps returns all the available endpoints for performing
// map-reduce queries.
func (agent *Agent) CapiEps() []string {
	return agent.httpMux.CapiEps()
}

// MgmtEps returns all the available endpoints for performing
// management queries.
func (agent *Agent) MgmtEps() []string {
	return agent.httpMux.MgmtEps()
}

// N1qlEps returns all the available endpoints for performing
// N1QL queries.
func (agent *Agent) N1qlEps() []string {
	return agent.httpMux.N1qlEps()
}

// FtsEps returns all the available endpoints for performing
// FTS queries.
func (agent *Agent) FtsEps() []string {
	return agent.httpMux.FtsEps()
}

// CbasEps returns all the available endpoints for performing
// CBAS queries.
func (agent *Agent) CbasEps() []string {
	return agent.httpMux.CbasEps()
}

// EventingEps returns all the available endpoints for managing/interacting with the Eventing Service.
func (agent *Agent) EventingEps() []string {
	return agent.httpMux.EventingEps()
}

// GSIEps returns all the available endpoints for managing/interacting with the GSI Service.
func (agent *Agent) GSIEps() []string {
	return agent.httpMux.GSIEps()
}

// BackupEps returns all the available endpoints for managing/interacting with the Backup Service.
func (agent *Agent) BackupEps() []string {
	return agent.httpMux.BackupEps()
}

// HasCollectionsSupport verifies whether or not collections are available on the agent.
func (agent *Agent) HasCollectionsSupport() bool {
	return agent.kvMux.SupportsCollections()
}

// IsSecure returns whether this client is connected via SSL.
func (agent *Agent) IsSecure() bool {
	return agent.kvMux.IsSecure()
}

// UsingGCCCP returns whether or not the Agent is currently using GCCCP polling.
func (agent *Agent) UsingGCCCP() bool {
	return agent.kvMux.SupportsGCCCP()
}

// HasSeenConfig returns whether or not the Agent has seen a valid cluster config. This does not mean that the agent
// currently has active connections.
// Volatile: This API is subject to change at any time.
func (agent *Agent) HasSeenConfig() (bool, error) {
	seen, err := agent.kvMux.ConfigRev()
	if err != nil {
		return false, err
	}

	return seen > -1, nil
}

// WaitUntilReady is used to verify that the SDK has been able to establish connections to the cluster.
// If no strategy is set then a fast fail retry strategy will be applied - only RetryReason that are set to always
// retry will be retried. This includes for WaitUntilReady, that is the SDK will wait until connections succeed
// or report a connection error - as soon as a connection error is reported WaitUntilReady will fail and return that
// error.
// Connection time errors are also be subject to KvConfig.ServerWaitBackoff. This is the period of time that the SDK
// will wait before attempting to reconnect to a node.
func (agent *Agent) WaitUntilReady(deadline time.Time, opts WaitUntilReadyOptions, cb WaitUntilReadyCallback) (PendingOp, error) {
	forceWait := true
	if len(opts.ServiceTypes) == 0 {
		forceWait = false
		opts.ServiceTypes = []ServiceType{MemdService}
	}

	return agent.diagnostics.WaitUntilReady(deadline, forceWait, opts, cb)
}

// ConfigSnapshot returns a snapshot of the underlying configuration currently in use.
func (agent *Agent) ConfigSnapshot() (*ConfigSnapshot, error) {
	return agent.kvMux.ConfigSnapshot()
}

// WaitForConfigSnapshot returns a snapshot of the underlying configuration currently in use, once one is available.
// Volatile: This API is subject to change at any time.
func (agent *Agent) WaitForConfigSnapshot(deadline time.Time, opts WaitForConfigSnapshotOptions, cb WaitForConfigSnapshotCallback) (PendingOp, error) {
	return agent.kvMux.WaitForConfigSnapshot(deadline, cb)
}

// BucketName returns the name of the bucket that the agent is using, if any.
// Uncommitted: This API may change in the future.
func (agent *Agent) BucketName() string {
	return agent.bucketName
}

// ForceReconnect gracefully rebuilds all connections being used by the agent.
// Any persistent in flight requests (e.g. DCP) will be terminated with ErrForcedReconnect.
//
// Internal: This should never be used and is not supported.
func (agent *Agent) ForceReconnect() {
	agent.connectionSettingsLock.Lock()
	auth := agent.auth
	mechs := agent.authMechanisms
	tlsConfig := agent.tlsConfig
	agent.connectionSettingsLock.Unlock()
	agent.kvMux.ForceReconnect(tlsConfig, mechs, auth, true)
}

// ReconfigureSecurityOptions are the options available to the ReconfigureSecurity function.
type ReconfigureSecurityOptions struct {
	UseTLS bool
	// If is nil will default to the TLSRootCAProvider already in use by the agent.
	TLSRootCAProvider func() *x509.CertPool

	Auth AuthProvider

	// AuthMechanisms is the list of mechanisms that the SDK can use to attempt authentication.
	// Note that if you add PLAIN to the list, this will cause credential leakage on the network
	// since PLAIN sends the credentials in cleartext. It is disabled by default to prevent downgrade attacks. We
	// recommend using a TLS connection if using PLAIN.
	// If is nil will default to the AuthMechanisms already in use by the Agent.
	AuthMechanisms []AuthMechanism
}

// ReconfigureSecurity updates the security configuration being used by the agent. This includes the ability to
// toggle TLS on and off.
//
// Calling this function will cause all underlying connections to be reconnected. The exception to this is the
// connection to the seed node (usually localhost), which will only be reconnected if the AuthProvider is provided
// on the options.
//
// This function can only be called when the seed poller is in use i.e. when the ns_server scheme is used.
// Internal: This should never be used and is not supported.
func (agent *Agent) ReconfigureSecurity(opts ReconfigureSecurityOptions) error {
	_, ok := agent.pollerController.(*seedConfigController)
	if !ok {
		return errors.New("reconfigure tls is only supported when the agent is in ns server mode")
	}

	var authProvided bool
	auth := opts.Auth
	mechs := opts.AuthMechanisms
	agent.connectionSettingsLock.Lock()
	if auth == nil {
		auth = agent.auth
	} else {
		authProvided = true
	}
	if len(mechs) == 0 {
		mechs = agent.authMechanisms
	}

	var tlsConfig *dynTLSConfig
	if opts.UseTLS {
		if opts.TLSRootCAProvider == nil {
			return wrapError(errInvalidArgument, "must provide TLSRootCAProvider when UseTLS is true")
		}
		tlsConfig = createTLSConfig(auth, opts.TLSRootCAProvider)
	}

	agent.auth = auth
	agent.authMechanisms = mechs
	agent.tlsConfig = tlsConfig
	agent.connectionSettingsLock.Unlock()

	agent.cfgManager.UseTLS(tlsConfig != nil)
	agent.kvMux.ForceReconnect(tlsConfig, mechs, auth, authProvided)
	agent.httpMux.UpdateTLS(tlsConfig, auth)
	return nil
}

func (agent *Agent) onCCCPUnsupported(err error) {
	// If this error is a legitimate fallback reason then we should immediately start the http poller.
	// This should always be a poller fallback error but lets just be sure.
	if agent.pollerController != nil && agent.isPollingFallbackError(err) {
		agent.pollerController.ForceHTTPPoller()
	}
}

func (agent *Agent) isPollingFallbackError(err error) bool {
	return isPollingFallbackError(err, agent.bucketName)
}

type srvAgent interface {
	srv() *srvDetails
	setSRVAddrs(routeEndpoints)
	routeConfigWatchers() []routeConfigWatcher
	resetConfig()
	IsSecure() bool
	stopped() <-chan struct{}
}

func (agent *Agent) srv() *srvDetails {
	return agent.srvDetails
}

func (agent *Agent) setSRVAddrs(addrs routeEndpoints) {
	agent.srvDetails.Addrs = addrs
}

func (agent *Agent) routeConfigWatchers() []routeConfigWatcher {
	return agent.cfgManager.Watchers()
}

func (agent *Agent) resetConfig() {
	// Reset the config manager to accept the next config that the poller fetches.
	// This is safe to do here, we're blocking the poller from fetching a config and if we're here then
	// we can't be performing ops.
	agent.cfgManager.ResetConfig()
	// Reset the dialer so that the next connections to bootstrap fetch a config and kick off the poller again.
	agent.dialer.ResetConfig()
}

func (agent *Agent) onCCCPNoConfigFromAnyNode(err error) {
	onCCCPNoConfigFromAnyNode(agent, err)
}

func (agent *Agent) stopped() <-chan struct{} {
	return agent.shutdownSig
}

// The CCCP poller suddenly becoming unable to fetch a config from any node in the cluster is the trigger
// for checking if we need to try refresh the DNS SRV record that we used to initially connect.
// Note that we don't need locking around of this because there is only one poller active at any given time
// and we're blocking it here.
func onCCCPNoConfigFromAnyNode(agent srvAgent, err error) {
	srvDetails := agent.srv()
	if srvDetails == nil {
		return
	}

	// We only want to refresh the SRV record under certain circumstances, namely that we can't connect to the cluster.
	var opErr *net.OpError
	if !errors.As(err, &opErr) {
		return
	}

	logInfof("Refreshing SRV record: %s", srvDetails.Record)

	var addrs []*net.SRV
	for {
		_, addrs, err = net.LookupSRV(srvDetails.Record.Scheme, srvDetails.Record.Proto, srvDetails.Record.Host)
		if err != nil {
			if isLogRedactionLevelFull() {
				logInfof("Failed to lookup SRV record: %s", redactSystemData(err))
			} else {
				logInfof("Failed to lookup SRV record: %s", err)
			}
		}

		if len(addrs) > 0 {
			break
		}

		select {
		case <-agent.stopped():
			return
		case <-time.After(10 * time.Second):
		}
	}

	// If any of the addresses in the SRV record match an address that we already know then we can say that the
	// cluster has not moved and bail out.
	useTLS := agent.IsSecure()
	var memdAddrs []routeEndpoint
	if useTLS {
		memdAddrs = srvDetails.Addrs.SSLEndpoints
	} else {
		memdAddrs = srvDetails.Addrs.NonSSLEndpoints
	}

	logAddrs := append([]routeEndpoint(nil), memdAddrs...)
	if isLogRedactionLevelFull() {
		for i, addr := range logAddrs {
			logAddrs[i].Address = redactSystemData(addr)
		}
	}
	logInfof("Found new addrs for SRV record: %v", logAddrs)

	for _, addr := range addrs {
		host := fmt.Sprintf("%s:%d", strings.TrimSuffix(addr.Target, "."), addr.Port)
		for _, seed := range memdAddrs {
			if host == seed.Address {
				logInfof("Found already known matching address, not refreshing system")
				return
			}
		}
	}
	logInfof("No matching address known, refreshing system")

	agent.resetConfig()

	kvServerList := routeEndpoints{}
	for _, seed := range addrs {
		host := fmt.Sprintf("%s:%d", strings.TrimSuffix(seed.Target, "."), seed.Port)
		if useTLS {
			kvServerList.SSLEndpoints = append(kvServerList.SSLEndpoints, routeEndpoint{
				Address:    host,
				IsSeedNode: true,
			})
		} else {
			kvServerList.NonSSLEndpoints = append(kvServerList.NonSSLEndpoints, routeEndpoint{
				Address:    host,
				IsSeedNode: true,
			})
		}
	}

	// Build a new fake config to kick off the pipelines again, this will make the kvmux stop the old pipelines
	// and create new connections to the new addresses.
	newCfg := &routeConfig{
		kvServerList: kvServerList,
		revID:        -1,
	}

	watchers := agent.routeConfigWatchers()
	for _, watcher := range watchers {
		watcher.OnNewRouteConfig(newCfg)
	}

	// Update the addresses we hold so that if the SRV changes again then we can correctly check the new vs old
	// addresses.
	agent.setSRVAddrs(kvServerList)
}

func authMechanismsFromConfig(authMechanisms []AuthMechanism, useTLS bool) []AuthMechanism {
	if len(authMechanisms) == 0 {
		if useTLS {
			authMechanisms = []AuthMechanism{PlainAuthMechanism}
		} else {
			// No user specified auth mechanisms so set our defaults.
			authMechanisms = []AuthMechanism{
				ScramSha512AuthMechanism,
				ScramSha256AuthMechanism,
				ScramSha1AuthMechanism}
		}
	} else if !useTLS {
		// The user has specified their own mechanisms and not using TLS so we check if they've set PLAIN.
		for _, mech := range authMechanisms {
			if mech == PlainAuthMechanism {
				logWarnf("PLAIN sends credentials in plaintext, this will cause credential leakage on the network")
			}
		}
	}
	return authMechanisms
}

func setupTLSConfig(addrs []string, config SecurityConfig) (*dynTLSConfig, error) {
	var tlsConfig *dynTLSConfig
	if config.UseTLS {
		if config.TLSRootCAProvider == nil {
			logDebugf("TLS enabled with no root ca provider - trusting system cert pool and Capella root CA")

			pool, err := x509.SystemCertPool()
			if err != nil {
				return nil, wrapError(err, "failed to load system cert pool")
			}
			pool.AppendCertsFromPEM(capellaRootCA)

			config.TLSRootCAProvider = func() *x509.CertPool {
				return pool
			}
		}
		tlsConfig = createTLSConfig(config.Auth, config.TLSRootCAProvider)
	} else {
		var endsInCloud bool
		for _, host := range addrs {
			if strings.HasSuffix(strings.Split(host, ":")[0], ".cloud.couchbase.com") {
				endsInCloud = true
				break
			}
		}
		if endsInCloud {
			logWarnf("TLS is required when connecting to Couchbase Capella. Please enable TLS by prefixing " +
				"the connection string with \"couchbases://\" (note the final 's').")
		}
	}

	return tlsConfig, nil
}

func parseSeedNode(addrs []string) (string, error) {
	if len(addrs) != 1 {
		return "", wrapError(errInvalidArgument, "must specify exactly one seed address")
	}

	host, _, err := net.SplitHostPort(addrs[0])
	if err != nil {
		return "", wrapError(err, "cannot split host port for seed address")
	}
	ip := net.ParseIP(host)
	if !ip.IsLoopback() {
		return "", wrapError(errInvalidArgument, "seed address must be loopback")
	}

	if ip.To4() == nil && ip.To16() != nil {
		// This is a valid IP v6 address but SplitHostPort strips out the [] surrounding the host, which leads
		// the later dial to fail with "too many colons in address".
		return "[" + host + "]", nil
	}

	return host, nil
}
