package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/subtle"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/errwrap"
	log "github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/go-kms-wrapping/wrappers/awskms/v2"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/mlock"
	"github.com/hashicorp/go-secure-stdlib/reloadutil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-secure-stdlib/tlsutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/helper/experiments"
	"github.com/hashicorp/vault/helper/identity/mfa"
	"github.com/hashicorp/vault/helper/locking"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/osutil"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/pathmanager"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	sr "github.com/hashicorp/vault/serviceregistration"
	"github.com/hashicorp/vault/shamir"
	"github.com/hashicorp/vault/vault/cluster"
	"github.com/hashicorp/vault/vault/eventbus"
	"github.com/hashicorp/vault/vault/quotas"
	vaultseal "github.com/hashicorp/vault/vault/seal"
	"github.com/hashicorp/vault/version"
	"github.com/patrickmn/go-cache"
	uberAtomic "go.uber.org/atomic"
	"google.golang.org/grpc"
)

const (
	// CoreLockPath is the path used to acquire a coordinating lock
	// for a highly-available deploy.
	CoreLockPath = "core/lock"

	// The poison pill is used as a check during certain scenarios to indicate
	// to standby nodes that they should seal
	poisonPillPath   = "core/poison-pill"
	poisonPillDRPath = "core/poison-pill-dr"

	// coreLeaderPrefix is the prefix used for the UUID that contains
	// the currently elected leader.
	coreLeaderPrefix = "core/leader/"

	// coreKeyringCanaryPath is used as a canary to indicate to replicated
	// clusters that they need to perform a rekey operation synchronously; this
	// isn't keyring-canary to avoid ignoring it when ignoring core/keyring
	coreKeyringCanaryPath = "core/canary-keyring"

	// coreGroupPolicyApplicationPath is used to store the behaviour for
	// how policies should be applied
	coreGroupPolicyApplicationPath = "core/group-policy-application-mode"

	// groupPolicyApplicationModeWithinNamespaceHierarchy is a configuration option for group
	// policy application modes, which allows only in-namespace-hierarchy policy application
	groupPolicyApplicationModeWithinNamespaceHierarchy = "within_namespace_hierarchy"

	// groupPolicyApplicationModeAny is a configuration option for group
	// policy application modes, which allows policy application irrespective of namespaces
	groupPolicyApplicationModeAny = "any"

	indexHeaderHMACKeyPath = "core/index-header-hmac-key"

	// defaultMFAAuthResponseTTL is the default duration that Vault caches the
	// MfaAuthResponse when the value is not specified in the server config
	defaultMFAAuthResponseTTL = 300 * time.Second

	// defaultMaxTOTPValidateAttempts is the default value for the number
	// of failed attempts to validate a request subject to TOTP MFA. If the
	// number of failed totp passcode validations exceeds this max value, the
	// user needs to wait until a fresh totp passcode is generated.
	defaultMaxTOTPValidateAttempts = 5

	// ForwardSSCTokenToActive is the value that must be set in the
	// forwardToActive to trigger forwarding if a perf standby encounters
	// an SSC Token that it does not have the WAL state for.
	ForwardSSCTokenToActive = "new_token"

	WrapperTypeHsmAutoDeprecated = wrapping.WrapperType("hsm-auto")

	// undoLogsAreSafeStoragePath is a storage path that we write once we know undo logs are
	// safe, so we don't have to keep checking all the time.
	undoLogsAreSafeStoragePath = "core/raft/undo_logs_are_safe"
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

	// ErrIntrospectionNotEnabled is returned if "introspection_endpoint" is not
	// enabled in the configuration file
	ErrIntrospectionNotEnabled = errors.New("The Vault configuration must set \"introspection_endpoint\" to true to enable this endpoint")

	// manualStepDownSleepPeriod is how long to sleep after a user-initiated
	// step down of the active node, to prevent instantly regrabbing the lock.
	// It's var not const so that tests can manipulate it.
	manualStepDownSleepPeriod = 10 * time.Second

	// Functions only in the Enterprise version
	enterprisePostUnseal         = enterprisePostUnsealImpl
	enterprisePreSeal            = enterprisePreSealImpl
	enterpriseSetupFilteredPaths = enterpriseSetupFilteredPathsImpl
	enterpriseSetupQuotas        = enterpriseSetupQuotasImpl
	enterpriseSetupAPILock       = setupAPILockImpl
	startReplication             = startReplicationImpl
	stopReplication              = stopReplicationImpl
	LastWAL                      = lastWALImpl
	LastPerformanceWAL           = lastPerformanceWALImpl
	LastDRWAL                    = lastDRWALImpl
	PerformanceMerkleRoot        = merkleRootImpl
	DRMerkleRoot                 = merkleRootImpl
	LastRemoteWAL                = lastRemoteWALImpl
	LastRemoteUpstreamWAL        = lastRemoteUpstreamWALImpl
	WaitUntilWALShipped          = waitUntilWALShippedImpl
	storedLicenseCheck           = func(c *Core, conf *CoreConfig) error { return nil }
	LicenseAutoloaded            = func(*Core) bool { return false }
	LicenseInitCheck             = func(*Core) error { return nil }
	LicenseSummary               = func(*Core) (*LicenseState, error) { return nil, nil }
	LicenseReload                = func(*Core) error { return nil }
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

// IsFatalError returns true if the given error is a fatal error.
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

type RegisterAuthFunc func(context.Context, time.Duration, string, *logical.Auth, string) error

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

type raftInformation struct {
	challenge           *wrapping.BlobInfo
	leaderClient        *api.Client
	leaderBarrierConfig *SealConfig
	nonVoter            bool
	joinInProgress      bool
}

type migrationInformation struct {
	// seal to use during a migration operation. It is the
	// seal we're migrating *from*.
	seal Seal

	// unsealKey was the unseal key provided for the migration seal.
	// This will be set as the recovery key when migrating from shamir to auto-seal.
	// We don't need to do anything with it when migrating auto->shamir because
	// we don't store the shamir combined key for shamir seals, nor when
	// migrating auto->auto because then the recovery key doesn't change.
	unsealKey []byte
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

	// storageType is the the storage type set in the storage configuration
	storageType string

	// redirectAddr is the address we advertise as leader if held
	redirectAddr string

	// clusterAddr is the address we use for clustering
	clusterAddr *atomic.Value

	// physical backend is the un-trusted backend with durable data
	physical physical.Backend

	// serviceRegistration is the ServiceRegistration network
	serviceRegistration sr.ServiceRegistration

	// hcpLinkStatus is a string describing the status of HCP link connection
	hcpLinkStatus HCPLinkStatus

	// underlyingPhysical will always point to the underlying backend
	// implementation. This is an un-trusted backend with durable data
	underlyingPhysical physical.Backend

	// seal is our seal, for seal configuration information
	seal Seal

	// raftJoinDoneCh is used by the raft retry join routine to inform unseal process
	// that the join is complete
	raftJoinDoneCh chan struct{}

	// postUnsealStarted informs the raft retry join routine that unseal key
	// validation is completed and post unseal has started so that it can complete
	// the join process when Shamir seal is in use
	postUnsealStarted *uint32

	// raftInfo will contain information required for this node to join as a
	// peer to an existing raft cluster. This is marked atomic to prevent data
	// races and casted to raftInformation wherever it is used.
	raftInfo *atomic.Value

	// migrationInfo is used during (and possibly after) a seal migration.
	// This contains information about the seal we are migrating *from*.  Even
	// post seal migration, provided the old seal is still in configuration
	// migrationInfo will be populated, which on enterprise may be necessary for
	// seal rewrap.
	migrationInfo     *migrationInformation
	sealMigrationDone *uint32

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
	stateLock locking.RWMutex
	sealed    *uint32

	standby              bool
	perfStandby          bool
	standbyDoneCh        chan struct{}
	standbyStopCh        *atomic.Value
	manualStepDownCh     chan struct{}
	keepHALockOnStepDown *uint32
	heldHALock           physical.Lock

	// shutdownDoneCh is used to notify when core.Shutdown() completes.
	// core.Shutdown() is typically issued in a goroutine to allow Vault to
	// release the stateLock. This channel is marked atomic to prevent race
	// conditions.
	shutdownDoneCh *atomic.Value

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

	// mountMigrationTracker tracks past and ongoing remount operations
	// against their migration ids
	mountMigrationTracker *sync.Map

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
	systemBackend   *SystemBackend
	loginMFABackend *LoginMFABackend

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

	// activityLog is used to track active client count
	activityLog *ActivityLog
	// activityLogLock protects the activityLog and activityLogConfig
	activityLogLock sync.RWMutex

	// metricsCh is used to stop the metrics streaming
	metricsCh chan struct{}

	// metricsMutex is used to prevent a race condition between
	// metrics emission and sealing leading to a nil pointer
	metricsMutex sync.Mutex

	// inFlightReqMap is used to store info about in-flight requests
	inFlightReqData *InFlightRequests

	// mfaResponseAuthQueue is used to cache the auth response per request ID
	mfaResponseAuthQueue     *LoginMFAPriorityQueue
	mfaResponseAuthQueueLock sync.Mutex

	// metricSink is the destination for all metrics that have
	// a cluster label.
	metricSink *metricsutil.ClusterMetricSink

	defaultLeaseTTL time.Duration
	maxLeaseTTL     time.Duration

	// baseLogger is used to avoid ResetNamed as it strips useful prefixes in
	// e.g. testing
	baseLogger log.Logger
	logger     log.Logger

	// log level provided by config, CLI flag, or env
	logLevel string

	// Disables the trace display for Sentinel checks
	sentinelTraceDisabled bool

	// cachingDisabled indicates whether caches are disabled
	cachingDisabled bool
	// Cache stores the actual cache; we always have this but may bypass it if
	// disabled
	physicalCache physical.ToggleablePurgemonster

	// logRequestsLevel indicates at which level requests should be logged
	logRequestsLevel *uberAtomic.Int32

	// reloadFuncs is a map containing reload functions
	reloadFuncs map[string][]reloadutil.ReloadFunc

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
	// ID
	clusterID uberAtomic.String
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
	// The UUID used to hold the leader lock. Only set on active node
	leaderUUID string

	// CORS Information
	corsConfig *CORSConfig

	// replicationState keeps the current replication state cached for quick
	// lookup; activeNodeReplicationState stores the active value on standbys
	replicationState           *uint32
	activeNodeReplicationState *uint32

	// uiConfig contains UI configuration
	uiConfig *UIConfig

	// rawEnabled indicates whether the Raw endpoint is enabled
	rawEnabled bool

	// inspectableEnabled indicates whether the Inspect endpoint is enabled
	introspectionEnabled     bool
	introspectionEnabledLock sync.Mutex

	// pluginDirectory is the location vault will look for plugin binaries
	pluginDirectory string

	// pluginFileUid is the uid of the plugin files and directory
	pluginFileUid int

	// pluginFilePermissions is the permissions of the plugin files and directory
	pluginFilePermissions int

	// pluginCatalog is used to manage plugin configurations
	pluginCatalog *PluginCatalog

	// The userFailedLoginInfo map has user failed login information.
	// It has user information (alias-name and mount accessor) as a key
	// and login counter, last failed login time as value
	userFailedLoginInfo map[FailedLoginUser]*FailedLoginInfo

	// userFailedLoginInfoLock controls access to the userFailedLoginInfoMap
	userFailedLoginInfoLock sync.RWMutex

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

	// Stores any funcs that should be run on successful barrier unseal in
	// recovery mode
	postRecoveryUnsealFuncs []func() error

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
	clusterListener *atomic.Value

	// customListenerHeader holds custom response headers for a listener
	customListenerHeader *atomic.Value

	// Telemetry objects
	metricsHelper *metricsutil.MetricsHelper

	// raftFollowerStates tracks information about all the raft follower nodes.
	raftFollowerStates *raft.FollowerStates
	// Stop channel for raft TLS rotations
	raftTLSRotationStopCh chan struct{}
	// Stores the pending peers we are waiting to give answers
	pendingRaftPeers *sync.Map

	// rawConfig stores the config as-is from the provided server configuration.
	rawConfig *atomic.Value

	coreNumber int

	// secureRandomReader is the reader used for CSP operations
	secureRandomReader io.Reader

	recoveryMode bool

	clusterNetworkLayer cluster.NetworkLayer

	// PR1103disabled is used to test upgrade workflows: when set to true,
	// the correct behaviour for namespaced cubbyholes is disabled, so we
	// can test an upgrade to a version that includes the fixes from
	// https://github.com/hashicorp/vault-enterprise/pull/1103
	PR1103disabled bool

	quotaManager *quotas.Manager

	clusterHeartbeatInterval time.Duration

	// activityLogConfig contains override values for the activity log
	// it is protected by activityLogLock
	activityLogConfig ActivityLogCoreConfig

	censusConfig atomic.Value

	// activeTime is set on active nodes indicating the time at which this node
	// became active.
	activeTime time.Time

	// KeyRotateGracePeriod is how long we allow an upgrade path
	// for standby instances before we delete the upgrade keys
	keyRotateGracePeriod *int64

	autoRotateCancel context.CancelFunc

	// number of workers to use for lease revocation in the expiration manager
	numExpirationWorkers int

	IndexHeaderHMACKey uberAtomic.Value

	// disableAutopilot is used to disable the autopilot subsystem in raft storage
	disableAutopilot bool

	// enable/disable identifying response headers
	enableResponseHeaderHostname   bool
	enableResponseHeaderRaftNodeID bool

	// disableSSCTokens is used to disable server side consistent token creation/usage
	disableSSCTokens bool

	// versionHistory is a map of vault versions to VaultVersion. The
	// VaultVersion.TimestampInstalled when the version will denote when the version
	// was first run. Note that because perf standbys should be upgraded first, and
	// only the active node will actually write the new version timestamp, a perf
	// standby shouldn't rely on the stored version timestamps being present.
	versionHistory map[string]VaultVersion

	// effectiveSDKVersion contains the SDK version that standby nodes should use when
	// heartbeating with the active node. Default to the current SDK version.
	effectiveSDKVersion string

	rollbackPeriod time.Duration

	experiments []string

	pendingRemovalMountsAllowed bool
	expirationRevokeRetryBase   time.Duration

	events *eventbus.EventBus

	// writeForwardedPaths are a set of storage paths which are GRPC forwarded
	// to the active node of the primary cluster, when present. This PathManager
	// contains absolute paths that we intend to forward (and template) when
	// we're on a secondary cluster.
	writeForwardedPaths *pathmanager.PathManager
}

// c.stateLock needs to be held in read mode before calling this function.
func (c *Core) HAState() consts.HAState {
	switch {
	case c.perfStandby:
		return consts.PerfStandby
	case c.standby:
		return consts.Standby
	default:
		return consts.Active
	}
}

func (c *Core) HAStateWithLock() consts.HAState {
	c.stateLock.RLock()
	c.stateLock.RUnlock()

	return c.HAState()
}

// CoreConfig is used to parameterize a core
type CoreConfig struct {
	entCoreConfig

	DevToken string

	BuiltinRegistry BuiltinRegistry

	LogicalBackends map[string]logical.Factory

	CredentialBackends map[string]logical.Factory

	AuditBackends map[string]audit.Factory

	Physical physical.Backend

	StorageType string

	// May be nil, which disables HA operations
	HAPhysical physical.HABackend

	ServiceRegistration sr.ServiceRegistration

	// Seal is the configured seal, or if none is configured explicitly, a
	// shamir seal.  In migration scenarios this is the new seal.
	Seal Seal

	// Unwrap seal is the optional seal marked "disabled"; this is the old
	// seal in migration scenarios.
	UnwrapSeal Seal

	SecureRandomReader io.Reader

	LogLevel string

	Logger log.Logger

	// Use the deadlocks library to detect deadlocks
	DetectDeadlocks string

	// Disables the trace display for Sentinel checks
	DisableSentinelTrace bool

	// Disables the LRU cache on the physical backend
	DisableCache bool

	// Disables mlock syscall
	DisableMlock bool

	// Custom cache size for the LRU cache on the physical backend, or zero for default
	CacheSize int

	// Set as the leader address for HA
	RedirectAddr string

	// Set as the cluster address for HA
	ClusterAddr string

	DefaultLeaseTTL time.Duration

	MaxLeaseTTL time.Duration

	ClusterName string

	ClusterCipherSuites string

	EnableUI bool

	// Enable the raw endpoint
	EnableRaw bool

	// Enable the introspection endpoint
	EnableIntrospection bool

	PluginDirectory string

	PluginFileUid int

	PluginFilePermissions int

	DisableSealWrap bool

	RawConfig *server.Config

	ReloadFuncs     *map[string][]reloadutil.ReloadFunc
	ReloadFuncsLock *sync.RWMutex

	// Licensing
	License         string
	LicensePath     string
	LicensingConfig *LicensingConfig

	// Configured Census Agent
	CensusAgent CensusReporter

	DisablePerformanceStandby bool
	DisableIndexing           bool
	DisableKeyEncodingChecks  bool

	AllLoggers []log.Logger

	// Telemetry objects
	MetricsHelper *metricsutil.MetricsHelper
	MetricSink    *metricsutil.ClusterMetricSink

	RecoveryMode bool

	ClusterNetworkLayer cluster.NetworkLayer

	ClusterHeartbeatInterval time.Duration

	// Activity log controls
	ActivityLogConfig ActivityLogCoreConfig

	// number of workers to use for lease revocation in the expiration manager
	NumExpirationWorkers int

	// DisableAutopilot is used to disable autopilot subsystem in raft storage
	DisableAutopilot bool

	// Whether to send headers in the HTTP response showing hostname or raft node ID
	EnableResponseHeaderHostname   bool
	EnableResponseHeaderRaftNodeID bool

	// DisableSSCTokens is used to disable the use of server side consistent tokens
	DisableSSCTokens bool

	EffectiveSDKVersion string

	RollbackPeriod time.Duration

	Experiments []string

	PendingRemovalMountsAllowed bool

	ExpirationRevokeRetryBase time.Duration
}

// GetServiceRegistration returns the config's ServiceRegistration, or nil if it does
// not exist.
func (c *CoreConfig) GetServiceRegistration() sr.ServiceRegistration {
	// Check whether there is a ServiceRegistration explicitly configured
	if c.ServiceRegistration != nil {
		return c.ServiceRegistration
	}

	// Check if HAPhysical is configured and implements ServiceRegistration
	if c.HAPhysical != nil && c.HAPhysical.HAEnabled() {
		if disc, ok := c.HAPhysical.(sr.ServiceRegistration); ok {
			return disc
		}
	}

	// No service discovery is available.
	return nil
}

// CreateCore conducts static validations on the Core Config
// and returns an uninitialized core.
func CreateCore(conf *CoreConfig) (*Core, error) {
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
			return nil, fmt.Errorf("redirect address is not valid url: %w", err)
		}

		if u.Scheme == "" {
			return nil, fmt.Errorf("redirect address must include scheme (ex. 'http')")
		}
	}

	// Make a default logger if not provided
	if conf.Logger == nil {
		conf.Logger = logging.NewVaultLogger(log.Trace)
	}

	// Make a default metric sink if not provided
	if conf.MetricSink == nil {
		conf.MetricSink = metricsutil.BlackholeSink()
	}

	// Instantiate a non-nil raw config if none is provided
	if conf.RawConfig == nil {
		conf.RawConfig = new(server.Config)
	}

	// secureRandomReader cannot be nil
	if conf.SecureRandomReader == nil {
		conf.SecureRandomReader = rand.Reader
	}

	clusterHeartbeatInterval := conf.ClusterHeartbeatInterval
	if clusterHeartbeatInterval == 0 {
		clusterHeartbeatInterval = 5 * time.Second
	}

	if conf.NumExpirationWorkers == 0 {
		conf.NumExpirationWorkers = numExpirationWorkersDefault
	}

	// Use imported logging deadlock if requested
	var stateLock locking.RWMutex
	if strings.Contains(conf.DetectDeadlocks, "statelock") {
		stateLock = &locking.DeadlockRWMutex{}
	} else {
		stateLock = &locking.SyncRWMutex{}
	}

	effectiveSDKVersion := conf.EffectiveSDKVersion
	if effectiveSDKVersion == "" {
		effectiveSDKVersion = version.GetVersion().Version
	}

	// Setup the core
	c := &Core{
		entCore:              entCore{},
		devToken:             conf.DevToken,
		physical:             conf.Physical,
		serviceRegistration:  conf.GetServiceRegistration(),
		underlyingPhysical:   conf.Physical,
		storageType:          conf.StorageType,
		redirectAddr:         conf.RedirectAddr,
		clusterAddr:          new(atomic.Value),
		clusterListener:      new(atomic.Value),
		customListenerHeader: new(atomic.Value),
		seal:                 conf.Seal,
		stateLock:            stateLock,
		router:               NewRouter(),
		sealed:               new(uint32),
		sealMigrationDone:    new(uint32),
		standby:              true,
		standbyStopCh:        new(atomic.Value),
		baseLogger:           conf.Logger,
		logger:               conf.Logger.Named("core"),
		logLevel:             conf.LogLevel,

		defaultLeaseTTL:                conf.DefaultLeaseTTL,
		maxLeaseTTL:                    conf.MaxLeaseTTL,
		sentinelTraceDisabled:          conf.DisableSentinelTrace,
		cachingDisabled:                conf.DisableCache,
		clusterName:                    conf.ClusterName,
		clusterNetworkLayer:            conf.ClusterNetworkLayer,
		clusterPeerClusterAddrsCache:   cache.New(3*clusterHeartbeatInterval, time.Second),
		enableMlock:                    !conf.DisableMlock,
		rawEnabled:                     conf.EnableRaw,
		introspectionEnabled:           conf.EnableIntrospection,
		shutdownDoneCh:                 new(atomic.Value),
		replicationState:               new(uint32),
		localClusterPrivateKey:         new(atomic.Value),
		localClusterCert:               new(atomic.Value),
		localClusterParsedCert:         new(atomic.Value),
		activeNodeReplicationState:     new(uint32),
		keepHALockOnStepDown:           new(uint32),
		replicationFailure:             new(uint32),
		disablePerfStandby:             true,
		activeContextCancelFunc:        new(atomic.Value),
		allLoggers:                     conf.AllLoggers,
		builtinRegistry:                conf.BuiltinRegistry,
		neverBecomeActive:              new(uint32),
		clusterLeaderParams:            new(atomic.Value),
		metricsHelper:                  conf.MetricsHelper,
		metricSink:                     conf.MetricSink,
		secureRandomReader:             conf.SecureRandomReader,
		rawConfig:                      new(atomic.Value),
		recoveryMode:                   conf.RecoveryMode,
		postUnsealStarted:              new(uint32),
		raftInfo:                       new(atomic.Value),
		raftJoinDoneCh:                 make(chan struct{}),
		clusterHeartbeatInterval:       clusterHeartbeatInterval,
		activityLogConfig:              conf.ActivityLogConfig,
		keyRotateGracePeriod:           new(int64),
		numExpirationWorkers:           conf.NumExpirationWorkers,
		raftFollowerStates:             raft.NewFollowerStates(),
		disableAutopilot:               conf.DisableAutopilot,
		enableResponseHeaderHostname:   conf.EnableResponseHeaderHostname,
		enableResponseHeaderRaftNodeID: conf.EnableResponseHeaderRaftNodeID,
		mountMigrationTracker:          &sync.Map{},
		disableSSCTokens:               conf.DisableSSCTokens,
		effectiveSDKVersion:            effectiveSDKVersion,
		userFailedLoginInfo:            make(map[FailedLoginUser]*FailedLoginInfo),
		experiments:                    conf.Experiments,
		pendingRemovalMountsAllowed:    conf.PendingRemovalMountsAllowed,
		expirationRevokeRetryBase:      conf.ExpirationRevokeRetryBase,
	}

	c.standbyStopCh.Store(make(chan struct{}))
	atomic.StoreUint32(c.sealed, 1)
	c.metricSink.SetGaugeWithLabels([]string{"core", "unsealed"}, 0, nil)

	c.shutdownDoneCh.Store(make(chan struct{}))

	c.allLoggers = append(c.allLoggers, c.logger)

	c.router.logger = c.logger.Named("router")
	c.allLoggers = append(c.allLoggers, c.router.logger)

	c.inFlightReqData = &InFlightRequests{
		InFlightReqMap:   &sync.Map{},
		InFlightReqCount: uberAtomic.NewUint64(0),
	}

	c.SetConfig(conf.RawConfig)

	atomic.StoreUint32(c.replicationState, uint32(consts.ReplicationDRDisabled|consts.ReplicationPerformanceDisabled))
	c.localClusterCert.Store(([]byte)(nil))
	c.localClusterParsedCert.Store((*x509.Certificate)(nil))
	c.localClusterPrivateKey.Store((*ecdsa.PrivateKey)(nil))

	c.clusterLeaderParams.Store((*ClusterLeaderParams)(nil))
	c.clusterAddr.Store(conf.ClusterAddr)
	c.activeContextCancelFunc.Store((context.CancelFunc)(nil))
	atomic.StoreInt64(c.keyRotateGracePeriod, int64(2*time.Minute))

	c.hcpLinkStatus = HCPLinkStatus{
		lock:             sync.RWMutex{},
		ConnectionStatus: "disconnected",
	}

	c.raftInfo.Store((*raftInformation)(nil))

	switch conf.ClusterCipherSuites {
	case "tls13", "tls12":
		// Do nothing, let Go use the default

	case "":
		// Add in forward compatible TLS 1.3 suites, followed by handpicked 1.2 suites
		c.clusterCipherSuites = []uint16{
			// 1.3
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
			tls.TLS_CHACHA20_POLY1305_SHA256,
			// 1.2
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
		}

	default:
		suites, err := tlsutil.ParseCiphers(conf.ClusterCipherSuites)
		if err != nil {
			return nil, fmt.Errorf("error parsing cluster cipher suites: %w", err)
		}
		c.clusterCipherSuites = suites
	}

	// Load CORS config and provide a value for the core field.
	c.corsConfig = &CORSConfig{
		core:    c,
		Enabled: new(uint32),
	}

	// Load write-forwarded path manager.
	c.writeForwardedPaths = pathmanager.New()

	// Load seal information.
	if c.seal == nil {
		wrapper := aeadwrapper.NewShamirWrapper()
		wrapper.SetConfig(context.Background(), awskms.WithLogger(c.logger.Named("shamir")))

		c.seal = NewDefaultSeal(&vaultseal.Access{
			Wrapper: wrapper,
		})
	}
	c.seal.SetCore(c)
	return c, nil
}

