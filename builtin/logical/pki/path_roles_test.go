// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/go-errors/errors"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/helper/testhelpers/schema"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPki_RoleGenerateLease(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

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

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
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

	// To test upgrade of generate_lease, we read the storage entry,
	// modify it to remove generate_lease, and rewrite it.
	entry, err := storage.Get(context.Background(), "role/testrole")
	if err != nil || entry == nil {
		t.Fatal(err)
	}

	var role issuing.RoleEntry
	if err := entry.DecodeJSON(&role); err != nil {
		t.Fatal(err)
	}

	role.GenerateLease = nil

	entry, err = logical.StorageEntryJSON("role/testrole", role)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// Reading should upgrade generate_lease
	resp, err = b.HandleRequest(context.Background(), roleReq)
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

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
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

func TestPki_RoleKeyUsage(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_domains": "myvault.com",
		"ttl":             "5h",
		"key_usage":       []string{"KeyEncipherment", "DigitalSignature"},
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route(roleReq.Path), logical.UpdateOperation), resp, true)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	schema.ValidateResponse(t, schema.GetResponseSchema(t, b.Route(roleReq.Path), logical.ReadOperation), resp, true)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	keyUsage := resp.Data["key_usage"].([]string)
	if len(keyUsage) != 2 {
		t.Fatalf("key_usage should have 2 values")
	}

	// To test the upgrade of KeyUsageOld into KeyUsage, we read
	// the storage entry, modify it to set KUO and unset KU, and
	// rewrite it.
	entry, err := storage.Get(context.Background(), "role/testrole")
	if err != nil || entry == nil {
		t.Fatal(err)
	}

	var role issuing.RoleEntry
	if err := entry.DecodeJSON(&role); err != nil {
		t.Fatal(err)
	}

	role.KeyUsageOld = "KeyEncipherment,DigitalSignature"
	role.KeyUsage = nil

	entry, err = logical.StorageEntryJSON("role/testrole", role)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// Reading should upgrade key_usage
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	keyUsage = resp.Data["key_usage"].([]string)
	if len(keyUsage) != 2 {
		t.Fatalf("key_usage should have 2 values")
	}

	// Read back from storage to ensure upgrade
	entry, err = storage.Get(context.Background(), "role/testrole")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("role should not be nil")
	}
	var result issuing.RoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		t.Fatalf("err: %v", err)
	}

	if result.KeyUsageOld != "" {
		t.Fatal("old key usage value should be blank")
	}
	if len(result.KeyUsage) != 2 {
		t.Fatal("key_usage should have 2 values")
	}
}

func TestPki_RoleOUOrganizationUpgrade(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_domains": "myvault.com",
		"ttl":             "5h",
		"ou":              []string{"abc", "123"},
		"organization":    []string{"org1", "org2"},
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	ou := resp.Data["ou"].([]string)
	if len(ou) != 2 {
		t.Fatalf("ou should have 2 values")
	}
	organization := resp.Data["organization"].([]string)
	if len(organization) != 2 {
		t.Fatalf("organization should have 2 values")
	}

	// To test upgrade of O/OU, we read the storage entry, modify it to set
	// the old O/OU value over the new one, and rewrite it.
	entry, err := storage.Get(context.Background(), "role/testrole")
	if err != nil || entry == nil {
		t.Fatal(err)
	}

	var role issuing.RoleEntry
	if err := entry.DecodeJSON(&role); err != nil {
		t.Fatal(err)
	}
	role.OUOld = "abc,123"
	role.OU = nil
	role.OrganizationOld = "org1,org2"
	role.Organization = nil

	entry, err = logical.StorageEntryJSON("role/testrole", role)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// Reading should upgrade key_usage
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	ou = resp.Data["ou"].([]string)
	if len(ou) != 2 {
		t.Fatalf("ou should have 2 values")
	}
	organization = resp.Data["organization"].([]string)
	if len(organization) != 2 {
		t.Fatalf("organization should have 2 values")
	}

	// Read back from storage to ensure upgrade
	entry, err = storage.Get(context.Background(), "role/testrole")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("role should not be nil")
	}
	var result issuing.RoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		t.Fatalf("err: %v", err)
	}

	if result.OUOld != "" {
		t.Fatal("old ou value should be blank")
	}
	if len(result.OU) != 2 {
		t.Fatal("ou should have 2 values")
	}
	if result.OrganizationOld != "" {
		t.Fatal("old organization value should be blank")
	}
	if len(result.Organization) != 2 {
		t.Fatal("organization should have 2 values")
	}
}

