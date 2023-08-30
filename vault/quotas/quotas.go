// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package quotas

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-memdb"
	"github.com/hashicorp/vault/helper/locking"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/pathmanager"
	"github.com/hashicorp/vault/sdk/logical"
)

// Type represents the quota kind
type Type string

const (
	// TypeRateLimit represents the rate limiting quota type
	TypeRateLimit Type = "rate-limit"

	// TypeLeaseCount represents the lease count limiting quota type
	TypeLeaseCount Type = "lease-count"
)

// LeaseAction is the action taken by the expiration manager on the lease. The
// quota manager will use this information to update the lease path cache and
// updating counters for relevant quota rules.
type LeaseAction uint32

// String converts each lease action into its string equivalent value
func (la LeaseAction) String() string {
	switch la {
	case LeaseActionLoaded:
		return "loaded"
	case LeaseActionCreated:
		return "created"
	case LeaseActionDeleted:
		return "deleted"
	case LeaseActionAllow:
		return "allow"
	}
	return "unknown"
}

const (
	_ LeaseAction = iota

	// LeaseActionLoaded indicates loading of lease in the expiration manager after
	// unseal.
	LeaseActionLoaded

	// LeaseActionCreated indicates that a lease is created in the expiration manager.
	LeaseActionCreated

	// LeaseActionDeleted indicates that is lease is expired and deleted in the
	// expiration manager.
	LeaseActionDeleted

	// LeaseActionAllow will be used to indicate the lease count checker that
	// incCounter is called from Allow(). All the rest of the actions indicate the
	// action took place on the lease in the expiration manager.
	LeaseActionAllow
)

type leaseWalkFunc func(context.Context, func(request *Request) bool) error

// String converts each quota type into its string equivalent value
func (q Type) String() string {
	switch q {
	case TypeLeaseCount:
		return "lease-count"
	case TypeRateLimit:
		return "rate-limit"
	}
	return "unknown"
}

const (
	indexID                 = "id"
	indexName               = "name"
	indexNamespace          = "ns"
	indexNamespaceMount     = "ns_mount"
	indexNamespaceMountPath = "ns_mount_path"
	indexNamespaceMountRole = "ns_mount_role"
)

const (
	// StoragePrefix is the prefix for the physical location where quota rules are
	// persisted.
	StoragePrefix = "quotas/"

	// ConfigPath is the physical location where the quota configuration is
	// persisted.
	ConfigPath = StoragePrefix + "config"

	// DefaultRateLimitExemptPathsToggle is the path to a toggle that allows us to
	// determine if a Vault operator explicitly modified the exempt paths set for
	// rate limit resource quotas. Specifically, when this toggle is false, we can
	// infer a Vault node is operating with an initial default set and on a subsequent
	// update to that set, we should not overwrite it on Setup.
	DefaultRateLimitExemptPathsToggle = StoragePrefix + "default_rate_limit_exempt_paths_toggle"
)

var (
	// ErrLeaseCountQuotaExceeded is returned when a request is rejected due to a lease
	// count quota being exceeded.
	ErrLeaseCountQuotaExceeded = errors.New("lease count quota exceeded")

	// ErrRateLimitQuotaExceeded is returned when a request is rejected due to a
	// rate limit quota being exceeded.
	ErrRateLimitQuotaExceeded = errors.New("rate limit quota exceeded")
)

var defaultExemptPaths = []string{
	"sys/generate-recovery-token/attempt",
	"sys/generate-recovery-token/update",
	"sys/generate-root/attempt",
	"sys/generate-root/update",
	"sys/health",
	"sys/seal-status",
	"sys/unseal",
}

// Access provides information to reach back to the quota checker.
type Access interface {
	// QuotaID is the identifier of the quota that issued this access.
	QuotaID() string
}

// Ensure that access implements the Access interface.
var _ Access = (*access)(nil)

// access implements the Access interface
type access struct {
	quotaID string
}

// QuotaID returns the identifier of the quota rule to which this access refers
// to.
func (a *access) QuotaID() string {
	return a.quotaID
}

// Manager holds all the existing quota rules. For any given input. the manager
// checks them against any applicable quota rules.
type Manager struct {
	entManager

	// db holds the in memory instances of all active quota rules indexed by
	// some of the quota properties.
	db *memdb.MemDB

	// config containing operator preferences and quota behaviors
	config *Config

	rateLimitPathManager *pathmanager.PathManager

	storage logical.Storage
	ctx     context.Context

	logger     log.Logger
	metricSink *metricsutil.ClusterMetricSink

	// quotaLock is a lock for manipulating quotas and anything not covered by a more specific lock
	quotaLock *locking.DeadlockRWMutex

	// quotaConfigLock is a lock for accessing config items, such as RateLimitExemptPaths
	quotaConfigLock *locking.DeadlockRWMutex

	// dbAndCacheLock is a lock for db and path caches that need to be reset during Reset()
	dbAndCacheLock *locking.DeadlockRWMutex
}

