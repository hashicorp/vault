package cli

import (
	"os"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"
	"github.com/hashicorp/vault/command/server"
	"github.com/hashicorp/vault/version"

	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credMarathon "github.com/hashicorp/vault/builtin/credential/marathon"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"

	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/builtin/logical/cassandra"
	"github.com/hashicorp/vault/builtin/logical/consul"
	"github.com/hashicorp/vault/builtin/logical/mssql"
	"github.com/hashicorp/vault/builtin/logical/mysql"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/postgresql"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/transit"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/logical"
	"github.com/mitchellh/cli"
)

// Commands returns the mapping of CLI commands for Vault. The meta
// parameter lets you set meta options for all commands.
func Commands(metaPtr *command.Meta) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = new(command.Meta)
	}

	meta := *metaPtr
	if meta.Ui == nil {
		meta.Ui = &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}

	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Meta: meta,
			}, nil
		},

		"server": func() (cli.Command, error) {
			return &command.ServerCommand{
				Meta: meta,
				AuditBackends: map[string]audit.Factory{
					"file":   auditFile.Factory,
					"syslog": auditSyslog.Factory,
				},
				CredentialBackends: map[string]logical.Factory{
					"cert":     credCert.Factory,
					"app-id":   credAppId.Factory,
					"github":   credGitHub.Factory,
					"userpass": credUserpass.Factory,
					"ldap":     credLdap.Factory,
					"marathon": credMarathon.Factory,
				},
				LogicalBackends: map[string]logical.Factory{
					"aws":        aws.Factory,
					"consul":     consul.Factory,
					"postgresql": postgresql.Factory,
					"cassandra":  cassandra.Factory,
					"pki":        pki.Factory,
					"transit":    transit.Factory,
					"mssql":      mssql.Factory,
					"mysql":      mysql.Factory,
					"ssh":        ssh.Factory,
				},
				ShutdownCh:  command.MakeShutdownCh(),
				SighupCh:    command.MakeSighupCh(),
				ReloadFuncs: map[string][]server.ReloadFunc{},
			}, nil
		},

		"ssh": func() (cli.Command, error) {
			return &command.SSHCommand{
				Meta: meta,
			}, nil
		},

		"path-help": func() (cli.Command, error) {
			return &command.PathHelpCommand{
				Meta: meta,
			}, nil
		},

		"auth": func() (cli.Command, error) {
			return &command.AuthCommand{
				Meta: meta,
				Handlers: map[string]command.AuthHandler{
					"github":   &credGitHub.CLIHandler{},
					"userpass": &credUserpass.CLIHandler{},
					"ldap":     &credLdap.CLIHandler{},
					"cert":     &credCert.CLIHandler{},
					"marathon": &credMarathon.CLIHandler{},
				},
			}, nil
		},

		"auth-enable": func() (cli.Command, error) {
			return &command.AuthEnableCommand{
				Meta: meta,
			}, nil
		},

		"auth-disable": func() (cli.Command, error) {
			return &command.AuthDisableCommand{
				Meta: meta,
			}, nil
		},

		"audit-list": func() (cli.Command, error) {
			return &command.AuditListCommand{
				Meta: meta,
			}, nil
		},

		"audit-disable": func() (cli.Command, error) {
			return &command.AuditDisableCommand{
				Meta: meta,
			}, nil
		},

		"audit-enable": func() (cli.Command, error) {
			return &command.AuditEnableCommand{
				Meta: meta,
			}, nil
		},

		"key-status": func() (cli.Command, error) {
			return &command.KeyStatusCommand{
				Meta: meta,
			}, nil
		},

		"policies": func() (cli.Command, error) {
			return &command.PolicyListCommand{
				Meta: meta,
			}, nil
		},

		"policy-delete": func() (cli.Command, error) {
			return &command.PolicyDeleteCommand{
				Meta: meta,
			}, nil
		},

		"policy-write": func() (cli.Command, error) {
			return &command.PolicyWriteCommand{
				Meta: meta,
			}, nil
		},

		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				Meta: meta,
			}, nil
		},

		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: meta,
			}, nil
		},

		"write": func() (cli.Command, error) {
			return &command.WriteCommand{
				Meta: meta,
			}, nil
		},

		"delete": func() (cli.Command, error) {
			return &command.DeleteCommand{
				Meta: meta,
			}, nil
		},

		"rekey": func() (cli.Command, error) {
			return &command.RekeyCommand{
				Meta: meta,
			}, nil
		},

		"generate-root": func() (cli.Command, error) {
			return &command.GenerateRootCommand{
				Meta: meta,
			}, nil
		},

		"renew": func() (cli.Command, error) {
			return &command.RenewCommand{
				Meta: meta,
			}, nil
		},

		"revoke": func() (cli.Command, error) {
			return &command.RevokeCommand{
				Meta: meta,
			}, nil
		},

		"seal": func() (cli.Command, error) {
			return &command.SealCommand{
				Meta: meta,
			}, nil
		},

		"status": func() (cli.Command, error) {
			return &command.StatusCommand{
				Meta: meta,
			}, nil
		},

		"unseal": func() (cli.Command, error) {
			return &command.UnsealCommand{
				Meta: meta,
			}, nil
		},

		"step-down": func() (cli.Command, error) {
			return &command.StepDownCommand{
				Meta: meta,
			}, nil
		},

		"mount": func() (cli.Command, error) {
			return &command.MountCommand{
				Meta: meta,
			}, nil
		},

		"mounts": func() (cli.Command, error) {
			return &command.MountsCommand{
				Meta: meta,
			}, nil
		},

		"mount-tune": func() (cli.Command, error) {
			return &command.MountTuneCommand{
				Meta: meta,
			}, nil
		},

		"remount": func() (cli.Command, error) {
			return &command.RemountCommand{
				Meta: meta,
			}, nil
		},

		"rotate": func() (cli.Command, error) {
			return &command.RotateCommand{
				Meta: meta,
			}, nil
		},

		"unmount": func() (cli.Command, error) {
			return &command.UnmountCommand{
				Meta: meta,
			}, nil
		},

		"token-create": func() (cli.Command, error) {
			return &command.TokenCreateCommand{
				Meta: meta,
			}, nil
		},

		"token-lookup": func() (cli.Command, error) {
			return &command.TokenLookupCommand{
				Meta: meta,
			}, nil
		},

		"token-renew": func() (cli.Command, error) {
			return &command.TokenRenewCommand{
				Meta: meta,
			}, nil
		},

		"token-revoke": func() (cli.Command, error) {
			return &command.TokenRevokeCommand{
				Meta: meta,
			}, nil
		},

		"capabilities": func() (cli.Command, error) {
			return &command.CapabilitiesCommand{
				Meta: meta,
			}, nil
		},

		"version": func() (cli.Command, error) {
			versionInfo := version.GetVersion()

			return &command.VersionCommand{
				VersionInfo: versionInfo,
				Ui:          meta.Ui,
			}, nil
		},
	}
}
