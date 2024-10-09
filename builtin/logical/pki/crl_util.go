// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/builtin/logical/pki/pki_backend"
	"github.com/hashicorp/vault/builtin/logical/pki/revocation"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/errutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	revokedPath                      = revocation.RevokedPath
	crossRevocationPrefix            = "cross-revocation-queue/"
	crossRevocationPath              = crossRevocationPrefix + "{{clusterId}}/"
	deltaWALLastBuildSerialName      = "last-build-serial"
	deltaWALLastRevokedSerialName    = "last-revoked-serial"
	localDeltaWALPath                = "delta-wal/"
	localDeltaWALLastBuildSerial     = localDeltaWALPath + deltaWALLastBuildSerialName
	localDeltaWALLastRevokedSerial   = localDeltaWALPath + deltaWALLastRevokedSerialName
	unifiedDeltaWALPrefix            = "unified-delta-wal/"
	unifiedDeltaWALPath              = "unified-delta-wal/{{clusterId}}/"
	unifiedDeltaWALLastBuildSerial   = unifiedDeltaWALPath + deltaWALLastBuildSerialName
	unifiedDeltaWALLastRevokedSerial = unifiedDeltaWALPath + deltaWALLastRevokedSerialName
)

type revocationRequest struct {
	RequestedAt time.Time `json:"requested_at"`
}

type revocationConfirmed struct {
	RevokedAt string `json:"revoked_at"`
	Source    string `json:"source"`
}

type revocationQueueEntry struct {
	Cluster string
	Serial  string
}

type (
	// Placeholder in case of migrations needing more data. Currently
	// we use the path name to store the serial number that was revoked.
	deltaWALInfo struct{}
	lastWALInfo  struct {
		// Info to write about the last WAL entry. This is the serial number
		// of the last revoked certificate.
		//
		// We write this below in revokedCert(...) and read it in
		// rebuildDeltaCRLsIfForced(...).
		Serial string `json:"serial"`
	}
	lastDeltaInfo struct {
		// Info to write about the last built delta CRL. This is the serial
		// number of the last revoked certificate that we saw prior to delta
		// CRL building.
		//
		// We write this below in buildAnyCRLs(...) and read it in
		// rebuildDeltaCRLsIfForced(...).
		Serial string `json:"serial"`
	}
)

// CrlBuilder is gatekeeper for controlling various read/write operations to the storage of the CRL.
// The extra complexity arises from secondary performance clusters seeing various writes to its storage
// without the actual API calls. During the storage invalidation process, we do not have the required state
// to actually rebuild the CRLs, so we need to schedule it in a deferred fashion. This allows either
// read or write calls to perform the operation if required, or have the flag reset upon a write operation
//
// The CRL builder also tracks the revocation configuration.
type CrlBuilder struct {
	_builder              sync.Mutex
	forceRebuild          *atomic.Bool
	canRebuild            bool
	lastDeltaRebuildCheck time.Time

	_config               sync.RWMutex
	dirty                 *atomic.Bool
	config                pki_backend.CrlConfig
	haveInitializedConfig bool

	// Whether to invalidate our LastModifiedTime due to write on the
	// global issuance config.
	invalidate *atomic.Bool

	// Global revocation queue entries get accepted by the invalidate func
	// and passed to the CrlBuilder for processing.
	haveInitializedQueue *atomic.Bool
	revQueue             *revocationQueue
	removalQueue         *revocationQueue
	crossQueue           *revocationQueue
}

var _ pki_backend.CrlBuilderType = (*CrlBuilder)(nil)

const (
	_ignoreForceFlag  = true
	_enforceForceFlag = false
)

func newCRLBuilder(canRebuild bool) *CrlBuilder {
	builder := &CrlBuilder{
		forceRebuild: &atomic.Bool{},
		canRebuild:   canRebuild,
		// Set the last delta rebuild window to now, delaying the first delta
		// rebuild by the first rebuild period to give us some time on startup
		// to stabilize.
		lastDeltaRebuildCheck: time.Now(),
		dirty:                 &atomic.Bool{},
		config:                pki_backend.DefaultCrlConfig,
		invalidate:            &atomic.Bool{},
		haveInitializedQueue:  &atomic.Bool{},
		revQueue:              newRevocationQueue(),
		removalQueue:          newRevocationQueue(),
		crossQueue:            newRevocationQueue(),
	}
	builder.dirty.Store(true)
	return builder
}

func (cb *CrlBuilder) SetLastDeltaRebuildCheckTime(t time.Time) {
	cb.lastDeltaRebuildCheck = t
}

func (cb *CrlBuilder) ShouldInvalidate() bool {
	return cb.invalidate.Load()
}

func (cb *CrlBuilder) markConfigDirty() {
	cb.dirty.Store(true)
}

func (cb *CrlBuilder) reloadConfigIfRequired(sc *storageContext) error {
	if cb.dirty.Load() {
		// Acquire a write lock.
		cb._config.Lock()
		defer cb._config.Unlock()

		if !cb.dirty.Load() {
			// Someone else might've been reloading the config; no need
			// to do it twice.
			return nil
		}

		config, err := sc.getRevocationConfig()
		if err != nil {
			return err
		}

		previousConfig := cb.config
		// Set the default config if none was returned to us.
		if config != nil {
			cb.config = *config
		} else {
			cb.config = pki_backend.DefaultCrlConfig
		}

		// Updated the config; unset dirty.
		cb.dirty.Store(false)
		triggerChangeNotification := true
		if !cb.haveInitializedConfig {
			cb.haveInitializedConfig = true
			triggerChangeNotification = false // do not trigger on the initial loading of configuration.
		}

		// Certain things need to be triggered on all server types when CrlConfig is loaded.
		if triggerChangeNotification {
			cb.notifyOnConfigChange(sc, previousConfig, cb.config)
		}
	}

	return nil
}

func (cb *CrlBuilder) notifyOnConfigChange(sc *storageContext, priorConfig pki_backend.CrlConfig, newConfig pki_backend.CrlConfig) {
	// If you need to hook into a CRL configuration change across different server types
	// such as primary clusters as well as performance replicas, it is easier to do here than
	// in two places (API layer and in invalidateFunc)
	if priorConfig.UnifiedCRL != newConfig.UnifiedCRL && newConfig.UnifiedCRL {
		sc.GetUnifiedTransferStatus().forceRun()
	}

	if priorConfig.UseGlobalQueue != newConfig.UseGlobalQueue && newConfig.UseGlobalQueue {
		cb.haveInitializedQueue.Store(false)
	}
}

func (cb *CrlBuilder) GetConfigWithUpdate(sc pki_backend.StorageContext) (*pki_backend.CrlConfig, error) {
	// Config may mutate immediately after accessing, but will be freshly
	// fetched if necessary.
	if err := cb.reloadConfigIfRequired(sc.(*storageContext)); err != nil {
		return nil, err
	}

	cb._config.RLock()
	defer cb._config.RUnlock()

	configCopy := cb.config
	return &configCopy, nil
}

func (cb *CrlBuilder) getConfigWithForcedUpdate(sc *storageContext) (*pki_backend.CrlConfig, error) {
	cb.markConfigDirty()
	return cb.GetConfigWithUpdate(sc)
}

func (cb *CrlBuilder) writeConfig(sc *storageContext, config *pki_backend.CrlConfig) (*pki_backend.CrlConfig, error) {
	cb._config.Lock()
	defer cb._config.Unlock()

	if err := sc.setRevocationConfig(config); err != nil {
		cb.markConfigDirty()
		return nil, fmt.Errorf("failed writing CRL config: %w", err)
	}

	previousConfig := cb.config
	if config != nil {
		cb.config = *config
	} else {
		cb.config = pki_backend.DefaultCrlConfig
	}

	triggerChangeNotification := true
	if !cb.haveInitializedConfig {
		cb.haveInitializedConfig = true
		triggerChangeNotification = false // do not trigger on the initial loading of configuration.
	}

	// Certain things need to be triggered on all server types when CrlConfig is loaded.
	if triggerChangeNotification {
		cb.notifyOnConfigChange(sc, previousConfig, cb.config)
	}

	return config, nil
}

func (cb *CrlBuilder) checkForAutoRebuild(sc *storageContext) error {
	cfg, err := cb.GetConfigWithUpdate(sc)
	if err != nil {
		return err
	}

	if cfg.Disable || !cfg.AutoRebuild || cb.forceRebuild.Load() {
		// Not enabled, not on auto-rebuilder, or we're already scheduled to
		// rebuild so there's no point to interrogate CRL values...
		return nil
	}

	// Auto-Rebuild is enabled. We need to check each issuer's CRL and see
	// if its about to expire. If it is, we've gotta rebuild it (and well,
	// every other CRL since we don't have a fine-toothed rebuilder).
	//
	// We store a list of all (unique) CRLs in the cluster-local CRL
	// configuration along with their expiration dates.
	internalCRLConfig, err := sc.getLocalCRLConfig()
	if err != nil {
		return fmt.Errorf("error checking for auto-rebuild status: unable to fetch cluster-local CRL configuration: %w", err)
	}

	// If there's no config, assume we've gotta rebuild it to get this
	// information.
	if internalCRLConfig == nil {
		cb.forceRebuild.Store(true)
		return nil
	}

	// If the map is empty, assume we need to upgrade and schedule a
	// rebuild.
	if len(internalCRLConfig.CRLExpirationMap) == 0 {
		cb.forceRebuild.Store(true)
		return nil
	}

	// Otherwise, check CRL's expirations and see if its zero or within
	// the grace period and act accordingly.
	now := time.Now()

	period, err := parseutil.ParseDurationSecond(cfg.AutoRebuildGracePeriod)
	if err != nil {
		// This may occur if the duration is empty; in that case
		// assume the default. The default should be valid and shouldn't
		// error.
		defaultPeriod, defaultErr := parseutil.ParseDurationSecond(pki_backend.DefaultCrlConfig.AutoRebuildGracePeriod)
		if defaultErr != nil {
			return fmt.Errorf("error checking for auto-rebuild status: unable to parse duration from both config's grace period (%v) and default grace period (%v):\n- config: %v\n- default: %w\n", cfg.AutoRebuildGracePeriod, pki_backend.DefaultCrlConfig.AutoRebuildGracePeriod, err, defaultErr)
		}

		period = defaultPeriod
	}

	for _, value := range internalCRLConfig.CRLExpirationMap {
		if value.IsZero() || now.After(value.Add(-1*period)) {
			cb.forceRebuild.Store(true)
			return nil
		}
	}

	return nil
}

// Mark the internal LastModifiedTime tracker invalid.
func (cb *CrlBuilder) invalidateCRLBuildTime() {
	cb.invalidate.Store(true)
}

