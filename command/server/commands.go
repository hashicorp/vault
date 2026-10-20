// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"github.com/hashicorp/cli"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/audit"
	logicalDb "github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/plugin"
	base "github.com/hashicorp/vault/command/base"
	client "github.com/hashicorp/vault/command/client"
	_ "github.com/hashicorp/vault/helper/builtinplugins"
	physRaft "github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
	sr "github.com/hashicorp/vault/serviceregistration"
	csr "github.com/hashicorp/vault/serviceregistration/consul"
	ksr "github.com/hashicorp/vault/serviceregistration/kubernetes"
)

// Run is the entry point for the full Vault server binary. It delegates to
// the shared CLI plumbing in command/base with the full command set.
func Run(args []string) int {
	return RunCustom(args, nil)
}

// RunCustom is like Run but allows passing a RunOptions value.
func RunCustom(args []string, runOpts *base.RunOptions) int {
	return base.RunCustomWithCommandsFn(args, runOpts, initCommands)
}

// vaultHandlers contains the handlers for creating the various Vault backends.
type vaultHandlers struct {
	physicalBackends     map[string]physical.Factory
	auditBackends        map[string]audit.Factory
	credentialBackends   map[string]logical.Factory
	logicalBackends      map[string]logical.Factory
	serviceRegistrations map[string]sr.Factory
}

// newMinimalVaultHandlers returns a new vaultHandlers that a minimal Vault would use.
func newMinimalVaultHandlers() *vaultHandlers {
	return &vaultHandlers{
		physicalBackends: map[string]physical.Factory{
			"inmem_ha":               physInmem.NewInmemHA,
			"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
			"inmem_transactional":    physInmem.NewTransactionalInmem,
			"inmem":                  physInmem.NewInmem,
			"raft":                   physRaft.NewRaftBackend,
		},
		auditBackends: map[string]audit.Factory{
			"file":   audit.NewFileBackend,
			"socket": audit.NewSocketBackend,
			"syslog": audit.NewSyslogBackend,
		},
		credentialBackends: map[string]logical.Factory{
			"plugin": plugin.Factory,
		},
		logicalBackends: map[string]logical.Factory{
			"plugin":   plugin.Factory,
			"database": logicalDb.Factory,
			// This is also available in the plugin catalog, but is here due to the need to
			// automatically mount it.
			"kv": logicalKv.Factory,
		},
		serviceRegistrations: map[string]sr.Factory{
			"consul":     csr.NewServiceRegistration,
			"kubernetes": ksr.NewServiceRegistration,
		},
	}
}

// newVaultHandlers returns a new vaultHandlers composed of newMinimalVaultHandlers()
// and any addon handlers from Vault CE and Vault Enterprise selected by Go build tags.
func newVaultHandlers() *vaultHandlers {
	handlers := newMinimalVaultHandlers()
	extendAddonHandlers(handlers)
	entExtendAddonHandlers(handlers)

	return handlers
}

