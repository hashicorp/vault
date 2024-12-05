// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/armon/go-metrics"
	"github.com/hashicorp/vault/helper/metricsutil"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/helper/pluginconsts"
	"github.com/hashicorp/vault/limits"
	"github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/logical"
)

func (c *Core) metricsLoop(stopCh chan struct{}) {
	emitTimer := time.Tick(time.Second)

	stopOrHAState := func() (bool, consts.HAState) {
		l := newLockGrabber(c.stateLock.RLock, c.stateLock.RUnlock, stopCh)
		go l.grab()
		if stopped := l.lockOrStop(); stopped {
			return true, 0
		}
		defer c.stateLock.RUnlock()
		return false, c.HAState()
	}

	identityCountTimer := time.Tick(time.Minute * 10)
	// Only emit on active node of cluster that is not a DR secondary.
	if stopped, haState := stopOrHAState(); stopped {
		return
	} else if haState == consts.Standby || c.IsDRSecondary() {
		identityCountTimer = nil
	}

	writeTimer := time.Tick(time.Second * 30)
	// Do not process the writeTimer on DR Secondary nodes
	if c.IsDRSecondary() {
		writeTimer = nil
	}

	// This loop covers
	// vault.expire.num_leases
	// vault.core.unsealed
	// vault.identity.num_entities
	// and the non-telemetry request counters shown in the UI.
	for {
		select {
		case <-emitTimer:
			stopped, haState := stopOrHAState()
			if stopped {
				return
			}
			if haState == consts.Active {
				c.metricsMutex.Lock()
				// Emit on active node only
				if c.expiration != nil {
					c.expiration.emitMetrics()
				}
				c.metricsMutex.Unlock()
			}

			// Refresh the sealed gauge, on all nodes
			if c.Sealed() {
				c.metricSink.SetGaugeWithLabels([]string{"core", "unsealed"}, 0, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "unsealed"}, 1, nil)
			}

			if c.UndoLogsEnabled() {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "write_undo_logs"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "write_undo_logs"}, 0, nil)
			}

			writeLimiter := c.GetRequestLimiter(limits.WriteLimiter)
			if writeLimiter != nil {
				c.metricSink.SetGaugeWithLabels([]string{
					"core", "limits", "concurrency", limits.WriteLimiter,
				}, float32(writeLimiter.EstimatedLimit()), nil)
			}

			pathLimiter := c.GetRequestLimiter(limits.SpecialPathLimiter)
			if pathLimiter != nil {
				c.metricSink.SetGaugeWithLabels([]string{
					"core", "limits", "concurrency", limits.SpecialPathLimiter,
				}, float32(pathLimiter.EstimatedLimit()), nil)
			}

			// Refresh the standby gauge, on all nodes
			if haState != consts.Active {
				c.metricSink.SetGaugeWithLabels([]string{"core", "active"}, 0, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "active"}, 1, nil)
			}

			if haState == consts.PerfStandby {
				c.metricSink.SetGaugeWithLabels([]string{"core", "performance_standby"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "performance_standby"}, 0, nil)
			}

			if c.ReplicationState().HasState(consts.ReplicationPerformancePrimary) {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "performance", "primary"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "performance", "primary"}, 0, nil)
			}

			if c.IsPerfSecondary() {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "performance", "secondary"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "performance", "secondary"}, 0, nil)
			}

			if c.ReplicationState().HasState(consts.ReplicationDRPrimary) {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "dr", "primary"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "dr", "primary"}, 0, nil)
			}

			if c.IsDRSecondary() {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "dr", "secondary"}, 1, nil)
			} else {
				c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "dr", "secondary"}, 0, nil)
			}

			if haState == consts.Active {
				reindexState := c.ReindexStage()
				if reindexState != nil {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "reindex_stage"}, float32(*reindexState), nil)
				} else {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "reindex_stage"}, 0, nil)
				}

				buildProgress := c.BuildProgress()
				if buildProgress != nil {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "build_progress"}, float32(*buildProgress), nil)
				} else {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "build_progress"}, 0, nil)
				}

				buildTotal := c.BuildTotal()
				if buildTotal != nil {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "build_total"}, float32(*buildTotal), nil)
				} else {
					c.metricSink.SetGaugeWithLabels([]string{"core", "replication", "build_total"}, 0, nil)
				}
			}

			// If we're using a raft backend, emit raft metrics
			if rb, ok := c.underlyingPhysical.(*raft.RaftBackend); ok {
				rb.CollectMetrics(c.MetricSink())
			}

			// Capture the total number of in-flight requests
			c.inFlightReqGaugeMetric()

			// Refresh gauge metrics that are looped
			c.cachedGaugeMetricsEmitter()
		case <-writeTimer:
			l := newLockGrabber(c.stateLock.RLock, c.stateLock.RUnlock, stopCh)
			go l.grab()
			if stopped := l.lockOrStop(); stopped {
				return
			}
			// Ship barrier encryption counts if a perf standby or the active node
			// on a performance secondary cluster
			if c.perfStandby || c.IsPerfSecondary() { // already have lock here, do not re-acquire
				err := syncBarrierEncryptionCounter(c)
				if err != nil {
					c.logger.Error("writing syncing encryption counters", "err", err)
				}
			}
			c.stateLock.RUnlock()
		case <-identityCountTimer:
			// TODO: this can be replaced by the identity gauge counter; we need to
			// sum across all namespaces.
			go func() {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
				defer cancel()
				entities, err := c.countActiveEntities(ctx)
				if err != nil {
					c.logger.Error("error counting identity entities", "err", err)
				} else {
					metrics.SetGauge([]string{"identity", "num_entities"}, float32(entities.Entities.Total))
				}
			}()
		case <-stopCh:
			return
		}
	}
}

