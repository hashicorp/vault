package vault

import "time"

type dynamicSystemView struct {
	core       *Core
	mountEntry *MountEntry
}

func (d dynamicSystemView) DefaultLeaseTTL() (time.Duration, error) {
	def, _, err := d.fetchTTLs()
	if err != nil {
		return 0, err
	}
	return def, nil
}

func (d dynamicSystemView) MaxLeaseTTL() (time.Duration, error) {
	_, max, err := d.fetchTTLs()
	if err != nil {
		return 0, err
	}
	return max, nil
}

// TTLsByPath returns the default and max TTLs corresponding to a particular
// mount point, or the system default
func (d dynamicSystemView) fetchTTLs() (def, max time.Duration, retErr error) {
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
