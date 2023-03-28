// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package command

import (
	"net/http"
	"reflect"
	"testing"
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