// Update the config to mark the modified CRL. See note in
// updateDefaultIssuerId about why this is necessary.
func (cb *CrlBuilder) flushCRLBuildTimeInvalidation(sc *storageContext) error {
	if cb.invalidate.CompareAndSwap(true, false) {
		// Flush out our invalidation.
		cfg, err := sc.getLocalCRLConfig()
		if err != nil {
			cb.invalidate.Store(true)
			return fmt.Errorf("unable to update local CRL config's modification time: error fetching: %w", err)
		}

		cfg.LastModified = time.Now().UTC()
		cfg.DeltaLastModified = time.Now().UTC()
		err = sc.setLocalCRLConfig(cfg)
		if err != nil {
			cb.invalidate.Store(true)
			return fmt.Errorf("unable to update local CRL config's modification time: error persisting: %w", err)
		}
	}

	return nil
}

// RebuildIfForced is to be called by readers or periodic functions that might need to trigger
// a refresh of the CRL before the read occurs.
func (cb *CrlBuilder) RebuildIfForced(sc pki_backend.StorageContext) ([]string, error) {
	if cb.forceRebuild.Load() {
		return cb._doRebuild(sc.(*storageContext), true, _enforceForceFlag)
	}

	return nil, nil
}

// rebuild is to be called by various write apis that know the CRL is to be updated and can be now.
func (cb *CrlBuilder) Rebuild(sc pki_backend.StorageContext, forceNew bool) ([]string, error) {
	return cb._doRebuild(sc.(*storageContext), forceNew, _ignoreForceFlag)
}

// requestRebuildIfActiveNode will schedule a rebuild of the CRL from the next read or write api call assuming we are the active node of a cluster
func (cb *CrlBuilder) requestRebuildIfActiveNode(b *backend) {
	// Only schedule us on active nodes, as the active node is the only node that can rebuild/write the CRL.
	// Note 1: The CRL is cluster specific, so this does need to run on the active node of a performance secondary cluster.
	// Note 2: This is called by the storage invalidation function, so it should not block.
	if !cb.canRebuild {
		b.Logger().Debug("Ignoring request to schedule a CRL rebuild, not on active node.")
		return
	}

	b.Logger().Info("Scheduling PKI CRL rebuild.")
	// Set the flag to 1, we don't care if we aren't the ones that actually swap it to 1.
	cb.forceRebuild.Store(true)
}

func (cb *CrlBuilder) _doRebuild(sc *storageContext, forceNew bool, ignoreForceFlag bool) ([]string, error) {
	cb._builder.Lock()
	defer cb._builder.Unlock()
	// Re-read the lock in case someone beat us to the punch between the previous load op.
	forceBuildFlag := cb.forceRebuild.Load()
	if forceBuildFlag || ignoreForceFlag {
		// Reset our original flag back to 0 before we start the rebuilding. This may lead to another round of
		// CRL building, but we want to avoid the race condition caused by clearing the flag after we completed (An
		// update/revocation occurred attempting to set the flag, after we listed the certs but before we wrote
		// the CRL, so we missed the update and cleared the flag.)
		cb.forceRebuild.Store(false)

		// if forceRebuild was requested, that should force a complete rebuild even if requested not too by forceNew
		myForceNew := forceBuildFlag || forceNew
		return buildCRLs(sc, myForceNew)
	}

	return nil, nil
}

func (cb *CrlBuilder) _getPresentDeltaWALForClearing(sc pki_backend.StorageContext, path string) ([]string, error) {
	// Clearing of the delta WAL occurs after a new complete CRL has been built.
	walSerials, err := sc.GetStorage().List(sc.GetContext(), path)
	if err != nil {
		return nil, fmt.Errorf("error fetching list of delta WAL certificates to clear: %w", err)
	}

	// We _should_ remove the special WAL entries here, but we don't really
	// want to traverse the list again (and also below in clearDeltaWAL). So
	// trust the latter does the right thing.
	return walSerials, nil
}

func (cb *CrlBuilder) GetPresentLocalDeltaWALForClearing(sc pki_backend.StorageContext) ([]string, error) {
	return cb._getPresentDeltaWALForClearing(sc, localDeltaWALPath)
}

func (cb *CrlBuilder) GetPresentUnifiedDeltaWALForClearing(sc pki_backend.StorageContext) ([]string, error) {
	walClusters, err := sc.GetStorage().List(sc.GetContext(), unifiedDeltaWALPrefix)
	if err != nil {
		return nil, fmt.Errorf("error fetching list of clusters with delta WAL entries: %w", err)
	}

	var allPaths []string
	for index, cluster := range walClusters {
		prefix := unifiedDeltaWALPrefix + cluster
		clusterPaths, err := cb._getPresentDeltaWALForClearing(sc, prefix)
		if err != nil {
			return nil, fmt.Errorf("error fetching delta WAL entries for cluster (%v / %v): %w", index, cluster, err)
		}

		// Here, we don't want to include the unifiedDeltaWALPrefix because
		// ClearUnifiedDeltaWAL handles that for us. Instead, just include
		// the cluster identifier.
		for _, clusterPath := range clusterPaths {
			allPaths = append(allPaths, cluster+clusterPath)
		}
	}

	return allPaths, nil
}

func (cb *CrlBuilder) _clearDeltaWAL(sc pki_backend.StorageContext, walSerials []string, path string) error {
	// Clearing of the delta WAL occurs after a new complete CRL has been built.
	for _, serial := range walSerials {
		// Don't remove our special entries!
		if strings.HasSuffix(serial, deltaWALLastBuildSerialName) || strings.HasSuffix(serial, deltaWALLastRevokedSerialName) {
			continue
		}

		if err := sc.GetStorage().Delete(sc.GetContext(), path+serial); err != nil {
			return fmt.Errorf("error clearing delta WAL certificate: %w", err)
		}
	}

	return nil
}

func (cb *CrlBuilder) ClearLocalDeltaWAL(sc pki_backend.StorageContext, walSerials []string) error {
	return cb._clearDeltaWAL(sc, walSerials, localDeltaWALPath)
}

func (cb *CrlBuilder) ClearUnifiedDeltaWAL(sc pki_backend.StorageContext, walSerials []string) error {
	return cb._clearDeltaWAL(sc, walSerials, unifiedDeltaWALPrefix)
}

func (cb *CrlBuilder) rebuildDeltaCRLsIfForced(sc *storageContext, override bool) ([]string, error) {
	// Delta CRLs use the same expiry duration as the complete CRL. Because
	// we always rebuild the complete CRL and then the delta CRL, we can
	// be assured that the delta CRL always expires after a complete CRL,
	// and that rebuilding the complete CRL will trigger a fresh delta CRL
	// build of its own.
	//
	// This guarantee means we can avoid checking delta CRL expiry. Thus,
	// we only need to rebuild the delta CRL when we have new revocations,
	// within our time window for updating it.
	cfg, err := cb.GetConfigWithUpdate(sc)
	if err != nil {
		return nil, err
	}

	if !cfg.EnableDelta {
		// We explicitly do not update the last check time here, as we
		// want to persist the last rebuild window if it hasn't been set.
		return nil, nil
	}

	deltaRebuildDuration, err := parseutil.ParseDurationSecond(cfg.DeltaRebuildInterval)
	if err != nil {
		return nil, err
	}

	// Acquire CRL building locks before we get too much further.
	cb._builder.Lock()
	defer cb._builder.Unlock()

	// Last is setup during newCRLBuilder(...), so we don't need to deal with
	// a zero condition.
	now := time.Now()
	last := cb.lastDeltaRebuildCheck
	nextRebuildCheck := last.Add(deltaRebuildDuration)
	if !override && now.Before(nextRebuildCheck) {
		// If we're still before the time of our next rebuild check, we can
		// safely return here even if we have certs. We'll wait for a bit,
		// retrigger this check, and then do the rebuild.
		return nil, nil
	}

	// Update our check time. If we bail out below (due to storage errors
	// or whatever), we'll delay the next CRL check (hopefully allowing
	// things to stabilize). Otherwise, we might not build a new Delta CRL
	// until our next complete CRL build.
	cb.lastDeltaRebuildCheck = now

	rebuildLocal, err := cb._shouldRebuildLocalCRLs(sc, override)
	if err != nil {
		return nil, fmt.Errorf("error determining if local CRLs should be rebuilt: %w", err)
	}

	rebuildUnified, err := cb._shouldRebuildUnifiedCRLs(sc, override)
	if err != nil {
		return nil, fmt.Errorf("error determining if unified CRLs should be rebuilt: %w", err)
	}

	if !rebuildLocal && !rebuildUnified {
		return nil, nil
	}

	// Finally, we must've needed to do the rebuild. Execute!
	return cb.RebuildDeltaCRLsHoldingLock(sc, false)
}

func (cb *CrlBuilder) _shouldRebuildLocalCRLs(sc *storageContext, override bool) (bool, error) {
	// Fetch two storage entries to see if we actually need to do this
	// rebuild, given we're within the window.
	lastWALEntry, err := sc.Storage.Get(sc.Context, localDeltaWALLastRevokedSerial)
	if err != nil || !override && (lastWALEntry == nil || lastWALEntry.Value == nil) {
		// If this entry does not exist, we don't need to rebuild the
		// delta WAL due to the expiration assumption above. There must
		// not have been any new revocations. Since err should be nil
		// in this case, we can safely return it.
		return false, err
	}

	lastBuildEntry, err := sc.Storage.Get(sc.Context, localDeltaWALLastBuildSerial)
	if err != nil {
		return false, err
	}

	if !override && lastBuildEntry != nil && lastBuildEntry.Value != nil {
		// If the last build entry doesn't exist, we still want to build a
		// new delta WAL, since this could be our very first time doing so.
		//
		// Otherwise, here, now that we know it exists, we want to check this
		// value against the other value. Since we previously guarded the WAL
		// entry being non-empty, we're good to decode everything within this
		// guard.
		var walInfo lastWALInfo
		if err := lastWALEntry.DecodeJSON(&walInfo); err != nil {
			return false, err
		}

		var deltaInfo lastDeltaInfo
		if err := lastBuildEntry.DecodeJSON(&deltaInfo); err != nil {
			return false, err
		}

		// Here, everything decoded properly and we know that no new certs
		// have been revoked since we built this last delta CRL. We can exit
		// without rebuilding then.
		if walInfo.Serial == deltaInfo.Serial {
			return false, nil
		}
	}

	return true, nil
}

