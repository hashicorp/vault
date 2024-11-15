package gocbcore

import (
	"context"
	"crypto/tls"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/couchbase/gocbcore/v10/memd"
)

type helloProps struct {
	MutationTokensEnabled          bool
	CollectionsEnabled             bool
	CompressionEnabled             bool
	DurationsEnabled               bool
	OutOfOrderEnabled              bool
	JSONFeatureEnabled             bool
	XErrorFeatureEnabled           bool
	SyncReplicationEnabled         bool
	PITRFeatureEnabled             bool
	ResourceUnitsEnabled           bool
	ClusterMapNotificationsEnabled bool
}

type bootstrapProps struct {
	Bucket        string
	UserAgent     string
	ErrMapManager *errMapComponent
	HelloProps    helloProps
}

type memdClientDialerComponent struct {
	kvConnectTimeout  time.Duration
	serverWaitTimeout time.Duration
	clientID          string
	breakerCfg        CircuitBreakerConfig

	compressionMinSize   int
	compressionMinRatio  float64
	disableDecompression bool
	connBufSize          uint

	serverFailuresLock sync.Mutex
	serverFailures     map[string]time.Time

	tracer       *tracerComponent
	zombieLogger *zombieLoggerComponent

	bootstrapProps bootstrapProps

	bootstrapFailHandlersLock sync.Mutex
	bootstrapFailHandlers     []memdBoostrapFailHandler

	cccpUnsupportedHandlersLock sync.Mutex
	cccpUnsupportedFailHandlers []memdBoostrapCCCPUnsupportedHandler

	configApplied uint32

	noTLSSeedNode bool

	dcpBootstrapProps *memdBootstrapDCPProps
	dcpQueueSize      int

	cfgManager *configManagementComponent
}

type memdBootstrapDCPProps struct {
	disableBufferAcknowledgement bool
	useOSOBackfill               bool
	useStreamID                  bool
	useChangeStreams             bool
	useExpiryOpcode              bool
	backfillOrderStr             string
	priorityStr                  string
	streamName                   string
	openFlags                    memd.DcpOpenFlag
	bufferSize                   int
}

type memdClientDialerProps struct {
	KVConnectTimeout     time.Duration
	ServerWaitTimeout    time.Duration
	ClientID             string
	CompressionMinSize   int
	CompressionMinRatio  float64
	DisableDecompression bool
	NoTLSSeedNode        bool
	ConnBufSize          uint

	DCPBootstrapProps *memdBootstrapDCPProps
	DCPQueueSize      int
}

type memdBoostrapFailHandler interface {
	onBootstrapFail(error)
}

type memdBoostrapCCCPUnsupportedHandler interface {
	onCCCPUnsupported(error)
}

func newMemdClientDialerComponent(props memdClientDialerProps, bSettings bootstrapProps, breakerCfg CircuitBreakerConfig,
	zLogger *zombieLoggerComponent, tracer *tracerComponent, cfgManager *configManagementComponent) *memdClientDialerComponent {
	dialer := &memdClientDialerComponent{
		kvConnectTimeout:  props.KVConnectTimeout,
		serverWaitTimeout: props.ServerWaitTimeout,
		clientID:          props.ClientID,
		breakerCfg:        breakerCfg,
		zombieLogger:      zLogger,
		tracer:            tracer,
		serverFailures:    make(map[string]time.Time),

		bootstrapProps: bSettings,

		dcpBootstrapProps:    props.DCPBootstrapProps,
		dcpQueueSize:         props.DCPQueueSize,
		compressionMinSize:   props.CompressionMinSize,
		compressionMinRatio:  props.CompressionMinRatio,
		disableDecompression: props.DisableDecompression,
		noTLSSeedNode:        props.NoTLSSeedNode,
		connBufSize:          props.ConnBufSize,

		cfgManager: cfgManager,
	}

	cfgManager.AddConfigWatcher(dialer)
	return dialer
}

func (mcc *memdClientDialerComponent) ResetConfig() {
	atomic.StoreUint32(&mcc.configApplied, 0)
	mcc.cfgManager.AddConfigWatcher(mcc)
}