// These wrappers are responsible for redirecting to the current instance of
// TokenStore; there is one per method because an additional level of abstraction
// seems confusing.
func (c *Core) tokenGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	if c.IsDRSecondary() {
		// there is no expiration manager on DR Secondaries
		return []metricsutil.GaugeLabelValues{}, nil
	}

	// stateLock or authLock protects the tokenStore pointer
	c.stateLock.RLock()
	ts := c.tokenStore
	c.stateLock.RUnlock()
	if ts == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil token store")
	}
	return ts.gaugeCollector(ctx)
}

func (c *Core) tokenGaugePolicyCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	if c.IsDRSecondary() {
		// there is no expiration manager on DR Secondaries
		return []metricsutil.GaugeLabelValues{}, nil
	}

	c.stateLock.RLock()
	ts := c.tokenStore
	c.stateLock.RUnlock()
	if ts == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil token store")
	}
	return ts.gaugeCollectorByPolicy(ctx)
}

func (c *Core) leaseExpiryGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	c.stateLock.RLock()
	e := c.expiration
	metricsConsts := c.MetricSink().TelemetryConsts
	c.stateLock.RUnlock()
	if e == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil expiration manager")
	}
	return e.leaseAggregationMetrics(ctx, metricsConsts)
}

func (c *Core) tokenGaugeMethodCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	if c.IsDRSecondary() {
		// there is no expiration manager on DR Secondaries
		return []metricsutil.GaugeLabelValues{}, nil
	}

	c.stateLock.RLock()
	ts := c.tokenStore
	c.stateLock.RUnlock()
	if ts == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil token store")
	}
	return ts.gaugeCollectorByMethod(ctx)
}

func (c *Core) tokenGaugeTtlCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	if c.IsDRSecondary() {
		// there is no expiration manager on DR Secondaries
		return []metricsutil.GaugeLabelValues{}, nil
	}

	c.stateLock.RLock()
	ts := c.tokenStore
	c.stateLock.RUnlock()
	if ts == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil token store")
	}
	return ts.gaugeCollectorByTtl(ctx)
}

