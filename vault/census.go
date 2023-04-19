//go:build !enterprise

package vault

import "time"

// CensusAgent is a stub for OSS
type CensusAgent struct {
	billingStart time.Time
}