func (cb *CrlBuilder) _shouldRebuildUnifiedCRLs(sc *storageContext, override bool) (bool, error) {
	// Unified CRL can only be built by the main cluster.
	sysView := sc.System()
	if sysView.ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!sysView.LocalMount() && sysView.ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		return false, nil
	}

	// If we're overriding whether we should build Delta CRLs, always return
	// true, even if storage errors might've happen.
	if override {
		return true, nil
	}

	// Fetch two storage entries to see if we actually need to do this
	// rebuild, given we're within the window. We need to fetch these
	// two entries per cluster.
	clusters, err := sc.Storage.List(sc.Context, unifiedDeltaWALPrefix)
	if err != nil {
		return false, fmt.Errorf("failed to get the list of clusters having written Delta WALs: %w", err)
	}

	// If any cluster tells us to rebuild, we should rebuild.
	shouldRebuild := false
	for index, cluster := range clusters {
		prefix := unifiedDeltaWALPrefix + cluster
		clusterUnifiedLastRevokedWALEntry := prefix + deltaWALLastRevokedSerialName
		clusterUnifiedLastBuiltWALEntry := prefix + deltaWALLastBuildSerialName

		lastWALEntry, err := sc.Storage.Get(sc.Context, clusterUnifiedLastRevokedWALEntry)
		if err != nil {
			return false, fmt.Errorf("failed fetching last revoked WAL entry for cluster (%v / %v): %w", index, cluster, err)
		}

		if lastWALEntry == nil || lastWALEntry.Value == nil {
			continue
		}

		lastBuildEntry, err := sc.Storage.Get(sc.Context, clusterUnifiedLastBuiltWALEntry)
		if err != nil {
			return false, fmt.Errorf("failed fetching last built CRL WAL entry for cluster (%v / %v): %w", index, cluster, err)
		}

		if lastBuildEntry == nil || lastBuildEntry.Value == nil {
			// If the last build entry doesn't exist, we still want to build a
			// new delta WAL, since this could be our very first time doing so.
			shouldRebuild = true
			break
		}

		// Otherwise, here, now that we know it exists, we want to check this
		// value against the other value. Since we previously guarded the WAL
		// entry being non-empty, we're good to decode everything within this
		// guard.
		var walInfo lastWALInfo
		if err := lastWALEntry.DecodeJSON(&walInfo); err != nil {
			return false, fmt.Errorf("failed decoding last revoked WAL entry for cluster (%v / %v): %w", index, cluster, err)
		}

		var deltaInfo lastDeltaInfo
		if err := lastBuildEntry.DecodeJSON(&deltaInfo); err != nil {
			return false, fmt.Errorf("failed decoding last built CRL WAL entry for cluster (%v / %v): %w", index, cluster, err)
		}

		if walInfo.Serial != deltaInfo.Serial {
			shouldRebuild = true
			break
		}
	}

	// No errors occurred, so return the result.
	return shouldRebuild, nil
}

func (cb *CrlBuilder) rebuildDeltaCRLs(sc *storageContext, forceNew bool) ([]string, error) {
	cb._builder.Lock()
	defer cb._builder.Unlock()

	return cb.RebuildDeltaCRLsHoldingLock(sc, forceNew)
}

func (cb *CrlBuilder) RebuildDeltaCRLsHoldingLock(sc pki_backend.StorageContext, forceNew bool) ([]string, error) {
	return buildAnyCRLs(sc.(*storageContext), forceNew, true /* building delta */)
}

func (cb *CrlBuilder) addCertForRevocationCheck(cluster, serial string) {
	entry := &revocationQueueEntry{
		Cluster: cluster,
		Serial:  serial,
	}
	cb.revQueue.Add(entry)
}

func (cb *CrlBuilder) addCertForRevocationRemoval(cluster, serial string) {
	entry := &revocationQueueEntry{
		Cluster: cluster,
		Serial:  serial,
	}
	cb.removalQueue.Add(entry)
}

func (cb *CrlBuilder) addCertFromCrossRevocation(cluster, serial string) {
	entry := &revocationQueueEntry{
		Cluster: cluster,
		Serial:  serial,
	}
	cb.crossQueue.Add(entry)
}

func (cb *CrlBuilder) maybeGatherQueueForFirstProcess(sc *storageContext, isNotPerfPrimary bool) error {
	// Assume holding lock.
	if cb.haveInitializedQueue.Load() {
		return nil
	}

	sc.Logger().Debug(fmt.Sprintf("gathering first time existing revocations"))

	clusters, err := sc.Storage.List(sc.Context, crossRevocationPrefix)
	if err != nil {
		return fmt.Errorf("failed to list cross-cluster revocation queue participating clusters: %w", err)
	}

	sc.Logger().Debug(fmt.Sprintf("found %v clusters: %v", len(clusters), clusters))

	for cIndex, cluster := range clusters {
		cluster = cluster[0 : len(cluster)-1]
		cPath := crossRevocationPrefix + cluster + "/"
		serials, err := sc.Storage.List(sc.Context, cPath)
		if err != nil {
			return fmt.Errorf("failed to list cross-cluster revocation queue entries for cluster %v (%v): %w", cluster, cIndex, err)
		}

		sc.Logger().Debug(fmt.Sprintf("found %v serials for cluster %v: %v", len(serials), cluster, serials))

		for _, serial := range serials {
			if serial[len(serial)-1] == '/' {
				serial = serial[0 : len(serial)-1]
			}

			ePath := cPath + serial
			eConfirmPath := ePath + "/confirmed"
			removalEntry, err := sc.Storage.Get(sc.Context, eConfirmPath)

			entry := &revocationQueueEntry{
				Cluster: cluster,
				Serial:  serial,
			}

			// No removal entry yet; add to regular queue. Otherwise, slate it
			// for removal if we're a perfPrimary.
			if err != nil || removalEntry == nil {
				cb.revQueue.Add(entry)
			} else if !isNotPerfPrimary {
				cb.removalQueue.Add(entry)
			} // Else, this is a confirmation but we're on a perf secondary so ignore it.

			// Overwrite the error; we don't really care about its contents
			// at this step.
			err = nil
		}
	}

	return nil
}

func (cb *CrlBuilder) processRevocationQueue(sc *storageContext) error {
	sc.Logger().Debug(fmt.Sprintf("starting to process revocation requests"))

	isNotPerfPrimary := sc.System().ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!sc.System().LocalMount() && sc.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary))

	if err := cb.maybeGatherQueueForFirstProcess(sc, isNotPerfPrimary); err != nil {
		return fmt.Errorf("failed to gather first queue: %w", err)
	}

	revQueue := cb.revQueue.Iterate()
	removalQueue := cb.removalQueue.Iterate()

	sc.Logger().Debug(fmt.Sprintf("gathered %v revocations and %v confirmation entries", len(revQueue), len(removalQueue)))

	crlConfig, err := cb.GetConfigWithUpdate(sc)
	if err != nil {
		return err
	}

	ourClusterId, err := sc.System().ClusterID(sc.Context)
	if err != nil {
		return fmt.Errorf("unable to fetch clusterID to ignore local revocation entries: %w", err)
	}

	for _, req := range revQueue {
		// Regardless of whether we're on the perf primary or a secondary
		// cluster, we can safely ignore revocation requests originating
		// from our node, because we've already checked them once (when
		// they were created).
		if ourClusterId != "" && ourClusterId == req.Cluster {
			continue
		}

		// Fetch the revocation entry to ensure it exists.
		rPath := crossRevocationPrefix + req.Cluster + "/" + req.Serial
		entry, err := sc.Storage.Get(sc.Context, rPath)
		if err != nil {
			return fmt.Errorf("failed to read cross-cluster revocation queue entry: %w", err)
		}
		if entry == nil {
			// Skipping this entry; it was likely an incorrect invalidation
			// caused by the primary cluster removing the confirmation.
			cb.revQueue.Remove(req)
			continue
		}

		resp, err := tryRevokeCertBySerial(sc, crlConfig, req.Serial)
		if err == nil && resp != nil && !resp.IsError() && resp.Data != nil && resp.Data["state"].(string) == "revoked" {
			if isNotPerfPrimary {
				// Write a revocation queue removal entry.
				confirmed := revocationConfirmed{
					RevokedAt: resp.Data["revocation_time_rfc3339"].(string),
					Source:    req.Cluster,
				}
				path := crossRevocationPath + req.Serial + "/confirmed"
				confirmedEntry, err := logical.StorageEntryJSON(path, confirmed)
				if err != nil {
					return fmt.Errorf("failed to create storage entry for cross-cluster revocation confirmed response: %w", err)
				}

				if err := sc.Storage.Put(sc.Context, confirmedEntry); err != nil {
					return fmt.Errorf("error persisting cross-cluster revocation confirmation: %w", err)
				}
			} else {
				// Since we're the active node of the primary cluster, go ahead
				// and just remove it.
				path := crossRevocationPrefix + req.Cluster + "/" + req.Serial
				if err := sc.Storage.Delete(sc.Context, path); err != nil {
					return fmt.Errorf("failed to delete processed revocation request: %w", err)
				}
			}
		} else if err != nil {
			// Because we fake being from a lease, we get the guarantee that
			// err == nil == resp if the cert was already revoked; this means
			// this err should actually be fatal.
			return err
		}
		cb.revQueue.Remove(req)
	}

	if isNotPerfPrimary {
		sc.Logger().Debug(fmt.Sprintf("not on perf primary so ignoring any revocation confirmations"))

		// See note in pki/backend.go; this should be empty.
		cb.removalQueue.RemoveAll()
		cb.haveInitializedQueue.Store(true)
		return nil
	}

	clusters, err := sc.Storage.List(sc.Context, crossRevocationPrefix)
	if err != nil {
		return err
	}

	for _, entry := range removalQueue {
		// First remove the revocation request.
		for cIndex, cluster := range clusters {
			eEntry := crossRevocationPrefix + cluster + entry.Serial
			if err := sc.Storage.Delete(sc.Context, eEntry); err != nil {
				return fmt.Errorf("failed to delete potential cross-cluster revocation entry for cluster %v (%v) and serial %v: %w", cluster, cIndex, entry.Serial, err)
			}
		}

		// Then remove the confirmation.
		if err := sc.Storage.Delete(sc.Context, crossRevocationPrefix+entry.Cluster+"/"+entry.Serial+"/confirmed"); err != nil {
			return fmt.Errorf("failed to delete cross-cluster revocation confirmation entry for cluster %v and serial %v: %w", entry.Cluster, entry.Serial, err)
		}

		cb.removalQueue.Remove(entry)
	}

	cb.haveInitializedQueue.Store(true)

	return nil
}

