// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"net"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/go-test/deep"
	"github.com/hashicorp/go-secure-stdlib/parseutil"
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/stretchr/testify/require"
)

func TestPki_FetchCertBySerial(t *testing.T) {
	t.Parallel()
	b, storage := CreateBackendWithStorage(t)
	sc := b.makeStorageContext(ctx, storage)

	cases := map[string]struct {
		Req    *logical.Request
		Prefix string
		Serial string
	}{
		"valid cert": {
			&logical.Request{
				Storage: storage,
			},
			issuing.PathCerts,
			"00:00:00:00:00:00:00:00",
		},
		"revoked cert": {
			&logical.Request{
				Storage: storage,
			},
			"revoked/",
			"11:11:11:11:11:11:11:11",
		},
	}

	// Test for colon-based paths in storage
	for name, tc := range cases {
		storageKey := fmt.Sprintf("%s%s", tc.Prefix, tc.Serial)
		err := storage.Put(context.Background(), &logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s colon-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(sc, tc.Prefix, tc.Serial)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path: %s", name, err)
		}

		// Check for non-nil on valid/revoked certs
		if certEntry == nil {
			t.Fatalf("nil on %s for colon-based storage path", name)
		}

		// Ensure that cert serials are converted/updated after fetch
		expectedKey := tc.Prefix + normalizeSerial(tc.Serial)
		se, err := storage.Get(context.Background(), expectedKey)
		if err != nil {
			t.Fatalf("error on %s for colon-based storage path:%s", name, err)
		}
		if strings.Compare(expectedKey, se.Key) != 0 {
			t.Fatalf("expected: %s, got: %s", expectedKey, certEntry.Key)
		}
	}

	// Reset storage
	storage = &logical.InmemStorage{}

	// Test for hyphen-base paths in storage
	for name, tc := range cases {
		storageKey := tc.Prefix + normalizeSerial(tc.Serial)
		err := storage.Put(context.Background(), &logical.StorageEntry{
			Key:   storageKey,
			Value: []byte("some data"),
		})
		if err != nil {
			t.Fatalf("error writing to storage on %s hyphen-based storage path: %s", name, err)
		}

		certEntry, err := fetchCertBySerial(sc, tc.Prefix, tc.Serial)
		if err != nil || certEntry == nil {
			t.Fatalf("error on %s for hyphen-based storage path: err: %v, entry: %v", name, err, certEntry)
		}
	}
}

// Demonstrate that multiple OUs in the name are handled in an
// order-preserving way.
func TestPki_MultipleOUs(t *testing.T) {
	t.Parallel()
	b, _ := CreateBackendWithStorage(t)
	fields := addCACommonFields(map[string]*framework.FieldSchema{})

	apiData := &framework.FieldData{
		Schema: fields,
		Raw: map[string]interface{}{
			"cn":  "example.com",
			"ttl": 3600,
		},
	}
	input := &inputBundle{
		apiData: apiData,
		role: &issuing.RoleEntry{
			MaxTTL: 3600,
			OU:     []string{"Z", "E", "V"},
		},
	}
	cb, _, err := generateCreationBundle(b.System(), input, nil, nil)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	expected := []string{"Z", "E", "V"}
	actual := cb.Params.Subject.OrganizationalUnit

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v, got %v", expected, actual)
	}
}

func TestPki_PermitFQDNs(t *testing.T) {
	t.Parallel()
	b, _ := CreateBackendWithStorage(t)
	fields := addCACommonFields(map[string]*framework.FieldSchema{})

	cases := map[string]struct {
		input            *inputBundle
		expectedDnsNames []string
		expectedEmails   []string
	}{
		"base valid case": {
			input: &inputBundle{
				apiData: &framework.FieldData{
					Schema: fields,
					Raw: map[string]interface{}{
						"common_name": "example.com.",
						"ttl":         3600,
					},
				},
				role: &issuing.RoleEntry{
					AllowAnyName:     true,
					MaxTTL:           3600,
					EnforceHostnames: true,
				},
			},
			expectedDnsNames: []string{"example.com."},
			expectedEmails:   []string{},
		},
		"case insensitivity validation": {
			input: &inputBundle{
				apiData: &framework.FieldData{
					Schema: fields,
					Raw: map[string]interface{}{
						"common_name": "Example.Net",
						"alt_names":   "eXaMPLe.COM",
						"ttl":         3600,
					},
				},
				role: &issuing.RoleEntry{
					AllowedDomains:   []string{"example.net", "EXAMPLE.COM"},
					AllowBareDomains: true,
					MaxTTL:           3600,
				},
			},
			expectedDnsNames: []string{"Example.Net", "eXaMPLe.COM"},
			expectedEmails:   []string{},
		},
		"case insensitivity subdomain validation": {
			input: &inputBundle{
				apiData: &framework.FieldData{
					Schema: fields,
					Raw: map[string]interface{}{
						"common_name": "SUB.EXAMPLE.COM",
						"ttl":         3600,
					},
				},
				role: &issuing.RoleEntry{
					AllowedDomains:   []string{"example.com", "*.Example.com"},
					AllowGlobDomains: true,
					MaxTTL:           3600,
				},
			},
			expectedDnsNames: []string{"SUB.EXAMPLE.COM"},
			expectedEmails:   []string{},
		},
		"case email as AllowedDomain with bare domains": {
			input: &inputBundle{
				apiData: &framework.FieldData{
					Schema: fields,
					Raw: map[string]interface{}{
						"common_name": "test@testemail.com",
						"ttl":         3600,
					},
				},
				role: &issuing.RoleEntry{
					AllowedDomains:   []string{"test@testemail.com"},
					AllowBareDomains: true,
					MaxTTL:           3600,
				},
			},
			expectedDnsNames: []string{},
			expectedEmails:   []string{"test@testemail.com"},
		},
		"case email common name with bare domains": {
			input: &inputBundle{
				apiData: &framework.FieldData{
					Schema: fields,
					Raw: map[string]interface{}{
						"common_name": "test@testemail.com",
						"ttl":         3600,
					},
				},
				role: &issuing.RoleEntry{
					AllowedDomains:   []string{"testemail.com"},
					AllowBareDomains: true,
					MaxTTL:           3600,
				},
			},
			expectedDnsNames: []string{},
			expectedEmails:   []string{"test@testemail.com"},
		},
	}

	for name, testCase := range cases {
		name := name
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			cb, _, err := generateCreationBundle(b.System(), testCase.input, nil, nil)
			if err != nil {
				t.Fatalf("Error: %v", err)
			}

			actualDnsNames := cb.Params.DNSNames

			if !reflect.DeepEqual(testCase.expectedDnsNames, actualDnsNames) {
				t.Fatalf("Expected dns names %v, got %v", testCase.expectedDnsNames, actualDnsNames)
			}

			actualEmails := cb.Params.EmailAddresses

			if !reflect.DeepEqual(testCase.expectedEmails, actualEmails) {
				t.Fatalf("Expected email addresses %v, got %v", testCase.expectedEmails, actualEmails)
			}
		})
	}
}

