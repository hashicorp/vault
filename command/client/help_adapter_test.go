// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	base "github.com/hashicorp/vault/command/base"
	"github.com/stretchr/testify/require"
)

// TestClientNaming_HCPHelp checks that we only change the binary name in HCP help
// and keep the parent commands' flags, synopsis, and optional CLI interfaces.
func TestClientNaming_HCPHelp(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", "vault", "vault-client", "renamed-client", "client-%s-{{.Name}}"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			connectUI, disconnectUI := cli.NewMockUi(), cli.NewMockUi()
			upstreamConnect := &hcpvlib.HCPConnectCommand{Ui: connectUI}
			upstreamDisconnect := &hcpvlib.HCPDisconnectCommand{Ui: disconnectUI}
			connect := &hcpConnectCommand{HCPConnectCommand: upstreamConnect, BinaryName: name}
			disconnect := &hcpDisconnectCommand{HCPDisconnectCommand: upstreamDisconnect, BinaryName: name}
			require.Same(t, upstreamConnect, connect.HCPConnectCommand, "connect retains its concrete receiver")
			require.Same(t, upstreamDisconnect, disconnect.HCPDisconnectCommand, "disconnect retains its concrete receiver")
			require.Same(t, connectUI, connect.Ui, "connect UI is not intercepted")
			require.Same(t, disconnectUI, disconnect.Ui, "disconnect UI is not intercepted")
			flags := connect.Flags()
			for _, flag := range []string{"client-id", "secret-id", "organization-id", "project-id", "cluster-id"} {
				require.NotNil(t, flags.Lookup(flag), "promoted flag %s", flag)
			}
			binaryName := name
			if binaryName == "" {
				binaryName = "vault"
			}
			for _, tc := range []struct {
				key      string
				upstream cli.Command
				wrapped  cli.Command
				examples []string
				prose    []string
			}{
				{
					key: "hcp connect", upstream: upstreamConnect, wrapped: connect,
					examples: []string{
						"      $ " + binaryName + " hcp connect -client-id=client-id-value -secret-id=secret-id-value",
						"      $ " + binaryName + " hcp connect -client-id=client-id-value -secret-id=secret-id-value -organization-id=org-UUID -project-id=proj-UUID -cluster-id=cluster-name",
					},
					prose: []string{
						"  Authenticates users or machines to HCP using either provided arguments or retrieved token through\n  browser login. A successful authentication results in an HCP token and an HCP Vault address being\n  locally cached. ",
						"  The default authentication method is an interactive one, redirecting users to the HCP login browser.\n  If a set of service principal credential is supplied, which is generated through the HCP portal with \n  the necessary capabilities to access the organization, project, and HCP Vault cluster chosen.",
						"  Additionally, the organization identification, project identification, and cluster name can be passed in to\n  directly connect to a specific HCP Vault cluster without interacting with the CLI.",
					},
				},
				{
					key: "hcp disconnect", upstream: upstreamDisconnect, wrapped: disconnect,
					examples: []string{"      $ " + binaryName + " hcp disconnect"},
					prose:    []string{"  Cleans up the cache with the HCP credentials used to connect to a HCP Vault cluster. "},
				},
			} {
				help := tc.wrapped.Help()
				if name == "" || name == "vault" {
					require.Equal(t, tc.upstream.Help(), help, "%s default help must remain byte-identical", tc.key)
				} else {
					require.Equal(t, tc.upstream.Help(), strings.ReplaceAll(help, name, "vault"), "%s help only changes the binary name", tc.key)
				}
				require.True(t, strings.HasPrefix(help, "Usage: "+binaryName+" "+tc.key+" [options]\n"), "%s usage names the executable literally", tc.key)
				for _, example := range tc.examples {
					require.Contains(t, strings.Split(help, "\n"), example, "%s retains example flags and values", tc.key)
				}
				for _, prose := range tc.prose {
					require.Contains(t, help, prose, "%s retains upstream prose", tc.key)
				}
				require.Equal(t, tc.upstream.Synopsis(), tc.wrapped.Synopsis(), "%s synopsis is promoted", tc.key)
				_, upstreamAutocomplete := tc.upstream.(cli.CommandAutocomplete)
				_, wrappedAutocomplete := tc.wrapped.(cli.CommandAutocomplete)
				require.Equal(t, upstreamAutocomplete, wrappedAutocomplete, "%s autocomplete capability is preserved", tc.key)
				_, upstreamTemplate := tc.upstream.(cli.CommandHelpTemplate)
				_, wrappedTemplate := tc.wrapped.(cli.CommandHelpTemplate)
				require.Equal(t, upstreamTemplate, wrappedTemplate, "%s help-template capability is preserved", tc.key)
				require.False(t, upstreamTemplate, "%s gained HelpTemplate: review whether it bypasses owned named help", tc.key)
			}
		})
	}
}

