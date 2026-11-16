// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// TestClientNaming_BaseHelp checks that shared flag help uses the name we were
// given and falls back to vault when no name is set.
func TestClientNaming_BaseHelp(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", "vault", "vault-client", "renamed-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			c := NewBaseCommand(cli.NewMockUi(), &RunOptions{BinaryName: name})
			if name == "" {
				name = "vault"
			}
			require.Equal(t, name, c.BinaryName, "constructor identity")
			help := strings.Join(strings.Fields(c.FlagSetForBit(FlagSetHTTP|FlagSetOutputFormat).Help()), " ")
			require.Contains(t, help, `"`+name+` unwrap"`, "wrapping instruction")
			require.Contains(t, help, "'"+name+" read'", "raw format instruction")
		})
	}
	require.Equal(t, "vault", NewBaseCommand(cli.NewMockUi(), nil).BinaryName, "nil options retain server identity")
}

// TestClientNaming_CurlString checks that -output-curl-string reads the token
// with the binary we were run as, because vault-client hosts may not have vault.
// vault's output must stay byte-identical, and request data that contains the
// same placeholder text must be printed as-is.
func TestClientNaming_CurlString(t *testing.T) {
	t.Parallel()
	body := `{"note":"-H \"X-Vault-Token: $(vault print token)\" $(vault print token)"}`
	options := `--insecure -X POST --cacert '/tmp/ca.pem' `
	rest := `-d '` + body + `' https://example.test/v1/secret/data/app` + "\n"
	for _, tc := range []struct {
		name string
		want string
	}{
		{"", "curl " + options + `-H "X-Vault-Token: $(vault print token)" ` + rest},
		{"vault", "curl " + options + `-H "X-Vault-Token: $(vault print token)" ` + rest},
		{"vault-client", `curl -H "X-Vault-Token: $(vault-client print token)" ` + options + rest},
		{"renamed-client", `curl -H "X-Vault-Token: $(renamed-client print token)" ` + options + rest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			request, err := retryablehttp.NewRequest(http.MethodPost, "https://example.test/v1/secret/data/app", []byte(body))
			require.NoError(t, err, "request")
			request.Header.Set(api.AuthHeaderName, "disposable-secret-token")
			current := &api.OutputStringError{Request: request, TLSSkipVerify: true, ClientCACert: "/tmp/ca.pem"}
			// Commands usually print the error first, which caches the vault rendering.
			require.Equal(t, api.ErrOutputStringRequest, current.Error(), "cached rendering")

			var stdout, stderr bytes.Buffer
			code := generateCurlString(2, &RunOptions{BinaryName: tc.name, Stdout: &stdout, Stderr: &stderr},
				&bytes.Buffer{}, nil, current)
			require.Zero(t, code, "curl status")
			require.Empty(t, stderr.String(), "no curl errors")
			require.Equal(t, tc.want, stdout.String(), "curl command")
			if tc.name == "" || tc.name == "vault" {
				cs, err := current.CurlString()
				require.NoError(t, err, "api rendering")
				require.Equal(t, cs+"\n", stdout.String(), "vault prints the api rendering unchanged")
			}
			require.NotContains(t, stdout.String(), "disposable-secret-token", "token stays private")
			require.Equal(t, "disposable-secret-token", request.Header.Get(api.AuthHeaderName), "original request is untouched")
		})
	}
}

// TestClientNaming_CurlStringWithoutToken checks that requests without a token
// header print the same curl command for every binary name.
func TestClientNaming_CurlStringWithoutToken(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", "vault", "vault-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			request, err := retryablehttp.NewRequest(http.MethodGet, "https://example.test/v1/sys/health", nil)
			require.NoError(t, err, "request")
			var stdout, stderr bytes.Buffer
			code := generateCurlString(2, &RunOptions{BinaryName: name, Stdout: &stdout, Stderr: &stderr},
				&bytes.Buffer{}, nil, &api.OutputStringError{Request: request})
			require.Zero(t, code, "curl status")
			require.Equal(t, "curl https://example.test/v1/sys/health\n", stdout.String(), "curl command")
		})
	}
}

// TestClientNaming_CurlStringFields checks that curlString knows every exported
// field of api.OutputStringError. curlString copies them by hand, so a field
// added to the api module would otherwise be dropped from vault-client output.
func TestClientNaming_CurlStringFields(t *testing.T) {
	t.Parallel()
	var fields []string
	typ := reflect.TypeOf(api.OutputStringError{})
	for i := range typ.NumField() {
		if field := typ.Field(i); field.IsExported() {
			fields = append(fields, field.Name)
		}
	}
	require.ElementsMatch(t, []string{"Request", "TLSSkipVerify", "ClientCACert", "ClientCAPath", "ClientCert", "ClientKey"}, fields,
		"api.OutputStringError fields changed: update curlString in runner.go to copy them, then update this list")
}

// TestClientNaming_RunnerRaw checks that raw output errors use this command's
// name without changing the shared formatter or the secret bytes we print.
func TestClientNaming_RunnerRaw(t *testing.T) {
	t.Parallel()
	for _, ui := range []cli.Ui{cli.NewMockUi(), &VaultUI{Ui: cli.NewMockUi(), Format: "raw"}} {
		require.ErrorContains(t, (RawFormatter{}).Output(ui, nil, "unsupported"), "only `vault read` does.", "bare/default UIs retain server diagnostic")
	}
	for _, name := range []string{"", "vault", "vault-client", "renamed-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			var commandUI, serverUI cli.Ui
			code := RunCustomWithCommandsFn([]string{"--help", "-format=raw"}, &RunOptions{
				BinaryName: name, Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
			}, func(ui, serverCmdUI cli.Ui, _ *RunOptions) map[string]cli.CommandFactory {
				commandUI, serverUI = ui, serverCmdUI
				return map[string]cli.CommandFactory{}
			})
			if name == "" {
				name = filepath.Base(os.Args[0])
			}
			require.Zero(t, code, "root help status")
			require.Contains(t, stdout.String(), "Usage: "+name+" <command> [args]", "runtime/explicit identity")
			stdout.Reset()
			require.Equal(t, 1, OutputWithFormat(commandUI, nil, "unsupported"), "unsupported raw data status")
			require.Contains(t, stderr.String(), "only `"+name+" read` does.", "per-run raw diagnostic")
			require.Empty(t, stdout.String(), "diagnostic stays on stderr")
			_, serverErr := Formatters["raw"].Format("unsupported")
			require.EqualError(t, serverErr, "This command does not support the -format=raw option; only `vault read` does.", "global formatter remains server-compatible")
			err := Formatters["raw"].Output(serverUI, nil, "unsupported")
			require.ErrorContains(t, err, "only `"+name+" read` does.", "server UI receives the same run identity")
			data := []byte("vault login fake-token https://example/vault/\x00\xff")
			require.Zero(t, OutputWithFormat(commandUI, nil, data), "raw bytes status")
			require.Equal(t, append(append([]byte{}, data...), '\n'), stdout.Bytes(), "only the existing UI newline is added")
		})
	}
}