func (cb *CrlBuilder) processCrossClusterRevocations(sc *storageContext) error {
	sc.Logger().Debug(fmt.Sprintf("starting to process unified revocations"))

	crlConfig, err := cb.GetConfigWithUpdate(sc)
	if err != nil {
		return err
	}

	if !crlConfig.UnifiedCRL {
		cb.crossQueue.RemoveAll()
		return nil
	}

	crossQueue := cb.crossQueue.Iterate()
	sc.Logger().Debug(fmt.Sprintf("gathered %v unified revocations entries", len(crossQueue)))

	ourClusterId, err := sc.System().ClusterID(sc.Context)
	if err != nil {
		return fmt.Errorf("unable to fetch clusterID to ignore local unified revocation entries: %w", err)
	}

	for _, req := range crossQueue {
		// Regardless of whether we're on the perf primary or a secondary
		// cluster, we can safely ignore revocation requests originating
		// from our node, because we've already checked them once (when
		// they were created).
		if ourClusterId != "" && ourClusterId == req.Cluster {
			continue
		}

		// Fetch the revocation entry to ensure it exists and this wasn't
		// a delete.
		rPath := unifiedRevocationReadPathPrefix + req.Cluster + "/" + req.Serial
		entry, err := sc.Storage.Get(sc.Context, rPath)
		if err != nil {
			return fmt.Errorf("failed to read unified revocation entry: %w", err)
		}
		if entry == nil {
			// Skip this entry: it was likely caused by the deletion of this
			// record during tidy.
			cb.crossQueue.Remove(req)
			continue
		}

		resp, err := tryRevokeCertBySerial(sc, crlConfig, req.Serial)
		if err == nil && resp != nil && !resp.IsError() && resp.Data != nil && resp.Data["state"].(string) == "revoked" {
			// We could theoretically save ourselves from writing a global
			// revocation entry during the above certificate revocation, as
			// we don't really need it to appear on either the unified CRL
			// or its delta CRL, but this would require more plumbing.
			cb.crossQueue.Remove(req)
		} else if err != nil {
			// Because we fake being from a lease, we get the guarantee that
			// err == nil == resp if the cert was already revoked; this means
			// this err should actually be fatal.
			return err
		}
	}

	return nil
}

// Revoke a certificate from a given serial number if it is present in local
// storage.
func tryRevokeCertBySerial(sc *storageContext, config *pki_backend.CrlConfig, serial string) (*logical.Response, error) {
	// revokeCert requires us to hold these locks before calling it.
	sc.GetRevokeStorageLock().Lock()
	defer sc.GetRevokeStorageLock().Unlock()

	certEntry, err := fetchCertBySerial(sc, issuing.PathCerts, serial)
	if err != nil {
		switch err.(type) {
		case errutil.UserError:
			return logical.ErrorResponse(err.Error()), nil
		default:
			return nil, err
		}
	}

	if certEntry == nil {
		return nil, nil
	}

	cert, err := x509.ParseCertificate(certEntry.Value)
	if err != nil {
		return nil, fmt.Errorf("error parsing certificate: %w", err)
	}

	return revokeCert(sc, config, cert)
}

// Revokes a cert, and tries to be smart about error recovery
func revokeCert(sc *storageContext, config *pki_backend.CrlConfig, cert *x509.Certificate) (*logical.Response, error) {
	// As this backend is self-contained and this function does not hook into
	// third parties to manage users or resources, if the mount is tainted,
	// revocation doesn't matter anyways -- the CRL that would be written will
	// be immediately blown away by the view being cleared. So we can simply
	// fast path a successful exit.
	if sc.System().Tainted() {
		return nil, nil
	}

	colonSerial := parsing.SerialFromCert(cert)
	hyphenSerial := parsing.NormalizeSerialForStorage(colonSerial)

	// Validate that no issuers match the serial number to be revoked. We need
	// to gracefully degrade to the legacy cert bundle when it is required, as
	// secondary PR clusters might not have been upgraded, but still need to
	// handle revoking certs.
	issuerIDCertMap, err := revocation.FetchIssuerMapForRevocationChecking(sc)
	if err != nil {
		return nil, err
	}

	// Ensure we don't revoke an issuer via this API; use /issuer/:issuer_ref/revoke
	// instead.
	for issuer, certificate := range issuerIDCertMap {
		if colonSerial == serialFromCert(certificate) {
			return logical.ErrorResponse(fmt.Sprintf("adding issuer (id: %v) to its own CRL is not allowed", issuer)), nil
		}
	}

	curRevInfo, err := fetchRevocationInfo(sc, colonSerial)
	if err != nil {
		return nil, err
	}
	if curRevInfo != nil {
		resp := &logical.Response{
			Data: map[string]interface{}{
				"revocation_time": curRevInfo.RevocationTime,
				"state":           "revoked",
			},
		}
		if !curRevInfo.RevocationTimeUTC.IsZero() {
			resp.Data["revocation_time_rfc3339"] = curRevInfo.RevocationTimeUTC.Format(time.RFC3339Nano)
		}

		return resp, nil
	}

	// Add a little wiggle room because leases are stored with a second
	// granularity
	if cert.NotAfter.Before(time.Now().Add(2 * time.Second)) {
		response := &logical.Response{}
		response.AddWarning(fmt.Sprintf("certificate with serial %s already expired; refusing to add to CRL", colonSerial))
		return response, nil
	}

	currTime := time.Now()
	revInfo := revocation.RevocationInfo{
		CertificateBytes:  cert.Raw,
		RevocationTime:    currTime.Unix(),
		RevocationTimeUTC: currTime.UTC(),
	}

	// We may not find an issuer with this certificate; that's fine so
	// ignore the return value.
	revInfo.AssociateRevokedCertWithIsssuer(cert, issuerIDCertMap)

	revEntry, err := logical.StorageEntryJSON(revokedPath+hyphenSerial, revInfo)
	if err != nil {
		return nil, fmt.Errorf("error creating revocation entry: %w", err)
	}

	certCounter := sc.GetCertificateCounter()
	certsCounted := certCounter.IsInitialized()
	err = sc.Storage.Put(sc.Context, revEntry)
	if err != nil {
		return nil, fmt.Errorf("error saving revoked certificate to new location: %w", err)
	}
	certCounter.IncrementTotalRevokedCertificatesCount(certsCounted, revEntry.Key)

	// From here on out, the certificate has been revoked locally. Any other
	// persistence issues might still err, but any other failure messages
	// should be added as warnings to the revocation.
	resp := &logical.Response{
		Data: map[string]interface{}{
			"revocation_time":         revInfo.RevocationTime,
			"revocation_time_rfc3339": revInfo.RevocationTimeUTC.Format(time.RFC3339Nano),
			"state":                   "revoked",
		},
	}

	// If this flag is enabled after the fact, existing local entries will be published to
	// the unified storage space through a periodic function.
	failedWritingUnifiedCRL := false
	if config.UnifiedCRL {
		entry := &revocation.UnifiedRevocationEntry{
			SerialNumber:      colonSerial,
			CertExpiration:    cert.NotAfter,
			RevocationTimeUTC: revInfo.RevocationTimeUTC,
			CertificateIssuer: revInfo.CertificateIssuer,
		}

		ignoreErr := revocation.WriteUnifiedRevocationEntry(sc.GetContext(), sc.GetStorage(), entry)
		if ignoreErr != nil {
			// Just log the error if we fail to write across clusters, a separate background
			// thread will reattempt it later on as we have the local write done.
			sc.Logger().Error("Failed to write unified revocation entry, will re-attempt later",
				"serial_number", colonSerial, "error", ignoreErr)
			sc.GetUnifiedTransferStatus().forceRun()

			resp.AddWarning(fmt.Sprintf("Failed to write unified revocation entry, will re-attempt later: %v", err))
			failedWritingUnifiedCRL = true
		}
	}

	if !config.AutoRebuild {
		// Note that writing the Delta WAL here isn't necessary; we've
		// already rebuilt the full CRL so the Delta WAL will be cleared
		// afterwards. Writing an entry only to immediately remove it
		// isn't necessary.
		warnings, crlErr := sc.CrlBuilder().Rebuild(sc, false)
		if crlErr != nil {
			switch crlErr.(type) {
			case errutil.UserError:
				return logical.ErrorResponse(fmt.Sprintf("Error during CRL building: %s", crlErr)), nil
			default:
				return nil, fmt.Errorf("error encountered during CRL building: %w", crlErr)
			}
		}
		for index, warning := range warnings {
			resp.AddWarning(fmt.Sprintf("Warning %d during CRL rebuild: %v", index+1, warning))
		}
	} else if config.EnableDelta {
		if err := writeRevocationDeltaWALs(sc, config, resp, failedWritingUnifiedCRL, hyphenSerial, colonSerial); err != nil {
			return nil, fmt.Errorf("failed to write WAL entries for Delta CRLs: %w", err)
		}
	}

	return resp, nil
}

func writeRevocationDeltaWALs(sc *storageContext, config *pki_backend.CrlConfig, resp *logical.Response, failedWritingUnifiedCRL bool, hyphenSerial string, colonSerial string) error {
	if err := writeSpecificRevocationDeltaWALs(sc, hyphenSerial, colonSerial, localDeltaWALPath); err != nil {
		return fmt.Errorf("failed to write local delta WAL entry: %w", err)
	}

	if config.UnifiedCRL && !failedWritingUnifiedCRL {
		// We only need to write cross-cluster unified Delta WAL entries when
		// it is enabled; in particular, because we rebuild CRLs when enabling
		// this flag, any revocations that happened prior to enabling unified
		// revocation will appear on the complete CRL (+/- synchronization:
		// in particular, if a perf replica revokes a cert prior to seeing
		// unified revocation enabled, but after the main node has done the
		// listing for the unified CRL rebuild, this revocation will not
		// appear on either the main or the next delta CRL, but will need to
		// wait for a subsequent complete CRL rebuild).
		//
		// Lastly, we don't attempt this if the unified CRL entry failed to
		// write, as we need that entry before the delta WAL entry will make
		// sense.
		if ignoredErr := writeSpecificRevocationDeltaWALs(sc, hyphenSerial, colonSerial, unifiedDeltaWALPath); ignoredErr != nil {
			// Just log the error if we fail to write across clusters, a separate background
			// thread will reattempt it later on as we have the local write done.
			sc.Logger().Error("Failed to write cross-cluster delta WAL entry, will re-attempt later",
				"serial_number", colonSerial, "error", ignoredErr)
			sc.GetUnifiedTransferStatus().forceRun()

			resp.AddWarning(fmt.Sprintf("Failed to write cross-cluster delta WAL entry, will re-attempt later: %v", ignoredErr))
		}
	} else if failedWritingUnifiedCRL {
		resp.AddWarning("Skipping cross-cluster delta WAL entry as cross-cluster revocation failed to write; will re-attempt later.")
	}

	return nil
}

