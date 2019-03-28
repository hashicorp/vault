package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/subtle"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/vault/helper/metricsutil"

	metrics "github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	multierror "github.com/hashicorp/go-multierror"
	uuid "github.com/hashicorp/go-uuid"
	cache "github.com/patrickmn/go-cache"

	"google.golang.org/grpc"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/certutil"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
	"github.com/hashicorp/vault/vault/seal"
)

const (
	// CoreLockPath is the path used to acquire a coordinating lock
	// for a highly-available deploy.
	CoreLockPath = "core/lock"

	// The poison pill is used as a check during certain scenarios to indicate
	// to standby nodes that they should seal
	poisonPillPath = "core/poison-pill"

	// coreLeaderPrefix is the prefix used for the UUID that contains
	// the currently elected leader.
	coreLeaderPrefix = "core/leader/"

	// knownPrimaryAddrsPrefix is used to store last-known cluster address
	// information for primaries
	knownPrimaryAddrsPrefix = "core/primary-addrs/"

	// coreKeyringCanaryPath is used as a canary to indicate to replicated
	// clusters that they need to perform a rekey operation synchronously; this
	// isn't keyring-canary to avoid ignoring it when ignoring core/keyring
	coreKeyringCanaryPath = "core/canary-keyring"
)

var (
	// ErrAlreadyInit is returned if the core is already
	// initialized. This prevents a re-initialization.
	ErrAlreadyInit = errors.New("Vault is already initialized")

	// ErrNotInit is returned if a non-initialized barrier
	// is attempted to be unsealed.
	ErrNotInit = errors.New("Vault is not initialized")

	// ErrInternalError is returned when we don't want to leak
	// any information about an internal error
	ErrInternalError = errors.New("internal error")

	// ErrHANotEnabled is returned if the operation only makes sense
	// in an HA setting
	ErrHANotEnabled = errors.New("Vault is not configured for highly-available mode")

	// manualStepDownSleepPeriod is how long to sleep after a user-initiated
	// step down of the active node, to prevent instantly regrabbing the lock.
	// It's var not const so that tests can manipulate it.
	manualStepDownSleepPeriod = 10 * time.Second

	// Functions only in the Enterprise version
	enterprisePostUnseal = enterprisePostUnsealImpl
	enterprisePreSeal    = enterprisePreSealImpl
	startReplication     = startReplicationImpl
	stopReplication      = stopReplicationImpl
	LastWAL              = lastWALImpl
	LastRemoteWAL        = lastRemoteWALImpl
	WaitUntilWALShipped  = waitUntilWALShippedImpl
)

// NonFatalError is an error that can be returned during NewCore that should be
// displayed but not cause a program exit
type NonFatalError struct {
	Err error
}

func (e *NonFatalError) WrappedErrors() []error {
	return []error{e.Err}
}

func (e *NonFatalError) Error() string {
	return e.Err.Error()
}

// NewNonFatalError returns a new non-fatal error.
func NewNonFatalError(err error) *NonFatalError {
	return &NonFatalError{Err: err}
}

// IsFatalError returns true if the given error is a non-fatal error.
func IsFatalError(err error) bool {
	return !errwrap.ContainsType(err, new(NonFatalError))
}

// ErrInvalidKey is returned if there is a user-based error with a provided
// unseal key. This will be shown to the user, so should not contain
// information that is sensitive.
type ErrInvalidKey struct {
	Reason string
}

func (e *ErrInvalidKey) Error() string {
	return fmt.Sprintf("invalid key: %v", e.Reason)
}

type RegisterAuthFunc func(context.Context, time.Duration, string, *logical.Auth) error

type activeAdvertisement struct {
	RedirectAddr     string                     `json:"redirect_addr"`
	ClusterAddr      string                     `json:"cluster_addr,omitempty"`
	ClusterCert      []byte                     `json:"cluster_cert,omitempty"`
	ClusterKeyParams *certutil.ClusterKeyParams `json:"cluster_key_params,omitempty"`
}

type unlockInformation struct {
	Parts [][]byte
	Nonce string
}

