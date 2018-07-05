package command

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ad "github.com/hashicorp/vault-plugin-secrets-ad/plugin"
	gcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
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

	credAzure "github.com/hashicorp/vault-plugin-auth-azure"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credJWT "github.com/hashicorp/vault-plugin-auth-jwt"
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
	physFoundationDB "github.com/hashicorp/vault/physical/foundationdb"
	physGCS "github.com/hashicorp/vault/physical/gcs"
	physInmem "github.com/hashicorp/vault/physical/inmem"
	physManta "github.com/hashicorp/vault/physical/manta"
	physMSSQL "github.com/hashicorp/vault/physical/mssql"
	physMySQL "github.com/hashicorp/vault/physical/mysql"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	physS3 "github.com/hashicorp/vault/physical/s3"
	physSpanner "github.com/hashicorp/vault/physical/spanner"
	physSwift "github.com/hashicorp/vault/physical/swift"
	physZooKeeper "github.com/hashicorp/vault/physical/zookeeper"
)

const (
	// EnvVaultCLINoColor is an env var that toggles colored UI output.
	EnvVaultCLINoColor = `VAULT_CLI_NO_COLOR`
	// EnvVaultFormat is the output format
	EnvVaultFormat = `VAULT_FORMAT`

	// flagNameAuditNonHMACRequestKeys is the flag name used for auth/secrets enable
	flagNameAuditNonHMACRequestKeys = "audit-non-hmac-request-keys"
	// flagNameAuditNonHMACResponseKeys is the flag name used for auth/secrets enable
	flagNameAuditNonHMACResponseKeys = "audit-non-hmac-response-keys"
	// flagNameDescription is the flag name used for tuning the secret and auth mount description parameter
	flagNameDescription = "description"
	// flagListingVisibility is the flag to toggle whether to show the mount in the UI-specific listing endpoint
	flagNameListingVisibility = "listing-visibility"
	// flagNamePassthroughRequestHeaders is the flag name used to set passthrough request headers to the backend
	flagNamePassthroughRequestHeaders = "passthrough-request-headers"
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
		"azure":      credAzure.Factory,
		"centrify":   credCentrify.Factory,
		"cert":       credCert.Factory,
		"gcp":        credGcp.Factory,
		"github":     credGitHub.Factory,
		"jwt":        credJWT.Factory,
		"kubernetes": credKube.Factory,
		"ldap":       credLdap.Factory,
		"okta":       credOkta.Factory,
		"plugin":     plugin.Factory,
		"radius":     credRadius.Factory,
		"userpass":   credUserpass.Factory,
	}

	logicalBackends = map[string]logical.Factory{
		"ad":         ad.Factory,
		"aws":        aws.Factory,
		"cassandra":  cassandra.Factory,
		"consul":     consul.Factory,
		"database":   database.Factory,
		"gcp":        gcp.Factory,
		"kv":         kv.Factory,
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
		"foundationdb":           physFoundationDB.NewFDBBackend,
		"gcs":                    physGCS.NewBackend,
		"inmem_ha":               physInmem.NewInmemHA,
		"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
		"inmem_transactional":    physInmem.NewTransactionalInmem,
		"inmem":                  physInmem.NewInmem,
		"manta":                  physManta.NewMantaBackend,
		"mssql":                  physMSSQL.NewMSSQLBackend,
		"mysql":                  physMySQL.NewMySQLBackend,
		"postgresql":             physPostgreSQL.NewPostgreSQLBackend,
		"s3":                     physS3.NewS3Backend,
		"spanner":                physSpanner.NewBackend,
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

func initCommands(ui, serverCmdUi cli.Ui, runOpts *RunOptions) {
	loginHandlers := map[string]LoginHandler{
		"aws":      &credAws.CLIHandler{},
		"centrify": &credCentrify.CLIHandler{},
		"cert":     &credCert.CLIHandler{},
		"gcp":      &credGcp.CLIHandler{},
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

	getBaseCommand := func() *BaseCommand {
		return &BaseCommand{
			UI:          ui,
			tokenHelper: runOpts.TokenHelper,
			flagAddress: runOpts.Address,
			client:      runOpts.Client,
		}
	}

	Commands = map[string]cli.CommandFactory{
		"agent": func() (cli.Command, error) {
			return &AgentCommand{
				BaseCommand: &BaseCommand{
					UI: serverCmdUi,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
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
				Handlers:    loginHandlers,
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
				Handlers:    loginHandlers,
			}, nil
		},
		"auth list": func() (cli.Command, error) {
			return &AuthListCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"delete": func() (cli.Command, error) {
			return &DeleteCommand{
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
				Handlers:    loginHandlers,
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
		"operator unseal": func() (cli.Command, error) {
			return &OperatorUnsealCommand{
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"path-help": func() (cli.Command, error) {
			return &PathHelpCommand{
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
			return &PluginRegisterCommand{
				BaseCommand: getBaseCommand(),
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
		"policy fmt": func() (cli.Command, error) {
			return &PolicyFmtCommand{
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
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: &BaseCommand{
					UI:          serverCmdUi,
					tokenHelper: runOpts.TokenHelper,
					flagAddress: runOpts.Address,
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
				BaseCommand: getBaseCommand(),
			}, nil
		},
		"status": func() (cli.Command, error) {
			return &StatusCommand{
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
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"audit-enable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "audit-enable",
				New: "audit enable",
				UI:  ui,
				Command: &AuditEnableCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"audit-list": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "audit-list",
				New: "audit list",
				UI:  ui,
				Command: &AuditListCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"auth-disable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "auth-disable",
				New: "auth disable",
				UI:  ui,
				Command: &AuthDisableCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"auth-enable": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "auth-enable",
				New: "auth enable",
				UI:  ui,
				Command: &AuthEnableCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"capabilities": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "capabilities",
				New: "token capabilities",
				UI:  ui,
				Command: &TokenCapabilitiesCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"generate-root": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "generate-root",
				New: "operator generate-root",
				UI:  ui,
				Command: &OperatorGenerateRootCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"init": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "init",
				New: "operator init",
				UI:  ui,
				Command: &OperatorInitCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"key-status": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "key-status",
				New: "operator key-status",
				UI:  ui,
				Command: &OperatorKeyStatusCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"renew": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "renew",
				New: "lease renew",
				UI:  ui,
				Command: &LeaseRenewCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"revoke": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "revoke",
				New: "lease revoke",
				UI:  ui,
				Command: &LeaseRevokeCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"mount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mount",
				New: "secrets enable",
				UI:  ui,
				Command: &SecretsEnableCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"mount-tune": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mount-tune",
				New: "secrets tune",
				UI:  ui,
				Command: &SecretsTuneCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"mounts": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "mounts",
				New: "secrets list",
				UI:  ui,
				Command: &SecretsListCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"policies": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policies",
				New: "policy read\" or \"vault policy list", // lol
				UI:  ui,
				Command: &PoliciesDeprecatedCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"policy-delete": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policy-delete",
				New: "policy delete",
				UI:  ui,
				Command: &PolicyDeleteCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"policy-write": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "policy-write",
				New: "policy write",
				UI:  ui,
				Command: &PolicyWriteCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"rekey": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "rekey",
				New: "operator rekey",
				UI:  ui,
				Command: &OperatorRekeyCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"remount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "remount",
				New: "secrets move",
				UI:  ui,
				Command: &SecretsMoveCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"rotate": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "rotate",
				New: "operator rotate",
				UI:  ui,
				Command: &OperatorRotateCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"seal": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "seal",
				New: "operator seal",
				UI:  ui,
				Command: &OperatorSealCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"step-down": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "step-down",
				New: "operator step-down",
				UI:  ui,
				Command: &OperatorStepDownCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"token-create": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-create",
				New: "token create",
				UI:  ui,
				Command: &TokenCreateCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"token-lookup": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-lookup",
				New: "token lookup",
				UI:  ui,
				Command: &TokenLookupCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"token-renew": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-renew",
				New: "token renew",
				UI:  ui,
				Command: &TokenRenewCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"token-revoke": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "token-revoke",
				New: "token revoke",
				UI:  ui,
				Command: &TokenRevokeCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"unmount": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "unmount",
				New: "secrets disable",
				UI:  ui,
				Command: &SecretsDisableCommand{
					BaseCommand: getBaseCommand(),
				},
			}, nil
		},

		"unseal": func() (cli.Command, error) {
			return &DeprecatedCommand{
				Old: "unseal",
				New: "operator unseal",
				UI:  ui,
				Command: &OperatorUnsealCommand{
					BaseCommand: getBaseCommand(),
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
