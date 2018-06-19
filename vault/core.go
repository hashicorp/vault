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

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"

	"google.golang.org/grpc"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/logging"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/helper/reload"
	"github.com/hashicorp/vault/helper/tlsutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
	cache "github.com/patrickmn/go-cache"
)

const (
	// coreLockPath is the path used to acquire a coordinating lock
	// for a highly-available deploy.
	coreLockPath = "core/lock"

	// The poison pill is used as a check during certain scenarios to indicate
	// to standby nodes that they should seal
	poisonPillPath = "core/poison-pill"

	// coreLeaderPrefix is the prefix used for the UUID that contains
	// the currently elected leader.
	coreLeaderPrefix = "core/leader/"

	// knownPrimaryAddrsPrefix is used to store last-known cluster address
	// information for primaries
	knownPrimaryAddrsPrefix = "core/primary-addrs/"

	// lockRetryInterval is the interval we re-attempt to acquire the
	// HA lock if an error is encountered
	lockRetryInterval = 10 * time.Second

	// leaderCheckInterval is how often a standby checks for a new leader
	leaderCheckInterval = 2500 * time.Millisecond

	// keyRotateCheckInterval is how often a standby checks for a key
	// rotation taking place.
	keyRotateCheckInterval = 30 * time.Second

	// keyRotateGracePeriod is how long we allow an upgrade path
	// for standby instances before we delete the upgrade keys
	keyRotateGracePeriod = 2 * time.Minute

	// leaderPrefixCleanDelay is how long to wait between deletions
	// of orphaned leader keys, to prevent slamming the backend.
	leaderPrefixCleanDelay = 200 * time.Millisecond

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
	LastRemoteWAL        = lastRemoteWALImpl
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

// ErrInvalidKey is returned if there is a user-based error with a provided
// unseal key. This will be shown to the user, so should not contain
// information that is sensitive.
type ErrInvalidKey struct {
	Reason string
}

func (e *ErrInvalidKey) Error() string {
	return fmt.Sprintf("invalid key: %v", e.Reason)
}

type activeAdvertisement struct {
	RedirectAddr     string            `json:"redirect_addr"`
	ClusterAddr      string            `json:"cluster_addr,omitempty"`
	ClusterCert      []byte            `json:"cluster_cert,omitempty"`
	ClusterKeyParams *clusterKeyParams `json:"cluster_key_params,omitempty"`
}

type unlockInformation struct {
	Parts [][]byte
	Nonce string
}

// Core is used as the central manager of Vault activity. It is the primary point of
// interface for API handlers and is responsible for managing the logical and physical
// backends, router, security barrier, and audit trails.
type Core struct {
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

	// Our Seal, for seal configuration information
	seal Seal

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
	sealed    bool

	standby              bool
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

	logger log.Logger

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
	// Tracks whether cluster listeners are running, e.g. it's safe to send a
	// shutdown down the channel
	clusterListenersRunning bool
	// Shutdown channel for the cluster listeners
	clusterListenerShutdownCh chan struct{}
	// Shutdown success channel. We need this to be done serially to ensure
	// that binds are removed before they might be reinstated.
	clusterListenerShutdownSuccessCh chan struct{}
	// Write lock used to ensure that we don't have multiple connections adjust
	// this value at the same time
	requestForwardingConnectionLock sync.RWMutex
	// Most recent leader UUID. Used to avoid repeatedly JSON parsing the same
	// values.
	clusterLeaderUUID string
	// Most recent leader redirect addr
	clusterLeaderRedirectAddr string
	// Most recent leader cluster addr
	clusterLeaderClusterAddr string
	// Lock for the cluster leader values
	clusterLeaderParamsLock sync.RWMutex
	// Info on cluster members
	clusterPeerClusterAddrsCache *cache.Cache
	// Stores whether we currently have a server running
	rpcServerActive *uint32
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
	activeContextCancelFunc context.CancelFunc

	// Stores the sealunwrapper for downgrade needs
	sealUnwrapper physical.Backend

	// Stores any funcs that should be run on successful postUnseal
	postUnsealFuncs []func()
}

// CoreConfig is used to parameterize a core
type CoreConfig struct {
	DevToken string `json:"dev_token" structs:"dev_token" mapstructure:"dev_token"`

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

	ReloadFuncs     *map[string][]reload.ReloadFunc
	ReloadFuncsLock *sync.RWMutex
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

	// Setup the core
	c := &Core{
		devToken:                         conf.DevToken,
		physical:                         conf.Physical,
		redirectAddr:                     conf.RedirectAddr,
		clusterAddr:                      conf.ClusterAddr,
		seal:                             conf.Seal,
		router:                           NewRouter(),
		sealed:                           true,
		standby:                          true,
		logger:                           conf.Logger.Named("core"),
		defaultLeaseTTL:                  conf.DefaultLeaseTTL,
		maxLeaseTTL:                      conf.MaxLeaseTTL,
		cachingDisabled:                  conf.DisableCache,
		clusterName:                      conf.ClusterName,
		clusterListenerShutdownCh:        make(chan struct{}),
		clusterListenerShutdownSuccessCh: make(chan struct{}),
		clusterPeerClusterAddrsCache:     cache.New(3*HeartbeatInterval, time.Second),
		enableMlock:                      !conf.DisableMlock,
		rawEnabled:                       conf.EnableRaw,
		replicationState:                 new(uint32),
		rpcServerActive:                  new(uint32),
		atomicPrimaryClusterAddrs:        new(atomic.Value),
		atomicPrimaryFailoverAddrs:       new(atomic.Value),
		localClusterPrivateKey:           new(atomic.Value),
		localClusterCert:                 new(atomic.Value),
		localClusterParsedCert:           new(atomic.Value),
		activeNodeReplicationState:       new(uint32),
		keepHALockOnStepDown:             new(uint32),
	}

	atomic.StoreUint32(c.replicationState, uint32(consts.ReplicationDRDisabled|consts.ReplicationPerformanceDisabled))
	c.localClusterCert.Store(([]byte)(nil))
	c.localClusterParsedCert.Store((*x509.Certificate)(nil))
	c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

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

	phys := conf.Physical
	_, txnOK := conf.Physical.(physical.Transactional)
	if c.seal == nil {
		c.seal = NewDefaultSeal()
	}
	c.seal.SetCore(c)

	c.sealUnwrapper = NewSealUnwrapper(phys, conf.Logger.ResetNamed("storage.sealunwrapper"))

	var ok bool

	// Wrap the physical backend in a cache layer if enabled
	if txnOK {
		c.physical = physical.NewTransactionalCache(c.sealUnwrapper, conf.CacheSize, conf.Logger.ResetNamed("storage.cache"))
	} else {
		c.physical = physical.NewCache(c.sealUnwrapper, conf.CacheSize, conf.Logger.ResetNamed("storage.cache"))
	}
	c.physicalCache = c.physical.(physical.ToggleablePurgemonster)

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

	// Setup the backends
	logicalBackends := make(map[string]logical.Factory)
	for k, f := range conf.LogicalBackends {
		logicalBackends[k] = f
	}
	_, ok = logicalBackends["kv"]
	if !ok {
		logicalBackends["kv"] = PassthroughBackendFactory
	}
	logicalBackends["cubbyhole"] = CubbyholeBackendFactory
	logicalBackends["system"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		b := NewSystemBackend(c, conf.Logger.Named("system"))
		if err := b.Setup(ctx, config); err != nil {
			return nil, err
		}
		return b, nil
	}

	logicalBackends["identity"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		return NewIdentityStore(ctx, c, config, conf.Logger.Named("identity"))
	}

	c.logicalBackends = logicalBackends

	credentialBackends := make(map[string]logical.Factory)
	for k, f := range conf.CredentialBackends {
		credentialBackends[k] = f
	}
	credentialBackends["token"] = func(ctx context.Context, config *logical.BackendConfig) (logical.Backend, error) {
		return NewTokenStore(ctx, conf.Logger.Named("token"), c, config)
	}
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
	c.stateLock.RLock()
	// Tell any requests that know about this to stop
	if c.activeContextCancelFunc != nil {
		c.activeContextCancelFunc()
	}
	c.stateLock.RUnlock()

	c.logger.Debug("shutdown initiating internal seal")
	// Seal the Vault, causes a leader stepdown
	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	c.logger.Debug("shutdown running internal seal")
	return c.sealInternal(false)
}

// CORSConfig returns the current CORS configuration
func (c *Core) CORSConfig() *CORSConfig {
	return c.corsConfig
}

func (c *Core) GetContext() (context.Context, context.CancelFunc) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()

	return context.WithCancel(c.activeContext)
}