// emitMetricsActiveNode is used to start all the periodic metrics; all of them should
// be shut down when stopCh is closed.  This code runs on the active node only.
func (c *Core) emitMetricsActiveNode(stopCh chan struct{}) {
	// The gauge collection processes are started and stopped here
	// because there's more than one TokenManager created during startup,
	// but we only want one set of gauges.
	metricsInit := []struct {
		MetricName    []string
		MetadataLabel []metrics.Label
		CollectorFunc metricsutil.GaugeCollector
		DisableEnvVar string
	}{
		{
			[]string{"token", "count"},
			[]metrics.Label{{"gauge", "token_by_namespace"}},
			c.tokenGaugeCollector,
			"",
		},
		{
			[]string{"token", "count", "by_policy"},
			[]metrics.Label{{"gauge", "token_by_policy"}},
			c.tokenGaugePolicyCollector,
			"",
		},
		{
			[]string{"expire", "leases", "by_expiration"},
			[]metrics.Label{{"gauge", "leases_by_expiration"}},
			c.leaseExpiryGaugeCollector,
			"",
		},
		{
			[]string{"token", "count", "by_auth"},
			[]metrics.Label{{"gauge", "token_by_auth"}},
			c.tokenGaugeMethodCollector,
			"",
		},
		{
			[]string{"token", "count", "by_ttl"},
			[]metrics.Label{{"gauge", "token_by_ttl"}},
			c.tokenGaugeTtlCollector,
			"",
		},
		{
			[]string{"secret", "kv", "count"},
			[]metrics.Label{{"gauge", "kv_secrets_by_mountpoint"}},
			c.kvSecretGaugeCollector,
			"VAULT_DISABLE_KV_GAUGE",
		},
		{
			[]string{"identity", "entity", "count"},
			[]metrics.Label{{"gauge", "identity_by_namespace"}},
			c.entityGaugeCollector,
			"",
		},
		{
			[]string{"identity", "entity", "alias", "count"},
			[]metrics.Label{{"gauge", "identity_by_mountpoint"}},
			c.entityGaugeCollectorByMount,
			"",
		},
		{
			[]string{"identity", "entity", "active", "partial_month"},
			[]metrics.Label{{"gauge", "identity_active_month"}},
			c.activeEntityGaugeCollector,
			"",
		},
		{
			[]string{"policy", "configured", "count"},
			[]metrics.Label{{"gauge", "number_policies_by_type"}},
			c.configuredPoliciesGaugeCollector,
			"",
		},
	}

	// Disable collection if configured, or if we're a performance standby
	// node or DR secondary cluster.
	if c.MetricSink().GaugeInterval == time.Duration(0) {
		c.logger.Info("usage gauge collection is disabled")
	} else if standby, _ := c.Standby(); !standby && !c.IsDRSecondary() {
		for _, init := range metricsInit {
			if init.DisableEnvVar != "" {
				if os.Getenv(init.DisableEnvVar) != "" {
					c.logger.Info("usage gauge collection is disabled for",
						"metric", init.MetricName)
					continue
				}
			}

			proc, err := c.MetricSink().NewGaugeCollectionProcess(
				init.MetricName,
				init.MetadataLabel,
				init.CollectorFunc,
				c.logger,
			)
			if err != nil {
				c.logger.Error("failed to start collector", "metric", init.MetricName, "error", err)
			} else {
				go proc.Run()
				defer proc.Stop()
			}
		}
	}

	// When this returns, all the defers set up above will fire.
	c.metricsLoop(stopCh)
}

type kvMount struct {
	Namespace  *namespace.Namespace
	MountPoint string
	Version    string
	NumSecrets int
}

func (c *Core) findKvMounts() []*kvMount {
	mounts := make([]*kvMount, 0)

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	// we don't grab the statelock, so this code might run during or after the seal process.
	// Therefore, we need to check if c.mounts is nil. If we do not, this will panic when
	// run after seal.
	if c.mounts == nil {
		return mounts
	}

	for _, entry := range c.mounts.Entries {
		if entry.Type == pluginconsts.SecretEngineKV || entry.Type == pluginconsts.SecretEngineGeneric {
			version, ok := entry.Options["version"]
			if !ok || version == "" {
				version = "1"
			}
			mounts = append(mounts, &kvMount{
				Namespace:  entry.namespace,
				MountPoint: entry.Path,
				Version:    version,
				NumSecrets: 0,
			})
		}
	}
	return mounts
}

