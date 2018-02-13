package command

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/vault/audit"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/version"
	"github.com/mitchellh/cli"

	"github.com/hashicorp/vault/builtin/logical/aws"
	"github.com/hashicorp/vault/builtin/logical/cassandra"
	"github.com/hashicorp/vault/builtin/logical/consul"
	"github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/logical/mongodb"
	"github.com/hashicorp/vault/builtin/logical/mssql"
	"github.com/hashicorp/vault/builtin/logical/mysql"
	"github.com/hashicorp/vault/builtin/logical/nomad"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/postgresql"
	"github.com/hashicorp/vault/builtin/logical/rabbitmq"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/builtin/plugin"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSocket "github.com/hashicorp/vault/builtin/audit/socket"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"

	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credKube "github.com/hashicorp/vault-plugin-auth-kubernetes"
	credAppId "github.com/hashicorp/vault/builtin/credential/app-id"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credAws "github.com/hashicorp/vault/builtin/credential/aws"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credGitHub "github.com/hashicorp/vault/builtin/credential/github"
	credLdap "github.com/hashicorp/vault/builtin/credential/ldap"
	credOkta "github.com/hashicorp/vault/builtin/credential/okta"
	credRadius "github.com/hashicorp/vault/builtin/credential/radius"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"

	physAzure "github.com/hashicorp/vault/physical/azure"
	physCassandra "github.com/hashicorp/vault/physical/cassandra"
	physCockroachDB "github.com/hashicorp/vault/physical/cockroachdb"
	physConsul "github.com/hashicorp/vault/physical/consul"
	physCouchDB "github.com/hashicorp/vault/physical/couchdb"
	physDynamoDB "github.com/hashicorp/vault/physical/dynamodb"
	physEtcd "github.com/hashicorp/vault/physical/etcd"
	physFile "github.com/hashicorp/vault/physical/file"
	physGCS "github.com/hashicorp/vault/physical/gcs"
	physInmem "github.com/hashicorp/vault/physical/inmem"
	physManta "github.com/hashicorp/vault/physical/manta"
	physMSSQL "github.com/hashicorp/vault/physical/mssql"
	physMySQL "github.com/hashicorp/vault/physical/mysql"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	physS3 "github.com/hashicorp/vault/physical/s3"
	physSwift "github.com/hashicorp/vault/physical/swift"
	physZooKeeper "github.com/hashicorp/vault/physical/zookeeper"
)

const (
	// EnvVaultCLINoColor is an env var that toggles colored UI output.
	EnvVaultCLINoColor = `VAULT_CLI_NO_COLOR`
	// EnvVaultFormat is the output format
	EnvVaultFormat = `VAULT_FORMAT`
)

var (
	auditBackends = map[string]audit.Factory{
		"file":   auditFile.Factory,
		"socket": auditSocket.Factory,
		"syslog": auditSyslog.Factory,
	}

	credentialBackends = map[string]logical.Factory{
		"app-id":     credAppId.Factory,
		"approle":    credAppRole.Factory,
		"aws":        credAws.Factory,
		"centrify":   credCentrify.Factory,
		"cert":       credCert.Factory,
		"gcp":        credGcp.Factory,
		"github":     credGitHub.Factory,
		"kubernetes": credKube.Factory,
		"ldap":       credLdap.Factory,
		"okta":       credOkta.Factory,
		"plugin":     plugin.Factory,
		"radius":     credRadius.Factory,
		"userpass":   credUserpass.Factory,
	}

	logicalBackends = map[string]logical.Factory{
		"aws":        aws.Factory,
		"cassandra":  cassandra.Factory,
		"consul":     consul.Factory,
		"database":   database.Factory,
		"mongodb":    mongodb.Factory,
		"mssql":      mssql.Factory,
		"mysql":      mysql.Factory,
		"nomad":      nomad.Factory,
		"pki":        pki.Factory,
		"plugin":     plugin.Factory,
		"postgresql": postgresql.Factory,
		"rabbitmq":   rabbitmq.Factory,
		"ssh":        ssh.Factory,
		"totp":       totp.Factory,
		"transit":    transit.Factory,
	}

	physicalBackends = map[string]physical.Factory{
		"azure":                  physAzure.NewAzureBackend,
		"cassandra":              physCassandra.NewCassandraBackend,
		"cockroachdb":            physCockroachDB.NewCockroachDBBackend,
		"consul":                 physConsul.NewConsulBackend,
		"couchdb_transactional":  physCouchDB.NewTransactionalCouchDBBackend,
		"couchdb":                physCouchDB.NewCouchDBBackend,
		"dynamodb":               physDynamoDB.NewDynamoDBBackend,
		"etcd":                   physEtcd.NewEtcdBackend,
		"file_transactional":     physFile.NewTransactionalFileBackend,
		"file":                   physFile.NewFileBackend,
		"gcs":                    physGCS.NewGCSBackend,
		"inmem_ha":               physInmem.NewInmemHA,
		"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
		"inmem_transactional":    physInmem.NewTransactionalInmem,
		"inmem":                  physInmem.NewInmem,
		"manta":                  physManta.NewMantaBackend,
		"mssql":                  physMSSQL.NewMSSQLBackend,
		"mysql":                  physMySQL.NewMySQLBackend,
		"postgresql":             physPostgreSQL.NewPostgreSQLBackend,
		"s3":                     physS3.NewS3Backend,
		"swift":                  physSwift.NewSwiftBackend,
		"zookeeper":              physZooKeeper.NewZooKeeperBackend,
	}
)