func TestPki_RoleAllowedDomains(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_domains": []string{"foobar.com", "*example.com"},
		"ttl":             "5h",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	allowedDomains := resp.Data["allowed_domains"].([]string)
	if len(allowedDomains) != 2 {
		t.Fatalf("allowed_domains should have 2 values")
	}

	// To test upgrade of allowed_domains, we read the storage entry,
	// set the old one, and rewrite it.
	entry, err := storage.Get(context.Background(), "role/testrole")
	if err != nil || entry == nil {
		t.Fatal(err)
	}

	var role issuing.RoleEntry
	if err := entry.DecodeJSON(&role); err != nil {
		t.Fatal(err)
	}
	role.AllowedDomainsOld = "foobar.com,*example.com"
	role.AllowedDomains = nil

	entry, err = logical.StorageEntryJSON("role/testrole", role)
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.Put(context.Background(), entry); err != nil {
		t.Fatal(err)
	}

	// Reading should upgrade key_usage
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	allowedDomains = resp.Data["allowed_domains"].([]string)
	if len(allowedDomains) != 2 {
		t.Fatalf("allowed_domains should have 2 values")
	}

	// Read back from storage to ensure upgrade
	entry, err = storage.Get(context.Background(), "role/testrole")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if entry == nil {
		t.Fatalf("role should not be nil")
	}
	var result issuing.RoleEntry
	if err := entry.DecodeJSON(&result); err != nil {
		t.Fatalf("err: %v", err)
	}

	if result.AllowedDomainsOld != "" {
		t.Fatal("old allowed_domains value should be blank")
	}
	if len(result.AllowedDomains) != 2 {
		t.Fatal("allowed_domains should have 2 values")
	}
}

