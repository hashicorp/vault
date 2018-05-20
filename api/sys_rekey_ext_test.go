package api_test

import (
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/vault"
)

func TestSysRekey_Verification(t *testing.T) {
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	vault.TestWaitActive(t, cluster.Cores[0].Core)
	client := cluster.Cores[0].Client

	vault.DefaultSealPretendsToAllowRecoveryKeys = true
	vault.DefaultSealPretendsToAllowStoredShares = true
	status, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		StoredShares:        1,
		RequireVerification: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if status == nil {
		t.Fatal("empty status")
	}

	/*
		cluster.EnsureCoresSealed(t)
		cluster.UnsealCores(t)
	*/
}