func (c *Core) kvCollectionErrorCount() {
	c.MetricSink().IncrCounterWithLabels(
		[]string{"metrics", "collection", "error"},
		1,
		[]metrics.Label{{"gauge", "kv_secrets_by_mountpoint"}},
	)
}

func (c *Core) walkKvMountSecrets(ctx context.Context, m *kvMount) {
	var subdirectories []string
	if m.Version == "1" {
		subdirectories = []string{m.Namespace.Path + m.MountPoint}
	} else {
		subdirectories = []string{m.Namespace.Path + m.MountPoint + "metadata/"}
	}

	for len(subdirectories) > 0 {
		// Check for cancellation
		select {
		case <-ctx.Done():
			return
		default:
			break
		}

		currentDirectory := subdirectories[0]
		subdirectories = subdirectories[1:]

		listRequest := &logical.Request{
			Operation: logical.ListOperation,
			Path:      currentDirectory,
		}
		resp, err := c.router.Route(ctx, listRequest)
		if err != nil {
			c.kvCollectionErrorCount()
			// ErrUnsupportedPath probably means that the mount is not there anymore,
			// don't log those cases.
			if !strings.Contains(err.Error(), logical.ErrUnsupportedPath.Error()) &&
				// ErrSetupReadOnly means the mount's currently being set up.
				// Nothing is wrong and there's no cause for alarm, just that we can't get data from it
				// yet. We also shouldn't log these cases
				!strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
				c.logger.Error("failed to perform internal KV list", "mount_point", m.MountPoint, "error", err)
				break
			}
			// Quit handling this mount point (but it'll still appear in the list)
			return
		}
		if resp == nil {
			continue
		}
		rawKeys, ok := resp.Data["keys"]
		if !ok {
			continue
		}
		keys, ok := rawKeys.([]string)
		if !ok {
			c.kvCollectionErrorCount()
			c.logger.Error("KV list keys are not a []string", "mount_point", m.MountPoint, "rawKeys", rawKeys)
			// Quit handling this mount point (but it'll still appear in the list)
			return
		}
		for _, path := range keys {
			if len(path) > 0 && path[len(path)-1] == '/' {
				subdirectories = append(subdirectories, currentDirectory+path)
			} else {
				m.NumSecrets += 1
			}
		}
	}
}

// GetTotalPkiRoles returns the total roles across all PKI mounts in Vault
func (c *Core) GetTotalPkiRoles(ctx context.Context) int {
	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	numRoles := 0

	for _, entry := range c.mounts.Entries {
		secretType := entry.Type
		if secretType == pluginconsts.SecretEnginePki {
			listRequest := &logical.Request{
				Operation: logical.ListOperation,
				Path:      entry.namespace.Path + entry.Path + "roles",
			}
			resp, err := c.router.Route(ctx, listRequest)
			if err != nil || resp == nil {
				continue
			}
			rawKeys, ok := resp.Data["keys"]
			if !ok {
				continue
			}
			keys, ok := rawKeys.([]string)
			if ok {
				numRoles += len(keys)
			}
		}
	}
	return numRoles
}

// GetTotalPkiIssuers returns the total issuers across all PKI mounts in Vault
func (c *Core) GetTotalPkiIssuers(ctx context.Context) int {
	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	numRoles := 0

	for _, entry := range c.mounts.Entries {
		secretType := entry.Type
		if secretType == pluginconsts.SecretEnginePki {
			listRequest := &logical.Request{
				Operation: logical.ListOperation,
				Path:      entry.namespace.Path + entry.Path + "issuers",
			}
			resp, err := c.router.Route(ctx, listRequest)
			if err != nil || resp == nil {
				continue
			}
			rawKeys, ok := resp.Data["keys"]
			if !ok {
				continue
			}
			keys, ok := rawKeys.([]string)
			if ok {
				numRoles += len(keys)
			}
		}
	}
	return numRoles
}

// getMinNamespaceSecrets is expected to be called on the output
// of GetKvUsageMetrics to get the min number of secrets in a single namespace.
func getMinNamespaceSecrets(mapOfNamespacesToSecrets map[string]int) int {
	currentMin := 0
	for _, n := range mapOfNamespacesToSecrets {
		if n < currentMin || currentMin == 0 {
			currentMin = n
		}
	}
	return currentMin
}

