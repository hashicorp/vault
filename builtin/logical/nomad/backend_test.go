package nomad

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	nomadapi "github.com/hashicorp/nomad/api"
	"github.com/hashicorp/vault/helper/testhelpers"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/ory/dockertest"
)

func prepareTestContainer(t *testing.T) (cleanup func(), retAddress string, nomadToken string) {
	nomadToken = os.Getenv("NOMAD_TOKEN")

	retAddress = os.Getenv("NOMAD_ADDR")

	if retAddress != "" {
		return func() {}, retAddress, nomadToken
	}

	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to connect to docker: %s", err)
	}

	dockerOptions := &dockertest.RunOptions{
		Repository: "catsby/nomad",
		Tag:        "0.8.4",
		Cmd:        []string{"agent", "-dev"},
		Env:        []string{`NOMAD_LOCAL_CONFIG=bind_addr = "0.0.0.0" acl { enabled = true }`},
	}
	resource, err := pool.RunWithOptions(dockerOptions)
	if err != nil {
		t.Fatalf("Could not start local Nomad docker container: %s", err)
	}

	cleanup = func() {
		err := pool.Purge(resource)
		if err != nil {
			t.Fatalf("Failed to cleanup local container: %s", err)
		}
	}

	retAddress = fmt.Sprintf("http://localhost:%s/", resource.GetPort("4646/tcp"))

	// exponential backoff-retry
	if err = pool.Retry(func() error {
		var err error
		nomadapiConfig := nomadapi.DefaultConfig()
		nomadapiConfig.Address = retAddress
		nomad, err := nomadapi.NewClient(nomadapiConfig)
		if err != nil {
			return err
		}
		aclbootstrap, _, err := nomad.ACLTokens().Bootstrap(nil)
		if err != nil {
			return err
		}
		nomadToken = aclbootstrap.SecretID
		t.Logf("[WARN] Generated Master token: %s", nomadToken)
		policy := &nomadapi.ACLPolicy{
			Name:        "test",
			Description: "test",
			Rules: `namespace "default" {
        policy = "read"
      }
      `,
		}
		anonPolicy := &nomadapi.ACLPolicy{
			Name:        "anonymous",
			Description: "Deny all access for anonymous requests",
			Rules: `namespace "default" {
            policy = "deny"
        }
        agent {
            policy = "deny"
        }
        node {
            policy = "deny"
        }
        `,
		}
		nomadAuthConfig := nomadapi.DefaultConfig()
		nomadAuthConfig.Address = retAddress
		nomadAuthConfig.SecretID = nomadToken
		nomadAuth, err := nomadapi.NewClient(nomadAuthConfig)
		_, err = nomadAuth.ACLPolicies().Upsert(policy, nil)
		if err != nil {
			return err
		}
		_, err = nomadAuth.ACLPolicies().Upsert(anonPolicy, nil)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		cleanup()
		t.Fatalf("Could not connect to docker: %s", err)
	}
	return cleanup, retAddress, nomadToken
}

func TestBackend_config_access(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainer(t)
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
		"address":               connData["address"].(string),
		"max_token_name_length": 0,
	}
	if !reflect.DeepEqual(expected, resp.Data) {
		t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
	}
	if resp.Data["token"] != nil {
		t.Fatalf("token should not be set in the response")
	}
}

func TestBackend_renew_revoke(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainer(t)
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

	req.Path = "role/test"
	req.Data = map[string]interface{}{
		"policies": []string{"policy"},
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
		Token    string `mapstructure:"secret_id"`
		Accessor string `mapstructure:"accessor_id"`
	}
	if err := mapstructure.Decode(resp.Data, &d); err != nil {
		t.Fatal(err)
	}
	t.Logf("[WARN] Generated token: %s with accessor %s", d.Token, d.Accessor)

	// Build a client and verify that the credentials work
	nomadapiConfig := nomadapi.DefaultConfig()
	nomadapiConfig.Address = connData["address"].(string)
	nomadapiConfig.SecretID = d.Token
	client, err := nomadapi.NewClient(nomadapiConfig)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("[WARN] Verifying that the generated token works...")
	_, err = client.Agent().Members, nil
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
	nomadmgmtConfig := nomadapi.DefaultConfig()
	nomadmgmtConfig.Address = connData["address"].(string)
	nomadmgmtConfig.SecretID = connData["token"].(string)
	mgmtclient, err := nomadapi.NewClient(nomadmgmtConfig)

	q := &nomadapi.QueryOptions{
		Namespace: "default",
	}

	t.Log("[WARN] Verifying that the generated token does not exist...")
	_, _, err = mgmtclient.ACLTokens().Info(d.Accessor, q)
	if err == nil {
		t.Fatal("err: expected error")
	}
}

