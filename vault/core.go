package vault

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	log "github.com/mgutz/logxi/v1"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/helper/errutil"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/helper/logformat"
	"github.com/hashicorp/vault/helper/mlock"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/shamir"
)

const (
	// coreLockPath is the path used to acquire a coordinating lock
	// for a highly-available deploy.
	coreLockPath = "core/lock"

	// coreLeaderPrefix is the prefix used for the UUID that contains
	// the currently elected leader.
	coreLeaderPrefix = "core/leader/"

	// lockRetryInterval is the interval we re-attempt to acquire the
	// HA lock if an error is encountered
	lockRetryInterval = 10 * time.Second

	// keyRotateCheckInterval is how often a standby checks for a key
	// rotation taking place.
	keyRotateCheckInterval = 30 * time.Second

	// keyRotateGracePeriod is how long we allow an upgrade path
	// for standby instances before we delete the upgrade keys
	keyRotateGracePeriod = 2 * time.Minute

	// leaderPrefixCleanDelay is how long to wait between deletions
	// of orphaned leader keys, to prevent slamming the backend.
	leaderPrefixCleanDelay = 200 * time.Millisecond
)

var (
	// ErrSealed is returned if an operation is performed on
	// a sealed barrier. No operation is expected to succeed before unsealing
	ErrSealed = errors.New("Vault is sealed")

	// ErrStandby is returned if an operation is performed on
	// a standby Vault. No operation is expected to succeed until active.
	ErrStandby = errors.New("Vault is in standby mode")

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
)

// ReloadFunc are functions that are called when a reload is requested.
type ReloadFunc func(map[string]string) error

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

// Core is used as the central manager of Vault activity. It is the primary point of
// interface for API handlers and is responsible for managing the logical and physical
// backends, router, security barrier, and audit trails.
type Core struct {
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

	standby          bool
	standbyDoneCh    chan struct{}
	standbyStopCh    chan struct{}
	manualStepDownCh chan struct{}

	// unlockParts has the keys provided to Unseal until
	// the threshold number of parts is available.
	unlockParts [][]byte

	// generateRootProgress holds the shares until we reach enough
	// to verify the master key
	generateRootConfig   *GenerateRootConfig
	generateRootProgress [][]byte
	generateRootLock     sync.Mutex

	// These variables holds the config and shares we have until we reach
	// enough to verify the appropriate master key. Note that the same lock is
	// used; this isn't time-critical so this shouldn't be a problem.
	barrierRekeyConfig    *SealConfig
	barrierRekeyProgress  [][]byte
	recoveryRekeyConfig   *SealConfig
	recoveryRekeyProgress [][]byte
	rekeyLock             sync.RWMutex

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

	// reloadFuncs is a map containing reload functions
	reloadFuncs map[string][]ReloadFunc

	// reloadFuncsLock controlls access to the funcs
	reloadFuncsLock sync.RWMutex

	//
	// Cluster information
	//
	// Name
	clusterName string
	// Used to modify cluster TLS params
	clusterParamsLock sync.RWMutex
	// The private key stored in the barrier used for establishing
	// mutually-authenticated connections between Vault cluster members
	localClusterPrivateKey crypto.Signer
	// The local cluster cert
	localClusterCert []byte
	// The cert pool containing the self-signed CA as a trusted CA
	localClusterCertPool *x509.CertPool
	// The TCP addresses we should use for clustering
	clusterListenerAddrs []*net.TCPAddr
	// The setup function that gives us the handler to use
	clusterHandlerSetupFunc func() (http.Handler, http.Handler)
	// Shutdown channel for the cluster listeners
	clusterListenerShutdownCh chan struct{}
	// Shutdown success channel. We need this to be done serially to ensure
	// that binds are removed before they might be reinstated.
	clusterListenerShutdownSuccessCh chan struct{}
	// Connection info containing a client and a current active address
	requestForwardingConnection *activeConnection
	// Write lock used to ensure that we don't have multiple connections adjust
	// this value at the same time
	requestForwardingConnectionLock sync.RWMutex
	// Most recent leader UUID. Used to avoid repeatedly JSON parsing the same
	// values.
	clusterLeaderUUID string
	// Most recent leader redirect addr
	clusterLeaderRedirectAddr string
	// The grpc Server that handles server RPC calls
	rpcServer *grpc.Server
	// The function for canceling the client connection
	rpcClientConnCancelFunc context.CancelFunc
	// The grpc ClientConn for RPC calls
	rpcClientConn *grpc.ClientConn
	// The grpc forwarding client
	rpcForwardingClient RequestForwardingClient
}

// CoreConfig is used to parameterize a core
type CoreConfig struct {
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

	ReloadFuncs     *map[string][]ReloadFunc
	ReloadFuncsLock *sync.RWMutex
}