// getMaxNamespaceSecrets is expected to be called on the output
// of GetKvUsageMetrics to get the max number of secrets in a single namespace.
func getMaxNamespaceSecrets(mapOfNamespacesToSecrets map[string]int) int {
	currentMax := 0
	for _, n := range mapOfNamespacesToSecrets {
		if n > currentMax {
			currentMax = n
		}
	}
	return currentMax
}

// getTotalSecretsAcrossAllNamespaces is expected to be called on the output
// of GetKvUsageMetrics to get the total number of secrets across namespaces.
func getTotalSecretsAcrossAllNamespaces(mapOfNamespacesToSecrets map[string]int) int {
	total := 0
	for _, n := range mapOfNamespacesToSecrets {
		total += n
	}
	return total
}

// getMeanNamespaceSecrets is expected to be called on the output
// of GetKvUsageMetrics to get the mean number of secrets across namespaces.
func getMeanNamespaceSecrets(mapOfNamespacesToSecrets map[string]int) int {
	length := len(mapOfNamespacesToSecrets)
	// Avoid divide by zero:
	if length == 0 {
		return length
	}
	return getTotalSecretsAcrossAllNamespaces(mapOfNamespacesToSecrets) / length
}

// GetSecretEngineUsageMetrics returns a map of secret engine mount types to the number of those mounts that exist.
func (c *Core) GetSecretEngineUsageMetrics() map[string]int {
	mounts := make(map[string]int)

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	for _, entry := range c.mounts.Entries {
		mountType := entry.Type
		if _, ok := mounts[mountType]; !ok {
			mounts[mountType] = 1
		} else {
			mounts[mountType] += 1
		}
	}
	return mounts
}

// GetAuthMethodUsageMetrics returns a map of auth mount types to the number of those mounts that exist.
func (c *Core) GetAuthMethodUsageMetrics() map[string]int {
	mounts := make(map[string]int)

	c.authLock.RLock()
	defer c.authLock.RUnlock()

	for _, entry := range c.auth.Entries {
		authType := entry.Type
		if _, ok := mounts[authType]; !ok {
			mounts[authType] = 1
		} else {
			mounts[authType] += 1
		}
	}
	return mounts
}

// GetKvUsageMetrics returns a map of namespace paths to KV secret counts within those namespaces.
func (c *Core) GetKvUsageMetrics(ctx context.Context, kvVersion string) (map[string]int, error) {
	mounts := c.findKvMounts()
	results := make(map[string]int)

	if kvVersion == "1" || kvVersion == "2" {
		var newMounts []*kvMount
		for _, mount := range mounts {
			if mount.Version == kvVersion {
				newMounts = append(newMounts, mount)
			}
		}
		mounts = newMounts
	} else if kvVersion != "0" {
		return results, fmt.Errorf("kv version %s not supported, must be 0, 1, or 2", kvVersion)
	}

	for _, m := range mounts {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("context expired")
		default:
			break
		}

		c.walkKvMountSecrets(ctx, m)

		_, ok := results[m.Namespace.Path]
		if ok {
			// we need to add, not overwrite
			results[m.Namespace.Path] += m.NumSecrets
		} else {
			results[m.Namespace.Path] = m.NumSecrets
		}
	}

	return results, nil
}

func (c *Core) kvSecretGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	// Find all KV mounts
	mounts := c.findKvMounts()
	results := make([]metricsutil.GaugeLabelValues, len(mounts))

	// Use a root namespace, so include namespace path
	// in any queries.
	ctx = namespace.RootContext(ctx)

	// Route list requests to all the identified mounts.
	// (All of these will show up as activity in the vault.route metric.)
	// Then we have to explore each subdirectory.
	for i, m := range mounts {
		// Check for cancellation, return empty array
		select {
		case <-ctx.Done():
			return []metricsutil.GaugeLabelValues{}, nil
		default:
			break
		}

		results[i].Labels = []metrics.Label{
			metricsutil.NamespaceLabel(m.Namespace),
			{"mount_point", m.MountPoint},
		}

		c.walkKvMountSecrets(ctx, m)
		results[i].Value = float32(m.NumSecrets)
	}

	return results, nil
}

