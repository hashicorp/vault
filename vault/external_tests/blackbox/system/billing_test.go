// Copyright IBM Corp. 2025, 2026
// SPDX-License-Identifier: BUSL-1.1

package system

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/sdk/helper/testcluster/blackbox"
	"github.com/stretchr/testify/require"
)

// TestBillingOverviewNamespaceRestrictions verifies that sys/billing/overview
// returns appropriate errors when called from different namespace levels in HVD.
// In HVD, tests run in admin/bbsdk-xxxxx, and this test verifies:
// - Calling from base namespace (admin) returns "unsupported path"
// - Calling from root namespace (empty) returns "permission denied"
func TestBillingOverviewNamespaceRestrictions(t *testing.T) {
	v := blackbox.New(t)

	// Verify cluster stability first
	v.AssertClusterHealthy()

	// Check if we're in HVD (has base namespace from VAULT_NAMESPACE)
	baseNS := v.GetParentNamespace()
	if baseNS == "" {
		t.Skip("Skipping namespace restriction tests - no base namespace configured (not in HVD)")
	}

	testCases := []struct {
		name              string
		namespaceSwitcher func(func() (*api.Secret, error)) (*api.Secret, error)
		expectedError     string
	}{
		{
			name:              "base_namespace_unsupported",
			namespaceSwitcher: v.WithParentNamespace,
			expectedError:     "unsupported path",
		},
		{
			name:              "root_namespace_permission_denied",
			namespaceSwitcher: v.WithRootNamespace,
			expectedError:     "permission denied",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var rawResp *api.Response
			var err error
			_, err = tc.namespaceSwitcher(func() (*api.Secret, error) {
				var readErr error
				rawResp, readErr = v.Client.Logical().ReadRawWithContext(context.Background(), "sys/billing/overview")
				if readErr != nil {
					return nil, readErr
				}
				// Parse the raw response to get the error
				return v.Client.Logical().ParseRawResponseAndCloseBody(rawResp, nil)
			})
			require.ErrorContains(t, err, tc.expectedError)
		})
	}
}