// NewCore is used to construct a new core
func NewCore(conf *CoreConfig) (*Core, error) {
	var err error
	c, err := CreateCore(conf)
	if err != nil {
		return nil, err
	}

	if err = coreInit(c, conf); err != nil {
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

	// Construct a new AES-GCM barrier
	c.barrier, err = NewAESGCMBarrier(c.physical)
	if err != nil {
		return nil, fmt.Errorf("barrier setup failed: %w", err)
	}

	if err := storedLicenseCheck(c, conf); err != nil {
		return nil, err
	}
	// We create the funcs here, then populate the given config with it so that
	// the caller can share state
	conf.ReloadFuncsLock = &c.reloadFuncsLock
	c.reloadFuncsLock.Lock()
	c.reloadFuncs = make(map[string][]reloadutil.ReloadFunc)
	c.reloadFuncsLock.Unlock()
	conf.ReloadFuncs = &c.reloadFuncs

	c.rollbackPeriod = conf.RollbackPeriod
	if conf.RollbackPeriod == 0 {
		c.rollbackPeriod = time.Minute
	}

	// All the things happening below this are not required in
	// recovery mode
	if c.recoveryMode {
		return c, nil
	}

	if conf.PluginDirectory != "" {
		c.pluginDirectory, err = filepath.Abs(conf.PluginDirectory)
		if err != nil {
			return nil, fmt.Errorf("core setup failed, could not verify plugin directory: %w", err)
		}
	}

	if conf.PluginFileUid != 0 {
		c.pluginFileUid = conf.PluginFileUid
	}
	if conf.PluginFilePermissions != 0 {
		c.pluginFilePermissions = conf.PluginFilePermissions
	}

	createSecondaries(c, conf)

	if conf.HAPhysical != nil && conf.HAPhysical.HAEnabled() {
		c.ha = conf.HAPhysical
	}

	c.loginMFABackend = NewLoginMFABackend(c, conf.Logger)

	if c.loginMFABackend.mfaLogger != nil {
		c.AddLogger(c.loginMFABackend.mfaLogger)
	}

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

	c.clusterListener.Store((*cluster.Listener)(nil))

	// for listeners with custom response headers, configuring customListenerHeader
	if conf.RawConfig.Listeners != nil {
		uiHeaders, err := c.UIHeaders()
		if err != nil {
			return nil, err
		}
		c.customListenerHeader.Store(NewListenerCustomHeader(conf.RawConfig.Listeners, c.logger, uiHeaders))
	} else {
		c.customListenerHeader.Store(([]*ListenerCustomHeaders)(nil))
	}

	logRequestsLevel := conf.RawConfig.LogRequestsLevel
	c.logRequestsLevel = uberAtomic.NewInt32(0)
	switch {
	case log.LevelFromString(logRequestsLevel) > log.NoLevel && log.LevelFromString(logRequestsLevel) < log.Off:
		c.logRequestsLevel.Store(int32(log.LevelFromString(logRequestsLevel)))
	case logRequestsLevel != "":
		c.logger.Warn("invalid log_requests_level", "level", conf.RawConfig.LogRequestsLevel)
	}

	quotasLogger := conf.Logger.Named("quotas")
	c.allLoggers = append(c.allLoggers, quotasLogger)
	c.quotaManager, err = quotas.NewManager(quotasLogger, c.quotaLeaseWalker, c.metricSink)
	if err != nil {
		return nil, err
	}

	err = c.adjustForSealMigration(conf.UnwrapSeal)
	if err != nil {
		return nil, err
	}

	if c.versionHistory == nil {
		c.logger.Info("Initializing version history cache for core")
		c.versionHistory = make(map[string]VaultVersion)
	}

	// start the event system
	eventsLogger := conf.Logger.Named("events")
	c.allLoggers = append(c.allLoggers, eventsLogger)
	events, err := eventbus.NewEventBus(eventsLogger)
	if err != nil {
		return nil, err
	}
	c.events = events
	if c.IsExperimentEnabled(experiments.VaultExperimentEventsAlpha1) {
		c.events.Start()
	}

	return c, nil
}

// handleVersionTimeStamps stores the current version at the current time to
// storage, and then loads all versions and upgrade timestamps out from storage.
func (c *Core) handleVersionTimeStamps(ctx context.Context) error {
	currentTime := time.Now().UTC()

	vaultVersion := &VaultVersion{
		TimestampInstalled: currentTime,
		Version:            version.Version,
		BuildDate:          version.BuildDate,
	}

	isUpdated, err := c.storeVersionEntry(ctx, vaultVersion, false)
	if err != nil {
		return fmt.Errorf("error storing vault version: %w", err)
	}
	if isUpdated {
		c.logger.Info("Recorded vault version", "vault version", version.Version, "upgrade time", currentTime, "build date", version.BuildDate)
	}

	// Finally, repopulate the version history cache
	err = c.loadVersionHistory(ctx)
	if err != nil {
		return err
	}
	return nil
}

// HostnameHeaderEnabled determines whether to add the X-Vault-Hostname header
// to HTTP responses.
func (c *Core) HostnameHeaderEnabled() bool {
	return c.enableResponseHeaderHostname
}

// RaftNodeIDHeaderEnabled determines whether to add the X-Vault-Raft-Node-ID header
// to HTTP responses.
func (c *Core) RaftNodeIDHeaderEnabled() bool {
	return c.enableResponseHeaderRaftNodeID
}

// DisableSSCTokens determines whether to use server side consistent tokens or not.
func (c *Core) DisableSSCTokens() bool {
	return c.disableSSCTokens
}

// ShutdownCoreError logs a shutdown error and shuts down the Vault core.
func (c *Core) ShutdownCoreError(err error) {
	c.Logger().Error("shutting down core", "error", err)
	if shutdownErr := c.ShutdownWait(); shutdownErr != nil {
		c.Logger().Error("failed to shutdown core", "error", shutdownErr)
	}
}

// Shutdown is invoked when the Vault instance is about to be terminated. It
// should not be accessible as part of an API call as it will cause an availability
// problem. It is only used to gracefully quit in the case of HA so that failover
// happens as quickly as possible.
func (c *Core) Shutdown() error {
	c.logger.Debug("shutdown called")
	err := c.sealInternal()

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	doneCh := c.shutdownDoneCh.Load().(chan struct{})
	if doneCh != nil {
		close(doneCh)
		c.shutdownDoneCh.Store((chan struct{})(nil))
	}

	return err
}

func (c *Core) ShutdownWait() error {
	donech := c.ShutdownDone()
	err := c.Shutdown()
	if donech != nil {
		<-donech
	}
	return err
}

// ShutdownDone returns a channel that will be closed after Shutdown completes
func (c *Core) ShutdownDone() <-chan struct{} {
	return c.shutdownDoneCh.Load().(chan struct{})
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

func (c *Core) UnsealMigrate(key []byte) (bool, error) {
	err := c.unsealFragment(key, true)
	return !c.Sealed(), err
}

// Unseal is used to provide one of the key parts to unseal the Vault.
func (c *Core) Unseal(key []byte) (bool, error) {
	err := c.unsealFragment(key, false)
	return !c.Sealed(), err
}

// unseal takes a key fragment and attempts to use it to unseal Vault.
// Vault may remain sealed afterwards even when no error is returned,
// depending on whether enough key fragments were provided to meet the
// target threshold.
//
// The provided key should be a recovery key fragment if the seal
// is an autoseal, or a regular seal key fragment for shamir.  In
// migration scenarios "seal" in the preceding sentence refers to
// the migration seal in c.migrationInfo.seal.
//
// We use getUnsealKey to work out if we have enough fragments,
// and if we don't have enough we return early.  Otherwise we get
// back the combined key.
//
// For legacy shamir the combined key *is* the master key.  For
// shamir the combined key is used to decrypt the master key
// read from storage.  For autoseal the combined key isn't used
// except to verify that the stored recovery key matches.
//
// In migration scenarios a side-effect of unsealing is that
// the members of c.migrationInfo are populated (excluding
// .seal, which must already be populated before unseal is called.)
func (c *Core) unsealFragment(key []byte, migrate bool) error {
	defer metrics.MeasureSince([]string{"core", "unseal"}, time.Now())

	c.stateLock.Lock()
	defer c.stateLock.Unlock()

	ctx := context.Background()

	if migrate && c.migrationInfo == nil {
		return fmt.Errorf("can't perform a seal migration, no migration seal found")
	}
	if migrate && c.isRaftUnseal() {
		return fmt.Errorf("can't perform a seal migration while joining a raft cluster")
	}
	if !migrate && c.migrationInfo != nil {
		done, err := c.sealMigrated(ctx)
		if err != nil {
			return fmt.Errorf("error checking to see if seal is migrated: %w", err)
		}
		if !done {
			return fmt.Errorf("migrate option not provided and seal migration is pending")
		}
	}

	c.logger.Debug("unseal key supplied", "migrate", migrate)

	// Explicitly check for init status. This also checks if the seal
	// configuration is valid (i.e. non-nil).
	init, err := c.Initialized(ctx)
	if err != nil {
		return err
	}
	if !init && !c.isRaftUnseal() {
		return ErrNotInit
	}

	// Verify the key length
	min, max := c.barrier.KeyLength()
	max += shamir.ShareOverhead
	if len(key) < min {
		return &ErrInvalidKey{fmt.Sprintf("key is shorter than minimum %d bytes", min)}
	}
	if len(key) > max {
		return &ErrInvalidKey{fmt.Sprintf("key is longer than maximum %d bytes", max)}
	}

	// Check if already unsealed
	if !c.Sealed() {
		return nil
	}

	sealToUse := c.seal
	if migrate {
		c.logger.Info("unsealing using migration seal")
		sealToUse = c.migrationInfo.seal
	}

	newKey, err := c.recordUnsealPart(key)
	if !newKey || err != nil {
		return err
	}

	// getUnsealKey returns either a recovery key (in the case of an autoseal)
	// or a master key (legacy shamir) or an unseal key (new-style shamir).
	combinedKey, err := c.getUnsealKey(ctx, sealToUse)
	if err != nil || combinedKey == nil {
		return err
	}
	if migrate {
		c.migrationInfo.unsealKey = combinedKey
	}

	if c.isRaftUnseal() {
		return c.unsealWithRaft(combinedKey)
	}
	masterKey, err := c.unsealKeyToMasterKeyPreUnseal(ctx, sealToUse, combinedKey)
	if err != nil {
		return err
	}
	return c.unsealInternal(ctx, masterKey)
}

func (c *Core) unsealWithRaft(combinedKey []byte) error {
	ctx := context.Background()

	if c.seal.BarrierType() == wrapping.WrapperTypeShamir {
		// If this is a legacy shamir seal this serves no purpose but it
		// doesn't hurt.
		err := c.seal.GetAccess().Wrapper.(*aeadwrapper.ShamirWrapper).SetAesGcmKeyBytes(combinedKey)
		if err != nil {
			return err
		}
	}

	raftInfo := c.raftInfo.Load().(*raftInformation)

	switch raftInfo.joinInProgress {
	case true:
		// JoinRaftCluster is already trying to perform a join based on retry_join configuration.
		// Inform that routine that unseal key validation is complete so that it can continue to
		// try and join possible leader nodes, and wait for it to complete.

		atomic.StoreUint32(c.postUnsealStarted, 1)

		c.logger.Info("waiting for raft retry join process to complete")
		<-c.raftJoinDoneCh

	default:
		// This is the case for manual raft join. Send the answer to the leader node and
		// wait for data to start streaming in.
		if err := c.joinRaftSendAnswer(ctx, c.seal.GetAccess(), raftInfo); err != nil {
			return err
		}
		// Reset the state
		c.raftInfo.Store((*raftInformation)(nil))
	}

	go func() {
		var masterKey []byte
		keyringFound := false

		// Wait until we at least have the keyring before we attempt to
		// unseal the node.
		for {
			if !keyringFound {
				keys, err := c.underlyingPhysical.List(ctx, keyringPrefix)
				if err != nil {
					c.logger.Error("failed to list physical keys", "error", err)
					return
				}
				if strutil.StrListContains(keys, "keyring") {
					keyringFound = true
				}
			}
			if keyringFound && len(masterKey) == 0 {
				var err error
				masterKey, err = c.unsealKeyToMasterKeyPreUnseal(ctx, c.seal, combinedKey)
				if err != nil {
					c.logger.Error("failed to read master key", "error", err)
					return
				}
			}
			if keyringFound && len(masterKey) > 0 {
				err := c.unsealInternal(ctx, masterKey)
				if err != nil {
					c.logger.Error("failed to unseal", "error", err)
				}
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return nil
}

// recordUnsealPart takes in a key fragment, and returns true if it's a new fragment.
func (c *Core) recordUnsealPart(key []byte) (bool, error) {
	// Check if we already have this piece
	if c.unlockInfo != nil {
		for _, existing := range c.unlockInfo.Parts {
			if subtle.ConstantTimeCompare(existing, key) == 1 {
				return false, nil
			}
		}
	} else {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return false, err
		}
		c.unlockInfo = &unlockInformation{
			Nonce: uuid,
		}
	}

	// Store this key
	c.unlockInfo.Parts = append(c.unlockInfo.Parts, key)
	return true, nil
}

// getUnsealKey uses key fragments recorded by recordUnsealPart and
// returns the combined key if the key share threshold is met.
// If the key fragments are part of a recovery key, also verify that
// it matches the stored recovery key on disk.
func (c *Core) getUnsealKey(ctx context.Context, seal Seal) ([]byte, error) {
	var config *SealConfig
	var err error

	raftInfo := c.raftInfo.Load().(*raftInformation)

	switch {
	case seal.RecoveryKeySupported():
		config, err = seal.RecoveryConfig(ctx)
	case c.isRaftUnseal():
		// Ignore follower's seal config and refer to leader's barrier
		// configuration.
		config = raftInfo.leaderBarrierConfig
	default:
		config, err = seal.BarrierConfig(ctx)
	}
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, fmt.Errorf("failed to obtain seal/recovery configuration")
	}

	// Check if we don't have enough keys to unlock, proceed through the rest of
	// the call only if we have met the threshold
	if len(c.unlockInfo.Parts) < config.SecretThreshold {
		if c.logger.IsDebug() {
			c.logger.Debug("cannot unseal, not enough keys", "keys", len(c.unlockInfo.Parts), "threshold", config.SecretThreshold, "nonce", c.unlockInfo.Nonce)
		}
		return nil, nil
	}

	defer func() {
		c.unlockInfo = nil
	}()

	// Recover the split key. recoveredKey is the shamir combined
	// key, or the single provided key if the threshold is 1.
	var unsealKey []byte
	if config.SecretThreshold == 1 {
		unsealKey = make([]byte, len(c.unlockInfo.Parts[0]))
		copy(unsealKey, c.unlockInfo.Parts[0])
	} else {
		unsealKey, err = shamir.Combine(c.unlockInfo.Parts)
		if err != nil {
			return nil, &ErrInvalidKey{fmt.Sprintf("failed to compute combined key: %v", err)}
		}
	}

	if seal.RecoveryKeySupported() {
		if err := seal.VerifyRecoveryKey(ctx, unsealKey); err != nil {
			return nil, &ErrInvalidKey{fmt.Sprintf("failed to verify recovery key: %v", err)}
		}
	}

	return unsealKey, nil
}

// sealMigrated must be called with the stateLock held.  It returns true if
// the seal configured in HCL and the seal configured in storage match.
// For the auto->auto same seal migration scenario, it will return false even
// if the preceding conditions are true but we cannot decrypt the master key
// in storage using the configured seal.
func (c *Core) sealMigrated(ctx context.Context) (bool, error) {
	if atomic.LoadUint32(c.sealMigrationDone) == 1 {
		return true, nil
	}

	existBarrierSealConfig, existRecoverySealConfig, err := c.PhysicalSealConfigs(ctx)
	if err != nil {
		return false, err
	}

	if existBarrierSealConfig.Type != c.seal.BarrierType().String() {
		return false, nil
	}
	if c.seal.RecoveryKeySupported() && existRecoverySealConfig.Type != c.seal.RecoveryType() {
		return false, nil
	}

	if c.seal.BarrierType() != c.migrationInfo.seal.BarrierType() {
		return true, nil
	}

	// The above checks can handle the auto->shamir and shamir->auto
	// and auto1->auto2 cases.  For auto1->auto1, we need to actually try
	// to read and decrypt the keys.

	keysMig, errMig := c.migrationInfo.seal.GetStoredKeys(ctx)
	keys, err := c.seal.GetStoredKeys(ctx)

	switch {
	case len(keys) > 0 && err == nil:
		return true, nil
	case len(keysMig) > 0 && errMig == nil:
		return false, nil
	case errors.Is(err, &ErrDecrypt{}) && errors.Is(errMig, &ErrDecrypt{}):
		return false, fmt.Errorf("decrypt error, neither the old nor new seal can read stored keys: old seal err=%v, new seal err=%v", errMig, err)
	default:
		return false, fmt.Errorf("neither the old nor new seal can read stored keys: old seal err=%v, new seal err=%v", errMig, err)
	}
}

// migrateSeal must be called with the stateLock held.
func (c *Core) migrateSeal(ctx context.Context) error {
	if c.migrationInfo == nil {
		return nil
	}

	ok, err := c.sealMigrated(ctx)
	if err != nil {
		return fmt.Errorf("error checking if seal is migrated or not: %w", err)
	}

	if ok {
		c.logger.Info("migration is already performed")
		return nil
	}

	c.logger.Info("seal migration initiated")

	switch {
	case c.migrationInfo.seal.RecoveryKeySupported() && c.seal.RecoveryKeySupported():
		c.logger.Info("migrating from one auto-unseal to another", "from",
			c.migrationInfo.seal.BarrierType(), "to", c.seal.BarrierType())

		// Set the recovery and barrier keys to be the same.
		recoveryKey, err := c.migrationInfo.seal.RecoveryKey(ctx)
		if err != nil {
			return fmt.Errorf("error getting recovery key to set on new seal: %w", err)
		}

		if err := c.seal.SetRecoveryKey(ctx, recoveryKey); err != nil {
			return fmt.Errorf("error setting new recovery key information during migrate: %w", err)
		}

		barrierKeys, err := c.migrationInfo.seal.GetStoredKeys(ctx)
		if err != nil {
			return fmt.Errorf("error getting stored keys to set on new seal: %w", err)
		}

		if err := c.seal.SetStoredKeys(ctx, barrierKeys); err != nil {
			return fmt.Errorf("error setting new barrier key information during migrate: %w", err)
		}

	case c.migrationInfo.seal.RecoveryKeySupported():
		c.logger.Info("migrating from one auto-unseal to shamir", "from", c.migrationInfo.seal.BarrierType())
		// Auto to Shamir, since recovery key isn't supported on new seal

		recoveryKey, err := c.migrationInfo.seal.RecoveryKey(ctx)
		if err != nil {
			return fmt.Errorf("error getting recovery key to set on new seal: %w", err)
		}

		// We have recovery keys; we're going to use them as the new shamir KeK.
		err = c.seal.GetAccess().Wrapper.(*aeadwrapper.ShamirWrapper).SetAesGcmKeyBytes(recoveryKey)
		if err != nil {
			return fmt.Errorf("failed to set master key in seal: %w", err)
		}

		barrierKeys, err := c.migrationInfo.seal.GetStoredKeys(ctx)
		if err != nil {
			return fmt.Errorf("error getting stored keys to set on new seal: %w", err)
		}

		if err := c.seal.SetStoredKeys(ctx, barrierKeys); err != nil {
			return fmt.Errorf("error setting new barrier key information during migrate: %w", err)
		}

	case c.seal.RecoveryKeySupported():
		c.logger.Info("migrating from shamir to auto-unseal", "to", c.seal.BarrierType())
		// Migration is happening from shamir -> auto. In this case use the shamir
		// combined key that was used to store the master key as the new recovery key.
		if err := c.seal.SetRecoveryKey(ctx, c.migrationInfo.unsealKey); err != nil {
			return fmt.Errorf("error setting new recovery key information: %w", err)
		}

		// Generate a new master key
		newMasterKey, err := c.barrier.GenerateKey(c.secureRandomReader)
		if err != nil {
			return fmt.Errorf("error generating new master key: %w", err)
		}

		// Rekey the barrier.  This handles the case where the shamir seal we're
		// migrating from was a legacy seal without a stored master key.
		if err := c.barrier.Rekey(ctx, newMasterKey); err != nil {
			return fmt.Errorf("error rekeying barrier during migration: %w", err)
		}

		// Store the new master key
		if err := c.seal.SetStoredKeys(ctx, [][]byte{newMasterKey}); err != nil {
			return fmt.Errorf("error storing new master key: %w", err)
		}

	default:
		return errors.New("unhandled migration case (shamir to shamir)")
	}

	err = c.migrateSealConfig(ctx)
	if err != nil {
		return fmt.Errorf("error storing new seal configs: %w", err)
	}

	// Flag migration performed for seal-rewrap later
	atomic.StoreUint32(c.sealMigrationDone, 1)

	c.logger.Info("seal migration complete")
	return nil
}

// unsealInternal takes in the master key and attempts to unseal the barrier.
// N.B.: This must be called with the state write lock held.
func (c *Core) unsealInternal(ctx context.Context, masterKey []byte) error {
	// Attempt to unlock
	if err := c.barrier.Unseal(ctx, masterKey); err != nil {
		return err
	}

	if err := preUnsealInternal(ctx, c); err != nil {
		return err
	}

	if err := c.startClusterListener(ctx); err != nil {
		return err
	}

	if err := c.startRaftBackend(ctx); err != nil {
		return err
	}

	if err := c.setupReplicationResolverHandler(); err != nil {
		c.logger.Warn("failed to start replication resolver server", "error", err)
	}

	// Do post-unseal setup if HA is not enabled
	if c.ha == nil {
		// We still need to set up cluster info even if it's not part of a
		// cluster right now. This also populates the cached cluster object.
		if err := c.setupCluster(ctx); err != nil {
			c.logger.Error("cluster setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			return err
		}

		if err := c.migrateSeal(ctx); err != nil {
			c.logger.Error("seal migration error", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			return err
		}

		ctx, ctxCancel := context.WithCancel(namespace.RootContext(nil))
		if err := c.postUnseal(ctx, ctxCancel, standardUnsealStrategy{}); err != nil {
			c.logger.Error("post-unseal setup failed", "error", err)
			c.barrier.Seal()
			c.logger.Warn("vault is sealed")
			return err
		}

		// Force a cache bust here, which will also run migration code
		if c.seal.RecoveryKeySupported() {
			c.seal.SetRecoveryConfig(ctx, nil)
		}

		c.standby = false
	} else {
		// Go to standby mode, wait until we are active to unseal
		c.standbyDoneCh = make(chan struct{})
		c.manualStepDownCh = make(chan struct{}, 1)
		c.standbyStopCh.Store(make(chan struct{}))
		go c.runStandby(c.standbyDoneCh, c.manualStepDownCh, c.standbyStopCh.Load().(chan struct{}))
	}

	// Success!
	atomic.StoreUint32(c.sealed, 0)
	c.metricSink.SetGaugeWithLabels([]string{"core", "unsealed"}, 1, nil)

	if c.logger.IsInfo() {
		c.logger.Info("vault is unsealed")
	}

	if c.serviceRegistration != nil {
		if err := c.serviceRegistration.NotifySealedStateChange(false); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("failed to notify unsealed status", "error", err)
			}
		}
		if err := c.serviceRegistration.NotifyInitializedStateChange(true); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("failed to notify initialized status", "error", err)
			}
		}
	}
	return nil
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

	var unlocked bool
	defer func() {
		if !unlocked {
			c.stateLock.RUnlock()
		}
	}()

	if req == nil {
		return errors.New("nil request to seal")
	}

	// Since there is no token store in standby nodes, sealing cannot be done.
	// Ideally, the request has to be forwarded to leader node for validation
	// and the operation should be performed. But for now, just returning with
	// an error and recommending a vault restart, which essentially does the
	// same thing.
	if c.standby {
		c.logger.Error("vault cannot seal when in standby mode; please restart instead")
		return errors.New("vault cannot seal when in standby mode; please restart instead")
	}

	err := c.PopulateTokenEntry(ctx, req)
	if err != nil {
		if errwrap.Contains(err, logical.ErrPermissionDenied.Error()) {
			return logical.ErrPermissionDenied
		}
		return logical.ErrInvalidRequest
	}
	acl, te, entity, identityPolicies, err := c.fetchACLTokenEntryAndEntity(ctx, req)
	if err != nil {
		return err
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

	logInput := &logical.LogInput{
		Auth:    auth,
		Request: req,
	}
	if err := c.auditBroker.LogRequest(ctx, logInput, c.auditedHeaders); err != nil {
		c.logger.Error("failed to audit request", "request_path", req.Path, "error", err)
		return errors.New("failed to audit request, cannot continue")
	}

	if entity != nil && entity.Disabled {
		c.logger.Warn("permission denied as the entity on the token is disabled")
		return logical.ErrPermissionDenied
	}
	if te != nil && te.EntityID != "" && entity == nil {
		c.logger.Warn("permission denied as the entity on the token is invalid")
		return logical.ErrPermissionDenied
	}

	// Attempt to use the token (decrement num_uses)
	// On error bail out; if the token has been revoked, bail out too
	if te != nil {
		te, err = c.tokenStore.UseToken(ctx, te)
		if err != nil {
			c.logger.Error("failed to use token", "error", err)
			return ErrInternalError
		}
		if te == nil {
			// Token is no longer valid
			return logical.ErrPermissionDenied
		}
	}

	// Verify that this operation is allowed
	authResults := c.performPolicyChecks(ctx, acl, te, req, entity, &PolicyCheckOpts{
		RootPrivsRequired: true,
	})
	if !authResults.Allowed {
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
	unlocked = true
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
	return c.sealInternalWithOptions(true, false, true)
}

func (c *Core) sealInternalWithOptions(grabStateLock, keepHALock, performCleanup bool) error {
	// Mark sealed, and if already marked return
	if swapped := atomic.CompareAndSwapUint32(c.sealed, 0, 1); !swapped {
		return nil
	}
	c.metricSink.SetGaugeWithLabels([]string{"core", "unsealed"}, 0, nil)

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
		close(c.standbyStopCh.Load().(chan struct{}))
		c.logger.Debug("finished triggering standbyStopCh for runStandby")

		// Wait for runStandby to stop
		<-c.standbyDoneCh
		atomic.StoreUint32(c.keepHALockOnStepDown, 0)
		c.logger.Debug("runStandby done")
	}

	c.teardownReplicationResolverHandler()

	// Perform additional cleanup upon sealing.
	if performCleanup {
		if raftBackend := c.getRaftBackend(); raftBackend != nil {
			if err := raftBackend.TeardownCluster(c.getClusterListener()); err != nil {
				c.logger.Error("error stopping storage cluster", "error", err)
				return err
			}
		}

		// Stop the cluster listener
		c.stopClusterListener()
	}

	c.logger.Debug("sealing barrier")
	if err := c.barrier.Seal(); err != nil {
		c.logger.Error("error sealing barrier", "error", err)
		return err
	}

	if c.serviceRegistration != nil {
		if err := c.serviceRegistration.NotifySealedStateChange(true); err != nil {
			if c.logger.IsWarn() {
				c.logger.Warn("failed to notify sealed status", "error", err)
			}
		}
	}

	if c.quotaManager != nil {
		if err := c.quotaManager.Reset(); err != nil {
			c.logger.Error("error resetting quota manager", "error", err)
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

	// Mark the active time. We do this first so it can be correlated to the logs
	// for the active startup.
	c.activeTime = time.Now().UTC()

	if err := postUnsealPhysical(c); err != nil {
		return err
	}

	if err := enterprisePostUnseal(c, false); err != nil {
		return err
	}
	if !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary | consts.ReplicationDRSecondary) {
		// Only perf primarys should write feature flags, but we do it by
		// excluding other states so that we don't have to change it when
		// a non-replicated cluster becomes a primary.
		if err := c.persistFeatureFlags(ctx); err != nil {
			return err
		}
	}

	if c.autoRotateCancel == nil {
		var autoRotateCtx context.Context
		autoRotateCtx, c.autoRotateCancel = context.WithCancel(c.activeContext)
		go c.autoRotateBarrierLoop(autoRotateCtx)
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
	if err := enterpriseSetupFilteredPaths(c); err != nil {
		return err
	}
	if err := c.setupMounts(ctx); err != nil {
		return err
	}
	if err := enterpriseSetupAPILock(c, ctx); err != nil {
		return err
	}
	if err := c.setupPolicyStore(ctx); err != nil {
		return err
	}
	if err := c.setupManagedKeyRegistry(); err != nil {
		return err
	}
	if err := c.loadCORSConfig(ctx); err != nil {
		return err
	}
	if err := c.loadCredentials(ctx); err != nil {
		return err
	}
	if err := enterpriseSetupFilteredPaths(c); err != nil {
		return err
	}
	if err := c.setupCredentials(ctx); err != nil {
		return err
	}
	if err := c.setupQuotas(ctx, false); err != nil {
		return err
	}
	if err := c.setupHeaderHMACKey(ctx, false); err != nil {
		return err
	}
	if err := c.runLockedUserEntryUpdates(ctx); err != nil {
		return err
	}
	c.updateLockedUserEntries()

	if !c.IsDRSecondary() {
		if err := c.startRollback(); err != nil {
			return err
		}
		if err := c.setupExpiration(expireLeaseStrategyFairsharing); err != nil {
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
		if err := loadPolicyMFAConfigs(ctx, c); err != nil {
			return err
		}
		c.setupCachedMFAResponseAuth()
		if err := c.loadLoginMFAConfigs(ctx); err != nil {
			return err
		}

		if err := c.setupAuditedHeadersConfig(ctx); err != nil {
			return err
		}

		if err := c.setupCensusAgent(); err != nil {
			c.logger.Error("skipping reporting for nil agent", "error", err)
		}

		// not waiting on wg to avoid changing existing behavior
		var wg sync.WaitGroup
		if err := c.setupActivityLog(ctx, &wg); err != nil {
			return err
		}
	} else {
		c.auditBroker = NewAuditBroker(c.logger)
	}

	if !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary | consts.ReplicationDRSecondary) {
		// Cannot do this above, as we need other resources like mounts to be setup
		if err := c.setupPluginReload(); err != nil {
			return err
		}
	}

	if c.getClusterListener() != nil && (c.ha != nil || shouldStartClusterListener(c)) {
		if err := c.setupRaftActiveNode(ctx); err != nil {
			return err
		}

		if err := c.startForwarding(ctx); err != nil {
			return err
		}

	}

	c.clusterParamsLock.Lock()
	defer c.clusterParamsLock.Unlock()
	if err := startReplication(c); err != nil {
		return err
	}

	c.metricsCh = make(chan struct{})
	go c.emitMetricsActiveNode(c.metricsCh)

	// Establish version timestamps at the end of unseal on active nodes only.
	if err := c.handleVersionTimeStamps(ctx); err != nil {
		return err
	}

	return nil
}

// postUnseal is invoked on the active node, and performance standby nodes,
// after the barrier is unsealed, but before
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
			_ = c.preSeal()
		}
	}()
	c.logger.Info("post-unseal setup starting")

	// Enable the cache
	c.physicalCache.Purge(ctx)
	if !c.cachingDisabled {
		c.physicalCache.SetEnabled(true)
	}

	// Purge these for safety in case of a rekey
	_ = c.seal.SetBarrierConfig(ctx, nil)
	if c.seal.RecoveryKeySupported() {
		_ = c.seal.SetRecoveryConfig(ctx, nil)
	}

	// Load prior un-updated store into version history cache to compare
	// previous state.
	if err := c.loadVersionHistory(ctx); err != nil {
		return err
	}

	if err := unsealer.unseal(ctx, c.logger, c); err != nil {
		return err
	}

	// Automatically re-encrypt the keys used for auto unsealing when the
	// seal's encryption key changes. The regular rotation of cryptographic
	// keys is a NIST recommendation. Access to prior keys for decryption
	// is normally supported for a configurable time period. Re-encrypting
	// the keys used for auto unsealing ensures Vault and its data will
	// continue to be accessible even after prior seal keys are destroyed.
	if seal, ok := c.seal.(*autoSeal); ok {
		if err := seal.UpgradeKeys(c.activeContext); err != nil {
			c.logger.Warn("post-unseal upgrade seal keys failed", "error", err)
		}

		// Start a periodic but infrequent heartbeat to detect auto-seal backend outages at runtime rather than being
		// surprised by this at the next need to unseal.
		seal.StartHealthCheck()
	}

	// This is intentionally the last block in this function. We want to allow
	// writes just before allowing client requests, to ensure everything has
	// been set up properly before any writes can have happened.
	//
	// Use a small temporary worker pool to run postUnsealFuncs in parallel
	postUnsealFuncConcurrency := runtime.NumCPU() * 2
	if v := os.Getenv("VAULT_POSTUNSEAL_FUNC_CONCURRENCY"); v != "" {
		pv, err := strconv.Atoi(v)
		if err != nil || pv < 1 {
			c.logger.Warn("invalid value for VAULT_POSTUNSEAL_FUNC_CURRENCY, must be a positive integer", "error", err, "value", pv)
		} else {
			postUnsealFuncConcurrency = pv
		}
	}
	if postUnsealFuncConcurrency <= 1 {
		// Out of paranoia, keep the old logic for parallism=1
		for _, v := range c.postUnsealFuncs {
			v()
		}
	} else {
		jobs := make(chan func())
		var wg sync.WaitGroup
		for i := 0; i < postUnsealFuncConcurrency; i++ {
			go func() {
				for v := range jobs {
					v()
					wg.Done()
				}
			}()
		}
		for _, v := range c.postUnsealFuncs {
			wg.Add(1)
			jobs <- v
		}
		wg.Wait()
		close(jobs)
	}

	if atomic.LoadUint32(c.sealMigrationDone) == 1 {
		if err := c.postSealMigration(ctx); err != nil {
			c.logger.Warn("post-unseal post seal migration failed", "error", err)
		}
	}

	if os.Getenv(EnvVaultDisableLocalAuthMountEntities) != "" {
		c.logger.Warn("disabling entities for local auth mounts through env var", "env", EnvVaultDisableLocalAuthMountEntities)
	}
	c.loginMFABackend.usedCodes = cache.New(0, 30*time.Second)
	if c.systemBackend != nil && c.systemBackend.mfaBackend != nil {
		c.systemBackend.mfaBackend.usedCodes = cache.New(0, 30*time.Second)
	}
	c.logger.Info("post-unseal setup complete")
	return nil
}

// preSeal is invoked before the barrier is sealed, allowing
// for any state teardown required.
func (c *Core) preSeal() error {
	defer metrics.MeasureSince([]string{"core", "pre_seal"}, time.Now())
	c.logger.Info("pre-seal teardown starting")

	if seal, ok := c.seal.(*autoSeal); ok {
		seal.StopHealthCheck()
	}
	// Clear any pending funcs
	c.postUnsealFuncs = nil
	c.activeTime = time.Time{}

	// Clear any rekey progress
	c.barrierRekeyConfig = nil
	c.recoveryRekeyConfig = nil

	if c.metricsCh != nil {
		close(c.metricsCh)
		c.metricsCh = nil
	}
	var result error

	c.stopForwarding()

	c.stopRaftActiveNode()

	c.clusterParamsLock.Lock()
	if err := stopReplication(c); err != nil {
		result = multierror.Append(result, fmt.Errorf("error stopping replication: %w", err))
	}
	c.clusterParamsLock.Unlock()

	if err := c.teardownAudits(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error tearing down audits: %w", err))
	}
	if err := c.stopExpiration(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error stopping expiration: %w", err))
	}
	c.stopActivityLog()
	// Clean up the censusAgent on seal
	if err := c.teardownCensusAgent(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error tearing down reporting agent: %w", err))
	}

	if err := c.teardownCredentials(context.Background()); err != nil {
		result = multierror.Append(result, fmt.Errorf("error tearing down credentials: %w", err))
	}
	if err := c.teardownPolicyStore(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error tearing down policy store: %w", err))
	}
	if err := c.stopRollback(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error stopping rollback: %w", err))
	}
	if err := c.unloadMounts(context.Background()); err != nil {
		result = multierror.Append(result, fmt.Errorf("error unloading mounts: %w", err))
	}

	if err := enterprisePreSeal(c); err != nil {
		result = multierror.Append(result, err)
	}

	if c.autoRotateCancel != nil {
		c.autoRotateCancel()
		c.autoRotateCancel = nil
	}

	if seal, ok := c.seal.(*autoSeal); ok {
		seal.StopHealthCheck()
	}

	if c.systemBackend != nil && c.systemBackend.mfaBackend != nil {
		c.systemBackend.mfaBackend.usedCodes = nil
	}
	if err := c.teardownLoginMFA(); err != nil {
		result = multierror.Append(result, fmt.Errorf("error tearing down login MFA, error: %w", err))
	}

	preSealPhysical(c)

	c.logger.Info("pre-seal teardown complete")
	return result
}

func enterprisePostUnsealImpl(c *Core, isStandby bool) error {
	return nil
}

func enterprisePreSealImpl(c *Core) error {
	return nil
}

func enterpriseSetupFilteredPathsImpl(c *Core) error {
	return nil
}

func enterpriseSetupQuotasImpl(ctx context.Context, c *Core) error {
	return nil
}

func startReplicationImpl(c *Core) error {
	return nil
}

func stopReplicationImpl(c *Core) error {
	return nil
}

func setupAPILockImpl(_ *Core, _ context.Context) error { return nil }

func (c *Core) ReplicationState() consts.ReplicationState {
	return consts.ReplicationState(atomic.LoadUint32(c.replicationState))
}

func (c *Core) ActiveNodeReplicationState() consts.ReplicationState {
	return consts.ReplicationState(atomic.LoadUint32(c.activeNodeReplicationState))
}

func (c *Core) SealAccess() *SealAccess {
	return NewSealAccess(c.seal)
}

// StorageType returns a string equal to the storage configuration's type.
func (c *Core) StorageType() string {
	return c.storageType
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

func merkleRootImpl(c *Core) string {
	return ""
}

func lastWALImpl(c *Core) uint64 {
	return 0
}

func lastPerformanceWALImpl(c *Core) uint64 {
	return 0
}

func lastDRWALImpl(c *Core) uint64 {
	return 0
}

func lastRemoteWALImpl(c *Core) uint64 {
	return 0
}

func lastRemoteUpstreamWALImpl(c *Core) uint64 {
	return 0
}

func (c *Core) PhysicalSealConfigs(ctx context.Context) (*SealConfig, *SealConfig, error) {
	pe, err := c.physical.Get(ctx, barrierSealConfigPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch barrier seal configuration at migration check time: %w", err)
	}
	if pe == nil {
		return nil, nil, nil
	}

	barrierConf := new(SealConfig)

	if err := jsonutil.DecodeJSON(pe.Value, barrierConf); err != nil {
		return nil, nil, fmt.Errorf("failed to decode barrier seal configuration at migration check time: %w", err)
	}
	err = barrierConf.Validate()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to validate barrier seal configuration at migration check time: %w", err)
	}
	// In older versions of vault the default seal would not store a type. This
	// is here to offer backwards compatibility for older seal configs.
	if barrierConf.Type == "" {
		barrierConf.Type = wrapping.WrapperTypeShamir.String()
	}

	var recoveryConf *SealConfig
	pe, err = c.physical.Get(ctx, recoverySealConfigPlaintextPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch seal configuration at migration check time: %w", err)
	}
	if pe != nil {
		recoveryConf = &SealConfig{}
		if err := jsonutil.DecodeJSON(pe.Value, recoveryConf); err != nil {
			return nil, nil, fmt.Errorf("failed to decode seal configuration at migration check time: %w", err)
		}
		err = recoveryConf.Validate()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to validate seal configuration at migration check time: %w", err)
		}
		// In older versions of vault the default seal would not store a type. This
		// is here to offer backwards compatibility for older seal configs.
		if recoveryConf.Type == "" {
			recoveryConf.Type = wrapping.WrapperTypeShamir.String()
		}
	}

	return barrierConf, recoveryConf, nil
}

