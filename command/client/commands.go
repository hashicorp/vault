// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/cli"
	hcpvlib "github.com/hashicorp/vault-hcp-lib"
	base "github.com/hashicorp/vault/command/base"
	"github.com/hashicorp/vault/version"
)

// Run is the entry point for the client binary.
func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom allows passing in base run options (e.g. a custom token helper).
func RunCustom(args []string, runOpts *base.RunOptions) int {
	return base.RunCustomWithCommandsFn(args, runOpts, InitClientCommands)
}

// InitClientCommands returns the set of CLI commands for the client binary.
// It registers all commands that interact with a running Vault cluster via the
// API, plus vault agent and vault proxy. It does not include vault server or
// any commands that require server-side internals (operator diagnose, operator
// migrate, operator raft snapshot inspect, policy fmt).
func InitClientCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions) map[string]cli.CommandFactory {
	getBaseCommand := func() *base.BaseCommand {
		return base.NewBaseCommand(ui, runOpts)
	}

	commands := map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &AgentCommand{
				BaseCommand: &base.BaseCommand{
					UI: serverCmdUi,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"agent generate-config": func() (cli.Command, error) {
			return &AgentGenerateConfigCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"audit": func() (cli.Command, error) {
			return &AuditCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"audit disable": func() (cli.Command, error) {
			return &AuditDisableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"audit enable": func() (cli.Command, error) {
			return &AuditEnableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"audit list": func() (cli.Command, error) {
			return &AuditListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth tune": func() (cli.Command, error) {
			return &AuthTuneCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth": func() (cli.Command, error) {
			return &AuthCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth disable": func() (cli.Command, error) {
			return &AuthDisableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth enable": func() (cli.Command, error) {
			return &AuthEnableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth help": func() (cli.Command, error) {
			return &AuthHelpCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth list": func() (cli.Command, error) {
			return &AuthListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"auth move": func() (cli.Command, error) {
			return &AuthMoveCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"debug": func() (cli.Command, error) {
			return &DebugCommand{
				BaseCommand: getBaseCommand(),
				ShutdownCh:  MakeShutdownCh(),
			}, nil
		},
		"delete": func() (cli.Command, error) {
			return &DeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"events subscribe": func() (cli.Command, error) {
			return &EventsSubscribeCommands{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease": func() (cli.Command, error) {
			return &LeaseCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease renew": func() (cli.Command, error) {
			return &LeaseRenewCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease lookup": func() (cli.Command, error) {
			return &LeaseLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"lease revoke": func() (cli.Command, error) {
			return &LeaseRevokeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace": func() (cli.Command, error) {
			return &NamespaceCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace list": func() (cli.Command, error) {
			return &NamespaceListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace lookup": func() (cli.Command, error) {
			return &NamespaceLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace create": func() (cli.Command, error) {
			return &NamespaceCreateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace patch": func() (cli.Command, error) {
			return &NamespacePatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace delete": func() (cli.Command, error) {
			return &NamespaceDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace lock": func() (cli.Command, error) {
			return &NamespaceAPILockCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"namespace unlock": func() (cli.Command, error) {
			return &NamespaceAPIUnlockCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator": func() (cli.Command, error) {
			return &OperatorCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator generate-root": func() (cli.Command, error) {
			return &OperatorGenerateRootCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator init": func() (cli.Command, error) {
			return &OperatorInitCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator key-status": func() (cli.Command, error) {
			return &OperatorKeyStatusCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft": func() (cli.Command, error) {
			return &OperatorRaftCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft autopilot get-config": func() (cli.Command, error) {
			return &OperatorRaftAutopilotGetConfigCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft autopilot set-config": func() (cli.Command, error) {
			return &OperatorRaftAutopilotSetConfigCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft autopilot state": func() (cli.Command, error) {
			return &OperatorRaftAutopilotStateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft list-peers": func() (cli.Command, error) {
			return &OperatorRaftListPeersCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft join": func() (cli.Command, error) {
			return &OperatorRaftJoinCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft remove-peer": func() (cli.Command, error) {
			return &OperatorRaftRemovePeerCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft snapshot": func() (cli.Command, error) {
			return &OperatorRaftSnapshotCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft snapshot restore": func() (cli.Command, error) {
			return &OperatorRaftSnapshotRestoreCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator raft snapshot save": func() (cli.Command, error) {
			return &OperatorRaftSnapshotSaveCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator rekey": func() (cli.Command, error) {
			return &OperatorRekeyCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator rotate": func() (cli.Command, error) {
			return &OperatorRotateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator seal": func() (cli.Command, error) {
			return &OperatorSealCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator step-down": func() (cli.Command, error) {
			return &OperatorStepDownCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator usage": func() (cli.Command, error) {
			return &OperatorUsageCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator utilization": func() (cli.Command, error) {
			return &OperatorUtilizationCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator unseal": func() (cli.Command, error) {
			return &OperatorUnsealCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator members": func() (cli.Command, error) {
			return &OperatorMembersCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"patch": func() (cli.Command, error) {
			return &PatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"path-help": func() (cli.Command, error) {
			return &PathHelpCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki": func() (cli.Command, error) {
			return &PKICommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki health-check": func() (cli.Command, error) {
			return &PKIHealthCheckCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki issue": func() (cli.Command, error) {
			return &PKIIssueCACommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki list-intermediates": func() (cli.Command, error) {
			return &PKIListIntermediateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki reissue": func() (cli.Command, error) {
			return &PKIReIssueCACommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"pki verify-sign": func() (cli.Command, error) {
			return &PKIVerifySignCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin": func() (cli.Command, error) {
			return &PluginCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin deregister": func() (cli.Command, error) {
			return &PluginDeregisterCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin info": func() (cli.Command, error) {
			return &PluginInfoCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin list": func() (cli.Command, error) {
			return &PluginListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin register": func() (cli.Command, error) {
			return NewPluginRegisterCommand(getBaseCommand()), nil
		},
		"plugin reload": func() (cli.Command, error) {
			return &PluginReloadCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin reload-status": func() (cli.Command, error) {
			return &PluginReloadStatusCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin runtime": func() (cli.Command, error) {
			return &PluginRuntimeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin runtime register": func() (cli.Command, error) {
			return &PluginRuntimeRegisterCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin runtime deregister": func() (cli.Command, error) {
			return &PluginRuntimeDeregisterCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin runtime info": func() (cli.Command, error) {
			return &PluginRuntimeInfoCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"plugin runtime list": func() (cli.Command, error) {
			return &PluginRuntimeListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"proxy": func() (cli.Command, error) {
			return &ProxyCommand{
				BaseCommand: &base.BaseCommand{
					UI: serverCmdUi,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
				SigUSR2Ch:  MakeSigUSR2Ch(),
			}, nil
		},
		"policy": func() (cli.Command, error) {
			return &PolicyCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"policy delete": func() (cli.Command, error) {
			return &PolicyDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"policy list": func() (cli.Command, error) {
			return &PolicyListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"policy read": func() (cli.Command, error) {
			return &PolicyReadCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"policy write": func() (cli.Command, error) {
			return &PolicyWriteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"print": func() (cli.Command, error) {
			return &PrintCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"print token": func() (cli.Command, error) {
			return &PrintTokenCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &ReadCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets": func() (cli.Command, error) {
			return &SecretsCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets disable": func() (cli.Command, error) {
			return &SecretsDisableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets enable": func() (cli.Command, error) {
			return &SecretsEnableCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets list": func() (cli.Command, error) {
			return &SecretsListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets move": func() (cli.Command, error) {
			return &SecretsMoveCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"secrets tune": func() (cli.Command, error) {
			return &SecretsTuneCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"ssh": func() (cli.Command, error) {
			return &SSHCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transform": func() (cli.Command, error) {
			return &TransformCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transform import": func() (cli.Command, error) {
			return &TransformImportCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transform import-version": func() (cli.Command, error) {
			return &TransformImportVersionCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transit": func() (cli.Command, error) {
			return &TransitCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transit import": func() (cli.Command, error) {
			return &TransitImportCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"transit import-version": func() (cli.Command, error) {
			return &TransitImportVersionCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &TokenCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token create": func() (cli.Command, error) {
			return &TokenCreateCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token capabilities": func() (cli.Command, error) {
			return &TokenCapabilitiesCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token lookup": func() (cli.Command, error) {
			return &TokenLookupCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token renew": func() (cli.Command, error) {
			return &TokenRenewCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"token revoke": func() (cli.Command, error) {
			return &TokenRevokeCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"unwrap": func() (cli.Command, error) {
			return &UnwrapCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				VersionInfo: version.GetVersion(),
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"version-history": func() (cli.Command, error) {
			return &VersionHistoryCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"write": func() (cli.Command, error) {
			return &WriteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv": func() (cli.Command, error) {
			return &KVCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv put": func() (cli.Command, error) {
			return &KVPutCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv patch": func() (cli.Command, error) {
			return &KVPatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv rollback": func() (cli.Command, error) {
			return &KVRollbackCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv get": func() (cli.Command, error) {
			return &KVGetCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv delete": func() (cli.Command, error) {
			return &KVDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv list": func() (cli.Command, error) {
			return &KVListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv destroy": func() (cli.Command, error) {
			return &KVDestroyCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv undelete": func() (cli.Command, error) {
			return &KVUndeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv enable-versioning": func() (cli.Command, error) {
			return &KVEnableVersioningCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata": func() (cli.Command, error) {
			return &KVMetadataCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata put": func() (cli.Command, error) {
			return &KVMetadataPutCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata patch": func() (cli.Command, error) {
			return &KVMetadataPatchCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata get": func() (cli.Command, error) {
			return &KVMetadataGetCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"kv metadata delete": func() (cli.Command, error) {
			return &KVMetadataDeleteCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"monitor": func() (cli.Command, error) {
			return &MonitorCommand{
				BaseCommand: getBaseCommand(),
				ShutdownCh:  MakeShutdownCh(),
			}, nil
		},
	}

	entInitCommands(ui, serverCmdUi, runOpts, commands)
	InitHCPCommands(ui, commands)

	return commands
}

// InitHCPCommands adds HCP-specific commands to the given commands map.
// It checks for conflicts and only adds a command if it doesn't already exist.
func InitHCPCommands(ui cli.Ui, commands map[string]cli.CommandFactory) {
	for cmd, cmdFactory := range hcpvlib.InitHCPCommand(ui) {
		// check for conflicts and only put command in the map in case it doesn't conflict with existing one
		_, ok := commands[cmd]
		if !ok {
			commands[cmd] = cmdFactory
		} else {
			ui.Error("Failed to initialize HCP commands.")
			break
		}
	}
}

// MakeShutdownCh returns a channel that can be used for shutdown
// notifications for commands. This channel will send a message for every
// SIGINT or SIGTERM received.
func MakeShutdownCh() chan struct{} {
	resultCh := make(chan struct{})

	shutdownCh := make(chan os.Signal, 4)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownCh
		close(resultCh)
	}()
	return resultCh
}

// MakeSighupCh returns a channel that can be used for SIGHUP
// reloading. This channel will send a message for every
// SIGHUP received.
func MakeSighupCh() chan struct{} {
	resultCh := make(chan struct{})

	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, syscall.SIGHUP)
	go func() {
		for {
			<-signalCh
			resultCh <- struct{}{}
		}
	}()
	return resultCh
}