func writeSpecificRevocationDeltaWALs(sc *storageContext, hyphenSerial string, colonSerial string, pathPrefix string) error {
	// Previously, regardless of whether or not we've presently enabled
	// Delta CRLs, we would always write the Delta WAL in case it is
	// enabled in the future. We though we could trigger another full CRL
	// rebuild instead (to avoid inconsistent state between the CRL and
	// missing Delta WAL entries), but writing extra (unused?) WAL entries
	// versus an expensive full CRL rebuild was thought of as being
	// probably a net wash.
	//
	// However, we've now added unified CRL building, adding cross-cluster
	// writes to the revocation path. Because this is relatively expensive,
	// we've opted to rebuild the complete+delta CRLs when toggling the
	// state of delta enabled, instead of always writing delta CRL entries.
	//
	// Thus Delta WAL building happens **only** when Delta CRLs are enabled.
	//
	// We should only do this when the cert hasn't already been revoked.
	// Otherwise, the re-revocation may appear on both an existing CRL and
	// on a delta CRL, or a serial may be skipped from the delta CRL if
	// there's an A->B->A revocation pattern and the delta was rebuilt
	// after the first cert.
	//
	// Currently we don't store any data in the WAL entry.
	var walInfo deltaWALInfo
	walEntry, err := logical.StorageEntryJSON(pathPrefix+hyphenSerial, walInfo)
	if err != nil {
		return fmt.Errorf("unable to create delta CRL WAL entry: %w", err)
	}

	if err = sc.Storage.Put(sc.Context, walEntry); err != nil {
		return fmt.Errorf("error saving delta CRL WAL entry: %w", err)
	}

	// In order for periodic delta rebuild to be mildly efficient, we
	// should write the last revoked delta WAL entry so we know if we
	// have new revocations that we should rebuild the delta WAL for.
	lastRevSerial := lastWALInfo{Serial: colonSerial}
	lastWALEntry, err := logical.StorageEntryJSON(pathPrefix+deltaWALLastRevokedSerialName, lastRevSerial)
	if err != nil {
		return fmt.Errorf("unable to create last delta CRL WAL entry: %w", err)
	}
	if err = sc.Storage.Put(sc.Context, lastWALEntry); err != nil {
		return fmt.Errorf("error saving last delta CRL WAL entry: %w", err)
	}

	return nil
}

func buildCRLs(sc *storageContext, forceNew bool) ([]string, error) {
	return buildAnyCRLs(sc, forceNew, false)
}

func buildAnyCRLs(sc *storageContext, forceNew bool, isDelta bool) ([]string, error) {
	// In order to build all CRLs, we need knowledge of all issuers. Any two
	// issuers with the same keys _and_ subject should have the same CRL since
	// they're functionally equivalent.
	//
	// When building CRLs, there's two types of CRLs: an "internal" CRL for
	// just certificates issued by this issuer, and a "default" CRL, which
	// not only contains certificates by this issuer, but also ones issued
	// by "unknown" or past issuers. This means we need knowledge of not
	// only all issuers (to tell whether or not to include these orphaned
	// certs) but whether the present issuer is the configured default.
	//
	// If a configured default is lacking, we won't provision these
	// certificates on any CRL.
	//
	// In order to know which CRL a given cert belongs on, we have to read
	// it into memory, identify the corresponding issuer, and update its
	// map with the revoked cert instance. If no such issuer is found, we'll
	// place it in the default issuer's CRL.
	//
	// By not relying on the _cert_'s storage, we allow issuers to come and
	// go (either by direct deletion, having their keys deleted, or by usage
	// restrictions) -- and when they return, we'll correctly place certs
	// on their CRLs.

	// See the message in revokedCert about rebuilding CRLs: we need to
	// gracefully handle revoking entries with the legacy cert bundle.
	var err error
	var issuers []issuing.IssuerID
	var wasLegacy bool

	// First, fetch an updated copy of the CRL config. We'll pass this into buildCRL.
	globalCRLConfig, err := sc.CrlBuilder().GetConfigWithUpdate(sc)
	if err != nil {
		return nil, fmt.Errorf("error building CRL: while updating config: %w", err)
	}

	if globalCRLConfig.Disable && !forceNew {
		// We build a single long-lived (but regular validity) empty CRL in
		// the event that we disable the CRL, but we don't keep updating it
		// with newer, more-valid empty CRLs in the event that we later
		// re-enable it. This is a historical behavior.
		//
		// So, since tidy can now associate issuers on revocation entries, we
		// can skip the rest of this function and exit early without updating
		// anything.
		return nil, nil
	}

	if !sc.UseLegacyBundleCaStorage() {
		issuers, err = sc.listIssuers()
		if err != nil {
			return nil, fmt.Errorf("error building CRL: while listing issuers: %w", err)
		}
	} else {
		// Here, we hard-code the legacy issuer entry instead of using the
		// default ref. This is because we need to hack some of the logic
		// below for revocation to handle the legacy bundle.
		issuers = []issuing.IssuerID{legacyBundleShimID}
		wasLegacy = true

		// Here, we avoid building a delta CRL with the legacy CRL bundle.
		//
		// Users should upgrade symmetrically, rather than attempting
		// backward compatibility for new features across disparate versions.
		if isDelta {
			return []string{"refusing to rebuild delta CRL with legacy bundle; finish migrating to newer issuer storage layout"}, nil
		}
	}

	issuersConfig, err := sc.getIssuersConfig()
	if err != nil {
		return nil, fmt.Errorf("error building CRLs: while getting the default config: %w", err)
	}

	// We map IssuerID->entry for fast lookup and also IssuerID->Cert for
	// signature verification and correlation of revoked certs.
	issuerIDEntryMap := make(map[issuing.IssuerID]*issuing.IssuerEntry, len(issuers))
	issuerIDCertMap := make(map[issuing.IssuerID]*x509.Certificate, len(issuers))

	// We use a double map (KeyID->subject->IssuerID) to store whether or not this
	// key+subject paring has been seen before. We can then iterate over each
	// key/subject and choose any representative issuer for that combination.
	keySubjectIssuersMap := make(map[issuing.KeyID]map[string][]issuing.IssuerID)
	for _, issuer := range issuers {
		// We don't strictly need this call, but by requesting the bundle, the
		// legacy path is automatically ignored.
		thisEntry, _, err := sc.fetchCertBundleByIssuerId(issuer, false)
		if err != nil {
			return nil, fmt.Errorf("error building CRLs: unable to fetch specified issuer (%v): %w", issuer, err)
		}

		if len(thisEntry.KeyID) == 0 {
			continue
		}

		// n.b.: issuer usage check has been delayed. This occurred because
		// we want to ensure any issuer (representative of a larger set) can
		// be used to associate revocation entries and we won't bother
		// rewriting that entry (causing churn) if the particular selected
		// issuer lacks CRL signing capabilities.
		//
		// The result is that this map (and the other maps) contain all the
		// issuers we know about, and only later do we check crlSigning before
		// choosing our representative.
		//
		// The other side effect (making this not compatible with Vault 1.11
		// behavior) is that _identified_ certificates whose issuer set is
		// not allowed for crlSigning will no longer appear on the default
		// issuer's CRL.
		issuerIDEntryMap[issuer] = thisEntry

		thisCert, err := thisEntry.GetCertificate()
		if err != nil {
			return nil, fmt.Errorf("error building CRLs: unable to parse issuer (%v)'s certificate: %w", issuer, err)
		}
		issuerIDCertMap[issuer] = thisCert

		subject := string(thisCert.RawSubject)
		if _, ok := keySubjectIssuersMap[thisEntry.KeyID]; !ok {
			keySubjectIssuersMap[thisEntry.KeyID] = make(map[string][]issuing.IssuerID)
		}

		keySubjectIssuersMap[thisEntry.KeyID][subject] = append(keySubjectIssuersMap[thisEntry.KeyID][subject], issuer)
	}

	// Now we do two calls: building the cluster-local CRL, and potentially
	// building the global CRL if we're on the active node of the performance
	// primary.
	currLocalDeltaSerials, localWarnings, err := buildAnyLocalCRLs(sc, issuersConfig, globalCRLConfig,
		issuers, issuerIDEntryMap,
		issuerIDCertMap, keySubjectIssuersMap,
		wasLegacy, forceNew, isDelta)
	if err != nil {
		return nil, err
	}
	currUnifiedDeltaSerials, unifiedWarnings, err := buildAnyUnifiedCRLs(sc, issuersConfig, globalCRLConfig,
		issuers, issuerIDEntryMap,
		issuerIDCertMap, keySubjectIssuersMap,
		wasLegacy, forceNew, isDelta)
	if err != nil {
		return nil, err
	}

	var warnings []string
	for _, warning := range localWarnings {
		warnings = append(warnings, fmt.Sprintf("warning from local CRL rebuild: %v", warning))
	}
	for _, warning := range unifiedWarnings {
		warnings = append(warnings, fmt.Sprintf("warning from unified CRL rebuild: %v", warning))
	}

	// Finally, we decide if we need to rebuild the Delta CRLs again, for both
	// global and local CRLs if necessary.
	if !isDelta {
		// After we've confirmed the primary CRLs have built OK, go ahead and
		// clear the delta CRL WAL and rebuild it.
		if err := sc.CrlBuilder().ClearLocalDeltaWAL(sc, currLocalDeltaSerials); err != nil {
			return nil, fmt.Errorf("error building CRLs: unable to clear Delta WAL: %w", err)
		}
		if err := sc.CrlBuilder().ClearUnifiedDeltaWAL(sc, currUnifiedDeltaSerials); err != nil {
			return nil, fmt.Errorf("error building CRLs: unable to clear Delta WAL: %w", err)
		}
		deltaWarnings, err := sc.CrlBuilder().RebuildDeltaCRLsHoldingLock(sc, forceNew)
		if err != nil {
			return nil, fmt.Errorf("error building CRLs: unable to rebuild empty Delta WAL: %w", err)
		}
		for _, warning := range deltaWarnings {
			warnings = append(warnings, fmt.Sprintf("warning from delta CRL rebuild: %v", warning))
		}
	}

	return warnings, nil
}

func getLastWALSerial(sc *storageContext, path string) (string, error) {
	lastWALEntry, err := sc.Storage.Get(sc.Context, path)
	if err != nil {
		return "", err
	}

	if lastWALEntry != nil && lastWALEntry.Value != nil {
		var walInfo lastWALInfo
		if err := lastWALEntry.DecodeJSON(&walInfo); err != nil {
			return "", err
		}

		return walInfo.Serial, nil
	}

	// No serial to return.
	return "", nil
}

