package framework

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestLeaseExtend(t *testing.T) {
	now := time.Now().UTC().Round(time.Hour)

	cases := map[string]struct {
		Max          time.Duration
		MaxSession   time.Duration
		Request      time.Duration
		Result       time.Duration
		MaxFromLease bool
		Error        bool
	}{
		"valid request, good bounds": {
			Max:     30 * time.Hour,
			Request: 1 * time.Hour,
			Result:  1 * time.Hour,
		},

		"valid request, zero max": {
			Max:     0,
			Request: 1 * time.Hour,
			Result:  1 * time.Hour,
		},

		"request is zero": {
			Max:     30 * time.Hour,
			Request: 0,
			Result:  30 * time.Hour,
		},

		"request is too long": {
			Max:     3 * time.Hour,
			Request: 7 * time.Hour,
			Result:  3 * time.Hour,
		},

		"request would go past max session": {
			Max:        9 * time.Hour,
			MaxSession: 5 * time.Hour,
			Request:    7 * time.Hour,
			Result:     5 * time.Hour,
		},

		"request within max session": {
			Max:        9 * time.Hour,
			MaxSession: 5 * time.Hour,
			Request:    4 * time.Hour,
			Result:     4 * time.Hour,
		},

		// Don't think core will allow this, but let's protect against
		// it at multiple layers anyways.
		"request is negative": {
			Max:     3 * time.Hour,
			Request: -7 * time.Hour,
			Error:   true,
		},

		"max form lease, request too large": {
			Request:      10 * time.Hour,
			MaxFromLease: true,
			Result:       time.Hour,
		},
	}

	for name, tc := range cases {
		req := &logical.Request{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					Lease:          1 * time.Hour,
					LeaseIssue:     now,
					LeaseIncrement: tc.Request,
				},
			},
		}

		callback := LeaseExtend(tc.Max, tc.MaxSession, tc.MaxFromLease)
		resp, err := callback(req, nil)
		if (err != nil) != tc.Error {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}
		if tc.Error {
			continue
		}

		// Round it to the nearest hour
		lease := now.Add(resp.Auth.Lease).Round(time.Hour).Sub(now)
		if lease != tc.Result {
			t.Fatalf("bad: %s\nlease: %s", name, lease)
		}
	}
}
