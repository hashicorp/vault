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
func LeaseExtend(max, maxSession time.Duration, maxFromLease bool) OperationFunc {
	return func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		lease := detectLease(req)
		if lease == nil {
			return nil, fmt.Errorf("no lease options for request")
		}

		// Check if we should limit max
		if maxFromLease {
			max = lease.Lease
		}

		// Sanity check the desired increment
		switch {
		// Protect against negative leases
		case lease.LeaseIncrement < 0:
			return logical.ErrorResponse(
				"increment must be greater than 0"), logical.ErrInvalidRequest

		// If no lease increment, or too large of an increment, use the max
		case max > 0 && lease.LeaseIncrement == 0, max > 0 && lease.LeaseIncrement > max:
			lease.LeaseIncrement = max
		}

		// Get the current time
		now := time.Now().UTC()

		// Check if we're passed the issue limit
		var maxSessionTime time.Time
		if maxSession > 0 {
			maxSessionTime = lease.LeaseIssue.Add(maxSession)
			if maxSessionTime.Before(now) {
				return logical.ErrorResponse(fmt.Sprintf(
					"lease can only be renewed up to %s past original issue",
					maxSession)), logical.ErrInvalidRequest
			}
		}

		// The new lease is the minimum of the requested LeaseIncrement
		// or the maxSessionTime
		requestedLease := now.Add(lease.LeaseIncrement)
		if !maxSessionTime.IsZero() && requestedLease.After(maxSessionTime) {
			requestedLease = maxSessionTime
		}

		// Determine the requested lease
		newLeaseDuration := requestedLease.Sub(now)

		// Set the lease
		lease.Lease = newLeaseDuration
		return &logical.Response{Auth: req.Auth, Secret: req.Secret}, nil
	}
}

func detectLease(req *logical.Request) *logical.LeaseOptions {
	if req.Auth != nil {
		return &req.Auth.LeaseOptions
	} else if req.Secret != nil {
		return &req.Secret.LeaseOptions
	}
	return nil
}