func buildAnyLocalCRLs(
	sc *storageContext,
	issuersConfig *issuing.IssuerConfigEntry,
	globalCRLConfig *pki_backend.CrlConfig,
	issuers []issuing.IssuerID,
	issuerIDEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	issuerIDCertMap map[issuing.IssuerID]*x509.Certificate,
	keySubjectIssuersMap map[issuing.KeyID]map[string][]issuing.IssuerID,
	wasLegacy bool,
	forceNew bool,
	isDelta bool,
) ([]string, []string, error) {
	var err error
	var warnings []string

	// Before we load cert entries, we want to store the last seen delta WAL
	// serial number. The subsequent List will have at LEAST that certificate
	// (and potentially more) in it; when we're done writing the delta CRL,
	// we'll write this serial as a sentinel to see if we need to rebuild it
	// in the future.
	var lastDeltaSerial string
	if isDelta {
		lastDeltaSerial, err = getLastWALSerial(sc, localDeltaWALLastRevokedSerial)
		if err != nil {
			return nil, nil, err
		}
	}

	// We fetch a list of delta WAL entries prior to generating the complete
	// CRL. This allows us to avoid a lock (to clear such storage): anything
	// visible now, should also be visible on the complete CRL we're writing.
	var currDeltaCerts []string
	if !isDelta {
		currDeltaCerts, err = sc.CrlBuilder().GetPresentLocalDeltaWALForClearing(sc)
		if err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to get present delta WAL entries for removal: %w", err)
		}
	}

	var unassignedCerts []pkix.RevokedCertificate
	var revokedCertsMap map[issuing.IssuerID][]pkix.RevokedCertificate

	// If the CRL is disabled do not bother reading in all the revoked certificates.
	if !globalCRLConfig.Disable {
		// Next, we load and parse all revoked certificates. We need to assign
		// these certificates to an issuer. Some certificates will not be
		// assignable (if they were issued by a since-deleted issuer), so we need
		// a separate pool for those.
		unassignedCerts, revokedCertsMap, err = getLocalRevokedCertEntries(sc, issuerIDCertMap, isDelta)
		if err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to get revoked certificate entries: %w", err)
		}

		if !isDelta {
			// Revoking an issuer forces us to rebuild our complete CRL,
			// regardless of whether or not we've enabled auto rebuilding or
			// delta CRLs. If we elide the above isDelta check, this results
			// in a non-empty delta CRL, containing the serial of the
			// now-revoked issuer, even though it was generated _after_ the
			// complete CRL with the issuer on it. There's no reason to
			// duplicate this serial number on the delta, hence the above
			// guard for isDelta.
			if err := augmentWithRevokedIssuers(issuerIDEntryMap, issuerIDCertMap, revokedCertsMap); err != nil {
				return nil, nil, fmt.Errorf("error building CRLs: unable to parse revoked issuers: %w", err)
			}
		}
	}

	// Fetch the cluster-local CRL mapping so we know where to write the
	// CRLs.
	internalCRLConfig, err := sc.getLocalCRLConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error building CRLs: unable to fetch cluster-local CRL configuration: %w", err)
	}

	rebuildWarnings, err := buildAnyCRLsWithCerts(sc, issuersConfig, globalCRLConfig, internalCRLConfig,
		issuers, issuerIDEntryMap, keySubjectIssuersMap,
		unassignedCerts, revokedCertsMap,
		forceNew, false /* isUnified */, isDelta)
	if err != nil {
		return nil, nil, fmt.Errorf("error building CRLs: %w", err)
	}
	if len(rebuildWarnings) > 0 {
		warnings = append(warnings, rebuildWarnings...)
	}

	// Finally, persist our potentially updated local CRL config. Only do this
	// if we didn't have a legacy CRL bundle.
	if !wasLegacy {
		if err := sc.setLocalCRLConfig(internalCRLConfig); err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to persist updated cluster-local CRL config: %w", err)
		}
	}

	if isDelta {
		// Update our last build time here so we avoid checking for new certs
		// for a while.
		sc.CrlBuilder().SetLastDeltaRebuildCheckTime(time.Now())

		if len(lastDeltaSerial) > 0 {
			// When we have a last delta serial, write out the relevant info
			// so we can skip extra CRL rebuilds.
			deltaInfo := lastDeltaInfo{Serial: lastDeltaSerial}

			lastDeltaBuildEntry, err := logical.StorageEntryJSON(localDeltaWALLastBuildSerial, deltaInfo)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating last delta CRL rebuild serial entry: %w", err)
			}

			err = sc.Storage.Put(sc.Context, lastDeltaBuildEntry)
			if err != nil {
				return nil, nil, fmt.Errorf("error persisting last delta CRL rebuild info: %w", err)
			}
		}
	}

	return currDeltaCerts, warnings, nil
}

func buildAnyUnifiedCRLs(
	sc *storageContext,
	issuersConfig *issuing.IssuerConfigEntry,
	globalCRLConfig *pki_backend.CrlConfig,
	issuers []issuing.IssuerID,
	issuerIDEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	issuerIDCertMap map[issuing.IssuerID]*x509.Certificate,
	keySubjectIssuersMap map[issuing.KeyID]map[string][]issuing.IssuerID,
	wasLegacy bool,
	forceNew bool,
	isDelta bool,
) ([]string, []string, error) {
	var err error
	var warnings []string

	// Unified CRL can only be built by the main cluster.
	sysView := sc.System()
	if sysView.ReplicationState().HasState(consts.ReplicationDRSecondary|consts.ReplicationPerformanceStandby) ||
		(!sysView.LocalMount() && sysView.ReplicationState().HasState(consts.ReplicationPerformanceSecondary)) {
		return nil, nil, nil
	}

	// Unified CRL should only be built if enabled.
	if !globalCRLConfig.UnifiedCRL && !forceNew {
		return nil, nil, nil
	}

	// Before we load cert entries, we want to store the last seen delta WAL
	// serial number. The subsequent List will have at LEAST that certificate
	// (and potentially more) in it; when we're done writing the delta CRL,
	// we'll write this serial as a sentinel to see if we need to rebuild it
	// in the future.
	//
	// We need to do this per-cluster.
	lastDeltaSerial := map[string]string{}
	if isDelta {
		clusters, err := sc.Storage.List(sc.Context, unifiedDeltaWALPrefix)
		if err != nil {
			return nil, nil, fmt.Errorf("error listing clusters for unified delta WAL building: %w", err)
		}

		for index, cluster := range clusters {
			path := unifiedDeltaWALPrefix + cluster + deltaWALLastRevokedSerialName
			serial, err := getLastWALSerial(sc, path)
			if err != nil {
				return nil, nil, fmt.Errorf("error getting last written Delta WAL serial for cluster (%v / %v): %w", index, cluster, err)
			}

			lastDeltaSerial[cluster] = serial
		}
	}

	// We fetch a list of delta WAL entries prior to generating the complete
	// CRL. This allows us to avoid a lock (to clear such storage): anything
	// visible now, should also be visible on the complete CRL we're writing.
	var currDeltaCerts []string
	if !isDelta {
		currDeltaCerts, err = sc.CrlBuilder().GetPresentUnifiedDeltaWALForClearing(sc)
		if err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to get present delta WAL entries for removal: %w", err)
		}
	}

	var unassignedCerts []pkix.RevokedCertificate
	var revokedCertsMap map[issuing.IssuerID][]pkix.RevokedCertificate

	// If the CRL is disabled do not bother reading in all the revoked certificates.
	if !globalCRLConfig.Disable {
		// Next, we load and parse all revoked certificates. We need to assign
		// these certificates to an issuer. Some certificates will not be
		// assignable (if they were issued by a since-deleted issuer), so we need
		// a separate pool for those.
		unassignedCerts, revokedCertsMap, err = getUnifiedRevokedCertEntries(sc, issuerIDCertMap, isDelta)
		if err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to get revoked certificate entries: %w", err)
		}

		if !isDelta {
			// Revoking an issuer forces us to rebuild our complete CRL,
			// regardless of whether or not we've enabled auto rebuilding or
			// delta CRLs. If we elide the above isDelta check, this results
			// in a non-empty delta CRL, containing the serial of the
			// now-revoked issuer, even though it was generated _after_ the
			// complete CRL with the issuer on it. There's no reason to
			// duplicate this serial number on the delta, hence the above
			// guard for isDelta.
			if err := augmentWithRevokedIssuers(issuerIDEntryMap, issuerIDCertMap, revokedCertsMap); err != nil {
				return nil, nil, fmt.Errorf("error building CRLs: unable to parse revoked issuers: %w", err)
			}
		}
	}

	// Fetch the cluster-local CRL mapping so we know where to write the
	// CRLs.
	internalCRLConfig, err := sc.getUnifiedCRLConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("error building CRLs: unable to fetch cluster-local CRL configuration: %w", err)
	}

	rebuildWarnings, err := buildAnyCRLsWithCerts(sc, issuersConfig, globalCRLConfig, internalCRLConfig,
		issuers, issuerIDEntryMap, keySubjectIssuersMap,
		unassignedCerts, revokedCertsMap,
		forceNew, true /* isUnified */, isDelta)
	if err != nil {
		return nil, nil, fmt.Errorf("error building CRLs: %w", err)
	}
	if len(rebuildWarnings) > 0 {
		warnings = append(warnings, rebuildWarnings...)
	}

	// Finally, persist our potentially updated local CRL config. Only do this
	// if we didn't have a legacy CRL bundle.
	if !wasLegacy {
		if err := sc.setUnifiedCRLConfig(internalCRLConfig); err != nil {
			return nil, nil, fmt.Errorf("error building CRLs: unable to persist updated cluster-local CRL config: %w", err)
		}
	}

	if isDelta {
		// Update our last build time here so we avoid checking for new certs
		// for a while.
		sc.CrlBuilder().SetLastDeltaRebuildCheckTime(time.Now())

		// Persist all of our known last revoked serial numbers here, as the
		// last seen serial during build. This will allow us to detect if any
		// new revocations have occurred, forcing us to rebuild the delta CRL.
		for cluster, serial := range lastDeltaSerial {
			if len(serial) == 0 {
				continue
			}

			// Make sure to use the cluster-specific path. Since we're on the
			// active node of the primary cluster, we own this entry and can
			// safely write it.
			path := unifiedDeltaWALPrefix + cluster + deltaWALLastBuildSerialName
			deltaInfo := lastDeltaInfo{Serial: serial}
			lastDeltaBuildEntry, err := logical.StorageEntryJSON(path, deltaInfo)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating last delta CRL rebuild serial entry: %w", err)
			}

			err = sc.Storage.Put(sc.Context, lastDeltaBuildEntry)
			if err != nil {
				return nil, nil, fmt.Errorf("error persisting last delta CRL rebuild info: %w", err)
			}
		}
	}

	return currDeltaCerts, warnings, nil
}