// Sealed checks if the Vault is current sealed
func (c *Core) Sealed() (bool, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.sealed, nil
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
	if !c.sealed {
		return
	}
	c.unlockInfo = nil
}

// Unseal is used to provide one of the key parts to unseal the Vault.
//
// They key given as a parameter will automatically be zerod after
// this method is done with it. If you want to keep the key around, a copy
// should be made.
func (c *Core) Unseal(key []byte) (bool, error) {
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

	// Get the barrier seal configuration
	config, err := c.seal.BarrierConfig(ctx)
	if err != nil {
		return false, err
	}

	// Check if already unsealed
	if !c.sealed {
		return true, nil
	}

	masterKey, err := c.unsealPart(ctx, config, key, false)
	if err != nil {
		return false, err
	}
	if masterKey != nil {
		return c.unsealInternal(ctx, masterKey)
	}

	return false, nil
}

// UnsealWithRecoveryKeys is used to provide one of the recovery key shares to
// unseal the Vault.
func (c *Core) UnsealWithRecoveryKeys(ctx context.Context, key []byte) (bool, error) {
	defer metrics.MeasureSince([]string{"core", "unseal_with_recovery_keys"}, time.Now())

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// Explicitly check for init status
	init, err := c.Initialized(ctx)
	if err != nil {
		return false, err
	}
	if !init {
		return false, ErrNotInit
	}

	var config *SealConfig
	// If recovery keys are supported then use recovery seal config to unseal
	if c.seal.RecoveryKeySupported() {
		config, err = c.seal.RecoveryConfig(ctx)
		if err != nil {
			return false, err
		}
	}

	// Check if already unsealed
	if !c.sealed {
		return true, nil
	}

	masterKey, err := c.unsealPart(ctx, config, key, true)
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
func (c *Core) unsealPart(ctx context.Context, config *SealConfig, key []byte, useRecoveryKeys bool) ([]byte, error) {
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
	var err error
	if config.SecretThreshold == 1 {
		recoveredKey = make([]byte, len(c.unlockInfo.Parts[0]))
		copy(recoveredKey, c.unlockInfo.Parts[0])
	} else {
		recoveredKey, err = shamir.Combine(c.unlockInfo.Parts)
		if err != nil {
			return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
		}
	}

	if c.seal.RecoveryKeySupported() && useRecoveryKeys {
		// Verify recovery key
		if err := c.seal.VerifyRecoveryKey(ctx, recoveredKey); err != nil {
			return nil, err
		}

		// Get stored keys and shamir combine into single master key. Unsealing with
		// recovery keys currently does not support: 1) mixed stored and non-stored
		// keys setup, nor 2) seals that support recovery keys but not stored keys.
		// If insufficient shares are provided, shamir.Combine will error, and if
		// no stored keys are found it will return masterKey as nil.
		var masterKey []byte
		if c.seal.StoredKeysSupported() {
			masterKeyShares, err := c.seal.GetStoredKeys(ctx)
			if err != nil {
				return nil, errwrap.Wrapf("unable to retrieve stored keys: {{err}}", err)
			}

			if len(masterKeyShares) == 1 {
				return masterKeyShares[0], nil
			}

			masterKey, err = shamir.Combine(masterKeyShares)
			if err != nil {
				return nil, errwrap.Wrapf("failed to compute master key: {{err}}", err)
			}
		}
		return masterKey, nil
	}

	// If this is not a recovery key-supported seal, then the recovered key is
	// the master key to be returned.
	return recoveredKey, nil
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

		if err := c.postUnseal(); err != nil {
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

	// Success!
	c.sealed = false

	// Force a cache bust here, which will also run migration code
	if c.seal.RecoveryKeySupported() {
		c.seal.SetRecoveryConfig(ctx, nil)
	}

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
func (c *Core) SealWithRequest(req *logical.Request) error {
	defer metrics.MeasureSince([]string{"core", "seal-with-request"}, time.Now())

	c.stateLock.RLock()

	if c.sealed {
		c.stateLock.RUnlock()
		return nil
	}

	// This will unlock the read lock
	// We use background context since we may not be active
	return c.sealInitCommon(context.Background(), req)
}

// Seal takes in a token and creates a logical.Request, acquires the lock, and
// passes through to sealInternal
func (c *Core) Seal(token string) error {
	defer metrics.MeasureSince([]string{"core", "seal"}, time.Now())

	c.stateLock.RLock()

	if c.sealed {
		c.stateLock.RUnlock()
		return nil
	}

	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/seal",
		ClientToken: token,
	}

	// This will unlock the read lock
	// We use background context since we may not be active
	return c.sealInitCommon(context.Background(), req)
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

	acl, te, entity, identityPolicies, err := c.fetchACLTokenEntryAndEntity(req)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		c.stateLock.RUnlock()
		return retErr
	}

	// Audit-log the request before going any further
	auth := &logical.Auth{
		ClientToken:      req.ClientToken,
		Policies:         identityPolicies,
		IdentityPolicies: identityPolicies,
	}
	if te != nil {
		auth.TokenPolicies = te.Policies
		auth.Policies = append(te.Policies, identityPolicies...)
		auth.Metadata = te.Meta
		auth.DisplayName = te.DisplayName
		auth.EntityID = te.EntityID
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
	if authResults.Error.ErrorOrNil() != nil {
		retErr = multierror.Append(retErr, authResults.Error)
		c.stateLock.RUnlock()
		return retErr
	}
	if !authResults.Allowed {
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		c.stateLock.RUnlock()
		return retErr
	}

	if te != nil && te.NumUses == tokenRevocationPending {
		// Token needs to be revoked. We do this immediately here because
		// we won't have a token store after sealing.
		leaseID, err := c.expiration.CreateOrFetchRevocationLeaseByToken(te)
		if err == nil {
			err = c.expiration.Revoke(leaseID)
		}
		if err != nil {
			c.logger.Error("token needed revocation before seal but failed to revoke", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
		}
	}

	// Tell any requests that know about this to stop
	if c.activeContextCancelFunc != nil {
		c.activeContextCancelFunc()
	}

	// Unlock from the request handling
	c.stateLock.RUnlock()

	//Seal the Vault
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	sealErr := c.sealInternal(false)

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
// any authorization checking. The stateLock must be held prior to calling.
func (c *Core) sealInternal(keepLock bool) error {
	if c.sealed {
		return nil
	}

	// Enable that we are sealed to prevent further transactions
	c.sealed = true

	c.logger.Debug("marked as sealed")

	// Clear forwarding clients
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	// Do pre-seal teardown if HA is not enabled
	if c.ha == nil {
		// Even in a non-HA context we key off of this for some things
		c.standby = true
		if err := c.preSeal(); err != nil {
			c.logger.Error("pre-seal teardown failed", "error", err)
			return fmt.Errorf("internal error")
		}
	} else {
		if keepLock {
			atomic.StoreUint32(c.keepHALockOnStepDown, 1)
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

	c.logger.Info("vault is sealed")

	return nil
}

// postUnseal is invoked after the barrier is unsealed, but before
// allowing any user operations. This allows us to setup any state that
// requires the Vault to be unsealed such as mount tables, logical backends,
// credential stores, etc.
func (c *Core) postUnseal() (retErr error) {
	defer metrics.MeasureSince([]string{"core", "post_unseal"}, time.Now())

	// Clear any out
	c.postUnsealFuncs = nil

	// Create a new request context
	c.activeContext, c.activeContextCancelFunc = context.WithCancel(context.Background())

	defer func() {
		if retErr != nil {
			c.activeContextCancelFunc()
			c.preSeal()
		}
	}()
	c.logger.Info("post-unseal setup starting")

	// Clear forwarding clients; we're active
	c.requestForwardingConnectionLock.Lock()
	c.clearForwardingClients()
	c.requestForwardingConnectionLock.Unlock()

	// Enable the cache
	c.physicalCache.Purge(c.activeContext)
	if !c.cachingDisabled {
		c.physicalCache.SetEnabled(true)
	}

	switch c.sealUnwrapper.(type) {
	case *sealUnwrapper:
		c.sealUnwrapper.(*sealUnwrapper).runUnwraps()
	case *transactionalSealUnwrapper:
		c.sealUnwrapper.(*transactionalSealUnwrapper).runUnwraps()
	}

	// Purge these for safety in case of a rekey
	c.seal.SetBarrierConfig(c.activeContext, nil)
	if c.seal.RecoveryKeySupported() {
		c.seal.SetRecoveryConfig(c.activeContext, nil)
	}

	if err := enterprisePostUnseal(c); err != nil {
		return err
	}
	if err := c.ensureWrappingKey(c.activeContext); err != nil {
		return err
	}
	if err := c.setupPluginCatalog(); err != nil {
		return err
	}
	if err := c.loadMounts(c.activeContext); err != nil {
		return err
	}
	if err := c.setupMounts(c.activeContext); err != nil {
		return err
	}
	if err := c.setupPolicyStore(c.activeContext); err != nil {
		return err
	}
	if err := c.loadCORSConfig(c.activeContext); err != nil {
		return err
	}
	if err := c.loadCredentials(c.activeContext); err != nil {
		return err
	}
	if err := c.setupCredentials(c.activeContext); err != nil {
		return err
	}
	if err := c.startRollback(); err != nil {
		return err
	}
	if err := c.setupExpiration(); err != nil {
		return err
	}
	if err := c.loadAudits(c.activeContext); err != nil {
		return err
	}
	if err := c.setupAudits(c.activeContext); err != nil {
		return err
	}
	if err := c.loadIdentityStoreArtifacts(c.activeContext); err != nil {
		return err
	}
	if err := c.setupAuditedHeadersConfig(c.activeContext); err != nil {
		return err
	}

	if c.ha != nil {
		if err := c.startClusterListener(c.activeContext); err != nil {
			return err
		}
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

	c.stopClusterListener()

	if err := c.teardownAudits(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down audits: {{err}}", err))
	}
	if err := c.stopExpiration(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping expiration: {{err}}", err))
	}
	if err := c.teardownCredentials(c.activeContext); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down credentials: {{err}}", err))
	}
	if err := c.teardownPolicyStore(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down policy store: {{err}}", err))
	}
	if err := c.stopRollback(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping rollback: {{err}}", err))
	}
	if err := c.unloadMounts(c.activeContext); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error unloading mounts: {{err}}", err))
	}
	if err := enterprisePreSeal(c); err != nil {
		result = multierror.Append(result, err)
	}

	switch c.sealUnwrapper.(type) {
	case *sealUnwrapper:
		c.sealUnwrapper.(*sealUnwrapper).stopUnwraps()
	case *transactionalSealUnwrapper:
		c.sealUnwrapper.(*transactionalSealUnwrapper).stopUnwraps()
	}

	// Purge the cache
	c.physicalCache.SetEnabled(false)
	c.physicalCache.Purge(c.activeContext)

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
	for {
		select {
		case <-time.After(time.Second):
			c.metricsMutex.Lock()
			if c.expiration != nil {
				c.expiration.emitMetrics()
			}
			c.metricsMutex.Unlock()
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

func lastRemoteWALImpl(c *Core) uint64 {
	return 0
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
