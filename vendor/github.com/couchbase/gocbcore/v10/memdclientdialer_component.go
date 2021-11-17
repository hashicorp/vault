package gocbcore

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"time"
)

type memdClientDialerComponent struct {
	kvConnectTimeout  time.Duration
	serverWaitTimeout time.Duration
	clientID          string
	breakerCfg        CircuitBreakerConfig
	tlsConfig         *dynTLSConfig

	dcpQueueSize         int
	compressionMinSize   int
	compressionMinRatio  float64
	disableDecompression bool

	serverFailuresLock sync.Mutex
	serverFailures     map[string]time.Time

	tracer       *tracerComponent
	zombieLogger *zombieLoggerComponent

	bootstrapProps         bootstrapProps
	bootstrapInitFunc      memdInitFunc
	InitialBootstrapNonTLS bool

	bootstrapFailHandlersLock sync.Mutex
	bootstrapFailHandlers     []memdBoostrapFailHandler
}

type memdClientDialerProps struct {
	KVConnectTimeout       time.Duration
	ServerWaitTimeout      time.Duration
	ClientID               string
	TLSConfig              *dynTLSConfig
	DCPQueueSize           int
	CompressionMinSize     int
	CompressionMinRatio    float64
	DisableDecompression   bool
	InitialBootstrapNonTLS bool
}

type memdBoostrapFailHandler interface {
	onBootstrapFail(error)
}

func newMemdClientDialerComponent(props memdClientDialerProps, bSettings bootstrapProps, breakerCfg CircuitBreakerConfig,
	zLogger *zombieLoggerComponent, tracer *tracerComponent, bootstrapInitFunc memdInitFunc) *memdClientDialerComponent {
	return &memdClientDialerComponent{
		kvConnectTimeout:  props.KVConnectTimeout,
		serverWaitTimeout: props.ServerWaitTimeout,
		clientID:          props.ClientID,
		tlsConfig:         props.TLSConfig,
		breakerCfg:        breakerCfg,
		zombieLogger:      zLogger,
		tracer:            tracer,
		serverFailures:    make(map[string]time.Time),

		bootstrapProps:         bSettings,
		bootstrapInitFunc:      bootstrapInitFunc,
		InitialBootstrapNonTLS: props.InitialBootstrapNonTLS,

		dcpQueueSize:         props.DCPQueueSize,
		compressionMinSize:   props.CompressionMinSize,
		compressionMinRatio:  props.CompressionMinRatio,
		disableDecompression: props.DisableDecompression,
	}
}

func (mcc *memdClientDialerComponent) AddBootstrapFailHandler(handler memdBoostrapFailHandler) {
	mcc.bootstrapFailHandlersLock.Lock()
	mcc.bootstrapFailHandlers = append(mcc.bootstrapFailHandlers, handler)
	mcc.bootstrapFailHandlersLock.Unlock()
}

func (mcc *memdClientDialerComponent) RemoveBootstrapFailHandler(handler memdBoostrapFailHandler) {
	var idx int
	mcc.bootstrapFailHandlersLock.Lock()
	for i, w := range mcc.bootstrapFailHandlers {
		if w == handler {
			idx = i
		}
	}

	if idx == len(mcc.bootstrapFailHandlers) {
		mcc.bootstrapFailHandlers = mcc.bootstrapFailHandlers[:idx]
	} else {
		mcc.bootstrapFailHandlers = append(mcc.bootstrapFailHandlers[:idx], mcc.bootstrapFailHandlers[idx+1:]...)
	}
	mcc.bootstrapFailHandlersLock.Unlock()
}

func (mcc *memdClientDialerComponent) SlowDialMemdClient(cancelSig <-chan struct{}, address string,
	postCompleteHandler postCompleteErrorHandler, firstConfig bool) (*memdClient, error) {
	mcc.serverFailuresLock.Lock()
	failureTime := mcc.serverFailures[address]
	mcc.serverFailuresLock.Unlock()

	if !failureTime.IsZero() {
		waitedTime := time.Since(failureTime)
		if waitedTime < mcc.serverWaitTimeout {
			select {
			case <-cancelSig:
				return nil, errRequestCanceled
			case <-time.After(mcc.serverWaitTimeout - waitedTime):
			}
		}
	}

	deadline := time.Now().Add(mcc.kvConnectTimeout)
	client, err := mcc.dialMemdClient(cancelSig, address, deadline, postCompleteHandler, firstConfig)
	if err != nil {
		if !errors.Is(err, ErrRequestCanceled) {
			mcc.serverFailuresLock.Lock()
			mcc.serverFailures[address] = time.Now()
			mcc.serverFailuresLock.Unlock()
		}

		return nil, err
	}

	err = client.Bootstrap(cancelSig, mcc.bootstrapProps, deadline, mcc.bootstrapInitFunc)
	if err != nil {
		closeErr := client.Close(nil)
		if closeErr != nil {
			logWarnf("Failed to close authentication client (%s)", closeErr)
		}
		if !errors.Is(err, ErrForcedReconnect) {
			mcc.serverFailuresLock.Lock()
			mcc.serverFailures[address] = time.Now()
			mcc.serverFailuresLock.Unlock()
		}

		mcc.bootstrapFailHandlersLock.Lock()
		handlers := make([]memdBoostrapFailHandler, len(mcc.bootstrapFailHandlers))
		copy(handlers, mcc.bootstrapFailHandlers)
		mcc.bootstrapFailHandlersLock.Unlock()
		for _, handler := range handlers {
			handler.onBootstrapFail(err)
		}

		return nil, err
	}

	return client, nil
}

func (mcc *memdClientDialerComponent) dialMemdClient(cancelSig <-chan struct{}, address string, deadline time.Time,
	postCompleteHandler postCompleteErrorHandler, firstConfig bool) (*memdClient, error) {
	// Copy the tls configuration since we need to provide the hostname for each
	// server that we connect to so that the certificate can be validated properly.
	var tlsConfig *tls.Config
	if mcc.tlsConfig != nil && !(mcc.InitialBootstrapNonTLS && firstConfig) {
		srvTLSConfig, err := mcc.tlsConfig.MakeForAddr(address)
		if err != nil {
			return nil, err
		}

		tlsConfig = srvTLSConfig
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-ctx.Done():
			return
		case <-cancelSig:
			cancel()
		}
	}()

	conn, err := dialMemdConn(ctx, address, tlsConfig, deadline)
	cancel()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			err = errRequestCanceled
		} else {
			err = wrapError(err, "check server ports and cluster encryption setting")
		}

		logDebugf("Failed to connect. %v", err)
		return nil, err
	}

	client := newMemdClient(
		memdClientProps{
			ClientID:             mcc.clientID,
			DCPQueueSize:         mcc.dcpQueueSize,
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