func TestPki_RoleAllowedURISANs(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

	roleData := map[string]interface{}{
		"allowed_uri_sans": []string{"http://foobar.com", "spiffe://*"},
		"ttl":              "5h",
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	allowedURISANs := resp.Data["allowed_uri_sans"].([]string)
	if len(allowedURISANs) != 2 {
		t.Fatalf("allowed_uri_sans should have 2 values")
	}
}

func TestPki_RolePkixFields(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

	roleData := map[string]interface{}{
		"ttl":            "5h",
		"country":        []string{"c1", "c2"},
		"ou":             []string{"abc", "123"},
		"organization":   []string{"org1", "org2"},
		"locality":       []string{"foocity", "bartown"},
		"province":       []string{"bar", "foo"},
		"street_address": []string{"123 foo street", "789 bar avenue"},
		"postal_code":    []string{"f00", "b4r"},
	}

	roleReq := &logical.Request{
		Operation: logical.UpdateOperation,
		Path:      "roles/testrole_pkixfields",
		Storage:   storage,
		Data:      roleData,
	}

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	origCountry := roleData["country"].([]string)
	respCountry := resp.Data["country"].([]string)
	if !strutil.StrListSubset(origCountry, respCountry) {
		t.Fatalf("country did not match values set in role")
	} else if len(origCountry) != len(respCountry) {
		t.Fatalf("country did not have same number of values set in role")
	}

	origOU := roleData["ou"].([]string)
	respOU := resp.Data["ou"].([]string)
	if !strutil.StrListSubset(origOU, respOU) {
		t.Fatalf("ou did not match values set in role")
	} else if len(origOU) != len(respOU) {
		t.Fatalf("ou did not have same number of values set in role")
	}

	origOrganization := roleData["organization"].([]string)
	respOrganization := resp.Data["organization"].([]string)
	if !strutil.StrListSubset(origOrganization, respOrganization) {
		t.Fatalf("organization did not match values set in role")
	} else if len(origOrganization) != len(respOrganization) {
		t.Fatalf("organization did not have same number of values set in role")
	}

	origLocality := roleData["locality"].([]string)
	respLocality := resp.Data["locality"].([]string)
	if !strutil.StrListSubset(origLocality, respLocality) {
		t.Fatalf("locality did not match values set in role")
	} else if len(origLocality) != len(respLocality) {
		t.Fatalf("locality did not have same number of values set in role: ")
	}

	origProvince := roleData["province"].([]string)
	respProvince := resp.Data["province"].([]string)
	if !strutil.StrListSubset(origProvince, respProvince) {
		t.Fatalf("province did not match values set in role")
	} else if len(origProvince) != len(respProvince) {
		t.Fatalf("province did not have same number of values set in role")
	}

	origStreetAddress := roleData["street_address"].([]string)
	respStreetAddress := resp.Data["street_address"].([]string)
	if !strutil.StrListSubset(origStreetAddress, respStreetAddress) {
		t.Fatalf("street_address did not match values set in role")
	} else if len(origStreetAddress) != len(respStreetAddress) {
		t.Fatalf("street_address did not have same number of values set in role")
	}

	origPostalCode := roleData["postal_code"].([]string)
	respPostalCode := resp.Data["postal_code"].([]string)
	if !strutil.StrListSubset(origPostalCode, respPostalCode) {
		t.Fatalf("postal_code did not match values set in role")
	} else if len(origPostalCode) != len(respPostalCode) {
		t.Fatalf("postal_code did not have same number of values set in role")
	}
}

func TestPki_RoleNoStore(t *testing.T) {
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

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

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// By default, no_store should be `false`
	noStore := resp.Data["no_store"].(bool)
	if noStore {
		t.Fatalf("no_store should not be set by default")
	}

	// By default, allowed_domains_template should be `false`
	allowedDomainsTemplate := resp.Data["allowed_domains_template"].(bool)
	if allowedDomainsTemplate {
		t.Fatalf("allowed_domains_template should not be set by default")
	}

	// By default, allowed_uri_sans_template should be `false`
	allowedURISANsTemplate := resp.Data["allowed_uri_sans_template"].(bool)
	if allowedURISANsTemplate {
		t.Fatalf("allowed_uri_sans_template should not be set by default")
	}

	// Make sure that setting no_store to `true` works properly
	roleReq.Operation = logical.UpdateOperation
	roleReq.Path = "roles/testrole_nostore"
	roleReq.Data["no_store"] = true
	roleReq.Data["allowed_domain"] = "myvault.com"
	roleReq.Data["allow_subdomains"] = true
	roleReq.Data["ttl"] = "5h"

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
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
	resp, err = b.HandleRequest(context.Background(), caReq)
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

	resp, err = b.HandleRequest(context.Background(), issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	// list certs
	resp, err = b.HandleRequest(context.Background(), &logical.Request{
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
	t.Parallel()
	var resp *logical.Response
	var err error
	b, storage := CreateBackendWithStorage(t)

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

	resp, err = b.HandleRequest(context.Background(), caReq)
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

	resp, err = b.HandleRequest(context.Background(), roleReq)
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

	resp, err = b.HandleRequest(context.Background(), issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	if resp.Secret != nil {
		t.Fatalf("expected a response that does not contain a secret")
	}

	// Turn on the lease generation and issue a certificate. The response
	// should have a `Secret` object populated.
	roleData["generate_lease"] = true

	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	resp, err = b.HandleRequest(context.Background(), issueReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	if resp.Secret == nil {
		t.Fatalf("expected a response that contains a secret")
	}
}

func TestPki_RolePatch(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Field   string
		Before  interface{}
		Patched interface{}
	}

	testCases := []TestCase{
		{
			Field:   "ttl",
			Before:  int64(5),
			Patched: int64(10),
		},
		{
			Field:   "max_ttl",
			Before:  int64(5),
			Patched: int64(10),
		},
		{
			Field:   "allow_localhost",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "allowed_domains",
			Before:  []string{"alex", "bob"},
			Patched: []string{"sam", "alex", "frank"},
		},
		{
			Field:   "allowed_domains_template",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "allow_bare_domains",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "allow_subdomains",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "allow_glob_domains",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "allow_wildcard_certificates",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "allow_any_name",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "enforce_hostnames",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "allow_ip_sans",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "allowed_uri_sans",
			Before:  []string{"gopher://*"},
			Patched: []string{"https://*"},
		},
		{
			Field:   "allowed_uri_sans_template",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "allowed_other_sans",
			Before:  []string{"1.2.3.4;UTF8:magic"},
			Patched: []string{"4.3.2.1;UTF8:cigam"},
		},
		{
			Field:   "allowed_serial_numbers",
			Before:  []string{"*"},
			Patched: []string{""},
		},
		{
			Field:   "server_flag",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "client_flag",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "code_signing_flag",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "email_protection_flag",
			Before:  false,
			Patched: true,
		},
		// key_type, key_bits, and signature_bits can't be tested in this setup
		// due to their non-default stored nature.
		{
			Field:   "key_usage",
			Before:  []string{"DigitialSignature"},
			Patched: []string{"DigitalSignature", "KeyAgreement"},
		},
		{
			Field:   "ext_key_usage",
			Before:  []string{"ServerAuth"},
			Patched: []string{"ClientAuth"},
		},
		{
			Field:   "ext_key_usage_oids",
			Before:  []string{"1.2.3.4"},
			Patched: []string{"4.3.2.1"},
		},
		{
			Field:   "use_csr_common_name",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "use_csr_sans",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "serial_number_source",
			Before:  "json-csr",
			Patched: "json",
		},
		{
			Field:   "ou",
			Before:  []string{"crypto"},
			Patched: []string{"cryptosec"},
		},
		{
			Field:   "organization",
			Before:  []string{"hashicorp"},
			Patched: []string{"dadgarcorp"},
		},
		{
			Field:   "country",
			Before:  []string{"US"},
			Patched: []string{"USA"},
		},
		{
			Field:   "locality",
			Before:  []string{"Orange"},
			Patched: []string{"Blue"},
		},
		{
			Field:   "province",
			Before:  []string{"CA"},
			Patched: []string{"AC"},
		},
		{
			Field:   "street_address",
			Before:  []string{"101 First"},
			Patched: []string{"202 Second", "Unit 020"},
		},
		{
			Field:   "postal_code",
			Before:  []string{"12345"},
			Patched: []string{"54321-1234"},
		},
		{
			Field:   "generate_lease",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "no_store",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "require_cn",
			Before:  false,
			Patched: true,
		},
		{
			Field:   "policy_identifiers",
			Before:  []string{"1.3.6.1.4.1.1.1"},
			Patched: []string{"1.3.6.1.4.1.1.2"},
		},
		{
			Field:   "basic_constraints_valid_for_non_ca",
			Before:  true,
			Patched: false,
		},
		{
			Field:   "not_before_duration",
			Before:  int64(30),
			Patched: int64(300),
		},
		{
			Field:   "not_after",
			Before:  "9999-12-31T23:59:59Z",
			Patched: "1230-12-31T23:59:59Z",
		},
		{
			Field:   "issuer_ref",
			Before:  "default",
			Patched: "missing",
		},
	}

	b, storage := CreateBackendWithStorage(t)

	for index, testCase := range testCases {
		var resp *logical.Response
		var roleDataResp *logical.Response
		var afterRoleDataResp *logical.Response
		var err error

		// Create the role
		roleData := map[string]interface{}{}
		roleData[testCase.Field] = testCase.Before

		roleReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "roles/testrole",
			Storage:   storage,
			Data:      roleData,
		}

		resp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad [%d/%v] create: err: %v resp: %#v", index, testCase.Field, err, resp)
		}

		// Read the role after creation
		roleReq.Operation = logical.ReadOperation
		roleDataResp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (roleDataResp != nil && roleDataResp.IsError()) {
			t.Fatalf("bad [%d/%v] read: err: %v resp: %#v", index, testCase.Field, err, resp)
		}

		beforeRoleData := roleDataResp.Data

		// Patch the role
		roleReq.Operation = logical.PatchOperation
		roleReq.Data[testCase.Field] = testCase.Patched
		resp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (resp != nil && resp.IsError()) {
			t.Fatalf("bad [%d/%v] patch: err: %v resp: %#v", index, testCase.Field, err, resp)
		}

		// Re-read and verify the role
		roleReq.Operation = logical.ReadOperation
		afterRoleDataResp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (afterRoleDataResp != nil && afterRoleDataResp.IsError()) {
			t.Fatalf("bad [%d/%v] read: err: %v resp: %#v", index, testCase.Field, err, resp)
		}

		afterRoleData := afterRoleDataResp.Data

		for field, before := range beforeRoleData {
			switch typed := before.(type) {
			case *bool:
				before = *typed
				afterRoleData[field] = *(afterRoleData[field].(*bool))
			}

			if field != testCase.Field {
				require.Equal(t, before, afterRoleData[field], fmt.Sprintf("bad [%d/%v] compare: non-modified field %v should not be changed", index, testCase.Field, field))
			} else {
				require.Equal(t, before, testCase.Before, fmt.Sprintf("bad [%d] compare: modified field %v before should be correct", index, field))
				require.Equal(t, afterRoleData[field], testCase.Patched, fmt.Sprintf("bad [%d] compare: modified field %v after should be correct", index, field))
			}
		}
	}
}

func TestPKI_RolePolicyInformation_Flat(t *testing.T) {
	t.Parallel()
	type TestCase struct {
		Input   interface{}
		ASN     interface{}
		OidList []string
	}

	expectedSimpleAsnExtension := "MBYwCQYHKwYBBAEBATAJBgcrBgEEAQEC"
	expectedSimpleOidList := append(*new([]string), "1.3.6.1.4.1.1.1", "1.3.6.1.4.1.1.2")

	testCases := []TestCase{
		{
			Input:   "1.3.6.1.4.1.1.1,1.3.6.1.4.1.1.2",
			ASN:     expectedSimpleAsnExtension,
			OidList: expectedSimpleOidList,
		},
		{
			Input:   "[{\"oid\":\"1.3.6.1.4.1.1.1\"},{\"oid\":\"1.3.6.1.4.1.1.2\"}]",
			ASN:     expectedSimpleAsnExtension,
			OidList: expectedSimpleOidList,
		},
		{
			Input:   "[{\"oid\":\"1.3.6.1.4.1.7.8\",\"notice\":\"I am a user Notice\"},{\"oid\":\"1.3.6.1.44947.1.2.4\",\"cps\":\"https://example.com\"}]",
			ASN:     "MF8wLQYHKwYBBAEHCDAiMCAGCCsGAQUFBwICMBQMEkkgYW0gYSB1c2VyIE5vdGljZTAuBgkrBgGC3xMBAgQwITAfBggrBgEFBQcCARYTaHR0cHM6Ly9leGFtcGxlLmNvbQ==",
			OidList: append(*new([]string), "1.3.6.1.4.1.7.8", "1.3.6.1.44947.1.2.4"),
		},
	}

	b, storage := CreateBackendWithStorage(t)

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
	caResp, err := b.HandleRequest(context.Background(), caReq)
	if err != nil || (caResp != nil && caResp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, caResp)
	}

	for index, testCase := range testCases {
		var roleResp *logical.Response
		var issueResp *logical.Response
		var err error

		// Create/update the role
		roleData := map[string]interface{}{}
		roleData[policyIdentifiersParam] = testCase.Input

		roleReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "roles/testrole",
			Storage:   storage,
			Data:      roleData,
		}

		roleResp, err = b.HandleRequest(context.Background(), roleReq)
		if err != nil || (roleResp != nil && roleResp.IsError()) {
			t.Fatalf("bad [%d], setting policy identifier %v err: %v resp: %#v", index, testCase.Input, err, roleResp)
		}

		// Issue Using this role
		issueData := map[string]interface{}{}
		issueData["common_name"] = "localhost"
		issueData["ttl"] = "2s"

		issueReq := &logical.Request{
			Operation: logical.UpdateOperation,
			Path:      "issue/testrole",
			Storage:   storage,
			Data:      issueData,
		}

		issueResp, err = b.HandleRequest(context.Background(), issueReq)
		if err != nil || (issueResp != nil && issueResp.IsError()) {
			t.Fatalf("bad [%d], setting policy identifier %v err: %v resp: %#v", index, testCase.Input, err, issueResp)
		}

		// Validate the OIDs
		policyIdentifiers, err := getPolicyIdentifiersOffCertificate(*issueResp)
		if err != nil {
			t.Fatalf("bad [%d], getting policy identifier from %v err: %v resp: %#v", index, testCase.Input, err, issueResp)
		}
		if len(policyIdentifiers) != len(testCase.OidList) {
			t.Fatalf("bad [%d], wrong certificate policy identifier from %v len expected: %d got %d", index, testCase.Input, len(testCase.OidList), len(policyIdentifiers))
		}
		for i, identifier := range policyIdentifiers {
			if identifier != testCase.OidList[i] {
				t.Fatalf("bad [%d], wrong certificate policy identifier from %v expected: %v got %v", index, testCase.Input, testCase.OidList[i], policyIdentifiers[i])
			}
		}
		// Validate the ASN
		certificateAsn, err := getPolicyInformationExtensionOffCertificate(*issueResp)
		if err != nil {
			t.Fatalf("bad [%d], getting extension from %v err: %v resp: %#v", index, testCase.Input, err, issueResp)
		}
		certificateB64 := make([]byte, len(certificateAsn)*2)
		base64.StdEncoding.Encode(certificateB64, certificateAsn)
		certificateString := string(certificateB64[:])
		assert.Contains(t, certificateString, testCase.ASN)
	}
}

func getPolicyIdentifiersOffCertificate(resp logical.Response) ([]string, error) {
	stringCertificate := resp.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(stringCertificate))
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	policyIdentifierStrings := make([]string, len(certificate.PolicyIdentifiers))
	for index, asnOid := range certificate.PolicyIdentifiers {
		policyIdentifierStrings[index] = asnOid.String()
	}
	return policyIdentifierStrings, nil
}

func getPolicyInformationExtensionOffCertificate(resp logical.Response) ([]byte, error) {
	stringCertificate := resp.Data["certificate"].(string)
	block, _ := pem.Decode([]byte(stringCertificate))
	certificate, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	for _, extension := range certificate.Extensions {
		if extension.Id.Equal(asn1.ObjectIdentifier{2, 5, 29, 32}) {
			return extension.Value, nil
		}
	}
	return *new([]byte), errors.New("No Policy Information Extension Found")
}