// Core is used as the central manager of Vault activity. It is the primary point of
// interface for API handlers and is responsible for managing the logical and physical
// backends, router, security barrier, and audit trails.
type Core struct {
	entCore

	// The registry of builtin plugins is passed in here as an interface because
	// if it's used directly, it results in import cycles.
	builtinRegistry BuiltinRegistry

	// N.B.: This is used to populate a dev token down replication, as
	// otherwise, after replication is started, a dev would have to go through
	// the generate-root process simply to talk to the new follower cluster.
	devToken string

	// HABackend may be available depending on the physical backend
	ha physical.HABackend

	// redirectAddr is the address we advertise as leader if held
	redirectAddr string

	// clusterAddr is the address we use for clustering
	clusterAddr string

	// physical backend is the un-trusted backend with durable data
	physical physical.Backend

	// seal is our seal, for seal configuration information
	seal Seal

	// migrationSeal is the seal to use during a migration operation. It is the
	// seal we're migrating *from*.
	migrationSeal Seal

	// unwrapSeal is the seal to use on Enterprise to unwrap values wrapped
	// with the previous seal.
	unwrapSeal Seal

	// barrier is the security barrier wrapping the physical backend
	barrier SecurityBarrier

	// router is responsible for managing the mount points for logical backends.
	router *Router

	// logicalBackends is the mapping of backends to use for this core
	logicalBackends map[string]logical.Factory

	// credentialBackends is the mapping of backends to use for this core
	credentialBackends map[string]logical.Factory

	// auditBackends is the mapping of backends to use for this core
	auditBackends map[string]audit.Factory

	// stateLock protects mutable state
	stateLock sync.RWMutex
	sealed    *uint32

	standby              bool
	perfStandby          bool
	standbyDoneCh        chan struct{}
	standbyStopCh        chan struct{}
	manualStepDownCh     chan struct{}
	keepHALockOnStepDown *uint32
	heldHALock           physical.Lock

	// unlockInfo has the keys provided to Unseal until the threshold number of parts is available, as well as the operation nonce
	unlockInfo *unlockInformation

	// generateRootProgress holds the shares until we reach enough
	// to verify the master key
	generateRootConfig   *GenerateRootConfig
	generateRootProgress [][]byte
	generateRootLock     sync.Mutex

	// These variables holds the config and shares we have until we reach
	// enough to verify the appropriate master key. Note that the same lock is
	// used; this isn't time-critical so this shouldn't be a problem.
	barrierRekeyConfig  *SealConfig
	recoveryRekeyConfig *SealConfig
	rekeyLock           sync.RWMutex

	// mounts is loaded after unseal since it is a protected
	// configuration
	mounts *MountTable

	// mountsLock is used to ensure that the mounts table does not
	// change underneath a calling function
	mountsLock sync.RWMutex

	// auth is loaded after unseal since it is a protected
	// configuration
	auth *MountTable

	// authLock is used to ensure that the auth table does not
	// change underneath a calling function
	authLock sync.RWMutex

	// audit is loaded after unseal since it is a protected
	// configuration
	audit *MountTable

	// auditLock is used to ensure that the audit table does not
	// change underneath a calling function
	auditLock sync.RWMutex

	// auditBroker is used to ingest the audit events and fan
	// out into the configured audit backends
	auditBroker *AuditBroker

	// auditedHeaders is used to configure which http headers
	// can be output in the audit logs
	auditedHeaders *AuditedHeadersConfig

	// systemBackend is the backend which is used to manage internal operations
	systemBackend *SystemBackend

	// cubbyholeBackend is the backend which manages the per-token storage
	cubbyholeBackend *CubbyholeBackend

	// systemBarrierView is the barrier view for the system backend
	systemBarrierView *BarrierView

	// expiration manager is used for managing LeaseIDs,
	// renewal, expiration and revocation
	expiration *ExpirationManager

	// rollback manager is used to run rollbacks periodically
	rollback *RollbackManager

	// policy store is used to manage named ACL policies
	policyStore *PolicyStore

	// token store is used to manage authentication tokens
	tokenStore *TokenStore

	// identityStore is used to manage client entities
	identityStore *IdentityStore

	// metricsCh is used to stop the metrics streaming
	metricsCh chan struct{}

	// metricsMutex is used to prevent a race condition between
	// metrics emission and sealing leading to a nil pointer
	metricsMutex sync.Mutex

	defaultLeaseTTL time.Duration
	maxLeaseTTL     time.Duration

	// baseLogger is used to avoid ResetNamed as it strips useful prefixes in
	// e.g. testing
	baseLogger log.Logger
	logger     log.Logger

	// cachingDisabled indicates whether caches are disabled
	cachingDisabled bool
	// Cache stores the actual cache; we always have this but may bypass it if
	// disabled
	physicalCache physical.ToggleablePurgemonster

	// reloadFuncs is a map containing reload functions
	reloadFuncs map[string][]reload.ReloadFunc

	// reloadFuncsLock controls access to the funcs
	reloadFuncsLock sync.RWMutex

	// wrappingJWTKey is the key used for generating JWTs containing response
	// wrapping information
	wrappingJWTKey *ecdsa.PrivateKey

	//
	// Cluster information
	//
	// Name
	clusterName string
	// Specific cipher suites to use for clustering, if any
	clusterCipherSuites []uint16
	// Used to modify cluster parameters
	clusterParamsLock sync.RWMutex
	// The private key stored in the barrier used for establishing
	// mutually-authenticated connections between Vault cluster members
	localClusterPrivateKey *atomic.Value
	// The local cluster cert
	localClusterCert *atomic.Value
	// The parsed form of the local cluster cert
	localClusterParsedCert *atomic.Value
	// The TCP addresses we should use for clustering
	clusterListenerAddrs []*net.TCPAddr
	// The handler to use for request forwarding
	clusterHandler http.Handler
	// Write lock used to ensure that we don't have multiple connections adjust
	// this value at the same time
	requestForwardingConnectionLock sync.RWMutex
	// Lock for the leader values, ensuring we don't run the parts of Leader()
	// that change things concurrently
	leaderParamsLock sync.RWMutex
	// Current cluster leader values
	clusterLeaderParams *atomic.Value
	// Info on cluster members
	clusterPeerClusterAddrsCache *cache.Cache
	// The context for the client
	rpcClientConnContext context.Context
	// The function for canceling the client connection
	rpcClientConnCancelFunc context.CancelFunc
	// The grpc ClientConn for RPC calls
	rpcClientConn *grpc.ClientConn
	// The grpc forwarding client
	rpcForwardingClient *forwardingClient

	// CORS Information
	corsConfig *CORSConfig

	// The active set of upstream cluster addresses; stored via the Echo
	// mechanism, loaded by the balancer
	atomicPrimaryClusterAddrs *atomic.Value

	atomicPrimaryFailoverAddrs *atomic.Value

	// replicationState keeps the current replication state cached for quick
	// lookup; activeNodeReplicationState stores the active value on standbys
	replicationState           *uint32
	activeNodeReplicationState *uint32

	// uiConfig contains UI configuration
	uiConfig *UIConfig

	// rawEnabled indicates whether the Raw endpoint is enabled
	rawEnabled bool

	// pluginDirectory is the location vault will look for plugin binaries
	pluginDirectory string

	// pluginCatalog is used to manage plugin configurations
	pluginCatalog *PluginCatalog

	enableMlock bool

	// This can be used to trigger operations to stop running when Vault is
	// going to be shut down, stepped down, or sealed
	activeContext           context.Context
	activeContextCancelFunc *atomic.Value

	// Stores the sealunwrapper for downgrade needs
	sealUnwrapper physical.Backend

	// unsealwithStoredKeysLock is a mutex that prevents multiple processes from
	// unsealing with stored keys are the same time.
	unsealWithStoredKeysLock sync.Mutex

	// Stores any funcs that should be run on successful postUnseal
	postUnsealFuncs []func()

	// replicationFailure is used to mark when replication has entered an
	// unrecoverable failure.
	replicationFailure *uint32

	// disablePerfStanby is used to tell a standby not to attempt to become a
	// perf standby
	disablePerfStandby bool

	licensingStopCh chan struct{}

	// Stores loggers so we can reset the level
	allLoggers     []log.Logger
	allLoggersLock sync.RWMutex

	// Can be toggled atomically to cause the core to never try to become
	// active, or give up active as soon as it gets it
	neverBecomeActive *uint32

	// loadCaseSensitiveIdentityStore enforces the loading of identity store
	// artifacts in a case sensitive manner. To be used only in testing.
	loadCaseSensitiveIdentityStore bool

	// clusterListener starts up and manages connections on the cluster ports
	clusterListener *ClusterListener

	// Telemetry objects
	metricsHelper *metricsutil.MetricsHelper

	// Stores request counters
	counters counters
}

