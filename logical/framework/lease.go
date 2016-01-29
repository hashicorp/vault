package framework

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
)

// LeaseExtend returns an OperationFunc that can be used to simply extend
// the lease of the auth/secret for the duration that was requested. Max
// is the max time past the _current_ time that a lease can be extended. i.e.
// setting it to 2 hours forces a renewal within the next 2 hours again.
//
// maxSession is the maximum session length allowed since the original
// issue time. If this is zero, it is ignored.
//
// maxFromLease controls if the maximum renewal period comes from the existing
// lease. This means the value of `max` will be replaced with the existing
// lease duration.
func LeaseExtend(backendIncrement, backendMax time.Duration, systemView logical.SystemView) OperationFunc {
	return func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		var leaseOpts *logical.LeaseOptions
		switch {
		case req.Auth != nil:
			leaseOpts = &req.Auth.LeaseOptions
		case req.Secret != nil:
			leaseOpts = &req.Secret.LeaseOptions
		default:
			return nil, fmt.Errorf("no lease options for request")
		}

		// Use the mount's configured max unless the backend specifies
		// something more restrictive (perhaps from a role configuration
		// parameter)
		max := systemView.MaxLeaseTTL()
		if backendMax > 0 && backendMax < max {
			max = backendMax
		}

		// Should never happen, but guard anyways
		if max < 0 {
			return nil, fmt.Errorf("max TTL is negative")
		}

		// We cannot go past this time
		maxValidTime := leaseOpts.IssueTime.UTC().Add(max)

		// Get the current time
		now := time.Now().UTC()

		// If we are past the max TTL, we shouldn't be in this function...but
		// fast path out if we are
		if maxValidTime.Before(now) {
			return nil, fmt.Errorf("past the max TTL, cannot renew")
		}

		// Basic max safety checks have passed, now let's figure out our
		// increment. We'll use the backend-provided value if possible, or the
		// mount/system default if not. We won't change the LeaseOpts value,
		// just adjust accordingly.
		increment := leaseOpts.Increment
		if backendIncrement > 0 {
			increment = backendIncrement
		}
		if increment <= 0 {
			increment = systemView.DefaultLeaseTTL()
		}

		proposedExpiration := leaseOpts.IssueTime.UTC().Add(increment)

		// If the proposed expiration is after the maximum TTL of the lease,
		// cap the increment to whatever is left
		if maxValidTime.Before(proposedExpiration) {
			increment = maxValidTime.Sub(now)
		}

		// Set the lease
		leaseOpts.TTL = increment

		return &logical.Response{Auth: req.Auth, Secret: req.Secret}, nil
	}
}
