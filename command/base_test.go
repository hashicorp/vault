// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"net/http"
	"reflect"
	"testing"

	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getDefaultCliHeaders(t *testing.T) http.Header {
	bc := &BaseCommand{}
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}
	return cli.Headers()
}

func TestClient_FlagHeader(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Input map[string]string
		Valid bool
	}{
		"empty": {
			map[string]string{},
			true,
		},
		"valid": {
			map[string]string{"foo": "bar", "header2": "value2"},
			true,
		},
		"invalid": {
			map[string]string{"X-Vault-foo": "bar", "header2": "value2"},
			false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			expectedHeaders := getDefaultCliHeaders(t)
			for key, val := range tc.Input {
				expectedHeaders.Add(key, val)
			}

			bc := &BaseCommand{flagHeader: tc.Input}
			cli, err := bc.Client()

			if err == nil && !tc.Valid {
				t.Errorf("No error for input[%#v], but not valid", tc.Input)
			}

			if err != nil {
				if tc.Valid {
					t.Errorf("Error[%v] with input[%#v], but valid", err, tc.Input)
				}
				return
			}

			if cli == nil {
				t.Error("client should not be nil")
			}

			actualHeaders := cli.Headers()
			if !reflect.DeepEqual(expectedHeaders, actualHeaders) {
				t.Errorf("expected [%#v] but got [%#v]", expectedHeaders, actualHeaders)
			}
		})
	}
}

// TestClient_HCPConfiguration tests that the HCP configuration is applied correctly when it exists in cache.
func TestClient_HCPConfiguration(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		Valid        bool
		ExpectedAddr string
	}{
		"valid hcp configuration": {
			Valid:        true,
			ExpectedAddr: "https://hcp-proxy.addr:8200",
		},
		"empty hcp configuration": {
			Valid:        false,
			ExpectedAddr: api.DefaultAddress,
		},
	}

	for n, tst := range cases {
		t.Run(n, func(t *testing.T) {
			t.Parallel()
			bc := &BaseCommand{hcpTokenHelper: &hcpvlib.TestingHCPTokenHelper{ValidCache: tst.Valid}}
			cli, err := bc.Client()
			assert.NoError(t, err)

			if tst.Valid {
				require.Equal(t, tst.ExpectedAddr, cli.Address())
				require.NotEmpty(t, cli.HCPCookie())
				require.Contains(t, cli.HCPCookie(), "hcp_access_token=Test.Access.Token")
			} else {
				require.Equal(t, tst.ExpectedAddr, cli.Address())
				require.Empty(t, cli.HCPCookie())
			}
		})
	}
}

// Test_FlagSet_StringVar_Normalizers verifies that the normalizer callbacks
// works as expected.
func Test_FlagSet_StringVar_Normalizers(t *testing.T) {
	appendA := func(in string) string { return in + "a" }
	prependB := func(in string) string { return "b" + in }

	for name, test := range map[string]struct {
		in            func() *StringVar
		envVars       map[string]string
		expectedValue string
	}{
		"no normalizers no env vars uses default value": {
			in: func() *StringVar {
				resT := ""
				return &StringVar{
					Name:    "test",
					Target:  &resT,
					EnvVar:  "VAULT_TEST",
					Default: "default",
				}
			},
			expectedValue: "default",
		},
		"one normalizer no env vars normalizes default value": {
			in: func() *StringVar {
				resT := ""
				return &StringVar{
					Name:        "test",
					Target:      &resT,
					EnvVar:      "VAULT_TEST",
					Default:     "default",
					Normalizers: []func(string) string{appendA},
				}
			},
			expectedValue: "defaulta",
		},
		"two normalizers no env vars normalizes default value with both": {
			in: func() *StringVar {
				resT := ""
				return &StringVar{
					Name:        "test",
					Target:      &resT,
					EnvVar:      "VAULT_TEST",
					Default:     "default",
					Normalizers: []func(string) string{appendA, prependB},
				}
			},
			expectedValue: "bdefaulta",
		},
		"two normalizers with env vars normalizes env var value with both": {
			in: func() *StringVar {
				resT := ""
				return &StringVar{
					Name:        "test",
					Target:      &resT,
					EnvVar:      "VAULT_TEST",
					Default:     "default",
					Normalizers: []func(string) string{appendA, prependB},
				}
			},
			envVars:       map[string]string{"VAULT_TEST": "env_override"},
			expectedValue: "benv_overridea",
		},
	} {
		t.Run(name, func(t *testing.T) {
			for k, v := range test.envVars {
				t.Setenv(k, v)
			}
			fsets := NewFlagSets(nil)
			fs := fsets.NewFlagSet("test")
			sv := test.in()
			fs.StringVar(sv)
			require.Equal(t, test.expectedValue, *sv.Target)
		})
	}
}
