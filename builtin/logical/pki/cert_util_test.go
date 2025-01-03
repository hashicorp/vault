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
// - allows: .example.com
// - excludes: bad.example.com
//
// and an intermediate that
// - forbids alsobad.example.com
//
// It verifies that the intermediate
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
	//     vault write pki/root/generate/internal ttl="999999h" common_name="myvault.com" permitted_dns_domains=".example.com" excluded_dns_domains="bad.example.com"
	//     vault write pki/intermediate/generate/exported common_name="myint.com"
	//     CSR=<csr>
	//     vault write pki/root/sign-intermediate common_name="myint.com" csr="$CSR" ttl="999998h" excluded_dns_domains="alsobad.example.com"

	var b *backend
	var s logical.Storage
	{
		// Intermediate Private Key; Intermediate Certificate and Root Certificate
		pemBundle := []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEA4xVskyF6ykn2yYMsK5VRemMpm2nlzsT3stiSyER7rzLzId5Q
9OjH3QOXR94hLTab/tJteukRtzZT8Rg5YmLIvSvNGyzjGUV4eErj1cRKUsOW1JtN
UXrKLIcjpGpoB8vDxxMHZhJL8JGOL+fx6qLuboWFuKmPJ1nE4tgS6yHR7x4oBp1V
oU0bUEPXTo2CYVIFuvIjK2mDpuUCprJtg0si9WbZi/59SCHJzAIti6xjJnGV8QnQ
/NHjiB2krL/0oeDkLRwS72wvmbYQ3PU7i23zVcEpcoOWUXXIeU645bupJY4dLC5Y
Viuj4UDuCiVrTb+Ea+pAJlrd5UIselSN+7KktQIDAQABAoIBAQCAQ0SdYiayBc8A
CTg0sdVgtIv2vXzRKo3iFdPqjEv0LGoJ8kF149mn63RSYpQIrrSz3PV7nBOmkWge
YJlhCfzqZMgoFlV7m7Ks91fzETkNwG38TnAAmsOBHR+zqWpzJNPDKOtf6uu7yOsw
Aemxpy/Xe1GJeTRjfJ/ppTQiXWrvjN+WPGKyOwwkyzZE08Xv5dnqkyCOX0TCGQfO
4CE7FoeMOm+A1Qxnl5OPQxjnS6IXOMWZNTuKWVs0s806ojG9WheA00WEbyIDIYSp
4XampDFjt3wQYwrv9hHV82f0Otqlys/4skq7Wp5Bl+/7QOOTlqzkRqk2CnQJXOA4
1e9mYfFtAoGBAPSNrTTCcFvVh3FBqt8kQjqyYDC86JyYmb+ZF3mpu9eDxrASYAT6
QoMAUUbvFSVdzcHufN+bymLDMPGPAkexFdlY+BxYF7d6CGXivgZBl+yQtZ6e7WkS
T+WB5wHySp5igIbV5SuVY+22pGGBVL0kqbPKFG/+MOlRqMBX15pPrSTjAoGBAO22
a0xYdqtF0jXgAHlKN4Z7UNdfsrkAWnMVaBt2Mi8FlomDGoXax9E6q5hmGthhFafV
FKqCZ8e148Wc3ry0Z/KYs4SRASUx/NOuQlpHre3DVLLSNawYXU3HqD06wEiQPWzR
fq+UoKsCiJC+qz8f/UD8UuiZZAa9EthOSpNxstuHAoGBAIIFgX13k31//c79dvfE
s2G5zOKczZ/UkooHvy90Sua+rTiXzG1ZEVvNI2lvW/LN+MOPJN1OW0A/PxpvSmsL
f+5bGy8WtyVZwHVLJHT3Eus31RhMrzUaA1imxEeIppunC2ak+n89oi+U17jvpjoZ
8BAi9NLGdwLV476/9WWZzxi5AoGATq0wyD0DUd6zG4e/QGWzCPypng8bfSXDyhFM
usIdC/kigPL2hVULC5IKl088FV/Upg7dXy34IV5vO8mW4wgm22F1ESxZH7Fyx7EG
XxEYXPhogSMYBpSt1P9/DHz0hU/QNMMF1iEwKEmXX6jrzuHMlYSuADQ8qgpMQXFw
N2rLUuMCgYABtwULLZJz8eMKRRFp6fJNcqqXFmCCaB4M3tGPGNnSpnTaqmBggYg6
NRAVwplparFZ0vg5ZSegCZzAAwCNLDHOBtlsakStDaXB2EdAhS5vKERh5HpWUeeR
iKyfuPEjfsS51MzFKunVWLotGBy33SWWhSevI03d/ECnC7E319+BcQ==
-----END RSA PRIVATE KEY-----
-----BEGIN CERTIFICATE-----
MIIDXDCCAkSgAwIBAgIUTOa1pFOYnf22Cu1g+SLtEdeXif4wDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wIBcNMjUwMTAzMjIxNzU0WhgPMjEz
OTAyMDIxMjE4MjRaMBQxEjAQBgNVBAMTCW15aW50LmNvbTCCASIwDQYJKoZIhvcN
AQEBBQADggEPADCCAQoCggEBAOMVbJMhespJ9smDLCuVUXpjKZtp5c7E97LYkshE
e68y8yHeUPTox90Dl0feIS02m/7SbXrpEbc2U/EYOWJiyL0rzRss4xlFeHhK49XE
SlLDltSbTVF6yiyHI6RqaAfLw8cTB2YSS/CRji/n8eqi7m6FhbipjydZxOLYEush
0e8eKAadVaFNG1BD106NgmFSBbryIytpg6blAqaybYNLIvVm2Yv+fUghycwCLYus
YyZxlfEJ0PzR44gdpKy/9KHg5C0cEu9sL5m2ENz1O4tt81XBKXKDllF1yHlOuOW7
qSWOHSwuWFYro+FA7gola02/hGvqQCZa3eVCLHpUjfuypLUCAwEAAaOBoTCBnjAO
BgNVHQ8BAf8EBAMCAQYwDwYDVR0TAQH/BAUwAwEB/zAdBgNVHQ4EFgQUIvhScHgB
Q3s9uuhM4wCr9n0JOrkwHwYDVR0jBBgwFoAUaRs3bAPVLwyZKS0PIKmoER4vreIw
FAYDVR0RBA0wC4IJbXlpbnQuY29tMCUGA1UdHgEB/wQbMBmhFzAVghNhbHNvYmFk
LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IBAQC/Yw+gctU5vEr5oS1gHCmm
9z6aB8fr8UTFAjHfVzggfXf3DKU6NJ9rqKQ5PNvHpVYNacgO3bCyDa0BFPQ9AsE/
vBDKspfF0ARgK7LFhuQi9hoUM6BfEt5CHOUVMYvx7MnYxLsdiNOCjbvnFjSz70Nr
+q63d0dISHAEjRW1+DbCMUihi/npOH0xKezDu23010Rfsos08ZVBJw9aqkAAe7zT
J7/RZ6ITnqQwJc2iMshIXbP8HrrLcEJ+9lyQK0h7gZ3kgcrObtg6Xr80HSDRus0U
z9J2vNwKYlEOXHVPU9X/zdNJhDuv+/OB469ehVcetyd4+q9dF5LNBwi5ZG6CLJXz
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDbjCCAlagAwIBAgIUH5XPxepR27cGAFwQsmNxQOEQhdEwDQYJKoZIhvcNAQEL
BQAwFjEUMBIGA1UEAxMLbXl2YXVsdC5jb20wIBcNMjUwMTAzMjIxNjM1WhgPMjEz
OTAyMDIxMzE3MDVaMBYxFDASBgNVBAMTC215dmF1bHQuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEAwway6sN9F5iQ7JBg5+5HOpfhHpUhfZyjwcIP
OWfRVU/+tD9WGdI42apMqMZu4lXhJGMnNKAZbTszyCcI4/a8INy+cTu/OJLT++r/
JCfLYcSnBBnCw8t/T0MvE3SWH/H0ToMIwSZPuMhBjVA6XlCopVzhZy8x+nM7oqNC
gqqW1e/oK2V6zZDlxD+hTDDdQnEhnFu9qyr0wOacAfqrtTwnR9G770CECF95mRzJ
ULmS5SJHkSxZzetRc+cTao4bar9JU2WUhih7qGs5OjHw6YmVosNWLAV38vX9DPy9
J2wCh7r4ttTir0GQdl1YaVDBtG5EQcuCEFpq79i+AgOK/LuMTQIDAQABo4GxMIGu
MA4GA1UdDwEB/wQEAwIBBjAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBRpGzds
A9UvDJkpLQ8gqagRHi+t4jAfBgNVHSMEGDAWgBRpGzdsA9UvDJkpLQ8gqagRHi+t
4jAWBgNVHREEDzANggtteXZhdWx0LmNvbTAzBgNVHR4BAf8EKTAnoBAwDoIMLmV4
YW1wbGUuY29toRMwEYIPYmFkLmV4YW1wbGUuY29tMA0GCSqGSIb3DQEBCwUAA4IB
AQBv65+kEvaAPYgLnfZjaXkFSysuVR5TI2sw9qMj23mCOfIm7mBzPkGKnB8g4XSg
tsOf+6UKcgQVkTvGSKqjavrHrtRtojgKvnXfra3yPY41//G2IWZ6juiQzqqj55Fm
rM952Wrr7iIjR5pLP6sZ74OCjxTr35Ky3Y2+k4blHlGA/JBQeGNhiyO6mYSqzY1i
tRI8ij2/5wSXYaPDl3DBSJSEDqYpSUINSUK2W1R5Rq10yD5HDREh96oeOYLKC3a0
C73G+J0nGaYffdR4MWeTEq6gfENZXrmlA5Fe4fyBgxlHfzxx5c0KWpZVJttoVGir
j8bwlXImoBeiKZL7lEGT2vDk
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
