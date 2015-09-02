package logical

import "time"

// SystemView exposes system configuration information in a safe way
// for logical backends to consume
type SystemView interface {
	// DefaultLeaseTTL returns the default lease TTL set in Vault configuration
	DefaultLeaseTTL() (time.Duration, error)

	// MaxLeaseTTL returns the max lease TTL set in Vault configuration; backend
	// authors should take care not to issue credentials that last longer than
	// this value, as Vault will revoke them
	MaxLeaseTTL() (time.Duration, error)
}

type StaticSystemView struct {
	DefaultLeaseTTLVal time.Duration
	MaxLeaseTTLVal     time.Duration
}

func (d StaticSystemView) DefaultLeaseTTL() (time.Duration, error) {
	return d.DefaultLeaseTTLVal, nil
}

func (d StaticSystemView) MaxLeaseTTL() (time.Duration, error) {
	return d.MaxLeaseTTLVal, nil
}