func TestBackend_CredsCreateEnvVar(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainer(t)
	defer cleanup()

	req := logical.TestRequest(t, logical.UpdateOperation, "role/test")
	req.Data = map[string]interface{}{
		"policies": []string{"policy"},
		"lease":    "6h",
	}
	resp, err := b.HandleRequest(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}

	os.Setenv("NOMAD_TOKEN", connToken)
	defer os.Unsetenv("NOMAD_TOKEN")
	os.Setenv("NOMAD_ADDR", connURL)
	defer os.Unsetenv("NOMAD_ADDR")

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
}

func TestBackend_max_token_name_length(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, connURL, connToken := prepareTestContainer(t)
	defer cleanup()

	testCases := []struct {
		title       string
		roleName    string
		tokenLength int
	}{
		{
			title: "Default",
		},
		{
			title:       "ConfigOverride",
			tokenLength: 64,
		},
		{
			title:       "ConfigOverride-LongName",
			roleName:    "testlongerrolenametoexceed64charsdddddddddddddddddddddddd",
			tokenLength: 64,
		},
		{
			title:    "Notrim",
			roleName: "testlongersubrolenametoexceed64charsdddddddddddddddddddddddd",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.title, func(t *testing.T) {
			// setup config/access
			connData := map[string]interface{}{
				"address":               connURL,
				"token":                 connToken,
				"max_token_name_length": tc.tokenLength,
			}
			expected := map[string]interface{}{
				"address":               connURL,
				"max_token_name_length": tc.tokenLength,
			}

			expectedMaxTokenNameLength := maxTokenNameLength
			if tc.tokenLength != 0 {
				expectedMaxTokenNameLength = tc.tokenLength
			}

			confReq := logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "config/access",
				Storage:   config.StorageView,
				Data:      connData,
			}

			resp, err := b.HandleRequest(context.Background(), &confReq)
			if err != nil || (resp != nil && resp.IsError()) || resp != nil {
				t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
			}
			confReq.Operation = logical.ReadOperation
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil || (resp != nil && resp.IsError()) {
				t.Fatalf("failed to write configuration: resp:%#v err:%s", resp, err)
			}

			// verify token length is returned in the config/access query
			if !reflect.DeepEqual(expected, resp.Data) {
				t.Fatalf("bad: expected:%#v\nactual:%#v\n", expected, resp.Data)
			}
			// verify token is not returned
			if resp.Data["token"] != nil {
				t.Fatalf("token should not be set in the response")
			}

			// create a role to create nomad credentials with
			// Seeds random with current timestamp

			if tc.roleName == "" {
				tc.roleName = "test"
			}
			roleTokenName := testhelpers.RandomWithPrefix(tc.roleName)

			confReq.Path = "role/" + roleTokenName
			confReq.Operation = logical.UpdateOperation
			confReq.Data = map[string]interface{}{
				"policies": []string{"policy"},
				"lease":    "6h",
			}
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil {
				t.Fatal(err)
			}

			confReq.Operation = logical.ReadOperation
			confReq.Path = "creds/" + roleTokenName
			resp, err = b.HandleRequest(context.Background(), &confReq)
			if err != nil {
				t.Fatal(err)
			}
			if resp == nil {
				t.Fatal("resp nil")
			}
			if resp.IsError() {
				t.Fatalf("resp is error: %v", resp.Error())
			}

			// extract the secret, so we can query nomad directly
			generatedSecret := resp.Secret
			generatedSecret.TTL = 6 * time.Hour

			var d struct {
				Token    string `mapstructure:"secret_id"`
				Accessor string `mapstructure:"accessor_id"`
			}
			if err := mapstructure.Decode(resp.Data, &d); err != nil {
				t.Fatal(err)
			}

			// Build a client and verify that the credentials work
			nomadapiConfig := nomadapi.DefaultConfig()
			nomadapiConfig.Address = connData["address"].(string)
			nomadapiConfig.SecretID = d.Token
			client, err := nomadapi.NewClient(nomadapiConfig)
			if err != nil {
				t.Fatal(err)
			}

			// default query options for Nomad queries ... not sure if needed
			qOpts := &nomadapi.QueryOptions{
				Namespace: "default",
			}

			// connect to Nomad and verify the token name does not exceed the
			// max_token_name_length
			token, _, err := client.ACLTokens().Self(qOpts)
			if err != nil {
				t.Fatal(err)
			}

			if len(token.Name) > expectedMaxTokenNameLength {
				t.Fatalf("token name exceeds max length (%d): %s (%d)", expectedMaxTokenNameLength, token.Name, len(token.Name))
			}
		})
	}
}
