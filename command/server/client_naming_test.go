// Copyright IBM Corp. 2026
// SPDX-License-Identifier: BUSL-1.1

package command_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
	clientcmd "github.com/hashicorp/vault/command/client"
	"github.com/hashicorp/vault/command/client/login/userpass"
	"github.com/hashicorp/vault/helper/testhelpers/minimal"
	"github.com/stretchr/testify/require"
)

// TestClientNaming_AgentSuggestion checks that Agent tells users to run proxy
// with the same binary they started Agent with.
func TestClientNaming_AgentSuggestion(t *testing.T) {
	t.Parallel()
	// Agent changes global telemetry on startup, so we need to stop it and wait
	// for it to exit before starting the next case.
	for _, name := range []string{"vault", "vault-client"} {
		t.Run(name, func(t *testing.T) {
			cluster := minimal.NewTestSoloCluster(t, nil)
			ui, cmd, control := clientcmd.MakeTestAgentCommandWithControl(hclog.NewNullLogger(), cluster.Cores[0].Client)
			cmd.BinaryName = name
			path := filepath.Join(t.TempDir(), "agent.hcl")
			require.NoError(t, os.WriteFile(path, []byte(`
listener "tcp" {
  address = "127.0.0.1:0"
  tls_disable = true
}
`), 0o600), "isolated Agent config")
			done := make(chan int, 1)
			t.Cleanup(func() {
				close(cmd.ShutdownCh)
				select {
				case code := <-done:
					require.Zero(t, code, "Agent shutdown status: %s", ui.ErrorWriter.String())
				case <-time.After(5 * time.Second):
					t.Error("Agent did not stop")
					return
				}
				require.Contains(t, ui.ErrorWriter.String(), "new subcommand, "+name+" proxy.", "migration invocation")
			})
			go func() { done <- cmd.Run([]string{"-config=" + path}) }()
			select {
			case <-control.StartedCh:
			case code := <-done:
				done <- code
				t.Fatalf("Agent exited before startup: %s", ui.ErrorWriter.String())
			case <-time.After(5 * time.Second):
				t.Fatal("Agent did not start")
			}
		})
	}
}

// TestClientNaming_MountedAuthHelp checks that we find the auth method for a mount
// and print its handler's help with the right binary name.
func TestClientNaming_MountedAuthHelp(t *testing.T) {
	t.Parallel()
	cluster := minimal.NewTestSoloCluster(t, nil)
	client := cluster.Cores[0].Client
	require.NoError(t, client.Sys().EnableAuthWithOptionsWithContext(t.Context(), "my-method", &api.EnableAuthOptions{
		Type: "userpass",
	}), "mount userpass at a path that doesn't match its type")
	ui := cli.NewMockUi()
	cmd, err := clientcmd.InitClientCommands(ui, ui, &base.RunOptions{
		BinaryName: "vault-client", Client: client, HCPTokenHelper: clientAuthNoHCPTokenHelper{},
	})["auth help"]()
	require.NoError(t, err, "auth help factory")
	require.Zero(t, cmd.Run([]string{"my-method/"}), "mounted alias help status")
	require.Equal(t, (&userpass.CLIHandler{BinaryName: "vault-client"}).Help()+"\n", ui.OutputWriter.String(), "resolved method help")
	require.Contains(t, ui.OutputWriter.String(), "Usage: vault-client login -method=userpass", "resolved method usage")
	require.Empty(t, ui.ErrorWriter.String(), "no mount lookup errors")
}