// QuotaLeaseInformation contains all of the information lease-count quotas require
// from a lease to uniquely identify the lease-count quota to increment/decrement
type QuotaLeaseInformation struct {
	// We can determine path and namespace from leaseId
	LeaseId string

	// We need the role as it's not part of the leaseId, and is required
	// to uniquely identify a lease count quota
	Role string
}

// Quota represents the common properties of every quota type
type Quota interface {
	// allow checks the if the request is allowed by the quota type implementation.
	allow(context.Context, *Request) (Response, error)

	// quotaID is the identifier of the quota rule
	quotaID() string

	// QuotaName is the name of the quota rule
	QuotaName() string

	// initialize sets up the fields in the quota type to begin operating
	initialize(log.Logger, *metricsutil.ClusterMetricSink) error

	// close defines any cleanup behavior that needs to be executed when a quota
	// rule is deleted.
	close(context.Context) error

	// Clone creates a clone of the calling quota
	Clone() Quota

	// Inheritable indicates if this quota can be applied to child namespaces
	IsInheritable() bool

	// handleRemount updates the mount and namesapce paths of the quota
	handleRemount(string, string)
}

// Response holds information about the result of the Allow() call. The response
// can optionally have the Access field set, which is used to reach back into
// the quota rule that sent this response.
type Response struct {
	// Allowed is set if the quota allows the request
	Allowed bool

	// Access is the handle to reach back into the quota rule that processed the
	// quota request. This may not be set all the time.
	Access Access

	// Headers defines any optional headers that may be returned by the quota rule
	// to clients.
	Headers map[string]string
}

// Config holds operator preferences around quota behaviors
type Config struct {
	// EnableRateLimitAuditLogging, if set, starts audit logging of the
	// request rejections that arise due to rate limit quota violations.
	EnableRateLimitAuditLogging bool `json:"enable_rate_limit_audit_logging"`

	// EnableRateLimitResponseHeaders dictates if rate limit quota HTTP headers
	// should be added to responses.
	EnableRateLimitResponseHeaders bool `json:"enable_rate_limit_response_headers"`

	// RateLimitExemptPaths defines the set of exempt paths used for all rate limit
	// quotas. Any request path that exists in this set is exempt from rate limiting.
	// If the set is empty, no paths are exempt.
	RateLimitExemptPaths []string `json:"rate_limit_exempt_paths"`
}

// Request contains information required by the quota manager to query and
// apply the quota rules.
type Request struct {
	// Type is the quota type
	Type Type

	// Path is the request path to which quota rules are being queried for
	Path string

	// Role is the role given as part of the request to a login endpoint
	Role string

	// NamespacePath is the namespace path to which the request belongs
	NamespacePath string

	// MountPath is the mount path to which the request is made
	MountPath string

	// ClientAddress is client unique addressable string (e.g. IP address). It can
	// be empty if the quota type does not need it.
	ClientAddress string
}

// NewManager creates and initializes a new quota manager to hold all the quota
// rules and to process incoming requests.
func NewManager(logger log.Logger, walkFunc leaseWalkFunc, ms *metricsutil.ClusterMetricSink) (*Manager, error) {
	db, err := memdb.NewMemDB(dbSchema())
	if err != nil {
		return nil, err
	}

	manager := &Manager{
		db:                   db,
		logger:               logger,
		metricSink:           ms,
		rateLimitPathManager: pathmanager.New(),
		config:               new(Config),
		quotaLock:            new(locking.DeadlockRWMutex),
		quotaConfigLock:      new(locking.DeadlockRWMutex),
		dbAndCacheLock:       new(locking.DeadlockRWMutex),
	}

	manager.init(walkFunc)

	return manager, nil
}

// SetQuota adds or updates a quota rule.
func (m *Manager) SetQuota(ctx context.Context, qType string, quota Quota, loading bool) error {
	m.quotaLock.Lock()
	m.dbAndCacheLock.RLock()
	defer m.quotaLock.Unlock()
	defer m.dbAndCacheLock.RUnlock()
	return m.setQuotaLocked(ctx, qType, quota, loading)
}