// NewCore is used to construct a new core
func NewCore(conf *CoreConfig) (*Core, error) {
	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		if conf.RedirectAddr == "" {
			return nil, fmt.Errorf("missing redirect address")
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
			return nil, fmt.Errorf("redirect address is not valid url: %s", err)
		}

		if u.Scheme == "" {
			return nil, fmt.Errorf("redirect address must include scheme (ex. 'http')")
		}
	}

	// Make a default logger if not provided
	if conf.Logger == nil {
		conf.Logger = logformat.NewVaultLogger(log.LevelTrace)
	}

	// Wrap the backend in a cache unless disabled
	if !conf.DisableCache {
		_, isCache := conf.Physical.(*physical.Cache)
		_, isInmem := conf.Physical.(*physical.InmemBackend)
		if !isCache && !isInmem {
			cache := physical.NewCache(conf.Physical, conf.CacheSize, conf.Logger)
			conf.Physical = cache
		}
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

	// Construct a new AES-GCM barrier
	barrier, err := NewAESGCMBarrier(conf.Physical)
	if err != nil {
		return nil, fmt.Errorf("barrier setup failed: %v", err)
	}

	// Setup the core
	c := &Core{
		redirectAddr:                     conf.RedirectAddr,
		clusterAddr:                      conf.ClusterAddr,
		physical:                         conf.Physical,
		seal:                             conf.Seal,
		barrier:                          barrier,
		router:                           NewRouter(),
		sealed:                           true,
		standby:                          true,
		logger:                           conf.Logger,
		defaultLeaseTTL:                  conf.DefaultLeaseTTL,
		maxLeaseTTL:                      conf.MaxLeaseTTL,
		cachingDisabled:                  conf.DisableCache,
		clusterName:                      conf.ClusterName,
		localClusterCertPool:             x509.NewCertPool(),
		clusterListenerShutdownCh:        make(chan struct{}),
		clusterListenerShutdownSuccessCh: make(chan struct{}),
	}

	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		c.ha = conf.HAPhysical
	}

	// We create the funcs here, then populate the given config with it so that
	// the caller can share state
	conf.ReloadFuncsLock = &c.reloadFuncsLock
	c.reloadFuncsLock.Lock()
	c.reloadFuncs = make(map[string][]ReloadFunc)
	c.reloadFuncsLock.Unlock()
	conf.ReloadFuncs = &c.reloadFuncs

	// Setup the backends
	logicalBackends := make(map[string]logical.Factory)
	for k, f := range conf.LogicalBackends {
		logicalBackends[k] = f
	}
	_, ok := logicalBackends["generic"]
	if !ok {
		logicalBackends["generic"] = PassthroughBackendFactory
	}
	logicalBackends["cubbyhole"] = CubbyholeBackendFactory
	logicalBackends["system"] = func(config *logical.BackendConfig) (logical.Backend, error) {
		return NewSystemBackend(c, config)
	}
	c.logicalBackends = logicalBackends

	credentialBackends := make(map[string]logical.Factory)
	for k, f := range conf.CredentialBackends {
		credentialBackends[k] = f
	}
	credentialBackends["token"] = func(config *logical.BackendConfig) (logical.Backend, error) {
		return NewTokenStore(c, config)
	}
	c.credentialBackends = credentialBackends

	auditBackends := make(map[string]audit.Factory)
	for k, f := range conf.AuditBackends {
		auditBackends[k] = f
	}
	c.auditBackends = auditBackends

	if c.seal == nil {
		c.seal = &DefaultSeal{}
	}
	c.seal.SetCore(c)

	// Attempt unsealing with stored keys; if there are no stored keys this
	// returns nil, otherwise returns nil or an error
	storedKeyErr := c.UnsealWithStoredKeys()

	return c, storedKeyErr
}

// Shutdown is invoked when the Vault instance is about to be terminated. It
// should not be accessible as part of an API call as it will cause an availability
// problem. It is only used to gracefully quit in the case of HA so that failover
// happens as quickly as possible.
func (c *Core) Shutdown() error {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if c.sealed {
		return nil
	}

	// Seal the Vault, causes a leader stepdown
	return c.sealInternal()
}