type parseCertificateTestCase struct {
	name       string
	data       map[string]interface{}
	roleData   map[string]interface{} // if a role is to be created
	ttl        time.Duration
	wantParams certutil.CreationParameters
	wantFields map[string]interface{}
	wantErr    bool
}

// TestDisableVerifyCertificateEnvVar verifies that env var VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION
// can be used to disable cert verification.
func TestDisableVerifyCertificateEnvVar(t *testing.T) {
	caData := map[string]any{
		// Copied from the "full CA" test case of TestParseCertificate,
		// with tweaked permitted_dns_domains and ttl
		"common_name":           "the common name",
		"alt_names":             "user@example.com,admin@example.com,example.com,www.example.com",
		"ip_sans":               "1.2.3.4,1.2.3.5",
		"uri_sans":              "https://example.com,https://www.example.com",
		"other_sans":            "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
		"ttl":                   "3h",
		"max_path_length":       2,
		"permitted_dns_domains": ".example.com,.www.example.com",
		"ou":                    "unit1, unit2",
		"organization":          "org1, org2",
		"country":               "US, CA",
		"locality":              "locality1, locality2",
		"province":              "province1, province2",
		"street_address":        "street_address1, street_address2",
		"postal_code":           "postal_code1, postal_code2",
		"not_before_duration":   "45s",
		"key_type":              "rsa",
		"use_pss":               true,
		"key_bits":              2048,
		"signature_bits":        384,
	}

	roleData := map[string]any{
		"allow_any_name":      true,
		"cn_validations":      "disabled",
		"allow_ip_sans":       true,
		"allowed_other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:*@example.com",
		"allowed_uri_sans":    "https://example.com,https://www.example.com",
		"allowed_user_ids":    "*",
		"not_before_duration": "45s",
		"signature_bits":      384,
		"key_usage":           "KeyAgreement",
		"ext_key_usage":       "ServerAuth",
		"ext_key_usage_oids":  "1.3.6.1.5.5.7.3.67,1.3.6.1.5.5.7.3.68",
		"client_flag":         false,
		"server_flag":         false,
		"policy_identifiers":  "1.2.3.4.5.6.7.8.9.0",
	}

	certData := map[string]any{
		// using the same order as in https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-certificate-and-key
		"common_name": "the common name non ca",
		"alt_names":   "user@example.com,admin@example.com,example.com,www.example.com",
		"ip_sans":     "1.2.3.4,1.2.3.5",
		"uri_sans":    "https://example.com,https://www.example.com",
		"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
		"ttl":         "2h",
		// format
		// private_key_format
		"exclude_cn_from_sans": true,
		// not_after
		// remove_roots_from_chain
		"user_ids": "humanoid,robot",
	}

	defer func() {
		os.Unsetenv("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION")
	}()

	b, s := CreateBackendWithStorage(t)

	// Create the CA
	resp, err := CBWrite(b, s, "root/generate/internal", caData)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Create the role
	resp, err = CBWrite(b, s, "roles/test", roleData)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Try to create the cert -- should fail verification, since example.com is not allowed
	t.Run("no VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION env var", func(t *testing.T) {
		resp, err = CBWrite(b, s, "issue/test", certData)
		require.ErrorContains(t, err, `DNS name "example.com" is not permitted by any constraint`)
	})

	// Try to create the cert -- should fail verification, since example.com is not allowed
	t.Run("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION=false", func(t *testing.T) {
		os.Setenv("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION", "false")
		resp, err = CBWrite(b, s, "issue/test", certData)
		require.ErrorContains(t, err, `DNS name "example.com" is not permitted by any constraint`)
	})

	// Create the cert, should succeed with the disable env var set
	t.Run("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION=true", func(t *testing.T) {
		os.Setenv("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION", "true")
		resp, err = CBWrite(b, s, "issue/test", certData)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})

	// Invalid env var
	t.Run("invalid VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION", func(t *testing.T) {
		os.Setenv("VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION", "invalid")
		resp, err = CBWrite(b, s, "issue/test", certData)
		require.ErrorContains(t, err, "failed parsing environment variable VAULT_DISABLE_PKI_CONSTRAINTS_VERIFICATION")
	})
}

