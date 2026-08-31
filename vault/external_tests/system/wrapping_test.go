// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package system

import (
	"strings"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestWrappingWrap_JWTFormatBlocked verifies that sys/wrapping/wrap ignores
// X-Vault-Wrap-Format: jwt and returns a standard Vault service wrapping token,
// preventing callers from obtaining a signed JWT wrapping token through this endpoint.
func TestWrappingWrap_JWTFormatBlocked(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client

	client.SetWrappingLookupFunc(func(operation, path string) string {
		return "300s"
	})
	client.AddHeader("X-Vault-Wrap-Format", "jwt")

	resp, err := client.Logical().Write("sys/wrapping/wrap", map[string]interface{}{
		"foo": "bar",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.NotNil(t, resp.WrapInfo)
	require.True(t, strings.HasPrefix(resp.WrapInfo.Token, "hvs."),
		"expected a Vault service token (hvs.) but got: %q", resp.WrapInfo.Token)
}
