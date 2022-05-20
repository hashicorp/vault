package pki

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/require"
)

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

	// Update values due to switching of ttl type
	resp.Data["ttl_duration"] = resp.Data["ttl"]
	resp.Data["ttl"] = (time.Duration(resp.Data["ttl"].(int64)) * time.Second).String()
	resp.Data["max_ttl_duration"] = resp.Data["max_ttl"]
	resp.Data["max_ttl"] = (time.Duration(resp.Data["max_ttl"].(int64)) * time.Second).String()
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
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	roleReq.Operation = logical.ReadOperation
	resp, err = b.HandleRequest(context.Background(), roleReq)
	if err != nil || (resp != nil && resp.IsError()) {
		t.Fatalf("bad: err: %v resp: %#v", err, resp)
	}

	keyUsage := resp.Data["key_usage"].([]string)
	if len(keyUsage) != 2 {
		t.Fatalf("key_usage should have 2 values")
	}

	// Update values due to switching of ttl type
	resp.Data["ttl_duration"] = resp.Data["ttl"]
	resp.Data["ttl"] = (time.Duration(resp.Data["ttl"].(int64)) * time.Second).String()
	resp.Data["max_ttl_duration"] = resp.Data["max_ttl"]
	resp.Data["max_ttl"] = (time.Duration(resp.Data["max_ttl"].(int64)) * time.Second).String()
	// Check that old key usage value is nil
	var role roleEntry
	err = mapstructure.Decode(resp.Data, &role)
	if err != nil {
		t.Fatal(err)
	}
	if role.KeyUsageOld != "" {
		t.Fatalf("old key usage storage value should be blank")
	}

	// Make it explicit
	role.KeyUsageOld = "KeyEncipherment,DigitalSignature"
	role.KeyUsage = nil

	entry, err := logical.StorageEntryJSON("role/testrole", role)
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
	var result roleEntry
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
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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

	// Update values due to switching of ttl type
	resp.Data["ttl_duration"] = resp.Data["ttl"]
	resp.Data["ttl"] = (time.Duration(resp.Data["ttl"].(int64)) * time.Second).String()
	resp.Data["max_ttl_duration"] = resp.Data["max_ttl"]
	resp.Data["max_ttl"] = (time.Duration(resp.Data["max_ttl"].(int64)) * time.Second).String()
	// Check that old key usage value is nil
	var role roleEntry
	err = mapstructure.Decode(resp.Data, &role)
	if err != nil {
		t.Fatal(err)
	}
	if role.OUOld != "" {
		t.Fatalf("old ou storage value should be blank")
	}
	if role.OrganizationOld != "" {
		t.Fatalf("old organization storage value should be blank")
	}

	// Make it explicit
	role.OUOld = "abc,123"
	role.OU = nil
	role.OrganizationOld = "org1,org2"
	role.Organization = nil

	entry, err := logical.StorageEntryJSON("role/testrole", role)
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
	var result roleEntry
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
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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

	// Update values due to switching of ttl type
	resp.Data["ttl_duration"] = resp.Data["ttl"]
	resp.Data["ttl"] = (time.Duration(resp.Data["ttl"].(int64)) * time.Second).String()
	resp.Data["max_ttl_duration"] = resp.Data["max_ttl"]
	resp.Data["max_ttl"] = (time.Duration(resp.Data["max_ttl"].(int64)) * time.Second).String()
	// Check that old key usage value is nil
	var role roleEntry
	err = mapstructure.Decode(resp.Data, &role)
	if err != nil {
		t.Fatal(err)
	}
	if role.AllowedDomainsOld != "" {
		t.Fatalf("old allowed_domains storage value should be blank")
	}

	// Make it explicit
	role.AllowedDomainsOld = "foobar.com,*example.com"
	role.AllowedDomains = nil

	entry, err := logical.StorageEntryJSON("role/testrole", role)
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
	var result roleEntry
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
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
	var resp *logical.Response
	var err error
	b, storage := createBackendWithStorage(t)

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
			Before:  []string{"1.2.3.4.5"},
			Patched: []string{"5.4.3.2.1"},
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

	b, storage := createBackendWithStorage(t)

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