func TestParseCertificate(t *testing.T) {
	t.Parallel()

	parseURL := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}

	tests := []*parseCertificateTestCase{
		{
			name: "simple CA",
			data: map[string]interface{}{
				"common_name":         "the common name",
				"key_type":            "ec",
				"key_bits":            384,
				"ttl":                 "1h",
				"not_before_duration": "30s",
				"street_address":      "",
			},
			ttl: 1 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName: "the common name",
				},
				DNSNames:                      nil,
				EmailAddresses:                nil,
				IPAddresses:                   nil,
				URIs:                          nil,
				OtherSANs:                     make(map[string][]string),
				IsCA:                          true,
				KeyType:                       "ec",
				KeyBits:                       384,
				NotAfter:                      time.Time{},
				KeyUsage:                      x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
				ExtKeyUsage:                   0,
				ExtKeyUsageOIDs:               nil,
				PolicyIdentifiers:             nil,
				BasicConstraintsValidForNonCA: false,
				SignatureBits:                 384,
				UsePSS:                        false,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 -1,
				NotBeforeDuration:             30,
				SKID:                          []byte("We'll assert that it is not nil as an special case"),
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name",
				"alt_names":             "",
				"ip_sans":               "",
				"uri_sans":              "",
				"other_sans":            "",
				"signature_bits":        384,
				"exclude_cn_from_sans":  true,
				"ou":                    "",
				"organization":          "",
				"country":               "",
				"locality":              "",
				"province":              "",
				"street_address":        "",
				"postal_code":           "",
				"serial_number":         "",
				"ttl":                   "1h0m30s",
				"max_path_length":       -1,
				"permitted_dns_domains": "",
				"use_pss":               false,
				"key_type":              "ec",
				"key_bits":              384,
				"skid":                  "We'll assert that it is not nil as an special case",
			},
			wantErr: false,
		},
		{
			// Note that this test's data is used to create the internal CA used by test "full non CA cert"
			name: "full CA",
			data: map[string]interface{}{
				// using the same order as in https://developer.hashicorp.com/vault/api-docs/secret/pki#sign-certificate
				"common_name":           "the common name",
				"alt_names":             "user@example.com,admin@example.com,example.com,www.example.com",
				"ip_sans":               "1.2.3.4,1.2.3.5",
				"uri_sans":              "https://example.com,https://www.example.com",
				"other_sans":            "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
				"ttl":                   "2h",
				"max_path_length":       2,
				"permitted_dns_domains": "example.com,.example.com,.www.example.com",
				"ou":                    "unit1, unit2",
				"organization":          "org1, org2",
				"country":               "US, CA",
				"locality":              "locality1, locality2",
				"province":              "province1, province2",
				"street_address":        "street_address1, street_address2",
				"postal_code":           "postal_code1, postal_code2",
				"not_before_duration":   "45s",
				"key_type":              "rsa",
				"use_pss":               true,
				"key_bits":              2048,
				"signature_bits":        384,
				// TODO(kitography): Specify key usage
			},
			ttl: 2 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName:         "the common name",
					OrganizationalUnit: []string{"unit1", "unit2"},
					Organization:       []string{"org1", "org2"},
					Country:            []string{"CA", "US"},
					Locality:           []string{"locality1", "locality2"},
					Province:           []string{"province1", "province2"},
					StreetAddress:      []string{"street_address1", "street_address2"},
					PostalCode:         []string{"postal_code1", "postal_code2"},
				},
				DNSNames:                      []string{"example.com", "www.example.com"},
				EmailAddresses:                []string{"admin@example.com", "user@example.com"},
				IPAddresses:                   []net.IP{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 5}},
				URIs:                          []*url.URL{parseURL("https://example.com"), parseURL("https://www.example.com")},
				OtherSANs:                     map[string][]string{"1.3.6.1.4.1.311.20.2.3": {"caadmin@example.com"}},
				IsCA:                          true,
				KeyType:                       "rsa",
				KeyBits:                       2048,
				NotAfter:                      time.Time{},
				KeyUsage:                      x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
				ExtKeyUsage:                   0,
				ExtKeyUsageOIDs:               nil,
				PolicyIdentifiers:             nil,
				BasicConstraintsValidForNonCA: false,
				SignatureBits:                 384,
				UsePSS:                        true,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           []string{"example.com", ".example.com", ".www.example.com"},
				URLs:                          nil,
				MaxPathLength:                 2,
				NotBeforeDuration:             45 * time.Second,
				SKID:                          []byte("We'll assert that it is not nil as an special case"),
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name",
				"alt_names":             "example.com,www.example.com,admin@example.com,user@example.com",
				"ip_sans":               "1.2.3.4,1.2.3.5",
				"uri_sans":              "https://example.com,https://www.example.com",
				"other_sans":            "1.3.6.1.4.1.311.20.2.3;UTF-8:caadmin@example.com",
				"signature_bits":        384,
				"exclude_cn_from_sans":  true,
				"ou":                    "unit1,unit2",
				"organization":          "org1,org2",
				"country":               "CA,US",
				"locality":              "locality1,locality2",
				"province":              "province1,province2",
				"street_address":        "street_address1,street_address2",
				"postal_code":           "postal_code1,postal_code2",
				"serial_number":         "",
				"ttl":                   "2h0m45s",
				"max_path_length":       2,
				"permitted_dns_domains": "example.com,.example.com,.www.example.com",
				"use_pss":               true,
				"key_type":              "rsa",
				"key_bits":              2048,
				"skid":                  "We'll assert that it is not nil as an special case",
			},
			wantErr: false,
		},
		{
			// Note that we use the data of test "full CA" to create the internal CA needed for this test
			name: "full non CA cert",
			data: map[string]interface{}{
				// using the same order as in https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-certificate-and-key
				"common_name": "the common name non ca",
				"alt_names":   "user@example.com,admin@example.com,example.com,www.example.com",
				"ip_sans":     "1.2.3.4,1.2.3.5",
				"uri_sans":    "https://example.com,https://www.example.com",
				"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
				"ttl":         "2h",
				// format
				// private_key_format
				"exclude_cn_from_sans": true,
				// not_after
				// remove_roots_from_chain
				"user_ids": "humanoid,robot",
			},
			roleData: map[string]interface{}{
				"allow_any_name":      true,
				"cn_validations":      "disabled",
				"allow_ip_sans":       true,
				"allowed_other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:*@example.com",
				"allowed_uri_sans":    "https://example.com,https://www.example.com",
				"allowed_user_ids":    "*",
				"not_before_duration": "45s",
				"signature_bits":      384,
				"key_usage":           "KeyAgreement",
				"ext_key_usage":       "ServerAuth",
				"ext_key_usage_oids":  "1.3.6.1.5.5.7.3.67,1.3.6.1.5.5.7.3.68",
				"client_flag":         false,
				"server_flag":         false,
				"policy_identifiers":  "1.2.3.4.5.6.7.8.9.0",
			},
			ttl: 2 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName: "the common name non ca",
				},
				DNSNames:                      []string{"example.com", "www.example.com"},
				EmailAddresses:                []string{"admin@example.com", "user@example.com"},
				IPAddresses:                   []net.IP{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 5}},
				URIs:                          []*url.URL{parseURL("https://example.com"), parseURL("https://www.example.com")},
				OtherSANs:                     map[string][]string{"1.3.6.1.4.1.311.20.2.3": {"caadmin@example.com"}},
				IsCA:                          false,
				KeyType:                       "rsa",
				KeyBits:                       2048,
				NotAfter:                      time.Time{},
				KeyUsage:                      x509.KeyUsageKeyAgreement,
				ExtKeyUsage:                   0, // Please Ignore
				ExtKeyUsageOIDs:               []string{"1.3.6.1.5.5.7.3.1", "1.3.6.1.5.5.7.3.67", "1.3.6.1.5.5.7.3.68"},
				PolicyIdentifiers:             []string{"1.2.3.4.5.6.7.8.9.0"},
				BasicConstraintsValidForNonCA: false,
				SignatureBits:                 384,
				UsePSS:                        false,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 0,
				NotBeforeDuration:             45,
				SKID:                          []byte("We'll assert that it is not nil as an special case"),
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name non ca",
				"alt_names":             "example.com,www.example.com,admin@example.com,user@example.com",
				"ip_sans":               "1.2.3.4,1.2.3.5",
				"uri_sans":              "https://example.com,https://www.example.com",
				"other_sans":            "1.3.6.1.4.1.311.20.2.3;UTF-8:caadmin@example.com",
				"signature_bits":        384,
				"exclude_cn_from_sans":  true,
				"ou":                    "",
				"organization":          "",
				"country":               "",
				"locality":              "",
				"province":              "",
				"street_address":        "",
				"postal_code":           "",
				"serial_number":         "",
				"ttl":                   "2h0m45s",
				"max_path_length":       0,
				"permitted_dns_domains": "",
				"use_pss":               false,
				"key_type":              "rsa",
				"key_bits":              2048,
				"skid":                  "We'll assert that it is not nil as an special case",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, s := CreateBackendWithStorage(t)

			var cert *x509.Certificate
			issueTime := time.Now()
			if tt.wantParams.IsCA {
				resp, err := CBWrite(b, s, "root/generate/internal", tt.data)
				require.NoError(t, err)
				require.NotNil(t, resp)

				certData := resp.Data["certificate"].(string)
				cert, err = parsing.ParseCertificateFromString(certData)
				require.NoError(t, err)
				require.NotNil(t, cert)
			} else {
				// use the "simple CA" data to create the internal CA
				caData := tests[1].data
				caData["ttl"] = "3h"
				resp, err := CBWrite(b, s, "root/generate/internal", caData)
				require.NoError(t, err)
				require.NotNil(t, resp)

				// create a role
				resp, err = CBWrite(b, s, "roles/test", tt.roleData)
				require.NoError(t, err)
				require.NotNil(t, resp)

				// create the cert
				resp, err = CBWrite(b, s, "issue/test", tt.data)
				require.NoError(t, err)
				require.NotNil(t, resp)

				certData := resp.Data["certificate"].(string)
				cert, err = parsing.ParseCertificateFromString(certData)
				require.NoError(t, err)
				require.NotNil(t, cert)
			}

			t.Run(tt.name+" parameters", func(t *testing.T) {
				testParseCertificateToCreationParameters(t, issueTime, tt, cert)
			})
			t.Run(tt.name+" fields", func(t *testing.T) {
				testParseCertificateToFields(t, issueTime, tt, cert)
			})
		})
	}
}

