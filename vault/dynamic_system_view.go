package vault

import (
	"time"

	"github.com/hashicorp/vault/logical"
)

type dynamicSystemView struct {
	core       *Core
	mountEntry *MountEntry
}

func (d dynamicSystemView) DefaultLeaseTTL() time.Duration {
	def, _ := d.fetchTTLs()
	return def
}

func (d dynamicSystemView) MaxLeaseTTL() time.Duration {
	_, max := d.fetchTTLs()
	return max
}

func (d dynamicSystemView) SudoPrivilege(path string, token string) bool {
	// Resolve the token policy
	te, err := d.core.tokenStore.Lookup(token)
	if err != nil {
		d.core.logger.Error("core: failed to lookup token", "error", err)
		return false
	}

	// Ensure the token is valid
	if te == nil {
		d.core.logger.Error("entry not found for given token")
		return false
	}

	// Construct the corresponding ACL object
	acl, err := d.core.policyStore.ACL(te.Policies...)
	if err != nil {
		d.core.logger.Error("failed to retrieve ACL for token's policies", "token_policies", te.Policies, "error", err)
		return false
	}

	// The operation type isn't important here as this is run from a path the
	// user has already been given access to; we only care about whether they
	// have sudo
	_, rootPrivs := acl.AllowOperation(logical.ReadOperation, path)
	return rootPrivs
}

// TTLsByPath returns the default and max TTLs corresponding to a particular
// mount point, or the system default
func (d dynamicSystemView) fetchTTLs() (def, max time.Duration) {
	def = d.core.defaultLeaseTTL
	max = d.core.maxLeaseTTL

	if d.mountEntry.Config.DefaultLeaseTTL != 0 {
		def = d.mountEntry.Config.DefaultLeaseTTL
	}
	if d.mountEntry.Config.MaxLeaseTTL != 0 {
		max = d.mountEntry.Config.MaxLeaseTTL
	}

	return
}

// Tainted indicates that the mount is in the process of being removed
func (d dynamicSystemView) Tainted() bool {
	return d.mountEntry.Tainted
}

// CachingDisabled indicates whether to use caching behavior
func (d dynamicSystemView) CachingDisabled() bool {
	return d.core.cachingDisabled
}
