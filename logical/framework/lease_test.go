package framework

import (
	"testing"
	"time"

	"github.com/hashicorp/vault/logical"
)

func TestLeaseExtend(t *testing.T) {
	now := time.Now().UTC().Round(time.Hour)

	cases := map[string]struct {
		Max     time.Duration
		Request time.Duration
		Result  time.Duration
	}{
		"valid request, good bounds": {
			Max:     30 * time.Hour,
			Request: 1 * time.Hour,
			Result:  1 * time.Hour,
		},

		"request is too long": {
			Max:     3 * time.Hour,
			Request: 7 * time.Hour,
			Result:  3 * time.Hour,
		},
	}

	for name, tc := range cases {
		req := &logical.Request{
			Auth: &logical.Auth{
				LeaseOptions: logical.LeaseOptions{
					Lease:          1 * time.Second,
					LeaseIssue:     now,
					LeaseIncrement: tc.Request,
				},
			},
		}

		callback := LeaseExtend(tc.Max)
		resp, err := callback(req, nil)
		if err != nil {
			t.Fatalf("bad: %s\nerr: %s", name, err)
		}

		// Round it to the nearest hour
		lease := now.Add(resp.Auth.Lease).Round(time.Hour).Sub(now)
		if lease != tc.Result {
			t.Fatalf("bad: %s\nlease: %s", name, lease)
		}
	}
}