func testParseCertificateToCreationParameters(t *testing.T, issueTime time.Time, tt *parseCertificateTestCase, cert *x509.Certificate) {
	params, err := certutil.ParseCertificateToCreationParameters(*cert)

	if tt.wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		ignoreBasicConstraintsValidForNonCA := tt.wantParams.IsCA

		var diff []string
		for _, d := range deep.Equal(tt.wantParams, params) {
			switch {
			case strings.HasPrefix(d, "SKID"):
				continue
			case strings.HasPrefix(d, "BasicConstraintsValidForNonCA") && ignoreBasicConstraintsValidForNonCA:
				continue
			case strings.HasPrefix(d, "NotBeforeDuration"):
				continue
			case strings.HasPrefix(d, "NotAfter"):
				continue
			}
			diff = append(diff, d)
		}
		if diff != nil {
			t.Errorf("testParseCertificateToCreationParameters() diff: %s", strings.Join(diff, "\n"))
		}

		require.NotNil(t, params.SKID)
		require.GreaterOrEqual(t, params.NotBeforeDuration, tt.wantParams.NotBeforeDuration,
			"NotBeforeDuration want: %s got: %s", tt.wantParams.NotBeforeDuration, params.NotBeforeDuration)

		require.GreaterOrEqual(t, params.NotAfter, issueTime.Add(tt.ttl).Add(-1*time.Minute),
			"NotAfter want: %s got: %s", tt.wantParams.NotAfter, params.NotAfter)
		require.LessOrEqual(t, params.NotAfter, issueTime.Add(tt.ttl).Add(1*time.Minute),
			"NotAfter want: %s got: %s", tt.wantParams.NotAfter, params.NotAfter)
	}
}