// DeprecatedCommand is a command that wraps an existing command and prints a
// deprecation notice and points the user to the new command. Deprecated
// commands are always hidden from help output.
type DeprecatedCommand struct {
	cli.Command
	UI cli.Ui

	// Old is the old command name, New is the new command name.
	Old, New string
}

// Help wraps the embedded Help command and prints a warning about deprecations.
func (c *DeprecatedCommand) Help() string {
	c.warn()
	return c.Command.Help()
}

// Run wraps the embedded Run command and prints a warning about deprecation.
func (c *DeprecatedCommand) Run(args []string) int {
	if Format(c.UI) == "table" {
		c.warn()
	}
	return c.Command.Run(args)
}

func (c *DeprecatedCommand) warn() {
	c.UI.Warn(wrapAtLength(fmt.Sprintf(
		"WARNING! The \"vault %s\" command is deprecated. Please use \"vault %s\" "+
			"instead. This command will be removed in Vault 0.11 (or later).",
		c.Old,
		c.New)))
	c.UI.Warn("")
}

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory
var DeprecatedCommands map[string]cli.CommandFactory

func initCommands(ui, serverCmdUi cli.Ui) {
	loginHandlers := map[string]LoginHandler{
		"aws":      &credAws.CLIHandler{},
		"centrify": &credCentrify.CLIHandler{},
		"cert":     &credCert.CLIHandler{},
		"github":   &credGitHub.CLIHandler{},
		"ldap":     &credLdap.CLIHandler{},
		"okta":     &credOkta.CLIHandler{},
		"radius": &credUserpass.CLIHandler{
			DefaultMount: "radius",
		},
		"token": &credToken.CLIHandler{},
		"userpass": &credUserpass.CLIHandler{
			DefaultMount: "userpass",
		},
	}

	Commands = map[string]cli.CommandFactory{
		"audit": func() (cli.Command, error) {
			return &AuditCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"audit disable": func() (cli.Command, error) {
			return &AuditDisableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"audit enable": func() (cli.Command, error) {
			return &AuditEnableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"audit list": func() (cli.Command, error) {
			return &AuditListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth tune": func() (cli.Command, error) {
			return &AuthTuneCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth": func() (cli.Command, error) {
			return &AuthCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
				Handlers: loginHandlers,
			}, nil
		},
		"auth disable": func() (cli.Command, error) {
			return &AuthDisableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth enable": func() (cli.Command, error) {
			return &AuthEnableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth help": func() (cli.Command, error) {
			return &AuthHelpCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
				Handlers: loginHandlers,
			}, nil
		},
		"auth list": func() (cli.Command, error) {
			return &AuthListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"delete": func() (cli.Command, error) {
			return &DeleteCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"lease": func() (cli.Command, error) {
			return &LeaseCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"lease renew": func() (cli.Command, error) {
			return &LeaseRenewCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"lease revoke": func() (cli.Command, error) {
			return &LeaseRevokeCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"login": func() (cli.Command, error) {
			return &LoginCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
				Handlers: loginHandlers,
			}, nil
		},
		"operator": func() (cli.Command, error) {
			return &OperatorCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator generate-root": func() (cli.Command, error) {
			return &OperatorGenerateRootCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator init": func() (cli.Command, error) {
			return &OperatorInitCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator key-status": func() (cli.Command, error) {
			return &OperatorKeyStatusCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator rekey": func() (cli.Command, error) {
			return &OperatorRekeyCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator rotate": func() (cli.Command, error) {
			return &OperatorRotateCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator seal": func() (cli.Command, error) {
			return &OperatorSealCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator step-down": func() (cli.Command, error) {
			return &OperatorStepDownCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"operator unseal": func() (cli.Command, error) {
			return &OperatorUnsealCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"path-help": func() (cli.Command, error) {
			return &PathHelpCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy": func() (cli.Command, error) {
			return &PolicyCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy delete": func() (cli.Command, error) {
			return &PolicyDeleteCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy fmt": func() (cli.Command, error) {
			return &PolicyFmtCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy list": func() (cli.Command, error) {
			return &PolicyListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy read": func() (cli.Command, error) {
			return &PolicyReadCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy write": func() (cli.Command, error) {
			return &PolicyWriteCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"read": func() (cli.Command, error) {
			return &ReadCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets": func() (cli.Command, error) {
			return &SecretsCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets disable": func() (cli.Command, error) {
			return &SecretsDisableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets enable": func() (cli.Command, error) {
			return &SecretsEnableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets list": func() (cli.Command, error) {
			return &SecretsListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets move": func() (cli.Command, error) {
			return &SecretsMoveCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"secrets tune": func() (cli.Command, error) {
			return &SecretsTuneCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: &BaseCommand{
					UI: serverCmdUi,
				},
				AuditBackends:      auditBackends,
				CredentialBackends: credentialBackends,
				LogicalBackends:    logicalBackends,
				PhysicalBackends:   physicalBackends,
				ShutdownCh:         MakeShutdownCh(),
				SighupCh:           MakeSighupCh(),
			}, nil
		},
		"ssh": func() (cli.Command, error) {
			return &SSHCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token": func() (cli.Command, error) {
			return &TokenCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token create": func() (cli.Command, error) {
			return &TokenCreateCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token capabilities": func() (cli.Command, error) {
			return &TokenCapabilitiesCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token lookup": func() (cli.Command, error) {
			return &TokenLookupCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token renew": func() (cli.Command, error) {
			return &TokenRenewCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token revoke": func() (cli.Command, error) {
			return &TokenRevokeCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"unwrap": func() (cli.Command, error) {
			return &UnwrapCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				VersionInfo: version.GetVersion(),
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"write": func() (cli.Command, error) {
			return &WriteCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
	}

	// Deprecated commands
	//
	// TODO: Remove not before 0.11.0
	DeprecatedCommands = map[string]cli.CommandFactory{
		"audit-disable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "audit-disable",
				New: "audit disable",
				UI:  ui,
				Command: &AuditDisableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"audit-enable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "audit-enable",
				New: "audit enable",
				UI:  ui,
				Command: &AuditEnableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"audit-list": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "audit-list",
				New: "audit list",
				UI:  ui,
				Command: &AuditListCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"auth-disable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "auth-disable",
				New: "auth disable",
				UI:  ui,
				Command: &AuthDisableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"auth-enable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "auth-enable",
				New: "auth enable",
				UI:  ui,
				Command: &AuthEnableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"capabilities": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "capabilities",
				New: "token capabilities",
				UI:  ui,
				Command: &TokenCapabilitiesCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"generate-root": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "generate-root",
				New: "operator generate-root",
				UI:  ui,
				Command: &OperatorGenerateRootCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"init": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "init",
				New: "operator init",
				UI:  ui,
				Command: &OperatorInitCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"key-status": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "key-status",
				New: "operator key-status",
				UI:  ui,
				Command: &OperatorKeyStatusCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"renew": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "renew",
				New: "lease renew",
				UI:  ui,
				Command: &LeaseRenewCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"revoke": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "revoke",
				New: "lease revoke",
				UI:  ui,
				Command: &LeaseRevokeCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"mount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mount",
				New: "secrets enable",
				UI:  ui,
				Command: &SecretsEnableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"mount-tune": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mount-tune",
				New: "secrets tune",
				UI:  ui,
				Command: &SecretsTuneCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"mounts": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mounts",
				New: "secrets list",
				UI:  ui,
				Command: &SecretsListCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"policies": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policies",
				New: "policy read\" or \"vault policy list", // lol
				UI:  ui,
				Command: &PoliciesDeprecatedCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"policy-delete": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policy-delete",
				New: "policy delete",
				UI:  ui,
				Command: &PolicyDeleteCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"policy-write": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policy-write",
				New: "policy write",
				UI:  ui,
				Command: &PolicyWriteCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"rekey": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "rekey",
				New: "operator rekey",
				UI:  ui,
				Command: &OperatorRekeyCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"remount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "remount",
				New: "secrets move",
				UI:  ui,
				Command: &SecretsMoveCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"rotate": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "rotate",
				New: "operator rotate",
				UI:  ui,
				Command: &OperatorRotateCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"seal": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "seal",
				New: "operator seal",
				UI:  ui,
				Command: &OperatorSealCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"step-down": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "step-down",
				New: "operator step-down",
				UI:  ui,
				Command: &OperatorStepDownCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"token-create": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-create",
				New: "token create",
				UI:  ui,
				Command: &TokenCreateCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"token-lookup": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-lookup",
				New: "token lookup",
				UI:  ui,
				Command: &TokenLookupCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"token-renew": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-renew",
				New: "token renew",
				UI:  ui,
				Command: &TokenRenewCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"token-revoke": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-revoke",
				New: "token revoke",
				UI:  ui,
				Command: &TokenRevokeCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"unmount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "unmount",
				New: "secrets disable",
				UI:  ui,
				Command: &SecretsDisableCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},

		"unseal": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "unseal",
				New: "operator unseal",
				UI:  ui,
				Command: &OperatorUnsealCommand{
					BaseCommand: &BaseCommand{
						UI: ui,
					},
				},
			}, nil
		},
	}

	// Add deprecated commands back to the main commands so they parse.
	for k, v := range DeprecatedCommands {
		if _, ok := Commands[k]; ok {
			// Can't deprecate an existing command...
			panic(fmt.Sprintf("command %q defined as deprecated and not at the same time!", k))
		}
		Commands[k] = v
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