// adjustForSealMigration takes the unwrapSeal, which is nil if (a) we're not
// configured for seal migration or (b) we might be doing a seal migration away
// from shamir.  It will only be non-nil if there is a configured seal with
// the config key disabled=true, which implies a migration away from autoseal.
//
// For case (a), the common case, we expect that the stored barrier
// config matches the seal type, in which case we simply return nil.  If they
// don't match, and the stored seal config is of type Shamir but the configured
// seal is not Shamir, that is case (b) and we make an unwrapSeal of type Shamir.
// Any other unwrapSeal=nil scenario is treated as an error.
//
// Given a non-nil unwrapSeal or case (b), we setup c.migrationInfo to prepare
// for a migration upon receiving a valid migration unseal request.  We cannot
// check at this time for already performed (or incomplete) migrations because
// we haven't yet been unsealed, so we have no way of checking whether a
// shamir seal works to read stored seal-encrypted data.
//
// The assumption throughout is that the very last step of seal migration is
// to write the new barrier/recovery stored seal config.
func (c *Core) adjustForSealMigration(unwrapSeal Seal) error {
	ctx := context.Background()
	existBarrierSealConfig, existRecoverySealConfig, err := c.PhysicalSealConfigs(ctx)
	if err != nil {
		return fmt.Errorf("Error checking for existing seal: %s", err)
	}

	// If we don't have an existing config or if it's the deprecated auto seal
	// which needs an upgrade, skip out
	if existBarrierSealConfig == nil || existBarrierSealConfig.Type == WrapperTypeHsmAutoDeprecated.String() {
		return nil
	}

	if unwrapSeal == nil {
		// With unwrapSeal==nil, either we're not migrating, or we're migrating
		// from shamir.

		switch {
		case existBarrierSealConfig.Type == c.seal.BarrierType().String():
			// We have the same barrier type and the unwrap seal is nil so we're not
			// migrating from same to same, IOW we assume it's not a migration.
			return nil
		case c.seal.BarrierType() == wrapping.WrapperTypeShamir:
			// The stored barrier config is not shamir, there is no disabled seal
			// in config, and either no configured seal (which equates to Shamir)
			// or an explicitly configured Shamir seal.
			return fmt.Errorf("cannot seal migrate from %q to Shamir, no disabled seal in configuration",
				existBarrierSealConfig.Type)
		case existBarrierSealConfig.Type == wrapping.WrapperTypeShamir.String():
			// The configured seal is not Shamir, the stored seal config is Shamir.
			// This is a migration away from Shamir.
			unwrapSeal = NewDefaultSeal(&vaultseal.Access{
				Wrapper: aeadwrapper.NewShamirWrapper(),
			})
		default:
			// We know at this point that there is a configured non-Shamir seal,
			// that it does not match the stored non-Shamir seal config, and that
			// there is no explicit disabled seal stanza.
			return fmt.Errorf("cannot seal migrate from %q to %q, no disabled seal in configuration",
				existBarrierSealConfig.Type, c.seal.BarrierType())
		}
	} else {
		// If we're not coming from Shamir we expect the previous seal to be
		// in the config and disabled.

		if unwrapSeal.BarrierType() == wrapping.WrapperTypeShamir {
			return errors.New("Shamir seals cannot be set disabled (they should simply not be set)")
		}
	}

	// If we've reached this point it's a migration attempt and we should have both
	// c.migrationInfo.seal (old seal) and c.seal (new seal) populated.
	unwrapSeal.SetCore(c)

	// No stored recovery seal config found, what about the legacy recovery config?
	if existBarrierSealConfig.Type != wrapping.WrapperTypeShamir.String() && existRecoverySealConfig == nil {
		entry, err := c.physical.Get(ctx, recoverySealConfigPath)
		if err != nil {
			return fmt.Errorf("failed to read %q recovery seal configuration: %w", existBarrierSealConfig.Type, err)
		}
		if entry == nil {
			return errors.New("Recovery seal configuration not found for existing seal")
		}
		return errors.New("Cannot migrate seals while using a legacy recovery seal config")
	}

	c.migrationInfo = &migrationInformation{
		seal: unwrapSeal,
	}
	if existBarrierSealConfig.Type != c.seal.BarrierType().String() {
		// It's unnecessary to call this when doing an auto->auto
		// same-seal-type migration, since they'll have the same configs before
		// and after migration.
		c.adjustSealConfigDuringMigration(existBarrierSealConfig, existRecoverySealConfig)
	}
	c.initSealsForMigration()
	c.logger.Warn("entering seal migration mode; Vault will not automatically unseal even if using an autoseal", "from_barrier_type", c.migrationInfo.seal.BarrierType(), "to_barrier_type", c.seal.BarrierType())

	return nil
}

