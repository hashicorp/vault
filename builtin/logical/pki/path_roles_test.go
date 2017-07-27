package pki

import (
	"testing"

	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/mapstructure"
)

func createBackendWithStorage(t *testing.T) (*backend, logical.Storage) {
	config := logical.TestBackendConfig()
	config.StorageView = &logical.InmemStorage{}

	var err error
	b := Backend()
	err = b.Setup(config)
	if err != nil {
		t.Fatal(err)
	}
	return b, config.StorageView
}

func TestPki_RoleGenerateLease(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_domains": "myvault.com",
		"ttl":             "5h",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// generate_lease cannot be nil. It either has to be set during role
	// creation or has to be filled in by the upgrade code
	generateLease := resp.Data["generate_lease"].(*bool)
	if generateLease == nil {
		t.Fatalf("generate_lease should not be nil")
	}

	// By default, generate_lease should be `false`
	if *generateLease {
		t.Fatalf("generate_lease should not be set by default")
	}

	// role.GenerateLease will be nil after the decode
	var role roleEntry
	err = mapstructure.Decode(resp.Data, &role)
	if err != nil {
		t.Fatal(err)
	}

	// Make it explicit
	role.GenerateLease = nil

	entry, err := logical.StorageEntryJSON("role/testrole", role)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(entry); err != nil {
		t.Fatal(err)
	}

	// Reading should upgrade generate_lease
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	generateLease = resp.Data["generate_lease"].(*bool)
	if generateLease == nil {
		t.Fatalf("generate_lease should not be nil")
	}

	// Upgrade should set generate_lease to `true`
	if !*generateLease {
		t.Fatalf("generate_lease should be set after an upgrade")
	}

	// Make sure that setting generate_lease to `true` works properly
	roleReq.Operation = logical.UpdateOperation
	roleReq.Path = "roles/testrole2"
	roleReq.Data["generate_lease"] = true

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	generateLease = resp.Data["generate_lease"].(*bool)
	if generateLease == nil {
		t.Fatalf("generate_lease should not be nil")
	}
	if !*generateLease {
		t.Fatalf("generate_lease should have been set")
	}
}

func TestPki_RoleNoStore(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_domains": "myvault.com",
		"ttl":             "5h",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// By default, no_store should be `false`
	noStore := resp.Data["no_store"].(bool)
	if noStore {
		t.Fatalf("no_store should not be set by default")
	}

	// Make sure that setting no_store to `true` works properly
	roleReq.Operation = logical.UpdateOperation
	roleReq.Path = "roles/testrole_nostore"
	roleReq.Data["no_store"] = true
	roleReq.Data["allowed_domain"] = "myvault.com"
	roleReq.Data["allow_subdomains"] = true
	roleReq.Data["ttl"] = "5h"

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	noStore = resp.Data["no_store"].(bool)
	if !noStore {
		t.Fatalf("no_store should have been set to true")
	}

	// issue a certificate and test that it's not stored
	caData := map[string]interface{}{
		"common_name": "myvault.com",
		"ttl":         "5h",
		"ip_sans":     "127.0.0.1",
	}
	caReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data:      caData,
	}
	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	issueData := map[string]interface{}{
		"common_name": "cert.myvault.com",
		"format":      "pem",
		"ip_sans":     "127.0.0.1",
		"ttl":         "1h",
	}
	issueReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/testrole_nostore",
		Storage:   storage,
		Data:      issueData,
	}

	resp, err = b.HandleRequest(issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// list certs
	resp, err = b.HandleRequest(&logical.Request{
		Operation: logical.ListOperation,
		Path:      "certs",
		Storage:   storage,
	})
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}
	if len(resp.Data["keys"].([]string)) != 1 {
		t.Fatalf("Only the CA certificate should be stored: %#v", resp)
	}
}

func TestPki_CertsLease(t *testing.T) {
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

	caData := map[string]interface{}{
		"common_name": "myvault.com",
		"ttl":         "5h",
		"ip_sans":     "127.0.0.1",
	}

	caReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "root/generate/internal",
		Storage:   storage,
		Data:      caData,
	}

	resp, err = b.HandleRequest(caReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleData := map[string]interface{}{
		"allowed_domains":  "myvault.com",
		"allow_subdomains": true,
		"ttl":              "2h",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	issueData := map[string]interface{}{
		"common_name": "cert.myvault.com",
		"format":      "pem",
		"ip_sans":     "127.0.0.1",
	}
	issueReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "issue/testrole",
		Storage:   storage,
		Data:      issueData,
	}

	resp, err = b.HandleRequest(issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	if resp.Secret != nil {
		t.Fatalf("expected a response that does not contain a secret")
	}

	// Turn on the lease generation and issue a certificate. The response
	// should have a `Secret` object populated.
	roleData["generate_lease"] = true

	resp, err = b.HandleRequest(roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	resp, err = b.HandleRequest(issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	if resp.Secret == nil {
		t.Fatalf("expected a response that contains a secret")
	}
}
