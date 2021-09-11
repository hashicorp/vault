package misc

import (
	"encoding/json"
	"testing"

	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

func TestKVV2_Patch_FieldFilters(t *testing.T) {
	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": logicalKv.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0]
	vault.TestWaitActive(t, core.Core)
	client := core.Client

	// Enable KVv2
	err := client.Sys().Mount("kv", &api.MountInput{
		Type: "kv-v2",
	})
	if err != nil {
		t.Fatal(err)
	}

	// create a policy with field filters that should emulate the behavior amex wants
	policy := `
path "kv/*" {
	capabilities = ["create", "patch", "list"]
}

path "kv/metadata/*" {
	capabilities = ["read", "list"]
}
`

	_, err = client.Logical().Write("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"bar": "baz", "quux": map[string]interface{}{"wibble": "wobble"}, "wibble": "wobble"}})
	if err != nil {
		t.Fatal(err)
	}

	err = client.Sys().PutPolicy("my-policy", policy)
	if err != nil {
		t.Fatal(err)
	}

	secret, err := client.Auth().Token().Create(&api.TokenCreateRequest{
		Policies: []string{"my-policy"},
	})
	if err != nil {
		t.Fatal(err)
	}

	token := secret.Auth.ClientToken
	client.SetToken(token)

	// run through the Amex use cases
	// 1. creating new versions of an existing keypair works
	secret, err = client.Logical().JSONMergePatch("kv/data/foo", map[string]interface{}{"data": map[string]interface{}{"baz": "quux"}})
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["version"].(json.Number) != json.Number("2") {
		t.Fatalf("expect version to be 2 but got %v", secret.Data["version"])
	}

	// 2. creating new keypairs works
	_, err = client.Logical().Write("kv/data/bar", map[string]interface{}{"data": map[string]interface{}{"foo": "baz"}})
	if err != nil {
		t.Fatal(err)
	}

	// 3. reading these keypairs fails
	_, err = client.Logical().Read("kv/data/foo")
	if err == nil {
		t.Fatal("expected to get an error but got none")
	}

	// 4. listing keys for a mount works, without being able to see the values (not being able to see the value is tested in #3 above)
	secret, err = client.Logical().List("kv/metadata")
	if err != nil {
		t.Fatal(err)
	}
	keys := make([]string, 0)
	for _, v := range secret.Data["keys"].([]interface{}) {
		keys = append(keys, v.(string))
	}
	if !strutil.EquivalentSlices([]string{"bar", "foo"}, keys) {
		t.Fatalf("expected to get a slice with bar and foo in it but got %v instead", secret.Data["keys"])
	}

	// 5. seeing the version count works
	secret, err = client.Logical().Read("kv/metadata/foo")
	if err != nil {
		t.Fatal(err)
	}
	if secret.Data["current_version"].(json.Number) != json.Number("2") {
		t.Fatalf("expected 2 versions but got %v", secret.Data["current_version"])
	}

	// 6. reading sub-keys from a top level key works
	// first update the policy to include field filters
	newPolicy := `
path "kv/*" {
	capabilities = ["create", "patch", "list"]

	field_filters = [
		{
			filter_on = ["read"]
			fields = ["/foo/bar", "/foo/quux/wibble"]
		}
	]
}

path "kv/metadata/*" {
	capabilities = ["read", "list"]
}
`

	err = client.Sys().PutPolicy("my-policy", newPolicy)
	if err != nil {
		t.Fatal(err)
	}

	// reading /foo/bar should work
	// reading /foo/quux/wibble should work
	// reading /foo/wibble should not work
}
