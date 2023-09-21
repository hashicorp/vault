// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package quotas

import (
	"context"
	"encoding/hex"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/armon/go-metrics"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/sethvargo/go-limiter"
	"github.com/sethvargo/go-limiter/httplimit"
	"github.com/sethvargo/go-limiter/memorystore"
)

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

	// Role is the role on an auth mount to apply the quota to upon /login requests
	// Not applicable for use with path suffixes
	Role string `json:"role"`

	// PathSuffix is the path suffix to which this quota is applicable
	PathSuffix string `json:"path_suffix"`

	// Inheritable indicates whether the quota will be inherited by child namespaces
	Inheritable bool `json:"inheritable"`

	// Rate defines the number of requests allowed per Interval.
	Rate float64 `json:"rate"`

	// Interval defines the duration to which rate limiting is applied.
	Interval time.Duration `json:"interval"`

	// BlockInterval defines the duration during which all requests are blocked for
	// a given client. This interval is enforced only if non-zero and a client
	// reaches the rate limit.
	BlockInterval time.Duration `json:"block_interval"`

	lock                *sync.RWMutex
	store               limiter.Store
	logger              log.Logger
	metricSink          *metricsutil.ClusterMetricSink
	purgeInterval       time.Duration
	staleAge            time.Duration
	blockedClients      sync.Map
	purgeBlocked        bool
	closePurgeBlockedCh chan struct{}
}

// NewRateLimitQuota creates a quota checker for imposing limits on the number
// of requests in a given interval. An interval time duration of zero may be
// provided, which will default to 1s when initialized. An optional block
// duration may be provided, where if set, when a client reaches the rate limit,
// subsequent requests will fail until the block duration has passed.
func NewRateLimitQuota(name, nsPath, mountPath, pathSuffix, role string, inheritable bool, interval, block time.Duration, rate float64) *RateLimitQuota {
	id, err := uuid.GenerateUUID()
	if err != nil {
		// Fall back to generating with a hash of the name, later in initialize
		id = ""
	}
	return &RateLimitQuota{
		Name:          name,
		ID:            id,
		Type:          TypeRateLimit,
		NamespacePath: nsPath,
		MountPath:     mountPath,
		Role:          role,
		PathSuffix:    pathSuffix,
		Inheritable:   inheritable,
		Rate:          rate,
		Interval:      interval,
		BlockInterval: block,
		purgeInterval: DefaultRateLimitPurgeInterval,
		staleAge:      DefaultRateLimitStaleAge,
	}
}

func (q *RateLimitQuota) Clone() Quota {
	rlq := &RateLimitQuota{
		ID:            q.ID,
		Name:          q.Name,
		MountPath:     q.MountPath,
		Role:          q.Role,
		Inheritable:   q.Inheritable,
		Type:          q.Type,
		NamespacePath: q.NamespacePath,
		PathSuffix:    q.PathSuffix,
		BlockInterval: q.BlockInterval,
		Rate:          q.Rate,
		Interval:      q.Interval,
	}
	return rlq
}

