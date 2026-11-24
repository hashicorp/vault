// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestClientNaming_Usage checks that help uses the right binary name without
// changing exit codes for typos or commands run without a subcommand.
func TestClientNaming_Usage(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"vault", "vault-client"} {
		for _, tc := range []struct {
			args   []string
			code   int
			header string
		}{
			{[]string{"--help"}, 0, " <command> [args]"},
			{[]string{"statu"}, 127, " <command> [args]"},
			{[]string{"auth", "--help"}, 0, " auth <subcommand> [options] [args]"},
			{[]string{"auth", "kist"}, 1, " auth <subcommand> [options] [args]"},
			{[]string{"operator", "--help"}, 0, " operator <subcommand> [options] [args]"},
			{[]string{"operator"}, 1, " operator <subcommand> [options] [args]"},
		} {
			t.Run(name+"/"+strings.Join(tc.args, " "), func(t *testing.T) {
				t.Parallel()
				var stdout, stderr bytes.Buffer
				code := RunCustom(tc.args, &base.RunOptions{
					BinaryName: name,
					Stdin:      strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
				})
				require.Equal(t, tc.code, code, "help/fallback status")
				out := stdout.String() + stderr.String()
				require.Contains(t, out, "Usage: "+name+tc.header, "configured executable")
				if tc.args[0] == "auth" {
					for _, example := range []string{"auth list", "auth enable userpass", "auth help github"} {
						require.Contains(t, out, "$ "+name+" "+example, "auth invocation example")
					}
				}
			})
		}
	}
}

// TestClientNaming_CommandHelp checks binary names in each registered command's
// help so we don't have to keep a separate list for CE and enterprise.
func TestClientNaming_CommandHelp(t *testing.T) {
	t.Parallel()
	bareInvocation := regexp.MustCompile(`(^|[\s"'` + "`" + `|])vault\s+(?:auth|login|token|tpm|hcp|kv|read|write|patch|unwrap|agent|proxy)\b`)
	for _, name := range []string{"vault", "vault-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ui := cli.NewMockUi()
			serverUI := cli.NewMockUi()
			for command, factory := range InitClientCommands(ui, serverUI, &base.RunOptions{BinaryName: name}) {
				cmd, err := factory()
				require.NoError(t, err, command)
				help := cmd.Help()
				switch c := cmd.(type) {
				case *AgentCommand:
					require.Same(t, serverUI, c.UI, "Agent retains its special UI")
				case *ProxyCommand:
					require.Same(t, serverUI, c.UI, "Proxy retains its special UI")
				}
				require.Contains(t, help, "Usage: "+name+" "+command, "command %s", command)
				if name == "vault-client" {
					require.False(t, bareInvocation.MatchString(help), "bare executable in %s help: %s", command, bareInvocation.FindString(help))
				}
			}
		})
	}
}

// TestClientNaming_AuthMethodHelp checks that `auth help TYPE` shows the right
// binary name for every login handler, whether it lives in this repo or comes
// from an auth plugin. Other names change nothing but the binary name, and
// vault's help stays exactly as the handler wrote it.
func TestClientNaming_AuthMethodHelp(t *testing.T) {
	t.Parallel()
	bareInvocation := regexp.MustCompile(`(^|[\s"'` + "`" + `|])vault\s+(?:auth|login|token|tpm|hcp|kv|read|write|patch|unwrap|agent|proxy)\b`)
	for _, name := range []string{"", "vault", "vault-client", "renamed-vault"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			for method, original := range LoginHandlers("vault") {
				ui := cli.NewMockUi()
				client, err := api.NewClient(api.DefaultConfig())
				require.NoError(t, err, method)
				commands := InitClientCommands(ui, ui, &base.RunOptions{
					BinaryName: name, Client: client, HCPTokenHelper: testHCPTokenHelper{},
				})
				authHelp, err := commands["auth help"]()
				require.NoError(t, err, method)
				require.Zero(t, authHelp.Run([]string{method}), "method %s", method)
				require.Empty(t, ui.ErrorWriter.String(), "method %s", method)
				help := ui.OutputWriter.String()
				if name == "" || name == "vault" {
					require.Equal(t, original.Help()+"\n", help, "method %s keeps its original help", method)
					require.Contains(t, help, "Usage: vault login", "method %s", method)
				} else {
					require.Equal(t, original.Help()+"\n", strings.ReplaceAll(help, name, "vault"), "method %s only changes the binary name", method)
					require.Contains(t, help, name+" login", "method %s", method)
					require.False(t, bareInvocation.MatchString(help), "bare executable in %s help: %s", method, bareInvocation.FindString(help))
				}
			}
		})
	}
}

// TestClientNaming_UnknownAuthMethod checks that we suggest the right auth list
// command without reading credentials, making requests, or printing extra arguments.
func TestClientNaming_UnknownAuthMethod(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"vault", "vault-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			var stdout, stderr bytes.Buffer
			helper := &namingTokenHelper{}
			// Nothing listens on port 1, so a request would fail the login
			// with a different error.
			code := RunCustom([]string{"login", "-method=not-a-real-method", "password=disposable-secret"}, &base.RunOptions{
				BinaryName: name, Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr, TokenHelper: helper, Address: "http://127.0.0.1:1",
			})
			require.Equal(t, 1, code, "unknown method status")
			require.Empty(t, stdout.String(), "error stays on stderr")
			require.Contains(t, stderr.String(), `Use "`+name+` auth list"`, "method-list suggestion")
			require.NotContains(t, stderr.String(), "disposable-secret", "extra arguments stay private")
			require.Zero(t, helper.calls, "no token helper operations")
		})
	}
}

// TestClientNaming_LoginSuccess checks that the login message uses the right
// binary name and that we still save the token and print the usual success message.
func TestClientNaming_LoginSuccess(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"vault", "vault-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			client, err := api.NewClient(api.DefaultConfig())
			require.NoError(t, err, "in-process API client")
			ui := cli.NewMockUi()
			helper := &namingTokenHelper{}
			cmd := &LoginCommand{BaseCommand: base.NewBaseCommand(&base.VaultUI{Ui: ui, Format: "table"}, &base.RunOptions{
				BinaryName: name, Client: client, TokenHelper: helper, HCPTokenHelper: testHCPTokenHelper{}, Stdin: strings.NewReader(""),
			}), Handlers: map[string]base.LoginHandler{"naming": namingLoginHandler{}}}
			require.Zero(t, cmd.Run([]string{"-method=naming"}), "login status")
			require.Equal(t, "disposable-test-token", helper.token, "authentication still stores the returned token")
			require.Equal(t, 1, helper.calls, "only the token store is needed")
			require.Contains(t, strings.Join(strings.Fields(ui.OutputWriter.String()), " "), `"`+name+` login" again.`, "next-login instruction")
		})
	}
}

type namingLoginHandler struct{}

func (namingLoginHandler) Help() string { return "" }
func (namingLoginHandler) Auth(*api.Client, map[string]string) (*api.Secret, error) {
	return &api.Secret{Auth: &api.SecretAuth{ClientToken: "disposable-test-token"}}, nil
}

type namingTokenHelper struct {
	calls int
	token string
}

func (h *namingTokenHelper) Path() string             { h.calls++; return "" }
func (h *namingTokenHelper) Get() (string, error)     { h.calls++; return h.token, nil }
func (h *namingTokenHelper) Erase() error             { h.calls++; h.token = ""; return nil }
func (h *namingTokenHelper) Store(token string) error { h.calls++; h.token = token; return nil }