// setQuotaLocked creates a transaction, passes it into setQuotaLockedWithTxn and manages its lifecycle
// along with updating lease quota counts
func (m *Manager) setQuotaLocked(ctx context.Context, qType string, quota Quota, loading bool) error {
	txn := m.db.Txn(true)
	defer txn.Abort()

	err := m.setQuotaLockedWithTxn(ctx, qType, quota, loading, txn)
	if err != nil {
		return err
	}

	if loading {
		txn.Commit()
		return nil
	}

	// For the lease count type, recompute the counters
	if !loading && qType == TypeLeaseCount.String() {
		if err := m.recomputeLeaseCounts(ctx, txn); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

// setQuotaLockedWithTxn adds or updates a quota rule, modifying the db as well as
// any runtime elements such as goroutines, using the transaction passed in
// It should be called with the write lock held.
func (m *Manager) setQuotaLockedWithTxn(ctx context.Context, qType string, quota Quota, loading bool, txn *memdb.Txn) error {
	if qType == TypeLeaseCount.String() {
		m.setIsPerfStandby(quota)
	}

	raw, err := txn.First(qType, indexID, quota.quotaID())
	if err != nil {
		return err
	}

	// If there already exists an entry in the db, remove that first.
	if raw != nil {
		quota := raw.(Quota)
		if err := quota.close(ctx); err != nil {
			return err
		}
		err = txn.Delete(qType, raw)
		if err != nil {
			return err
		}
	}

	// Initialize the quota type implementation
	if err := quota.initialize(m.logger, m.metricSink); err != nil {
		return err
	}

	// Add the initialized quota type implementation to the db
	if err := txn.Insert(qType, quota); err != nil {
		return err
	}

	return nil
}

// QuotaNames returns the names of all the quota rules for a given type
func (m *Manager) QuotaNames(qType Type) ([]string, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	return m.quotaNamesLocked(qType)
}

// quotaNamesLocked returns the names of all the quota rules for a given type, and must be called with the lock
func (m *Manager) quotaNamesLocked(qType Type) ([]string, error) {
	txn := m.db.Txn(false)
	iter, err := txn.Get(qType.String(), indexID)
	if err != nil {
		return nil, err
	}
	var names []string
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		names = append(names, raw.(Quota).QuotaName())
	}
	return names, nil
}

// QuotaByID queries for a quota rule in the db for a given quota ID
func (m *Manager) QuotaByID(qType string, id string) (Quota, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	txn := m.db.Txn(false)

	quotaRaw, err := txn.First(qType, indexID, id)
	if err != nil {
		return nil, err
	}
	if quotaRaw == nil {
		return nil, nil
	}

	return quotaRaw.(Quota), nil
}

// QuotaByName queries for a quota rule in the db for a given quota name
func (m *Manager) QuotaByName(qType string, name string) (Quota, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	return m.quotaByNameLocked(qType, name)
}

// quotaByNameLocked queries for a quota rule in the db for a given quota name, and must be called with the lock
func (m *Manager) quotaByNameLocked(qType string, name string) (Quota, error) {
	txn := m.db.Txn(false)

	quotaRaw, err := txn.First(qType, indexName, name)
	if err != nil {
		return nil, err
	}
	if quotaRaw == nil {
		return nil, nil
	}

	return quotaRaw.(Quota), nil
}

// QuotaByFactors returns the quota rule that matches the provided factors
func (m *Manager) QuotaByFactors(ctx context.Context, qType, nsPath, mountPath, pathSuffix, role string) (Quota, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	// nsPath would have been made non-empty during insertion. Use non-empty value
	// during query as well.
	if nsPath == "" {
		nsPath = "root"
	}

	idx := indexNamespace
	args := []interface{}{nsPath, false, false, false}
	if mountPath != "" {
		if pathSuffix != "" {
			idx = indexNamespaceMountPath
			args = []interface{}{nsPath, mountPath, pathSuffix, false}
		} else if role != "" {
			idx = indexNamespaceMountRole
			args = []interface{}{nsPath, mountPath, false, role}
		} else {
			idx = indexNamespaceMount
			args = []interface{}{nsPath, mountPath, false, false}
		}
	}

	txn := m.db.Txn(false)
	iter, err := txn.Get(qType, idx, args...)
	if err != nil {
		return nil, err
	}
	var quotas []Quota
	for raw := iter.Next(); raw != nil; raw = iter.Next() {
		quotas = append(quotas, raw.(Quota))
	}
	if len(quotas) > 1 {
		m.logger.Debug("conflicting quotas in QuotaByFactors", "matching_quotas", quotas)
		return nil, fmt.Errorf("conflicting quota definitions detected")
	}
	if len(quotas) == 0 {
		return nil, nil
	}

	return quotas[0], nil
}

// QueryQuota returns the most specific applicable quota for a given request.
func (m *Manager) QueryQuota(req *Request) (Quota, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	return m.queryQuota(nil, req)
}

// queryQuota returns the quota rule that is applicable for the given request. It
// queries all the quota rules that are defined against request values and finds
// the quota rule that takes priority.
//
// Priority rules are as follows:
// - namespace specific quota takes precedence over global quota
// - mount specific quota takes precedence over namespace specific quota
// - path suffix specific quota takes precedence over mount specific quota
// - role based quota takes precedence over path suffix/mount specific quota
func (m *Manager) queryQuota(txn *memdb.Txn, req *Request) (Quota, error) {
	if txn == nil {
		txn = m.db.Txn(false)
	}

	// ns would have been made non-empty during insertion. Use non-empty
	// value during query as well.
	if req.NamespacePath == "" {
		req.NamespacePath = "root"
	}

	//
	// Find a match from most specific applicable quota rule to less specific one.
	//
	quotaFetchFunc := func(idx string, args ...interface{}) (Quota, error) {
		iter, err := txn.Get(req.Type.String(), idx, args...)
		if err != nil {
			return nil, err
		}
		var quotas []Quota
		for raw := iter.Next(); raw != nil; raw = iter.Next() {
			quota := raw.(Quota)
			quotas = append(quotas, quota)
		}
		if len(quotas) > 1 {
			m.logger.Debug("conflicting quotas in queryQuota", "matching_quotas", quotas, "args", args)
			return nil, fmt.Errorf("conflicting quota definitions detected")
		}
		if len(quotas) == 0 {
			return nil, nil
		}

		return quotas[0], nil
	}

	// Fetch role suffix quota
	quota, err := quotaFetchFunc(indexNamespaceMountRole, req.NamespacePath, req.MountPath, false, req.Role)
	if err != nil {
		return nil, err
	}
	if quota != nil {
		return quota, nil
	}

	// Fetch path suffix quota
	pathSuffix := strings.TrimSuffix(strings.TrimPrefix(strings.TrimPrefix(req.Path, req.NamespacePath), req.MountPath), "/")
	quota, err = quotaFetchFunc(indexNamespaceMountPath, req.NamespacePath, req.MountPath, pathSuffix, false)
	if err != nil {
		return nil, err
	}
	if quota != nil {
		return quota, nil
	}

	// Fetch path suffix quotas with globbing
	// Request paths which match the resulting glob (i.e. share the same prefix prior to the glob) are in scope for the quota
	for i := 0; i <= len(pathSuffix); i++ {
		trimmedSuffixWithGlob := pathSuffix[:len(pathSuffix)-i] + "*"
		// Check to see if a quota exists with this particular pattern
		quota, err = quotaFetchFunc(indexNamespaceMountPath, req.NamespacePath, req.MountPath, trimmedSuffixWithGlob, false)
		if err != nil {
			return nil, err
		}
		if quota != nil {
			return quota, nil
		}
	}

	// Fetch mount quota
	quota, err = quotaFetchFunc(indexNamespaceMount, req.NamespacePath, req.MountPath, false, false)
	if err != nil {
		return nil, err
	}
	if quota != nil {
		return quota, nil
	}

	// Fetch ns quota. If NamespacePath is root, this will return the global quota.
	quota, err = quotaFetchFunc(indexNamespace, req.NamespacePath, false, false, false)
	if err != nil {
		return nil, err
	}
	if quota != nil {
		return quota, nil
	}

	// Fetch parent ns quotas
	curNsSplitPath := strings.SplitAfter(namespace.Canonicalize(req.NamespacePath), "/")
	for len(curNsSplitPath) > 2 {
		parentNs := strings.Join(curNsSplitPath[0:len(curNsSplitPath)-2], "")
		parentQuota, err := quotaFetchFunc(indexNamespace, parentNs, false, false, false)
		if err != nil {
			return nil, err
		}
		if parentQuota != nil && parentQuota.IsInheritable() {
			return parentQuota, nil
		}
		curNsSplitPath = strings.SplitAfter(parentNs, "/")
	}

	// If the request belongs to "root" namespace, then we have already looked at
	// global quotas when fetching namespace specific quota rule. When the request
	// belongs to a non-root namespace, and when there are no namespace specific
	// quota rules present, we fall back on the global quotas.
	if req.NamespacePath == "root" {
		return nil, nil
	}

	// Fetch global quota
	quota, err = quotaFetchFunc(indexNamespace, "root", false, false, false)
	if err != nil {
		return nil, err
	}
	if quota != nil {
		return quota, nil
	}

	return nil, nil
}

// QueryResolveRoleQuotas checks if there's a quota for the request mount path
// which requires ResolveRoleOperation.
func (m *Manager) QueryResolveRoleQuotas(req *Request) (bool, error) {
	m.dbAndCacheLock.RLock()
	defer m.dbAndCacheLock.RUnlock()

	txn := m.db.Txn(false)

	// ns would have been made non-empty during insertion. Use non-empty
	// value during query as well.
	if req.NamespacePath == "" {
		req.NamespacePath = "root"
	}

	// Check for any role-based quotas on the request namespaces/mount path.
	for _, qType := range quotaTypes() {
		// Use the namespace and mount as indexes and find all matches with a
		// set Role field (see: 'true' as the last argument). We can't use
		// indexNamespaceMountRole for this, because Role is a StringFieldIndex,
		// which won't match on an empty string.
		quota, err := txn.First(qType, indexNamespaceMount, req.NamespacePath, req.MountPath, false, true)
		if err != nil {
			return false, err
		}
		if quota != nil {
			return true, nil
		}
	}
	return false, nil
}

// DeleteQuota removes a quota rule from the db for a given name
func (m *Manager) DeleteQuota(ctx context.Context, qType string, name string) error {
	m.quotaLock.Lock()
	m.dbAndCacheLock.RLock()
	defer m.quotaLock.Unlock()
	defer m.dbAndCacheLock.RUnlock()

	txn := m.db.Txn(true)
	defer txn.Abort()

	raw, err := txn.First(qType, indexName, name)
	if err != nil {
		return err
	}
	if raw == nil {
		return nil
	}

	quota := raw.(Quota)
	if err := quota.close(ctx); err != nil {
		return err
	}

	err = txn.Delete(qType, raw)
	if err != nil {
		return err
	}

	// For the lease count type, recompute the counters
	if qType == TypeLeaseCount.String() {
		if err := m.recomputeLeaseCounts(ctx, txn); err != nil {
			return err
		}
	}

	txn.Commit()
	return nil
}

// ApplyQuota runs the request against any quota rule that is applicable to it. If
// there are multiple quota rule that matches the request parameters, rule that
// takes precedence will be used to allow/reject the request.
func (m *Manager) ApplyQuota(ctx context.Context, req *Request) (Response, error) {
	var resp Response

	quota, err := m.QueryQuota(req)
	if err != nil {
		return resp, err
	}

	// If there is no quota defined, allow the request.
	if quota == nil {
		resp.Allowed = true
		return resp, nil
	}

	// If the quota type is lease count, and if the path is not known to
	// generate leases, allow the request.
	if req.Type == TypeLeaseCount && !m.inLeasePathCache(req.Path) {
		resp.Allowed = true
		return resp, nil
	}

	return quota.allow(ctx, req)
}

// SetEnableRateLimitAuditLogging updates the operator preference regarding the
// audit logging behavior.
func (m *Manager) SetEnableRateLimitAuditLogging(val bool) {
	m.quotaLock.Lock()
	defer m.quotaLock.Unlock()
	m.setEnableRateLimitAuditLoggingLocked(val)
}

func (m *Manager) setEnableRateLimitAuditLoggingLocked(val bool) {
	m.config.EnableRateLimitAuditLogging = val
}

// SetEnableRateLimitResponseHeaders updates the operator preference regarding
// the rate limit quota HTTP header behavior.
func (m *Manager) SetEnableRateLimitResponseHeaders(val bool) {
	m.quotaConfigLock.Lock()
	defer m.quotaConfigLock.Unlock()
	m.setEnableRateLimitResponseHeadersLocked(val)
}

func (m *Manager) setEnableRateLimitResponseHeadersLocked(val bool) {
	m.config.EnableRateLimitResponseHeaders = val
}

// SetRateLimitExemptPaths updates the rate limit exempt paths in the Manager's
// configuration in addition to updating the path manager. Every call to
// SetRateLimitExemptPaths will wipe out the existing path manager and set the
// paths based on the provided argument.
func (m *Manager) SetRateLimitExemptPaths(vals []string) {
	m.quotaConfigLock.Lock()
	defer m.quotaConfigLock.Unlock()
	m.setRateLimitExemptPathsLocked(vals)
}

func (m *Manager) setRateLimitExemptPathsLocked(vals []string) {
	if vals == nil {
		vals = []string{}
	}
	m.config.RateLimitExemptPaths = vals
	m.rateLimitPathManager = pathmanager.New()
	m.rateLimitPathManager.AddPaths(vals)
}

// RateLimitAuditLoggingEnabled returns if the quota configuration allows audit
// logging of request rejections due to rate limiting quota rule violations.
func (m *Manager) RateLimitAuditLoggingEnabled() bool {
	m.quotaConfigLock.RLock()
	defer m.quotaConfigLock.RUnlock()

	return m.config.EnableRateLimitAuditLogging
}

// RateLimitResponseHeadersEnabled returns if the quota configuration allows for
// rate limit quota HTTP headers to be added to responses.
func (m *Manager) RateLimitResponseHeadersEnabled() bool {
	m.quotaConfigLock.RLock()
	defer m.quotaConfigLock.RUnlock()

	return m.config.EnableRateLimitResponseHeaders
}

// RateLimitPathExempt returns a boolean dictating if a given path is exempt from
// any rate limit quota. If not rate limit path manager is defined, false is
// returned.
func (m *Manager) RateLimitPathExempt(path string) bool {
	m.quotaConfigLock.RLock()
	defer m.quotaConfigLock.RUnlock()

	if m.rateLimitPathManager == nil {
		return false
	}

	return m.rateLimitPathManager.HasPath(path)
}

// Config returns the operator preferences in the quota manager
func (m *Manager) Config() *Config {
	return m.config
}

// Reset will clear all the quotas from the db and clear the lease path cache.
func (m *Manager) Reset() error {
	m.quotaLock.Lock()
	defer m.quotaLock.Unlock()
	m.dbAndCacheLock.Lock()
	defer m.dbAndCacheLock.Unlock()

	err := m.resetCache()
	if err != nil {
		return err
	}
	m.storage = nil
	m.ctx = nil

	return m.entManager.Reset()
}

// Must be called with the lock held
func (m *Manager) resetCache() error {
	names, err := m.quotaNamesLocked(TypeRateLimit)
	if err != nil {
		return err
	}
	for _, name := range names {
		quota, err := m.quotaByNameLocked(TypeRateLimit.String(), name)
		if err != nil {
			return err
		}
		if quota != nil {
			rlq := quota.(*RateLimitQuota)
			err = rlq.store.Close(context.Background())
			if err != nil {
				return err
			}
		}
	}
	db, err := memdb.NewMemDB(dbSchema())
	if err != nil {
		return err
	}
	m.db = db
	return nil
}

// dbSchema creates a DB schema for holding all the quota rules. It creates
// table for each supported type of quota.
func dbSchema() *memdb.DBSchema {
	schema := &memdb.DBSchema{
		Tables: make(map[string]*memdb.TableSchema),
	}

	commonSchema := func(name string) *memdb.TableSchema {
		return &memdb.TableSchema{
			Name: name,
			Indexes: map[string]*memdb.IndexSchema{
				indexID: {
					Name:   indexID,
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "ID",
					},
				},
				indexName: {
					Name:   indexName,
					Unique: true,
					Indexer: &memdb.StringFieldIndex{
						Field: "Name",
					},
				},
				indexNamespace: {
					Name: indexNamespace,
					Indexer: &memdb.CompoundMultiIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{
								Field: "NamespacePath",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "MountPath",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "PathSuffix",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "Role",
							},
						},
					},
				},
				indexNamespaceMount: {
					Name:         indexNamespaceMount,
					AllowMissing: true,
					Indexer: &memdb.CompoundMultiIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{
								Field: "NamespacePath",
							},
							&memdb.StringFieldIndex{
								Field: "MountPath",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "PathSuffix",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "Role",
							},
						},
					},
				},
				indexNamespaceMountRole: {
					Name:         indexNamespaceMountRole,
					AllowMissing: true,
					Indexer: &memdb.CompoundMultiIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{
								Field: "NamespacePath",
							},
							&memdb.StringFieldIndex{
								Field: "MountPath",
							},
							// By sending false as the query parameter, we can
							// query just the role specific quota.
							&memdb.FieldSetIndex{
								Field: "PathSuffix",
							},
							&memdb.StringFieldIndex{
								Field: "Role",
							},
						},
					},
				},
				indexNamespaceMountPath: {
					Name:         indexNamespaceMountPath,
					AllowMissing: true,
					Indexer: &memdb.CompoundMultiIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{
								Field: "NamespacePath",
							},
							&memdb.StringFieldIndex{
								Field: "MountPath",
							},
							&memdb.StringFieldIndex{
								Field: "PathSuffix",
							},
							// By sending false as the query parameter, we can
							// query just the namespace specific quota.
							&memdb.FieldSetIndex{
								Field: "Role",
							},
						},
					},
				},
			},
		}
	}

	// Create a table per quota type. This allows names to be reused between
	// different quota types and querying a bit easier.
	for _, name := range quotaTypes() {
		schema.Tables[name] = commonSchema(name)
	}

	return schema
}