// TestClientNaming_HCPRegistry checks that we wrap every HCP command that needs
// a different binary name, and that hcp connect prints through the renaming
// Ui, but leave the original types and Ui alone for vault.
func TestClientNaming_HCPRegistry(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", "vault", "vault-client", "renamed-client"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ui := cli.NewMockUi()
			commands := map[string]cli.CommandFactory{}
			InitHCPCommands(ui, commands, name)
			factories := hcpvlib.InitHCPCommand(ui)
			require.Len(t, commands, len(factories), "registration set is unchanged")
			for key, factory := range factories {
				require.Contains(t, commands, key, "upstream registration is retained")
				upstream, err := factory()
				require.NoError(t, err, key)
				command, err := commands[key]()
				require.NoError(t, err, key)
				if name == "" || name == "vault" {
					require.IsType(t, upstream, command, "%s default command is unwrapped", key)
					require.Equal(t, upstream.Help(), command.Help(), "%s default help", key)
				} else {
					switch upstream.(type) {
					case *hcpvlib.HCPConnectCommand:
						require.IsType(t, &hcpConnectCommand{}, command, "%s requires a concrete naming wrapper", key)
						wrapped := command.(*hcpConnectCommand)
						require.NotNil(t, wrapped.HCPConnectCommand, "%s retains the upstream receiver", key)
						require.Equal(t, name, wrapped.BinaryName, "%s configured name", key)
					case *hcpvlib.HCPDisconnectCommand:
						require.IsType(t, &hcpDisconnectCommand{}, command, "%s requires a concrete naming wrapper", key)
						wrapped := command.(*hcpDisconnectCommand)
						require.NotNil(t, wrapped.HCPDisconnectCommand, "%s retains the upstream receiver", key)
						require.Equal(t, name, wrapped.BinaryName, "%s configured name", key)
					default:
						require.FailNow(t, "new upstream HCP command requires naming review", "%s: %T", key, upstream)
					}
				}
				switch c := command.(type) {
				case *hcpvlib.HCPConnectCommand:
					require.Same(t, ui, c.Ui, "%s original UI", key)
				case *hcpvlib.HCPDisconnectCommand:
					require.Same(t, ui, c.Ui, "%s original UI", key)
				case *hcpConnectCommand:
					require.Equal(t, &hcpConnectUi{Ui: ui, BinaryName: name}, c.Ui, "%s renaming UI", key)
				case *hcpDisconnectCommand:
					require.Same(t, ui, c.Ui, "%s original UI", key)
				}
			}
			require.Empty(t, ui.OutputWriter.String(), "construction emits no output")
			require.Empty(t, ui.ErrorWriter.String(), "construction emits no error")
		})
	}
}

// TestClientNaming_HCPCollision checks that we report a name conflict without
// replacing an existing command, no matter which HCP command is registered first.
func TestClientNaming_HCPCollision(t *testing.T) {
	t.Parallel()
	for key := range hcpvlib.InitHCPCommand(cli.NewMockUi()) {
		t.Run(key, func(t *testing.T) {
			t.Parallel()
			ui := cli.NewMockUi()
			existing := &cli.MockCommand{}
			commands := map[string]cli.CommandFactory{key: func() (cli.Command, error) { return existing, nil }}
			InitHCPCommands(ui, commands, "vault")
			command, err := commands[key]()
			require.NoError(t, err, "existing factory")
			require.Same(t, existing, command, "colliding registration is not overwritten")
			require.Equal(t, "Failed to initialize HCP commands.\n", ui.ErrorWriter.String(), "existing collision diagnostic")
			require.Empty(t, ui.OutputWriter.String(), "collision stays on stderr")
		})
	}
}

// TestClientNaming_HCPCLIHelp checks that HCP help uses the binary name we run
// with. Asking for help shouldn't contact HCP or touch its cache.
func TestClientNaming_HCPCLIHelp(t *testing.T) {
	t.Parallel()
	for _, name := range []string{"", "vault", "vault-client", "renamed-client"} {
		for key := range hcpvlib.InitHCPCommand(cli.NewMockUi()) {
			t.Run(name+"/"+key, func(t *testing.T) {
				t.Parallel()
				var stdout, stderr bytes.Buffer
				args := append(strings.Fields(key), "--help")
				code := RunCustom(args, &base.RunOptions{
					BinaryName: name, Stdin: strings.NewReader(""), Stdout: &stdout, Stderr: &stderr,
				})
				require.Zero(t, code, "CLI help status")
				require.Empty(t, stderr.String(), "CLI help stays on stdout")
				binaryName := name
				if binaryName == "" {
					binaryName = filepath.Base(os.Args[0])
				}
				commands := map[string]cli.CommandFactory{}
				InitHCPCommands(cli.NewMockUi(), commands, binaryName)
				command, err := commands[key]()
				require.NoError(t, err, "expected HCP factory")
				require.Equal(t, command.Help()+"\n", stdout.String(), "CLI renders the registered help unchanged")
				require.True(t, strings.HasPrefix(stdout.String(), "Usage: "+binaryName+" "+key+" [options]\n"), "CLI uses the effective executable name")
				require.Contains(t, stdout.String(), "$ "+binaryName+" "+key, "CLI example uses the effective executable name")
			})
		}
	}
}