func (q *RateLimitQuota) IsInheritable() bool {
	return q.Inheritable
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
		return fmt.Errorf("invalid rate: %v", rlq.Rate)
	}

	if rlq.BlockInterval < 0 {
		return fmt.Errorf("invalid block interval: %v", rlq.BlockInterval)
	}

	if logger != nil {
		rlq.logger = logger
	}

	if rlq.metricSink == nil {
		rlq.metricSink = ms
	}

	if rlq.ID == "" {
		// A lease which was created with a blank ID may have been persisted
		// to storage already (this is the case up to release 1.6.2.)
		// So, performance standby nodes could call initialize() on their copy
		// of the lease; for consistency we need to generate an ID that is
		// deterministic. That ensures later invalidation removes the original
		// lease from the memdb, instead of creating a duplicate.
		rlq.ID = hex.EncodeToString(cryptoutil.Blake2b256Hash(rlq.Name))
	}

	// Set purgeInterval if coming from a previous version where purgeInterval was
	// not defined.
	if rlq.purgeInterval == 0 {
		rlq.purgeInterval = DefaultRateLimitPurgeInterval
	}

	// Set staleAge if coming from a previous version where staleAge was not defined.
	if rlq.staleAge == 0 {
		rlq.staleAge = DefaultRateLimitStaleAge
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
	rlq.blockedClients = sync.Map{}

	if rlq.BlockInterval > 0 && !rlq.purgeBlocked {
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
	if rlq.purgeInterval <= 0 {
		rlq.purgeInterval = DefaultRateLimitPurgeInterval
	}
	ticker := time.NewTicker(rlq.purgeInterval)
	rlq.lock.RUnlock()

	for {
		select {
		case t := <-ticker.C:
			rlq.blockedClients.Range(func(key, value interface{}) bool {
				blockedAt := value.(time.Time)
				if t.Sub(blockedAt) >= rlq.BlockInterval {
					rlq.blockedClients.Delete(key)
				}

				return true
			})

		case <-rlq.closePurgeBlockedCh:
			ticker.Stop()

			rlq.lock.Lock()
			rlq.purgeBlocked = false
			rlq.lock.Unlock()

			return
		}
	}
}

func (rlq *RateLimitQuota) getPurgeBlocked() bool {
	rlq.lock.RLock()
	defer rlq.lock.RUnlock()
	return rlq.purgeBlocked
}

func (rlq *RateLimitQuota) numBlockedClients() int {
	rlq.lock.RLock()
	defer rlq.lock.RUnlock()

	size := 0
	rlq.blockedClients.Range(func(_, _ interface{}) bool {
		size++
		return true
	})

	return size
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
func (rlq *RateLimitQuota) allow(ctx context.Context, req *Request) (Response, error) {
	resp := Response{
		Headers: make(map[string]string),
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

	// Check if the client is currently blocked and if so, deny the request. Note,
	// we cannot simply rely on the presence of the client in the map as the timing
	// of purging blocked clients may not yield a false negative. In other words,
	// a client may no longer be considered blocked whereas the purging interval
	// has yet to run.
	if v, ok := rlq.blockedClients.Load(req.ClientAddress); ok {
		blockedAt := v.(time.Time)
		if time.Since(blockedAt) >= rlq.BlockInterval {
			// allow the request and remove the blocked client
			rlq.blockedClients.Delete(req.ClientAddress)
		} else {
			// deny the request and return early
			resp.Allowed = false
			retryAfter = strconv.Itoa(int(time.Until(blockedAt.Add(rlq.BlockInterval)).Seconds()))
			return resp, nil
		}
	}

	limit, remaining, reset, allow, err := rlq.store.Take(ctx, req.ClientAddress)
	if err != nil {
		return resp, err
	}

	resp.Allowed = allow
	resp.Headers[httplimit.HeaderRateLimitLimit] = strconv.FormatUint(limit, 10)
	resp.Headers[httplimit.HeaderRateLimitRemaining] = strconv.FormatUint(remaining, 10)
	resp.Headers[httplimit.HeaderRateLimitReset] = strconv.Itoa(int(time.Until(time.Unix(0, int64(reset))).Seconds()))
	retryAfter = resp.Headers[httplimit.HeaderRateLimitReset]

	// If the request is not allowed (i.e. rate limit threshold reached) and blocking
	// is enabled, we add the client to the set of blocked clients.
	if !resp.Allowed && rlq.purgeBlocked {
		blockedAt := time.Now()
		retryAfter = strconv.Itoa(int(time.Until(blockedAt.Add(rlq.BlockInterval)).Seconds()))
		rlq.blockedClients.Store(req.ClientAddress, blockedAt)
	}

	return resp, nil
}

// close stops the current running client purge loop.
// It should be called with the write lock held.
func (rlq *RateLimitQuota) close(ctx context.Context) error {
	if rlq.purgeBlocked {
		close(rlq.closePurgeBlockedCh)
	}

	if rlq.store != nil {
		return rlq.store.Close(ctx)
	}

	return nil
}

func (rlq *RateLimitQuota) handleRemount(mountpath, nspath string) {
	rlq.MountPath = mountpath
	rlq.NamespacePath = nspath
}
