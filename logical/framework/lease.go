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
func LeaseExtend(max time.Duration) OperationFunc {
	return func(req *logical.Request, data *FieldData) (*logical.Response, error) {
		lease := detectLease(req)
		if lease == nil {
			return nil, fmt.Errorf("no lease options for request")
		}

		// Protect against negative leases
		if lease.LeaseIncrement < 0 {
			return logical.ErrorResponse(
				"increment must be greater than 0"), logical.ErrInvalidRequest
		}

		// If the lease is zero, then assume max
		if lease.LeaseIncrement == 0 {
			lease.LeaseIncrement = max
		}

		// Determine the requested lease
		newLease := lease.IncrementedLease(lease.LeaseIncrement)

		if max > 0 {
			// Determine if the requested lease is too long
			now := time.Now().UTC()
			maxExpiration := now.Add(max)
			newExpiration := now.Add(newLease)
			if newExpiration.Sub(maxExpiration) > 0 {
				// The new expiration is past the max expiration. In this
				// case, admit the longest lease we can.
				newLease = maxExpiration.Sub(lease.ExpirationTime())
			}
		}

		// Set the lease
		lease.Lease = newLease
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
