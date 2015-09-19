package vault

import "time"

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

func (d dynamicSystemView) SudoPrivilege(path string, policies []string) bool {
	acl, err := d.core.policy.ACL(policies...)
	if err != nil {
		d.core.logger.Printf("[ERR] Failed to retrieve ACL for policies [%#v]: %s", policies, err)
		return false
	}
	return acl.RootPrivilege(path)
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
