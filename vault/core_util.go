// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/helper/license"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/vault/quotas"
	"github.com/hashicorp/vault/vault/replication"
)

const (
	activityLogEnabledDefault      = false
	activityLogEnabledDefaultValue = "default-disabled"
)

type (
	entCore       struct{}
	entCoreConfig struct{}
)

func (e entCoreConfig) Clone() entCoreConfig {
	return entCoreConfig{}
}

type LicensingConfig struct {
	AdditionalPublicKeys []interface{}
}

func coreInit(c *Core, conf *CoreConfig) error {
	phys := conf.Physical
	_, txnOK := phys.(physical.Transactional)
	sealUnwrapperLogger := conf.Logger.Named("storage.sealunwrapper")
	c.sealUnwrapper = NewSealUnwrapper(phys, sealUnwrapperLogger)
	// Wrap the physical backend in a cache layer if enabled
	cacheLogger := c.baseLogger.Named("storage.cache")
	if txnOK {
		c.physical = physical.NewTransactionalCache(c.sealUnwrapper, conf.CacheSize, cacheLogger, c.MetricSink().Sink)
	} else {
		c.physical = physical.NewCache(c.sealUnwrapper, conf.CacheSize, cacheLogger, c.MetricSink().Sink)
	}
	c.physicalCache = c.physical.(physical.ToggleablePurgemonster)

	// Wrap in encoding checks
	if !conf.DisableKeyEncodingChecks {
		c.physical = physical.NewStorageEncoding(c.physical)
	}

	return nil
}

func (c *Core) setupReplicationResolverHandler() error {
	return nil
}

func NewPolicyMFABackend(core *Core, logger hclog.Logger) *PolicyMFABackend { return nil }

func (c *Core) barrierViewForNamespace(namespaceId string) (*BarrierView, error) {
	if namespaceId != namespace.RootNamespaceID {
		return nil, fmt.Errorf("failed to find barrier view for non-root namespace")
	}

	return c.systemBarrierView, nil
}

func (c *Core) UndoLogsEnabled() bool            { return false }
func (c *Core) UndoLogsPersisted() (bool, error) { return false, nil }
func (c *Core) PersistUndoLogs() error           { return nil }

func (c *Core) teardownReplicationResolverHandler() {}
func createSecondaries(*Core, *CoreConfig)          {}

func addExtraLogicalBackends(*Core, map[string]logical.Factory, string) {}

func addExtraCredentialBackends(*Core, map[string]logical.Factory) {}

func preUnsealInternal(context.Context, *Core) error { return nil }

func postSealInternal(*Core) {}

func preSealPhysical(c *Core) {
	switch c.sealUnwrapper.(type) {
	case *sealUnwrapper:
		c.sealUnwrapper.(*sealUnwrapper).stopUnwraps()
	case *transactionalSealUnwrapper:
		c.sealUnwrapper.(*transactionalSealUnwrapper).stopUnwraps()
	}

	// Purge the cache
	c.physicalCache.SetEnabled(false)
	c.physicalCache.Purge(context.Background())
}

func postUnsealPhysical(c *Core) error {
	switch c.sealUnwrapper.(type) {
	case *sealUnwrapper:
		c.sealUnwrapper.(*sealUnwrapper).runUnwraps()
	case *transactionalSealUnwrapper:
		c.sealUnwrapper.(*transactionalSealUnwrapper).runUnwraps()
	}
	return nil
}

func loadPolicyMFAConfigs(context.Context, *Core) error { return nil }

func shouldStartClusterListener(*Core) bool { return true }

func hasNamespaces(*Core) bool { return false }

func (c *Core) Features() license.Features {
	return license.FeatureNone
}

func (c *Core) HasFeature(license.Features) bool {
	return false
}

func (c *Core) collectNamespaces() []*namespace.Namespace {
	return []*namespace.Namespace{
		namespace.RootNamespace,
	}
}

func (c *Core) HasWALState(required *logical.WALState, perfStandby bool) bool {
	return true
}

func (c *Core) setupReplicatedClusterPrimary(*replication.Cluster) error { return nil }

func (c *Core) perfStandbyCount() int { return 0 }

func (c *Core) removePathFromFilteredPaths(context.Context, string, string) error {
	return nil
}

func (c *Core) checkReplicatedFiltering(context.Context, *MountEntry, string) (bool, error) {
	return false, nil
}

func (c *Core) invalidateSentinelPolicy(PolicyType, string) {}

func (c *Core) removePerfStandbySecondary(context.Context, string) {}

func (c *Core) removeAllPerfStandbySecondaries() {}

func (c *Core) perfStandbyClusterHandler() (*replication.Cluster, chan struct{}, error) {
	return nil, make(chan struct{}), nil
}

func (c *Core) initSealsForMigration() {}

func (c *Core) postSealMigration(ctx context.Context) error { return nil }

