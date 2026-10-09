// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

package command

import (
	"github.com/hashicorp/cli"
	credOIDC "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/audit"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
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
	loginHandlers        map[string]base.LoginHandler
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
		loginHandlers: map[string]base.LoginHandler{
			"cert":  &credCert.CLIHandler{},
			"oidc":  &credOIDC.CLIHandler{},
			"token": &credToken.CLIHandler{},
			"userpass": &credUserpass.CLIHandler{
				DefaultMount: "userpass",
			},
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

	commands["auth help"] = func() (cli.Command, error) {
		return &client.AuthHelpCommand{
			BaseCommand: getBaseCommand(),
			Handlers:    handlers.loginHandlers,
		}, nil
	}
	commands["login"] = func() (cli.Command, error) {
		return &client.LoginCommand{
			BaseCommand: getBaseCommand(),
			Handlers:    handlers.loginHandlers,
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

	return commands
}
