package vault

import (
	"fmt"
	"strings"
	"time"
)

type dynamicSystemView struct {
	core *Core
	path string
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
	// Ensure we end the path in a slash
	if !strings.HasSuffix(d.path, "/") {
		d.path += "/"
	}

	me := d.core.router.MatchingMountEntry(d.path)
	if me == nil {
		return 0, 0, fmt.Errorf("[ERR] core: failed to get mount entry for %s", d.path)
	}

	def = d.core.defaultLeaseTTL
	max = d.core.maxLeaseTTL

	if me.Config.DefaultLeaseTTL != nil && *me.Config.DefaultLeaseTTL != 0 {
		def = *me.Config.DefaultLeaseTTL
	}
	if me.Config.MaxLeaseTTL != nil && *me.Config.MaxLeaseTTL != 0 {
		max = *me.Config.MaxLeaseTTL
	}

	return
}
