package command

import (
	"testing"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/vault"

	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
)

func TestPathMap_Upgrade_API(t *testing.T) {
	var err error
	coreConfig := &vault.CoreConfig{
		DisableMlock: true,
		DisableCache: true,
		Logger:       log.NewNullLogger(),
		CredentialBackends: map[string]logical.Factory{
			"app-id": credAppId.Factory,
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

	// Enable the app-id method
	err = client.Sys().EnableAuthWithOptions("app-id", &api.EnableAuthOptions{
		Type: "app-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create an app-id
	_, err = client.Logical().Write("auth/app-id/map/app-id/test-app-id", map[string]interface{}{
		"policy": "test-policy",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Create a user-id
	_, err = client.Logical().Write("auth/app-id/map/user-id/test-user-id", map[string]interface{}{
		"value": "test-app-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Perform a login. It should succeed.
	_, err = client.Logical().Write("auth/app-id/login", map[string]interface{}{
		"app_id":  "test-app-id",
		"user_id": "test-user-id",
	})
	if err != nil {
		t.Fatal(err)
	}

	// List the hashed app-ids in the storage
	secret, err := client.Logical().List("auth/app-id/map/app-id")
	if err != nil {
		t.Fatal(err)
	}
	hashedAppID := secret.Data["keys"].([]interface{})[0].(string)

	// Try reading it. This used to cause an issue which is fixed in [GH-3806].
	_, err = client.Logical().Read("auth/app-id/map/app-id/" + hashedAppID)
	if err != nil {
		t.Fatal(err)
	}

	// Ensure that there was no issue by performing another login
	_, err = client.Logical().Write("auth/app-id/login", map[string]interface{}{
		"app_id":  "test-app-id",
		"user_id": "test-user-id",
	})
	if err != nil {
		t.Fatal(err)
	}
}
