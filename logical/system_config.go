package logical

import "time"

// SystemConfig exposes system configuration information in a safe way
// for other logical backends to consume
type SystemConfig struct {

	// DefaultLeaseTTL returns the default lease TTL set in Vault configuration
	DefaultLeaseTTL func() time.Duration

	// MaxLeaseTTL returns the max lease TTL set in Vault configuration; backend
	// authors should take care not to issue credentials that last longer than
	// this value, as Vault will revoke them
	MaxLeaseTTL func() time.Duration
}
