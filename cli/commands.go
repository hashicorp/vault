package cli

import (
	"os"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSocket "github.com/hashicorp/vault/builtin/audit/socket"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"
	"github.com/hashicorp/vault/version"

	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credRadius "github.com/hashicorp/vault/builtin/credential/radius"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"

	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/builtin/logical/cassandra"
	"github.com/hashicorp/vault/builtin/logical/consul"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/mongodb"
	"github.com/hashicorp/vault/builtin/logical/mssql"
	"github.com/hashicorp/vault/builtin/logical/mysql"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/postgresql"
	"github.com/hashicorp/vault/builtin/logical/rabbitmq"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/builtin/logical/transit"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/command"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/meta"
	"github.com/mitchellh/cli"
)

// Commands returns the mapping of CLI commands for Vault. The meta
// parameter lets you set meta options for all commands.
func Commands(metaPtr *meta.Meta) map[string]cli.CommandFactory {
	if metaPtr == nil {
		metaPtr = &meta.Meta{
			TokenHelper: command.DefaultTokenHelper,
		}
	}

	if metaPtr.Ui == nil {
		metaPtr.Ui = &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		}
	}

	return map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &command.InitCommand{
				Meta: *metaPtr,
			}, nil
		},

		"server": func() (cli.Command, error) {
			return &command.ServerCommand{
				Meta: *metaPtr,
				AuditBackends: map[string]audit.Factory{
					"file":   auditFile.Factory,
					"syslog": auditSyslog.Factory,
					"socket": auditSocket.Factory,
				},
				CredentialBackends: map[string]logical.Factory{
					"approle":  credAppRole.Factory,
					"cert":     credCert.Factory,
					"aws":      credAws.Factory,
					"app-id":   credAppId.Factory,
					"github":   credGitHub.Factory,
					"userpass": credUserpass.Factory,
					"ldap":     credLdap.Factory,
					"okta":     credOkta.Factory,
					"radius":   credRadius.Factory,
				},
				LogicalBackends: map[string]logical.Factory{
					"aws":        aws.Factory,
					"consul":     consul.Factory,
					"postgresql": postgresql.Factory,
					"cassandra":  cassandra.Factory,
					"pki":        pki.Factory,
					"transit":    transit.Factory,
					"mongodb":    mongodb.Factory,
					"mssql":      mssql.Factory,
					"mysql":      mysql.Factory,
					"ssh":        ssh.Factory,
					"rabbitmq":   rabbitmq.Factory,
					"database":   database.Factory,
					"totp":       totp.Factory,
				},
				ShutdownCh: command.MakeShutdownCh(),
				SighupCh:   command.MakeSighupCh(),
			}, nil
		},

		"ssh": func() (cli.Command, error) {
			return &command.SSHCommand{
				Meta: *metaPtr,
			}, nil
		},

		"path-help": func() (cli.Command, error) {
			return &command.PathHelpCommand{
				Meta: *metaPtr,
			}, nil
		},

		"auth": func() (cli.Command, error) {
			return &command.AuthCommand{
				Meta: *metaPtr,
				Handlers: map[string]command.AuthHandler{
					"github":   &credGitHub.CLIHandler{},
					"userpass": &credUserpass.CLIHandler{DefaultMount: "userpass"},
					"ldap":     &credLdap.CLIHandler{},
					"okta":     &credOkta.CLIHandler{},
					"cert":     &credCert.CLIHandler{},
					"aws":      &credAws.CLIHandler{},
					"radius":   &credUserpass.CLIHandler{DefaultMount: "radius"},
				},
			}, nil
		},

		"auth-enable": func() (cli.Command, error) {
			return &command.AuthEnableCommand{
				Meta: *metaPtr,
			}, nil
		},

		"auth-disable": func() (cli.Command, error) {
			return &command.AuthDisableCommand{
				Meta: *metaPtr,
			}, nil
		},

		"audit-list": func() (cli.Command, error) {
			return &command.AuditListCommand{
				Meta: *metaPtr,
			}, nil
		},

		"audit-disable": func() (cli.Command, error) {
			return &command.AuditDisableCommand{
				Meta: *metaPtr,
			}, nil
		},

		"audit-enable": func() (cli.Command, error) {
			return &command.AuditEnableCommand{
				Meta: *metaPtr,
			}, nil
		},

		"key-status": func() (cli.Command, error) {
			return &command.KeyStatusCommand{
				Meta: *metaPtr,
			}, nil
		},

		"policies": func() (cli.Command, error) {
			return &command.PolicyListCommand{
				Meta: *metaPtr,
			}, nil
		},

		"policy-delete": func() (cli.Command, error) {
			return &command.PolicyDeleteCommand{
				Meta: *metaPtr,
			}, nil
		},

		"policy-write": func() (cli.Command, error) {
			return &command.PolicyWriteCommand{
				Meta: *metaPtr,
			}, nil
		},

		"read": func() (cli.Command, error) {
			return &command.ReadCommand{
				Meta: *metaPtr,
			}, nil
		},

		"unwrap": func() (cli.Command, error) {
			return &command.UnwrapCommand{
				Meta: *metaPtr,
			}, nil
		},

		"list": func() (cli.Command, error) {
			return &command.ListCommand{
				Meta: *metaPtr,
			}, nil
		},

		"write": func() (cli.Command, error) {
			return &command.WriteCommand{
				Meta: *metaPtr,
			}, nil
		},

		"delete": func() (cli.Command, error) {
			return &command.DeleteCommand{
				Meta: *metaPtr,
			}, nil
		},

		"rekey": func() (cli.Command, error) {
			return &command.RekeyCommand{
				Meta: *metaPtr,
			}, nil
		},

		"generate-root": func() (cli.Command, error) {
			return &command.GenerateRootCommand{
				Meta: *metaPtr,
			}, nil
		},

		"renew": func() (cli.Command, error) {
			return &command.RenewCommand{
				Meta: *metaPtr,
			}, nil
		},

		"revoke": func() (cli.Command, error) {
			return &command.RevokeCommand{
				Meta: *metaPtr,
			}, nil
		},

		"seal": func() (cli.Command, error) {
			return &command.SealCommand{
				Meta: *metaPtr,
			}, nil
		},

		"status": func() (cli.Command, error) {
			return &command.StatusCommand{
				Meta: *metaPtr,
			}, nil
		},

		"unseal": func() (cli.Command, error) {
			return &command.UnsealCommand{
				Meta: *metaPtr,
			}, nil
		},

		"step-down": func() (cli.Command, error) {
			return &command.StepDownCommand{
				Meta: *metaPtr,
			}, nil
		},

		"mount": func() (cli.Command, error) {
			return &command.MountCommand{
				Meta: *metaPtr,
			}, nil
		},

		"mounts": func() (cli.Command, error) {
			return &command.MountsCommand{
				Meta: *metaPtr,
			}, nil
		},

		"mount-tune": func() (cli.Command, error) {
			return &command.MountTuneCommand{
				Meta: *metaPtr,
			}, nil
		},

		"remount": func() (cli.Command, error) {
			return &command.RemountCommand{
				Meta: *metaPtr,
			}, nil
		},

		"rotate": func() (cli.Command, error) {
			return &command.RotateCommand{
				Meta: *metaPtr,
			}, nil
		},

		"unmount": func() (cli.Command, error) {
			return &command.UnmountCommand{
				Meta: *metaPtr,
			}, nil
		},

		"token-create": func() (cli.Command, error) {
			return &command.TokenCreateCommand{
				Meta: *metaPtr,
			}, nil
		},

		"token-lookup": func() (cli.Command, error) {
			return &command.TokenLookupCommand{
				Meta: *metaPtr,
			}, nil
		},

		"token-renew": func() (cli.Command, error) {
			return &command.TokenRenewCommand{
				Meta: *metaPtr,
			}, nil
		},

		"token-revoke": func() (cli.Command, error) {
			return &command.TokenRevokeCommand{
				Meta: *metaPtr,
			}, nil
		},

		"capabilities": func() (cli.Command, error) {
			return &command.CapabilitiesCommand{
				Meta: *metaPtr,
			}, nil
		},

		"version": func() (cli.Command, error) {
			versionInfo := version.GetVersion()

			return &command.VersionCommand{
				VersionInfo: versionInfo,
				Ui:          metaPtr.Ui,
			}, nil
		},
	}
}
