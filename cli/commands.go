package cli

import (
	"os"

	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSocket "github.com/hashicorp/vault/builtin/audit/socket"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"
	"github.com/hashicorp/vault/physical"
	"github.com/hashicorp/vault/version"

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
	physMSSQL "github.com/hashicorp/vault/physical/mssql"
	physMySQL "github.com/hashicorp/vault/physical/mysql"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	physS3 "github.com/hashicorp/vault/physical/s3"
	physSwift "github.com/hashicorp/vault/physical/swift"
	physZooKeeper "github.com/hashicorp/vault/physical/zookeeper"

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
	"github.com/hashicorp/vault/builtin/plugin"

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
			c := &command.ServerCommand{
				Meta: *metaPtr,
				AuditBackends: map[string]audit.Factory{
					"file":   auditFile.Factory,
					"syslog": auditSyslog.Factory,
					"socket": auditSocket.Factory,
				},
				CredentialBackends: map[string]logical.Factory{
					"approle":    credAppRole.Factory,
					"cert":       credCert.Factory,
					"aws":        credAws.Factory,
					"app-id":     credAppId.Factory,
					"gcp":        credGcp.Factory,
					"github":     credGitHub.Factory,
					"userpass":   credUserpass.Factory,
					"ldap":       credLdap.Factory,
					"okta":       credOkta.Factory,
					"radius":     credRadius.Factory,
					"kubernetes": credKube.Factory,
					"plugin":     plugin.Factory,
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
					"plugin":     plugin.Factory,
				},

				ShutdownCh: command.MakeShutdownCh(),
				SighupCh:   command.MakeSighupCh(),
			}

			c.PhysicalBackends = map[string]physical.Factory{
				"azure":                  physAzure.NewAzureBackend,
				"cassandra":              physCassandra.NewCassandraBackend,
				"cockroachdb":            physCockroachDB.NewCockroachDBBackend,
				"consul":                 physConsul.NewConsulBackend,
				"couchdb":                physCouchDB.NewCouchDBBackend,
				"couchdb_transactional":  physCouchDB.NewTransactionalCouchDBBackend,
				"dynamodb":               physDynamoDB.NewDynamoDBBackend,
				"etcd":                   physEtcd.NewEtcdBackend,
				"file":                   physFile.NewFileBackend,
				"file_transactional":     physFile.NewTransactionalFileBackend,
				"gcs":                    physGCS.NewGCSBackend,
				"inmem":                  physInmem.NewInmem,
				"inmem_ha":               physInmem.NewInmemHA,
				"inmem_transactional":    physInmem.NewTransactionalInmem,
				"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
				"mssql":                  physMSSQL.NewMSSQLBackend,
				"mysql":                  physMySQL.NewMySQLBackend,
				"postgresql":             physPostgreSQL.NewPostgreSQLBackend,
				"s3":                     physS3.NewS3Backend,
				"swift":                  physSwift.NewSwiftBackend,
				"zookeeper":              physZooKeeper.NewZooKeeperBackend,
			}

			return c, nil
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
