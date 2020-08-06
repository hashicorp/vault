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
// number of requests in a given interval for a namespace or mount.
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

	// Block defines the duration during which all requests are blocked for a given
	// client. This interval is enforced only if non-zero and a client reaches the
	// rate limit.
	Block time.Duration `json:"block"`

	lock                *sync.RWMutex
	store               limiter.Store
	logger              log.Logger
	metricSink          *metricsutil.ClusterMetricSink
	purgeInterval       time.Duration
	staleAge            time.Duration
	blockedClients      map[string]time.Time
	purgeBlocked        bool
	closePurgeBlockedCh chan struct{}
}

// NewRateLimitQuota creates a quota checker for imposing limits on the number
// of requests in a given interval. An interval time duration of zero may be
// provided, which will default to 1s when initialized. An optional block
// duration may be provided, where if set, when a client reaches the rate limit,
// subsequent requests will fail until the block duration has passed.
func NewRateLimitQuota(name, nsPath, mountPath string, rate float64, interval, block time.Duration) *RateLimitQuota {
	return &RateLimitQuota{
		Name:          name,
		Type:          TypeRateLimit,
		NamespacePath: nsPath,
		MountPath:     mountPath,
		Rate:          rate,
		Interval:      interval,
		Block:         block,
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
		return fmt.Errorf("invalid 'rate': %v", rlq.Rate)
	}

	if rlq.Block < 0 {
		return fmt.Errorf("invalid 'block': %v", rlq.Block)
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
	rlq.blockedClients = make(map[string]time.Time)

	if rlq.Block > 0 && !rlq.purgeBlocked {
		rlq.purgeBlocked = true
		rlq.closePurgeBlockedCh = make(chan struct{})
		go rlq.purgeBlockedClients()
	}

	return nil
}

// purgeBlockedClients performs a blocking process where every purgeInterval
// duration, we look at all blocked clients to potentially remove from the blocked
// clients map.
//
// A blocked client will only be removed if the current time minus the time the
// client was blocked at is greater than or equal to the block duration. The loop
// will continue to run indefinitely until a value is	sent on the closePurgeBlockedCh
// in which we stop the ticker and return.
func (rlq *RateLimitQuota) purgeBlockedClients() {
	rlq.lock.RLock()
	ticker := time.NewTicker(rlq.purgeInterval)
	rlq.lock.RUnlock()

	for {
		select {
		case t := <-ticker.C:
			rlq.lock.Lock()

			for client, blockedAt := range rlq.blockedClients {
				if t.UTC().Sub(blockedAt) >= rlq.Block {
					delete(rlq.blockedClients, client)
				}
			}

			rlq.lock.Unlock()

		case <-rlq.closePurgeBlockedCh:
			ticker.Stop()

			rlq.lock.Lock()
			rlq.purgeBlocked = false
			rlq.lock.Unlock()

			return
		}
	}
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
	resp := Response{
		Headers: make(map[string]string),
	}

	// Skip rate limit checks for paths that are exempt from rate limiting.
	if rateLimitExemptPaths.HasPath(req.Path) {
		resp.Allowed = true
		return resp, nil
	}

	if req.ClientAddress == "" {
		return resp, fmt.Errorf("missing request client address in quota request")
	}

	var retryAfter string

	defer func() {
		if !resp.Allowed {
			resp.Headers[httplimit.HeaderRetryAfter] = retryAfter
			rlq.metricSink.IncrCounterWithLabels([]string{"quota", "rate_limit", "violation"}, 1, []metrics.Label{{"name", rlq.Name}})
		}
	}()

	rlq.lock.Lock()
	defer rlq.lock.Unlock()

	// Check if the client is currently blocked and if so, deny the request. Note,
	// we cannot simply rely on the presence of the client in the map as the timing
	// of purging blocked clients may not yield a false negative. In other words,
	// a client may no longer be considered blocked whereas the purging interval
	// has yet to run.
	if blockedAt, ok := rlq.blockedClients[req.ClientAddress]; ok {
		if time.Since(blockedAt) >= rlq.Block {
			// allow the request and remove the blocked client
			delete(rlq.blockedClients, req.ClientAddress)
		} else {
			// deny the request and return early
			resp.Allowed = false
			retryAfter = blockedAt.Add(rlq.Block).Format(time.RFC1123)
			return resp, nil
		}
	}

	limit, remaining, reset, allow := rlq.store.Take(req.ClientAddress)
	resp.Allowed = allow
	resp.Headers[httplimit.HeaderRateLimitLimit] = strconv.FormatUint(limit, 10)
	resp.Headers[httplimit.HeaderRateLimitRemaining] = strconv.FormatUint(remaining, 10)
	resp.Headers[httplimit.HeaderRateLimitReset] = time.Unix(0, int64(reset)).UTC().Format(time.RFC1123)
	retryAfter = resp.Headers[httplimit.HeaderRateLimitReset]

	return resp, nil
}

// close stops the current running client purge loop.
func (rlq *RateLimitQuota) close() error {
	if rlq.purgeBlocked {
		close(rlq.closePurgeBlockedCh)
	}

	if rlq.store != nil {
		return rlq.store.Close()
	}

	return nil
}

func (rlq *RateLimitQuota) handleRemount(toPath string) {
	rlq.MountPath = toPath
}