func (c *Core) applyLeaseCountQuota(_ context.Context, in *quotas.Request) (*quotas.Response, error) {
	return &quotas.Response{Allowed: true}, nil
}

func (c *Core) ackLeaseQuota(access quotas.Access, leaseGenerated bool) error {
	return nil
}

func (c *Core) quotaLeaseWalker(ctx context.Context, callback func(request *quotas.Request) bool) error {
	return nil
}

func (c *Core) quotasHandleLeases(ctx context.Context, action quotas.LeaseAction, leases []*quotas.QuotaLeaseInformation) error {
	return nil
}

func (c *Core) namespaceByPath(path string) *namespace.Namespace {
	return namespace.RootNamespace
}

func (c *Core) AllowForwardingViaHeader() bool {
	return false
}

func (c *Core) ForwardToActive() string {
	return ""
}

func (c *Core) MissingRequiredState(raw []string, perfStandby bool) bool {
	return false
}

func DiagnoseCheckLicense(ctx context.Context, vaultCore *Core, coreConfig CoreConfig, generate bool) (bool, []string) {
	return false, nil
}

func (c *Core) initializeBarrier(ctx context.Context, barrierConfig *SealConfig) (*InitResult, error) {
	rootKey, rootKeyShares, err := c.generateShares(barrierConfig)
	if err != nil {
		c.logger.Error("error generating shares", "error", err)
		return nil, err
	}

	var sealKey []byte
	var sealKeyShares [][]byte

	if barrierConfig.StoredShares == 1 && c.seal.BarrierType() == wrapping.WrapperTypeShamir {
		sealKey, sealKeyShares, err = c.generateShares(barrierConfig)
		if err != nil {
			c.logger.Error("error generating shares", "error", err)
			return nil, err
		}
	}

	// Initialize the barrier
	if err := c.barrier.Initialize(ctx, rootKey, sealKey, c.secureRandomReader); err != nil {
		c.logger.Error("failed to initialize barrier", "error", err)
		return nil, fmt.Errorf("failed to initialize barrier: %w", err)
	}
	if c.logger.IsInfo() {
		c.logger.Info("security barrier initialized", "stored", barrierConfig.StoredShares, "shares", barrierConfig.SecretShares, "threshold", barrierConfig.SecretThreshold)
	}

	// Unseal the barrier
	if err := c.barrier.Unseal(ctx, rootKey); err != nil {
		c.logger.Error("failed to unseal barrier", "error", err)
		return nil, fmt.Errorf("failed to unseal barrier: %w", err)
	}

	err := c.seal.SetBarrierConfig(ctx, barrierConfig)
	if err != nil {
		c.logger.Error("failed to save barrier configuration", "error", err)
		return fmt.Errorf("barrier configuration saving failed: %w", err)
	}

	results := &InitResult{
		SecretShares: [][]byte{},
	}

	// If we are storing shares, pop them out of the returned results and push
	// them through the seal
	switch c.seal.StoredKeysSupported() {
	case seal.StoredKeysSupportedShamirRoot:
		keysToStore := [][]byte{barrierKey}
		if err := c.seal.GetAccess().SetShamirSealKey(sealKey); err != nil {
			c.logger.Error("failed to set seal key", "error", err)
			return nil, fmt.Errorf("failed to set seal key: %w", err)
		}
		if err := c.seal.SetStoredKeys(ctx, keysToStore); err != nil {
			c.logger.Error("failed to store keys", "error", err)
			return nil, fmt.Errorf("failed to store keys: %w", err)
		}
		results.SecretShares = sealKeyShares
	case seal.StoredKeysSupportedGeneric:
		keysToStore := [][]byte{barrierKey}
		if err := c.seal.SetStoredKeys(ctx, keysToStore); err != nil {
			c.logger.Error("failed to store keys", "error", err)
			return nil, fmt.Errorf("failed to store keys: %w", err)
		}
	default:
		// We don't support initializing an old-style Shamir seal anymore, so
		// this case is only reachable by tests.
		results.SecretShares = barrierKeyShares
	}

	results := &InitResult{}

	// Save the configuration regardless, but only generate a key if it's not
	// disabled. When using recovery keys they are stored in the barrier, so
	// this must happen post-unseal.
	if c.seal.RecoveryKeySupported() {
		err = c.seal.SetRecoveryConfig(ctx, recoveryConfig)
		if err != nil {
			c.logger.Error("failed to save recovery configuration", "error", err)
			return nil, fmt.Errorf("recovery configuration saving failed: %w", err)
		}

		if recoveryConfig.SecretShares > 0 {
			recoveryKey, recoveryUnsealKeys, err := c.generateShares(recoveryConfig)
			if err != nil {
				c.logger.Error("failed to generate recovery shares", "error", err)
				return nil, err
			}

			err = c.seal.SetRecoveryKey(ctx, recoveryKey)
			if err != nil {
				return nil, err
			}

			results.RecoveryShares[recoveryConfig.Name] = recoveryUnsealKeys
		}
	}

	return results, nil
}