// CoreConfig is used to parameterize a core
type CoreConfig struct {
	DevToken string `json:"dev_token" structs:"dev_token" mapstructure:"dev_token"`

	BuiltinRegistry BuiltinRegistry `json:"builtin_registry" structs:"builtin_registry" mapstructure:"builtin_registry"`

	LogicalBackends map[string]logical.Factory `json:"logical_backends" structs:"logical_backends" mapstructure:"logical_backends"`

	CredentialBackends map[string]logical.Factory `json:"credential_backends" structs:"credential_backends" mapstructure:"credential_backends"`

	AuditBackends map[string]audit.Factory `json:"audit_backends" structs:"audit_backends" mapstructure:"audit_backends"`

	Physical physical.Backend `json:"physical" structs:"physical" mapstructure:"physical"`

	// May be nil, which disables HA operations
	HAPhysical physical.HABackend `json:"ha_physical" structs:"ha_physical" mapstructure:"ha_physical"`

	Seal Seal `json:"seal" structs:"seal" mapstructure:"seal"`

	Logger log.Logger `json:"logger" structs:"logger" mapstructure:"logger"`

	// Disables the LRU cache on the physical backend
	DisableCache bool `json:"disable_cache" structs:"disable_cache" mapstructure:"disable_cache"`

	// Disables mlock syscall
	DisableMlock bool `json:"disable_mlock" structs:"disable_mlock" mapstructure:"disable_mlock"`

	// Custom cache size for the LRU cache on the physical backend, or zero for default
	CacheSize int `json:"cache_size" structs:"cache_size" mapstructure:"cache_size"`

	// Set as the leader address for HA
	RedirectAddr string `json:"redirect_addr" structs:"redirect_addr" mapstructure:"redirect_addr"`

	// Set as the cluster address for HA
	ClusterAddr string `json:"cluster_addr" structs:"cluster_addr" mapstructure:"cluster_addr"`

	DefaultLeaseTTL time.Duration `json:"default_lease_ttl" structs:"default_lease_ttl" mapstructure:"default_lease_ttl"`

	MaxLeaseTTL time.Duration `json:"max_lease_ttl" structs:"max_lease_ttl" mapstructure:"max_lease_ttl"`

	ClusterName string `json:"cluster_name" structs:"cluster_name" mapstructure:"cluster_name"`

	ClusterCipherSuites string `json:"cluster_cipher_suites" structs:"cluster_cipher_suites" mapstructure:"cluster_cipher_suites"`

	EnableUI bool `json:"ui" structs:"ui" mapstructure:"ui"`

	// Enable the raw endpoint
	EnableRaw bool `json:"enable_raw" structs:"enable_raw" mapstructure:"enable_raw"`

	PluginDirectory string `json:"plugin_directory" structs:"plugin_directory" mapstructure:"plugin_directory"`

	DisableSealWrap bool `json:"disable_sealwrap" structs:"disable_sealwrap" mapstructure:"disable_sealwrap"`

	ReloadFuncs     *map[string][]reload.ReloadFunc
	ReloadFuncsLock *sync.RWMutex

	// Licensing
	LicensingConfig *LicensingConfig
	// Don't set this unless in dev mode, ideally only when using inmem
	DevLicenseDuration time.Duration

	DisablePerformanceStandby bool
	DisableIndexing           bool
	DisableKeyEncodingChecks  bool

	AllLoggers []log.Logger

	// Telemetry objects
	MetricsHelper *metricsutil.MetricsHelper

	CounterSyncInterval time.Duration
}

func (c *CoreConfig) Clone() *CoreConfig {
	return &CoreConfig{
		DevToken:                  c.DevToken,
		LogicalBackends:           c.LogicalBackends,
		CredentialBackends:        c.CredentialBackends,
		AuditBackends:             c.AuditBackends,
		Physical:                  c.Physical,
		HAPhysical:                c.HAPhysical,
		Seal:                      c.Seal,
		Logger:                    c.Logger,
		DisableCache:              c.DisableCache,
		DisableMlock:              c.DisableMlock,
		CacheSize:                 c.CacheSize,
		RedirectAddr:              c.RedirectAddr,
		ClusterAddr:               c.ClusterAddr,
		DefaultLeaseTTL:           c.DefaultLeaseTTL,
		MaxLeaseTTL:               c.MaxLeaseTTL,
		ClusterName:               c.ClusterName,
		ClusterCipherSuites:       c.ClusterCipherSuites,
		EnableUI:                  c.EnableUI,
		EnableRaw:                 c.EnableRaw,
		PluginDirectory:           c.PluginDirectory,
		DisableSealWrap:           c.DisableSealWrap,
		ReloadFuncs:               c.ReloadFuncs,
		ReloadFuncsLock:           c.ReloadFuncsLock,
		LicensingConfig:           c.LicensingConfig,
		DevLicenseDuration:        c.DevLicenseDuration,
		DisablePerformanceStandby: c.DisablePerformanceStandby,
		DisableIndexing:           c.DisableIndexing,
		AllLoggers:                c.AllLoggers,
		CounterSyncInterval:       c.CounterSyncInterval,
	}
}

