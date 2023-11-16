package ldap

import (
	"context"
	"os"
	"testing"

	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
)

// This test relies on an external ldap server with a suitable person object (cn=admin,dc=planetexpress,dc=com)
// with bindpassword "admin". - see the backend_test for more details.
// This test will not run unless VAULT_ACC is set to something
func TestRotateRoot(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()

	b, store := createBackendWithStorage(t)
	defer b.Cleanup(ctx)

	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":      "ldap://localhost:389",
			"binddn":   "cn=admin,dc=planetexpress,dc=com",
			"bindpass": "admin",
			"userdn":   "dc=planetexpress,dc=com",
		},
	}

	_, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}

	req = &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config/rotate-root",
		Storage:   store,
	}

	_, err = b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to rotate password: %s", err)
	}

	cfg, err := b.Config(ctx, req)
	if cfg.BindDN != "cn=admin,dc=planetexpress,dc=com" {
		t.Fatalf("a value in config that should have stayed the same changed: %s", cfg.BindDN)
	}
	if cfg.BindPassword == "admin" {
		t.Fatalf("the password should have changed, but it didn't")
	}
}
