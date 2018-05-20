package api_test

import (
	"strings"
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
	client.SetMaxRetries(0)

	// This first block verifies that if we are using recovery keys to force a
	// rekey of a stored-shares barrier that verification is not allowed since
	// the keys aren't returned
	vault.DefaultSealPretendsToAllowRecoveryKeys = true
	vault.DefaultSealPretendsToAllowStoredShares = true
	vault.DefaultSealPretendRecoveryConfig = &vault.SealConfig{}
	status, err := client.Sys().RekeyInit(&api.RekeyInitRequest{
		StoredShares:        1,
		RequireVerification: true,
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "requiring verification not supported") {
		t.Fatalf("unexpected error: %v", err)
	}

	// Now we set things back and start a normal rekey with the verification process
	vault.DefaultSealPretendsToAllowRecoveryKeys = false
	vault.DefaultSealPretendsToAllowStoredShares = false
	vault.DefaultSealPretendRecoveryConfig = nil
	status, err = client.Sys().RekeyInit(&api.RekeyInitRequest{
		SecretShares:        5,
		SecretThreshold:     3,
		RequireVerification: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if status == nil {
		t.Fatal("nil status")
	}
	/*
		cluster.EnsureCoresSealed(t)
		cluster.UnsealCores(t)
	*/
}