// NewCore is used to construct a new core
func NewCore(conf *CoreConfig) (*Core, error) {
	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		if conf.RedirectAddr == "" {
			return nil, fmt.Errorf("missing API address, please set in configuration or via environment")
		}
	}

	if conf.DefaultLeaseTTL == 0 {
		conf.DefaultLeaseTTL = defaultLeaseTTL
	}
	if conf.MaxLeaseTTL == 0 {
		conf.MaxLeaseTTL = maxLeaseTTL
	}
	if conf.DefaultLeaseTTL > conf.MaxLeaseTTL {
		return nil, fmt.Errorf("cannot have DefaultLeaseTTL larger than MaxLeaseTTL")
	}

	// Validate the advertise addr if its given to us
	if conf.RedirectAddr != "" {
		u, err := url.Parse(conf.RedirectAddr)
		if err != nil {
			return nil, errwrap.Wrapf("redirect address is not valid url: {{err}}", err)
		}

		if u.Scheme == "" {
			return nil, fmt.Errorf("redirect address must include scheme (ex. 'http')")
		}
	}

	// Make a default logger if not provided
	if conf.Logger == nil {
		conf.Logger = logging.NewVaultLogger(log.Trace)
	}

	syncInterval := conf.CounterSyncInterval
	if syncInterval.Nanoseconds() == 0 {
		syncInterval = 30 * time.Second
	}

	// Setup the core
	c := &Core{
		entCore:                      entCore{},
		devToken:                     conf.DevToken,
		physical:                     conf.Physical,
		redirectAddr:                 conf.RedirectAddr,
		clusterAddr:                  conf.ClusterAddr,
		seal:                         conf.Seal,
		router:                       NewRouter(),
		sealed:                       new(uint32),
		standby:                      true,
		baseLogger:                   conf.Logger,
		logger:                       conf.Logger.Named("core"),
		defaultLeaseTTL:              conf.DefaultLeaseTTL,
		maxLeaseTTL:                  conf.MaxLeaseTTL,
		cachingDisabled:              conf.DisableCache,
		clusterName:                  conf.ClusterName,
		clusterPeerClusterAddrsCache: cache.New(3*HeartbeatInterval, time.Second),
		enableMlock:                  !conf.DisableMlock,
		rawEnabled:                   conf.EnableRaw,
		replicationState:             new(uint32),
		atomicPrimaryClusterAddrs:    new(atomic.Value),
		atomicPrimaryFailoverAddrs:   new(atomic.Value),
		localClusterPrivateKey:       new(atomic.Value),
		localClusterCert:             new(atomic.Value),
		localClusterParsedCert:       new(atomic.Value),
		activeNodeReplicationState:   new(uint32),
		keepHALockOnStepDown:         new(uint32),
		replicationFailure:           new(uint32),
		disablePerfStandby:           true,
		activeContextCancelFunc:      new(atomic.Value),
		allLoggers:                   conf.AllLoggers,
		builtinRegistry:              conf.BuiltinRegistry,
		neverBecomeActive:            new(uint32),
		clusterLeaderParams:          new(atomic.Value),
		metricsHelper:                conf.MetricsHelper,
		counters: counters{
			requests:     new(uint64),
			syncInterval: syncInterval,
		},
	}

	atomic.StoreUint32(c.sealed, 1)
	c.allLoggers = append(c.allLoggers, c.logger)

	atomic.StoreUint32(c.replicationState, uint32(consts.ReplicationDRDisabled|consts.ReplicationPerformanceDisabled))
	c.localClusterCert.Store(([]byte)(nil))
	c.localClusterParsedCert.Store((*x509.Certificate)(nil))
	c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

	c.clusterLeaderParams.Store((*ClusterLeaderParams)(nil))

	c.activeContextCancelFunc.Store((context.CancelFunc)(nil))

	if conf.ClusterCipherSuites != "" {
		suites, err := tlsutil.ParseCiphers(conf.ClusterCipherSuites)
		if err != nil {
			return nil, errwrap.Wrapf("error parsing cluster cipher suites: {{err}}", err)
		}
		c.clusterCipherSuites = suites
	}

	// Load CORS config and provide a value for the core field.
	c.corsConfig = &CORSConfig{
		core:    c,
		Enabled: new(uint32),
	}

	if c.seal == nil {
		c.seal = NewDefaultSeal()
	}
	c.seal.SetCore(c)

	if err := coreInit(c, conf); err != nil {
		return nil, err
	}

	if !conf.DisableMlock {
		// Ensure our memory usage is locked into physical RAM
		if err := mlock.LockMemory(); err != nil {
			return nil, fmt.Errorf(
				"Failed to lock memory: %v\n\n"+
					"This usually means that the mlock syscall is not available.\n"+
					"Vault uses mlock to prevent memory from being swapped to\n"+
					"disk. This requires root privileges as well as a machine\n"+
					"that supports mlock. Please enable mlock on your system or\n"+
					"disable Vault from using it. To disable Vault from using it,\n"+
					"set the `disable_mlock` configuration option in your configuration\n"+
					"file.",
				err)
		}
	}

	var err error

	if conf.PluginDirectory != "" {
		c.pluginDirectory, err = filepath.Abs(conf.PluginDirectory)
		if err != nil {
			return nil, errwrap.Wrapf("core setup failed, could not verify plugin directory: {{err}}", err)
		}
	}

	// Construct a new AES-GCM barrier
	c.barrier, err = NewAESGCMBarrier(c.physical)
	if err != nil {
		return nil, errwrap.Wrapf("barrier setup failed: {{err}}", err)
	}

	createSecondaries(c, conf)

	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		c.ha = conf.HAPhysical
	}

	// We create the funcs here, then populate the given config with it so that
	// the caller can share state
	conf.ReloadFuncsLock = &c.reloadFuncsLock
	c.reloadFuncsLock.Lock()
	c.reloadFuncs = make(map[string][]reload.ReloadFunc)
	c.reloadFuncsLock.Unlock()
	conf.ReloadFuncs = &c.reloadFuncs

	logicalBackends := make(map[string]logical.Factory)
	for k, f := range conf.LogicalBackends {
		logicalBackends[k] = f
	}
	_, ok := logicalBackends["kv"]
	if !ok {
		logicalBackends["kv"] = PassthroughBackendFactory
	}

	logicalBackends["cubbyhole"] = CubbyholeBackendFactory
	logicalBackends[systemMountType] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		sysBackendLogger := conf.Logger.Named("system")
		c.AddLogger(sysBackendLogger)
		b := NewSystemBackend(c, sysBackendLogger)
		if err := b.Setup(ctx, config); err != nil {
			return nil, err
		}
		return b, nil
	}
	logicalBackends["identity"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		identityLogger := conf.Logger.Named("identity")
		c.AddLogger(identityLogger)
		return NewIdentityStore(ctx, c, config, identityLogger)
	}
	addExtraLogicalBackends(c, logicalBackends)
	c.logicalBackends = logicalBackends

	credentialBackends := make(map[string]logical.Factory)
	for k, f := range conf.CredentialBackends {
		credentialBackends[k] = f
	}
	credentialBackends["token"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		tsLogger := conf.Logger.Named("token")
		c.AddLogger(tsLogger)
		return NewTokenStore(ctx, tsLogger, c, config)
	}
	addExtraCredentialBackends(c, credentialBackends)
	c.credentialBackends = credentialBackends

	auditBackends := make(map[string]audit.Factory)
	for k, f := range conf.AuditBackends {
		auditBackends[k] = f
	}
	c.auditBackends = auditBackends

	uiStoragePrefix := systemBarrierPrefix + "ui"
	c.uiConfig = NewUIConfig(conf.EnableUI, physical.NewView(c.physical, uiStoragePrefix), NewBarrierView(c.barrier, uiStoragePrefix))

	return c, nil
}

// Shutdown is invoked when the Vault instance is about to be terminated. It
// should not be accessible as part of an API call as it will cause an availability
// problem. It is only used to gracefully quit in the case of HA so that failover
// happens as quickly as possible.
func (c *Core) Shutdown() error {
	c.logger.Debug("shutdown called")
	return c.sealInternal()
}

// CORSConfig returns the current CORS configuration
func (c *Core) CORSConfig() *CORSConfig {
	return c.corsConfig
}

func (c *Core) GetContext() (context.Context, context.CancelFunc) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	return context.WithCancel(namespace.RootContext(c.activeContext))
}

// Sealed checks if the Vault is current sealed
func (c *Core) Sealed() bool {
	return atomic.LoadUint32(c.sealed) == 1
}

// SecretProgress returns the number of keys provided so far
func (c *Core) SecretProgress() (int, string) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	switch c.unlockInfo {
	case nil:
		return 0, ""
	default:
		return len(c.unlockInfo.Parts), c.unlockInfo.Nonce
	}
}

// ResetUnsealProcess removes the current unlock parts from memory, to reset
// the unsealing process
func (c *Core) ResetUnsealProcess() {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	c.unlockInfo = nil
}

// Unseal is used to provide one of the key parts to unseal the Vault.
//
// They key given as a parameter will automatically be zerod after
// this method is done with it. If you want to keep the key around, a copy
// should be made.
func (c *Core) Unseal(key []byte) (bool, error) {
	return c.unseal(key, false)
}

func (c *Core) UnsealWithRecoveryKeys(key []byte) (bool, error) {
	return c.unseal(key, true)
}