// Invalidate receives notifications from the replication sub-system when a key
// is updated in the storage. This function will read the key from storage and
// updates the caches and data structures to reflect those updates.
func (m *Manager) Invalidate(key string) {
	switch key {
	case "config":
		config, err := LoadConfig(m.ctx, m.storage)
		if err != nil {
			m.logger.Error("failed to invalidate quota config", "error", err)
			return
		}

		m.SetEnableRateLimitAuditLogging(config.EnableRateLimitAuditLogging)
		m.SetEnableRateLimitResponseHeaders(config.EnableRateLimitResponseHeaders)
		m.SetRateLimitExemptPaths(config.RateLimitExemptPaths)

	default:
		splitKeys := strings.Split(key, "/")
		if len(splitKeys) != 2 {
			m.logger.Error("incorrect key while invalidating quota rule", "key", key)
			return
		}
		qType := splitKeys[0]
		name := splitKeys[1]

		if qType == TypeLeaseCount.String() && m.isDRSecondary {
			// lease count invalidation not supported on DR Secondary
			return
		}

		// Read quota rule from storage
		quota, err := Load(m.ctx, m.storage, qType, name)
		if err != nil {
			m.logger.Error("failed to read invalidated quota rule", "error", err)
			return
		}

		switch {
		case quota == nil:
			// Handle quota deletion
			if err := m.DeleteQuota(m.ctx, qType, name); err != nil {
				m.logger.Error("failed to delete invalidated quota rule", "error", err)
				return
			}
		default:
			// Handle quota update
			if err := m.SetQuota(m.ctx, qType, quota, false); err != nil {
				m.logger.Error("failed to update invalidated quota rule", "error", err)
				return
			}
		}
	}
}

