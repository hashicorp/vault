// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package mountsauthtune

import (
	"fmt"
	"strings"
	"testing"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestMountsAuthPathsRequireSudo verifies that sys/mounts/auth/<path> and
// sys/mounts/auth/<path>/tune require sudo, matching sys/auth/<path>, for auth
// methods mounted at single-segment and nested paths. A token that can read and
// update those paths is denied without sudo and allowed with it.
func TestMountsAuthPathsRequireSudo(t *testing.T) {
	t.Parallel()

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": credUserpass.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	rootClient := cluster.Cores[0].Client

	for _, mountPath := range []string{"userpass", "team/userpass", "a/b/c/userpass"} {
		testName := strings.ReplaceAll(mountPath, "/", "-")
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			require.NoError(t, rootClient.Sys().EnableAuthWithOptions(mountPath, &api.EnableAuthOptions{
				Type: "userpass",
			}))
			mountInfoPath := "sys/mounts/auth/" + mountPath
			tunePath := mountInfoPath + "/tune"

			require.NoError(t, rootClient.Sys().PutPolicy(testName+"-no-sudo", fmt.Sprintf(`
path %q {
	capabilities = ["read", "update"]
}
path %q {
	capabilities = ["read"]
}
`, tunePath, mountInfoPath)))
			noSudoToken, err := rootClient.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{testName + "-no-sudo"},
			})
			require.NoError(t, err)
			noSudoClient, err := rootClient.Clone()
			require.NoError(t, err)
			noSudoClient.SetToken(noSudoToken.Auth.ClientToken)

			_, err = noSudoClient.Logical().Write(tunePath, map[string]interface{}{
				"max_lease_ttl": "2h",
			})
			require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())
			_, err = noSudoClient.Logical().Read(tunePath)
			require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())
			_, err = noSudoClient.Logical().Read(mountInfoPath)
			require.ErrorContains(t, err, logical.ErrPermissionDenied.Error())

			require.NoError(t, rootClient.Sys().PutPolicy(testName+"-sudo", fmt.Sprintf(`
path %q {
	capabilities = ["read", "update", "sudo"]
}
path %q {
	capabilities = ["read", "sudo"]
}
`, tunePath, mountInfoPath)))
			sudoToken, err := rootClient.Auth().Token().Create(&api.TokenCreateRequest{
				Policies: []string{testName + "-sudo"},
			})
			require.NoError(t, err)
			sudoClient, err := rootClient.Clone()
			require.NoError(t, err)
			sudoClient.SetToken(sudoToken.Auth.ClientToken)

			_, err = sudoClient.Logical().Write(tunePath, map[string]interface{}{
				"max_lease_ttl": "3h",
			})
			require.NoError(t, err)
			tuneResp, err := sudoClient.Logical().Read(tunePath)
			require.NoError(t, err)
			require.Equal(t, "10800", fmt.Sprint(tuneResp.Data["max_lease_ttl"]))
			mountInfoResp, err := sudoClient.Logical().Read(mountInfoPath)
			require.NoError(t, err)
			require.Equal(t, "userpass", mountInfoResp.Data["type"])
		})
	}
}

// TestMountsTuneSecretsEngineWithoutSudo verifies that tuning a secrets engine
// through sys/mounts/<path>/tune still needs only read and update, so the sudo
// requirement for auth methods doesn't extend to secrets engines.
func TestMountsTuneSecretsEngineWithoutSudo(t *testing.T) {
	t.Parallel()

	cluster := vault.NewTestCluster(t, &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
		NumCores:    1,
	})
	rootClient := cluster.Cores[0].Client

	require.NoError(t, rootClient.Sys().Mount("team/kv", &api.MountInput{Type: "kv"}))
	require.NoError(t, rootClient.Sys().PutPolicy("kv-tune", `
path "sys/mounts/team/kv/tune" {
	capabilities = ["read", "update"]
}
`))
	token, err := rootClient.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"kv-tune"},
	})
	require.NoError(t, err)
	client, err := rootClient.Clone()
	require.NoError(t, err)
	client.SetToken(token.Auth.ClientToken)

	_, err = client.Logical().Write("sys/mounts/team/kv/tune", map[string]interface{}{
		"max_lease_ttl": "3h",
	})
	require.NoError(t, err)
	tuneResp, err := client.Logical().Read("sys/mounts/team/kv/tune")
	require.NoError(t, err)
	require.Equal(t, "10800", fmt.Sprint(tuneResp.Data["max_lease_ttl"]))
}