func testParseCertificateToFields(t *testing.T, issueTime time.Time, tt *parseCertificateTestCase, cert *x509.Certificate) {
	fields, err := certutil.ParseCertificateToFields(*cert)
	if tt.wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		require.NotNil(t, fields["skid"])
		delete(fields, "skid")
		delete(tt.wantFields, "skid")

		{
			// Sometimes TTL comes back as 1s off, so we'll allow that
			expectedTTL, err := parseutil.ParseDurationSecond(tt.wantFields["ttl"].(string))
			require.NoError(t, err)
			actualTTL, err := parseutil.ParseDurationSecond(fields["ttl"].(string))
			require.NoError(t, err)

			diff := expectedTTL - actualTTL
			require.LessOrEqual(t, actualTTL, expectedTTL, // NotAfter is generated before NotBefore so the time.Now of notBefore may be later, shrinking our calculated TTL during very slow tests
				"ttl should be, if off, smaller than expected want: %s got: %s", tt.wantFields["ttl"], fields["ttl"])
			require.LessOrEqual(t, diff, 30*time.Second, // Test can be slow, allow more off in the other direction
				"ttl must be at most 30s off, want: %s got: %s", tt.wantFields["ttl"], fields["ttl"])
			delete(fields, "ttl")
			delete(tt.wantFields, "ttl")
		}

		if diff := deep.Equal(tt.wantFields, fields); diff != nil {
			t.Errorf("testParseCertificateToFields() diff: %s", strings.ReplaceAll(strings.Join(diff, "\n"), "map", "\nmap"))
		}
	}
}