func (c *Core) unseal(key []byte, useRecoveryKeys bool) (bool, error) {
	defer metrics.MeasureSince([]string{"core", "unseal"}, time.Now())

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	ctx := context.Background()

	// Explicitly check for init status. This also checks if the seal
	// configuration is valid (i.e. non-nil).
	init, err := c.Initialized(ctx)
	if err != nil {
		return false, err
	}
	if !init {
		return false, ErrNotInit
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return false, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return false, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Check if already unsealed
	if !c.Sealed() {
		return true, nil
	}

	sealToUse := c.seal
	if c.migrationSeal != nil {
		sealToUse = c.migrationSeal
	}

	masterKey, err := c.unsealPart(ctx, sealToUse, key, useRecoveryKeys)
	if err != nil {
		return false, err
	}
	if masterKey != nil {
		return c.unsealInternal(ctx, masterKey)
	}

	return false, nil
}

// unsealPart takes in a key share, and returns the master key if the threshold
// is met. If recovery keys are supported, recovery key shares may be provided.
func (c *Core) unsealPart(ctx context.Context, seal Seal, key []byte, useRecoveryKeys bool) ([]byte, error) {
	// Check if we already have this piece
	if c.unlockInfo != nil {
		for _, existing := range c.unlockInfo.Parts {
			if subtle.ConstantTimeCompare(existing, key) == 1 {
				return nil, nil
			}
		}
	} else {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		c.unlockInfo = &unlockInformation{
			Nonce: uuid,
		}
	}

	// Store this key
	c.unlockInfo.Parts = append(c.unlockInfo.Parts, key)

	var config *SealConfig
	var err error
	if seal.RecoveryKeySupported() && (useRecoveryKeys || c.migrationSeal != nil) {
		config, err = seal.RecoveryConfig(ctx)
	} else {
		config, err = seal.BarrierConfig(ctx)
	}
	if err != nil {
		return nil, err
	}

	// Check if we don't have enough keys to unlock, proceed through the rest of
	// the call only if we have met the threshold
	if len(c.unlockInfo.Parts) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot unseal, not enough keys", "keys", len(c.unlockInfo.Parts), "threshold", config.SecretThreshold, "nonce", c.unlockInfo.Nonce)
		}
		return nil, nil
	}

	// Best-effort memzero of unlock parts once we're done with them
	defer func() {
		for i := range c.unlockInfo.Parts {
			memzero(c.unlockInfo.Parts[i])
		}
		c.unlockInfo = nil
	}()

	// Recover the split key. recoveredKey is the shamir combined
	// key, or the single provided key if the threshold is 1.
	var recoveredKey []byte
	var masterKey []byte
	var recoveryKey []byte
	if config.SecretThreshold == 1 {
		recoveredKey = make([]byte, len(c.unlockInfo.Parts[0]))
		copy(recoveredKey, c.unlockInfo.Parts[0])
	} else {
		recoveredKey, err = shamir.Combine(c.unlockInfo.Parts)
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
		}
	}

	if seal.RecoveryKeySupported() && (useRecoveryKeys || c.migrationSeal != nil) {
		// Verify recovery key
		if err := seal.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			return nil, err
		}
		recoveryKey = recoveredKey

		// Get stored keys and shamir combine into single master key. Unsealing with
		// recovery keys currently does not support: 1) mixed stored and non-stored
		// keys setup, nor 2) seals that support recovery keys but not stored keys.
		// If insufficient shares are provided, shamir.Combine will error, and if
		// no stored keys are found it will return masterKey as nil.
		if seal.StoredKeysSupported() {
			masterKeyShares, err := seal.GetStoredKeys(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("unable to retrieve stored keys: {{err}}", err)
			}

			switch len(masterKeyShares) {
			case 0:
				return nil, errors.New("seal returned no master key shares")
			case 1:
				masterKey = masterKeyShares[0]
			default:
				masterKey, err = shamir.Combine(masterKeyShares)
				if err != nil {
					return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
				}
			}
		}
	} else {
		masterKey = recoveredKey
	}

	// If we have a migration seal, now's the time!
	if c.migrationSeal != nil {
		// Unseal the barrier so we can rekey
		if err := c.barrier.Unseal(ctx, masterKey); err != nil {
			return nil, errwrap.Wrapf("error unsealing barrier with constructed master key: {{err}}", err)
		}
		defer c.barrier.Seal()

		switch {
		case c.migrationSeal.RecoveryKeySupported() && c.seal.RecoveryKeySupported():
			// Set the recovery and barrier keys to be the same.
			recoveryKey, err := c.migrationSeal.RecoveryKey(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("error getting recovery key to set on new seal: {{err}}", err)
			}

			if err := c.seal.SetRecoveryKey(ctx, recoveryKey); err != nil {
				return nil, errwrap.Wrapf("error setting new recovery key information during migrate: {{err}}", err)
			}

			barrierKeys, err := c.migrationSeal.GetStoredKeys(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("error getting stored keys to set on new seal: {{err}}", err)
			}

			if err := c.seal.SetStoredKeys(ctx, barrierKeys); err != nil {
				return nil, errwrap.Wrapf("error setting new barrier key information during migrate: {{err}}", err)
			}

		case c.migrationSeal.RecoveryKeySupported():
			// Auto to Shamir, since recovery key isn't supported on new seal

			// In this case we have to ensure that the recovery information was
			// set properly.
			if recoveryKey == nil {
				return nil, errors.New("did not get expected recovery information to set new seal during migration")
			}

			// We have recovery keys; we're going to use them as the new
			// barrier key.
			if err := c.barrier.Rekey(ctx, recoveryKey); err != nil {
				return nil, errwrap.Wrapf("error rekeying barrier during migration: {{err}}", err)
			}

			if err := c.barrier.Delete(ctx, StoredBarrierKeysPath); err != nil {
				// Don't actually exit here as successful deletion isn't critical
				c.logger.Error("error deleting stored barrier keys after migration; continuing anyways", "error", err)
			}

			masterKey = recoveryKey

		case c.seal.RecoveryKeySupported():
			// The new seal will have recovery keys; we set it to the existing
			// master key, so barrier key shares -> recovery key shares
			if err := c.seal.SetRecoveryKey(ctx, masterKey); err != nil {
				return nil, errwrap.Wrapf("error setting new recovery key information: {{err}}", err)
			}

			// Generate a new master key
			newMasterKey, err := c.barrier.GenerateKey()
			if err != nil {
				return nil, errwrap.Wrapf("error generating new master key: {{err}}", err)
			}

			// Rekey the barrier
			if err := c.barrier.Rekey(ctx, newMasterKey); err != nil {
				return nil, errwrap.Wrapf("error rekeying barrier during migration: {{err}}", err)
			}

			// Store the new master key
			if err := c.seal.SetStoredKeys(ctx, [][]byte{newMasterKey}); err != nil {
				return nil, errwrap.Wrapf("error storing new master key: {[err}}", err)
			}

			// Return the new key so it can be used to unlock the barrier
			masterKey = newMasterKey

		default:
			return nil, errors.New("unhandled migration case (shamir to shamir)")
		}

		// At this point we've swapped things around and need to ensure we
		// don't migrate again
		c.migrationSeal = nil

		// Ensure we populate the new values
		bc, err := c.seal.BarrierConfig(ctx)
		if err != nil {
			return nil, errwrap.Wrapf("error fetching barrier config after migration: {{err}}", err)
		}
		if err := c.seal.SetBarrierConfig(ctx, bc); err != nil {
			return nil, errwrap.Wrapf("error storing barrier config after migration: {{err}}", err)
		}

		if c.seal.RecoveryKeySupported() {
			rc, err := c.seal.RecoveryConfig(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("error fetching recovery config after migration: {{err}}", err)
			}
			if err := c.seal.SetRecoveryConfig(ctx, rc); err != nil {
				return nil, errwrap.Wrapf("error storing recovery config after migration: {{err}}", err)
			}
		}
	}

	return masterKey, nil
}

