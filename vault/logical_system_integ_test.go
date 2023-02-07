package vault_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/physical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/hashicorp/vault/version"
)

func TestSystemBackend_InternalUIResultantACL(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		NumCores:                1,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	resp, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"default"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Auth == nil {
		t.Fatal("nil auth")
	}
	if resp.Auth.ClientToken == "" {
		t.Fatal("empty client token")
	}

	client.SetToken(resp.Auth.ClientToken)

	resp, err = client.Logical().Read("sys/internal/ui/resultant-acl")
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("nil response")
	}
	if resp.Data == nil {
		t.Fatal("nil data")
	}

	exp := map[string]interface{}{
		"exact_paths": map[string]interface{}{
			"auth/token/lookup-self": map[string]interface{}{
				"capabilities": []interface{}{
					"read",
				},
			},
			"auth/token/renew-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"auth/token/revoke-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/capabilities-self": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/control-group/request": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/internal/ui/resultant-acl": map[string]interface{}{
				"capabilities": []interface{}{
					"read",
				},
			},
			"sys/leases/lookup": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/leases/renew": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/renew": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/tools/hash": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/lookup": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/unwrap": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
			"sys/wrapping/wrap": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
		},
		"glob_paths": map[string]interface{}{
			"cubbyhole/": map[string]interface{}{
				"capabilities": []interface{}{
					"create",
					"delete",
					"list",
					"read",
					"update",
				},
			},
			"sys/tools/hash/": map[string]interface{}{
				"capabilities": []interface{}{
					"update",
				},
			},
		},
		"root": false,
	}

	if diff := deep.Equal(resp.Data, exp); diff != nil {
		t.Fatal(diff)
	}
}

func TestSystemBackend_HAStatus(t *testing.T) {
	logger := logging.NewVaultLogger(hclog.Trace)
	inm, err := inmem.NewTransactionalInmem(nil, logger)
	if err != nil {
		t.Fatal(err)
	}
	inmha, err := inmem.NewInmemHA(nil, logger)
	if err != nil {
		t.Fatal(err)
	}

	conf := &vault.CoreConfig{
		Physical:   inm,
		HAPhysical: inmha.(physical.HABackend),
	}
	opts := &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	}
	cluster := vault.NewTestCluster(t, conf, opts)
	cluster.Start()
	defer cluster.Cleanup()

	corehelpers.RetryUntil(t, 15*time.Second, func() error {
		// Use standby deliberately to make sure it forwards
		client := cluster.Cores[1].Client
		resp, err := client.Sys().HAStatus()
		if err != nil {
			t.Fatal(err)
		}

		if len(resp.Nodes) != len(cluster.Cores) {
			return fmt.Errorf("expected %d nodes, got %d", len(cluster.Cores), len(resp.Nodes))
		}
		return nil
	})
}

// TestSystemBackend_VersionHistory_unauthenticated tests the sys/version-history
// endpoint without providing a token. Requests to the endpoint must be
// authenticated and thus a 403 response is expected.
func TestSystemBackend_VersionHistory_unauthenticated(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		NumCores:                1,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	client.SetToken("")
	resp, err := client.Logical().List("sys/version-history")

	if resp != nil {
		t.Fatalf("expected nil response, resp: %#v", resp)
	}

	respErr, ok := err.(*api.ResponseError)
	if !ok {
		t.Fatalf("unexpected error type: err: %#v", err)
	}

	if respErr.StatusCode != 403 {
		t.Fatalf("expected response status to be 403, actual: %d", respErr.StatusCode)
	}
}

// TestSystemBackend_VersionHistory_authenticated tests the sys/version-history
// endpoint with authentication. Without synthetically altering the underlying
// core/versions storage entries, a single version entry should exist.
func TestSystemBackend_VersionHistory_authenticated(t *testing.T) {
	t.Parallel()
	cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
		HandlerFunc:             vaulthttp.Handler,
		NumCores:                1,
		RequestResponseCallback: schema.ResponseValidatingCallback(t),
	})
	cluster.Start()
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client

	resp, err := client.Logical().List("sys/version-history")
	if err != nil || resp == nil {
		t.Fatalf("request failed, err: %v, resp: %#v", err, resp)
	}

	var ok bool
	var keys []interface{}
	var keyInfo map[string]interface{}

	if keys, ok = resp.Data["keys"].([]interface{}); !ok {
		t.Fatalf("expected keys to be array, actual: %#v", resp.Data["keys"])
	}

	if keyInfo, ok = resp.Data["key_info"].(map[string]interface{}); !ok {
		t.Fatalf("expected key_info to be map, actual: %#v", resp.Data["key_info"])
	}

	if len(keys) != 1 {
		t.Fatalf("expected single version history entry for %q", version.Version)
	}

	if keyInfo[version.Version] == nil {
		t.Fatalf("expected version %s to be present in key_info, actual: %#v", version.Version, keyInfo)
	}
}
