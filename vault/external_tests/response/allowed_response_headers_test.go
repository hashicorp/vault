package identity

import (
	"context"
	"testing"

	"github.com/hashicorp/vault/api"
	vaulthttp "github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/vault"
	"github.com/kr/pretty"
)

func TestIdentityStore_EntityDisabled(t *testing.T) {
	be := &framework.Backend{
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login*",
			},
		},
		Paths: []*framework.Path{
			{
				Pattern: "login",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
						return &logical.Response{
							Headers: map[string][]string{
								"www-authenticate": []string{"Negotiate"},
							},
						}, logical.CodedError(401, "authentication required")
					},
				},
			},
			{
				Pattern: "loginnoerror",
				Callbacks: map[logical.Operation]framework.OperationFunc{
					logical.ReadOperation: func(context.Context, *logical.Request, *framework.FieldData) (*logical.Response, error) {
						return &logical.Response{
							Auth: &logical.Auth{},
							Headers: map[string][]string{
								"www-authenticate": []string{"Negotiate"},
							},
						}, nil
					},
				},
			},
		},
		BackendType: logical.TypeCredential,
	}

	// Use a TestCluster and the approle backend to get a token and entity for testing
	coreConfig := &vault.CoreConfig{
		CredentialBackends: map[string]logical.Factory{
			"headtest": func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
				err := be.Setup(ctx, conf)
				if err != nil {
					return nil, err
				}
				return be, nil
			},
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: vaulthttp.Handler,
	})
	cluster.Start()
	defer cluster.Cleanup()

	core := cluster.Cores[0].Core
	vault.TestWaitActive(t, core)
	client := cluster.Cores[0].Client

	// Mount the auth backend
	err := client.Sys().EnableAuthWithOptions("headtest", &api.EnableAuthOptions{
		Type: "headtest",
	})
	if err != nil {
		t.Fatal(err)
	}

	// Here, should suceed but we should not see the header since it's
	// not in the allowed list
	req := client.NewRequest("GET", "/v1/auth/headtest/loginnoerror")
	resp, err := client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected code 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("www-authenticate") != "" {
		t.Fatal("expected header to not be allowed")
	}

	// Should fail but we should not see the header since it's
	// not in the allowed list
	req = client.NewRequest("GET", "/v1/auth/headtest/login")
	resp, err = client.RawRequest(req)
	if err == nil {
		t.Fatal("expected error")
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected code 401, got %d", resp.StatusCode)
	}
	if resp.Header.Get("www-authenticate") != "" {
		t.Fatal("expected header to not be allowed")
	}

	// Tune the mount
	err = client.Sys().TuneMount("auth/headtest", api.MountConfigInput{
		AllowedResponseHeaders: []string{"WwW-AuthenTicate"},
	})
	if err != nil {
		t.Fatal(err)
	}

	req = client.NewRequest("GET", "/v1/auth/headtest/loginnoerror")
	resp, err = client.RawRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected code 200, got %d", resp.StatusCode)
	}
	if resp.Header.Get("www-authenticate") != "Negotiate" {
		t.Fatalf("expected negotiate header; headers:\n%s", pretty.Sprint(resp.Header))
	}

	req = client.NewRequest("GET", "/v1/auth/headtest/login")
	resp, err = client.RawRequest(req)
	if err == nil {
		t.Fatal("expected error")
	}
	if resp.StatusCode != 401 {
		t.Fatalf("expected code 401, got %d", resp.StatusCode)
	}
	if resp.Header.Get("www-authenticate") != "Negotiate" {
		t.Fatalf("expected negotiate header; headers:\n%s", pretty.Sprint(resp.Header))
	}
}