// TestClientNaming_RenameVaultInvocations checks both renaming rules. Help
// renames "vault login" and "vault hcp" wherever vault starts a word. hcp
// connect output renames only quoted ones, because it also prints HCP org and
// project names that users type back. Lookalikes never change.
func TestClientNaming_RenameVaultInvocations(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		text   string
		help   string
		output string
	}{
		{"auth/vault login", "auth/vault login", "auth/vault login"},
		{"myvault login", "myvault login", "myvault login"},
		{"vault logins", "vault logins", "vault logins"},
		{"vault hcpx", "vault hcpx", "vault hcpx"},
		{"Vault login", "Vault login", "Vault login"},
		{"x-vault-aws", "x-vault-aws", "x-vault-aws"},
		{"'vault logins' `vault hcpx`", "'vault logins' `vault hcpx`", "'vault logins' `vault hcpx`"},
		{"'vault login' 'vault login'", "'vault-client login' 'vault-client login'", "'vault-client login' 'vault-client login'"},
		{"vault login vault login", "vault-client login vault-client login", "vault login vault login"},
		{"\t$ vault hcp connect", "\t$ vault-client hcp connect", "\t$ vault hcp connect"},
		{"then run vault login", "then run vault-client login", "then run vault login"},
		{"Then vault login.", "Then vault-client login.", "Then vault login."},
		{`Use "vault hcp connect".`, `Use "vault-client hcp connect".`, `Use "vault-client hcp connect".`},
		{"Project name: vault hcp-prod", "Project name: vault-client hcp-prod", "Project name: vault hcp-prod"},
	} {
		require.Equal(t, tc.help, renameVaultInvocations(vaultInvocation, tc.text, "vault-client"), "help rule for %q", tc.text)
		require.Equal(t, tc.output, renameVaultInvocations(quotedVaultInvocation, tc.text, "vault-client"), "output rule for %q", tc.text)
		require.Equal(t, tc.text, renameVaultInvocations(vaultInvocation, tc.text, "vault"), "vault keeps %q", tc.text)
		require.Equal(t, tc.text, renameVaultInvocations(quotedVaultInvocation, tc.text, ""), "empty name keeps %q", tc.text)
	}
}

// TestClientNaming_HCPConnectUi checks that hcp connect's Ui renames quoted
// vault invocations in anything it prints, whatever upstream's wording.
// Unquoted text such as HCP org and project names, other mentions of vault,
// and prompts pass through, and the vault binary's Ui changes nothing.
func TestClientNaming_HCPConnectUi(t *testing.T) {
	t.Parallel()
	messages := []string{
		"Run 'vault login' to authenticate.",
		"Use `vault hcp connect` to reconnect.",
		"vault login -method=oidc opens a browser.",
		"Project name: vault hcp-prod",
		"Organization name: vault login team",
		"Vault at https://vault.example.com/v1/auth/vault/login reads VAULT_TOKEN, vault-plugin-auth-jwt, and /etc/vault/login.hcl.",
	}
	for _, name := range []string{"", "vault", "vault-client", "renamed-vault"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			want := slices.Clone(messages)
			if name != "" && name != "vault" {
				want[0] = "Run '" + name + " login' to authenticate."
				want[1] = "Use `" + name + " hcp connect` to reconnect."
			}
			for _, tc := range []struct {
				method   string
				print    func(cli.Ui, string)
				toStderr bool
			}{
				{"Output", cli.Ui.Output, false},
				{"Info", cli.Ui.Info, false},
				{"Warn", cli.Ui.Warn, true},
				{"Error", cli.Ui.Error, true},
			} {
				inner := cli.NewMockUi()
				ui := &hcpConnectUi{Ui: inner, BinaryName: name}
				for _, message := range messages {
					tc.print(ui, message)
				}
				written, unused := inner.OutputWriter, inner.ErrorWriter
				if tc.toStderr {
					written, unused = unused, written
				}
				require.Equal(t, strings.Join(want, "\n")+"\n", written.String(), "%s renames only vault login invocations", tc.method)
				require.Empty(t, unused.String(), "%s writes to one stream", tc.method)
			}

			inner := cli.NewMockUi()
			ui := &hcpConnectUi{Ui: inner, BinaryName: name}
			inner.InputReader = strings.NewReader("1\n")
			answer, err := ui.Ask("Run 'vault login' after choosing: ")
			require.NoError(t, err, "prompt")
			require.Equal(t, "1", answer, "prompt answer")
			inner.InputReader = strings.NewReader("disposable-secret\n")
			secret, err := ui.AskSecret("vault login secret: ")
			require.NoError(t, err, "secret prompt")
			require.Equal(t, "disposable-secret", secret, "secret prompt answer")
			require.Equal(t, "Run 'vault login' after choosing: vault login secret: ", inner.OutputWriter.String(), "prompts pass through")
		})
	}
}
