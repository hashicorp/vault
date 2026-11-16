// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	vaulthttp "github.com/hashicorp/vault/http"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/stretchr/testify/require"
)

// TestServerNaming_CommandHelp checks that every command in the vault binary,
// including the server-only ones, shows the name the binary was run as in its
// usage line and examples instead of a hard-coded "vault".
func TestServerNaming_CommandHelp(t *testing.T) {
	t.Parallel()
	// "server" has its own pattern because flag help mentions "the vault
	// server" in prose. Only quoted or example invocations count for it.
	bareInvocation := regexp.MustCompile(`(^|[\s"'` + "`" + `|])vault\s+(?:agent|audit|auth|debug|delete|events|hcp|kv|lease|license|list|login|monitor|namespace|operator|patch|path-help|pki|plugin|policy|print|proxy|read|secrets|ssh|status|token|tpm|transform|transit|unwrap|version|write)\b`)
	bareServerInvocation := regexp.MustCompile(`(["'` + "`" + `]|\$ |Usage: )vault server\b`)
	for _, name := range []string{"vault", "renamed-vault"} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ui := cli.NewMockUi()
			for command, factory := range initCommands(ui, ui, &base.RunOptions{BinaryName: name}) {
				cmd, err := factory()
				require.NoError(t, err, command)
				help := cmd.Help()
				// go-licensing's usage line names the subcommand but no binary.
				if command != "license inspect" {
					require.Contains(t, help, "Usage: "+name+" "+command, "command %s", command)
				}
				if name != "vault" {
					require.False(t, bareInvocation.MatchString(help), "bare executable in %s help: %s", command, bareInvocation.FindString(help))
					require.False(t, bareServerInvocation.MatchString(help), "bare executable in %s help: %s", command, bareServerInvocation.FindString(help))
				}
			}
		})
	}
}

// TestServerNaming_StartSuggestion checks that when a key is provided but no
// rekey or root generation is in progress, the error tells users to start one
// with the binary they ran rather than "vault".
func TestServerNaming_StartSuggestion(t *testing.T) {
	t.Parallel()
	for _, command := range []string{"operator rekey", "operator generate-root"} {
		t.Run(command, func(t *testing.T) {
			t.Parallel()
			cluster := minimal.NewTestSoloCluster(t, nil)
			ui := cli.NewMockUi()
			cmd, err := initCommands(ui, ui, &base.RunOptions{
				BinaryName: "renamed-vault", Client: cluster.Cores[0].Client, HCPTokenHelper: noHCPTokenHelper{},
			})[command]()
			require.NoError(t, err, "command factory")
			require.Equal(t, 1, cmd.Run([]string{"disposable-key"}), "nothing in progress")
			require.Contains(t, strings.Join(strings.Fields(ui.ErrorWriter.String()), " "), `"renamed-vault `+command+` -init"`, "start suggestion")
		})
	}
}

