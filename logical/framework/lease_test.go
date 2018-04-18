package framework

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestCalculateTTL(t *testing.T) {
	testSysView := logical.StaticSystemView{
		DefaultLeaseTTLVal: 5 * time.Hour,
		MaxLeaseTTLVal:     30 * time.Hour,
	}

	cases := map[string]struct {
		Increment      time.Duration
		BackendDefault time.Duration
		BackendMax     time.Duration
		Period         time.Duration
		ExplicitMaxTTL time.Duration
		Result         time.Duration
		Warnings       int
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
			Warnings:       1,
		},

		"all request values are larger than the system view, so the system view limits": {
			BackendDefault: 40 * time.Hour,
			BackendMax:     50 * time.Hour,
			Increment:      40 * time.Hour,
			Result:         30 * time.Hour,
			Warnings:       1,
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
			Warnings:       1,
		},

		"request is negative, no backend default, use sysview": {
			Increment: -7 * time.Hour,
			Result:    5 * time.Hour,
		},

		"lease increment too large": {
			Increment: 40 * time.Hour,
			Result:    30 * time.Hour,
			Warnings:  1,
		},

		"periodic, good request, period is preferred": {
			Increment:      3 * time.Hour,
			BackendDefault: 4 * time.Hour,
			BackendMax:     2 * time.Hour,
			Period:         1 * time.Hour,
			Result:         1 * time.Hour,
		},

		"period too large, explicit max ttl is preferred": {
			Period:         2 * time.Hour,
			ExplicitMaxTTL: 1 * time.Hour,
			Result:         1 * time.Hour,
			Warnings:       1,
		},

		"period too large, capped by backend max": {
			Period:     2 * time.Hour,
			BackendMax: 1 * time.Hour,
			Result:     1 * time.Hour,
			Warnings:   1,
		},
	}

	for name, tc := range cases {
		ttl, warnings, err := CalculateTTL(testSysView, tc.Increment, tc.BackendDefault, tc.Period, tc.BackendMax, tc.ExplicitMaxTTL, time.Time{})
		if (err != nil) != tc.Error {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}
		if tc.Error {
			continue
		}

		// Round it to the nearest hour
		now := time.Now().Round(time.Hour)
		lease := now.Add(ttl).Round(time.Hour).Sub(now)
		if lease != tc.Result {
			t.Fatalf("bad: %s\nlease: %s", name, lease)
		}

		if tc.Warnings != len(warnings) {
			t.Fatalf("bad: %s\nwarning count mismatch, expect %d, got %d: %#v", name, tc.Warnings, len(warnings), warnings)
		}
	}
}
