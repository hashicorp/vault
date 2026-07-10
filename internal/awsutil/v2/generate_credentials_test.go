// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package awsutil

import (
	"context"
	"testing"
)

// TestSTSLoginEndpoint verifies that STSLoginEndpoint resolves the correct
// regional STS endpoint URL for a given region, including the global and
// China partition cases.
func TestSTSLoginEndpoint(t *testing.T) {
	cases := map[string]string{
		"us-east-1":  "https://sts.amazonaws.com/",
		"us-west-2":  "https://sts.us-west-2.amazonaws.com/",
		"eu-west-1":  "https://sts.eu-west-1.amazonaws.com/",
		"cn-north-1": "https://sts.cn-north-1.amazonaws.com.cn/",
	}
	for region, want := range cases {
		t.Run(region, func(t *testing.T) {
			got, err := STSLoginEndpoint(context.Background(), region)
			if err != nil {
				t.Fatalf("unexpected error for region %q: %v", region, err)
			}
			if got != want {
				t.Fatalf("region %q: got %q, want %q", region, got, want)
			}
		})
	}
}
