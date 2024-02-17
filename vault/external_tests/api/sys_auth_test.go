// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package api

import (
	"testing"

	"github.com/hashicorp/vault/api"
)

func TestGetAuth(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		mountName   string
		authInput   *api.EnableAuthOptions
		expected    *api.AuthMount
		shouldMount bool
		expectErr   bool
	}{
		{
			name:      "get-default-auth-mount-success",
			mountName: "token",
			authInput: nil,
			expected: &api.AuthMount{
				Type: "token",
			},
			shouldMount: false,
			expectErr:   false,
		},
		{
			name:      "get-manual-auth-mount-success",
			mountName: "userpass",
			authInput: &api.EnableAuthOptions{
				Type: "userpass",
			},
			expected: &api.AuthMount{
				Type: "userpass",
			},
			shouldMount: true,
			expectErr:   false,
		},
		{
			name:      "error-not-found",
			mountName: "not-found",
			authInput: nil,
			expected: &api.AuthMount{
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
				err := client.Sys().EnableAuthWithOptions(tc.mountName, tc.authInput)
				if err != nil {
					t.Fatal(err)
				}
			}

			mount, err := client.Sys().GetAuth(tc.mountName)
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