// LoadConfig reads the quota configuration from the underlying storage
func LoadConfig(ctx context.Context, storage logical.Storage) (*Config, error) {
	var config Config
	entry, err := storage.Get(ctx, ConfigPath)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return &config, nil
	}

	err = entry.DecodeJSON(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// Load reads the quota rule from the underlying storage
func Load(ctx context.Context, storage logical.Storage, qType, name string) (Quota, error) {
	var quota Quota
	entry, err := storage.Get(ctx, QuotaStoragePath(qType, name))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	switch qType {
	case TypeRateLimit.String():
		quota = &RateLimitQuota{}
	case TypeLeaseCount.String():
		quota = &LeaseCountQuota{}
	default:
		return nil, fmt.Errorf("unsupported type: %v", qType)
	}

	err = entry.DecodeJSON(quota)
	if err != nil {
		return nil, err
	}

	return quota, nil
}

// Setup loads the quota configuration and all the quota rules into the
// quota manager.
func (m *Manager) Setup(ctx context.Context, storage logical.Storage, isPerfStandby, isDRSecondary bool) error {
	m.quotaLock.Lock()
	m.quotaConfigLock.Lock()
	m.dbAndCacheLock.Lock()
	defer m.quotaLock.Unlock()
	defer m.quotaConfigLock.Unlock()
	defer m.dbAndCacheLock.Unlock()

	m.storage = storage
	m.ctx = ctx
	m.isPerfStandby = isPerfStandby
	m.isDRSecondary = isDRSecondary

	// Load the quota configuration from storage and load it into the quota
	// manager.
	config, err := LoadConfig(ctx, storage)
	if err != nil {
		return err
	}

	entry, err := storage.Get(ctx, DefaultRateLimitExemptPathsToggle)
	if err != nil {
		return err
	}

	// Determine if we need to set the default set of exempt paths for rate limit
	// resource quotas. We use a default set introduced in 1.5 when the toggle
	// entry does not exist in storage or is false. The toggle is flipped , i.e.
	// set to true when SetRateLimitExemptPaths is called during a config update.
	var toggle bool
	if entry != nil {
		if err := entry.DecodeJSON(&toggle); err != nil {
			return err
		}
	}

	exemptPaths := defaultExemptPaths
	if toggle {
		exemptPaths = config.RateLimitExemptPaths
	}

	m.setEnableRateLimitAuditLoggingLocked(config.EnableRateLimitAuditLogging)
	m.setEnableRateLimitResponseHeadersLocked(config.EnableRateLimitResponseHeaders)
	m.setRateLimitExemptPathsLocked(exemptPaths)
	if err = m.resetCache(); err != nil {
		return err
	}

	for _, qType := range quotaTypes() {
		m.setupQuotaType(ctx, storage, qType)
	}

	return nil
}

func (m *Manager) setupQuotaType(ctx context.Context, storage logical.Storage, quotaType string) error {
	if quotaType == TypeLeaseCount.String() && m.isDRSecondary {
		m.logger.Trace("lease count quotas are not processed on DR Secondaries")
		return nil
	}

	names, err := logical.CollectKeys(ctx, logical.NewStorageView(storage, StoragePrefix+quotaType+"/"))
	if err != nil {
		return err
	}
	for _, name := range names {
		quota, err := Load(ctx, m.storage, quotaType, name)
		if err != nil {
			return err
		}

		if quota == nil {
			continue
		}

		err = m.setQuotaLocked(ctx, quotaType, quota, true)
		if err != nil {
			return err
		}
	}

	return nil
}

// QuotaStoragePath returns the storage path suffix for persisting the quota
// rule.
func QuotaStoragePath(quotaType, name string) string {
	return path.Join(StoragePrefix+quotaType, name)
}

// HandleRemount updates the quota subsystem about the remount operation that
// took place. Quota manager will trigger the quota specific updates including
// the mount path update and the namespace update
func (m *Manager) HandleRemount(ctx context.Context, from, to namespace.MountPathDetails) error {
	m.quotaLock.Lock()
	m.dbAndCacheLock.RLock()
	defer m.quotaLock.Unlock()
	defer m.dbAndCacheLock.RUnlock()

	// Grab a write transaction, as we want to save the updated quota in memdb
	txn := m.db.Txn(true)
	defer txn.Abort()

	// quota namespace would have been made non-empty during insertion. Use non-empty value
	// during query as well.
	fromNs := from.Namespace.Path
	if fromNs == "" {
		fromNs = namespace.RootNamespaceID
	}

	toNs := to.Namespace.Path
	if toNs == "" {
		toNs = namespace.RootNamespaceID
	}

	leaseQuotaUpdated := false

	updateMounts := func(idx string, args ...interface{}) error {
		for _, quotaType := range quotaTypes() {
			iter, err := txn.Get(quotaType, idx, args...)
			if err != nil {
				return err
			}
			for raw := iter.Next(); raw != nil; raw = iter.Next() {
				quota := raw.(Quota)

				// Clone the object and update it
				clonedQuota := quota.Clone()
				clonedQuota.handleRemount(to.MountPath, toNs)
				// Update both underlying storage and memdb with the quota change
				entry, err := logical.StorageEntryJSON(QuotaStoragePath(quotaType, quota.QuotaName()), quota)
				if err != nil {
					return err
				}
				if err := m.storage.Put(ctx, entry); err != nil {
					return err
				}
				if err := m.setQuotaLockedWithTxn(ctx, quotaType, clonedQuota, false, txn); err != nil {
					return err
				}
				if quotaType == TypeLeaseCount.String() {
					leaseQuotaUpdated = true
				}
			}
		}
		return nil
	}

	// Update mounts for everything without a path prefix or role
	err := updateMounts(indexNamespaceMount, fromNs, from.MountPath, false, false)
	if err != nil {
		return err
	}

	// Update mounts for everything with a path prefix
	err = updateMounts(indexNamespaceMount, fromNs, from.MountPath, true, false)
	if err != nil {
		return err
	}

	// Update mounts for everything with a role
	err = updateMounts(indexNamespaceMount, fromNs, from.MountPath, false, true)
	if err != nil {
		return err
	}

	if leaseQuotaUpdated {
		if err := m.recomputeLeaseCounts(ctx, txn); err != nil {
			return err
		}
	}

	txn.Commit()

	return nil
}

// HandleBackendDisabling updates the quota subsystem with the disabling of auth
// or secret engine disabling. This should only be called on the primary cluster
// node.
func (m *Manager) HandleBackendDisabling(ctx context.Context, nsPath, mountPath string) error {
	m.quotaLock.Lock()
	m.dbAndCacheLock.RLock()
	defer m.quotaLock.Unlock()
	defer m.dbAndCacheLock.RUnlock()

	txn := m.db.Txn(true)
	defer txn.Abort()

	// nsPath would have been made non-empty during insertion. Use non-empty value
	// during query as well.
	if nsPath == "" {
		nsPath = "root"
	}

	leaseQuotaDeleted := false

	updateMounts := func(idx string, args ...interface{}) error {
		for _, quotaType := range quotaTypes() {
			iter, err := txn.Get(quotaType, idx, args...)
			if err != nil {
				return err
			}
			for raw := iter.Next(); raw != nil; raw = iter.Next() {
				if err := txn.Delete(quotaType, raw); err != nil {
					return fmt.Errorf("failed to delete quota from db after mount disabling; namespace %q, err %v", nsPath, err)
				}
				quota := raw.(Quota)
				if err := m.storage.Delete(ctx, QuotaStoragePath(quotaType, quota.QuotaName())); err != nil {
					return fmt.Errorf("failed to delete quota from storage after mount disabling; namespace %q, err %v", nsPath, err)
				}
				if quotaType == TypeLeaseCount.String() {
					leaseQuotaDeleted = true
				}
			}
		}
		return nil
	}

	// Update mounts for everything without a path prefix or role
	err := updateMounts(indexNamespaceMount, nsPath, mountPath, false, false)
	if err != nil {
		return err
	}

	// Update mounts for everything with a path prefix
	err = updateMounts(indexNamespaceMount, nsPath, mountPath, true, false)
	if err != nil {
		return err
	}

	// Update mounts for everything with a role
	err = updateMounts(indexNamespaceMount, nsPath, mountPath, false, true)
	if err != nil {
		return err
	}

	if leaseQuotaDeleted {
		if err := m.recomputeLeaseCounts(ctx, txn); err != nil {
			return err
		}
	}

	txn.Commit()

	return nil
}
