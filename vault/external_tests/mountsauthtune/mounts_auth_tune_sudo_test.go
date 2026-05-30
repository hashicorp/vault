// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mountsauthtune

import (
	"fmt"
	"testing"

	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestMountsAuthTuneRequiresSudo ensures sys/mounts/auth/<path>/tune requires
// sudo-equivalent privilege while sys/mounts/auth/<path> remains readable
// without sudo when ACL allows it.
func TestMountsAuthTuneRequiresSudo(t *testing.T) {
	t.Parallel()

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})

	client := cluster.Cores[0].Client
	rootToken := cluster.RootToken

	require.NoError(t, client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	}))

	noSudoPolicy := `
path "sys/mounts/auth/userpass/tune" {
	capabilities = ["read", "update"]
}
path "sys/mounts/auth/userpass*" {
	capabilities = ["read"]
}
`
	require.NoError(t, client.Sys().PutPolicy("mounts-auth-tune-no-sudo", noSudoPolicy))

	noSudoTokenResp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"mounts-auth-tune-no-sudo"},
	})
	require.NoError(t, err)
	require.NotNil(t, noSudoTokenResp)
	require.NotNil(t, noSudoTokenResp.Auth)
	require.NotEmpty(t, noSudoTokenResp.Auth.ClientToken)

	client.SetToken(noSudoTokenResp.Auth.ClientToken)

	// Non-tune path should remain readable without sudo.
	nonTuneResp, err := client.Logical().Read("sys/mounts/auth/userpass")
	require.NoError(t, err)
	require.NotNil(t, nonTuneResp)

	// Tune endpoints should fail without sudo.
	_, err = client.Logical().Write("sys/mounts/auth/userpass/tune", map[string]interface{}{
		"max_lease_ttl": "2h",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), logical.ErrPermissionDenied.Error())

	_, err = client.Logical().Read("sys/mounts/auth/userpass/tune")
	require.Error(t, err)
	require.Contains(t, err.Error(), logical.ErrPermissionDenied.Error())

	client.SetToken(rootToken)

	sudoPolicy := `
path "sys/mounts/auth/userpass/tune" {
	capabilities = ["sudo", "read", "update"]
}
`
	require.NoError(t, client.Sys().PutPolicy("mounts-auth-tune-with-sudo", sudoPolicy))

	withSudoTokenResp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"mounts-auth-tune-with-sudo"},
	})
	require.NoError(t, err)
	require.NotNil(t, withSudoTokenResp)
	require.NotNil(t, withSudoTokenResp.Auth)
	require.NotEmpty(t, withSudoTokenResp.Auth.ClientToken)

	client.SetToken(withSudoTokenResp.Auth.ClientToken)

	_, err = client.Logical().Write("sys/mounts/auth/userpass/tune", map[string]interface{}{
		"max_lease_ttl": "3h",
	})
	require.NoError(t, err)

	tuneResp, err := client.Logical().Read("sys/mounts/auth/userpass/tune")
	require.NoError(t, err)
	require.NotNil(t, tuneResp)
	require.Contains(t, tuneResp.Data, "max_lease_ttl")
	require.Equal(t, "10800", fmt.Sprint(tuneResp.Data["max_lease_ttl"]))
}
