package vault

import (
	"log"
	"time"
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

func (d dynamicSystemView) SudoPrivilege(path, policy string) bool {
	// Special "root" policy name can never be overwritten and it always will
	// have all the privileges
	if policy == "root" {
		return true
	}

	// Get the associated policy from core's PolicyStore
	p, err := d.core.policy.GetPolicy(policy)
	if err != nil {
		log.Printf("[WARN] Failed to retrieve policy '%s': %s", policy, err)
		return false
	}

	// Look all the paths in the policy object to find an entry for given path
	// and check its respective policy.
	for _, item := range p.Paths {
		if item.Prefix == path && item.Policy == PathPolicySudo {
			return true
		}
	}

	return false
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
