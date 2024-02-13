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
	defaultHeaders := getDefaultCliHeaders(t)

	cases := []struct {
		Input map[string]string
		Valid bool
	}{
		{
			map[string]string{},
			true,
		},
		{
			map[string]string{"foo": "bar", "header2": "value2"},
			true,
		},
		{
			map[string]string{"X-Vault-foo": "bar", "header2": "value2"},
			false,
		},
	}

	for _, tc := range cases {
		expectedHeaders := defaultHeaders.Clone()
		for key, val := range tc.Input {
			expectedHeaders.Add(key, val)
		}

		bc := &BaseCommand{flagHeader: tc.Input}
		cli, err := bc.Client()

		if err == nil && !tc.Valid {
			t.Errorf("No error for input[%#v], but not valid", tc.Input)
			continue
		}

		if err != nil {
			if tc.Valid {
				t.Errorf("Error[%v] with input[%#v], but valid", err, tc.Input)
			}
			continue
		}

		if cli == nil {
			t.Error("client should not be nil")
		}

		actualHeaders := cli.Headers()
		if !reflect.DeepEqual(expectedHeaders, actualHeaders) {
			t.Errorf("expected [%#v] but got [%#v]", expectedHeaders, actualHeaders)
		}
	}
}

// TestClient_HCPConfiguration tests that the HCP configuration is applied correctly when it exists in cache.
func TestClient_HCPConfiguration(t *testing.T) {
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
			bc := &BaseCommand{hcpTokenHelper: &hcpvlib.TestingHCPTokenHelper{tst.Valid}}
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
