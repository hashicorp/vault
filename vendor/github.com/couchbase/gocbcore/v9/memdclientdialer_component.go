package gocbcore

import (
	"crypto/tls"
	"sync"
	"time"
)

type memdClientDialerComponent struct {
	kvConnectTimeout  time.Duration
	serverWaitTimeout time.Duration
	clientID          string
	breakerCfg        CircuitBreakerConfig
	tlsConfig         *dynTLSConfig

	compressionMinSize   int
	compressionMinRatio  float64
	disableDecompression bool

	serverFailuresLock sync.Mutex
	serverFailures     map[string]time.Time

	tracer       *tracerComponent
	zombieLogger *zombieLoggerComponent

	bootstrapProps bootstrapProps
	bootstrapCB    memdInitFunc
}

type memdClientDialerProps struct {
	KVConnectTimeout     time.Duration
	ServerWaitTimeout    time.Duration
	ClientID             string
	TLSConfig            *dynTLSConfig
	CompressionMinSize   int
	CompressionMinRatio  float64
	DisableDecompression bool
}

func newMemdClientDialerComponent(props memdClientDialerProps, bSettings bootstrapProps, breakerCfg CircuitBreakerConfig,
	zLogger *zombieLoggerComponent, tracer *tracerComponent, bootstrapCB memdInitFunc) *memdClientDialerComponent {
	return &memdClientDialerComponent{
		kvConnectTimeout:  props.KVConnectTimeout,
		serverWaitTimeout: props.ServerWaitTimeout,
		clientID:          props.ClientID,
		tlsConfig:         props.TLSConfig,
		breakerCfg:        breakerCfg,
		zombieLogger:      zLogger,
		tracer:            tracer,
		serverFailures:    make(map[string]time.Time),

		bootstrapProps: bSettings,
		bootstrapCB:    bootstrapCB,

		compressionMinSize:   props.CompressionMinSize,
		compressionMinRatio:  props.CompressionMinRatio,
		disableDecompression: props.DisableDecompression,
	}
}

func (mcc *memdClientDialerComponent) SlowDialMemdClient(address string, postCompleteHandler postCompleteErrorHandler) (*memdClient, error) {
	mcc.serverFailuresLock.Lock()
	failureTime := mcc.serverFailures[address]
	mcc.serverFailuresLock.Unlock()

	if !failureTime.IsZero() {
		waitedTime := time.Since(failureTime)
		if waitedTime < mcc.serverWaitTimeout {
			time.Sleep(mcc.serverWaitTimeout - waitedTime)
		}
	}

	deadline := time.Now().Add(mcc.kvConnectTimeout)
	client, err := mcc.dialMemdClient(address, deadline, postCompleteHandler)
	if err != nil {
		mcc.serverFailuresLock.Lock()
		mcc.serverFailures[address] = time.Now()
		mcc.serverFailuresLock.Unlock()

		return nil, err
	}

	err = client.Bootstrap(mcc.bootstrapProps, deadline, mcc.bootstrapCB)
	if err != nil {
		closeErr := client.Close()
		if closeErr != nil {
			logWarnf("Failed to close authentication client (%s)", closeErr)
		}
		mcc.serverFailuresLock.Lock()
		mcc.serverFailures[address] = time.Now()
		mcc.serverFailuresLock.Unlock()

		return nil, err
	}

	return client, nil
}

func (mcc *memdClientDialerComponent) dialMemdClient(address string, deadline time.Time,
	postCompleteHandler postCompleteErrorHandler) (*memdClient, error) {
	// Copy the tls configuration since we need to provide the hostname for each
	// server that we connect to so that the certificate can be validated properly.
	var tlsConfig *tls.Config
	if mcc.tlsConfig != nil {
		srvTLSConfig, err := mcc.tlsConfig.MakeForAddr(address)
		if err != nil {
			return nil, err
		}

		tlsConfig = srvTLSConfig
	}

	conn, err := dialMemdConn(address, tlsConfig, deadline)
	if err != nil {
		logDebugf("Failed to connect. %v", err)
		return nil, err
	}

	client := newMemdClient(
		memdClientProps{
			ClientID:             mcc.clientID,
			DisableDecompression: mcc.disableDecompression,
			CompressionMinRatio:  mcc.compressionMinRatio,
			CompressionMinSize:   mcc.compressionMinSize,
		},
		conn,
		mcc.breakerCfg,
		postCompleteHandler,
		mcc.tracer,
		mcc.zombieLogger,
	)

	return client, err
}
