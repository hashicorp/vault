// +build !enterprise

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/helper/license"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/vault/replication"
	cache "github.com/patrickmn/go-cache"
)

type entCore struct{}

type LicensingConfig struct{}

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
		c.physical = physical.NewTransactionalCache(c.sealUnwrapper, conf.CacheSize, cacheLogger)
	} else {
		c.physical = physical.NewCache(c.sealUnwrapper, conf.CacheSize, cacheLogger)
	}
	c.physicalCache = c.physical.(physical.ToggleablePurgemonster)

	// Wrap in encoding checks
	if !conf.DisableKeyEncodingChecks {
		c.physical = physical.NewStorageEncoding(c.physical)
	}
	return nil
}

func createSecondaries(*Core, *CoreConfig) {}

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

func (c *Core) namepaceByPath(string) *namespace.Namespace {
	return namespace.RootNamespace
}

func (c *Core) setupReplicatedClusterPrimary(*replication.Cluster) error { return nil }

func (c *Core) perfStandbyCount() int { return 0 }

func (c *Core) removePrefixFromFilteredPaths(context.Context, string) error {
	return nil
}

func (c *Core) checkReplicatedFiltering(context.Context, *MountEntry, string) (bool, error) {
	return false, nil
}

func (c *Core) invalidateSentinelPolicy(PolicyType, string) {}

func (c *Core) removePerfStandbySecondary(context.Context, string) {}

func (c *Core) perfStandbyClusterHandler() (*replication.Cluster, *cache.Cache, chan struct{}, error) {
	return nil, cache.New(2*HeartbeatInterval, 1*time.Second), make(chan struct{}), nil
}
