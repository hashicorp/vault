package quotas

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/pathmanager"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

var rateLimitExemptPaths = pathmanager.New()

const (
	// DefaultRateLimitPurgeInterval defines the default purge interval used by a
	// RateLimitQuota to remove stale client rate limiters.
	DefaultRateLimitPurgeInterval = time.Minute

	// DefaultRateLimitStaleAge defines the default stale age of a client limiter.
	DefaultRateLimitStaleAge = 3 * time.Minute

	// EnvVaultEnableRateLimitAuditLogging is used to enable audit logging of
	// requests that get rejected due to rate limit quota violations.
	EnvVaultEnableRateLimitAuditLogging = "VAULT_ENABLE_RATE_LIMIT_AUDIT_LOGGING"
)

func init() {
	rateLimitExemptPaths.AddPaths([]string{
		"sys/internal/ui/mounts",
		"sys/generate-recovery-token/attempt",
		"sys/generate-recovery-token/update",
		"sys/generate-root/attempt",
		"sys/generate-root/update",
		"sys/health",
		"sys/seal-status",
		"sys/unseal",
	})
}

// Ensure that RateLimitQuota implements the Quota interface
var _ Quota = (*RateLimitQuota)(nil)

// RateLimitQuota represents the quota rule properties that is used to limit the
// number of requests per second for a namespace or mount.
type RateLimitQuota struct {
	// ID is the identifier of the quota
	ID string `json:"id"`

	// Type of quota this represents
	Type Type `json:"type"`

	// Name of the quota rule
	Name string `json:"name"`

	// NamespacePath is the path of the namespace to which this quota is
	// applicable.
	NamespacePath string `json:"namespace_path"`

	// MountPath is the path of the mount to which this quota is applicable
	MountPath string `json:"mount_path"`

	// Rate defines the number of requests allowed per Interval.
	Rate float64 `json:"rate"`

	// Interval defines the duration to which rate limiting is applied.
	Interval time.Duration `json:"interval"`

	lock          *sync.RWMutex
	store         limiter.Store
	logger        log.Logger
	metricSink    *metricsutil.ClusterMetricSink
	purgeInterval time.Duration
	staleAge      time.Duration
}

// NewRateLimitQuota creates a quota checker for imposing limits on the number
// of requests per second.
func NewRateLimitQuota(name, nsPath, mountPath string, rate float64, interval time.Duration) *RateLimitQuota {
	return &RateLimitQuota{
		Name:          name,
		Type:          TypeRateLimit,
		NamespacePath: nsPath,
		MountPath:     mountPath,
		Rate:          rate,
		Interval:      interval,
		purgeInterval: DefaultRateLimitPurgeInterval,
		staleAge:      DefaultRateLimitStaleAge,
	}
}

// initialize ensures the namespace and max requests are initialized, sets the ID
// if it's currently empty, sets the purge interval and stale age to default
// values, and finally starts the client purge go routine if it has been started
// already. Note, initialize will reset the internal rateQuotas mapping.
func (rlq *RateLimitQuota) initialize(logger log.Logger, ms *metricsutil.ClusterMetricSink) error {
	if rlq.lock == nil {
		rlq.lock = new(sync.RWMutex)
	}

	rlq.lock.Lock()
	defer rlq.lock.Unlock()

	// Memdb requires a non-empty value for indexing
	if rlq.NamespacePath == "" {
		rlq.NamespacePath = "root"
	}

	if rlq.Interval == 0 {
		rlq.Interval = time.Second
	}

	if rlq.Rate <= 0 {
		return fmt.Errorf("invalid avg rps: %v", rlq.Rate)
	}

	if logger != nil {
		rlq.logger = logger
	}

	if rlq.metricSink == nil {
		rlq.metricSink = ms
	}

	if rlq.ID == "" {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return err
		}

		rlq.ID = id
	}

	rlStore, err := memorystore.New(&memorystore.Config{
		Tokens:        uint64(math.Round(rlq.Rate)), // allow 'rlq.Rate' number of requests per 'Interval'
		Interval:      rlq.Interval,                 // time interval in which to enforce rate limiting
		SweepInterval: rlq.purgeInterval,            // how often stale clients are removed
		SweepMinTTL:   rlq.staleAge,                 // how long since the last request a client is considered stale
	})
	if err != nil {
		return err
	}

	rlq.store = rlStore

	return nil
}

// quotaID returns the identifier of the quota rule
func (rlq *RateLimitQuota) quotaID() string {
	return rlq.ID
}

// QuotaName returns the name of the quota rule
func (rlq *RateLimitQuota) QuotaName() string {
	return rlq.Name
}

// allow decides if the request is allowed by the quota. An error will be
// returned if the request ID or address is empty. If the path is exempt, the
// quota will not be evaluated. Otherwise, the client rate limiter is retrieved
// by address and the rate limit quota is checked against that limiter.
func (rlq *RateLimitQuota) allow(req *Request) (Response, error) {
	var resp Response

	// Skip rate limit checks for paths that are exempt from rate limiting.
	if rateLimitExemptPaths.HasPath(req.Path) {
		resp.Allowed = true
		return resp, nil
	}

	if req.ClientAddress == "" {
		return resp, fmt.Errorf("missing request client address in quota request")
	}

	limit, remaining, reset, allow := rlq.store.Take(req.ClientAddress)
	resp.Allowed = allow
	resp.Headers = map[string]string{
		httplimit.HeaderRateLimitLimit:     strconv.FormatUint(limit, 10),
		httplimit.HeaderRateLimitRemaining: strconv.FormatUint(remaining, 10),
		httplimit.HeaderRateLimitReset:     time.Unix(0, int64(reset)).UTC().Format(time.RFC1123),
	}

	if !resp.Allowed {
		resp.Headers[httplimit.HeaderRetryAfter] = resp.Headers[httplimit.HeaderRateLimitReset]
		rlq.metricSink.IncrCounterWithLabels([]string{"quota", "rate_limit", "violation"}, 1, []metrics.Label{{"name", rlq.Name}})
	}

	return resp, nil
}

// close stops the current running client purge loop.
func (rlq *RateLimitQuota) close() error {
	if rlq.store != nil {
		return rlq.store.Close()
	}

	return nil
}

func (rlq *RateLimitQuota) handleRemount(toPath string) {
	rlq.MountPath = toPath
}
