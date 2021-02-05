package gocbcore

import (
	"fmt"
	"time"

	"github.com/couchbase/gocbcore/v9/memd"
)

// DCPAgent represents the base client handling DCP connections to a Couchbase Server.
type DCPAgent struct {
	clientID   string
	bucketName string
	tlsConfig  *dynTLSConfig
	initFn     memdInitFunc

	pollerController *pollerController
	kvMux            *kvMux
	httpMux          *httpMux

	cfgManager  *configManagementComponent
	errMap      *errMapComponent
	tracer      *tracerComponent
	diagnostics *diagnosticsComponent
	dcp         *dcpComponent
	http        *httpComponent
}

// CreateDcpAgent creates an agent for performing DCP operations.
func CreateDcpAgent(config *DCPAgentConfig, dcpStreamName string, openFlags memd.DcpOpenFlag) (*DCPAgent, error) {
	logInfof("SDK Version: gocbcore/%s", goCbCoreVersionStr)
	logInfof("Creating new dcp agent: %+v", config)

	auth := config.Auth
	userAgent := config.UserAgent
	disableDecompression := config.DisableDecompression
	useCompression := config.UseCompression
	useCollections := config.UseCollections
	useJSONHello := !config.DisableJSONHello
	useXErrorHello := !config.DisableXErrorHello
	useSyncReplicationHello := !config.DisableSyncReplicationHello
	dcpBufferSize := 8 * 1024 * 1024
	compressionMinSize := 32
	compressionMinRatio := 0.83
	dcpBackfillOrderStr := ""
	dcpPriorityStr := ""

	kvConnectTimeout := 7000 * time.Millisecond
	if config.KVConnectTimeout > 0 {
		kvConnectTimeout = config.KVConnectTimeout
	}

	serverWaitTimeout := 5 * time.Second

	kvPoolSize := 1
	if config.KvPoolSize > 0 {
		kvPoolSize = config.KvPoolSize
	}

	maxQueueSize := 2048
	if config.MaxQueueSize > 0 {
		maxQueueSize = config.MaxQueueSize
	}

	confCccpMaxWait := 3 * time.Second
	if config.CccpMaxWait > 0 {
		confCccpMaxWait = config.CccpMaxWait
	}

	confCccpPollPeriod := 2500 * time.Millisecond
	if config.CccpPollPeriod > 0 {
		confCccpPollPeriod = config.CccpPollPeriod
	}

	confHTTPRetryDelay := 10 * time.Second
	if config.HTTPRetryDelay > 0 {
		confHTTPRetryDelay = config.HTTPRetryDelay
	}

	confHTTPRedialPeriod := 10 * time.Second
	if config.HTTPRedialPeriod > 0 {
		confHTTPRedialPeriod = config.HTTPRedialPeriod
	}

	if config.CompressionMinSize > 0 {
		compressionMinSize = config.CompressionMinSize
	}
	if config.CompressionMinRatio > 0 {
		compressionMinRatio = config.CompressionMinRatio
		if compressionMinRatio >= 1.0 {
			compressionMinRatio = 1.0
		}
	}

	if config.DCPBufferSize > 0 {
		dcpBufferSize = config.DCPBufferSize
	}
	dcpQueueSize := (dcpBufferSize + 23) / 24

	switch config.AgentPriority {
	case DcpAgentPriorityLow:
		dcpPriorityStr = "low"
	case DcpAgentPriorityMed:
		dcpPriorityStr = "medium"
	case DcpAgentPriorityHigh:
		dcpPriorityStr = "high"
	}

	// If the user doesn't explicitly set the backfill order, the DCP control flag will not be sent to the cluster
	// and the default will implicitly be used (which is 'round-robin').
	switch config.BackfillOrder {
	case DCPBackfillOrderRoundRobin:
		dcpBackfillOrderStr = "round-robin"
	case DCPBackfillOrderSequential:
		dcpBackfillOrderStr = "sequential"
	}

	authMechanisms := []AuthMechanism{
		ScramSha512AuthMechanism,
		ScramSha256AuthMechanism,
		ScramSha1AuthMechanism}

	// PLAIN authentication is only supported over TLS
	if config.UseTLS {
		authMechanisms = append(authMechanisms, PlainAuthMechanism)
	}

	var tlsConfig *dynTLSConfig
	if config.UseTLS {
		tlsConfig = createTLSConfig(config.Auth, config.TLSRootCAProvider)
	}

	httpCli := createHTTPClient(config.HTTPMaxIdleConns, config.HTTPMaxIdleConnsPerHost,
		config.HTTPIdleConnectionTimeout, tlsConfig)

	tracerCmpt := newTracerComponent(noopTracer{}, config.BucketName, false)

	// We wrap the authorization system to force DCP channel opening
	//   as part of the "initialization" for any servers.
	initFn := func(client *memdClient, deadline time.Time) error {
		sclient := &syncClient{client: client}
		if err := sclient.ExecOpenDcpConsumer(dcpStreamName, openFlags, deadline); err != nil {
			return err
		}

		if err := sclient.ExecEnableDcpNoop(180*time.Second, deadline); err != nil {
			return err
		}

		if dcpPriorityStr != "" {
			if err := sclient.ExecDcpControl("set_priority", dcpPriorityStr, deadline); err != nil {
				return err
			}
		}

		if config.UseExpiryOpcode {
			if err := sclient.ExecDcpControl("enable_expiry_opcode", "true", deadline); err != nil {
				return err
			}
		}

		if config.UseStreamID {
			if err := sclient.ExecDcpControl("enable_stream_id", "true", deadline); err != nil {
				return err
			}
		}

		if config.UseOSOBackfill {
			if err := sclient.ExecDcpControl("enable_out_of_order_snapshots", "true", deadline); err != nil {
				return err
			}
		}

		if dcpBackfillOrderStr != "" {
			if err := sclient.ExecDcpControl("backfill_order", dcpBackfillOrderStr, deadline); err != nil {
				return err
			}
		}

		if err := sclient.ExecEnableDcpBufferAck(dcpBufferSize, deadline); err != nil {
			return err
		}

		return sclient.ExecEnableDcpClientEnd(deadline)
	}

	c := &DCPAgent{
		clientID:   formatCbUID(randomCbUID()),
		bucketName: config.BucketName,
		tlsConfig:  tlsConfig,
		initFn:     initFn,
		tracer:     tracerCmpt,

		errMap: newErrMapManager(config.BucketName),
	}

	circuitBreakerConfig := CircuitBreakerConfig{
		Enabled: false,
	}

	authHandler := buildAuthHandler(auth)

	var httpEpList []string
	for _, hostPort := range config.HTTPAddrs {
		if !c.IsSecure() {
			httpEpList = append(httpEpList, fmt.Sprintf("http://%s", hostPort))
		} else {
			httpEpList = append(httpEpList, fmt.Sprintf("https://%s", hostPort))
		}
	}

	c.cfgManager = newConfigManager(
		configManagerProperties{
			NetworkType:  config.NetworkType,
			UseSSL:       config.UseTLS,
			SrcMemdAddrs: config.MemdAddrs,
			SrcHTTPAddrs: []string{},
		},
	)

	dialer := newMemdClientDialerComponent(
		memdClientDialerProps{
			ServerWaitTimeout:    serverWaitTimeout,
			KVConnectTimeout:     kvConnectTimeout,
			ClientID:             c.clientID,
			TLSConfig:            c.tlsConfig,
			DCPQueueSize:         dcpQueueSize,
			CompressionMinSize:   compressionMinSize,
			CompressionMinRatio:  compressionMinRatio,
			DisableDecompression: disableDecompression,
		},
		bootstrapProps{
			HelloProps: helloProps{
				CollectionsEnabled:     useCollections,
				CompressionEnabled:     useCompression,
				JSONFeatureEnabled:     useJSONHello,
				XErrorFeatureEnabled:   useXErrorHello,
				SyncReplicationEnabled: useSyncReplicationHello,
			},
			Bucket:         c.bucketName,
			UserAgent:      userAgent,
			AuthMechanisms: authMechanisms,
			AuthHandler:    authHandler,
			ErrMapManager:  c.errMap,
		},
		circuitBreakerConfig,
		nil,
		c.tracer,
		initFn,
		c,
	)
	c.kvMux = newKVMux(
		kvMuxProps{
			QueueSize:          maxQueueSize,
			PoolSize:           kvPoolSize,
			CollectionsEnabled: useCollections,
		},
		c.cfgManager,
		c.errMap,
		c.tracer,
		dialer,
	)
	c.httpMux = newHTTPMux(circuitBreakerConfig, c.cfgManager)
	c.http = newHTTPComponent(
		httpComponentProps{
			UserAgent:            userAgent,
			DefaultRetryStrategy: &failFastRetryStrategy{},
		},
		httpCli,
		c.httpMux,
		auth,
		c.tracer,
	)

	c.pollerController = newPollerController(
		newCCCPConfigController(
			cccpPollerProperties{
				confCccpMaxWait:    confCccpMaxWait,
				confCccpPollPeriod: confCccpPollPeriod,
			},
			c.kvMux,
			c.cfgManager,
		),
		newHTTPConfigController(
			c.bucketName,
			httpPollerProperties{
				httpComponent:        c.http,
				confHTTPRetryDelay:   confHTTPRetryDelay,
				confHTTPRedialPeriod: confHTTPRedialPeriod,
			},
			c.httpMux,
			c.cfgManager,
		),
		c.cfgManager,
	)

	c.diagnostics = newDiagnosticsComponent(c.kvMux, nil, nil, c.bucketName, newFailFastRetryStrategy(), c.pollerController)
	c.dcp = newDcpComponent(c.kvMux, config.UseStreamID)

	// Kick everything off.
	cfg := &routeConfig{
		kvServerList: config.MemdAddrs,
		mgmtEpList:   httpEpList,
		revID:        -1,
	}

	c.httpMux.OnNewRouteConfig(cfg)
	c.kvMux.OnNewRouteConfig(cfg)

	go c.pollerController.Start()

	return c, nil
}

