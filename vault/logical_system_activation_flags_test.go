// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/vault/helper/activationflags"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

// TestActivationFlags_Read tests the read operation for the activation flags.
func TestActivationFlags_Read(t *testing.T) {
	t.Run("given an initial state then read flags and expect all to be unactivated", func(t *testing.T) {
		core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

		resp, err := core.systemBackend.HandleRequest(
			context.Background(),
			&logical.Request{
				Operation: logical.ReadOperation,
				Path:      prefixActivationFlags,
				Storage:   core.systemBarrierView,
			},
		)

		require.NoError(t, err)
		require.Equal(t, resp.Data, map[string]interface{}{
			"activated": []string{},
		})
	})
}

// TestActivationFlags_BadFeatureName tests a nonexistent feature name or a missing feature name
// in the activation-flags path API call.
func TestActivationFlags_BadFeatureName(t *testing.T) {
	core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

	tests := map[string]struct {
		featureName string
	}{
		"if no feature name is provided then expect unsupported path": {
			featureName: "",
		},
		"if an invalid feature name is provided then expect unsupported path": {
			featureName: "fake-feature",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := core.router.Route(
				namespace.ContextWithNamespace(context.Background(), namespace.RootNamespace),
				&logical.Request{
					Operation: logical.UpdateOperation,
					Path:      fmt.Sprintf("sys/%s/%s/%s", prefixActivationFlags, tt.featureName, verbActivationFlagsActivate),
					Storage:   core.systemBarrierView,
				},
			)

			require.Error(t, err)
			require.Nil(t, resp)
			require.Equal(t, err, logical.ErrUnsupportedPath)
		})
	}
}

// TestActivationFlags_Write tests the write operations for the activation flags
func TestActivationFlags_Write(t *testing.T) {
	t.Run("given an initial state then write an activation test flag and expect no errors", func(t *testing.T) {
		core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

		_, err := core.systemBackend.HandleRequest(
			context.Background(),
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("%s/%s/%s", prefixActivationFlags, activationFlagTest, verbActivationFlagsActivate),
				Storage:   core.systemBarrierView,
			},
		)

		require.NoError(t, err)
	})

	t.Run("activate identity cleanup flag", func(t *testing.T) {
		core, _, _ := TestCoreUnsealedWithConfig(t, &CoreConfig{})

		resp, err := core.systemBackend.HandleRequest(
			context.Background(),
			&logical.Request{
				Operation: logical.UpdateOperation,
				Path:      fmt.Sprintf("%s/%s/%s", prefixActivationFlags, activationflags.IdentityDeduplication, verbActivationFlagsActivate),
				Storage:   core.systemBarrierView,
			},
		)

		require.NoError(t, err)
		require.NotNil(t, resp)
		require.NotEmpty(t, resp.Data)
		require.NotNil(t, resp.Data["activated"])
		require.Contains(t, resp.Data["activated"], activationflags.IdentityDeduplication)
	})
}
