// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"github.com/hashicorp/cli"
	"github.com/hashicorp/go-hclog"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	"github.com/hashicorp/vault/api"
	base "github.com/hashicorp/vault/command/base"
)

type testHCPTokenHelper struct{}

func (testHCPTokenHelper) GetHCPToken(string) (*hcpvlib.HCPToken, error) {
	return nil, nil
}

// TestCommandControl exposes lifecycle notifications needed by integration
// tests without exposing mutable command internals.
type TestCommandControl struct {
	StartedCh  <-chan struct{}
	ReloadedCh <-chan struct{}
}

// MakeTestAgentCommand constructs an AgentCommand wired to a MockUi for use in
// tests. If client is non-nil it is pre-configured as the command's Vault client.
func MakeTestAgentCommand(logger hclog.Logger, client *api.Client) (*cli.MockUi, *AgentCommand) {
	ui, cmd, _ := MakeTestAgentCommandWithControl(logger, client)
	return ui, cmd
}

// MakeTestAgentCommandWithControl constructs an AgentCommand and returns
// receive-only lifecycle notifications for integration-test synchronization.
func MakeTestAgentCommandWithControl(logger hclog.Logger, client *api.Client) (*cli.MockUi, *AgentCommand, TestCommandControl) {
	ui := cli.NewMockUi()
	startedCh := make(chan struct{}, 5)
	reloadedCh := make(chan struct{}, 5)
	cmd := &AgentCommand{
		BaseCommand: makeTestBaseCommand(ui, client),
		ShutdownCh:  MakeShutdownCh(),
		SighupCh:    MakeSighupCh(),
		SigUSR2Ch:   MakeSigUSR2Ch(),
		logger:      logger,
		startedCh:   startedCh,
		reloadedCh:  reloadedCh,
	}
	return ui, cmd, TestCommandControl{
		StartedCh:  startedCh,
		ReloadedCh: reloadedCh,
	}
}

// MakeTestProxyCommandWithControl constructs a ProxyCommand and returns
// receive-only lifecycle notifications for integration-test synchronization.
func MakeTestProxyCommandWithControl(logger hclog.Logger, client *api.Client) (*cli.MockUi, *ProxyCommand, TestCommandControl) {
	ui := cli.NewMockUi()
	startedCh := make(chan struct{}, 5)
	reloadedCh := make(chan struct{}, 5)
	cmd := &ProxyCommand{
		BaseCommand: makeTestBaseCommand(ui, client),
		ShutdownCh:  MakeShutdownCh(),
		SighupCh:    MakeSighupCh(),
		SigUSR2Ch:   MakeSigUSR2Ch(),
		logger:      logger,
		startedCh:   startedCh,
		reloadedCh:  reloadedCh,
	}
	return ui, cmd, TestCommandControl{
		StartedCh:  startedCh,
		ReloadedCh: reloadedCh,
	}
}

// MakeTestDebugCommand constructs a DebugCommand with shortened capture timing
// for integration tests.
func MakeTestDebugCommand(client *api.Client) (*cli.MockUi, *DebugCommand) {
	ui := cli.NewMockUi()
	return ui, &DebugCommand{
		BaseCommand:      makeTestBaseCommand(ui, client),
		SkipTimingChecks: true,
	}
}

func makeTestBaseCommand(ui cli.Ui, client *api.Client) *base.BaseCommand {
	cmd := &base.BaseCommand{UI: ui}
	cmd.SetHCPTokenHelper(testHCPTokenHelper{})
	if client != nil {
		cmd.SetClient(client)
	}
	return cmd
}
