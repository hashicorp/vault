package framework

import (
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
)

// LeaseExtend returns an OperationFunc that can be used to simply extend the
// lease of the auth/secret for the duration that was requested.
//
// backendIncrement is the backend's requested increment -- perhaps from a user
// request, perhaps from a role/config value. If not set, uses the mount/system
// value.
//
// backendMax is the backend's requested increment -- this can be more
// restrictive than the mount/system value but not less.
//
// systemView is the system view from the calling backend, used to determine
// and/or correct default/max times.
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
		maxValidTime := leaseOpts.IssueTime.Add(max)

		// Get the current time
		now := time.Now()

		// If we are past the max TTL, we shouldn't be in this function...but
		// fast path out if we are
		if maxValidTime.Before(now) {
			return nil, fmt.Errorf("past the max TTL, cannot renew")
		}

		// Basic max safety checks have passed, now let's figure out our
		// increment. We'll use the user-supplied value first, then backend-provided default if possible, or the
		// mount/system default if not.
		increment := leaseOpts.Increment
		if increment <= 0 {
			if backendIncrement > 0 {
				increment = backendIncrement
			} else {
				increment = systemView.DefaultLeaseTTL()
			}
		}

		// We are proposing a time of the current time plus the increment
		proposedExpiration := now.Add(increment)

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