// unsealInternal takes in the master key and attempts to unseal the barrier.
// N.B.: This must be called with the state write lock held.
func (c *Core) unsealInternal(ctx context.Context, masterKey []byte) (bool, error) {
	defer memzero(masterKey)

	// Attempt to unlock
	if err := c.barrier.Unseal(ctx, masterKey); err != nil {
		return false, err
	}
	if c.logger.IsInfo() {
		c.logger.Info("vault is unsealed")
	}

	if err := preUnsealInternal(ctx, c); err != nil {
		return false, err
	}

	if err := c.startClusterListener(ctx); err != nil {
		return false, err
	}

	// Do post-unseal setup if HA is not enabled
	if c.ha == nil {
		// We still need to set up cluster info even if it's not part of a
		// cluster right now. This also populates the cached cluster object.
		if err := c.setupCluster(ctx); err != nil {
			c.logger.Error("cluster setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			return false, err
		}

		ctx, ctxCancel := context.WithCancel(namespace.RootContext(nil))
		if err := c.postUnseal(ctx, ctxCancel, standardUnsealStrategy{}); err != nil {
			c.logger.Error("post-unseal setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			return false, err
		}

		c.standby = false
	} else {
		// Go to standby mode, wait until we are active to unseal
		c.standbyDoneCh = make(chan struct{})
		c.manualStepDownCh = make(chan struct{})
		c.standbyStopCh = make(chan struct{})
		go c.runStandby(c.standbyDoneCh, c.manualStepDownCh, c.standbyStopCh)
	}

	// Force a cache bust here, which will also run migration code
	if c.seal.RecoveryKeySupported() {
		c.seal.SetRecoveryConfig(ctx, nil)
	}

	// Success!
	atomic.StoreUint32(c.sealed, 0)

	if c.ha != nil {
		sd, ok := c.ha.(physical.ServiceDiscovery)
		if ok {
			if err := sd.NotifySealedStateChange(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn("failed to notify unsealed status", "error", err)
				}
			}
		}
	}
	return true, nil
}

// SealWithRequest takes in a logical.Request, acquires the lock, and passes
// through to sealInternal
func (c *Core) SealWithRequest(httpCtx context.Context, req *logical.Request) error {
	defer metrics.MeasureSince([]string{"core", "seal-with-request"}, time.Now())

	if c.Sealed() {
		return nil
	}

	c.stateLock.RLock()

	// This will unlock the read lock
	// We use background context since we may not be active
	ctx, cancel := context.WithCancel(namespace.RootContext(nil))
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-httpCtx.Done():
			cancel()
		}
	}()

	// This will unlock the read lock
	return c.sealInitCommon(ctx, req)
}

// Seal takes in a token and creates a logical.Request, acquires the lock, and
// passes through to sealInternal
func (c *Core) Seal(token string) error {
	defer metrics.MeasureSince([]string{"core", "seal"}, time.Now())

	if c.Sealed() {
		return nil
	}

	c.stateLock.RLock()

	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/seal",
		ClientToken: token,
	}

	// This will unlock the read lock
	// We use background context since we may not be active
	return c.sealInitCommon(namespace.RootContext(nil), req)
}

// sealInitCommon is common logic for Seal and SealWithRequest and is used to
// re-seal the Vault. This requires the Vault to be unsealed again to perform
// any further operations. Note: this function will read-unlock the state lock.
func (c *Core) sealInitCommon(ctx context.Context, req *logical.Request) (retErr error) {
	defer metrics.MeasureSince([]string{"core", "seal-internal"}, time.Now())

	if req == nil {
		retErr = multierror.Append(retErr, errors.New("nil request to seal"))
		c.stateLock.RUnlock()
		return retErr
	}

	// Since there is no token store in standby nodes, sealing cannot be done.
	// Ideally, the request has to be forwarded to leader node for validation
	// and the operation should be performed. But for now, just returning with
	// an error and recommending a vault restart, which essentially does the
	// same thing.
	if c.standby {
		c.logger.Error("vault cannot seal when in standby mode; please restart instead")
		retErr = multierror.Append(retErr, errors.New("vault cannot seal when in standby mode; please restart instead"))
		c.stateLock.RUnlock()
		return retErr
	}

	acl, te, entity, identityPolicies, err := c.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		c.stateLock.RUnlock()
		return retErr
	}

	// Audit-log the request before going any further
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Accessor:    req.ClientTokenAccessor,
	}
	if te != nil {
		auth.IdentityPolicies = identityPolicies[te.NamespaceID]
		delete(identityPolicies, te.NamespaceID)
		auth.ExternalNamespacePolicies = identityPolicies
		auth.TokenPolicies = te.Policies
		auth.Policies = append(te.Policies, identityPolicies[te.NamespaceID]...)
		auth.Metadata = te.Meta
		auth.DisplayName = te.DisplayName
		auth.EntityID = te.EntityID
		auth.TokenType = te.Type
	}

	logInput := &audit.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
		c.logger.Error("failed to audit request", "request_path", req.Path, "error", err)
		retErr = multierror.Append(retErr, errors.New("failed to audit request, cannot continue"))
		c.stateLock.RUnlock()
		return retErr
	}

	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		c.stateLock.RUnlock()
		return retErr
	}
	if te != nil && te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		c.stateLock.RUnlock()
		return retErr
	}

	// Attempt to use the token (decrement num_uses)
	// On error bail out; if the token has been revoked, bail out too
	if te != nil {
		te, err = c.tokenStore.UseToken(ctx, te)
		if err != nil {
			c.logger.Error("failed to use token", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			c.stateLock.RUnlock()
			return retErr
		}
		if te == nil {
			// Token is no longer valid
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			c.stateLock.RUnlock()
			return retErr
		}
	}

	// Verify that this operation is allowed
	authResults := c.performPolicyChecks(ctx, acl, te, req, entity, &PolicyCheckOpts{
		RootPrivsRequired: true,
	})
	if !authResults.Allowed {
		c.stateLock.RUnlock()
		retErr = multierror.Append(retErr, authResults.Error)
		if authResults.Error.ErrorOrNil() == nil || authResults.DeniedError {
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		}
		return retErr
	}

	if te != nil && te.NumUses == tokenRevocationPending {
		// Token needs to be revoked. We do this immediately here because
		// we won't have a token store after sealing.
		leaseID, err := c.expiration.CreateOrFetchRevocationLeaseByToken(c.activeContext, te)
		if err == nil {
			err = c.expiration.Revoke(c.activeContext, leaseID)
		}
		if err != nil {
			c.logger.Error("token needed revocation before seal but failed to revoke", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
		}
	}

	// Unlock; sealing will grab the lock when needed
	c.stateLock.RUnlock()

	sealErr := c.sealInternal()

	if sealErr != nil {
		retErr = multierror.Append(retErr, sealErr)
	}

	return
}

// UIEnabled returns if the UI is enabled
func (c *Core) UIEnabled() bool {
	return c.uiConfig.Enabled()
}

// UIHeaders returns configured UI headers
func (c *Core) UIHeaders() (http.Header, error) {
	return c.uiConfig.Headers(context.Background())
}

// sealInternal is an internal method used to seal the vault.  It does not do
// any authorization checking.
func (c *Core) sealInternal() error {
	return c.sealInternalWithOptions(true, false)
}

func (c *Core) sealInternalWithOptions(grabStateLock, keepHALock bool) error {
	// Mark sealed, and if already marked return
	if swapped := atomic.CompareAndSwapUint32(c.sealed, 0, 1); !swapped {
		return nil
	}

	c.logger.Info("marked as sealed")

	// Clear forwarding clients
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	activeCtxCancel := c.activeContextCancelFunc.Load().(context.CancelFunc)
	cancelCtxAndLock := func() {
		doneCh := make(chan struct{})
		go func() {
			select {
			case <-doneCh:
			// Attempt to drain any inflight requests
			case <-time.After(DefaultMaxRequestDuration):
				if activeCtxCancel != nil {
					activeCtxCancel()
				}
			}
		}()

		c.stateLock.Lock()
		close(doneCh)
		// Stop requests from processing
		if activeCtxCancel != nil {
			activeCtxCancel()
		}
	}

	// Do pre-seal teardown if HA is not enabled
	if c.ha == nil {
		if grabStateLock {
			cancelCtxAndLock()
			defer c.stateLock.Unlock()
		}
		// Even in a non-HA context we key off of this for some things
		c.standby = true

		// Stop requests from processing
		if activeCtxCancel != nil {
			activeCtxCancel()
		}

		if err := c.preSeal(); err != nil {
			c.logger.Error("pre-seal teardown failed", "error", err)
			return fmt.Errorf("internal error")
		}
	} else {
		// If we are keeping the lock we already have the state write lock
		// held. Otherwise grab it here so that when stopCh is triggered we are
		// locked.
		if keepHALock {
			atomic.StoreUint32(c.keepHALockOnStepDown, 1)
		}
		if grabStateLock {
			cancelCtxAndLock()
			defer c.stateLock.Unlock()
		}

		// If we are trying to acquire the lock, force it to return with nil so
		// runStandby will exit
		// If we are active, signal the standby goroutine to shut down and wait
		// for completion. We have the state lock here so nothing else should
		// be toggling standby status.
		close(c.standbyStopCh)
		c.logger.Debug("finished triggering standbyStopCh for runStandby")

		// Wait for runStandby to stop
		<-c.standbyDoneCh
		atomic.StoreUint32(c.keepHALockOnStepDown, 0)
		c.logger.Debug("runStandby done")
	}

	// Stop the cluster listener
	c.stopClusterListener()

	c.logger.Debug("sealing barrier")
	if err := c.barrier.Seal(); err != nil {
		c.logger.Error("error sealing barrier", "error", err)
		return err
	}

	if c.ha != nil {
		sd, ok := c.ha.(physical.ServiceDiscovery)
		if ok {
			if err := sd.NotifySealedStateChange(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn("failed to notify sealed status", "error", err)
				}
			}
		}
	}

	postSealInternal(c)

	c.logger.Info("vault is sealed")

	return nil
}

