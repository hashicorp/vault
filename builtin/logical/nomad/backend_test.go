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
	"github.com/hashicorp/vault/helper/testhelpers/docker"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	docker.ServiceURL
	Token string
}

func (c *Config) APIConfig() *nomadapi.Config {
	apiConfig := nomadapi.DefaultConfig()
	apiConfig.Address = c.URL().String()
	apiConfig.SecretID = c.Token
	return apiConfig
}

func prepareTestContainer(t *testing.T) (func(), *Config) {
	if retAddress := os.Getenv("NOMAD_ADDR"); retAddress != "" {
		s, err := docker.NewServiceURLParse(retAddress)
		if err != nil {
			t.Fatal(err)
		}
		return func() {}, &Config{*s, os.Getenv("NOMAD_TOKEN")}
	}

	runner, err := docker.NewServiceRunner(docker.RunOptions{
		ImageRepo:     "catsby/nomad",
		ImageTag:      "0.8.4",
		ContainerName: "nomad",
		Ports:         []string{"4646/tcp"},
		Cmd:           []string{"agent", "-dev"},
		Env:           []string{`NOMAD_LOCAL_CONFIG=bind_addr = "0.0.0.0" acl { enabled = true }`},
	})
	if err != nil {
		t.Fatalf("Could not start docker Nomad: %s", err)
	}

	var nomadToken string
	svc, err := runner.StartService(context.Background(), func(ctx context.Context, host string, port int) (docker.ServiceConfig, error) {
		var err error
		nomadapiConfig := nomadapi.DefaultConfig()
		nomadapiConfig.Address = fmt.Sprintf("http://%s:%d/", host, port)
		nomad, err := nomadapi.NewClient(nomadapiConfig)
		if err != nil {
			return nil, err
		}
		aclbootstrap, _, err := nomad.ACLTokens().Bootstrap(nil)
		if err != nil {
			return nil, err
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
		nomadAuthConfig.Address = nomad.Address()
		nomadAuthConfig.SecretID = nomadToken
		nomadAuth, err := nomadapi.NewClient(nomadAuthConfig)
		if err != nil {
			return nil, err
		}
		_, err = nomadAuth.ACLPolicies().Upsert(policy, nil)
		if err != nil {
			return nil, err
		}
		_, err = nomadAuth.ACLPolicies().Upsert(anonPolicy, nil)
		if err != nil {
			return nil, err
		}
		u, _ := docker.NewServiceURLParse(nomadapiConfig.Address)
		return &Config{
			ServiceURL: *u,
			Token:      nomadToken,
		}, nil
	})
	if err != nil {
		t.Fatalf("Could not start docker Nomad: %s", err)
	}

	return svc.Cleanup, svc.Config.(*Config)
}

func TestBackend_config_access(t *testing.T) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}
	b, err := Factory(context.Background(), config)
	if err != nil {
		t.Fatal(err)
	}

	cleanup, svccfg := prepareTestContainer(t)
	defer cleanup()

	connData := map[string]interface{}{
		"address": svccfg.URL().String(),
		"token":   svccfg.Token,
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

	cleanup, svccfg := prepareTestContainer(t)
	defer cleanup()

	connData := map[string]interface{}{
		"address": svccfg.URL().String(),
		"token":   svccfg.Token,
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
	if err != nil {
		t.Fatal(err)
	}

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

	cleanup, svccfg := prepareTestContainer(t)
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

	os.Setenv("NOMAD_TOKEN", svccfg.Token)
	defer os.Unsetenv("NOMAD_TOKEN")
	os.Setenv("NOMAD_ADDR", svccfg.URL().String())
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

	cleanup, svccfg := prepareTestContainer(t)
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
				"address":               svccfg.URL().String(),
				"token":                 svccfg.Token,
				"max_token_name_length": tc.tokenLength,
			}
			expected := map[string]interface{}{
				"address":               svccfg.URL().String(),
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
