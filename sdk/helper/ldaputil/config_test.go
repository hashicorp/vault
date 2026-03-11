// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/hashicorp/vault/sdk/framework"
)

func TestCertificateValidation(t *testing.T) {
	// certificate should default to "" without error if it doesn't exist
	config := testConfig(t)
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	if config.Certificate != "" {
		t.Fatalf("expected no certificate but received %s", config.Certificate)
	}

	// certificate should cause an error if a bad one is provided
	config.Certificate = "cats"
	if err := config.Validate(); err == nil {
		t.Fatal("should err due to bad cert")
	}

	// valid certificates should pass inspection
	config.Certificate = validCertificate
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestNewConfigEntry(t *testing.T) {
	s := &framework.FieldData{Schema: ConfigFields()}
	config, err := NewConfigEntry(nil, s)
	if err != nil {
		t.Fatal("error getting default config")
	}
	configFromJSON := testJSONConfig(t, jsonConfigDefault)

	t.Run("equality_check", func(t *testing.T) {
		if diff := deep.Equal(config, configFromJSON); len(diff) > 0 {
			t.Fatalf("bad, diff: %#v", diff)
		}
	})
}

func TestConfig(t *testing.T) {
	config := testConfig(t)
	configFromJSON := testJSONConfig(t, jsonConfig)

	t.Run("equality_check", func(t *testing.T) {
		if diff := deep.Equal(config, configFromJSON); len(diff) > 0 {
			t.Fatalf("bad, diff: %#v", diff)
		}
	})

	t.Run("default_use_token_groups", func(t *testing.T) {
		if config.UseTokenGroups {
			t.Errorf("expected false UseTokenGroups but got %t", config.UseTokenGroups)
		}

		if configFromJSON.UseTokenGroups {
			t.Errorf("expected false UseTokenGroups from JSON but got %t", configFromJSON.UseTokenGroups)
		}
	})

	t.Run("default_schema_type", func(t *testing.T) {
		if config.Schema != SchemaOpenLDAP {
			t.Errorf("expected default Schema %s but got %s", SchemaOpenLDAP, config.Schema)
		}
	})
}

func TestNormalizedSchema(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty_string_defaults_to_openldap",
			input:    "",
			expected: SchemaOpenLDAP,
		},
		{
			name:     "lowercase_ad",
			input:    "ad",
			expected: SchemaAD,
		},
		{
			name:     "uppercase_AD",
			input:    "AD",
			expected: SchemaAD,
		},
		{
			name:     "mixed_case_Ad",
			input:    "Ad",
			expected: SchemaAD,
		},
		{
			name:     "lowercase_openldap",
			input:    "openldap",
			expected: SchemaOpenLDAP,
		},
		{
			name:     "uppercase_OPENLDAP",
			input:    "OPENLDAP",
			expected: SchemaOpenLDAP,
		},
		{
			name:     "mixed_case_OpenLDAP",
			input:    "OpenLDAP",
			expected: SchemaOpenLDAP,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizedSchema(tt.input)
			if result != tt.expected {
				t.Errorf("NormalizedSchema(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSchemaInConfigEntry(t *testing.T) {
	t.Run("new_config_defaults_to_openldap", func(t *testing.T) {
		s := &framework.FieldData{Schema: ConfigFields()}
		config, err := NewConfigEntry(nil, s)
		if err != nil {
			t.Fatal("error getting default config")
		}

		if config.Schema != SchemaOpenLDAP {
			t.Errorf("expected default Schema %s but got %s", SchemaOpenLDAP, config.Schema)
		}
	})

	t.Run("schema_normalized_on_create", func(t *testing.T) {
		schema := ConfigFields()
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				"schema": "AD",
			},
			Schema: schema,
		}

		config, err := NewConfigEntry(nil, data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if config.Schema != SchemaAD {
			t.Errorf("expected normalized Schema 'ad' but got %s", config.Schema)
		}
	})

	t.Run("schema_in_passwordless_map", func(t *testing.T) {
		config := &ConfigEntry{
			Schema: SchemaAD,
		}

		m := config.PasswordlessMap()
		schema, ok := m["schema"]
		if !ok {
			t.Error("schema not found in PasswordlessMap")
		}

		if schema != SchemaAD {
			t.Errorf("expected schema %s in map but got %v", SchemaAD, schema)
		}
	})

	t.Run("schema_update_existing_config", func(t *testing.T) {
		existingConfig := &ConfigEntry{
			Url:    "ldap://127.0.0.1",
			Schema: SchemaOpenLDAP,
		}

		schema := ConfigFields()
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				"schema": "AD",
			},
			Schema: schema,
		}

		updatedConfig, err := NewConfigEntry(existingConfig, data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if updatedConfig.Schema != SchemaAD {
			t.Errorf("expected updated Schema 'ad' but got %s", updatedConfig.Schema)
		}
	})
}

func TestSupportedSchemas(t *testing.T) {
	schemas := SupportedSchemas()

	expectedSchemas := []string{SchemaOpenLDAP, SchemaAD}
	if len(schemas) != len(expectedSchemas) {
		t.Errorf("expected %d schemas but got %d", len(expectedSchemas), len(schemas))
	}

	for _, expected := range expectedSchemas {
		found := false
		for _, schema := range schemas {
			if schema == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected schema %q not found in SupportedSchemas()", expected)
		}
	}
}

func TestSchemaValidation(t *testing.T) {
	t.Run("valid_openldap_schema", func(t *testing.T) {
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        SchemaOpenLDAP,
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		if err := config.Validate(); err != nil {
			t.Errorf("expected no error for valid openldap schema, got: %v", err)
		}
	})

	t.Run("valid_ad_schema", func(t *testing.T) {
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        SchemaAD,
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		if err := config.Validate(); err != nil {
			t.Errorf("expected no error for valid ad schema, got: %v", err)
		}
	})

	t.Run("empty_schema_passes_validation", func(t *testing.T) {
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        "",
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		if err := config.Validate(); err != nil {
			t.Errorf("expected no error for empty schema (defaults to openldap), got: %v", err)
		}
	})

	t.Run("generic_unsupported_schema_fails_validation", func(t *testing.T) {
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        "unsupported_schema",
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		err := config.Validate()
		if err == nil {
			t.Error("expected error for unsupported schema, got nil")
		}

		if err != nil && !strings.Contains(err.Error(), "unsupported schema") {
			t.Errorf("expected error message to contain 'unsupported schema', got: %v", err)
		}
	})

	t.Run("freeipa_schema_fails_validation", func(t *testing.T) {
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        "freeipa",
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		err := config.Validate()
		if err == nil {
			t.Error("expected error for unsupported schema 'freeipa', got nil")
		}

		if err != nil && !strings.Contains(err.Error(), "freeipa") {
			t.Errorf("expected error message to include 'freeipa', got: %v", err)
		}
	})

	t.Run("typo_in_schema_fails_validation", func(t *testing.T) {
		// Test common typos like "adc" instead of "ad"
		config := &ConfigEntry{
			Url:           "ldap://127.0.0.1",
			UserDN:        "ou=users,dc=example,dc=com",
			Schema:        "adc",
			TLSMinVersion: "tls12",
			TLSMaxVersion: "tls12",
		}

		err := config.Validate()
		if err == nil {
			t.Error("expected error for typo schema 'adc', got nil")
		}

		// Verify error message is helpful
		if err != nil && !strings.Contains(err.Error(), "unsupported schema") {
			t.Errorf("expected error message to contain 'unsupported schema', got: %v", err)
		}
		if err != nil && !strings.Contains(err.Error(), "adc") {
			t.Errorf("expected error message to include the unsupported value 'adc', got: %v", err)
		}
	})

	t.Run("normalized_unsupported_schema_fails_validation", func(t *testing.T) {
		// Test that even after normalization (lowercase), unsupported schemas fail
		schema := ConfigFields()
		data := &framework.FieldData{
			Raw: map[string]interface{}{
				"url":    "ldap://127.0.0.1",
				"userdn": "ou=users,dc=example,dc=com",
				"schema": "UNSUPPORTED_SCHEMAS",
			},
			Schema: schema,
		}

		// NewConfigEntry should fail because schema validation happens during creation
		_, err := NewConfigEntry(nil, data)
		if err == nil {
			t.Error("expected error from NewConfigEntry for unsupported schema, got nil")
		}

		// Verify the error message is helpful
		if err != nil {
			if !strings.Contains(err.Error(), "unsupported schema") {
				t.Errorf("expected error message to contain %q, got: %v", "unsupported schema", err)
			}
		}
	})
}

func testConfig(t *testing.T) *ConfigEntry {
	t.Helper()

	return &ConfigEntry{
		Url:               "ldap://138.91.247.105",
		UserDN:            "example,com",
		BindDN:            "kitty",
		BindPassword:      "cats",
		TLSMaxVersion:     "tls12",
		TLSMinVersion:     "tls12",
		RequestTimeout:    30,
		ConnectionTimeout: 15,
		ClientTLSCert:     "",
		ClientTLSKey:      "",
		Schema:            SchemaOpenLDAP,
	}
}

func testJSONConfig(t *testing.T, rawJson []byte) *ConfigEntry {
	t.Helper()

	config := new(ConfigEntry)
	if err := json.Unmarshal(rawJson, config); err != nil {
		t.Fatal(err)
	}
	return config
}

const validCertificate = `
-----BEGIN CERTIFICATE-----
MIIF7zCCA9egAwIBAgIJAOY2qjn64Qq5MA0GCSqGSIb3DQEBCwUAMIGNMQswCQYD
VQQGEwJVUzEQMA4GA1UECAwHTm93aGVyZTERMA8GA1UEBwwIVGltYnVrdHUxEjAQ
BgNVBAoMCVRlc3QgRmFrZTENMAsGA1UECwwETm9uZTEPMA0GA1UEAwwGTm9ib2R5
MSUwIwYJKoZIhvcNAQkBFhZkb25vdHRydXN0QG5vd2hlcmUuY29tMB4XDTE4MDQw
MzIwNDQwOFoXDTE5MDQwMzIwNDQwOFowgY0xCzAJBgNVBAYTAlVTMRAwDgYDVQQI
DAdOb3doZXJlMREwDwYDVQQHDAhUaW1idWt0dTESMBAGA1UECgwJVGVzdCBGYWtl
MQ0wCwYDVQQLDAROb25lMQ8wDQYDVQQDDAZOb2JvZHkxJTAjBgkqhkiG9w0BCQEW
FmRvbm90dHJ1c3RAbm93aGVyZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAw
ggIKAoICAQDzQPGErqjaoFcuUV6QFpSMU6w8wO8F0othik+rrlKERmrGonUGsoum
WqRe6L4ZnxBvCKB6EWjvf894TXOF2cpUnjDAyBePISyPkRBEJS6VS2SEC4AJzmVu
a+P+fZr4Hf7/bEcUr7Ax37yGVZ5i5ByNHgZkBlPxKiGWSmAqIDRZLp9gbu2EkG9q
NOjNLPU+QI2ov6U/laGS1vbE2LahTYeT5yscu9LpllxzFv4lM1f4wYEaM3HuOxzT
l86cGmEr9Q2N4PZ2T0O/s6D4but7c6Bz2XPXy9nWb5bqu0n5bJEpbRFrkryW1ozh
L9uVVz4dyW10pFBJtE42bqA4PRCDQsUof7UfsQF11D1ThrDfKsQa8PxrYdGUHUG9
GFF1MdTTwaoT90RI582p+6XYV+LNlXcdfyNZO9bMThu9fnCvT7Ey0TKU4MfPrlfT
aIhZmyaHt6mL5p881UPDIvy7paTLgL+C1orLjZAiT//c4Zn+0qG0//Cirxr020UF
3YiEFk2H0bBVwOHoOGw4w5HrvLdyy0ZLDSPQbzkSZ0RusHb5TjiyhtTk/h9vvJv7
u1fKJub4MzgrBRi16ejFdiWoVuMXRC6fu/ERy3+9DH6LURerbPrdroYypUmTe9N6
XPeaF1Tc+WO7O/yW96mV7X/D211qjkOtwboZC5kjogVbaZgGzjHCVwIDAQABo1Aw
TjAdBgNVHQ4EFgQU2zWT3HeiMBzusz7AggVqVEL5g0UwHwYDVR0jBBgwFoAU2zWT
3HeiMBzusz7AggVqVEL5g0UwDAYDVR0TBAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AgEAwTGcppY86mNRE43uOimeApTfqHJv+lGDTjEoJCZZmzmtxFe6O9+Vk4bH/8/i
gVQvqzBpaWXRt9OhqlFMK7OkX4ZvqXmnShmxib1dz1XxGhbwSec9ca8bill59Jqa
bIOq2SXVMcFD0GwFxfJRBVzHHuB6AwV9B2QN61zeB1oxNGJrUOo80jVkB7+MWMyD
bQqiFCHWGMa6BG4N91KGOTveZCGdBvvVw5j6lt731KjbvL2hB1UHioucOweKLfa4
QWDImTEjgV68699wKERNL0DCpeD7PcP/L3SY2RJzdyC1CSR7O8yU4lQK7uZGusgB
Mgup+yUaSjxasIqYMebNDDocr5kdwG0+2r2gQdRwc5zLX6YDBn6NLSWjRnY04ZuK
P1cF68rWteWpzJu8bmkJ5r2cqskqrnVK+zz8xMQyEaj548Bnt51ARLHOftR9jkSU
NJWh7zOLZ1r2UUKdDlrMoh3GQO3rvnCJJ16NBM1dB7TUyhMhtF6UOE62BSKdHtQn
d6TqelcRw9WnDsb9IPxRwaXhvGljnYVAgXXlJEI/6nxj2T4wdmL1LWAr6C7DuWGz
8qIvxc4oAau4DsZs2+BwolCFtYc98OjWGcBStBfZz/YYXM+2hKjbONKFxWdEPxGR
Beq3QOqp2+dga36IzQybzPQ8QtotrpSJ3q82zztEvyWiJ7E=
-----END CERTIFICATE-----
`

var jsonConfig = []byte(`{
	"url": "ldap://138.91.247.105",
	"userdn": "example,com",
	"binddn": "kitty",
	"bindpass": "cats",
	"tls_max_version": "tls12",
	"tls_min_version": "tls12",
	"request_timeout": 30,
	"connection_timeout": 15,
	"ClientTLSCert":  "",
	"ClientTLSKey":   "",
	"schema": "openldap"
}`)

var jsonConfigDefault = []byte(`
{
  "url": "ldap://127.0.0.1",
  "userdn": "",
  "anonymous_group_search": false,
  "groupdn": "",
  "groupfilter": "(|(memberUid={{.Username}})(member={{.UserDN}})(uniqueMember={{.UserDN}}))",
  "groupattr": "cn",
  "upndomain": "",
  "userattr": "cn",
  "userfilter": "({{.UserAttr}}={{.Username}})",
  "certificate": "",
  "client_tls_cert": "",
  "client_tsl_key": "",
  "insecure_tls": false,
  "starttls": false,
  "binddn": "",
  "bindpass": "",
  "deny_null_bind": true,
  "discoverdn": false,
  "tls_min_version": "tls12",
  "tls_max_version": "tls12",
  "use_token_groups": false,
  "use_pre111_group_cn_behavior": null,
  "username_as_alias": false,
  "request_timeout": 90,
  "connection_timeout": 30,
  "dereference_aliases": "never",
  "max_page_size": 0,
  "CaseSensitiveNames": false,
  "ClientTLSCert": "",
  "ClientTLSKey": "",
  "enable_samaccountname_login": false,
  "schema": "openldap"
}
`)
