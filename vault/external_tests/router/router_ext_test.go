package router

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/builtin/credential/userpass"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"
)

func TestRouter_MountSubpath_Checks(t *testing.T) {
	testRouter_MountSubpath(t, []string{"auth/abcd/123", "abcd/123"})
	testRouter_MountSubpath(t, []string{"abcd/123", "auth/abcd/123"})
	testRouter_MountSubpath(t, []string{"auth/abcd/123", "abcd/123"})
}

func testRouter_MountSubpath(t *testing.T, mountPoints []string) {
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"userpass": userpass.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	authInput := &api.EnableAuthOptions{
		Type: "userpass",
	}

	for _, mp := range mountPoints {
		t.Logf("mounting %s", mp)
		var err error
		err = client.Sys().EnableAuthWithOptions(mp, authInput)
		if err != nil {
			t.Fatalf("err: %v", err)
		}
	}

	cluster.EnsureCoresSealed(t)

	cluster.UnsealCores(t)

	t.Logf("Done: %#v", mountPoints)
}
