// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

// TestGetMount tests that we can get a single secret mount
func TestGetMount(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		mountName   string
		mountInput  *api.MountInput
		expected    *api.MountOutput
		shouldMount bool
		expectErr   bool
	}{
		{
			name:       "get-default-mount-success",
			mountName:  "secret",
			mountInput: nil,
			expected: &api.MountOutput{
				Type: "kv",
			},
			shouldMount: false,
			expectErr:   false,
		},
		{
			name:      "get-manual-mount-success",
			mountName: "pki",
			mountInput: &api.MountInput{
				Type: "pki",
			},
			expected: &api.MountOutput{
				Type: "pki",
			},
			shouldMount: true,
			expectErr:   false,
		},
		{
			name:       "error-not-found",
			mountName:  "not-found",
			mountInput: nil,
			expected: &api.MountOutput{
				Type: "not-found",
			},
			shouldMount: false,
			expectErr:   true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			client, closer := testVaultServer(t)
			defer closer()

			if tc.shouldMount {
				err := client.Sys().Mount(tc.mountName, tc.mountInput)
				if err != nil {
					t.Fatal(err)
				}
			}

			mount, err := client.Sys().GetMount(tc.mountName)
			if !tc.expectErr && err != nil {
				t.Fatal(err)
			}

			if !tc.expectErr {
				if tc.expected.Type != mount.Type || tc.expected.PluginVersion != mount.PluginVersion {
					t.Errorf("mount did not match: expected %+v but got %+v", tc.expected, mount)
				}
			} else {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
			}
		})
	}
}