func (c *Core) entityGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	// Protect against concurrent changes during seal
	c.stateLock.RLock()
	identityStore := c.identityStore
	c.stateLock.RUnlock()
	if identityStore == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil identity store")
	}

	byNamespace, err := identityStore.countEntitiesByNamespace(ctx)
	if err != nil {
		return []metricsutil.GaugeLabelValues{}, err
	}

	// No check for expiration here; the bulk of the work should be in
	// counting the entities.
	allNamespaces := c.collectNamespaces()
	values := make([]metricsutil.GaugeLabelValues, len(allNamespaces))
	for i := range values {
		values[i].Labels = []metrics.Label{
			metricsutil.NamespaceLabel(allNamespaces[i]),
		}
		values[i].Value = float32(byNamespace[allNamespaces[i].ID])
	}

	return values, nil
}

func (c *Core) entityGaugeCollectorByMount(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	c.stateLock.RLock()
	identityStore := c.identityStore
	c.stateLock.RUnlock()
	if identityStore == nil {
		return []metricsutil.GaugeLabelValues{}, errors.New("nil identity store")
	}

	byAccessor, err := identityStore.countEntitiesByMountAccessor(ctx)
	if err != nil {
		return []metricsutil.GaugeLabelValues{}, err
	}

	values := make([]metricsutil.GaugeLabelValues, 0)
	for accessor, count := range byAccessor {
		// Terminate if taking too long to do the translation
		select {
		case <-ctx.Done():
			return values, errors.New("context cancelled")
		default:
			break
		}

		c.stateLock.RLock()
		mountEntry := c.router.MatchingMountByAccessor(accessor)
		c.stateLock.RUnlock()
		if mountEntry == nil {
			continue
		}
		values = append(values, metricsutil.GaugeLabelValues{
			Labels: []metrics.Label{
				metricsutil.NamespaceLabel(mountEntry.namespace),
				{"auth_method", mountEntry.Type},
				{"mount_point", "auth/" + mountEntry.Path},
			},
			Value: float32(count),
		})
	}

	return values, nil
}

func (c *Core) cachedGaugeMetricsEmitter() {
	if c.metricsHelper == nil {
		return
	}

	loopMetrics := &c.metricsHelper.LoopMetrics.Metrics

	emit := func(key interface{}, value interface{}) bool {
		metricValue := value.(metricsutil.GaugeMetric)
		c.metricSink.SetGaugeWithLabels(metricValue.Key, metricValue.Value, metricValue.Labels)
		return true
	}

	loopMetrics.Range(emit)
}

func (c *Core) inFlightReqGaugeMetric() {
	totalInFlightReq := c.inFlightReqData.InFlightReqCount.Load()
	// Adding a gauge metric to capture total number of inflight requests
	c.metricSink.SetGaugeWithLabels([]string{"core", "in_flight_requests"}, float32(totalInFlightReq), nil)
}

// configuredPoliciesGaugeCollector is used to collect gauge label values for the `vault.policy.configured.count` metric
func (c *Core) configuredPoliciesGaugeCollector(ctx context.Context) ([]metricsutil.GaugeLabelValues, error) {
	c.stateLock.RLock()
	policyStore := c.policyStore
	c.stateLock.RUnlock()

	if policyStore == nil {
		return []metricsutil.GaugeLabelValues{}, nil
	}

	ctx = namespace.RootContext(ctx)
	namespaces := c.collectNamespaces()

	policyTypes := []PolicyType{
		PolicyTypeACL,
		PolicyTypeRGP,
		PolicyTypeEGP,
	}
	var values []metricsutil.GaugeLabelValues

	for _, pt := range policyTypes {
		policies, err := policyStore.policiesByNamespaces(ctx, pt, namespaces)
		if err != nil {
			return []metricsutil.GaugeLabelValues{}, err
		}

		v := metricsutil.GaugeLabelValues{}
		v.Labels = []metricsutil.Label{{
			"policy_type",
			pt.String(),
		}}
		v.Value = float32(len(policies))
		values = append(values, v)
	}

	return values, nil
}
