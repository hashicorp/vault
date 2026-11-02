// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package clientlogin

import (
	"bytes"
	"testing"

	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/command/token"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// noHCPTokenHelper reports no HCP token so the login doesn't depend on any
// HCP credentials cached on the machine running the test.
type noHCPTokenHelper struct{}

func (noHCPTokenHelper) GetHCPToken(string) (*hcpvlib.HCPToken, error) {
	return nil, nil
}

// TestClientLogin_Userpass runs `login -method=userpass` through the
// vault-client entry point against a real server. vault-client once
// registered `login` without any auth method handlers, so this failed with
// "Unknown auth method: userpass" even though the full vault binary worked.
func TestClientLogin_Userpass(t *testing.T) {
	t.Parallel()

	cluster := minimal.NewTestSoloCluster(t, &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	})
	client := cluster.Cores[0].Client

	require.NoError(t, client.Sys().EnableAuthWithOptions("userpass", &api.EnableAuthOptions{
		Type: "userpass",
	}))
	_, err := client.Logical().Write("auth/userpass/users/vc-test-user", map[string]any{
		"password":       "test-password",
		"token_policies": "default",
	})
	require.NoError(t, err)

	loginClient, err := client.Clone()
	require.NoError(t, err)
	loginClient.ClearToken()

	tokenHelper := token.NewTestingTokenHelper()
	var stdout, stderr bytes.Buffer
	code := clientcmd.RunCustom([]string{
		"login",
		"-method=userpass",
		"username=vc-test-user",
		"password=test-password",
	}, &base.RunOptions{
		TokenHelper:    tokenHelper,
		HCPTokenHelper: noHCPTokenHelper{},
		Client:         loginClient,
		Stdout:         &stdout,
		Stderr:         &stderr,
	})
	require.Equal(t, 0, code, "stdout: %s\nstderr: %s", stdout.String(), stderr.String())

	storedToken, err := tokenHelper.Get()
	require.NoError(t, err)
	require.NotEmpty(t, storedToken)

	client, err = client.Clone()
	require.NoError(t, err)
	client.SetToken(storedToken)
	secret, err := client.Auth().Token().LookupSelf()
	require.NoError(t, err)
	require.Equal(t, "auth/userpass/login/vc-test-user", secret.Data["path"])
}
