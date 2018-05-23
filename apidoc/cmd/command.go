package cmd

import (
	"context"
	"encoding/json"
	"os"

	log "github.com/hashicorp/go-hclog"
	credAzure "github.com/hashicorp/vault-plugin-auth-azure/plugin"
	credCentrify "github.com/hashicorp/vault-plugin-auth-centrify"
	credGcp "github.com/hashicorp/vault-plugin-auth-gcp/plugin"
	credKube "github.com/hashicorp/vault-plugin-auth-kubernetes"
	gcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
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
	"github.com/hashicorp/vault/builtin/logical/nomad"
	"github.com/hashicorp/vault/builtin/logical/pki"
	"github.com/hashicorp/vault/builtin/logical/postgresql"
	"github.com/hashicorp/vault/builtin/logical/rabbitmq"
	"github.com/hashicorp/vault/builtin/logical/ssh"
	"github.com/hashicorp/vault/builtin/logical/totp"
	"github.com/hashicorp/vault/builtin/logical/transit"
	"github.com/hashicorp/vault/helper/oas"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"github.com/hashicorp/vault/vault"
)

type backender func(*logical.BackendConfig) (interface{}, error)

func getBackend(b backender) *framework.Backend {
	var ret *framework.Backend = nil

	return ret
}

func initLBs() map[string]logical.Backend {
	conf := &logical.BackendConfig{
		Logger: log.NewNullLogger(),
		System: logical.StaticSystemView{},
	}

	factories := map[string]logical.Factory{
		// logical
		"aws":       aws.Factory,
		"cassandra": cassandra.Factory,
		"consul":    consul.Factory,
		"database":  database.Factory,
		"gcp":       gcp.Factory,
		"kv":        kv.Factory,
		"mongodb":   mongodb.Factory,
		"mssql":     mssql.Factory,
		"mysql":     mysql.Factory,
		"nomad":     nomad.Factory,
		"pki":       pki.Factory,
		//"plugin":     plugin.Factory,
		"postgresql": postgresql.Factory,
		"rabbitmq":   rabbitmq.Factory,
		"ssh":        ssh.Factory,
		"totp":       totp.Factory,
		"transit":    transit.Factory,

		// credential
		"app-id":     credAppId.Factory,
		"approle":    credAppRole.Factory,
		"awscred":    credAws.Factory,
		"azure":      credAzure.Factory,
		"centrify":   credCentrify.Factory,
		"cert":       credCert.Factory,
		"gcpcred":    credGcp.Factory,
		"github":     credGitHub.Factory,
		"kubernetes": credKube.Factory,
		"ldap":       credLdap.Factory,
		"okta":       credOkta.Factory,
		//"plugin":     plugin.Factory,
		"radius":   credRadius.Factory,
		"userpass": credUserpass.Factory,
	}

	ret := make(map[string]logical.Backend)
	for key, factory := range factories {
		b, err := factory(context.Background(), conf)
		if err != nil {
			panic(err)
		}
		ret[key] = b
	}

	return ret
}

func Run() int {
	docs := make(map[string]*oas.OASDoc)

	backend := vault.NewSystemBackend(&vault.Core{}, nil).Backend
	docs["sys"] = backend.Describe()

	for mount, backend := range initLBs() {
		docs[mount] = backend.Describe()
	}

	output, _ := json.MarshalIndent(docs, "", "  ")
	os.Stdout.Write(output)

	return 0
}

func buildDoc(backend *framework.Backend) *oas.OASDoc {
	doc := oas.NewOASDoc()
	framework.DocumentPaths(backend, &doc)
	return &doc
}
