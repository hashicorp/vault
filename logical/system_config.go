package logical

import "time"

// SystemView exposes system configuration information in a safe way
// for logical backends to consume
type SystemView interface {

	// DefaultLeaseTTL returns the default lease TTL set in Vault configuration
	DefaultLeaseTTL() time.Duration

	// MaxLeaseTTL returns the max lease TTL set in Vault configuration; backend
	// authors should take care not to issue credentials that last longer than
	// this value, as Vault will revoke them
	MaxLeaseTTL() time.Duration
}

type DefaultSystemView struct {
	DefaultLeaseTTLFunc func() time.Duration
	MaxLeaseTTLFunc     func() time.Duration
}

func (d *DefaultSystemView) DefaultLeaseTTL() time.Duration {
	return d.DefaultLeaseTTLFunc()
}

func (d *DefaultSystemView) MaxLeaseTTL() time.Duration {
	return d.MaxLeaseTTLFunc()
}