func (c *Core) migrateSealConfig(ctx context.Context) error {
	existBarrierSealConfig, existRecoverySealConfig, err := c.PhysicalSealConfigs(ctx)
	if err != nil {
		return fmt.Errorf("failed to read existing seal configuration during migration: %v", err)
	}

	var bc, rc *SealConfig

	switch {
	case c.migrationInfo.seal.RecoveryKeySupported() && c.seal.RecoveryKeySupported():
		// Migrating from auto->auto, copy the configs over
		bc, rc = existBarrierSealConfig, existRecoverySealConfig
	case c.migrationInfo.seal.RecoveryKeySupported():
		// Migrating from auto->shamir, clone auto's recovery config and set
		// stored keys to 1.
		bc = existRecoverySealConfig.Clone()
		bc.StoredShares = 1
	case c.seal.RecoveryKeySupported():
		// Migrating from shamir->auto, set a new barrier config and set
		// recovery config to a clone of shamir's barrier config with stored
		// keys set to 0.
		bc = &SealConfig{
			Type:            c.seal.BarrierType().String(),
			SecretShares:    1,
			SecretThreshold: 1,
			StoredShares:    1,
		}

		rc = existBarrierSealConfig.Clone()
		rc.StoredShares = 0
	}

	if err := c.seal.SetBarrierConfig(ctx, bc); err != nil {
		return fmt.Errorf("error storing barrier config after migration: %w", err)
	}

	if c.seal.RecoveryKeySupported() {
		if err := c.seal.SetRecoveryConfig(ctx, rc); err != nil {
			return fmt.Errorf("error storing recovery config after migration: %w", err)
		}
	} else if err := c.physical.Delete(ctx, recoverySealConfigPlaintextPath); err != nil {
		return fmt.Errorf("failed to delete old recovery seal configuration during migration: %w", err)
	}

	return nil
}