func (mcc *memdClientDialerComponent) OnNewRouteConfig(cfg *routeConfig) {
	if cfg.revID == -1 {
		return
	}

	atomic.StoreUint32(&mcc.configApplied, 1)
	mcc.cfgManager.RemoveConfigWatcher(mcc)
}

func (mcc *memdClientDialerComponent) AddBootstrapFailHandler(handler memdBoostrapFailHandler) {
	mcc.bootstrapFailHandlersLock.Lock()
	mcc.bootstrapFailHandlers = append(mcc.bootstrapFailHandlers, handler)
	mcc.bootstrapFailHandlersLock.Unlock()
}

func (mcc *memdClientDialerComponent) AddCCCPUnsupportedHandler(handler memdBoostrapCCCPUnsupportedHandler) {
	mcc.cccpUnsupportedHandlersLock.Lock()
	mcc.cccpUnsupportedFailHandlers = append(mcc.cccpUnsupportedFailHandlers, handler)
	mcc.cccpUnsupportedHandlersLock.Unlock()
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

func (mcc *memdClientDialerComponent) SlowDialMemdClient(cancelSig <-chan struct{}, address routeEndpoint, tlsConfig *dynTLSConfig,
	auth AuthProvider, authMechanisms []AuthMechanism, postCompleteHandler postCompleteErrorHandler,
	serverRequestHandler serverRequestHandler) (*memdClient, error) {
	mcc.serverFailuresLock.Lock()
	failureTime := mcc.serverFailures[address.Address]
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
	client, err := mcc.dialMemdClient(cancelSig, address, deadline, postCompleteHandler, tlsConfig, serverRequestHandler)
	if err != nil {
		if !errors.Is(err, ErrRequestCanceled) {
			mcc.serverFailuresLock.Lock()
			mcc.serverFailures[address.Address] = time.Now()
			mcc.serverFailuresLock.Unlock()
		}

		return nil, err
	}

	bClient := newMemdBootstrapClient(client, cancelSig)
	if mcc.dcpBootstrapProps == nil {
		err = mcc.bootstrap(bClient, deadline, authMechanisms, auth)
	} else {
		err = mcc.dcpBootstrap(newDCPBootstrapClient(bClient), deadline, authMechanisms, auth)
	}
	if err != nil {
		closeErr := client.Close()
		if closeErr != nil {
			logWarnf("Failed to close authentication client (%s)", closeErr)
		}
		if !errors.Is(err, ErrForcedReconnect) {
			mcc.serverFailuresLock.Lock()
			mcc.serverFailures[address.Address] = time.Now()
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

func (mcc *memdClientDialerComponent) dialMemdClient(cancelSig <-chan struct{}, address routeEndpoint, deadline time.Time,
	postCompleteHandler postCompleteErrorHandler, dynTls *dynTLSConfig, serverRequestHandler serverRequestHandler) (*memdClient, error) {
	// Copy the tls configuration since we need to provide the hostname for each
	// server that we connect to so that the certificate can be validated properly.
	var tlsConfig *tls.Config
	if dynTls != nil && !(mcc.noTLSSeedNode && address.IsSeedNode) {
		srvTLSConfig, err := dynTls.MakeForAddr(address.Address)
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

	conn, err := dialMemdConn(ctx, address.Address, tlsConfig, deadline, mcc.connBufSize)
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
		serverRequestHandler,
	)

	return client, err
}

func (mcc *memdClientDialerComponent) dcpBootstrap(client *dcpBootstrapClient, deadline time.Time,
	authMechanisms []AuthMechanism, authProvider AuthProvider) error {
	if err := mcc.bootstrap(client, deadline, authMechanisms, authProvider); err != nil {
		return err
	}

	if err := client.ExecOpenDcpConsumer(mcc.dcpBootstrapProps.streamName, mcc.dcpBootstrapProps.openFlags, deadline); err != nil {
		return err
	}

	if err := client.ExecEnableDcpNoop(180*time.Second, deadline); err != nil {
		return err
	}

	if mcc.dcpBootstrapProps.priorityStr != "" {
		if err := client.ExecDcpControl("set_priority", mcc.dcpBootstrapProps.priorityStr, deadline); err != nil {
			return err
		}
	}

	if mcc.dcpBootstrapProps.useChangeStreams {
		if err := client.ExecDcpControl("change_streams", "true", deadline); err != nil {
			return err
		}
	}

	if mcc.dcpBootstrapProps.useExpiryOpcode {
		if err := client.ExecDcpControl("enable_expiry_opcode", "true", deadline); err != nil {
			return err
		}
	}

	if mcc.dcpBootstrapProps.useStreamID {
		if err := client.ExecDcpControl("enable_stream_id", "true", deadline); err != nil {
			return err
		}
	}

	if mcc.dcpBootstrapProps.useOSOBackfill {
		if err := client.ExecDcpControl("enable_out_of_order_snapshots", "true_with_seqno_advanced", deadline); err != nil {
			return err
		}
	}

	if mcc.dcpBootstrapProps.backfillOrderStr != "" {
		if err := client.ExecDcpControl("backfill_order", mcc.dcpBootstrapProps.backfillOrderStr, deadline); err != nil {
			return err
		}
	}

	if !mcc.dcpBootstrapProps.disableBufferAcknowledgement {
		if err := client.ExecEnableDcpBufferAck(mcc.dcpBootstrapProps.bufferSize, deadline); err != nil {
			return err
		}
	}

	return client.ExecEnableDcpClientEnd(deadline)
}

func (mcc *memdClientDialerComponent) bootstrap(client bootstrapClient, deadline time.Time,
	authMechanisms []AuthMechanism, authProvider AuthProvider) error {
	logDebugf("Memdclient %s Fetching cluster client data", client.LoggerID())

	bucket := mcc.bootstrapProps.Bucket
	features := helloFeatures(mcc.bootstrapProps.HelloProps)
	clientInfoStr := clientInfoString(client.ConnID(), mcc.bootstrapProps.UserAgent)

	helloCh, err := client.ExecHello(clientInfoStr, features, deadline)
	if err != nil {
		logDebugf("Memdclient %s Failed to execute HELLO (%v)", client.LoggerID(), err)
		return err
	}

	errMapCh, err := client.ExecGetErrorMap(2, deadline)
	if err != nil {
		// GetErrorMap isn't integral to bootstrap succeeding
		logDebugf("Memdclient %s Failed to execute Get error map (%v)", client.LoggerID(), err)
	}

	var listMechsCh chan SaslListMechsCompleted
	var completedAuthCh chan error
	var continueAuthCh chan bool

	firstAuthMethod := mcc.buildAuthHandler(client, authProvider, deadline, authMechanisms[0])

	if firstAuthMethod != nil {
		// If the auth method is nil then we don't actually need to do any auth so no need to Get the mechanisms.
		listMechsCh = make(chan SaslListMechsCompleted, 1)
		err = client.SaslListMechs(deadline, func(mechs []AuthMechanism, err error) {
			if err != nil {
				logDebugf("Memdclient %s Failed to fetch list auth mechs (%v)", client.LoggerID(), err)
			}
			listMechsCh <- SaslListMechsCompleted{
				Err:   err,
				Mechs: mechs,
			}
		})
		if err != nil {
			logDebugf("Memdclient %s Failed to execute list auth mechs (%v)", client.LoggerID(), err)
		}

		completedAuthCh, continueAuthCh, err = firstAuthMethod()
		if err != nil {
			logDebugf("Memdclient %s Failed to execute auth (%v)", client.LoggerID(), err)
			return err
		}
	}

	var selectCh chan error
	var configCh chan getConfigResponse
	// If there's no bucket then we don't need to do select bucket, we also don't need to wait for the continue channel,
	// as it will never be read and will be garbage collected.
	if continueAuthCh == nil {
		if bucket != "" {
			selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
			if err != nil {
				logDebugf("Memdclient %s Failed to execute select bucket (%v)", client.LoggerID(), err)
				return err
			}
		}
		if atomic.LoadUint32(&mcc.configApplied) == 0 {
			configCh, err = client.ExecGetConfig(deadline)
			if err != nil {
				// Getting a config isn't essential to bootstrap.
				logDebugf("Memdclient %s Failed to execute get config (%v)", client.LoggerID(), err)
			}
		}
	} else {
		selectCh, configCh = mcc.continueAfterAuth(client, bucket, continueAuthCh, deadline)
	}

	helloResp := <-helloCh
	if helloResp.Err != nil {
		logDebugf("Memdclient %s Failed to hello with server (%v)", client.LoggerID(), helloResp.Err)
		return helloResp.Err
	}

	if errMapCh != nil {
		errMapResp := <-errMapCh
		if errMapResp.Err == nil {
			mcc.bootstrapProps.ErrMapManager.StoreErrorMap(errMapResp.Bytes)
		} else {
			logDebugf("Memdclient %s Failed to fetch kv error map (%s)", client.LoggerID(), errMapResp.Err)
		}
	}

	var serverAuthMechanisms []AuthMechanism
	if listMechsCh != nil {
		listMechsResp := <-listMechsCh
		if listMechsResp.Err == nil {
			serverAuthMechanisms = listMechsResp.Mechs
			logDebugf("Memdclient %s Server supported auth mechanisms: %v", client.LoggerID(), serverAuthMechanisms)
		} else {
			logDebugf("Memdclient %s Failed to fetch auth mechs from server (%v)", client.LoggerID(), listMechsResp.Err)
		}
	}

	// If completedAuthCh isn't nil then we have attempted to do auth so we need to wait on the result of that.
	if completedAuthCh != nil {
		authErr := <-completedAuthCh
		if authErr != nil {
			logDebugf("Memdclient %s Failed to perform auth against server (%v)", client.LoggerID(), authErr)
			if errors.Is(authErr, ErrRequestCanceled) {
				// There's no point in us trying different mechanisms if something has cancelled bootstrapping.
				return authErr
			} else if errors.Is(authErr, ErrAuthenticationFailure) {
				// If there's only one auth mechanism then we can just fail.
				if len(authMechanisms) == 1 {
					return authErr
				}
				// If the server supports the mechanism we've tried then this auth error can't be due to an unsupported
				// mechanism.
				for _, mech := range serverAuthMechanisms {
					if mech == authMechanisms[0] {
						return authErr
					}
				}

				// If we've got here then the auth mechanism we tried is unsupported so let's keep trying with the next
				// supported mechanism.
				logInfof("Memdclient `%s` Unsupported authentication mechanism, will attempt to find next supported mechanism", client.ConnID())
			}

			for {
				var found bool
				var mech AuthMechanism
				found, mech, authMechanisms = findNextAuthMechanism(authMechanisms, serverAuthMechanisms)
				if !found {
					logDebugf("Memdclient %s Failed to authenticate, all options exhausted", client.LoggerID())
					return authErr
				}

				logDebugf("Memdclient %s Retrying authentication with found supported mechanism: %s", client.LoggerID(), mech)
				nextAuthFunc := mcc.buildAuthHandler(client, authProvider, deadline, mech)
				if nextAuthFunc == nil {
					// This can't really happen but just in case it somehow does.
					logInfof("Memdclient `%p` Failed to authenticate, no available credentials", client)
					return authErr
				}
				completedAuthCh, continueAuthCh, err = nextAuthFunc()
				if err != nil {
					logDebugf("Memdclient %s Failed to execute auth (%v)", client.LoggerID(), err)
					return err
				}
				if continueAuthCh == nil {
					if bucket != "" {
						selectCh, err = client.ExecSelectBucket([]byte(bucket), deadline)
						if err != nil {
							logDebugf("Memdclient %s Failed to execute select bucket (%v)", client.LoggerID(), err)
							return err
						}
					}
					if atomic.LoadUint32(&mcc.configApplied) == 0 {
						configCh, err = client.ExecGetConfig(deadline)
						if err != nil {
							// Getting a config isn't essential to bootstrap.
							logDebugf("Memdclient %s Failed to execute get config (%v)", client.LoggerID(), err)
						}
					}
				} else {
					selectCh, configCh = mcc.continueAfterAuth(client, bucket, continueAuthCh, deadline)
				}
				authErr = <-completedAuthCh
				if authErr == nil {
					break
				}

				logDebugf("Memdclient %s Failed to perform auth against server (%v)", client.LoggerID(), authErr)
				if errors.Is(authErr, ErrAuthenticationFailure) || errors.Is(err, ErrRequestCanceled) {
					return authErr
				}
			}
		}
		logDebugf("Memdclient %s Authenticated successfully", client.LoggerID())
	}

	if selectCh != nil {
		selectErr := <-selectCh
		if selectErr != nil {
			logDebugf("Memdclient %s Failed to perform select bucket against server (%v)", client.LoggerID(), selectErr)
			return selectErr
		}
	}

	if configCh != nil {
		configResp := <-configCh
		err = configResp.Err
		if err == nil {
			// We don't want this to block us completing bootstrap.
			go mcc.cfgManager.OnNewConfig(configResp.Config)
		} else {
			logDebugf("Memdclient %s Failed to perform config fetch against server (%v)", client.LoggerID(), err)
			if errors.Is(err, ErrDocumentNotFound) {
				logDebugf("Memdclient %s detected that CCCP is unsupported, informing upstream", client.LoggerID())
				mcc.sendErrorToCCCPUnsupportedHandlers()
			}
		}
	}

	client.Features(helloResp.SrvFeatures)

	logDebugf("Memdclient %s Client Features: %+v", client.LoggerID(), features)
	logDebugf("Memdclient %s Server Features: %+v", client.LoggerID(), helloResp.SrvFeatures)

	return nil
}

func (mcc *memdClientDialerComponent) continueAfterAuth(client bootstrapClient, bucketName string, continueAuthCh chan bool,
	deadline time.Time) (chan error, chan getConfigResponse) {

	var selectCh chan error
	if bucketName != "" {
		selectCh = make(chan error, 1)
	}

	var configCh chan getConfigResponse
	if atomic.LoadUint32(&mcc.configApplied) == 0 {
		configCh = make(chan getConfigResponse, 1)
	}

	go func() {
		success := <-continueAuthCh
		if !success {
			if selectCh != nil {
				close(selectCh)
			}
			if configCh != nil {
				close(configCh)
			}
			return
		}
		var execCh chan error
		if selectCh != nil {
			var err error
			execCh, err = client.ExecSelectBucket([]byte(bucketName), deadline)
			if err != nil {
				logDebugf("Memdclient %s Failed to execute select bucket (%v)", client.LoggerID(), err)
				selectCh <- err
				return
			}
		}

		var execConfigCh chan getConfigResponse
		if configCh != nil {
			var err error
			execConfigCh, err = client.ExecGetConfig(deadline)
			if err != nil {
				// Getting a config isn't essential to bootstrap.
				logDebugf("Memdclient %s Failed to execute get config (%v)", client.LoggerID(), err)
				close(configCh)
				return
			}
		}

		if selectCh != nil {
			execErr := <-execCh
			selectCh <- execErr
		}

		if configCh != nil {
			configResp := <-execConfigCh
			configCh <- configResp
		}
	}()

	return selectCh, configCh
}

type authFunc func() (continueCh chan error, completedCb chan bool, err error)

func (mcc *memdClientDialerComponent) buildAuthHandler(client bootstrapClient, auth AuthProvider, deadline time.Time,
	mechanism AuthMechanism) authFunc {
	creds, err := getKvAuthCreds(auth, client.Address())
	if err != nil {
		return nil
	}

	if creds.Username != "" || creds.Password != "" {
		return func() (chan error, chan bool, error) {
			continueCh := make(chan bool, 1)
			completedCh := make(chan error, 1)
			hasContinued := int32(0)
			callErr := saslMethod(mechanism, creds.Username, creds.Password, client, deadline, func() {
				// hasContinued should never be 1 here but let's guard against it.
				if atomic.CompareAndSwapInt32(&hasContinued, 0, 1) {
					continueCh <- true
				}
			}, func(err error) {
				if atomic.CompareAndSwapInt32(&hasContinued, 0, 1) {
					sendContinue := true
					if err != nil {
						sendContinue = false
					}
					continueCh <- sendContinue
				}
				completedCh <- err
			})
			if callErr != nil {
				return nil, nil, callErr
			}
			return completedCh, continueCh, nil
		}
	}

	return nil
}

func (mcc *memdClientDialerComponent) sendErrorToCCCPUnsupportedHandlers() {
	mcc.cccpUnsupportedHandlersLock.Lock()
	handlers := make([]memdBoostrapCCCPUnsupportedHandler, len(mcc.cccpUnsupportedFailHandlers))
	copy(handlers, mcc.cccpUnsupportedFailHandlers)
	mcc.cccpUnsupportedHandlersLock.Unlock()
	for _, h := range handlers {
		h.onCCCPUnsupported(ErrUnsupportedOperation)
	}
}

func checkSupportsFeature(srvFeatures []memd.HelloFeature, feature memd.HelloFeature) bool {
	for _, srvFeature := range srvFeatures {
		if srvFeature == feature {
			return true
		}
	}
	return false
}

func findNextAuthMechanism(authMechanisms []AuthMechanism, serverAuthMechanisms []AuthMechanism) (bool, AuthMechanism, []AuthMechanism) {
	for {
		if len(authMechanisms) <= 1 {
			break
		}
		authMechanisms = authMechanisms[1:]
		mech := authMechanisms[0]
		for _, serverMech := range serverAuthMechanisms {
			if mech == serverMech {
				return true, mech, authMechanisms
			}
		}
	}

	return false, "", authMechanisms
}

func helloFeatures(props helloProps) []memd.HelloFeature {
	var features []memd.HelloFeature

	// Send the TLS flag, which has unknown effects.
	features = append(features, memd.FeatureTLS)

	// Indicate that we understand XATTRs
	features = append(features, memd.FeatureXattr)

	// Indicates that we understand select buckets.
	features = append(features, memd.FeatureSelectBucket)

	// Indicates that we understand nmv responses containing no config map.
	features = append(features, memd.FeatureDedupeNotMyVbucketClustermap)

	// Indicates that we understand known version cluster map requests.
	features = append(features, memd.FeatureClusterMapKnownVersion)

	// Indicates that we understand duplex communication.
	features = append(features, memd.FeatureDuplex)

	// If the user wants to use KV Error maps, lets enable them
	if props.XErrorFeatureEnabled {
		features = append(features, memd.FeatureXerror)
	}

	// Indicate that we understand JSON
	if props.JSONFeatureEnabled {
		features = append(features, memd.FeatureJSON)
	}

	// Indicate that we understand Point in Time
	if props.PITRFeatureEnabled {
		features = append(features, memd.FeaturePITR)
	}

	// If the user wants to use mutation tokens, lets enable them
	if props.MutationTokensEnabled {
		features = append(features, memd.FeatureSeqNo)
	}

	// If the user wants on-the-wire compression, lets try to enable it
	if props.CompressionEnabled {
		features = append(features, memd.FeatureSnappy)
		features = append(features, memd.FeatureSnappyEverywhere)
	}

	if props.DurationsEnabled {
		features = append(features, memd.FeatureDurations)
	}

	if props.CollectionsEnabled {
		features = append(features, memd.FeatureCollections)
	}

	if props.OutOfOrderEnabled {
		features = append(features, memd.FeatureUnorderedExec)
	}

	if props.ClusterMapNotificationsEnabled {
		features = append(features, memd.FeatureClustermapChangeNotificationBrief)
	}

	// These flags are informational so don't actually enable anything
	features = append(features, memd.FeatureAltRequests)
	features = append(features, memd.FeatureCreateAsDeleted)
	features = append(features, memd.FeatureReplaceBodyWithXattr)
	features = append(features, memd.FeaturePreserveExpiry)
	features = append(features, memd.FeatureSubdocReplicaRead)

	if props.SyncReplicationEnabled {
		features = append(features, memd.FeatureSyncReplication)
	}

	if props.ResourceUnitsEnabled {
		features = append(features, memd.FeatureResourceUnits)
	}

	return features
}