func initCommands(ui, serverCmdUi cli.Ui, runOpts *base.RunOptions) map[string]cli.CommandFactory {
	handlers := newVaultHandlers()

	getBaseCommand := func() *base.BaseCommand {
		return base.NewBaseCommand(ui, runOpts)
	}

	commands := client.InitClientCommands(ui, serverCmdUi, runOpts)

	// Server-only per PRD VLT-2015; could be a client command.
	commands["audit"] = func() (cli.Command, error) {
		return &AuditCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["audit disable"] = func() (cli.Command, error) {
		return &AuditDisableCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["audit enable"] = func() (cli.Command, error) {
		return &AuditEnableCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["audit list"] = func() (cli.Command, error) {
		return &AuditListCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["debug"] = func() (cli.Command, error) {
		return &DebugCommand{
			BaseCommand: getBaseCommand(),
			ShutdownCh:  client.MakeShutdownCh(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["events subscribe"] = func() (cli.Command, error) {
		return &EventsSubscribeCommands{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator generate-root"] = func() (cli.Command, error) {
		return &OperatorGenerateRootCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator init"] = func() (cli.Command, error) {
		return &OperatorInitCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft"] = func() (cli.Command, error) {
		return &OperatorRaftCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft autopilot get-config"] = func() (cli.Command, error) {
		return &OperatorRaftAutopilotGetConfigCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft autopilot set-config"] = func() (cli.Command, error) {
		return &OperatorRaftAutopilotSetConfigCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft autopilot state"] = func() (cli.Command, error) {
		return &OperatorRaftAutopilotStateCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft list-peers"] = func() (cli.Command, error) {
		return &OperatorRaftListPeersCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft join"] = func() (cli.Command, error) {
		return &OperatorRaftJoinCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft remove-peer"] = func() (cli.Command, error) {
		return &OperatorRaftRemovePeerCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft snapshot"] = func() (cli.Command, error) {
		return &OperatorRaftSnapshotCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft snapshot restore"] = func() (cli.Command, error) {
		return &OperatorRaftSnapshotRestoreCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator raft snapshot save"] = func() (cli.Command, error) {
		return &OperatorRaftSnapshotSaveCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator rekey"] = func() (cli.Command, error) {
		return &OperatorRekeyCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator seal"] = func() (cli.Command, error) {
		return &OperatorSealCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator step-down"] = func() (cli.Command, error) {
		return &OperatorStepDownCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator usage"] = func() (cli.Command, error) {
		return &OperatorUsageCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator utilization"] = func() (cli.Command, error) {
		return &OperatorUtilizationCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["operator members"] = func() (cli.Command, error) {
		return &OperatorMembersCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki"] = func() (cli.Command, error) {
		return &PKICommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki health-check"] = func() (cli.Command, error) {
		return &PKIHealthCheckCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki issue"] = func() (cli.Command, error) {
		return &PKIIssueCACommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki list-intermediates"] = func() (cli.Command, error) {
		return &PKIListIntermediateCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki reissue"] = func() (cli.Command, error) {
		return &PKIReIssueCACommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["pki verify-sign"] = func() (cli.Command, error) {
		return &PKIVerifySignCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin"] = func() (cli.Command, error) {
		return &PluginCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin deregister"] = func() (cli.Command, error) {
		return &PluginDeregisterCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin info"] = func() (cli.Command, error) {
		return &PluginInfoCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin list"] = func() (cli.Command, error) {
		return &PluginListCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin register"] = func() (cli.Command, error) {
		return NewPluginRegisterCommand(getBaseCommand()), nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin reload"] = func() (cli.Command, error) {
		return &PluginReloadCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin reload-status"] = func() (cli.Command, error) {
		return &PluginReloadStatusCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin runtime"] = func() (cli.Command, error) {
		return &PluginRuntimeCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin runtime register"] = func() (cli.Command, error) {
		return &PluginRuntimeRegisterCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin runtime deregister"] = func() (cli.Command, error) {
		return &PluginRuntimeDeregisterCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin runtime info"] = func() (cli.Command, error) {
		return &PluginRuntimeInfoCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["plugin runtime list"] = func() (cli.Command, error) {
		return &PluginRuntimeListCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["policy"] = func() (cli.Command, error) {
		return &PolicyCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["policy delete"] = func() (cli.Command, error) {
		return &PolicyDeleteCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["policy list"] = func() (cli.Command, error) {
		return &PolicyListCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["policy read"] = func() (cli.Command, error) {
		return &PolicyReadCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["policy write"] = func() (cli.Command, error) {
		return &PolicyWriteCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["ssh"] = func() (cli.Command, error) {
		return &SSHCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transform"] = func() (cli.Command, error) {
		return &TransformCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transform import"] = func() (cli.Command, error) {
		return &TransformImportCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transform import-version"] = func() (cli.Command, error) {
		return &TransformImportVersionCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit"] = func() (cli.Command, error) {
		return &TransitCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit envelope"] = func() (cli.Command, error) {
		return &TransitEnvelopeCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit envelope encrypt"] = func() (cli.Command, error) {
		return &TransitEncryptCommand{
			TransitEnvEncCommand: TransitEnvEncCommand{
				BaseCommand: getBaseCommand(),
			},
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit envelope decrypt"] = func() (cli.Command, error) {
		return &TransitDecryptCommand{
			TransitEnvEncCommand: TransitEnvEncCommand{
				BaseCommand: getBaseCommand(),
			},
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit envelope header"] = func() (cli.Command, error) {
		return &TransitEnvHeaderCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit import"] = func() (cli.Command, error) {
		return &TransitImportCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["transit import-version"] = func() (cli.Command, error) {
		return &TransitImportVersionCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["version-history"] = func() (cli.Command, error) {
		return &VersionHistoryCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	// Server-only per PRD VLT-2015; could be a client command.
	commands["monitor"] = func() (cli.Command, error) {
		return &MonitorCommand{
			BaseCommand: getBaseCommand(),
			ShutdownCh:  client.MakeShutdownCh(),
		}, nil
	}

	commands["operator diagnose"] = func() (cli.Command, error) {
		return &OperatorDiagnoseCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	commands["operator migrate"] = func() (cli.Command, error) {
		return &OperatorMigrateCommand{
			BaseCommand:      getBaseCommand(),
			PhysicalBackends: handlers.physicalBackends,
			ShutdownCh:       client.MakeShutdownCh(),
		}, nil
	}
	commands["operator raft snapshot inspect"] = func() (cli.Command, error) {
		return &OperatorRaftSnapshotInspectCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	commands["policy fmt"] = func() (cli.Command, error) {
		return &PolicyFmtCommand{
			BaseCommand: getBaseCommand(),
		}, nil
	}
	commands["server"] = func() (cli.Command, error) {
		return &ServerCommand{
			BaseCommand: base.NewBaseCommand(serverCmdUi, &base.RunOptions{
				TokenHelper: runOpts.TokenHelper,
				Address:     runOpts.Address,
			}),
			AuditBackends:        handlers.auditBackends,
			CredentialBackends:   handlers.credentialBackends,
			LogicalBackends:      handlers.logicalBackends,
			PhysicalBackends:     handlers.physicalBackends,
			ServiceRegistrations: handlers.serviceRegistrations,

			ShutdownCh: client.MakeShutdownCh(),
			SighupCh:   client.MakeSighupCh(),
			SigUSR2Ch:  client.MakeSigUSR2Ch(),
		}, nil
	}

	entInitServerOnlyCommands(ui, runOpts, commands)

	return commands
}