func (c *Core) adjustSealConfigDuringMigration(existBarrierSealConfig, existRecoverySealConfig *SealConfig) {
	switch {
	case c.migrationInfo.seal.RecoveryKeySupported() && existRecoverySealConfig != nil:
		// Migrating from auto->shamir, clone auto's recovery config and set
		// stored keys to 1.  Unless the recover config doesn't exist, in which
		// case the migration is assumed to already have been performed.
		newSealConfig := existRecoverySealConfig.Clone()
		newSealConfig.StoredShares = 1
		c.seal.SetCachedBarrierConfig(newSealConfig)
	case !c.migrationInfo.seal.RecoveryKeySupported() && c.seal.RecoveryKeySupported():
		// Migrating from shamir->auto, set a new barrier config and set
		// recovery config to a clone of shamir's barrier config with stored
		// keys set to 0.
		newBarrierSealConfig := &SealConfig{
			Type:            c.seal.BarrierType().String(),
			SecretShares:    1,
			SecretThreshold: 1,
			StoredShares:    1,
		}
		c.seal.SetCachedBarrierConfig(newBarrierSealConfig)

		newRecoveryConfig := existBarrierSealConfig.Clone()
		newRecoveryConfig.StoredShares = 0
		c.seal.SetCachedRecoveryConfig(newRecoveryConfig)
	}
}

