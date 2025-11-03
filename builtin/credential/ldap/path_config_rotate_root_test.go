// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-2.0

package ldap

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/vault/helper/testhelpers/ldap"
	logicaltest "github.com/hashicorp/vault/helper/testhelpers/logical"
	"github.com/hashicorp/vault/sdk/helper/ldaputil"
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

// TestRotateRootWithRotationUrl relies on a docker ldap server with a suitable person object (cn=admin,dc=planetexpress,dc=com)
// with bindpassword "admin". `PrepareTestContainer` does this for us. - see the backend_test for more details
// It checks that rotation url is being used instead of the main URL and assures that setting rotation url does't
// replace main URL
func TestRotateRootWithRotationUrl(t *testing.T) {
	if os.Getenv(logicaltest.TestEnvVar) == "" {
		t.Skip("skipping rotate root tests because VAULT_ACC is unset")
	}
	ctx := context.Background()

	b, store := createBackendWithStorage(t)
	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	const mainDummyUrl = "ldap://example.com:389"
	// set up auth config
	req := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Storage:   store,
		Data: map[string]interface{}{
			"url":          mainDummyUrl,
			"binddn":       cfg.BindDN,
			"bindpass":     cfg.BindPassword,
			"userdn":       cfg.UserDN,
			"rotation_url": cfg.Url,
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
	if err != nil {
		t.Fatalf("failed to read config after rotation: %s", err)
	}
	if newCFG.BindDN != cfg.BindDN {
		t.Fatalf("BindDN %q changed unexpectedly, found new value %q", cfg.BindDN, newCFG.BindDN)
	}
	if newCFG.BindPassword == cfg.BindPassword {
		t.Fatalf("the password should have changed, but it didn't")
	}
	// expecting the newCFG url to be "ldap://example.com:389"
	if newCFG.Url != mainDummyUrl {
		t.Fatalf("URL %q changed unexpectedly, found new value %q", mainDummyUrl, newCFG.Url)
	}
}

// TestGetModifyRequest tests that the correct LDAP modify requests are generated
// for different rotation schemas and credential types.
func TestGetModifyRequest(t *testing.T) {
	b, _ := createBackendWithStorage(t)
	cfgE := new(ldaputil.ConfigEntry)
	cfgE.BindDN = "cn=admin,dc=planetexpress,dc=com"
	cfg := &ldapConfigEntry{
		ConfigEntry:            cfgE,
		RotationSchema:         schemaOpenLDAP,
		RotationCredentialType: credentialTypePassword,
	}
	dummyPassword := "newpassword123"
	// Test OpenLDAP schema
	lreq, err := b.getModifyRequest(cfg, dummyPassword)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(lreq.Changes) != 1 {
		t.Fatalf("expected 1 change, got: %d", len(lreq.Changes))
	}
	if lreq.Changes[0].Modification.Type != "userPassword" {
		t.Fatalf("expected userPassword attribute to be modified, got: %s", lreq.Changes[0].Modification.Type)
	}
	if lreq.Changes[0].Modification.Vals[0] != dummyPassword {
		t.Fatalf("expected new password to be set to newpassword123, got: %s", lreq.Changes[0].Modification.Vals[0])
	}
	// Test Active Directory schema
	cfg.RotationSchema = schemaAD
	lreq, err = b.getModifyRequest(cfg, dummyPassword)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(lreq.Changes) != 1 {
		t.Fatalf("expected 1 change, got: %d", len(lreq.Changes))
	}
	if lreq.Changes[0].Modification.Type != "unicodePwd" {
		t.Fatalf("expected unicodePwd attribute to be modified, got: %s", lreq.Changes[0].Modification.Type)
	}
	pwdEncoded, err := formatPassword(dummyPassword)
	if err != nil {
		t.Fatalf("unexpected error encoding password: %s", err)
	}
	if lreq.Changes[0].Modification.Vals[0] != pwdEncoded {
		t.Fatalf("expected new password to be encoded, got: %s", lreq.Changes[0].Modification.Vals[0])
	}
	// Test RACF schema with password type
	cfg.RotationSchema = schemaRACF
	lreq, err = b.getModifyRequest(cfg, dummyPassword)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(lreq.Changes) != 2 {
		t.Fatalf("expected 2 change, got: %d", len(lreq.Changes))
	}
	if lreq.Changes[0].Modification.Type != "racfPassword" {
		t.Fatalf("expected racfPassword attribute to be modified, got: %s", lreq.Changes[0].Modification.Type)
	}
	if lreq.Changes[0].Modification.Vals[0] != dummyPassword {
		t.Fatalf("expected new password to be set to newpassword123, got: %s", lreq.Changes[0].Modification.Vals[0])
	}
	if lreq.Changes[1].Modification.Type != "racfAttributes" {
		t.Fatalf("expected racfAttributes attribute to be modified, got: %s", lreq.Changes[1].Modification.Type)
	}
	if lreq.Changes[1].Modification.Vals[0] != "noexpired" {
		t.Fatalf("expected racfAttributes to be set to noexpired, got: %s", lreq.Changes[1].Modification.Vals[0])
	}
	// Test RACF schema with passphrase type
	cfg.RotationCredentialType = credentialTypePhrase
	lreq, err = b.getModifyRequest(cfg, dummyPassword)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(lreq.Changes) != 2 {
		t.Fatalf("expected 2 change, got: %d", len(lreq.Changes))
	}
	if lreq.Changes[0].Modification.Type != "racfPassPhrase" {
		t.Fatalf("expected racfPassPhrase attribute to be modified, got: %s", lreq.Changes[0].Modification.Type)
	}
	if lreq.Changes[0].Modification.Vals[0] != dummyPassword {
		t.Fatalf("expected new passphrase to be set to newpassword123, got: %s", lreq.Changes[0].Modification.Vals[0])
	}
	if lreq.Changes[1].Modification.Type != "racfAttributes" {
		t.Fatalf("expected racfAttributes attribute to be modified, got: %s", lreq.Changes[1].Modification.Type)
	}
	if lreq.Changes[1].Modification.Vals[0] != "noexpired" {
		t.Fatalf("expected racfAttributes to be set to noexpired, got: %s", lreq.Changes[1].Modification.Vals[0])
	}
	// Test invalid schema
	cfg.RotationSchema = "invalidSchema"
	_, err = b.getModifyRequest(cfg, dummyPassword)
	if err == nil {
		t.Fatalf("expected error due to invalid schema, got none")
	}
}

// TestRotateRootConfigs tests that the rotate-root configuration options
// are correctly set and validated.
func TestRotateRootConfigs(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	ctx := context.Background()

	cleanup, cfg := ldap.PrepareTestContainer(t, ldap.DefaultVersion)
	defer cleanup()
	testCases := []struct {
		name           string
		rotateSchema   string
		rotateCredType string
		expectError    bool
	}{
		{
			name:           "With rotation_schema=openldap and rotation_credential_type=password",
			rotateSchema:   schemaOpenLDAP,
			rotateCredType: credentialTypePassword,
			expectError:    false,
		},
		{
			name:           "With rotation_schema=ad and rotation_credential_type=password",
			rotateSchema:   schemaAD,
			rotateCredType: credentialTypePassword,
			expectError:    false,
		},
		{
			name:           "With rotation_schema=racf and rotation_credential_type=password",
			rotateSchema:   schemaRACF,
			rotateCredType: credentialTypePassword,
			expectError:    false,
		},
		{
			name:           "With rotation_schema=racf and rotation_credential_type=passphrase",
			rotateSchema:   schemaRACF,
			rotateCredType: credentialTypePhrase,
			expectError:    false,
		},
		{
			name:           "With invalid rotation_schema",
			rotateSchema:   "invalidschema",
			rotateCredType: credentialTypePassword,
			expectError:    true,
		},
		{
			name:           "With invalid rotation_credential_type",
			rotateSchema:   schemaOpenLDAP,
			rotateCredType: "invalidcredtype",
			expectError:    true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure the backend
			configReq := &logical.Request{
				Operation: logical.UpdateOperation,
				Path:      "config",
				Data: map[string]interface{}{
					"url":                         cfg.Url,
					"userattr":                    cfg.UserAttr,
					"userdn":                      cfg.UserDN,
					"userfilter":                  cfg.UserFilter,
					"groupdn":                     cfg.GroupDN,
					"groupattr":                   cfg.GroupAttr,
					"binddn":                      cfg.BindDN,
					"bindpass":                    cfg.BindPassword,
					"token_period":                "5m",
					"token_explicit_max_ttl":      "24h",
					"request_timeout":             cfg.RequestTimeout,
					"max_page_size":               cfg.MaximumPageSize,
					"connection_timeout":          cfg.ConnectionTimeout,
					rootRotationSchemaKey:         tc.rotateSchema,
					rootRotationCredentialTypeKey: tc.rotateCredType,
				},
				Storage:    storage,
				Connection: &logical.Connection{},
			}
			resp, err = b.HandleRequest(ctx, configReq)
			if tc.expectError && (err == nil && (resp == nil || !resp.IsError())) {
				t.Fatalf("expected error but got none")
			}
			if !tc.expectError && (err != nil || (resp != nil && resp.IsError())) {
				t.Fatalf("did not expect error but got err:%v resp:%#v", err, resp)
			}
		})
	}
	// test defaults
	configReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "config",
		Data: map[string]interface{}{
			"url":                    cfg.Url,
			"userattr":               cfg.UserAttr,
			"userdn":                 cfg.UserDN,
			"userfilter":             cfg.UserFilter,
			"groupdn":                cfg.GroupDN,
			"groupattr":              cfg.GroupAttr,
			"binddn":                 cfg.BindDN,
			"bindpass":               cfg.BindPassword,
			"token_period":           "5m",
			"token_explicit_max_ttl": "24h",
			"request_timeout":        cfg.RequestTimeout,
			"max_page_size":          cfg.MaximumPageSize,
			"connection_timeout":     cfg.ConnectionTimeout,
		},
		Storage:    storage,
		Connection: &logical.Connection{},
	}
	resp, err = b.HandleRequest(ctx, configReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("did not expect error but got err:%v resp:%#v", err, resp)
	}
	// verify defaults
	resp, err = b.HandleRequest(ctx, &logical.Request{
		Operation: logical.ReadOperation,
		Path:      "config",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("did not expect error but got err:%v resp:%#v", err, resp)
	}
	if resp.Data[rootRotationSchemaKey].(string) != schemaOpenLDAP {
		t.Fatalf("expected default rotation schema to be openldap, got: %s", resp.Data[rootRotationSchemaKey].(string))
	}
	if resp.Data[rootRotationCredentialTypeKey].(string) != credentialTypePassword {
		t.Fatalf("expected default rotation credential type to be password, got: %s", resp.Data[rootRotationCredentialTypeKey].(string))
	}
}