// TestServerNaming_RuntimeHints checks that follow-up commands printed while
// server-only commands run use the name the binary was run as. Each case
// produces one hint from a test cluster, in-memory storage, or the code that
// builds the message.
func TestServerNaming_RuntimeHints(t *testing.T) {
	t.Parallel()
	bareInvocation := regexp.MustCompile(`(^|[\s"'` + "`" + `|])vault\s+(?:audit|delete|events|login|operator|patch|pki|plugin|policy|read|server|token|write)\b`)
	runOpts := &base.RunOptions{BinaryName: "renamed-vault"}
	for _, tc := range []struct {
		name string
		want []string
		run  func(t *testing.T) string
	}{
		{
			name: "operator init rekey hint",
			want: []string{`See "renamed-vault operator rekey" for more information.`},
			run: func(t *testing.T) string {
				cluster := vault.NewTestCluster(t, nil, &vault.TestClusterOptions{
					NumCores: 1, HandlerFunc: vaulthttp.Handler, SkipInit: true,
				})
				ui := cli.NewMockUi()
				cmd := &OperatorInitCommand{BaseCommand: base.NewBaseCommand(&base.VaultUI{Ui: ui, Format: "table"}, runOpts)}
				require.Zero(t, cmd.init(cluster.Cores[0].Client, &api.InitRequest{SecretShares: 1, SecretThreshold: 1}), "init status")
				return ui.OutputWriter.String()
			},
		},
		{
			name: "operator rekey barrier backup",
			want: []string{`using "renamed-vault operator rekey -backup-delete".`},
			run: func(t *testing.T) string {
				ui := cli.NewMockUi()
				cmd := &OperatorRekeyCommand{BaseCommand: base.NewBaseCommand(&base.VaultUI{Ui: ui, Format: "table"}, runOpts), flagTarget: "barrier"}
				require.Zero(t, cmd.printUnsealKeys(nil, &api.RekeyStatusResponse{N: 1, T: 1}, &api.RekeyUpdateResponse{
					Keys: []string{"disposable-key"}, PGPFingerprints: []string{"disposable-fingerprint"}, Backup: true,
				}), "print status")
				return ui.OutputWriter.String()
			},
		},
		{
			name: "operator rekey recovery backup",
			want: []string{`using "renamed-vault operator rekey -backup-delete -target=recovery".`},
			run: func(t *testing.T) string {
				ui := cli.NewMockUi()
				cmd := &OperatorRekeyCommand{BaseCommand: base.NewBaseCommand(&base.VaultUI{Ui: ui, Format: "table"}, runOpts), flagTarget: "recovery"}
				require.Zero(t, cmd.printUnsealKeys(nil, &api.RekeyStatusResponse{N: 1, T: 1}, &api.RekeyUpdateResponse{
					Keys: []string{"disposable-key"}, PGPFingerprints: []string{"disposable-fingerprint"}, Backup: true,
				}), "print status")
				return ui.OutputWriter.String()
			},
		},
		{
			name: "server storage migration lock",
			want: []string{"Use 'renamed-vault operator migrate -reset'"},
			run: func(t *testing.T) string {
				backend, err := physInmem.NewInmem(nil, hclog.NewNullLogger())
				require.NoError(t, err, "in-memory storage")
				require.NoError(t, SetStorageMigration(backend, true), "migration lock")
				ui := cli.NewMockUi()
				cmd := &ServerCommand{BaseCommand: base.NewBaseCommand(ui, runOpts)}
				require.True(t, cmd.storageMigrationActive(backend), "migration blocks startup")
				return ui.ErrorWriter.String()
			},
		},
		{
			name: "pki issue continue and delete key",
			want: []string{
				"`renamed-vault write pki/issuer/default/sign-intermediate csr=@pki_int.csr ...`",
				"`renamed-vault write pki_int/issuers/import/cert pem_bundle=@chain.pem`",
				"`renamed-vault patch pki_int/issuer/<issuer-uuid> issuer_name=new-int`",
				"`renamed-vault delete pki_int/key/disposable-key`",
			},
			run: func(t *testing.T) string {
				return inCaseOfFailure{
					csrGenerated: true, createdKeyId: "disposable-key", intermediateMount: "pki_int",
					parentMount: "pki", parentIssuer: "pki/issuer/default", newName: "new-int", binaryName: "renamed-vault",
				}.generateFailureMessage()
			},
		},
		{
			name: "pki issue revoke signed certificate",
			want: []string{"`renamed-vault write pki/revoke/01:02`"},
			run: func(t *testing.T) string {
				return inCaseOfFailure{
					csrGenerated: true, csrSigned: true, certSerialNumber: "01:02", intermediateMount: "pki_int",
					parentMount: "pki", parentIssuer: "pki/issuer/default", binaryName: "renamed-vault",
				}.generateFailureMessage()
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			output := strings.Join(strings.Fields(tc.run(t)), " ")
			for _, want := range tc.want {
				require.Contains(t, output, want, "hint uses the binary name")
			}
			require.False(t, bareInvocation.MatchString(output), "bare executable in hint: %s", bareInvocation.FindString(output))
		})
	}
}
