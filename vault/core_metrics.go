// Copyright IBM Corp. 2016, 2025
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
	"github.com/hashicorp/vault/sdk/helper/pluginutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	KVv2MetadataPath = "metadata"
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
		MetricName       []string
		MetadataLabel    []metrics.Label
		CollectorFunc    metricsutil.GaugeCollector
		DisableEnvVar    string
		IsEnterpriseOnly bool
	}{
		{
			[]string{"token", "count"},
			[]metrics.Label{{"gauge", "token_by_namespace"}},
			c.tokenGaugeCollector,
			"",
			false,
		},
		{
			[]string{"token", "count", "by_policy"},
			[]metrics.Label{{"gauge", "token_by_policy"}},
			c.tokenGaugePolicyCollector,
			"",
			false,
		},
		{
			[]string{"expire", "leases", "by_expiration"},
			[]metrics.Label{{"gauge", "leases_by_expiration"}},
			c.leaseExpiryGaugeCollector,
			"",
			false,
		},
		{
			[]string{"token", "count", "by_auth"},
			[]metrics.Label{{"gauge", "token_by_auth"}},
			c.tokenGaugeMethodCollector,
			"",
			false,
		},
		{
			[]string{"token", "count", "by_ttl"},
			[]metrics.Label{{"gauge", "token_by_ttl"}},
			c.tokenGaugeTtlCollector,
			"",
			false,
		},
		{
			[]string{"secret", "kv", "count"},
			[]metrics.Label{{"gauge", "kv_secrets_by_mountpoint"}},
			c.kvSecretGaugeCollector,
			"VAULT_DISABLE_KV_GAUGE",
			false,
		},
		{
			[]string{"identity", "entity", "count"},
			[]metrics.Label{{"gauge", "identity_by_namespace"}},
			c.entityGaugeCollector,
			"",
			false,
		},
		{
			[]string{"identity", "entity", "alias", "count"},
			[]metrics.Label{{"gauge", "identity_by_mountpoint"}},
			c.entityGaugeCollectorByMount,
			"",
			false,
		},
		{
			[]string{"identity", "entity", "active", "partial_month"},
			[]metrics.Label{{"gauge", "identity_active_month"}},
			c.activeEntityGaugeCollector,
			"",
			false,
		},
		{
			[]string{"policy", "configured", "count"},
			[]metrics.Label{{"gauge", "number_policies_by_type"}},
			c.configuredPoliciesGaugeCollector,
			"",
			false,
		},
		{
			[]string{"client", "billing_period", "activity"},
			[]metrics.Label{{"gauge", "clients_current_billing_period"}},
			c.clientsGaugeCollectorCurrentBillingPeriod,
			"",
			true,
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

			// Billing start date is always 0 on CE
			if init.IsEnterpriseOnly && c.BillingStart().IsZero() {
				continue
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
	Namespace            *namespace.Namespace
	MountPoint           string
	MountAccessor        string
	Version              string
	Local                bool
	NumSecrets           int
	RunningPluginVersion string
}

// findOfficialKvMounts differs from findKvMounts in that it will ignore any sideloaded
// or externally compiled KV mounts that are still of type KV.
// It's a simple function that's slightly reimplemented to prevent needing a context
// in findKvMounts.
func (c *Core) findOfficialKvMounts(ctx context.Context) []*kvMount {
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

			pluginName := getAdjustedPluginType(entry)
			if pluginName == "" {
				continue
			}

			pluginVersion := entry.RunningVersion
			runner, err := c.pluginCatalog.Get(ctx, pluginName, consts.PluginTypeSecrets, pluginVersion)
			if err != nil {
				continue
			}

			if !(isOfficialOrBuiltin(runner)) {
				continue
			}

			mounts = append(mounts, &kvMount{
				Namespace:            entry.namespace,
				MountPoint:           entry.Path,
				MountAccessor:        entry.Accessor,
				Version:              version,
				NumSecrets:           0,
				Local:                entry.Local,
				RunningPluginVersion: entry.RunningVersion,
			})
		}
	}
	return mounts
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
				Namespace:            entry.namespace,
				MountPoint:           entry.Path,
				MountAccessor:        entry.Accessor,
				Version:              version,
				NumSecrets:           0,
				Local:                entry.Local,
				RunningPluginVersion: entry.RunningVersion,
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

func (c *Core) walkKvSecrets(
	ctx context.Context,
	rootDirs []string,
	m *kvMount,
	onSecret func(ctx context.Context, fullPath string) error,
) error {
	subdirectories := rootDirs

	for len(subdirectories) > 0 {
		// Context cancellation check
		select {
		case <-ctx.Done():
			return nil
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
			return err
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
			return fmt.Errorf("KV list keys are not a []string")
		}

		for _, path := range keys {
			fullPath := currentDirectory + path
			if strings.HasSuffix(path, "/") {
				subdirectories = append(subdirectories, fullPath)
			} else {
				if callBackErr := onSecret(ctx, fullPath); callBackErr != nil {
					c.logger.Error("failed to get metadata for KVv2 secret", "path", fullPath, "error", err)
					return callBackErr
				}
			}
		}
	}
	return nil
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

func (c *Core) walkKvMountSecrets(ctx context.Context, m *kvMount) {
	var startDirs []string
	if m.Version == "1" {
		startDirs = []string{m.Namespace.Path + m.MountPoint}
	} else {
		startDirs = []string{m.Namespace.Path + m.MountPoint + KVv2MetadataPath + "/"}
	}

	err := c.walkKvSecrets(ctx, startDirs, m, func(ctx context.Context, fullPath string) error {
		m.NumSecrets++
		return nil
	})
	if err != nil {
		// ErrUnsupportedPath probably means that the mount is not there anymore,
		// don't log those cases.
		if !strings.Contains(err.Error(), logical.ErrUnsupportedPath.Error()) &&
			// ErrSetupReadOnly means the mount's currently being set up.
			// Nothing is wrong and there's no cause for alarm, just that we can't get data from it
			// yet. We also shouldn't log these cases
			!strings.Contains(err.Error(), logical.ErrSetupReadOnly.Error()) {
			c.logger.Error("failed to walk KV mount", "mount_point", m.MountPoint, "error", err)
		}
	}
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

type RoleCounts struct {
	AWSDynamicRoles            int `json:"aws_dynamic_roles"`
	AWSStaticRoles             int `json:"aws_static_roles"`
	AzureDynamicRoles          int `json:"azure_dynamic_roles"`
	AzureStaticRoles           int `json:"azure_static_roles"`
	DatabaseDynamicRoles       int `json:"database_dynamic_roles"`
	DatabaseStaticRoles        int `json:"database_static_roles"`
	GCPRolesets                int `json:"gcp_rolesets"`
	GCPStaticAccounts          int `json:"gcp_static_accounts"`
	GCPImpersonatedAccounts    int `json:"gcp_impersonated_accounts"`
	LDAPDynamicRoles           int `json:"ldap_dynamic_roles"`
	LDAPStaticRoles            int `json:"ldap_static_roles"`
	OpenLDAPDynamicRoles       int `json:"openldap_dynamic_roles"`
	OpenLDAPStaticRoles        int `json:"openldap_static_roles"`
	AlicloudDynamicRoles       int `json:"alicloud_dynamic_roles"`
	RabbitMQDynamicRoles       int `json:"rabbitmq_dynamic_roles"`
	ConsulDynamicRoles         int `json:"consul_dynamic_roles"`
	NomadDynamicRoles          int `json:"nomad_dynamic_roles"`
	KubernetesDynamicRoles     int `json:"kubernetes_dynamic_roles"`
	MongoDBAtlasDynamicRoles   int `json:"mongodb_atlas_dynamic_roles"`
	TerraformCloudDynamicRoles int `json:"terraformcloud_dynamic_roles"`
}

type ManagedKeyCounts struct {
	TotpKeys int `json:"totp_keys"`
}

// getRoleAndManagedKeyCountsInternal gets the role counts for plugins and managed key counts
// includeLocal determines if local mounts are included
// includeReplicated determines if replicated mounts are included
// officialPluginsOnly determines if this function should include only plugins that are official,
// which would exclude, for example, a custom built version of these plugins.
func (c *Core) getRoleAndManagedKeyCountsInternal(includeLocal bool, includeReplicated bool, officialPluginsOnly bool) (*RoleCounts, *ManagedKeyCounts) {
	if c.Sealed() {
		c.logger.Debug("core is sealed, cannot access mounts table")
		return nil, nil
	}

	ctx := namespace.RootContext(c.activeContext)
	apiList := func(entry *MountEntry, apiPath string) []string {
		listRequest := &logical.Request{
			Operation: logical.ListOperation,
			Path:      entry.namespace.Path + entry.Path + apiPath,
		}

		resp, err := c.router.Route(ctx, listRequest)
		if err != nil || resp == nil {
			return nil
		}
		rawKeys, ok := resp.Data["keys"]
		if !ok {
			return nil
		}
		keys, ok := rawKeys.([]string)
		if !ok {
			return nil
		}
		return keys
	}

	c.mountsLock.RLock()
	defer c.mountsLock.RUnlock()

	var roles RoleCounts
	var keyCounts ManagedKeyCounts
	for _, entry := range c.mounts.Entries {
		if !entry.Local && !includeReplicated {
			continue
		}
		if entry.Local && !includeLocal {
			continue
		}

		pluginName := getAdjustedPluginType(entry)
		if pluginName == "" {
			continue
		}

		pluginVersion := entry.RunningVersion

		if officialPluginsOnly {
			runner, err := c.pluginCatalog.Get(ctx, pluginName, consts.PluginTypeSecrets, pluginVersion)
			if err != nil {
				continue
			}

			if !(isOfficialOrBuiltin(runner)) {
				continue
			}
		}

		switch pluginName {
		case pluginconsts.SecretEngineAWS:
			dynamicRoles := apiList(entry, "roles")
			roles.AWSDynamicRoles += len(dynamicRoles)
			staticRoles := apiList(entry, "static-roles")
			roles.AWSStaticRoles += len(staticRoles)

		case pluginconsts.SecretEngineAzure:
			dynamicRoles := apiList(entry, "roles")
			roles.AzureDynamicRoles += len(dynamicRoles)
			staticRoles := apiList(entry, "static-roles")
			roles.AzureStaticRoles += len(staticRoles)

		case pluginconsts.SecretEngineDatabase:
			dynamicRoles := apiList(entry, "roles")
			roles.DatabaseDynamicRoles += len(dynamicRoles)
			staticRoles := apiList(entry, "static-roles")
			roles.DatabaseStaticRoles += len(staticRoles)

		case pluginconsts.SecretEngineGCP:
			rolesets := apiList(entry, "rolesets")
			roles.GCPRolesets += len(rolesets)
			staticAccounts := apiList(entry, "static-accounts")
			roles.GCPStaticAccounts += len(staticAccounts)
			impersonatedAccounts := apiList(entry, "impersonated-accounts")
			roles.GCPImpersonatedAccounts += len(impersonatedAccounts)

		case pluginconsts.SecretEngineLDAP:
			dynamicRoles := apiList(entry, "role")
			roles.LDAPDynamicRoles += len(dynamicRoles)
			staticRoles := apiList(entry, "static-role")
			roles.LDAPStaticRoles += len(staticRoles)

		case pluginconsts.SecretEngineOpenLDAP:
			dynamicRoles := apiList(entry, "role")
			roles.OpenLDAPDynamicRoles += len(dynamicRoles)
			staticRoles := apiList(entry, "static-role")
			roles.OpenLDAPStaticRoles += len(staticRoles)

		case pluginconsts.SecretEngineAlicloud:
			dynamicRoles := apiList(entry, "role")
			roles.AlicloudDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineRabbitMQ:
			dynamicRoles := apiList(entry, "roles")
			roles.RabbitMQDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineConsul:
			dynamicRoles := apiList(entry, "roles")
			roles.ConsulDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineNomad:
			dynamicRoles := apiList(entry, "role")
			roles.NomadDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineKubernetes:
			dynamicRoles := apiList(entry, "roles")
			roles.KubernetesDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineMongoDBAtlas:
			dynamicRoles := apiList(entry, "roles")
			roles.MongoDBAtlasDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineTerraform:
			dynamicRoles := apiList(entry, "role")
			roles.TerraformCloudDynamicRoles += len(dynamicRoles)

		case pluginconsts.SecretEngineTOTP:
			keyCountPerEntry := apiList(entry, "keys")
			keyCounts.TotpKeys += len(keyCountPerEntry)
		}
	}

	return &roles, &keyCounts
}

func (c *Core) GetRoleCounts() *RoleCounts {
	roleCounts, _ := c.getRoleAndManagedKeyCountsInternal(true, true, false)
	return roleCounts
}

func (c *Core) GetRoleCountsForCluster() *RoleCounts {
	roleCounts, _ := c.getRoleAndManagedKeyCountsInternal(true, c.isPrimary(), false)
	return roleCounts
}

// GetKvUsageMetrics returns a map of namespace paths to KV secret counts.
func (c *Core) GetKvUsageMetrics(ctx context.Context, kvVersion string) (map[string]int, error) {
	return c.GetKvUsageMetricsByNamespace(ctx, kvVersion, "", true, true, true)
}

// GetKvUsageMetricsByNamespace returns a map of namespace paths to KV secret counts within a specific namespace.
func (c *Core) GetKvUsageMetricsByNamespace(ctx context.Context, kvVersion string, nsPath string, includeLocal bool, includeReplicated bool, includeUnofficial bool) (map[string]int, error) {
	mounts := c.findKvMounts()
	if !includeUnofficial {
		mounts = c.findOfficialKvMounts(ctx)
	}
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
		if !includeLocal && m.Local {
			continue
		}
		if !includeReplicated && !m.Local {
			continue
		}

		if nsPath != "" && !strings.HasPrefix(m.Namespace.Path, nsPath) {
			continue
		}

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

// isOfficialOrBuiltin determines if a plugin is official based on its runner.
// We treat it as official if runner is nil to avoid overcharging, but ensure
// that it is properly scanned if it _is_ an official mount.
func isOfficialOrBuiltin(runner *pluginutil.PluginRunner) bool {
	return runner == nil || runner.Builtin || runner.Tier == consts.PluginTierOfficial
}

// ListOfficialAndExternalSecretPlugins gets a list of all secret plugins, official and external.
// The union of both sets is the set of all secret plugins.
// Returns a list of official plugins, external plugins, and error, in that order.
func (c *Core) ListOfficialAndExternalSecretPlugins(ctx context.Context) ([]*MountEntry, []*MountEntry, error) {
	if c == nil || c.pluginCatalog == nil {
		return nil, nil, fmt.Errorf("core or plugin catalog is nil")
	}

	mounts, err := c.ListMounts()
	if err != nil {
		return nil, nil, fmt.Errorf("error listing mounts: %w", err)
	}

	var official []*MountEntry
	var external []*MountEntry
	for _, entry := range mounts {
		if entry == nil {
			continue
		}

		// Only secrets-engine mounts live in the mounts table. Exclude the known
		// non-secrets mounts and database mounts (PluginTypeDatabase).
		if entry.Table != mountTableType {
			continue
		}

		pluginName := getAdjustedPluginType(entry)
		if pluginName == "" {
			continue
		}

		pluginVersion := entry.RunningVersion

		runner, err := c.pluginCatalog.Get(ctx, pluginName, consts.PluginTypeSecrets, pluginVersion)
		if err != nil {
			continue
		}

		if isOfficialOrBuiltin(runner) {
			official = append(official, entry)
		} else {
			external = append(external, entry)
		}
	}

	return official, external, nil
}

// ListOfficialSecretPlugins gets a list of all 'official'/builtin secret plugins.
func (c *Core) ListOfficialSecretPlugins(ctx context.Context) ([]*MountEntry, error) {
	internalPlugins, _, err := c.ListOfficialAndExternalSecretPlugins(ctx)
	if err != nil {
		return nil, err
	}
	return internalPlugins, nil
}

// getAdjustedPluginType gets the adjusted plugin type for an entry. In most cases
// this will be entry.Type, but it will correctly return the type for legacy (pre-Vault 1.0) plugins.
func getAdjustedPluginType(entry *MountEntry) string {
	if entry == nil {
		return ""
	}
	pluginName := entry.Type
	if pluginName == mountTypePlugin && entry.Config.PluginName != "" {
		pluginName = entry.Config.PluginName
	}
	return pluginName
}

// ListDeduplicatedExternalSecretPlugins returns the enabled secret engines
// that are not builtin and not official-tier.
//
// This is useful for identifying "third-party" secrets mounts (e.g. community or
// partner tier external plugins) while excluding builtins and official HashiCorp
// plugins.
// Note: This will include all mounts that have been built externally (even if they are
// Hashicorp owned). This will happen if the plugin was built from a Github repo or from an
// artifact.
func (c *Core) ListDeduplicatedExternalSecretPlugins(ctx context.Context) ([]*MountEntry, error) {
	_, externalPlugins, err := c.ListOfficialAndExternalSecretPlugins(ctx)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]struct{})
	var result []*MountEntry
	for _, entry := range externalPlugins {
		if entry == nil {
			continue
		}

		pluginName := getAdjustedPluginType(entry)
		if pluginName == "" {
			continue
		}

		pluginVersion := entry.RunningVersion

		// De-dupe: multiple mounts can point at the same underlying plugin+version.
		// We want to charge for each unique plugin+version pair.
		key := pluginName + "\x00" + pluginVersion
		if _, ok := seen[key]; ok {
			continue
		}

		result = append(result, entry)
		seen[key] = struct{}{}
	}

	return result, nil
}