func buildAnyCRLsWithCerts(
	sc *storageContext,
	issuersConfig *issuing.IssuerConfigEntry,
	globalCRLConfig *pki_backend.CrlConfig,
	internalCRLConfig *issuing.InternalCRLConfigEntry,
	issuers []issuing.IssuerID,
	issuerIDEntryMap map[issuing.IssuerID]*issuing.IssuerEntry,
	keySubjectIssuersMap map[issuing.KeyID]map[string][]issuing.IssuerID,
	unassignedCerts []pkix.RevokedCertificate,
	revokedCertsMap map[issuing.IssuerID][]pkix.RevokedCertificate,
	forceNew bool,
	isUnified bool,
	isDelta bool,
) ([]string, error) {
	// Now we can call buildCRL once, on an arbitrary/representative issuer
	// from each of these (KeyID, subject) sets.
	var warnings []string
	for _, subjectIssuersMap := range keySubjectIssuersMap {
		for _, issuersSet := range subjectIssuersMap {
			if len(issuersSet) == 0 {
				continue
			}

			var revokedCerts []pkix.RevokedCertificate
			representative := issuing.IssuerID("")
			var crlIdentifier issuing.CrlID
			var crlIdIssuer issuing.IssuerID
			for _, issuerId := range issuersSet {
				// Skip entries which aren't enabled for CRL signing. We don't
				// particularly care which issuer is ultimately chosen as the
				// set representative for signing at this point, other than
				// that it has crl-signing usage.
				if err := issuerIDEntryMap[issuerId].EnsureUsage(issuing.CRLSigningUsage); err != nil {
					continue
				}

				// Prefer to use the default as the representative of this
				// set, if it is a member.
				//
				// If it is, we'll also pull in the unassigned certs to remain
				// compatible with Vault's earlier, potentially questionable
				// behavior.
				if issuerId == issuersConfig.DefaultIssuerId {
					if len(unassignedCerts) > 0 {
						revokedCerts = append(revokedCerts, unassignedCerts...)
					}

					representative = issuerId
				}

				// Otherwise, use any other random issuer if we've not yet
				// chosen one.
				if representative == issuing.IssuerID("") {
					representative = issuerId
				}

				// Pull in the revoked certs associated with this member.
				if thisRevoked, ok := revokedCertsMap[issuerId]; ok && len(thisRevoked) > 0 {
					revokedCerts = append(revokedCerts, thisRevoked...)
				}

				// Finally, check our crlIdentifier.
				if thisCRLId, ok := internalCRLConfig.IssuerIDCRLMap[issuerId]; ok && len(thisCRLId) > 0 {
					if len(crlIdentifier) > 0 && crlIdentifier != thisCRLId {
						return nil, fmt.Errorf("error building CRLs: two issuers with same keys/subjects (%v vs %v) have different internal CRL IDs: %v vs %v", issuerId, crlIdIssuer, thisCRLId, crlIdentifier)
					}

					crlIdentifier = thisCRLId
					crlIdIssuer = issuerId
				}
			}

			if representative == "" {
				// Skip this set for the time being; while we have valid
				// issuers and associated keys, this occurred because we lack
				// crl-signing usage on all issuers in this set.
				//
				// But, tell the user about this, so they can either correct
				// this by reissuing the CA certificate or adding an equivalent
				// version with KU bits if the CA cert lacks KU altogether.
				//
				// See also: https://github.com/hashicorp/vault/issues/20137
				warning := "Issuer equivalency set with associated keys lacked an issuer with CRL Signing KeyUsage; refusing to rebuild CRL for this group of issuers: "
				var issuers []string
				for _, issuerId := range issuersSet {
					issuers = append(issuers, issuerId.String())
				}
				warning += strings.Join(issuers, ",")

				// We only need this warning once. :-)
				if !isUnified && !isDelta {
					warnings = append(warnings, warning)
				}

				continue
			}

			if len(crlIdentifier) == 0 {
				// Create a new random UUID for this CRL if none exists.
				crlIdentifier = genCRLId()
				internalCRLConfig.CRLNumberMap[crlIdentifier] = 1
			}

			// Update all issuers in this group to set the CRL Issuer
			for _, issuerId := range issuersSet {
				internalCRLConfig.IssuerIDCRLMap[issuerId] = crlIdentifier
			}

			// We always update the CRL Number since we never want to
			// duplicate numbers and missing numbers is fine.
			crlNumber := internalCRLConfig.CRLNumberMap[crlIdentifier]
			internalCRLConfig.CRLNumberMap[crlIdentifier] += 1

			// CRLs (regardless of complete vs delta) are incrementally
			// numbered. But delta CRLs need to know the number of the
			// last complete CRL. We assume that's the previous identifier
			// if no value presently exists.
			lastCompleteNumber, haveLast := internalCRLConfig.LastCompleteNumberMap[crlIdentifier]
			if !haveLast {
				// We use the value of crlNumber for the current CRL, so
				// decrement it by one to find the last one.
				lastCompleteNumber = crlNumber - 1
			}

			// Update `LastModified`
			if isDelta {
				internalCRLConfig.DeltaLastModified = time.Now().UTC()
			} else {
				internalCRLConfig.LastModified = time.Now().UTC()
			}

			// Enforce the max CRL size guard before building the actual CRL
			if globalCRLConfig.MaxCRLEntries > -1 {
				limit := maxCRLEntriesOrDefault(globalCRLConfig.MaxCRLEntries)
				revokedCount := len(revokedCerts)
				if revokedCount > limit {
					// Also log a nasty error to get the operator's attention
					sc.Logger().Error("CRL was not updated, as it exceeds the configured max size.  The CRL now does not contain all revoked certificates!  This may be indicative of a runaway issuance/revocation pattern.", "limit", limit)
					return nil, fmt.Errorf("error building CRL: revocation list size (%d) exceeds configured maximum (%d)", revokedCount, limit)
				}
				if revokedCount > int(float32(limit)*0.90) {
					sc.Logger().Warn("warning, revoked certificate count is within 10% of the configured maximum CRL size", "revoked_certs", revokedCount, "limit", limit)
				}
			}

			// Lastly, build the CRL.
			nextUpdate, err := buildCRL(sc, globalCRLConfig, forceNew, representative, revokedCerts, crlIdentifier, crlNumber, isUnified, isDelta, lastCompleteNumber)
			if err != nil {
				return nil, fmt.Errorf("error building CRLs: unable to build CRL for issuer (%v): %w", representative, err)
			}

			internalCRLConfig.CRLExpirationMap[crlIdentifier] = *nextUpdate
			if !isDelta {
				internalCRLConfig.LastCompleteNumberMap[crlIdentifier] = crlNumber
			} else if !haveLast {
				// Since we're writing this config anyways, save our guess
				// as to the last CRL number.
				internalCRLConfig.LastCompleteNumberMap[crlIdentifier] = lastCompleteNumber
			}
		}
	}

	// Before persisting our updated CRL config, check to see if we have
	// any dangling references. If we have any issuers that don't exist,
	// remove them, remembering their CRLs IDs. If we've completely removed
	// all issuers pointing to that CRL number, we can remove it from the
	// number map and from storage.
	//
	// Note that we persist the last generated CRL for a specified issuer
	// if it is later disabled for CRL generation. This mirrors the old
	// root deletion behavior, but using soft issuer deletes. If there is an
	// alternate, equivalent issuer however, we'll keep updating the shared
	// CRL; all equivalent issuers must have their CRLs disabled.
	for mapIssuerId := range internalCRLConfig.IssuerIDCRLMap {
		stillHaveIssuer := false
		for _, listedIssuerId := range issuers {
			if mapIssuerId == listedIssuerId {
				stillHaveIssuer = true
				break
			}
		}

		if !stillHaveIssuer {
			delete(internalCRLConfig.IssuerIDCRLMap, mapIssuerId)
		}
	}
	for crlId := range internalCRLConfig.CRLNumberMap {
		stillHaveIssuerForID := false
		for _, remainingCRL := range internalCRLConfig.IssuerIDCRLMap {
			if remainingCRL == crlId {
				stillHaveIssuerForID = true
				break
			}
		}

		if !stillHaveIssuerForID {
			if err := sc.Storage.Delete(sc.Context, issuing.PathCrls+crlId.String()); err != nil {
				return nil, fmt.Errorf("error building CRLs: unable to clean up deleted issuers' CRL: %w", err)
			}
		}
	}

	// All good :-)
	return warnings, nil
}

func isRevInfoIssuerValid(revInfo *revocation.RevocationInfo, issuerIDCertMap map[issuing.IssuerID]*x509.Certificate) bool {
	if len(revInfo.CertificateIssuer) > 0 {
		issuerId := revInfo.CertificateIssuer
		if _, issuerExists := issuerIDCertMap[issuerId]; issuerExists {
			return true
		}
	}

	return false
}

func getLocalRevokedCertEntries(sc *storageContext, issuerIDCertMap map[issuing.IssuerID]*x509.Certificate, isDelta bool) ([]pkix.RevokedCertificate, map[issuing.IssuerID][]pkix.RevokedCertificate, error) {
	var unassignedCerts []pkix.RevokedCertificate
	revokedCertsMap := make(map[issuing.IssuerID][]pkix.RevokedCertificate)

	listingPath := revokedPath
	if isDelta {
		listingPath = localDeltaWALPath
	}

	revokedSerials, err := sc.Storage.List(sc.Context, listingPath)
	if err != nil {
		return nil, nil, errutil.InternalError{Err: fmt.Sprintf("error fetching list of revoked certs: %s", err)}
	}

	// Build a mapping of issuer serial -> certificate.
	issuerSerialCertMap := make(map[string][]*x509.Certificate, len(issuerIDCertMap))
	for _, cert := range issuerIDCertMap {
		serialStr := serialFromCert(cert)
		issuerSerialCertMap[serialStr] = append(issuerSerialCertMap[serialStr], cert)
	}

	for _, serial := range revokedSerials {
		if isDelta && (serial == deltaWALLastBuildSerialName || serial == deltaWALLastRevokedSerialName) {
			// Skip our placeholder entries...
			continue
		}

		var revInfo revocation.RevocationInfo
		revokedEntry, err := sc.Storage.Get(sc.Context, revokedPath+serial)
		if err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unable to fetch revoked cert with serial %s: %s", serial, err)}
		}

		if revokedEntry == nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("revoked certificate entry for serial %s is nil", serial)}
		}
		if revokedEntry.Value == nil || len(revokedEntry.Value) == 0 {
			// TODO: In this case, remove it and continue? How likely is this to
			// happen? Alternately, could skip it entirely, or could implement a
			// delete function so that there is a way to remove these
			return nil, nil, errutil.InternalError{Err: "found revoked serial but actual certificate is empty"}
		}

		err = revokedEntry.DecodeJSON(&revInfo)
		if err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("error decoding revocation entry for serial %s: %s", serial, err)}
		}

		revokedCert, err := x509.ParseCertificate(revInfo.CertificateBytes)
		if err != nil {
			return nil, nil, errutil.InternalError{Err: fmt.Sprintf("unable to parse stored revoked certificate with serial %s: %s", serial, err)}
		}

		// We want to skip issuer certificate's revocationEntries for two
		// reasons:
		//
		// 1. We canonically use augmentWithRevokedIssuers to handle this
		//    case and this entry is just a backup. This prevents the issue
		//    of duplicate serial numbers on the CRL from both paths.
		// 2. We want to avoid a root's serial from appearing on its own
		//    CRL. If it is a cross-signed or re-issued variant, this is OK,
		//    but in the case we mark the root itself as "revoked", we want
		//    to avoid it appearing on the CRL as that is definitely
		//    undefined/little-supported behavior.
		//
		// This hash map lookup should be faster than byte comparison against
		// each issuer proactively.
		if candidates, present := issuerSerialCertMap[serialFromCert(revokedCert)]; present {
			revokedCertIsIssuer := false
			for _, candidate := range candidates {
				if bytes.Equal(candidate.Raw, revokedCert.Raw) {
					revokedCertIsIssuer = true
					break
				}
			}

			if revokedCertIsIssuer {
				continue
			}
		}

		// NOTE: We have to change this to UTC time because the CRL standard
		// mandates it but Go will happily encode the CRL without this.
		newRevCert := pkix.RevokedCertificate{
			SerialNumber: revokedCert.SerialNumber,
		}
		if !revInfo.RevocationTimeUTC.IsZero() {
			newRevCert.RevocationTime = revInfo.RevocationTimeUTC
		} else {
			newRevCert.RevocationTime = time.Unix(revInfo.RevocationTime, 0).UTC()
		}

		// If we have a CertificateIssuer field on the revocation entry,
		// prefer it to manually checking each issuer signature, assuming it
		// appears valid. It's highly unlikely for two different issuers
		// to have the same id (after the first was deleted).
		if isRevInfoIssuerValid(&revInfo, issuerIDCertMap) {
			revokedCertsMap[revInfo.CertificateIssuer] = append(revokedCertsMap[revInfo.CertificateIssuer], newRevCert)
			continue

			// Otherwise, fall through and update the entry.
		}

		// Now we need to assign the revoked certificate to an issuer.
		foundParent := revInfo.AssociateRevokedCertWithIsssuer(revokedCert, issuerIDCertMap)
		if !foundParent {
			// If the parent isn't found, add it to the unassigned bucket.
			unassignedCerts = append(unassignedCerts, newRevCert)
		} else {
			revokedCertsMap[revInfo.CertificateIssuer] = append(revokedCertsMap[revInfo.CertificateIssuer], newRevCert)

			// When the CertificateIssuer field wasn't found on the existing
			// entry (or was invalid), and we've found a new value for it,
			// we should update the entry to make future CRL builds faster.
			revokedEntry, err = logical.StorageEntryJSON(revokedPath+serial, revInfo)
			if err != nil {
				return nil, nil, fmt.Errorf("error creating revocation entry for existing cert: %v: %w", serial, err)
			}

			err = sc.Storage.Put(sc.Context, revokedEntry)
			if err != nil {
				return nil, nil, fmt.Errorf("error updating revoked certificate at existing location: %v: %w", serial, err)
			}
		}
	}

	return unassignedCerts, revokedCertsMap, nil
}

