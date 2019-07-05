package awsauth

import (
	"context"
	"testing"
	"time"

	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestBackend_E2E_Initialize(t *testing.T) {

	ctx := context.Background()

	// Set up the cluster.  This will trigger an Initialize(); we sleep briefly
	// awaiting its completion.
	cluster := setupAwsTestCluster(t, ctx)
	defer cluster.Cleanup()
	time.Sleep(time.Second)
	core := cluster.Cores[0]

	// Fetch the aws auth's path in storage.  This is a uuid that is different
	// every time we run the test
	authUuids, err := core.UnderlyingStorage.List(ctx, "auth/")
	if err != nil {
		t.Fatal(err)
	}
	if len(authUuids) != 1 {
		t.Fatalf("expected exactly one auth path")
	}
	awsPath := "auth/" + authUuids[0]

	// Make sure that the upgrade happened, by fishing the 'config/version'
	// entry out of storage.  We can't use core.Client.Logical().Read() to do
	// this, because 'config/version' hasn't been exposed as a path.
	// TODO: should we expose 'config/version' as a path?
	version, err := core.UnderlyingStorage.Get(ctx, awsPath+"config/version")
	if err != nil {
		t.Fatal(err)
	}
	if version == nil {
		t.Fatalf("no config found")
	}

	// Nuke the version, so we can pretend that Initialize() has never been run
	if err := core.UnderlyingStorage.Delete(ctx, awsPath+"config/version"); err != nil {
		t.Fatal(err)
	}
	version, err = core.UnderlyingStorage.Get(ctx, awsPath+"config/version")
	if err != nil {
		t.Fatal(err)
	}
	if version != nil {
		t.Fatalf("version found")
	}

	// Create a role
	data := map[string]interface{}{
		"auth_type":       "ec2",
		"policies":        "default",
		"bound_subnet_id": "subnet-abcdef"}
	if _, err := core.Client.Logical().Write("auth/aws/role/test-role", data); err != nil {
		t.Fatal(err)
	}
	role, err := core.Client.Logical().Read("auth/aws/role/test-role")
	if err != nil {
		t.Fatal(err)
	}
	if role == nil {
		t.Fatalf("no role found")
	}

	// There should _still_ be no config version
	version, err = core.UnderlyingStorage.Get(ctx, awsPath+"config/version")
	if err != nil {
		t.Fatal(err)
	}
	if version != nil {
		t.Fatalf("version found")
	}

	// Seal, and then Unseal. This will once again trigger an Initialize(),
	// only this time there will be a role present during the upgrade.
	core.Seal(t)
	cluster.UnsealCores(t)
	time.Sleep(time.Second)

	// Now the config version should be there again
	version, err = core.UnderlyingStorage.Get(ctx, awsPath+"config/version")
	if err != nil {
		t.Fatal(err)
	}
	if version == nil {
		t.Fatalf("no version found")
	}
}

func setupAwsTestCluster(t *testing.T, ctx context.Context) *vault.TestCluster {

	// create a cluster with the aws auth backend built-in
	logger := logging.NewVaultLogger(hclog.Trace)
	coreConfig := &vault.CoreConfig{
		Logger: logger,
		CredentialBackends: map[string]logical.Factory{
			"aws": Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: vaulthttp.Handler,
	})

	cluster.Start()
	if len(cluster.Cores) != 1 {
		t.Fatalf("expected exactly one core")
	}
	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)

	// load the auth plugin
	if err := core.Client.Sys().EnableAuthWithOptions("aws", &api.EnableAuthOptions{
		Type: "aws",
	}); err != nil {
		t.Fatal(err)
	}

	return cluster
}
