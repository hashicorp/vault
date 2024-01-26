// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"context"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"github.com/go-test/deep"
	"github.com/hashicorp/vault/builtin/logical/pki/parsing"
	"github.com/hashicorp/vault/sdk/helper/certutil"
	"github.com/stretchr/testify/require"
	"net"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
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
			"certs/",
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
	cb, _, err := generateCreationBundle(b, input, nil, nil)
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
			cb, _, err := generateCreationBundle(b, testCase.input, nil, nil)
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
	ttl        time.Duration
	wantParams certutil.CreationParameters
	wantFields map[string]interface{}
	wantErr    bool
}

func TestParseCertificateToCreationParameters(t *testing.T) {
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
				"common_name":    "the common name",
				"key_type":       "ec",
				"key_bits":       384,
				"ttl":            "1h",
				"street_address": "",
			},
			ttl: 1 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName: "the common name", // add logic to make this work
				},
				DNSNames:                      nil,
				EmailAddresses:                nil,
				IPAddresses:                   nil,
				URIs:                          nil,
				OtherSANs:                     make(map[string][]string),
				IsCA:                          true,
				KeyType:                       "ec", // cover all the getKeytype() types
				KeyBits:                       384,
				NotAfter:                      time.Time{},
				KeyUsage:                      x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
				ExtKeyUsage:                   0,
				ExtKeyUsageOIDs:               nil,
				PolicyIdentifiers:             nil,
				BasicConstraintsValidForNonCA: false, // don't assert for CAs
				SignatureBits:                 384,   // look at findSignatureBits
				UsePSS:                        false, // look at isPSS
				ForceAppendCaChain:            false,
				UseCSRValues:                  false,
				PermittedDNSDomains:           nil,
				URLs:                          nil,
				MaxPathLength:                 -1,
				NotBeforeDuration:             0,   // assert that it is greater than 30 s (the default)
				SKID:                          nil, // assert that it is not nil
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
				"skid":                  "",
				"key_type":              "ec",
				"key_bits":              384,
			},
			wantErr: false,
		},
		{
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
				"permitted_dns_domains": ".example.com,.www.example.com",
				"ou":                    "unit1, unit2",
				"organization":          "org1, org2",
				"country":               "US, CA",
				"locality":              "locality1, locality2",
				"province":              "province1, province2",
				"street_address":        "street_address1, street_address2",
				"postal_code":           "postal_code1, postal_code2",
				"not_before_duration":   "30s",
				"key_type":              "rsa",
				"use_pss":               true,
				"key_bits":              2048,
				"signature_bits":        384,
				// TODO(kitography): Specify key usage
			},
			ttl: 2 * time.Hour,
			wantParams: certutil.CreationParameters{
				Subject: pkix.Name{
					CommonName:         "the common name", // add logic to make this work
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
				OtherSANs:                     map[string][]string{"1.3.6.1.4.1.311.20.2.3": []string{"caadmin@example.com"}},
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
				PermittedDNSDomains:           []string{".example.com", ".www.example.com"},
				URLs:                          nil,
				MaxPathLength:                 2,
				NotBeforeDuration:             0,
				SKID:                          nil,
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
				"ttl":                   "1h0m30s",
				"max_path_length":       2,
				"permitted_dns_domains": ".example.com,.www.example.com",
				"use_pss":               true,
				"skid":                  "",
				"key_type":              "rsa",
				"key_bits":              2048,
			},
			wantErr: false,
		},
		// need a test for non CA
		// need a test with a different ttl
	}
	for _, tt := range tests {

		b, s := CreateBackendWithStorage(t)

		issueTime := time.Now()
		resp, err := CBWrite(b, s, "root/generate/internal", tt.data)
		require.NoError(t, err)
		require.NotNil(t, resp)

		certData := resp.Data["certificate"].(string)
		cert, err := parsing.ParseCertificateFromString(certData)
		require.NoError(t, err)
		require.NotNil(t, cert)

		t.Run(tt.name+" parameters", func(t *testing.T) {
			testParseCertificateToCreationParameters(t, issueTime, tt, cert)
		})
		t.Run(tt.name+" fields", func(t *testing.T) {
			testParseCertificateToFields(t, issueTime, tt, cert)
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
		require.GreaterOrEqual(t, params.NotBeforeDuration, 30*time.Second)

		// with ttl=1h
		require.GreaterOrEqual(t, params.NotAfter, issueTime.Add(59*time.Minute))
		require.LessOrEqual(t, params.NotAfter, issueTime.Add(61*time.Minute))
	}
}

func testParseCertificateToFields(t *testing.T, issueTime time.Time, tt *parseCertificateTestCase, cert *x509.Certificate) {
	fields, err := certutil.ParseCertificateToFields(*cert)
	t.Log(fields)
	if tt.wantErr {
		require.Error(t, err)
	} else {
		require.NoError(t, err)

		require.NotNil(t, fields["skid"])
		fieldsToDelete := []string{"skid"}
		for _, f := range fieldsToDelete {
			delete(fields, f)
			delete(tt.wantFields, f)
		}

		if diff := deep.Equal(tt.wantFields, fields); diff != nil {
			t.Errorf("testParseCertificateToFields() diff: %s", strings.ReplaceAll(strings.Join(diff, "\n"), "map", "\nmap"))
		}

		//require.NotNil(t, params.SKID)
		//require.GreaterOrEqual(t, params.NotBeforeDuration, 30*time.Second)
		//
		//// with ttl=1h
		//require.GreaterOrEqual(t, params.NotAfter, issueTime.Add(59*time.Minute))
		//require.LessOrEqual(t, params.NotAfter, issueTime.Add(61*time.Minute))
	}
}