func TestParseCsr(t *testing.T) {
	t.Parallel()

	parseURL := func(s string) *url.URL {
		u, err := url.Parse(s)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}

	tests := []*parseCertificateTestCase{
		{
			name: "simple CSR",
			data: map[string]interface{}{
				"common_name":         "the common name",
				"key_type":            "ec",
				"key_bits":            384,
				"ttl":                 "1h",
				"not_before_duration": "30s",
				"street_address":      "",
			},
			ttl: 1 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName: "the common name",
				},
				DNSNames:                      nil,
				EmailAddresses:                nil,
				IPAddresses:                   nil,
				URIs:                          nil,
				OtherSANs:                     make(map[string][]string),
				IsCA:                          false,
				KeyType:                       "ec",
				KeyBits:                       384,
				NotAfter:                      time.Time{},
				KeyUsage:                      0,
				ExtKeyUsage:                   0,
				ExtKeyUsageOIDs:               nil,
				PolicyIdentifiers:             nil,
				BasicConstraintsValidForNonCA: false,
				SignatureBits:                 384,
				UsePSS:                        false,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 0,
				NotBeforeDuration:             0,
				SKID:                          nil,
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name",
				"ou":                    "",
				"organization":          "",
				"country":               "",
				"locality":              "",
				"province":              "",
				"street_address":        "",
				"postal_code":           "",
				"alt_names":             "",
				"ip_sans":               "",
				"uri_sans":              "",
				"other_sans":            "",
				"exclude_cn_from_sans":  true,
				"key_type":              "ec",
				"key_bits":              384,
				"signature_bits":        384,
				"use_pss":               false,
				"serial_number":         "",
				"add_basic_constraints": false,
			},
			wantErr: false,
		},
		{
			name: "full CSR with basic constraints",
			data: map[string]interface{}{
				// using the same order as in https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-intermediate-csr
				"common_name": "the common name",
				"alt_names":   "user@example.com,admin@example.com,example.com,www.example.com",
				"ip_sans":     "1.2.3.4,1.2.3.5",
				"uri_sans":    "https://example.com,https://www.example.com",
				"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
				// format
				// private_key_format
				"key_type": "rsa",
				"key_bits": 2048,
				"key_name": "the-key-name",
				// key_ref
				"signature_bits": 384,
				// exclude_cn_from_sans
				"ou":                    "unit1, unit2",
				"organization":          "org1, org2",
				"country":               "US, CA",
				"locality":              "locality1, locality2",
				"province":              "province1, province2",
				"street_address":        "street_address1, street_address2",
				"postal_code":           "postal_code1, postal_code2",
				"serial_number":         "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				"add_basic_constraints": true,
			},
			ttl: 2 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName:         "the common name",
					OrganizationalUnit: []string{"unit1", "unit2"},
					Organization:       []string{"org1", "org2"},
					Country:            []string{"CA", "US"},
					Locality:           []string{"locality1", "locality2"},
					Province:           []string{"province1", "province2"},
					StreetAddress:      []string{"street_address1", "street_address2"},
					PostalCode:         []string{"postal_code1", "postal_code2"},
					SerialNumber:       "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				},
				DNSNames:                      []string{"example.com", "www.example.com"},
				EmailAddresses:                []string{"admin@example.com", "user@example.com"},
				IPAddresses:                   []net.IP{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 5}},
				URIs:                          []*url.URL{parseURL("https://example.com"), parseURL("https://www.example.com")},
				OtherSANs:                     map[string][]string{"1.3.6.1.4.1.311.20.2.3": {"caadmin@example.com"}},
				IsCA:                          true,
				KeyType:                       "rsa",
				KeyBits:                       2048,
				NotAfter:                      time.Time{},
				KeyUsage:                      0,   // TODO(kitography): Verify with Kit
				ExtKeyUsage:                   0,   // TODO(kitography): Verify with Kit
				ExtKeyUsageOIDs:               nil, // TODO(kitography): Verify with Kit
				PolicyIdentifiers:             nil, // TODO(kitography): Verify with Kit
				BasicConstraintsValidForNonCA: true,
				SignatureBits:                 384,
				UsePSS:                        false,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 -1,
				NotBeforeDuration:             0,
				SKID:                          nil,
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name",
				"ou":                    "unit1,unit2",
				"organization":          "org1,org2",
				"country":               "CA,US",
				"locality":              "locality1,locality2",
				"province":              "province1,province2",
				"street_address":        "street_address1,street_address2",
				"postal_code":           "postal_code1,postal_code2",
				"alt_names":             "example.com,www.example.com,admin@example.com,user@example.com",
				"ip_sans":               "1.2.3.4,1.2.3.5",
				"uri_sans":              "https://example.com,https://www.example.com",
				"other_sans":            "1.3.6.1.4.1.311.20.2.3;UTF-8:caadmin@example.com",
				"exclude_cn_from_sans":  true,
				"key_type":              "rsa",
				"key_bits":              2048,
				"signature_bits":        384,
				"use_pss":               false,
				"serial_number":         "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				"add_basic_constraints": true,
			},
			wantErr: false,
		},
		{
			name: "full CSR without basic constraints",
			data: map[string]interface{}{
				// using the same order as in https://developer.hashicorp.com/vault/api-docs/secret/pki#generate-intermediate-csr
				"common_name": "the common name",
				"alt_names":   "user@example.com,admin@example.com,example.com,www.example.com",
				"ip_sans":     "1.2.3.4,1.2.3.5",
				"uri_sans":    "https://example.com,https://www.example.com",
				"other_sans":  "1.3.6.1.4.1.311.20.2.3;utf8:caadmin@example.com",
				// format
				// private_key_format
				"key_type": "rsa",
				"key_bits": 2048,
				"key_name": "the-key-name",
				// key_ref
				"signature_bits": 384,
				// exclude_cn_from_sans
				"ou":                    "unit1, unit2",
				"organization":          "org1, org2",
				"country":               "CA,US",
				"locality":              "locality1, locality2",
				"province":              "province1, province2",
				"street_address":        "street_address1, street_address2",
				"postal_code":           "postal_code1, postal_code2",
				"serial_number":         "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				"add_basic_constraints": false,
			},
			ttl: 2 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName:         "the common name",
					OrganizationalUnit: []string{"unit1", "unit2"},
					Organization:       []string{"org1", "org2"},
					Country:            []string{"CA", "US"},
					Locality:           []string{"locality1", "locality2"},
					Province:           []string{"province1", "province2"},
					StreetAddress:      []string{"street_address1", "street_address2"},
					PostalCode:         []string{"postal_code1", "postal_code2"},
					SerialNumber:       "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				},
				DNSNames:                      []string{"example.com", "www.example.com"},
				EmailAddresses:                []string{"admin@example.com", "user@example.com"},
				IPAddresses:                   []net.IP{[]byte{1, 2, 3, 4}, []byte{1, 2, 3, 5}},
				URIs:                          []*url.URL{parseURL("https://example.com"), parseURL("https://www.example.com")},
				OtherSANs:                     map[string][]string{"1.3.6.1.4.1.311.20.2.3": {"caadmin@example.com"}},
				IsCA:                          false,
				KeyType:                       "rsa",
				KeyBits:                       2048,
				NotAfter:                      time.Time{},
				KeyUsage:                      0,
				ExtKeyUsage:                   0,
				ExtKeyUsageOIDs:               nil,
				PolicyIdentifiers:             nil,
				BasicConstraintsValidForNonCA: false,
				SignatureBits:                 384,
				UsePSS:                        false,
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 0,
				NotBeforeDuration:             0,
				SKID:                          nil,
			},
			wantFields: map[string]interface{}{
				"common_name":           "the common name",
				"ou":                    "unit1,unit2",
				"organization":          "org1,org2",
				"country":               "CA,US",
				"locality":              "locality1,locality2",
				"province":              "province1,province2",
				"street_address":        "street_address1,street_address2",
				"postal_code":           "postal_code1,postal_code2",
				"alt_names":             "example.com,www.example.com,admin@example.com,user@example.com",
				"ip_sans":               "1.2.3.4,1.2.3.5",
				"uri_sans":              "https://example.com,https://www.example.com",
				"other_sans":            "1.3.6.1.4.1.311.20.2.3;UTF-8:caadmin@example.com",
				"exclude_cn_from_sans":  true,
				"key_type":              "rsa",
				"key_bits":              2048,
				"signature_bits":        384,
				"use_pss":               false,
				"serial_number":         "37:60:16:e4:85:d5:96:38:3a:ed:31:06:8d:ed:7a:46:d4:22:63:d8",
				"add_basic_constraints": false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {

		b, s := CreateBackendWithStorage(t)

		issueTime := time.Now()
		resp, err := CBWrite(b, s, "intermediate/generate/internal", tt.data)
		require.NoError(t, err)
		require.NotNil(t, resp)

		csrData := resp.Data["csr"].(string)
		csr, err := parsing.ParseCertificateRequestFromString(csrData)
		require.NoError(t, err)
		require.NotNil(t, csr)

		t.Run(tt.name+" parameters", func(t *testing.T) {
			testParseCsrToCreationParameters(t, issueTime, tt, csr)
		})
		t.Run(tt.name+" fields", func(t *testing.T) {
			testParseCsrToFields(t, issueTime, tt, csr)
		})
	}
}

func testParseCsrToCreationParameters(t *testing.T, issueTime time.Time, tt *parseCertificateTestCase, csr *x509.CertificateRequest) {
	params, err := certutil.ParseCsrToCreationParameters(*csr)

	if tt.wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		if diff := deep.Equal(tt.wantParams, params); diff != nil {
			t.Errorf("testParseCertificateToCreationParameters() diff: %s", strings.ReplaceAll(strings.Join(diff, "\n"), "map", "\nmap"))
		}
	}
}

func testParseCsrToFields(t *testing.T, issueTime time.Time, tt *parseCertificateTestCase, csr *x509.CertificateRequest) {
	fields, err := certutil.ParseCsrToFields(*csr)
	if tt.wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		if diff := deep.Equal(tt.wantFields, fields); diff != nil {
			t.Errorf("testParseCertificateToFields() diff: %s", strings.ReplaceAll(strings.Join(diff, "\n"), "map", "\nmap"))
		}
	}
}

