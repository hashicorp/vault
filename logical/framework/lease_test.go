package framework

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestLeaseExtend(t *testing.T) {

	testSysView := logical.StaticSystemView{
		DefaultLeaseTTLVal: 5 * time.Hour,
		MaxLeaseTTLVal:     30 * time.Hour,
	}

	now := time.Now().Round(time.Hour)

	cases := map[string]struct {
		BackendDefault time.Duration
		BackendMax     time.Duration
		Increment      time.Duration
		Result         time.Duration
		Error          bool
	}{
		"valid request, good bounds, increment is preferred": {
			BackendDefault: 30 * time.Hour,
			Increment:      1 * time.Hour,
			Result:         1 * time.Hour,
		},

		"valid request, zero backend default, uses increment": {
			BackendDefault: 0,
			Increment:      1 * time.Hour,
			Result:         1 * time.Hour,
		},

		"lease increment is zero, uses backend default": {
			BackendDefault: 30 * time.Hour,
			Increment:      0,
			Result:         30 * time.Hour,
		},

		"lease increment and default are zero, uses systemview": {
			BackendDefault: 0,
			Increment:      0,
			Result:         5 * time.Hour,
		},

		"backend max and associated request are too long": {
			BackendDefault: 40 * time.Hour,
			BackendMax:     45 * time.Hour,
			Result:         30 * time.Hour,
		},

		"all request values are larger than the system view, so the system view limits": {
			BackendDefault: 40 * time.Hour,
			BackendMax:     50 * time.Hour,
			Increment:      40 * time.Hour,
			Result:         30 * time.Hour,
		},

		"request within backend max": {
			BackendDefault: 9 * time.Hour,
			BackendMax:     5 * time.Hour,
			Increment:      4 * time.Hour,
			Result:         4 * time.Hour,
		},

		"request outside backend max": {
			BackendDefault: 9 * time.Hour,
			BackendMax:     4 * time.Hour,
			Increment:      5 * time.Hour,
			Result:         4 * time.Hour,
		},

		"request is negative, no backend default, use sysview": {
			Increment: -7 * time.Hour,
			Result:    5 * time.Hour,
		},

		"lease increment too large": {
			Increment: 40 * time.Hour,
			Result:    30 * time.Hour,
		},
	}

	for name, tc := range cases {
		req := &logical.Request{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					TTL:       1 * time.Hour,
					IssueTime: now,
					Increment: tc.Increment,
				},
			},
		}

		callback := LeaseExtend(tc.BackendDefault, tc.BackendMax, testSysView)
		resp, err := callback(req, nil)
		if (err != nil) != tc.Error {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}
		if tc.Error {
			continue
		}

		// Round it to the nearest hour
		lease := now.Add(resp.Auth.TTL).Round(time.Hour).Sub(now)
		if lease != tc.Result {
			t.Fatalf("bad: %s\nlease: %s", name, lease)
		}
	}
}
