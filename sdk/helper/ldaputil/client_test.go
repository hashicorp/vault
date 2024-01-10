// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ldaputil

import (
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDialLDAP duplicates a potential panic that was
// present in the previous version of TestDialLDAP,
// then confirms its fix by passing.
func TestDialLDAP(t *testing.T) {
	ldapClient := Client{
		Logger: hclog.NewNullLogger(),
		LDAP:   NewLDAP(),
	}

	ce := &ConfigEntry{
		Url:            "ldap://localhost:384654786",
		RequestTimeout: 3,
	}
	if _, err := ldapClient.DialLDAP(ce); err == nil {
		t.Fatal("expected error")
	}
}

func TestLDAPEscape(t *testing.T) {
	testcases := map[string]string{
		"#test":                      "\\#test",
		"test,hello":                 "test\\,hello",
		"test,hel+lo":                "test\\,hel\\+lo",
		"test\\hello":                "test\\\\hello",
		"  test  ":                   "\\  test \\ ",
		"":                           "",
		`\`:                          `\\`,
		"trailing\000":               `trailing\00`,
		"mid\000dle":                 `mid\00dle`,
		"\000":                       `\00`,
		"multiple\000\000":           `multiple\00\00`,
		"backlash-before-null\\\000": `backlash-before-null\\\00`,
		"trailing\\":                 `trailing\\`,
		"double-escaping\\>":         `double-escaping\\\>`,
	}

	for test, answer := range testcases {
		res := EscapeLDAPValue(test)
		if res != answer {
			t.Errorf("Failed to escape %s: %s != %s\n", test, res, answer)
		}
	}
}

func TestGetTLSConfigs(t *testing.T) {
	config := testConfig(t)
	if err := config.Validate(); err != nil {
		t.Fatal(err)
	}
	tlsConfig, err := getTLSConfig(config, "138.91.247.105")
	if err != nil {
		t.Fatal(err)
	}
	if tlsConfig == nil {
		t.Fatal("expected 1 TLS config because there's 1 url")
	}
	if tlsConfig.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify should be false because we should default to the most secure connection")
	}
	if tlsConfig.ServerName != "138.91.247.105" {
		t.Fatalf("expected ServerName of \"138.91.247.105\" but received %q", tlsConfig.ServerName)
	}
	expected := uint16(771)
	if tlsConfig.MinVersion != expected || tlsConfig.MaxVersion != expected {
		t.Fatal("expected TLS min and max version of 771 which corresponds with TLS 1.2 since TLS 1.1 and 1.0 have known vulnerabilities")
	}
}

func TestSIDBytesToString(t *testing.T) {
	testcases := map[string][]byte{
		"S-1-5-21-2127521184-1604012920-1887927527-72713": {0x01, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05, 0x15, 0x00, 0x00, 0x00, 0xA0, 0x65, 0xCF, 0x7E, 0x78, 0x4B, 0x9B, 0x5F, 0xE7, 0x7C, 0x87, 0x70, 0x09, 0x1C, 0x01, 0x00},
		"S-1-1-0": {0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
		"S-1-5":   {0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x05},
	}

	for answer, test := range testcases {
		res, err := sidBytesToString(test)
		if err != nil {
			t.Errorf("Failed to conver %#v: %s", test, err)
		} else if answer != res {
			t.Errorf("Failed to convert %#v: %s != %s", test, res, answer)
		}
	}
}

func TestClient_renderUserSearchFilter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		conf        *ConfigEntry
		username    string
		want        string
		errContains string
	}{
		{
			name:     "valid-default",
			username: "alice",
			conf: &ConfigEntry{
				UserAttr: "cn",
			},
			want: "(cn=alice)",
		},
		{
			name:     "escaped-malicious-filter",
			username: "foo@example.com)((((((((((((((((((((((((((((((((((((((userPrincipalName=foo",
			conf: &ConfigEntry{
				UPNDomain:  "example.com",
				UserFilter: "(&({{.UserAttr}}={{.Username}})({{.UserAttr}}=admin@example.com))",
			},
			want: "(&(userPrincipalName=foo@example.com\\29\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28\\28userPrincipalName=foo@example.com)(userPrincipalName=admin@example.com))",
		},
		{
			name:     "bad-filter-unclosed-action",
			username: "alice",
			conf: &ConfigEntry{
				UserFilter: "hello{{range",
			},
			errContains: "search failed due to template compilation error",
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := Client{
				Logger: hclog.NewNullLogger(),
				LDAP:   NewLDAP(),
			}

			f, err := c.RenderUserSearchFilter(tc.conf, tc.username)
			if tc.errContains != "" {
				require.Error(t, err)
				assert.ErrorContains(t, err, tc.errContains)
				return
			}
			require.NoError(t, err)
			assert.NotEmpty(t, f)
			assert.Equal(t, tc.want, f)
		})
	}
}
