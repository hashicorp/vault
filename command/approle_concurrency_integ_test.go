package command

import (
	"context"
	"sync"
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestAppRole_Integ_ConcurrentLogins(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"approle": credAppRole.Factory,
		},
	}

	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.Cores

	vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err = client.Sys().EnableAuthWithOptions("approle", &api.EnableAuthOptions{
		Type: "approle",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("auth/approle/role/role1", map[string]interface{}{
		"bind_secret_id": "true",
		"period":         "300",
	})
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Logical().Write("auth/approle/role/role1/secret-id", nil)
	if err != nil {
		t.Fatal(err)
	}
	secretID := secret.Data["secret_id"].(string)

	secret, err = client.Logical().Read("auth/approle/role/role1/role-id")
	if err != nil {
		t.Fatal(err)
	}
	roleID := secret.Data["role_id"].(string)

	wg := &sync.WaitGroup{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			appRoleAuth, err := auth.NewAppRoleAuth(roleID, &auth.SecretID{FromString: secretID})
			if err != nil {
				t.Error(err)
				return
			}
			secret, err := client.Auth().Login(context.TODO(), appRoleAuth)
			if err != nil {
				t.Error(err)
				return
			}
			if secret.Auth.ClientToken == "" {
				t.Error("expected a successful login")
				return
			}
		}()

	}
	wg.Wait()
}