// IsSecure returns whether this client is connected via SSL.
func (agent *DCPAgent) IsSecure() bool {
	return agent.tlsConfig != nil
}

// Close shuts down the agent, disconnecting from all servers and failing
// any outstanding operations with ErrShutdown.
func (agent *DCPAgent) Close() error {
	routeCloseErr := agent.kvMux.Close()
	agent.pollerController.Stop()

	// Wait for our external looper goroutines to finish, note that if the
	// specific looper wasn't used, it will be a nil value otherwise it
	// will be an open channel till its closed to signal completion.
	<-agent.pollerController.Done()

	return routeCloseErr
}

// WaitUntilReady returns whether or not the Agent has seen a valid cluster config.
func (agent *DCPAgent) WaitUntilReady(deadline time.Time, opts WaitUntilReadyOptions,
	cb WaitUntilReadyCallback) (PendingOp, error) {
	return agent.diagnostics.WaitUntilReady(deadline, opts, cb)
}

// OpenStream opens a DCP stream for a particular VBucket and, optionally, filter.
func (agent *DCPAgent) OpenStream(vbID uint16, flags memd.DcpStreamAddFlag, vbUUID VbUUID, startSeqNo,
	endSeqNo, snapStartSeqNo, snapEndSeqNo SeqNo, evtHandler StreamObserver, opts OpenStreamOptions,
	cb OpenStreamCallback) (PendingOp, error) {
	return agent.dcp.OpenStream(vbID, flags, vbUUID, startSeqNo, endSeqNo, snapStartSeqNo, snapEndSeqNo, evtHandler, opts, cb)
}