type UnsealStrategy interface {
	unseal(context.Context, log.Logger, *Core) error
}

type standardUnsealStrategy struct{}

func (s standardUnsealStrategy) unseal(ctx context.Context, logger log.Logger, c *Core) error {
	// Clear forwarding clients; we're active
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	if err := postUnsealPhysical(c); err != nil {
		return err
	}

	if err := enterprisePostUnseal(c); err != nil {
		return err
	}

	if !c.IsDRSecondary() {
		if err := c.ensureWrappingKey(ctx); err != nil {
			return err
		}
	}
	if err := c.setupPluginCatalog(ctx); err != nil {
		return err
	}
	if err := c.loadMounts(ctx); err != nil {
		return err
	}
	if err := c.setupMounts(ctx); err != nil {
		return err
	}
	if err := c.setupPolicyStore(ctx); err != nil {
		return err
	}
	if err := c.loadCORSConfig(ctx); err != nil {
		return err
	}
	if err := c.loadCurrentRequestCounters(ctx, time.Now()); err != nil {
		return err
	}
	if err := c.loadCredentials(ctx); err != nil {
		return err
	}
	if err := c.setupCredentials(ctx); err != nil {
		return err
	}
	if !c.IsDRSecondary() {
		if err := c.startRollback(); err != nil {
			return err
		}
		if err := c.setupExpiration(expireLeaseStrategyRevoke); err != nil {
			return err
		}
		if err := c.loadAudits(ctx); err != nil {
			return err
		}
		if err := c.setupAudits(ctx); err != nil {
			return err
		}
		if err := c.loadIdentityStoreArtifacts(ctx); err != nil {
			return err
		}
		if err := loadMFAConfigs(ctx, c); err != nil {
			return err
		}
		if err := c.setupAuditedHeadersConfig(ctx); err != nil {
			return err
		}
	} else {
		c.auditBroker = NewAuditBroker(c.logger)
	}

	if c.clusterListener != nil && (c.ha != nil || shouldStartClusterListener(c)) {
		if err := c.startForwarding(ctx); err != nil {
			return err
		}
	}

	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()
	if err := startReplication(c); err != nil {
		return err
	}

	return nil
}

// postUnseal is invoked after the barrier is unsealed, but before
// allowing any user operations. This allows us to setup any state that
// requires the Vault to be unsealed such as mount tables, logical backends,
// credential stores, etc.
func (c *Core) postUnseal(ctx context.Context, ctxCancelFunc context.CancelFunc, unsealer UnsealStrategy) (retErr error) {
	defer metrics.MeasureSince([]string{"core", "post_unseal"}, time.Now())

	// Clear any out
	c.postUnsealFuncs = nil

	// Create a new request context
	c.activeContext = ctx
	c.activeContextCancelFunc.Store(ctxCancelFunc)

	defer func() {
		if retErr != nil {
			ctxCancelFunc()
			c.preSeal()
		}
	}()
	c.logger.Info("post-unseal setup starting")

	// Enable the cache
	c.physicalCache.Purge(ctx)
	if !c.cachingDisabled {
		c.physicalCache.SetEnabled(true)
	}

	// Purge these for safety in case of a rekey
	c.seal.SetBarrierConfig(ctx, nil)
	if c.seal.RecoveryKeySupported() {
		c.seal.SetRecoveryConfig(ctx, nil)
	}

	if err := unsealer.unseal(ctx, c.logger, c); err != nil {
		return err
	}

	c.metricsCh = make(chan struct{})
	go c.emitMetrics(c.metricsCh)

	// This is intentionally the last block in this function. We want to allow
	// writes just before allowing client requests, to ensure everything has
	// been set up properly before any writes can have happened.
	for _, v := range c.postUnsealFuncs {
		v()
	}

	c.logger.Info("post-unseal setup complete")
	return nil
}

