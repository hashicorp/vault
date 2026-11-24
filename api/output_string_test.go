// Copyright IBM Corp. 2026
// SPDX-License-Identifier: MPL-2.0

package api_test

import (
	"net/http"
	"testing"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// TestOutputStringError_CurlString verifies that the token header reads the
// token with BinaryName, falls back to vault when it's empty, picks up a name
// set after Error() has cached a curl string, and never prints the token.
func TestOutputStringError_CurlString(t *testing.T) {
	t.Parallel()

	request, err := retryablehttp.NewRequest(http.MethodGet, "https://example.test/v1/sys/health", nil)
	require.NoError(t, err)
	request.Header.Set(api.AuthHeaderName, "s.secret")
	outputErr := &api.OutputStringError{Request: request}

	_ = outputErr.Error() // Cache the default curl string
	cs, err := outputErr.CurlString()
	require.NoError(t, err)
	require.Contains(t, cs, "$(vault print token)")
	require.NotContains(t, cs, "s.secret")

	outputErr.BinaryName = "vault-client"
	cs, err = outputErr.CurlString()
	require.NoError(t, err)
	require.Contains(t, cs, "$(vault-client print token)")
	require.NotContains(t, cs, "$(vault print token)")
	require.NotContains(t, cs, "s.secret")
}
