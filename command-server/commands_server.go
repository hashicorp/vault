package command_server

import (
	"github.com/hashicorp/vault/command"
	"github.com/mitchellh/cli"
)

func initServerCommands(serverCmdUi cli.Ui, runOpts *command.RunOptions) map[string]cli.CommandFactory {
	getBaseCommand := func() *command.BaseCommand {
		return &command.BaseCommand{
			TkHelper:    runOpts.TokenHelper,
			FlagAddress: runOpts.Address,
			ApiClient:   runOpts.Client,
		}
	}

	commands := map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &AgentCommand{
				BaseCommand: &command.BaseCommand{
					UI: serverCmdUi,
				},
				ShutdownCh: command.MakeShutdownCh(),
				SighupCh:   command.MakeSighupCh(),
			}, nil
		},
		"agent generate-config": func() (cli.Command, error) {
			return &AgentGenerateConfigCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator": func() (cli.Command, error) {
			return &OperatorCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"operator diagnose": func() (cli.Command, error) {
			return &OperatorDiagnoseCommand{
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
		"operator migrate": func() (cli.Command, error) {
			return &OperatorMigrateCommand{
				BaseCommand:      getBaseCommand(),
				PhysicalBackends: command.PhysicalBackends,
				ShutdownCh:       command.MakeShutdownCh(),
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
		"operator raft snapshot inspect": func() (cli.Command, error) {
			return &OperatorRaftSnapshotInspectCommand{
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
		"proxy": func() (cli.Command, error) {
			return &ProxyCommand{
				BaseCommand: &command.BaseCommand{
					UI: serverCmdUi,
				},
				ShutdownCh: command.MakeShutdownCh(),
				SighupCh:   command.MakeSighupCh(),
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: &command.BaseCommand{
					UI:          serverCmdUi,
					TkHelper:    runOpts.TokenHelper,
					FlagAddress: runOpts.Address,
				},
				AuditBackends:      command.AuditBackends,
				CredentialBackends: command.CredentialBackends,
				LogicalBackends:    command.LogicalBackends,
				PhysicalBackends:   command.PhysicalBackends,

				ServiceRegistrations: command.ServiceRegistrations,

				ShutdownCh: command.MakeShutdownCh(),
				SighupCh:   command.MakeSighupCh(),
				SigUSR2Ch:  command.MakeSigUSR2Ch(),
			}, nil
		},
	}

	return commands
}
