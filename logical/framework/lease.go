package framework

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/vault/logical"
)

// LeaseExtend is left for backwards compatibility for plugins. This function
// now just passes back the data that was passed into it to be processed in core.
// DEPRECATED
func LeaseExtend(backendIncrement, backendMax time.Duration, systemView logical.SystemView) OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *FieldData) (*logical.Response, error) {
		switch {
		case req.Auth != nil:
			req.Auth.TTL = backendIncrement
			req.Auth.MaxTTL = backendMax
			return &logical.Response{Auth: req.Auth}, nil
		case req.Secret != nil:
			req.Secret.TTL = backendIncrement
			req.Secret.MaxTTL = backendMax
			return &logical.Response{Secret: req.Secret}, nil
		}
		return nil, fmt.Errorf("no lease options for request")
	}
}

// CalculateTTL takes all the user-specified, backend, and system inputs and calculates
// a TTL for a lease
func CalculateTTL(sysView logical.SystemView, increment, backendTTL, period, backendMaxTTL, explicitMaxTTL time.Duration, startTime time.Time) (ttl time.Duration, warnings []string, errors error) {
	// Truncate all times to the second since that is the lowest precision for
	// TTLs
	now := time.Now().Truncate(time.Second)
	if startTime.IsZero() {
		startTime = now
	} else {
		startTime = startTime.Truncate(time.Second)
	}

	// Use the mount's configured max unless the backend specifies
	// something more restrictive (perhaps from a role configuration
	// parameter)
	maxTTL := sysView.MaxLeaseTTL()
	if backendMaxTTL > 0 && backendMaxTTL < maxTTL {
		maxTTL = backendMaxTTL
	}
	if explicitMaxTTL > 0 && explicitMaxTTL < maxTTL {
		maxTTL = explicitMaxTTL
	}

	// Should never happen, but guard anyways
	if maxTTL <= 0 {
		return 0, nil, fmt.Errorf("max TTL must be greater than zero")
	}

	var maxValidTime time.Time
	switch {
	case period > 0:
		// Cap the period value to the sys max_ttl value
		if period > maxTTL {
			warnings = append(warnings,
				fmt.Sprintf("period of %q exceeded the effective max_ttl of %q; period value is capped accordingly", period, maxTTL))
			period = maxTTL
		}
		ttl = period

		if explicitMaxTTL > 0 {
			maxValidTime = startTime.Add(explicitMaxTTL)
		}
	default:
		switch {
		case increment > 0:
			ttl = increment
		case backendTTL > 0:
			ttl = backendTTL
		default:
			ttl = sysView.DefaultLeaseTTL()
		}

		// We cannot go past this time
		maxValidTime = startTime.Add(maxTTL)
	}

	if !maxValidTime.IsZero() {
		// Determine the max valid TTL
		maxValidTTL := maxValidTime.Sub(now)

		// If we are past the max TTL, we shouldn't be in this function...but
		// fast path out if we are
		if maxValidTTL < 0 {
			return 0, nil, fmt.Errorf("past the max TTL, cannot renew")
		}

		// If the proposed expiration is after the maximum TTL of the lease,
		// cap the increment to whatever is left
		if maxValidTTL-ttl < 0 {
			warnings = append(warnings,
				fmt.Sprintf("TTL of %q exceeded the effective max_ttl of %q; TTL value is capped accordingly", ttl, maxValidTTL))
			ttl = maxValidTTL
		}
	}

	return ttl, warnings, nil
}