// TestVerify_chained_name_constraints verifies that we perform name constraints certificate validation using the
// entire CA chain.
//
// This test constructs a root CA that
// - allows: .example.com and myint.com
// - excludes: bad.example.com
//
// and an intermediate that
// - forbids alsobad.example.com
//
// By importing the intermediate chain in the "wrong" order, it validates parsePEM
// is using the updated verification logic, not the historical strict one which
// prevented multi-chains.
//
// It then checks verification by issuing from the intermediate certificate:
// - can issue certs like good.example.com
// - rejects names like notanexample.com since they are not in the namespace of names permitted by the root CA
// - rejects bad.example.com, since the root CA excludes it
// - rejects alsobad.example.com, since the intermediate CA excludes it.
func TestVerify_chained_name_constraints(t *testing.T) {
	t.Parallel()

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Setup
	// In 1.19 and above, this setup is done using the new "excluded_dns_domains" field, not in this version
	// Instead, the certificates have been pre-generated on 1.19 with a 1,000,000 hour (>100 years) ttl
	// If they ever expire for some reason, they can be regenerated as follows:
	//     vault secrets enable pki
	//     vault secrets tune -default-lease-ttl=1000000h -max-lease-ttl=1000000h pki
	//     vault write pki/root/generate/internal ttl="999999h" common_name="myvault.com" permitted_dns_domains=".example.com,myint.com" excluded_dns_domains="bad.example.com"
	//     vault write pki/intermediate/generate/exported common_name="myint.com"
	//     CSR=<csr>
	//     vault write pki/root/sign-intermediate common_name="myint.com" csr="$CSR" ttl="999998h" excluded_dns_domains="alsobad.example.com"

	var b *backend
	var s logical.Storage
	{
		// Intermediate Private Key; Intermediate Certificate and Root Certificate
		pemBundle := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEowIBAAKCAQEAsHQ0W6x4Dbh7E/G8lAHsnXnjDX0ZHkxLutZp/rQW4XvxqYqJ
gXULzkg7lHWb/YAfvW7MgG2Uhs93tAF21j57Oh1yEA6U/6uVr4rePbnZJfuLYknV
JjkES/+3sbYyNUdvvHb6NrWkpquL94KmUiZIN6yfd3xvnHqckiOarRX3zk9Oi6MY
azbhYDZpyBtPj0FS9UENoaPkGUXz3J648w41MG83BzPsjE/F2Awcd0dRqJOdxdo2
y/or+aXHZOqYumGt1XllCIgBMFEdf6iV1SmREim86BNfqojulacD2QWKoSw0Ju6t
qdcLtpjDV/jQwT9h56Z53TK9HSCPRQH3W07JsQIDAQABAoIBAErlN+gFX3urZwpD
G+DSpET/HEF4bLXwemBWa+0BMYkJm7xzogPB033+KNPNy6Ugwdgc9E46tMU3wD42
NOOg55r2LpjkJ5gRk1pMHiKjBKlbiIDyh7dOkiVT5/tVMyBsLuAWvAOGu9lGjQVE
fiX2O8btzU3hElqfrF9XH2sCHQmC6ON8IT5UdVUMITWd7YhhjNOSjuEX2bBuzdH8
cLF2L2sFnGczEI0ob7mcDDYwpBaJ2JBZFIkud/2VBdjD2Gs8+5SHTN6kW4Wvomih
WEjmdA47CN6H+IVg6R616vUDSjyZLstS7h94edguBUm+GovxQHhkmF8HyF/dGVE6
B02D98UCgYEA55vwZMw5J+WkBtoDTb3TEaRkblVVxvuGSHLFOahR/f9+8sbsLkc1
6Xrg3Z8RFuYTSa6cMIIaddA96QLVQ+6dnN6Fh5UW2me30zaA2/dcfDULqH/AFmlP
XNWnkANI7nI7S67rGpLS/x3u4l3u2rGVWGyBH733PIF9N7oT2OXWeBcCgYEAwwlP
UiXWFx1eZXQ57+yyzBY8MPnED7uAiN4I/zSVrwfiOApbBEG5P6xTeTU1pa27D6Vw
8iQKT4nfG6r7sKHJizG7uRl/ukCMoD+kO/nrL1f8aglJkTpwM/c00uJaD+AYVIRB
qcMuHP10xe5nF8ou8WsfdIBSndlix/3Au+0uIXcCgYBGu8IkFpbSD3kTuptNr855
Udp2M8uZlJGDKMIBBN3dk50CqivRSmA1qRhptr2yX5u+YfxDeEh6eJmfwwS81W//
S5o+ORMcqpZig//Arf/p0w13bWSNs6XQNVQSUF6CrBJd39Lfb73OuiaNraBX2+o/
cgvRVOPf76/9R+/tyuvACwKBgQCU81A72Z+WiZxk+Q2iWW/3g2e2s5lztmwkmwXd
+cfzI58101rNdnBG25HoeKWfAX5iJaSLmNCzDAVzOPL9Y6HrhzDp3Tp6AQAlBlWk
ZGqj/ptMVAl2O14jocEa6TmE9E1Ahr0rLF/UThFXdIoRVhzhzUD5lWMDIFlbiUUf
gLUBYwKBgGnl1KVj9IAQHj4o52Nwwc9q3R70FCt81s3w6Y0Ixm4atc1LwV+U6QM4
k+UxzUBQedLIjMlHmbipx6lpCpGQBZxAO9daZ2Qdks24wi2KNGEWFQyMeYv1r4XO
mXDJbHDJYOYZ29g3GHDds3SkiEe2iLD20hnkgszdsm2FqbujekHn
-----END RSA PRIVATE KEY-----
-----BEGIN CERTIFICATE-----
MIIDXDCCAkSgAwIBAgIUGwwoD7Thgv/vz2TQYuMtpLt9pWQwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wIBcNMjUwMjAzMTQ1ODMwWhgPMjEz
OTAzMDUwNDU5MDBaMBQxEjAQBgNVBAMTCW15aW50LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBALB0NFuseA24exPxvJQB7J154w19GR5MS7rWaf60
FuF78amKiYF1C85IO5R1m/2AH71uzIBtlIbPd7QBdtY+ezodchAOlP+rla+K3j25
2SX7i2JJ1SY5BEv/t7G2MjVHb7x2+ja1pKari/eCplImSDesn3d8b5x6nJIjmq0V
985PToujGGs24WA2acgbT49BUvVBDaGj5BlF89yeuPMONTBvNwcz7IxPxdgMHHdH
UaiTncXaNsv6K/mlx2TqmLphrdV5ZQiIATBRHX+oldUpkRIpvOgTX6qI7pWnA9kF
iqEsNCburanXC7aYw1f40ME/Yeemed0yvR0gj0UB91tOybECAwEAAaOBoTCBnjAO
BgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQU9Sn3+AJg
lLHCRMJZWuXwK2Quk38wHwYDVR0jBBgwFoAU2Ub/bdY7gdjZm9yT1/CxPh8N0VUw
FAYDVR0RBA0wC4IJbXlpbnQuY29tMCUGA1UdHgEB/wQbMBmhFzAVghNhbHNvYmFk
LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQB/ucgSAJaHl+rKt+Uzh0Np
WfQSgimFHE+5IVwe9GjUC/8xMHCphU28obbMeORh8c9RpYflvHcf0AZCSljwzcf+
U7KzDFRahHHUuuOCub+OW0P1O7iVKrvw24u/eQ2U6umJHnr6iED15tc2LTIynWQ2
0t6B8w+AKOU7dB6h0Ar6Vr1DD6Q/HpugjZpEKrKEqDtB4V6pUTGj7quHeDNeX4hK
foRgA1xHY0zgrqzPDXXs7NbGlF+2SDImGO7XwCieJgnOWpQDcFCws65bn3xQmo2m
rElegmDHtGE+J0raGswYBBUTedQ6fGliQVadev+0xR7YojOrS15LRVmxy7UQ/Vzd
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDezCCAmOgAwIBAgIUFufV7PUhESnMe3jsbaFOaJdGmBcwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wIBcNMjUwMjAzMTQ1NzE4WhgPMjEz
OTAzMDUwNTU3NDdaMBYxFDASBgNVBAMTC215dmF1bHQuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAz7wyAmSIrnb9MPWXmHW1cT03hoDkSWZK76Oi
MoRVnCFTsY2eKrofhUpO6Cxkl4PhUCST6oFdnkTbm5Oxg5q7i1iRLi2s1l/XYOpk
GQl+WeG59mxeQ6j98K9EAHmN6zm5K4xjej9uHkZK6Nlm2qzyjEaq94pQl56jMH97
mXsifvqkZUJV0EvW1BYNl+VqMUfMYkuRsS+OOGFY6qinlEL/K1OmDNP8+2RsCo5h
HbTAf90y1skAyOA99vYkySNubIxSYbl84CrsmhuWY9+oBRBRxLEeTfSbTP4nb2Gh
VbyEiQhcp94mg9rR9Jd5vN3QjsbNBCP7mdFK734acpfu4YYIowIDAQABo4G+MIG7
MA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBTZRv9t
1juB2Nmb3JPX8LE+Hw3RVTAfBgNVHSMEGDAWgBTZRv9t1juB2Nmb3JPX8LE+Hw3R
VTAWBgNVHREEDzANggtteXZhdWx0LmNvbTBABgNVHR4BAf8ENjA0oB0wDoIMLmV4
YW1wbGUuY29tMAuCCW15aW50LmNvbaETMBGCD2JhZC5leGFtcGxlLmNvbTANBgkq
hkiG9w0BAQsFAAOCAQEAPNVRHHCkBfq6qGOdevMIjCQEsRoIhwAx398+oMUrgjKd
vslOP6dqXoWmKJnFtPKFWpxIkUm4v83Kx7TNNlMwAzf1Hyhjq17OaH0jDbe8LcmY
CCRrK6yxeN9iJZInJKwmWHtVVAvdbxY2iwLwtqY6TlajP76aO8Zi5rasvEBurhmX
97acn2LgbmxrDnNEaU2gT60ft0noaGXUrK9WX1LoT0+0dtFSqr5tR/S0sXVioVnG
ypyP4Fe/DXvXHrLw/v+Kldu5ibKRgx4iCT1vC6CX93/m/AQ7tWawLY+BymFLbrNh
zOYF7zwf65ExfN+o1q38bFDTOi4idCcL07e7uk2uDQ==
-----END CERTIFICATE-----
`)

		// Create the Certificates by Importing It into the Mount
		b, s = CreateBackendWithStorage(t)

		// Note that we append the root CA certificate to the signed intermediate, so that
		// the entire chain is stored by set-signed.
		resp, err := CBWrite(b, s, "issuers/import/bundle", map[string]interface{}{
			"pem_bundle": pemBundle,
		})
		require.NoError(t, err)

		// Create a Role in the Mount
		resp, err = CBWrite(b, s, "roles/test", map[string]interface{}{
			"allow_bare_domains": true,
			"allow_subdomains":   true,
			"allow_any_name":     true,
		})
		require.NoError(t, err)
		require.NotNil(t, resp)
	}

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Tests

	testCases := []struct {
		commonName string
		wantError  string
	}{
		{
			commonName: "good.example.com",
		},
		{
			commonName: "notanexample.com",
			wantError:  "should not be permitted by root CA",
		},
		{
			commonName: "bad.example.com",
			wantError:  "should be rejected by the root CA",
		},
		{
			commonName: "alsobad.example.com",
			wantError:  "should be rejected by the intermediate CA",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.commonName, func(t *testing.T) {
			resp, err := CBWrite(b, s, "issue/test", map[string]any{
				"common_name": tc.commonName,
			})
			if tc.wantError != "" {
				require.Error(t, err, tc.wantError)
				require.ErrorContains(t, err, "certificate is not authorized to sign for this name")
				require.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NoError(t, resp.Error())
			}
		})
	}
}
