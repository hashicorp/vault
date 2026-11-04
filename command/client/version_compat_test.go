// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestIsCompatibleVersion verifies that a vault-client build (ent.client)
// matches enterprise servers that run the same version, so Agent and Proxy do
// not report a version mismatch. Any other difference in version, prerelease,
// or edition must still count as a mismatch.
func TestIsCompatibleVersion(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name   string
		client string
		server string
		want   bool
	}{
		{name: "identical", client: "2.2.2+ent", server: "2.2.2+ent", want: true},
		{name: "different_edition", client: "2.2.2+ent.hsm", server: "2.2.2+ent", want: false},
		{name: "client_ent", client: "2.2.2+ent.client", server: "2.2.2+ent", want: true},
		{name: "client_ent_hsm", client: "2.2.2+ent.client", server: "2.2.2+ent.hsm", want: true},
		{name: "client_ent_fips1403", client: "2.2.2+ent.client", server: "2.2.2+ent.fips1403", want: true},
		{name: "client_ent_hsm_fips1403", client: "2.2.2+ent.client", server: "2.2.2+ent.hsm.fips1403", want: true},
		{name: "client_prerelease", client: "2.3.0-beta1+ent.client", server: "2.3.0-beta1+ent", want: true},
		{name: "client_different_version", client: "2.2.2+ent.client", server: "2.2.3+ent", want: false},
		{name: "client_different_prerelease", client: "2.3.0-beta1+ent.client", server: "2.3.0+ent", want: false},
		{name: "client_ce_server", client: "2.2.2+ent.client", server: "2.2.2", want: false},
		{name: "client_on_server_side", client: "2.2.2+ent", server: "2.2.2+ent.client", want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tc.want, IsCompatibleVersion(tc.client, tc.server),
				"unexpected match result for client %q and server %q", tc.client, tc.server)
		})
	}
}
