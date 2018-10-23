package consul

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/dockertest"
)

func prepareTestContainerNewACL(t *testing.T) (cleanup func(), retAddress string, consulToken string) {
	consulToken = os.Getenv("CONSUL_HTTP_TOKEN")

	retAddress = os.Getenv("CONSUL_HTTP_ADDR")

	if retAddress != "" {
		return func() {}, retAddress, consulToken
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "ncorrare/consul",
		Tag:        "v1.4.0-rc1",
		Cmd:        []string{"agent", "-dev", "-client", "0.0.0.0", "-hcl", `acl { enabled = true default_policy = "deny" }`},
		Env:        []string{"CONSUL_BIND_INTERFACE=eth0"},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local Consul docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retAddress = fmt.Sprintf("localhost:%s", resource.GetPort("8500/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		consulConfig := consulapi.DefaultNonPooledConfig()
		consulConfig.Address = retAddress
		consul, err := consulapi.NewClient(consulConfig)
		if err != nil {
			return err
		}
		aclbootstrap, _, err := consul.ACL().Bootstrap()
		if err != nil {
			return err
		}
		consulToken = aclbootstrap.SecretID
		t.Logf("[WARN] Generated Master token: %s", consulToken)
		policy := &consulapi.ACLPolicy{
			Name:        "test",
			Description: "test",
			Rules: `node_prefix "" {
                policy = "write"
              }

              service_prefix "" {
                policy = "read"
              }
      `,
		}
		q := &consulapi.WriteOptions{
			Token: consulToken,
		}
		_, _, err = consul.ACL().PolicyCreate(policy, q)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retAddress, consulToken
}

func TestBackend_old_config_access(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainerNewACL(t)
	defer cleanup()

	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	confReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Storage:   config.StorageView,
		Data:      connData,
	}

	resp, err := b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) || resp != nil {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	confReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), confReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
	}

	expected := map[string]interface{}{
		"address": connData["address"].(string),
		"scheme":  "http",
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_old_renew_revoke(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainerNewACL(t)
	defer cleanup()
	connData := map[string]interface{}{
		"address": connURL,
		"token":   connToken,
	}

	req := &logical.Request{
		Storage:   config.StorageView,
		Operation: logical.UpdateOperation,
		Path:      "config/access",
		Data:      connData,
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Path = "roles/test"
	req.Data = map[string]interface{}{
		"policies": []string{"test"},
		"lease":    "6h",
	}
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.ReadOperation
	req.Path = "creds/test"
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("resp nil")
	}
	if resp.IsError() {
		t.Fatalf("resp is error: %v", resp.Error())
	}

	generatedSecret := resp.Secret
	generatedSecret.TTL = 6 * time.Hour

	var d struct {
		Token    string `mapstructure:"token"`
		Accessor string `mapstructure:"accessor"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("[WARN] Generated token: %s with accessor %s", d.Token, d.Accessor)

	// Build a client and verify that the credentials work
	consulapiConfig := consulapi.DefaultNonPooledConfig()
	consulapiConfig.Address = connData["address"].(string)
	consulapiConfig.Token = d.Token
	client, err := consulapi.NewClient(consulapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[WARN] Verifying that the generated token works...")
	_, err = client.Catalog(), nil
	if err != nil {
		t.Fatal(err)
	}

	req.Operation = logical.RenewOperation
	req.Secret = generatedSecret
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp == nil {
		t.Fatal("got nil response from renew")
	}

	req.Operation = logical.RevokeOperation
	resp, err = b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	// Build a management client and verify that the token does not exist anymore
	consulmgmtConfig := consulapi.DefaultNonPooledConfig()
	consulmgmtConfig.Address = connData["address"].(string)
	consulmgmtConfig.Token = connData["token"].(string)
	mgmtclient, err := consulapi.NewClient(consulmgmtConfig)

	q := &consulapi.QueryOptions{
		Datacenter: "DC1",
	}

	t.Log("[WARN] Verifying that the generated token does not exist...")
	_, _, err = mgmtclient.ACL().TokenRead(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}