func (c *Core) unsealKeyToRootKeyPostUnseal(ctx context.Context, combinedKey []byte) ([]byte, error) {
	return c.unsealKeyToMasterKey(ctx, c.seal, combinedKey, true, false)
}

func (c *Core) unsealKeyToMasterKeyPreUnseal(ctx context.Context, seal Seal, combinedKey []byte) ([]byte, error) {
	return c.unsealKeyToMasterKey(ctx, seal, combinedKey, false, true)
}

// unsealKeyToMasterKey takes a key provided by the user, either a recovery key
// if using an autoseal or an unseal key with Shamir.  It returns a nil error
// if the key is valid and an error otherwise. It also returns the master key
// that can be used to unseal the barrier.
// If useTestSeal is true, seal will not be modified; this is used when not
// invoked as part of an unseal process.  Otherwise in the non-legacy shamir
// case the combinedKey will be set in the seal, which means subsequent attempts
// to use the seal to read the master key will succeed, assuming combinedKey is
// valid.
// If allowMissing is true, a failure to find the master key in storage results
// in a nil error and a nil master key being returned.
func (c *Core) unsealKeyToMasterKey(ctx context.Context, seal Seal, combinedKey []byte, useTestSeal bool, allowMissing bool) ([]byte, error) {
	switch seal.StoredKeysSupported() {
	case vaultseal.StoredKeysSupportedGeneric:
		if err := seal.VerifyRecoveryKey(ctx, combinedKey); err != nil {
			return nil, fmt.Errorf("recovery key verification failed: %w", err)
		}

		storedKeys, err := seal.GetStoredKeys(ctx)
		if storedKeys == nil && err == nil && allowMissing {
			return nil, nil
		}

		if err == nil && len(storedKeys) != 1 {
			err = fmt.Errorf("expected exactly one stored key, got %d", len(storedKeys))
		}
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve stored keys: %w", err)
		}
		return storedKeys[0], nil

	case vaultseal.StoredKeysSupportedShamirRoot:
		if useTestSeal {
			testseal := NewDefaultSeal(&vaultseal.Access{
				Wrapper: aeadwrapper.NewShamirWrapper(),
			})
			testseal.SetCore(c)
			cfg, err := seal.BarrierConfig(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to setup test barrier config: %w", err)
			}
			testseal.SetCachedBarrierConfig(cfg)
			seal = testseal
		}

		err := seal.GetAccess().Wrapper.(*aeadwrapper.ShamirWrapper).SetAesGcmKeyBytes(combinedKey)
		if err != nil {
			return nil, &ErrInvalidKey{fmt.Sprintf("failed to setup unseal key: %v", err)}
		}
		storedKeys, err := seal.GetStoredKeys(ctx)
		if storedKeys == nil && err == nil && allowMissing {
			return nil, nil
		}
		if err == nil && len(storedKeys) != 1 {
			err = fmt.Errorf("expected exactly one stored key, got %d", len(storedKeys))
		}
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve stored keys: %w", err)
		}
		return storedKeys[0], nil

	case vaultseal.StoredKeysNotSupported:
		return combinedKey, nil
	}
	return nil, fmt.Errorf("invalid seal")
}

// IsInSealMigrationMode returns true if we're configured to perform a seal migration,
// meaning either that we have a disabled seal in HCL configuration or the seal
// configuration in storage is Shamir but the seal in HCL is not.  In this
// mode we should not auto-unseal (even if the migration is done) and we will
// accept unseal requests with and without the `migrate` option, though the migrate
// option is required if we haven't yet performed the seal migration.
func (c *Core) IsInSealMigrationMode() bool {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	return c.migrationInfo != nil
}

// IsSealMigrated returns true if we're in seal migration mode but migration
// has already been performed (possibly by another node, or prior to this node's
// current invocation.)
func (c *Core) IsSealMigrated() bool {
	if !c.IsInSealMigrationMode() {
		return false
	}
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	done, _ := c.sealMigrated(context.Background())
	return done
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

func (c *Core) IsPerfSecondary() bool {
	return c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary)
}

func (c *Core) AddLogger(logger log.Logger) {
	c.allLoggersLock.Lock()
	defer c.allLoggersLock.Unlock()
	c.allLoggers = append(c.allLoggers, logger)
}

// SetLogLevel sets logging level for all tracked loggers to the level provided
func (c *Core) SetLogLevel(level log.Level) {
	c.allLoggersLock.RLock()
	defer c.allLoggersLock.RUnlock()
	for _, logger := range c.allLoggers {
		logger.SetLevel(level)
	}
}

// SetLogLevelByName sets the logging level of named logger to level provided
// if it exists. Core.allLoggers is a slice and as such it is entirely possible
// that multiple entries exist for the same name. Each instance will be modified.
func (c *Core) SetLogLevelByName(name string, level log.Level) bool {
	c.allLoggersLock.RLock()
	defer c.allLoggersLock.RUnlock()

	found := false
	for _, logger := range c.allLoggers {
		if logger.Name() == name {
			logger.SetLevel(level)
			found = true
		}
	}

	return found
}

// SetConfig sets core's config object to the newly provided config.
func (c *Core) SetConfig(conf *server.Config) {
	c.rawConfig.Store(conf)
	bz, err := json.Marshal(c.SanitizedConfig())
	if err != nil {
		c.logger.Error("error serializing sanitized config", "error", err)
		return
	}

	c.logger.Debug("set config", "sanitized config", string(bz))
}

func (c *Core) GetListenerCustomResponseHeaders(listenerAdd string) *ListenerCustomHeaders {
	customHeaders := c.customListenerHeader.Load()
	if customHeaders == nil {
		return nil
	}

	customHeadersList, ok := customHeaders.([]*ListenerCustomHeaders)
	if customHeadersList == nil || !ok {
		return nil
	}

	for _, l := range customHeadersList {
		if l.Address == listenerAdd {
			return l
		}
	}
	return nil
}

// ExistCustomResponseHeader checks if a custom header is configured in any
// listener's stanza
func (c *Core) ExistCustomResponseHeader(header string) bool {
	customHeaders := c.customListenerHeader.Load()
	if customHeaders == nil {
		return false
	}

	customHeadersList, ok := customHeaders.([]*ListenerCustomHeaders)
	if customHeadersList == nil || !ok {
		return false
	}

	for _, l := range customHeadersList {
		exist := l.ExistCustomResponseHeader(header)
		if exist {
			return true
		}
	}

	return false
}

func (c *Core) ReloadCustomResponseHeaders() error {
	conf := c.rawConfig.Load()
	if conf == nil {
		return fmt.Errorf("failed to load core raw config")
	}
	lns := conf.(*server.Config).Listeners
	if lns == nil {
		return fmt.Errorf("no listener configured")
	}

	uiHeaders, err := c.UIHeaders()
	if err != nil {
		return err
	}
	c.customListenerHeader.Store(NewListenerCustomHeader(lns, c.logger, uiHeaders))

	return nil
}

// SanitizedConfig returns a sanitized version of the current config.
// See server.Config.Sanitized for specific values omitted.
func (c *Core) SanitizedConfig() map[string]interface{} {
	conf := c.rawConfig.Load()
	if conf == nil {
		return nil
	}
	return conf.(*server.Config).Sanitized()
}

// LogFormat returns the log format current in use.
func (c *Core) LogFormat() string {
	conf := c.rawConfig.Load()
	return conf.(*server.Config).LogFormat
}

// LogLevel returns the log level provided by level provided by config, CLI flag, or env
func (c *Core) LogLevel() string {
	return c.logLevel
}

// MetricsHelper returns the global metrics helper which allows external
// packages to access Vault's internal metrics.
func (c *Core) MetricsHelper() *metricsutil.MetricsHelper {
	return c.metricsHelper
}

// MetricSink returns the metrics wrapper with which Core has been configured.
func (c *Core) MetricSink() *metricsutil.ClusterMetricSink {
	return c.metricSink
}

// BuiltinRegistry is an interface that allows the "vault" package to use
// the registry of builtin plugins without getting an import cycle. It
// also allows for mocking the registry easily.
type BuiltinRegistry interface {
	Contains(name string, pluginType consts.PluginType) bool
	Get(name string, pluginType consts.PluginType) (func() (interface{}, error), bool)
	Keys(pluginType consts.PluginType) []string
	DeprecationStatus(name string, pluginType consts.PluginType) (consts.DeprecationStatus, bool)
}

func (c *Core) AuditLogger() AuditLogger {
	return &basicAuditor{c: c}
}

type FeatureFlags struct {
	NamespacesCubbyholesLocal bool `json:"namespace_cubbyholes_local"`
}

