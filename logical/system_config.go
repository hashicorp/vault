package logical

import "time"

// LogicalSystemConfig exposes system configuration information in a safe way
// for other logical backends to consume
type SystemConfig struct {
	DefaultLeaseTTL func() time.Duration
	MaxLeaseTTL     func() time.Duration
}
