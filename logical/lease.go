package logical

import (
	"time"
)

// LeaseOptions is an embeddable struct to capture common lease
// settings between a Secret and Auth
type LeaseOptions struct {
	// TTL is the duration that this secret is valid for. Vault
	// will automatically revoke it after the duration.
	TTL time.Duration `json:"lease"`

	// MaxTTL is the maximum duration that this secret is valid for.
	MaxTTL time.Duration `json:"max_ttl"`

	// Renewable, if true, means that this secret can be renewed.
	Renewable bool `json:"renewable"`

	// Increment will be the lease increment that the user requested.
	// This is only available on a Renew operation and has no effect
	// when returning a response.
	Increment time.Duration `json:"-"`

	// IssueTime is the time of issue for the original lease. This is
	// only available on a Renew operation and has no effect when returning
	// a response. It can be used to enforce maximum lease periods by
	// a logical backend.
	IssueTime time.Time `json:"-"`
}

// LeaseEnabled checks if leasing is enabled
func (l *LeaseOptions) LeaseEnabled() bool {
	return l.TTL > 0
}

// LeaseTotal is the lease duration with a guard against a negative TTL
func (l *LeaseOptions) LeaseTotal() time.Duration {
	if l.TTL <= 0 {
		return 0
	}

	return l.TTL
}

// ExpirationTime computes the time until expiration including the grace period
func (l *LeaseOptions) ExpirationTime() time.Time {
	var expireTime time.Time
	if l.LeaseEnabled() {
		expireTime = time.Now().Add(l.LeaseTotal())
	}
	return expireTime
}