// preSeal is invoked before the barrier is sealed, allowing
// for any state teardown required.
func (c *Core) preSeal() error {
	defer metrics.MeasureSince([]string{"core", "pre_seal"}, time.Now())
	c.logger.Info("pre-seal teardown starting")

	// Clear any pending funcs
	c.postUnsealFuncs = nil

	// Clear any rekey progress
	c.barrierRekeyConfig = nil
	c.recoveryRekeyConfig = nil

	if c.metricsCh != nil {
		close(c.metricsCh)
		c.metricsCh = nil
	}
	var result error

	c.clusterParamsLock.Lock()
	if err := stopReplication(c); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping replication: {{err}}", err))
	}
	c.clusterParamsLock.Unlock()

	c.stopForwarding()

	if err := c.teardownAudits(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down audits: {{err}}", err))
	}
	if err := c.stopExpiration(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping expiration: {{err}}", err))
	}
	if err := c.teardownCredentials(context.Background()); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down credentials: {{err}}", err))
	}
	if err := c.teardownPolicyStore(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down policy store: {{err}}", err))
	}
	if err := c.stopRollback(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping rollback: {{err}}", err))
	}
	if err := c.unloadMounts(context.Background()); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error unloading mounts: {{err}}", err))
	}
	if err := enterprisePreSeal(c); err != nil {
		result = multierror.Append(result, err)
	}

	preSealPhysical(c)

	c.logger.Info("pre-seal teardown complete")
	return result
}

func enterprisePostUnsealImpl(c *Core) error {
	return nil
}

func enterprisePreSealImpl(c *Core) error {
	return nil
}

func startReplicationImpl(c *Core) error {
	return nil
}

func stopReplicationImpl(c *Core) error {
	return nil
}

// emitMetrics is used to periodically expose metrics while running
func (c *Core) emitMetrics(stopCh chan struct{}) {
	emitTimer := time.Tick(time.Second)
	writeTimer := time.Tick(c.counters.syncInterval)

	for {
		select {
		case <-emitTimer:
			c.metricsMutex.Lock()
			if c.expiration != nil {
				c.expiration.emitMetrics()
			}
			c.metricsMutex.Unlock()

		case <-writeTimer:
			if c.perfStandby {
				syncCounter(c)
			} else {
				err := c.saveCurrentRequestCounters(context.Background(), time.Now())
				if err != nil {
					c.logger.Error("writing request counters to barrier", "err", err)
				}
			}

		case <-stopCh:
			return
		}
	}
}

func (c *Core) ReplicationState() consts.ReplicationState {
	return consts.ReplicationState(atomic.LoadUint32(c.replicationState))
}

func (c *Core) ActiveNodeReplicationState() consts.ReplicationState {
	return consts.ReplicationState(atomic.LoadUint32(c.activeNodeReplicationState))
}

func (c *Core) SealAccess() *SealAccess {
	return NewSealAccess(c.seal)
}

func (c *Core) Logger() log.Logger {
	return c.logger
}

func (c *Core) BarrierKeyLength() (min, max int) {
	min, max = c.barrier.KeyLength()
	max += shamir.ShareOverhead
	return
}

func (c *Core) AuditedHeadersConfig() *AuditedHeadersConfig {
	return c.auditedHeaders
}

func waitUntilWALShippedImpl(ctx context.Context, c *Core, index uint64) bool {
	return true
}

func lastWALImpl(c *Core) uint64 {
	return 0
}

func lastRemoteWALImpl(c *Core) uint64 {
	return 0
}

func (c *Core) PhysicalSealConfigs(ctx context.Context) (*SealConfig, *SealConfig, error) {
	pe, err := c.physical.Get(ctx, barrierSealConfigPath)
	if err != nil {
		return nil, nil, errwrap.Wrapf("failed to fetch barrier seal configuration at migration check time: {{err}}", err)
	}
	if pe == nil {
		return nil, nil, nil
	}

	barrierConf := new(SealConfig)

	if err := jsonutil.DecodeJSON(pe.Value, barrierConf); err != nil {
		return nil, nil, errwrap.Wrapf("failed to decode barrier seal configuration at migration check time: {{err}}", err)
	}
	err = barrierConf.Validate()
	if err != nil {
		return nil, nil, errwrap.Wrapf("failed to validate barrier seal configuration at migration check time: {{err}}", err)
	}
	// In older versions of vault the default seal would not store a type. This
	// is here to offer backwards compatibility for older seal configs.
	if barrierConf.Type == "" {
		barrierConf.Type = seal.Shamir
	}

	var recoveryConf *SealConfig
	pe, err = c.physical.Get(ctx, recoverySealConfigPlaintextPath)
	if err != nil {
		return nil, nil, errwrap.Wrapf("failed to fetch seal configuration at migration check time: {{err}}", err)
	}
	if pe != nil {
		recoveryConf = &SealConfig{}
		if err := jsonutil.DecodeJSON(pe.Value, recoveryConf); err != nil {
			return nil, nil, errwrap.Wrapf("failed to decode seal configuration at migration check time: {{err}}", err)
		}
		err = recoveryConf.Validate()
		if err != nil {
			return nil, nil, errwrap.Wrapf("failed to validate seal configuration at migration check time: {{err}}", err)
		}
		// In older versions of vault the default seal would not store a type. This
		// is here to offer backwards compatibility for older seal configs.
		if recoveryConf.Type == "" {
			recoveryConf.Type = seal.Shamir
		}
	}

	return barrierConf, recoveryConf, nil
}

func (c *Core) SetSealsForMigration(migrationSeal, newSeal, unwrapSeal Seal) {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	c.unwrapSeal = unwrapSeal
	if c.unwrapSeal != nil {
		c.unwrapSeal.SetCore(c)
	}
	if newSeal != nil && migrationSeal != nil {
		c.migrationSeal = migrationSeal
		c.migrationSeal.SetCore(c)
		c.seal = newSeal
		c.seal.SetCore(c)
		c.logger.Warn("entering seal migration mode; Vault will not automatically unseal even if using an autoseal", "from_barrier_type", c.migrationSeal.BarrierType(), "to_barrier_type", c.seal.BarrierType())
	}
}

func (c *Core) IsInSealMigration() bool {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.migrationSeal != nil
}

func (c *Core) BarrierEncryptorAccess() *BarrierEncryptorAccess {
	return NewBarrierEncryptorAccess(c.barrier)
}

func (c *Core) PhysicalAccess() *physical.PhysicalAccess {
	return physical.NewPhysicalAccess(c.physical)
}

func (c *Core) RouterAccess() *RouterAccess {
	return NewRouterAccess(c)
}

// IsDRSecondary returns if the current cluster state is a DR secondary.
func (c *Core) IsDRSecondary() bool {
	return c.ReplicationState().HasState(consts.ReplicationDRSecondary)
}

func (c *Core) AddLogger(logger log.Logger) {
	c.allLoggersLock.Lock()
	defer c.allLoggersLock.Unlock()
	c.allLoggers = append(c.allLoggers, logger)
}

func (c *Core) SetLogLevel(level log.Level) {
	c.allLoggersLock.RLock()
	defer c.allLoggersLock.RUnlock()
	for _, logger := range c.allLoggers {
		logger.SetLevel(level)
	}
}

// BuiltinRegistry is an interface that allows the "vault" package to use
// the registry of builtin plugins without getting an import cycle. It
// also allows for mocking the registry easily.
type BuiltinRegistry interface {
	Contains(name string, pluginType consts.PluginType) bool
	Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool)
	Keys(pluginType consts.PluginType) []string
}
