// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"encoding/json"
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
	"ClientTLSKey":   ""
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
  "enable_samaccountname_login": false
}
`)