func (c *Core) persistFeatureFlags(ctx context.Context) error {
	if !c.PR1103disabled {
		c.logger.Debug("persisting feature flags")
		json, err := jsonutil.EncodeJSON(&FeatureFlags{NamespacesCubbyholesLocal: !c.PR1103disabled})
		if err != nil {
			return err
		}
		return c.barrier.Put(ctx, &logical.StorageEntry{
			Key:   consts.CoreFeatureFlagPath,
			Value: json,
		})
	}
	return nil
}

func (c *Core) readFeatureFlags(ctx context.Context) (*FeatureFlags, error) {
	entry, err := c.barrier.Get(ctx, consts.CoreFeatureFlagPath)
	if err != nil {
		return nil, err
	}
	var flags FeatureFlags
	if entry != nil {
		err = jsonutil.DecodeJSON(entry.Value, &flags)
		if err != nil {
			return nil, err
		}
	}
	return &flags, nil
}

// isMountable tells us whether or not we can continue mounting a plugin-based
// mount entry after failing to instantiate a backend. We do this to preserve
// the storage and path when a plugin is missing or has otherwise been
// misconfigured. This allows users to recover from errors when starting Vault
// with misconfigured plugins. It should not be possible for existing builtins
// to be misconfigured, so that is a fatal error.
func (c *Core) isMountable(ctx context.Context, entry *MountEntry, pluginType consts.PluginType) bool {
	return !c.isMountEntryBuiltin(ctx, entry, pluginType)
}

// isMountEntryBuiltin determines whether a mount entry is associated with a
// builtin of the specified plugin type.
func (c *Core) isMountEntryBuiltin(ctx context.Context, entry *MountEntry, pluginType consts.PluginType) bool {
	// Prevent a panic early on
	if entry == nil || c.pluginCatalog == nil {
		return false
	}

	// Allow type to be determined from mount entry when not otherwise specified
	if pluginType == consts.PluginTypeUnknown {
		pluginType = c.builtinTypeFromMountEntry(ctx, entry)
	}

	// Handle aliases
	pluginName := entry.Type
	if alias, ok := mountAliases[pluginName]; ok {
		pluginName = alias
	}

	plug, err := c.pluginCatalog.Get(ctx, pluginName, pluginType, entry.Version)
	if err != nil || plug == nil {
		return false
	}

	return plug.Builtin
}

// MatchingMount returns the path of the mount that will be responsible for
// handling the given request path.
func (c *Core) MatchingMount(ctx context.Context, reqPath string) string {
	return c.router.MatchingMount(ctx, reqPath)
}

func (c *Core) setupQuotas(ctx context.Context, isPerfStandby bool) error {
	if c.quotaManager == nil {
		return nil
	}

	return c.quotaManager.Setup(ctx, c.systemBarrierView, isPerfStandby, c.IsDRSecondary())
}

// ApplyRateLimitQuota checks the request against all the applicable quota rules.
// If the given request's path is exempt, no rate limiting will be applied.
func (c *Core) ApplyRateLimitQuota(ctx context.Context, req *quotas.Request) (quotas.Response, error) {
	req.Type = quotas.TypeRateLimit

	resp := quotas.Response{
		Allowed: true,
		Headers: make(map[string]string),
	}

	if c.quotaManager != nil {
		// skip rate limit checks for paths that are exempt from rate limiting
		if c.quotaManager.RateLimitPathExempt(req.Path) {
			return resp, nil
		}

		return c.quotaManager.ApplyQuota(ctx, req)
	}

	return resp, nil
}

// RateLimitAuditLoggingEnabled returns if the quota configuration allows audit
// logging of request rejections due to rate limiting quota rule violations.
func (c *Core) RateLimitAuditLoggingEnabled() bool {
	if c.quotaManager != nil {
		return c.quotaManager.RateLimitAuditLoggingEnabled()
	}

	return false
}

// RateLimitResponseHeadersEnabled returns if the quota configuration allows for
// rate limit quota HTTP headers to be added to responses.
func (c *Core) RateLimitResponseHeadersEnabled() bool {
	if c.quotaManager != nil {
		return c.quotaManager.RateLimitResponseHeadersEnabled()
	}

	return false
}

func (c *Core) KeyRotateGracePeriod() time.Duration {
	return time.Duration(atomic.LoadInt64(c.keyRotateGracePeriod))
}

func (c *Core) SetKeyRotateGracePeriod(t time.Duration) {
	atomic.StoreInt64(c.keyRotateGracePeriod, int64(t))
}

// Periodically test whether to automatically rotate the barrier key
func (c *Core) autoRotateBarrierLoop(ctx context.Context) {
	t := time.NewTicker(autoRotateCheckInterval)
	for {
		select {
		case <-t.C:
			c.checkBarrierAutoRotate(ctx)
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func (c *Core) checkBarrierAutoRotate(ctx context.Context) {
	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.isPrimary() {
		reason, err := c.barrier.CheckBarrierAutoRotate(ctx)
		if err != nil {
			lf := c.logger.Error
			if strings.HasSuffix(err.Error(), "context canceled") {
				lf = c.logger.Debug
			}
			lf("error in barrier auto rotation", "error", err)
			return
		}
		if reason != "" {
			// Time to rotate.  Invoke the rotation handler in order to both rotate and create
			// the replication canary
			c.logger.Info("automatic barrier key rotation triggered", "reason", reason)

			_, err := c.systemBackend.handleRotate(ctx, nil, nil)
			if err != nil {
				c.logger.Error("error automatically rotating barrier key", "error", err)
			} else {
				metrics.IncrCounter(barrierRotationsMetric, 1)
			}
		}
	}
}

func (c *Core) isPrimary() bool {
	return !c.ReplicationState().HasState(consts.ReplicationPerformanceSecondary | consts.ReplicationDRSecondary)
}

type LicenseState struct {
	State      string
	ExpiryTime time.Time
	Terminated bool
}

func (c *Core) loadLoginMFAConfigs(ctx context.Context) error {
	eConfigs := make([]*mfa.MFAEnforcementConfig, 0)
	allNamespaces := c.collectNamespaces()
	for _, ns := range allNamespaces {
		err := c.loginMFABackend.loadMFAMethodConfigs(ctx, ns)
		if err != nil {
			return fmt.Errorf("error loading MFA method Config, namespaceid %s, error: %w", ns.ID, err)
		}

		loadedConfigs, err := c.loginMFABackend.loadMFAEnforcementConfigs(ctx, ns)
		if err != nil {
			return fmt.Errorf("error loading MFA enforcement Config, namespaceid %s, error: %w", ns.ID, err)
		}

		eConfigs = append(eConfigs, loadedConfigs...)
	}

	for _, conf := range eConfigs {
		if err := c.loginMFABackend.loginMFAMethodExistenceCheck(conf); err != nil {
			c.loginMFABackend.mfaLogger.Error("failed to find all MFA methods that exist in MFA enforcement configs", "configID", conf.ID, "namespaceID", conf.NamespaceID, "error", err.Error())
		}
	}
	return nil
}

type MFACachedAuthResponse struct {
	CachedAuth            *logical.Auth
	RequestPath           string
	RequestNSID           string
	RequestNSPath         string
	RequestConnRemoteAddr string
	TimeOfStorage         time.Time
	RequestID             string
}

func (c *Core) setupCachedMFAResponseAuth() {
	c.mfaResponseAuthQueueLock.Lock()
	c.mfaResponseAuthQueue = NewLoginMFAPriorityQueue()
	mfaQueue := c.mfaResponseAuthQueue
	c.mfaResponseAuthQueueLock.Unlock()

	ctx := c.activeContext

	go func() {
		ticker := time.Tick(5 * time.Second)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker:
				err := mfaQueue.RemoveExpiredMfaAuthResponse(defaultMFAAuthResponseTTL, time.Now())
				if err != nil {
					c.Logger().Error("failed to remove stale MFA auth response", "error", err)
				}
			}
		}
	}()
	return
}

// updateLockedUserEntries runs every 15 mins to remove stale user entries from storage
// it also updates the userFailedLoginInfo map with correct information for locked users if incorrect
func (c *Core) updateLockedUserEntries() {
	ctx := c.activeContext
	go func() {
		ticker := time.NewTicker(15 * time.Minute)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				if err := c.runLockedUserEntryUpdates(ctx); err != nil {
					c.Logger().Error("failed to run locked user entry updates", "error", err)
				}
			}
		}
	}()
	return
}

// runLockedUserEntryUpdates runs updates for locked user storage entries and userFailedLoginInfo map
func (c *Core) runLockedUserEntryUpdates(ctx context.Context) error {
	// check environment variable to see if user lockout workflow is disabled
	var disableUserLockout bool
	if disableUserLockoutEnv := os.Getenv(consts.VaultDisableUserLockout); disableUserLockoutEnv != "" {
		var err error
		disableUserLockout, err = strconv.ParseBool(disableUserLockoutEnv)
		if err != nil {
			c.Logger().Error("Error parsing the environment variable VAULT_DISABLE_USER_LOCKOUT", "error", err)
		}
	}
	if disableUserLockout {
		return nil
	}

	// get the list of namespaces of locked users from locked users path in storage
	nsIDs, err := c.barrier.List(ctx, coreLockedUsersPath)
	if err != nil {
		return err
	}

	totalLockedUsersCount := 0
	for _, nsID := range nsIDs {
		// get the list of mount accessors of locked users for each namespace
		mountAccessors, err := c.barrier.List(ctx, coreLockedUsersPath+nsID)
		if err != nil {
			return err
		}

		// update the entries for locked users for each mount accessor
		// if storage entry is stale i.e; the lockout duration has passed
		// remove this entry from storage and userFailedLoginInfo map
		// else check if the userFailedLoginInfo map has correct failed login information
		// if incorrect, update the entry in userFailedLoginInfo map
		for _, mountAccessorPath := range mountAccessors {
			mountAccessor := strings.TrimSuffix(mountAccessorPath, "/")
			lockedAliasesCount, err := c.runLockedUserEntryUpdatesForMountAccessor(ctx, mountAccessor, coreLockedUsersPath+nsID+mountAccessorPath)
			if err != nil {
				return err
			}
			totalLockedUsersCount = totalLockedUsersCount + lockedAliasesCount
		}
	}

	// emit locked user count metrics
	metrics.SetGaugeWithLabels([]string{"core", "locked_users"}, float32(totalLockedUsersCount), nil)
	return nil
}

// runLockedUserEntryUpdatesForMountAccessor updates the storage entry for each locked user (alias name)
// if the entry is stale, it removes it from storage and userFailedLoginInfo map if present
// if the entry is not stale, it updates the userFailedLoginInfo map with correct values for entry if incorrect
func (c *Core) runLockedUserEntryUpdatesForMountAccessor(ctx context.Context, mountAccessor string, path string) (int, error) {
	// get mount entry for mountAccessor
	mountEntry := c.router.MatchingMountByAccessor(mountAccessor)
	if mountEntry == nil {
		mountEntry = &MountEntry{}
	}
	// get configuration for mount entry
	userLockoutConfiguration := c.getUserLockoutConfiguration(mountEntry)

	// get the list of aliases for mount accessor
	aliases, err := c.barrier.List(ctx, path)
	if err != nil {
		return 0, err
	}

	lockedAliasesCount := len(aliases)

	// check storage entry for each alias to update
	for _, alias := range aliases {
		loginUserInfoKey := FailedLoginUser{
			aliasName:     alias,
			mountAccessor: mountAccessor,
		}

		existingEntry, err := c.barrier.Get(ctx, path+alias)
		if err != nil {
			return 0, err
		}

		if existingEntry == nil {
			continue
		}

		var lastLoginTime int
		err = jsonutil.DecodeJSON(existingEntry.Value, &lastLoginTime)
		if err != nil {
			return 0, err
		}

		lastFailedLoginTimeFromStorageEntry := time.Unix(int64(lastLoginTime), 0)
		lockoutDurationFromConfiguration := userLockoutConfiguration.LockoutDuration

		// get the entry for the locked user from userFailedLoginInfo map
		failedLoginInfoFromMap := c.LocalGetUserFailedLoginInfo(ctx, loginUserInfoKey)

		// check if the storage entry for locked user is stale
		if time.Now().After(lastFailedLoginTimeFromStorageEntry.Add(lockoutDurationFromConfiguration)) {
			// stale entry, remove from storage
			// leaving this as it is as this happens on the active node
			// also handles case where namespace is deleted
			if err := c.barrier.Delete(ctx, path+alias); err != nil {
				return 0, err
			}
			// remove entry for this user from userFailedLoginInfo map if present as the user is not locked
			if failedLoginInfoFromMap != nil {
				if err = updateUserFailedLoginInfo(ctx, c, loginUserInfoKey, nil, true); err != nil {
					return 0, err
				}
			}
			lockedAliasesCount -= 1
			continue
		}

		// this is not a stale entry
		// update the map with actual failed login information
		actualFailedLoginInfo := FailedLoginInfo{
			lastFailedLoginTime: lastLoginTime,
			count:               uint(userLockoutConfiguration.LockoutThreshold),
		}

		if failedLoginInfoFromMap != &actualFailedLoginInfo {
			// entry is invalid, updating the entry in userFailedLoginMap with correct information
			if err = updateUserFailedLoginInfo(ctx, c, loginUserInfoKey, &actualFailedLoginInfo, false); err != nil {
				return 0, err
			}
		}
	}
	return lockedAliasesCount, nil
}

// PopMFAResponseAuthByID pops an item from the mfaResponseAuthQueue by ID
// it returns the cached auth response or an error
func (c *Core) PopMFAResponseAuthByID(reqID string) (*MFACachedAuthResponse, error) {
	c.mfaResponseAuthQueueLock.Lock()
	defer c.mfaResponseAuthQueueLock.Unlock()
	return c.mfaResponseAuthQueue.PopByKey(reqID)
}

// SaveMFAResponseAuth pushes an MFACachedAuthResponse to the mfaResponseAuthQueue.
// it returns an error in case of failure
func (c *Core) SaveMFAResponseAuth(respAuth *MFACachedAuthResponse) error {
	c.mfaResponseAuthQueueLock.Lock()
	defer c.mfaResponseAuthQueueLock.Unlock()
	return c.mfaResponseAuthQueue.Push(respAuth)
}