func (c *Core) fetchACLandTokenEntry(req *logical.Request) (*ACL, *TokenEntry, error) {
	defer metrics.MeasureSince([]string{"core", "fetch_acl_and_token"}, time.Now())

	// Ensure there is a client token
	if req.ClientToken == "" {
		return nil, nil, fmt.Errorf("missing client token")
	}

	if c.tokenStore == nil {
		c.logger.Error("core: token store is unavailable")
		return nil, nil, ErrInternalError
	}

	// Resolve the token policy
	te, err := c.tokenStore.Lookup(req.ClientToken)
	if err != nil {
		c.logger.Error("core: failed to lookup token", "error", err)
		return nil, nil, ErrInternalError
	}

	// Ensure the token is valid
	if te == nil {
		return nil, nil, logical.ErrPermissionDenied
	}

	// Construct the corresponding ACL object
	acl, err := c.policyStore.ACL(te.Policies...)
	if err != nil {
		c.logger.Error("core: failed to construct ACL", "error", err)
		return nil, nil, ErrInternalError
	}

	return acl, te, nil
}

func (c *Core) checkToken(req *logical.Request) (*logical.Auth, *TokenEntry, error) {
	defer metrics.MeasureSince([]string{"core", "check_token"}, time.Now())

	acl, te, err := c.fetchACLandTokenEntry(req)
	if err != nil {
		return nil, te, err
	}

	// Check if this is a root protected path
	rootPath := c.router.RootPath(req.Path)

	// When we receive a write of either type, rather than require clients to
	// PUT/POST and trust the operation, we ask the backend to give us the real
	// skinny -- if the backend implements an existence check, it can tell us
	// whether a particular resource exists. Then we can mark it as an update
	// or creation as appropriate.
	if req.Operation == logical.CreateOperation || req.Operation == logical.UpdateOperation {
		checkExists, resourceExists, err := c.router.RouteExistenceCheck(req)
		switch err {
		case logical.ErrUnsupportedPath:
			// fail later via bad path to avoid confusing items in the log
			checkExists = false
		case nil:
			// Continue on
		default:
			c.logger.Error("core: failed to run existence check", "error", err)
			if _, ok := err.(errutil.UserError); ok {
				return nil, nil, err
			} else {
				return nil, nil, ErrInternalError
			}
		}

		switch {
		case checkExists == false:
			// No existence check, so always treate it as an update operation, which is how it is pre 0.5
			req.Operation = logical.UpdateOperation
		case resourceExists == true:
			// It exists, so force an update operation
			req.Operation = logical.UpdateOperation
		case resourceExists == false:
			// It doesn't exist, force a create operation
			req.Operation = logical.CreateOperation
		default:
			panic("unreachable code")
		}
	}

	// Check the standard non-root ACLs. Return the token entry if it's not
	// allowed so we can decrement the use count.
	allowed, rootPrivs := acl.AllowOperation(req.Operation, req.Path)
	if !allowed {
		return nil, te, logical.ErrPermissionDenied
	}
	if rootPath && !rootPrivs {
		return nil, te, logical.ErrPermissionDenied
	}

	// Create the auth response
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		DisplayName: te.DisplayName,
	}
	return auth, te, nil
}

// Sealed checks if the Vault is current sealed
func (c *Core) Sealed() (bool, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.sealed, nil
}

// Standby checks if the Vault is in standby mode
func (c *Core) Standby() (bool, error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.standby, nil
}

// Leader is used to get the current active leader
func (c *Core) Leader() (isLeader bool, leaderAddr string, err error) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	// Check if HA enabled
	if c.ha == nil {
		return false, "", ErrHANotEnabled
	}

	// Check if sealed
	if c.sealed {
		return false, "", ErrSealed
	}

	// Check if we are the leader
	if !c.standby {
		// If we have connections from talking to a previous leader, close them
		// out to free resources
		if c.requestForwardingConnection != nil {
			c.requestForwardingConnectionLock.Lock()
			// Verify that the condition hasn't changed
			if c.requestForwardingConnection != nil {
				c.requestForwardingConnection.transport.CloseIdleConnections()
			}
			c.requestForwardingConnection = nil
			c.requestForwardingConnectionLock.Unlock()
		}
		return true, c.redirectAddr, nil
	}

	// Initialize a lock
	lock, err := c.ha.LockWith(coreLockPath, "read")
	if err != nil {
		return false, "", err
	}

	// Read the value
	held, leaderUUID, err := lock.Value()
	if err != nil {
		return false, "", err
	}
	if !held {
		return false, "", nil
	}

	// If the leader hasn't changed, return the cached value; nothing changes
	// mid-leadership, and the barrier caches anyways
	if leaderUUID == c.clusterLeaderUUID && c.clusterLeaderRedirectAddr != "" {
		return false, c.clusterLeaderRedirectAddr, nil
	}

	key := coreLeaderPrefix + leaderUUID
	entry, err := c.barrier.Get(key)
	if err != nil {
		return false, "", err
	}
	if entry == nil {
		return false, "", nil
	}

	var oldAdv bool

	var adv activeAdvertisement
	err = jsonutil.DecodeJSON(entry.Value, &adv)
	if err != nil {
		// Fall back to pre-struct handling
		adv.RedirectAddr = string(entry.Value)
		oldAdv = true
	}

	if !oldAdv {
		// Ensure we are using current values
		err = c.loadClusterTLS(adv)
		if err != nil {
			return false, "", err
		}

		// This will ensure that we both have a connection at the ready and that
		// the address is the current known value
		err = c.refreshRequestForwardingConnection(adv.ClusterAddr)
		if err != nil {
			return false, "", err
		}
	}

	// Don't set these until everything has been parsed successfully or we'll
	// never try again
	c.clusterLeaderRedirectAddr = adv.RedirectAddr
	c.clusterLeaderUUID = leaderUUID

	return false, c.clusterLeaderRedirectAddr, nil
}

