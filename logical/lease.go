package logical

import "time"

// LeaseOptions is an embeddable struct to capture common lease
// settings between a Secret and Auth
type LeaseOptions struct {
	// Lease is the duration that this secret is valid for. Vault
	// will automatically revoke it after the duration + grace period.
	Lease            time.Duration `json:"lease"`
	LeaseGracePeriod time.Duration `json:"lease_grace_period"`

	// Renewable, if true, means that this secret can be renewed.
	Renewable bool `json:"renewable"`

	// LeaseIncrement will be the lease increment that the user requested.
	// This is only available on a Renew operation and has no effect
	// when returning a response.
	LeaseIncrement time.Duration `json:"-"`

	// LeaseIssue is the time of issue for the original lease. This is
	// only available on a Renew operation and has no effect when returning
	// a response. It can be used to enforce maximum lease periods by
	// a logical backend. This time will always be in UTC.
	LeaseIssue time.Time `json:"-"`
}

// LeaseEnabled checks if leasing is enabled
func (l *LeaseOptions) LeaseEnabled() bool {
	return l.Lease > 0
}

// LeaseTotal is the total lease time including the grace period
func (l *LeaseOptions) LeaseTotal() time.Duration {
	if l.Lease <= 0 {
		return 0
	}

	if l.LeaseGracePeriod < 0 {
		return l.Lease
	}

	return l.Lease + l.LeaseGracePeriod
}

// ExpirationTime computes the time until expiration including the grace period
func (l *LeaseOptions) ExpirationTime() time.Time {
	var expireTime time.Time
	if !l.LeaseIssue.IsZero() && l.Lease > 0 {
		expireTime = l.LeaseIssue.UTC().Add(l.LeaseTotal())
	}

	return expireTime
}

// IncrementedLease returns the lease duration that would need to set
// in order to increment the _current_ lease by the given duration
// if the auth were re-issued right now.
func (l *LeaseOptions) IncrementedLease(inc time.Duration) time.Duration {
	var result time.Duration
	expireTime := l.ExpirationTime()
	if expireTime.IsZero() {
		return result
	}

	return expireTime.Add(inc).Sub(time.Now().UTC())
}