type InFlightRequests struct {
	InFlightReqMap   *sync.Map
	InFlightReqCount *uberAtomic.Uint64
}

type InFlightReqData struct {
	StartTime        time.Time `json:"start_time"`
	ClientRemoteAddr string    `json:"client_remote_address"`
	ReqPath          string    `json:"request_path"`
	Method           string    `json:"request_method"`
	ClientID         string    `json:"client_id"`
}

func (c *Core) StoreInFlightReqData(reqID string, data InFlightReqData) {
	c.inFlightReqData.InFlightReqMap.Store(reqID, data)
	c.inFlightReqData.InFlightReqCount.Inc()
}

// FinalizeInFlightReqData is going log the completed request if the
// corresponding server config option is enabled. It also removes the
// request from the inFlightReqMap and decrement the number of in-flight
// requests by one.
func (c *Core) FinalizeInFlightReqData(reqID string, statusCode int) {
	if c.logRequestsLevel != nil && c.logRequestsLevel.Load() != 0 {
		c.LogCompletedRequests(reqID, statusCode)
	}

	c.inFlightReqData.InFlightReqMap.Delete(reqID)
	c.inFlightReqData.InFlightReqCount.Dec()
}

// LoadInFlightReqData creates a snapshot map of the current
// in-flight requests
func (c *Core) LoadInFlightReqData() map[string]InFlightReqData {
	currentInFlightReqMap := make(map[string]InFlightReqData)
	c.inFlightReqData.InFlightReqMap.Range(func(key, value interface{}) bool {
		// there is only one writer to this map, so skip checking for errors
		v := value.(InFlightReqData)
		currentInFlightReqMap[key.(string)] = v
		return true
	})

	return currentInFlightReqMap
}

// UpdateInFlightReqData updates the data for a specific reqID with
// the clientID
func (c *Core) UpdateInFlightReqData(reqID, clientID string) {
	v, ok := c.inFlightReqData.InFlightReqMap.Load(reqID)
	if !ok {
		c.Logger().Trace("failed to retrieve request with ID", "request_id", reqID)
		return
	}

	// there is only one writer to this map, so skip checking for errors
	reqData := v.(InFlightReqData)
	reqData.ClientID = clientID
	c.inFlightReqData.InFlightReqMap.Store(reqID, reqData)
}

// LogCompletedRequests Logs the completed request to the server logs
func (c *Core) LogCompletedRequests(reqID string, statusCode int) {
	logLevel := log.Level(c.logRequestsLevel.Load())
	v, ok := c.inFlightReqData.InFlightReqMap.Load(reqID)
	if !ok {
		c.logger.Log(logLevel, fmt.Sprintf("failed to retrieve request with ID %v", reqID))
		return
	}

	// there is only one writer to this map, so skip checking for errors
	reqData := v.(InFlightReqData)
	c.logger.Log(logLevel, "completed_request",
		"start_time", reqData.StartTime.Format(time.RFC3339),
		"duration", fmt.Sprintf("%dms", time.Now().Sub(reqData.StartTime).Milliseconds()),
		"client_id", reqData.ClientID,
		"client_address", reqData.ClientRemoteAddr, "status_code", statusCode, "request_path", reqData.ReqPath,
		"request_method", reqData.Method)
}

func (c *Core) ReloadLogRequestsLevel() {
	conf := c.rawConfig.Load()
	if conf == nil {
		return
	}

	infoLevel := conf.(*server.Config).LogRequestsLevel
	switch {
	case log.LevelFromString(infoLevel) > log.NoLevel && log.LevelFromString(infoLevel) < log.Off:
		c.logRequestsLevel.Store(int32(log.LevelFromString(infoLevel)))
	case infoLevel != "":
		c.logger.Warn("invalid log_requests_level", "level", infoLevel)
	}
}

func (c *Core) ReloadIntrospectionEndpointEnabled() {
	conf := c.rawConfig.Load()
	if conf == nil {
		return
	}
	c.introspectionEnabledLock.Lock()
	defer c.introspectionEnabledLock.Unlock()
	c.introspectionEnabled = conf.(*server.Config).EnableIntrospectionEndpoint
}

type PeerNode struct {
	Hostname       string    `json:"hostname"`
	APIAddress     string    `json:"api_address"`
	ClusterAddress string    `json:"cluster_address"`
	Version        string    `json:"version"`
	LastEcho       time.Time `json:"last_echo"`
	UpgradeVersion string    `json:"upgrade_version,omitempty"`
	RedundancyZone string    `json:"redundancy_zone,omitempty"`
}

// GetHAPeerNodesCached returns the nodes that've sent us Echo requests recently.
func (c *Core) GetHAPeerNodesCached() []PeerNode {
	var nodes []PeerNode
	for itemClusterAddr, item := range c.clusterPeerClusterAddrsCache.Items() {
		info := item.Object.(nodeHAConnectionInfo)
		nodes = append(nodes, PeerNode{
			Hostname:       info.nodeInfo.Hostname,
			APIAddress:     info.nodeInfo.ApiAddr,
			ClusterAddress: itemClusterAddr,
			LastEcho:       info.lastHeartbeat,
			Version:        info.version,
			UpgradeVersion: info.upgradeVersion,
			RedundancyZone: info.redundancyZone,
		})
	}
	return nodes
}

func (c *Core) CheckPluginPerms(pluginName string) (err error) {
	var enableFilePermissionsCheck bool
	if enableFilePermissionsCheckEnv := os.Getenv(consts.VaultEnableFilePermissionsCheckEnv); enableFilePermissionsCheckEnv != "" {
		var err error
		enableFilePermissionsCheck, err = strconv.ParseBool(enableFilePermissionsCheckEnv)
		if err != nil {
			return errors.New("Error parsing the environment variable VAULT_ENABLE_FILE_PERMISSIONS_CHECK")
		}
	}

	if c.pluginDirectory != "" && enableFilePermissionsCheck {
		err = osutil.OwnerPermissionsMatch(c.pluginDirectory, c.pluginFileUid, c.pluginFilePermissions)
		if err != nil {
			return err
		}
		fullPath := filepath.Join(c.pluginDirectory, pluginName)
		err = osutil.OwnerPermissionsMatch(fullPath, c.pluginFileUid, c.pluginFilePermissions)
		if err != nil {
			return err
		}
	}
	return err
}

func (c *Core) LoadNodeID() (string, error) {
	raftNodeID := c.GetRaftNodeID()
	if raftNodeID != "" {
		return raftNodeID, nil
	}
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	return hostname, nil
}

// DetermineRoleFromLoginRequestFromBytes will determine the role that should be applied to a quota for a given
// login request, accepting a byte payload
func (c *Core) DetermineRoleFromLoginRequestFromBytes(mountPoint string, payload []byte, ctx context.Context) string {
	data := make(map[string]interface{})
	err := jsonutil.DecodeJSON(payload, &data)
	if err != nil {
		// Cannot discern a role from a request we cannot parse
		return ""
	}

	return c.DetermineRoleFromLoginRequest(mountPoint, data, ctx)
}

// DetermineRoleFromLoginRequest will determine the role that should be applied to a quota for a given
// login request
func (c *Core) DetermineRoleFromLoginRequest(mountPoint string, data map[string]interface{}, ctx context.Context) string {
	c.authLock.RLock()
	defer c.authLock.RUnlock()
	matchingBackend := c.router.MatchingBackend(ctx, mountPoint)
	if matchingBackend == nil || matchingBackend.Type() != logical.TypeCredential {
		// Role based quotas do not apply to this request
		return ""
	}

	resp, err := matchingBackend.HandleRequest(ctx, &logical.Request{
		MountPoint: mountPoint,
		Path:       "login",
		Operation:  logical.ResolveRoleOperation,
		Data:       data,
		Storage:    c.router.MatchingStorageByAPIPath(ctx, mountPoint+"login"),
	})
	if err != nil || resp.Data["role"] == nil {
		return ""
	}
	return resp.Data["role"].(string)
}

// aliasNameFromLoginRequest will determine the aliasName from the login Request
func (c *Core) aliasNameFromLoginRequest(ctx context.Context, req *logical.Request) (string, error) {
	c.authLock.RLock()
	defer c.authLock.RUnlock()
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return "", err
	}

	// ns path is added while checking matching backend
	mountPath := strings.TrimPrefix(req.MountPoint, ns.Path)

	matchingBackend := c.router.MatchingBackend(ctx, mountPath)
	if matchingBackend == nil || matchingBackend.Type() != logical.TypeCredential {
		// pathLoginAliasLookAhead operation does not apply to this request
		return "", nil
	}

	path := strings.ReplaceAll(req.Path, mountPath, "")

	resp, err := matchingBackend.HandleRequest(ctx, &logical.Request{
		MountPoint: req.MountPoint,
		Path:       path,
		Operation:  logical.AliasLookaheadOperation,
		Data:       req.Data,
		Storage:    c.router.MatchingStorageByAPIPath(ctx, req.Path),
	})
	if err != nil || resp.Auth.Alias == nil {
		return "", nil
	}
	return resp.Auth.Alias.Name, nil
}

// ListMounts will provide a slice containing a deep copy each mount entry
func (c *Core) ListMounts() ([]*MountEntry, error) {
	if c.Sealed() {
		return nil, fmt.Errorf("vault is sealed")
	}

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	var entries []*MountEntry

	for _, entry := range c.mounts.Entries {
		clone, err := entry.Clone()
		if err != nil {
			return nil, err
		}

		entries = append(entries, clone)
	}

	return entries, nil
}

// ListAuths will provide a slice containing a deep copy each auth entry
func (c *Core) ListAuths() ([]*MountEntry, error) {
	if c.Sealed() {
		return nil, fmt.Errorf("vault is sealed")
	}

	c.authLock.RLock()
	defer c.authLock.RUnlock()

	var entries []*MountEntry

	for _, entry := range c.auth.Entries {
		clone, err := entry.Clone()
		if err != nil {
			return nil, err
		}

		entries = append(entries, clone)
	}

	return entries, nil
}

type GroupPolicyApplicationMode struct {
	GroupPolicyApplicationMode string `json:"group_policy_application_mode"`
}

func (c *Core) GetGroupPolicyApplicationMode(ctx context.Context) (string, error) {
	se, err := c.barrier.Get(ctx, coreGroupPolicyApplicationPath)
	if err != nil {
		return "", err
	}
	if se == nil {
		return groupPolicyApplicationModeWithinNamespaceHierarchy, nil
	}

	var modeStruct GroupPolicyApplicationMode

	err = jsonutil.DecodeJSON(se.Value, &modeStruct)
	if err != nil {
		return "", err
	}
	mode := modeStruct.GroupPolicyApplicationMode
	if mode == "" {
		mode = groupPolicyApplicationModeWithinNamespaceHierarchy
	}

	return mode, nil
}

func (c *Core) SetGroupPolicyApplicationMode(ctx context.Context, mode string) error {
	json, err := jsonutil.EncodeJSON(&GroupPolicyApplicationMode{GroupPolicyApplicationMode: mode})
	if err != nil {
		return err
	}
	return c.barrier.Put(ctx, &logical.StorageEntry{
		Key:   coreGroupPolicyApplicationPath,
		Value: json,
	})
}

type HCPLinkStatus struct {
	lock             sync.RWMutex
	ConnectionStatus string `json:"hcp_link_status,omitempty"`
	ResourceIDOnHCP  string `json:"resource_ID_on_hcp,omitempty"`
}

func (c *Core) SetHCPLinkStatus(status, resourceID string) {
	c.hcpLinkStatus.lock.Lock()
	defer c.hcpLinkStatus.lock.Unlock()
	c.hcpLinkStatus.ConnectionStatus = status
	c.hcpLinkStatus.ResourceIDOnHCP = resourceID
}

func (c *Core) GetHCPLinkStatus() (string, string) {
	c.hcpLinkStatus.lock.RLock()
	defer c.hcpLinkStatus.lock.RUnlock()

	status := c.hcpLinkStatus.ConnectionStatus
	resourceID := c.hcpLinkStatus.ResourceIDOnHCP

	return status, resourceID
}

// IsExperimentEnabled is true if the experiment is enabled in the core.
func (c *Core) IsExperimentEnabled(experiment string) bool {
	return strutil.StrListContains(c.experiments, experiment)
}

// ListenerAddresses provides a slice of configured listener addresses
func (c *Core) ListenerAddresses() ([]string, error) {
	addresses := make([]string, 0)

	conf := c.rawConfig.Load()
	if conf == nil {
		return nil, fmt.Errorf("failed to load core raw config")
	}

	listeners := conf.(*server.Config).Listeners
	if listeners == nil {
		return nil, fmt.Errorf("no listener configured")
	}

	for _, listener := range listeners {
		addresses = append(addresses, listener.Address)
	}

	return addresses, nil
}

// IsRaftVoter specifies whether the node is a raft voter which is
// always false if raft storage is not in use.
func (c *Core) IsRaftVoter() bool {
	raftInfo := c.raftInfo.Load().(*raftInformation)

	if raftInfo == nil {
		return false
	}

	return !raftInfo.nonVoter
}

func (c *Core) HAEnabled() bool {
	return c.ha != nil && c.ha.HAEnabled()
}

func (c *Core) GetRaftConfiguration(ctx context.Context) (*raft.RaftConfigurationResponse, error) {
	raftBackend := c.getRaftBackend()

	if raftBackend == nil {
		return nil, nil
	}

	return raftBackend.GetConfiguration(ctx)
}

func (c *Core) GetRaftAutopilotState(ctx context.Context) (*raft.AutopilotState, error) {
	raftBackend := c.getRaftBackend()
	if raftBackend == nil {
		return nil, nil
	}

	return raftBackend.GetAutopilotServerState(ctx)
}

// Events returns a reference to the common event bus for sending and subscribint to events.
func (c *Core) Events() *eventbus.EventBus {
	return c.events
}
