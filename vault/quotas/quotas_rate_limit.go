package quotas

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/pathmanager"
	"golang.org/x/time/rate"
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
		"/v1/sys/generate-recovery-token/attempt",
		"/v1/sys/generate-recovery-token/update",
		"/v1/sys/generate-root/attempt",
		"/v1/sys/generate-root/update",
		"/v1/sys/health",
		"/v1/sys/seal-status",
		"/v1/sys/unseal",
	})
}

// ClientRateLimiter defines a token bucket based rate limiter for a unique
// addressable client (e.g. IP address). Whenever this client attempts to make
// a request, the lastSeen value will be updated.
type ClientRateLimiter struct {
	// lastSeen defines the UNIX timestamp the client last made a request.
	lastSeen time.Time

	// limiter represents an instance of a token bucket based rate limiter.
	limiter *rate.Limiter
}

// newClientRateLimiter returns a token bucket based rate limiter for a client
// that is uniquely addressable, where maxRequests defines the requests-per-second
// and burstSize defines the maximum burst allowed. A caller may provide -1 for
// burstSize to allow the burst value to be roughly equivalent to the RPS. Note,
// the underlying rate limiter is already thread-safe.
func newClientRateLimiter(maxRequests float64, burstSize int) *ClientRateLimiter {
	if burstSize < 0 {
		burstSize = int(math.Ceil(maxRequests))
	}

	return &ClientRateLimiter{
		lastSeen: time.Now().UTC(),
		limiter:  rate.NewLimiter(rate.Limit(maxRequests), burstSize),
	}
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

	// Rate defines the rate of which allowed requests are refilled per second.
	Rate float64 `json:"rate"`

	// Burst defines maximum number of requests at any given moment to be allowed.
	Burst int `json:"burst"`

	lock         *sync.RWMutex
	logger       log.Logger
	metricSink   *metricsutil.ClusterMetricSink
	purgeEnabled bool

	// purgeInterval defines the interval in seconds in which the RateLimitQuota
	// attempts to remove stale entries from the rateQuotas mapping.
	purgeInterval time.Duration
	closeCh       chan struct{}

	// staleAge defines the age in seconds in which a clientRateLimiter is
	// considered stale. A clientRateLimiter is considered stale if the delta
	// between the current purge time and its lastSeen timestamp is greater than
	// this value.
	staleAge time.Duration

	// rateQuotas contains a mapping from a unique addressable client (e.g. IP address)
	// to a clientRateLimiter reference. Every purgeInterval seconds, the RateLimitQuota
	// will attempt to remove stale entries from the mapping.
	rateQuotas map[string]*ClientRateLimiter
}

// NewRateLimitQuota creates a quota checker for imposing limits on the number
// of requests per second.
func NewRateLimitQuota(name, nsPath, mountPath string, rate float64, burst int) *RateLimitQuota {
	return &RateLimitQuota{
		Name:          name,
		Type:          TypeRateLimit,
		NamespacePath: nsPath,
		MountPath:     mountPath,
		Rate:          rate,
		Burst:         burst,
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

	if rlq.Rate <= 0 {
		return fmt.Errorf("invalid avg rps: %v", rlq.Rate)
	}

	if rlq.Burst < int(rlq.Rate) {
		return fmt.Errorf("burst size (%v) must be greater than or equal to average rps (%v)", rlq.Burst, rlq.Rate)
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

	rlq.purgeInterval = DefaultRateLimitPurgeInterval
	rlq.staleAge = DefaultRateLimitStaleAge
	rlq.rateQuotas = make(map[string]*ClientRateLimiter)

	if !rlq.purgeEnabled {
		rlq.purgeEnabled = true
		rlq.closeCh = make(chan struct{})
		go rlq.purgeClientsLoop()
	}

	return nil
}

func (rlq *RateLimitQuota) hasClient(addr string) bool {
	rlq.lock.RLock()
	defer rlq.lock.RUnlock()
	rlc, ok := rlq.rateQuotas[addr]
	return ok && rlc != nil
}

func (rlq *RateLimitQuota) numClients() int {
	rlq.lock.RLock()
	defer rlq.lock.RUnlock()
	return len(rlq.rateQuotas)
}

func (rlq *RateLimitQuota) getPurgeEnabled() bool {
	rlq.lock.RLock()
	defer rlq.lock.RUnlock()
	return rlq.purgeEnabled
}

// quotaID returns the identifier of the quota rule
func (rlq *RateLimitQuota) quotaID() string {
	return rlq.ID
}

// QuotaName returns the name of the quota rule
func (rlq *RateLimitQuota) QuotaName() string {
	return rlq.Name
}

// purgeClientsLoop performs a blocking process where every purgeInterval
// duration, we look for stale clients to remove from the rateQuotas map.
// A ClientRateLimiter is considered stale if its lastSeen timestamp exceeds the
// current time. The loop will continue to run indefinitely until a value is
// sent on the closeCh in which we stop the ticker and exit.
func (rlq *RateLimitQuota) purgeClientsLoop() {
	rlq.lock.RLock()
	ticker := time.NewTicker(rlq.purgeInterval)
	rlq.lock.RUnlock()

	for {
		select {
		case t := <-ticker.C:
			rlq.lock.Lock()

			for client, crl := range rlq.rateQuotas {
				if t.UTC().Sub(crl.lastSeen) >= rlq.staleAge {
					delete(rlq.rateQuotas, client)
				}
			}

			rlq.lock.Unlock()

		case <-rlq.closeCh:
			ticker.Stop()
			rlq.lock.Lock()
			rlq.purgeEnabled = false
			rlq.lock.Unlock()
			return
		}
	}
}

// clientRateLimiter returns a reference to a ClientRateLimiter based on a
// provided client address (e.g. IP address). If the ClientRateLimiter does not
// exist in the RateLimitQuota's mapping, one will be created and set. The
// created RateLimitQuota will have its requests-per-second set to
// RateLimitQuota.AverageRps. If the ClientRateLimiter already exists, the
// lastSeen timestamp will be updated.
func (rlq *RateLimitQuota) clientRateLimiter(addr string) *ClientRateLimiter {
	rlq.lock.Lock()
	defer rlq.lock.Unlock()

	crl, ok := rlq.rateQuotas[addr]
	if !ok {
		limiter := newClientRateLimiter(rlq.Rate, rlq.Burst)
		rlq.rateQuotas[addr] = limiter
		return limiter
	}

	crl.lastSeen = time.Now().UTC()
	return crl
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

	resp.Allowed = rlq.clientRateLimiter(req.ClientAddress).limiter.Allow()
	if !resp.Allowed {
		rlq.metricSink.IncrCounterWithLabels([]string{"quota", "rate_limit", "violation"}, 1, []metrics.Label{{"name", rlq.Name}})
	}

	return resp, nil
}

// close stops the current running client purge loop.
func (rlq *RateLimitQuota) close() error {
	close(rlq.closeCh)
	return nil
}

func (rlq *RateLimitQuota) handleRemount(toPath string) {
	rlq.MountPath = toPath
}
