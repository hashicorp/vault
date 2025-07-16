// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/ldap"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/logical"
)

// This test relies on a docker ldap server with a suitable person object (cn=admin,dc=planetexpress,dc=com)
// with bindpassword "admin". `PrepareTestContainer` does this for us. - see the backend_test for more details
func TestRotateRoot(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()

	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":      cfg.Url,
			"binddn":   cfg.BindDN,
			"bindpass": cfg.BindPassword,
			"userdn":   cfg.UserDN,
		},
	}

	resp, err := b.HandleRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to initialize ldap auth config: %s", err)
	}
	if resp != nil && resp.IsError() {
		t.Fatalf("failed to initialize ldap auth config: %s", resp.Data["error"])
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

	newCFG, err := b.Config(ctx, req)
	if newCFG.BindDN != cfg.BindDN {
		t.Fatalf("a value in config that should have stayed the same changed: %s", cfg.BindDN)
	}
	if newCFG.BindPassword == cfg.BindPassword {
		t.Fatalf("the password should have changed, but it didn't")
	}
}