// CloseStream shuts down an open stream for the specified VBucket.
func (agent *DCPAgent) CloseStream(vbID uint16, opts CloseStreamOptions, cb CloseStreamCallback) (PendingOp, error) {
	return agent.dcp.CloseStream(vbID, opts, cb)
}

// GetFailoverLog retrieves the fail-over log for a particular VBucket.  This is used
// to resume an interrupted stream after a node fail-over has occurred.
func (agent *DCPAgent) GetFailoverLog(vbID uint16, cb GetFailoverLogCallback) (PendingOp, error) {
	return agent.dcp.GetFailoverLog(vbID, cb)
}

// GetVbucketSeqnos returns the last checkpoint for a particular VBucket.  This is useful
// for starting a DCP stream from wherever the server currently is.
func (agent *DCPAgent) GetVbucketSeqnos(serverIdx int, state memd.VbucketState, opts GetVbucketSeqnoOptions,
	cb GetVBucketSeqnosCallback) (PendingOp, error) {
	return agent.dcp.GetVbucketSeqnos(serverIdx, state, opts, cb)
}

// HasCollectionsSupport verifies whether or not collections are available on the agent.
func (agent *DCPAgent) HasCollectionsSupport() bool {
	return agent.kvMux.SupportsCollections()
}

// ConfigSnapshot returns a snapshot of the underlying configuration currently in use.
func (agent *DCPAgent) ConfigSnapshot() (*ConfigSnapshot, error) {
	return agent.kvMux.ConfigSnapshot()
}

func (agent *DCPAgent) onBootstrapFail(err error) {
	// If this error is a legitimate fallback reason then we should immediately start the http poller.
	if agent.pollerController != nil && isPollingFallbackError(err) {
		agent.pollerController.ForceHTTPPoller()
	}
}