func getUnifiedRevokedCertEntries(sc *storageContext, issuerIDCertMap map[issuing.IssuerID]*x509.Certificate, isDelta bool) ([]pkix.RevokedCertificate, map[issuing.IssuerID][]pkix.RevokedCertificate, error) {
	// Getting unified revocation entries is a bit different than getting
	// the local ones. In particular, the full copy of the certificate is
	// unavailable, so we'll be able to avoid parsing the stored certificate,
	// at the expense of potentially having incorrect issuer mappings.
	var unassignedCerts []pkix.RevokedCertificate
	revokedCertsMap := make(map[issuing.IssuerID][]pkix.RevokedCertificate)

	listingPath := unifiedRevocationReadPathPrefix
	if isDelta {
		listingPath = unifiedDeltaWALPrefix
	}

	// First, we find all clusters that have written certificates.
	clusterIds, err := sc.Storage.List(sc.Context, listingPath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list clusters for unified CRL building: %w", err)
	}

	// We wish to prevent duplicate revocations on separate clusters from
	// being added multiple times to the CRL. While we can't guarantee these
	// are the same certificate, it doesn't matter as (as long as they have
	// the same issuer), it'd imply issuance of two certs with the same
	// serial which'd be an intentional violation of RFC 5280 before importing
	// an issuer into Vault, and would be highly unlikely within Vault, due
	// to 120-bit random serial numbers.
	foundSerials := make(map[string]bool)

	// Then for every cluster, we find its revoked certificates...
	for _, clusterId := range clusterIds {
		if !strings.HasSuffix(clusterId, "/") {
			// No entries
			continue
		}

		clusterPath := listingPath + clusterId
		serials, err := sc.Storage.List(sc.Context, clusterPath)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to list serials in cluster (%v) for unified CRL building: %w", clusterId, err)
		}

		// At this point, we need the storage entry. Rather than using the
		// clusterPath and adding the serial, we need to use the true
		// cross-cluster revocation entry (as, our above listing might have
		// used delta WAL entires without the full revocation info).
		serialPrefix := unifiedRevocationReadPathPrefix + clusterId
		for _, serial := range serials {
			if isDelta && (serial == deltaWALLastBuildSerialName || serial == deltaWALLastRevokedSerialName) {
				// Skip our placeholder entries...
				continue
			}

			serialPath := serialPrefix + serial
			entryRaw, err := sc.Storage.Get(sc.Context, serialPath)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to read unified revocation entry in cluster (%v) for unified CRL building: %w", clusterId, err)
			}
			if entryRaw == nil {
				// Skip empty entries. We'll eventually tidy them.
				continue
			}

			var xRevEntry revocation.UnifiedRevocationEntry
			if err := entryRaw.DecodeJSON(&xRevEntry); err != nil {
				return nil, nil, fmt.Errorf("failed json decoding of unified revocation entry at path %v: %w ", serialPath, err)
			}

			// Convert to pkix.RevokedCertificate entries.
			var revEntry pkix.RevokedCertificate
			var ok bool
			revEntry.SerialNumber, ok = serialToBigInt(serial)
			if !ok {
				return nil, nil, fmt.Errorf("failed to encode serial for CRL building: %v", serial)
			}

			revEntry.RevocationTime = xRevEntry.RevocationTimeUTC

			if found, inFoundMap := foundSerials[normalizeSerial(serial)]; found && inFoundMap {
				// Serial has already been added to the CRL.
				continue
			}
			foundSerials[normalizeSerial(serial)] = true

			// Finally, add it to the correct mapping.
			_, present := issuerIDCertMap[xRevEntry.CertificateIssuer]
			if !present {
				unassignedCerts = append(unassignedCerts, revEntry)
			} else {
				revokedCertsMap[xRevEntry.CertificateIssuer] = append(revokedCertsMap[xRevEntry.CertificateIssuer], revEntry)
			}
		}
	}

	return unassignedCerts, revokedCertsMap, nil
}

func augmentWithRevokedIssuers(issuerIDEntryMap map[issuing.IssuerID]*issuing.IssuerEntry, issuerIDCertMap map[issuing.IssuerID]*x509.Certificate, revokedCertsMap map[issuing.IssuerID][]pkix.RevokedCertificate) error {
	// When setup our maps with the legacy CA bundle, we only have a
	// single entry here. This entry is never revoked, so the outer loop
	// will exit quickly.
	for ourIssuerID, ourIssuer := range issuerIDEntryMap {
		if !ourIssuer.Revoked {
			continue
		}

		ourCert := issuerIDCertMap[ourIssuerID]
		ourRevCert := pkix.RevokedCertificate{
			SerialNumber:   ourCert.SerialNumber,
			RevocationTime: ourIssuer.RevocationTimeUTC,
		}

		for otherIssuerID := range issuerIDEntryMap {
			if otherIssuerID == ourIssuerID {
				continue
			}

			// Find all _other_ certificates which verify this issuer,
			// allowing us to add this revoked issuer to this issuer's
			// CRL.
			otherCert := issuerIDCertMap[otherIssuerID]
			if err := ourCert.CheckSignatureFrom(otherCert); err == nil {
				// Valid signature; add our result.
				revokedCertsMap[otherIssuerID] = append(revokedCertsMap[otherIssuerID], ourRevCert)
			}
		}
	}

	return nil
}

// Builds a CRL by going through the list of revoked certificates and building
// a new CRL with the stored revocation times and serial numbers.
func buildCRL(sc *storageContext, crlInfo *pki_backend.CrlConfig, forceNew bool, thisIssuerId issuing.IssuerID, revoked []pkix.RevokedCertificate, identifier issuing.CrlID, crlNumber int64, isUnified bool, isDelta bool, lastCompleteNumber int64) (*time.Time, error) {
	var revokedCerts []pkix.RevokedCertificate

	crlLifetime, err := parseutil.ParseDurationSecond(crlInfo.Expiry)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error parsing CRL duration of %s", crlInfo.Expiry)}
	}

	if crlInfo.Disable {
		if !forceNew {
			// In the event of a disabled CRL, we'll have the next time set
			// to the zero time as a sentinel in case we get re-enabled.
			return &time.Time{}, nil
		}

		// NOTE: in this case, the passed argument (revoked) is not added
		// to the revokedCerts list. This is because we want to sign an
		// **empty** CRL (as the CRL was disabled but we've specified the
		// forceNew option). In previous versions of Vault (1.10 series and
		// earlier), we'd have queried the certs below, whereas we now have
		// an assignment from a pre-queried list.
		goto WRITE
	}

	revokedCerts = revoked

WRITE:
	signingBundle, caErr := sc.fetchCAInfoByIssuerId(thisIssuerId, issuing.CRLSigningUsage)
	if caErr != nil {
		switch caErr.(type) {
		case errutil.UserError:
			return nil, errutil.UserError{Err: fmt.Sprintf("could not fetch the CA certificate: %s", caErr)}
		default:
			return nil, errutil.InternalError{Err: fmt.Sprintf("error fetching CA certificate: %s", caErr)}
		}
	}

	now := time.Now()
	nextUpdate := now.Add(crlLifetime)

	var extensions []pkix.Extension
	if isDelta {
		ext, err := certutil.CreateDeltaCRLIndicatorExt(lastCompleteNumber)
		if err != nil {
			return nil, fmt.Errorf("could not create crl delta indicator extension: %w", err)
		}
		extensions = []pkix.Extension{ext}
	}

	revocationListTemplate := &x509.RevocationList{
		RevokedCertificates: revokedCerts,
		Number:              big.NewInt(crlNumber),
		ThisUpdate:          now,
		NextUpdate:          nextUpdate,
		SignatureAlgorithm:  signingBundle.RevocationSigAlg,
		ExtraExtensions:     extensions,
	}

	crlBytes, err := x509.CreateRevocationList(rand.Reader, revocationListTemplate, signingBundle.Certificate, signingBundle.PrivateKey)
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error creating new CRL: %s", err)}
	}

	writePath := issuing.PathCrls + identifier.String()
	if thisIssuerId == legacyBundleShimID {
		// Ignore the CRL ID as it won't be persisted anyways; hard-code the
		// old legacy path and allow it to be updated.
		writePath = legacyCRLPath
	} else {
		if isUnified {
			writePath = unifiedCRLPathPrefix + writePath
		}

		if isDelta {
			// Write the delta CRL to a unique storage location.
			writePath += deltaCRLPathSuffix
		}
	}

	err = sc.Storage.Put(sc.Context, &logical.StorageEntry{
		Key:   writePath,
		Value: crlBytes,
	})
	if err != nil {
		return nil, errutil.InternalError{Err: fmt.Sprintf("error storing CRL: %s", err)}
	}

	return &nextUpdate, nil
}

// shouldLocalPathsUseUnified assuming a legacy path for a CRL/OCSP request, does our
// configuration say we should be returning the unified response or not
func shouldLocalPathsUseUnified(cfg *pki_backend.CrlConfig) bool {
	return cfg.UnifiedCRL && cfg.UnifiedCRLOnExistingPaths
}
