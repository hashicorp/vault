// +build !enterprise

package vault

import (
	"context"

	"github.com/hashicorp/vault/command/server"
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
	c.allLoggers = append(c.allLoggers, sealUnwrapperLogger)
	c.sealUnwrapper = NewSealUnwrapper(phys, sealUnwrapperLogger)
	// Wrap the physical backend in a cache layer if enabled
	cacheLogger := c.baseLogger.Named("storage.cache")
	c.allLoggers = append(c.allLoggers, cacheLogger)
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

// GetCoreConfigInternal returns the server configuration
// in struct format.
func (c *Core) GetCoreConfigInternal() *server.Config {
	conf := c.rawConfig.Load()
	if conf == nil {
		return nil
	}
	return conf.(*server.Config)
}

func (c *Core) teardownReplicationResolverHandler() {}
func createSecondaries(*Core, *CoreConfig)          {}

func addExtraLogicalBackends(*Core, map[string]logical.Factory) {}

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

func loadMFAConfigs(context.Context, *Core) error { return nil }

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

func (c *Core) namepaceByPath(string) *namespace.Namespace {
	return namespace.RootNamespace
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

func (c *Core) applyLeaseCountQuota(in *quotas.Request) (*quotas.Response, error) {
	return &quotas.Response{Allowed: true}, nil
}

func (c *Core) ackLeaseQuota(access quotas.Access, leaseGenerated bool) error {
	return nil
}

func (c *Core) quotaLeaseWalker(ctx context.Context, callback func(request *quotas.Request) bool) error {
	return nil
}

func (c *Core) quotasHandleLeases(ctx context.Context, action quotas.LeaseAction, leaseIDs []string) error {
	return nil
}

func (c *Core) namespaceByPath(path string) *namespace.Namespace {
	return namespace.RootNamespace
}

func (c *Core) AllowForwardingViaHeader() bool {
	return false
}

func (c *Core) MissingRequiredState(raw []string, perfStandby bool) bool {
	return false
}

func DiagnoseCheckLicense(ctx context.Context, vaultCore *Core, coreConfig CoreConfig, generate bool) (bool, []string) {
	return false, nil
}