// SecretProgress returns the number of keys provided so far
func (c *Core) SecretProgress() int {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return len(c.unlockParts)
}

// ResetUnsealProcess removes the current unlock parts from memory, to reset
// the unsealing process
func (c *Core) ResetUnsealProcess() {
	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if !c.sealed {
		return
	}
	c.unlockParts = nil
}

// Unseal is used to provide one of the key parts to unseal the Vault.
//
// They key given as a parameter will automatically be zerod after
// this method is done with it. If you want to keep the key around, a copy
// should be made.
func (c *Core) Unseal(key []byte) (bool, error) {
	defer metrics.MeasureSince([]string{"core", "unseal"}, time.Now())

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return false, &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return false, &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Get the seal configuration
	config, err := c.seal.BarrierConfig()
	if err != nil {
		return false, err
	}

	// Ensure the barrier is initialized
	if config == nil {
		return false, ErrNotInit
	}

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	// Check if already unsealed
	if !c.sealed {
		return true, nil
	}

	// Check if we already have this piece
	for _, existing := range c.unlockParts {
		if bytes.Equal(existing, key) {
			return false, nil
		}
	}

	// Store this key
	c.unlockParts = append(c.unlockParts, key)

	// Check if we don't have enough keys to unlock
	if len(c.unlockParts) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("core: cannot unseal, not enough keys", "keys", len(c.unlockParts), "threshold", config.SecretThreshold)
		}
		return false, nil
	}

	// Recover the master key
	var masterKey []byte
	if config.SecretThreshold == 1 {
		masterKey = c.unlockParts[0]
		c.unlockParts = nil
	} else {
		masterKey, err = shamir.Combine(c.unlockParts)
		c.unlockParts = nil
		if err != nil {
			return false, fmt.Errorf("failed to compute master key: %v", err)
		}
	}
	defer memzero(masterKey)

	// Attempt to unlock
	if err := c.barrier.Unseal(masterKey); err != nil {
		return false, err
	}
	if c.logger.IsInfo() {
		c.logger.Info("core: vault is unsealed")
	}

	// Do post-unseal setup if HA is not enabled
	if c.ha == nil {
		// We still need to set up cluster info even if it's not part of a
		// cluster right now
		if err := c.setupCluster(); err != nil {
			c.logger.Error("core: cluster setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("core: vault is sealed")
			return false, err
		}
		if err := c.postUnseal(); err != nil {
			c.logger.Error("core: post-unseal setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("core: vault is sealed")
			return false, err
		}
		c.standby = false
	} else {
		// Go to standby mode, wait until we are active to unseal
		c.standbyDoneCh = make(chan struct{})
		c.standbyStopCh = make(chan struct{})
		c.manualStepDownCh = make(chan struct{})
		go c.runStandby(c.standbyDoneCh, c.standbyStopCh, c.manualStepDownCh)
	}

	// Success!
	c.sealed = false
	if c.ha != nil {
		sd, ok := c.ha.(physical.ServiceDiscovery)
		if ok {
			if err := sd.NotifySealedStateChange(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn("core: failed to notify unsealed status", "error", err)
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

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	if c.sealed {
		return nil
	}

	return c.sealInitCommon(req)
}

// Seal takes in a token and creates a logical.Request, acquires the lock, and
// passes through to sealInternal
func (c *Core) Seal(token string) error {
	defer metrics.MeasureSince([]string{"core", "seal"}, time.Now())

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	if c.sealed {
		return nil
	}

	req := &logical.Request{
		Operation:   logical.UpdateOperation,
		Path:        "sys/seal",
		ClientToken: token,
	}

	return c.sealInitCommon(req)
}

// sealInitCommon is common logic for Seal and SealWithRequest and is used to
// re-seal the Vault. This requires the Vault to be unsealed again to perform
// any further operations.
func (c *Core) sealInitCommon(req *logical.Request) (retErr error) {
	defer metrics.MeasureSince([]string{"core", "seal-internal"}, time.Now())

	if req == nil {
		retErr = multierror.Append(retErr, errors.New("nil request to seal"))
		return retErr
	}

	// Validate the token is a root token
	acl, te, err := c.fetchACLandTokenEntry(req)
	if err != nil {
		// Since there is no token store in standby nodes, sealing cannot
		// be done. Ideally, the request has to be forwarded to leader node
		// for validation and the operation should be performed. But for now,
		// just returning with an error and recommending a vault restart, which
		// essentially does the same thing.
		if c.standby {
			c.logger.Error("core: vault cannot seal when in standby mode; please restart instead")
			retErr = multierror.Append(retErr, errors.New("vault cannot seal when in standby mode; please restart instead"))
			return retErr
		}
		retErr = multierror.Append(retErr, err)
		return retErr
	}

	// Audit-log the request before going any further
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		DisplayName: te.DisplayName,
	}

	if err := c.auditBroker.LogRequest(auth, req, nil); err != nil {
		c.logger.Error("core: failed to audit request", "request_path", req.Path, "error", err)
		retErr = multierror.Append(retErr, errors.New("failed to audit request, cannot continue"))
		return retErr
	}

	// Attempt to use the token (decrement num_uses)
	// On error bail out; if the token has been revoked, bail out too
	if te != nil {
		te, err = c.tokenStore.UseToken(te)
		if err != nil {
			c.logger.Error("core: failed to use token", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return retErr
		}
		if te == nil {
			// Token is no longer valid
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			return retErr
		}
		if te.NumUses == -1 {
			// Token needs to be revoked
			defer func(id string) {
				err = c.tokenStore.Revoke(id)
				if err != nil {
					c.logger.Error("core: token needed revocation after seal but failed to revoke", "error", err)
					retErr = multierror.Append(retErr, ErrInternalError)
				}
			}(te.ID)
		}
	}

	// Verify that this operation is allowed
	allowed, rootPrivs := acl.AllowOperation(req.Operation, req.Path)
	if !allowed {
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		return retErr
	}

	// We always require root privileges for this operation
	if !rootPrivs {
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		return retErr
	}

	//Seal the Vault
	err = c.sealInternal()
	if err != nil {
		retErr = multierror.Append(retErr, err)
	}

	return retErr
}

// StepDown is used to step down from leadership
func (c *Core) StepDown(req *logical.Request) (retErr error) {
	defer metrics.MeasureSince([]string{"core", "step_down"}, time.Now())

	if req == nil {
		retErr = multierror.Append(retErr, errors.New("nil request to step-down"))
		return retErr
	}

	c.stateLock.Lock()
	defer c.stateLock.Unlock()
	if c.sealed {
		return nil
	}
	if c.ha == nil || c.standby {
		return nil
	}

	acl, te, err := c.fetchACLandTokenEntry(req)
	if err != nil {
		retErr = multierror.Append(retErr, err)
		return retErr
	}

	// Audit-log the request before going any further
	auth := &logical.Auth{
		ClientToken: req.ClientToken,
		Policies:    te.Policies,
		Metadata:    te.Meta,
		DisplayName: te.DisplayName,
	}

	if err := c.auditBroker.LogRequest(auth, req, nil); err != nil {
		c.logger.Error("core: failed to audit request", "request_path", req.Path, "error", err)
		retErr = multierror.Append(retErr, errors.New("failed to audit request, cannot continue"))
		return retErr
	}

	// Attempt to use the token (decrement num_uses)
	if te != nil {
		te, err = c.tokenStore.UseToken(te)
		if err != nil {
			c.logger.Error("core: failed to use token", "error", err)
			retErr = multierror.Append(retErr, ErrInternalError)
			return retErr
		}
		if te == nil {
			// Token has been revoked
			retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
			return retErr
		}
		if te.NumUses == -1 {
			// Token needs to be revoked
			defer func(id string) {
				err = c.tokenStore.Revoke(id)
				if err != nil {
					c.logger.Error("core: token needed revocation after step-down but failed to revoke", "error", err)
					retErr = multierror.Append(retErr, ErrInternalError)
				}
			}(te.ID)
		}
	}

	// Verify that this operation is allowed
	allowed, rootPrivs := acl.AllowOperation(req.Operation, req.Path)
	if !allowed {
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		return retErr
	}

	// We always require root privileges for this operation
	if !rootPrivs {
		retErr = multierror.Append(retErr, logical.ErrPermissionDenied)
		return retErr
	}

	select {
	case c.manualStepDownCh <- struct{}{}:
	default:
		c.logger.Warn("core: manual step-down operation already queued")
	}

	return retErr
}

// sealInternal is an internal method used to seal the vault.  It does not do
// any authorization checking. The stateLock must be held prior to calling.
func (c *Core) sealInternal() error {
	// Enable that we are sealed to prevent furthur transactions
	c.sealed = true

	// Do pre-seal teardown if HA is not enabled
	if c.ha == nil {
		if err := c.preSeal(); err != nil {
			c.logger.Error("core: pre-seal teardown failed", "error", err)
			return fmt.Errorf("internal error")
		}
	} else {
		// Signal the standby goroutine to shutdown, wait for completion
		close(c.standbyStopCh)

		// Release the lock while we wait to avoid deadlocking
		c.stateLock.Unlock()
		<-c.standbyDoneCh
		c.stateLock.Lock()
	}

	if err := c.barrier.Seal(); err != nil {
		return err
	}
	c.logger.Info("core: vault is sealed")

	if c.ha != nil {
		sd, ok := c.ha.(physical.ServiceDiscovery)
		if ok {
			if err := sd.NotifySealedStateChange(); err != nil {
				if c.logger.IsWarn() {
					c.logger.Warn("core: failed to notify sealed status", "error", err)
				}
			}
		}
	}

	return nil
}

// postUnseal is invoked after the barrier is unsealed, but before
// allowing any user operations. This allows us to setup any state that
// requires the Vault to be unsealed such as mount tables, logical backends,
// credential stores, etc.
func (c *Core) postUnseal() (retErr error) {
	defer metrics.MeasureSince([]string{"core", "post_unseal"}, time.Now())
	defer func() {
		if retErr != nil {
			c.preSeal()
		}
	}()
	c.logger.Info("core: post-unseal setup starting")
	if cache, ok := c.physical.(*physical.Cache); ok {
		cache.Purge()
	}
	// HA mode requires us to handle keyring rotation and rekeying
	if c.ha != nil {
		if err := c.checkKeyUpgrades(); err != nil {
			return err
		}
		if err := c.barrier.ReloadMasterKey(); err != nil {
			return err
		}
		if err := c.barrier.ReloadKeyring(); err != nil {
			return err
		}
		if err := c.scheduleUpgradeCleanup(); err != nil {
			return err
		}
	}
	if err := c.loadMounts(); err != nil {
		return err
	}
	if err := c.setupMounts(); err != nil {
		return err
	}
	if err := c.startRollback(); err != nil {
		return err
	}
	if err := c.setupPolicyStore(); err != nil {
		return err
	}
	if err := c.loadCredentials(); err != nil {
		return err
	}
	if err := c.setupCredentials(); err != nil {
		return err
	}
	if err := c.setupExpiration(); err != nil {
		return err
	}
	if err := c.loadAudits(); err != nil {
		return err
	}
	if err := c.setupAudits(); err != nil {
		return err
	}
	if c.ha != nil {
		if err := c.startClusterListener(); err != nil {
			return err
		}
	}
	c.metricsCh = make(chan struct{})
	go c.emitMetrics(c.metricsCh)
	c.logger.Info("core: post-unseal setup complete")
	return nil
}

// preSeal is invoked before the barrier is sealed, allowing
// for any state teardown required.
func (c *Core) preSeal() error {
	defer metrics.MeasureSince([]string{"core", "pre_seal"}, time.Now())
	c.logger.Info("core: pre-seal teardown starting")

	// Clear any rekey progress
	c.barrierRekeyConfig = nil
	c.barrierRekeyProgress = nil
	c.recoveryRekeyConfig = nil
	c.recoveryRekeyProgress = nil

	if c.metricsCh != nil {
		close(c.metricsCh)
		c.metricsCh = nil
	}
	var result error
	if c.ha != nil {
		c.stopClusterListener()
	}

	if err := c.teardownAudits(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down audits: {{err}}", err))
	}
	if err := c.stopExpiration(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping expiration: {{err}}", err))
	}
	if err := c.teardownCredentials(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down credentials: {{err}}", err))
	}
	if err := c.teardownPolicyStore(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error tearing down policy store: {{err}}", err))
	}
	if err := c.stopRollback(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error stopping rollback: {{err}}", err))
	}
	if err := c.unloadMounts(); err != nil {
		result = multierror.Append(result, errwrap.Wrapf("error unloading mounts: {{err}}", err))
	}
	if cache, ok := c.physical.(*physical.Cache); ok {
		cache.Purge()
	}
	c.logger.Info("core: pre-seal teardown complete")
	return result
}

// runStandby is a long running routine that is used when an HA backend
// is enabled. It waits until we are leader and switches this Vault to
// active.
func (c *Core) runStandby(doneCh, stopCh, manualStepDownCh chan struct{}) {
	defer close(doneCh)
	defer close(manualStepDownCh)
	c.logger.Info("core: entering standby mode")

	// Monitor for key rotation
	keyRotateDone := make(chan struct{})
	keyRotateStop := make(chan struct{})
	go c.periodicCheckKeyUpgrade(keyRotateDone, keyRotateStop)
	defer func() {
		close(keyRotateStop)
		<-keyRotateDone
	}()

	for {
		// Check for a shutdown
		select {
		case <-stopCh:
			return
		default:
		}

		// Create a lock
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			c.logger.Error("core: failed to generate uuid", "error", err)
			return
		}
		lock, err := c.ha.LockWith(coreLockPath, uuid)
		if err != nil {
			c.logger.Error("core: failed to create lock", "error", err)
			return
		}

		// Attempt the acquisition
		leaderLostCh := c.acquireLock(lock, stopCh)

		// Bail if we are being shutdown
		if leaderLostCh == nil {
			return
		}
		c.logger.Info("core: acquired lock, enabling active operation")

		// This is used later to log a metrics event; this can be helpful to
		// detect flapping
		activeTime := time.Now()

		// Grab the lock as we need it for cluster setup, which needs to happen
		// before advertising
		c.stateLock.Lock()
		if err := c.setupCluster(); err != nil {
			c.stateLock.Unlock()
			c.logger.Error("core: cluster setup failed", "error", err)
			lock.Unlock()
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			continue
		}

		// Advertise as leader
		if err := c.advertiseLeader(uuid, leaderLostCh); err != nil {
			c.stateLock.Unlock()
			c.logger.Error("core: leader advertisement setup failed", "error", err)
			lock.Unlock()
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			continue
		}

		// Attempt the post-unseal process
		err = c.postUnseal()
		if err == nil {
			c.standby = false
		}
		c.stateLock.Unlock()

		// Handle a failure to unseal
		if err != nil {
			c.logger.Error("core: post-unseal setup failed", "error", err)
			lock.Unlock()
			metrics.MeasureSince([]string{"core", "leadership_setup_failed"}, activeTime)
			continue
		}

		// Monitor a loss of leadership
		var manualStepDown bool
		select {
		case <-leaderLostCh:
			c.logger.Warn("core: leadership lost, stopping active operation")
		case <-stopCh:
			c.logger.Warn("core: stopping active operation")
		case <-manualStepDownCh:
			c.logger.Warn("core: stepping down from active operation to standby")
			manualStepDown = true
		}

		metrics.MeasureSince([]string{"core", "leadership_lost"}, activeTime)

		// Clear ourself as leader
		if err := c.clearLeader(uuid); err != nil {
			c.logger.Error("core: clearing leader advertisement failed", "error", err)
		}

		// Attempt the pre-seal process
		c.stateLock.Lock()
		c.standby = true
		preSealErr := c.preSeal()
		c.stateLock.Unlock()

		// Give up leadership
		lock.Unlock()

		// Check for a failure to prepare to seal
		if preSealErr != nil {
			c.logger.Error("core: pre-seal teardown failed", "error", err)
		}

		// If we've merely stepped down, we could instantly grab the lock
		// again. Give the other nodes a chance.
		if manualStepDown {
			time.Sleep(manualStepDownSleepPeriod)
		}
	}
}

// periodicCheckKeyUpgrade is used to watch for key rotation events as a standby
func (c *Core) periodicCheckKeyUpgrade(doneCh, stopCh chan struct{}) {
	defer close(doneCh)
	for {
		select {
		case <-time.After(keyRotateCheckInterval):
			// Only check if we are a standby
			c.stateLock.RLock()
			standby := c.standby
			c.stateLock.RUnlock()
			if !standby {
				continue
			}

			if err := c.checkKeyUpgrades(); err != nil {
				c.logger.Error("core: key rotation periodic upgrade check failed", "error", err)
			}
		case <-stopCh:
			return
		}
	}
}

// checkKeyUpgrades is used to check if there have been any key rotations
// and if there is a chain of upgrades available
func (c *Core) checkKeyUpgrades() error {
	for {
		// Check for an upgrade
		didUpgrade, newTerm, err := c.barrier.CheckUpgrade()
		if err != nil {
			return err
		}

		// Nothing to do if no upgrade
		if !didUpgrade {
			break
		}
		if c.logger.IsInfo() {
			c.logger.Info("core: upgraded to new key term", "term", newTerm)
		}
	}
	return nil
}

// scheduleUpgradeCleanup is used to ensure that all the upgrade paths
// are cleaned up in a timely manner if a leader failover takes place
func (c *Core) scheduleUpgradeCleanup() error {
	// List the upgrades
	upgrades, err := c.barrier.List(keyringUpgradePrefix)
	if err != nil {
		return fmt.Errorf("failed to list upgrades: %v", err)
	}

	// Nothing to do if no upgrades
	if len(upgrades) == 0 {
		return nil
	}

	// Schedule cleanup for all of them
	time.AfterFunc(keyRotateGracePeriod, func() {
		for _, upgrade := range upgrades {
			path := fmt.Sprintf("%s%s", keyringUpgradePrefix, upgrade)
			if err := c.barrier.Delete(path); err != nil {
				c.logger.Error("core: failed to cleanup upgrade", "path", path, "error", err)
			}
		}
	})
	return nil
}

// acquireLock blocks until the lock is acquired, returning the leaderLostCh
func (c *Core) acquireLock(lock physical.Lock, stopCh <-chan struct{}) <-chan struct{} {
	for {
		// Attempt lock acquisition
		leaderLostCh, err := lock.Lock(stopCh)
		if err == nil {
			return leaderLostCh
		}

		// Retry the acquisition
		c.logger.Error("core: failed to acquire lock", "error", err)
		select {
		case <-time.After(lockRetryInterval):
		case <-stopCh:
			return nil
		}
	}
}

// advertiseLeader is used to advertise the current node as leader
func (c *Core) advertiseLeader(uuid string, leaderLostCh <-chan struct{}) error {
	go c.cleanLeaderPrefix(uuid, leaderLostCh)

	var key *ecdsa.PrivateKey
	switch c.localClusterPrivateKey.(type) {
	case *ecdsa.PrivateKey:
		key = c.localClusterPrivateKey.(*ecdsa.PrivateKey)
	default:
		c.logger.Error("core: unknown cluster private key type", "key_type", fmt.Sprintf("%T", c.localClusterPrivateKey))
		return fmt.Errorf("unknown cluster private key type %T", c.localClusterPrivateKey)
	}

	keyParams := &clusterKeyParams{
		Type: corePrivateKeyTypeP521,
		X:    key.X,
		Y:    key.Y,
		D:    key.D,
	}

	adv := &activeAdvertisement{
		RedirectAddr:     c.redirectAddr,
		ClusterAddr:      c.clusterAddr,
		ClusterCert:      c.localClusterCert,
		ClusterKeyParams: keyParams,
	}
	val, err := jsonutil.EncodeJSON(adv)
	if err != nil {
		return err
	}
	ent := &Entry{
		Key:   coreLeaderPrefix + uuid,
		Value: val,
	}
	err = c.barrier.Put(ent)
	if err != nil {
		return err
	}

	sd, ok := c.ha.(physical.ServiceDiscovery)
	if ok {
		if err := sd.NotifyActiveStateChange(); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("core: failed to notify active status", "error", err)
			}
		}
	}
	return nil
}

func (c *Core) cleanLeaderPrefix(uuid string, leaderLostCh <-chan struct{}) {
	keys, err := c.barrier.List(coreLeaderPrefix)
	if err != nil {
		c.logger.Error("core: failed to list entries in core/leader", "error", err)
		return
	}
	for len(keys) > 0 {
		select {
		case <-time.After(leaderPrefixCleanDelay):
			if keys[0] != uuid {
				c.barrier.Delete(coreLeaderPrefix + keys[0])
			}
			keys = keys[1:]
		case <-leaderLostCh:
			return
		}
	}
}

// clearLeader is used to clear our leadership entry
func (c *Core) clearLeader(uuid string) error {
	key := coreLeaderPrefix + uuid
	err := c.barrier.Delete(key)

	// Advertise ourselves as a standby
	sd, ok := c.ha.(physical.ServiceDiscovery)
	if ok {
		if err := sd.NotifyActiveStateChange(); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("core: failed to notify standby status", "error", err)
			}
		}
	}

	return err
}

// emitMetrics is used to periodically expose metrics while runnig
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

func (c *Core) SealAccess() *SealAccess {
	sa := &SealAccess{}
	sa.SetSeal(c.seal)
	return sa
}

func (c *Core) Logger() log.Logger {
	return c.logger
}

func (c *Core) BarrierKeyLength() (min, max int) {
	min, max = c.barrier.KeyLength()
	max += shamir.ShareOverhead
	return
}

func (c *Core) ValidateWrappingToken(token string) (bool, error) {
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	te, err := c.tokenStore.Lookup(token)
	if err != nil {
		return false, err
	}
	if te == nil {
		return false, nil
	}

	if len(te.Policies) != 1 {
		return false, nil
	}

	if te.Policies[0] != responseWrappingPolicyName {
		return false, nil
	}

	return true, nil
}
