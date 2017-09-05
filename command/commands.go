package command

import (
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

	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
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
	physMSSQL "github.com/hashicorp/vault/physical/mssql"
	physMySQL "github.com/hashicorp/vault/physical/mysql"
	physPostgreSQL "github.com/hashicorp/vault/physical/postgresql"
	physS3 "github.com/hashicorp/vault/physical/s3"
	physSwift "github.com/hashicorp/vault/physical/swift"
	physZooKeeper "github.com/hashicorp/vault/physical/zookeeper"
)

// Commands is the mapping of all the available commands.
var Commands map[string]cli.CommandFactory

func init() {
	ui := &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		Ui: &cli.BasicUi{
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	authHandlers := map[string]AuthHandler{
		"aws":    &credAws.CLIHandler{},
		"cert":   &credCert.CLIHandler{},
		"github": &credGitHub.CLIHandler{},
		"ldap":   &credLdap.CLIHandler{},
		"okta":   &credOkta.CLIHandler{},
		"radius": &credUserpass.CLIHandler{
			DefaultMount: "radius",
		},
		"token": &credToken.CLIHandler{},
		"userpass": &credUserpass.CLIHandler{
			DefaultMount: "userpass",
		},
	}

	Commands = map[string]cli.CommandFactory{
		"audit-disable": func() (cli.Command, error) {
			return &AuditDisableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"audit-enable": func() (cli.Command, error) {
			return &AuditEnableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"audit-list": func() (cli.Command, error) {
			return &AuditListCommand{
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
				Handlers: authHandlers,
			}, nil
		},
		"auth-disable": func() (cli.Command, error) {
			return &AuthDisableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth-enable": func() (cli.Command, error) {
			return &AuthEnableCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"auth-help": func() (cli.Command, error) {
			return &AuthHelpCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
				Handlers: authHandlers,
			}, nil
		},
		"auth-list": func() (cli.Command, error) {
			return &AuthListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"capabilities": func() (cli.Command, error) {
			return &CapabilitiesCommand{
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
		"generate-root": func() (cli.Command, error) {
			return &GenerateRootCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"init": func() (cli.Command, error) {
			return &InitCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"key-status": func() (cli.Command, error) {
			return &KeyStatusCommand{
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
		"mount": func() (cli.Command, error) {
			return &MountCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"mounts": func() (cli.Command, error) {
			return &MountsCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"mount-tune": func() (cli.Command, error) {
			return &MountTuneCommand{
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
		"policies": func() (cli.Command, error) {
			return &PolicyListCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy-delete": func() (cli.Command, error) {
			return &PolicyDeleteCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"policy-write": func() (cli.Command, error) {
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
		"rekey": func() (cli.Command, error) {
			return &RekeyCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"remount": func() (cli.Command, error) {
			return &RemountCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"renew": func() (cli.Command, error) {
			return &RenewCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"revoke": func() (cli.Command, error) {
			return &RevokeCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"rotate": func() (cli.Command, error) {
			return &RotateCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"seal": func() (cli.Command, error) {
			return &SealCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
				AuditBackends: map[string]audit.Factory{
					"file":   auditFile.Factory,
					"socket": auditSocket.Factory,
					"syslog": auditSyslog.Factory,
				},
				CredentialBackends: map[string]logical.Factory{
					"app-id":   credAppId.Factory,
					"approle":  credAppRole.Factory,
					"aws":      credAws.Factory,
					"cert":     credCert.Factory,
					"gcp":      credGcp.Factory,
					"github":   credGitHub.Factory,
					"ldap":     credLdap.Factory,
					"okta":     credOkta.Factory,
					"plugin":   plugin.Factory,
					"radius":   credRadius.Factory,
					"userpass": credUserpass.Factory,
				},
				LogicalBackends: map[string]logical.Factory{
					"aws":        aws.Factory,
					"cassandra":  cassandra.Factory,
					"consul":     consul.Factory,
					"database":   database.Factory,
					"mongodb":    mongodb.Factory,
					"mssql":      mssql.Factory,
					"mysql":      mysql.Factory,
					"pki":        pki.Factory,
					"plugin":     plugin.Factory,
					"postgresql": postgresql.Factory,
					"rabbitmq":   rabbitmq.Factory,
					"ssh":        ssh.Factory,
					"totp":       totp.Factory,
					"transit":    transit.Factory,
				},
				PhysicalBackends: map[string]physical.Factory{
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
					"mssql":                  physMSSQL.NewMSSQLBackend,
					"mysql":                  physMySQL.NewMySQLBackend,
					"postgresql":             physPostgreSQL.NewPostgreSQLBackend,
					"s3":                     physS3.NewS3Backend,
					"swift":                  physSwift.NewSwiftBackend,
					"zookeeper":              physZooKeeper.NewZooKeeperBackend,
				},
				ShutdownCh: MakeShutdownCh(),
				SighupCh:   MakeSighupCh(),
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
		"step-down": func() (cli.Command, error) {
			return &StepDownCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token-create": func() (cli.Command, error) {
			return &TokenCreateCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token-lookup": func() (cli.Command, error) {
			return &TokenLookupCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token-renew": func() (cli.Command, error) {
			return &TokenRenewCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"token-revoke": func() (cli.Command, error) {
			return &TokenRevokeCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"unseal": func() (cli.Command, error) {
			return &UnsealCommand{
				BaseCommand: &BaseCommand{
					UI: ui,
				},
			}, nil
		},
		"unmount": func() (cli.Command, error) {
			return &UnmountCommand{
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
