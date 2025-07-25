// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package mounts

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

func TestMountTuneRemoveHeaders(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	c := cluster.Cores[0].Client

	// Mount a PKI backend with certain allowed response headers
	input := &api.MountInput{
		Type:        "pki",
		Description: "my rad pki mount",
		Config: api.MountConfigInput{
			DefaultLeaseTTL: "10m",
			MaxLeaseTTL:     "10m",
			AllowedResponseHeaders: &[]string{
				"Content-Transfer-Encoding",
				"Content-Length",
				"WWW-Authenticate",
			},
		},
		Local:                 false,
		SealWrap:              true,
		ExternalEntropyAccess: false,
	}

	err := c.Sys().Mount("lol", input)
	require.NoError(t, err)

	validateHeaders := func(t *testing.T, m *api.MountOutput, e error) {
		require.NoError(t, err)
		require.NotNil(t, m)
		require.NotNil(t, m.Config.AllowedResponseHeaders)
		require.Equal(t, len(m.Config.AllowedResponseHeaders), 3)
		headers := m.Config.AllowedResponseHeaders
		require.Equal(t, headers[0], "Content-Transfer-Encoding")
		require.Equal(t, headers[1], "Content-Length")
		require.Equal(t, headers[2], "WWW-Authenticate")
	}

	// Confirm the allowed response headers are present
	mount, err := c.Sys().GetMount("lol")
	validateHeaders(t, mount, err)

	// Tune the mount and update unrelated fields, leaving the allowed response headers untouched
	tuneInput := api.MountConfigInput{
		MaxLeaseTTL: "20m",
	}

	err = c.Sys().TuneMount("lol", tuneInput)
	require.NoError(t, err)

	// Verify the original allowed response headers are still present
	mount, err = c.Sys().GetMount("lol")
	validateHeaders(t, mount, err)

	// Tune the mount and remove those headers by explicitly setting them to an empty slice
	tuneInput = api.MountConfigInput{
		AllowedResponseHeaders: &[]string{},
	}

	err = c.Sys().TuneMount("lol", tuneInput)
	require.NoError(t, err)

	// Confirm the allowed response headers are now empty
	mount, err = c.Sys().GetMount("lol")
	require.NoError(t, err)
	require.NotNil(t, mount)
	require.Equal(t, len(mount.Config.AllowedResponseHeaders), 0)
}
