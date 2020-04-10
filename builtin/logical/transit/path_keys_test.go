package transit_test

import (
	"testing"

	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/builtin/audit/file"
	"github.com/hashicorp/vault/builtin/logical/transit"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/acctest"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestTransit_Issue_2958(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": file.Factory,
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

	err := client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "/dev/null",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foo", map[string]interface{}{
		"type": "ecdsa-p256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foobar", map[string]interface{}{
		"type": "ecdsa-p384",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/bar", map[string]interface{}{
		"type": "ed25519",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foobar")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/bar")
	if err != nil {
		t.Fatal(err)
	}
}

// TODO: remove this POC
func TestTransit_Issue_2958_Docker(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"transit": transit.Factory,
		},
		AuditBackends: map[string]audit.Factory{
			"file": file.Factory,
		},
	}

	cluster, err := acctest.NewDockerCluster(t.Name(), coreConfig, nil)
	if err != nil {
		t.Fatal(err)
	}

	// cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
	// 	HandlerFunc: vaulthttp.Handler,
	// })
	// cluster.Start()
	defer cluster.Cleanup()

	cores := cluster.ClusterNodes

	// vault.TestWaitActive(t, cores[0].Core)

	client := cores[0].Client

	err = client.Sys().EnableAuditWithOptions("file", &api.EnableAuditOptions{
		Type: "file",
		Options: map[string]string{
			"file_path": "/dev/null",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().Mount("transit", &api.MountInput{
		Type: "transit",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foo", map[string]interface{}{
		"type": "ecdsa-p256",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/foobar", map[string]interface{}{
		"type": "ecdsa-p384",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Write("transit/keys/bar", map[string]interface{}{
		"type": "ed25519",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foo")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/foobar")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Logical().Read("transit/keys/bar")
	if err != nil {
		t.Fatal(err)
	}
}
