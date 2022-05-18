package command

import (
	"net/http"
	"reflect"
	"testing"

	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/stretchr/testify/assert"
)

func getDefaultCliHeaders(t *testing.T) http.Header {
	bc := &BaseCommand{}
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}
	return cli.Headers()
}

func TestClient_NamespaceFlagSetToValue(t *testing.T) {
	bc := &BaseCommand{}
	bc.flagSet(FlagSetHTTP)
	bc.flagNamespace = "juan"
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}

	h := cli.Headers()
	v := h.Get(consts.NamespaceHeaderName)

	// Expect the `namespace.Canonicalize` method will add a trailing slash
	assert.Equal(t, "juan/", v, "Expected namespace not found")
}

func TestClient_NSFlagSetToValue(t *testing.T) {
	bc := &BaseCommand{}
	bc.flagSet(FlagSetHTTP)
	bc.flagNS = "juan"
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}

	h := cli.Headers()
	v := h.Get(consts.NamespaceHeaderName)

	assert.Equal(t, "juan/", v, "Expected ns not found")
}

func TestClient_NSFlagSetToValueAndNamespaceFlagSetToDifferentValue(t *testing.T) {
	bc := &BaseCommand{}
	bc.flagSet(FlagSetHTTP)
	bc.flagNS = "juan"
	bc.flagNamespace = "john"
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}

	h := cli.Headers()
	v := h.Get(consts.NamespaceHeaderName)

	assert.Equal(t, "juan/", v, "Expected ns not found")
}

func TestClient_NSFlagSetToEmptyString(t *testing.T) {
	bc := &BaseCommand{}
	bc.flagSet(FlagSetHTTP)
	bc.flagNS = ""
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}

	h := cli.Headers()
	v := h.Get(consts.NamespaceHeaderName)

	assert.Equal(t, "", v, "Expected ns was configured rather than ignored")
}

func TestClient_NamespaceFlagSetToEmptyString(t *testing.T) {
	bc := &BaseCommand{}
	bc.flagSet(FlagSetHTTP)
	bc.flagNamespace = ""
	cli, err := bc.Client()
	if err != nil {
		t.Fatal(err)
	}

	h := cli.Headers()
	v := h.Get(consts.NamespaceHeaderName)

	assert.Equal(t, "", v, "Expected namespace was configured rather than ignored")
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
