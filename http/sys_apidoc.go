package http

import (
	"context"
	"encoding/json"
	"net/http"
	"path"

	log "github.com/hashicorp/go-hclog"
	gcp "github.com/hashicorp/vault-plugin-secrets-gcp/plugin"
	kv "github.com/hashicorp/vault-plugin-secrets-kv"
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
	"github.com/hashicorp/vault/vault"

	credAzure "github.com/hashicorp/vault-plugin-auth-azure/plugin"
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
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
)

var (
	backends = map[string]logical.Factory{
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
		"radius":     credRadius.Factory,
		"userpass":   credUserpass.Factory,

		// logical
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
		"postgresql": postgresql.Factory,
		"rabbitmq":   rabbitmq.Factory,
		"ssh":        ssh.Factory,
		"totp":       totp.Factory,
		"transit":    transit.Factory,

		// sys
		"sys": func(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
			return vault.NewSystemBackend(&vault.Core{}, nil).Backend, nil
		},
	}
)

func handleSysApiDoc(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var doc interface{}

		if r.Method != "GET" {
			respondError(w, http.StatusMethodNotAllowed, nil)
			return
		}
		_, path := path.Split(r.URL.Path)

		if path != "" {
			doc = genDoc(path)
		} else {
			doc = genAllDocs()
		}

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		enc.Encode(doc)
	})
}

func genAllDocs() map[string]*oas.OASDoc {
	docs := make(map[string]*oas.OASDoc)

	for b := range backends {
		d := genDoc(b)
		if d != nil {
			docs[b] = d
		}
	}
	return docs
}

func genDoc(mount string) *oas.OASDoc {
	conf := &logical.BackendConfig{
		Logger: log.NewNullLogger(),
		System: logical.StaticSystemView{},
	}

	f, ok := backends[mount]
	if !ok {
		return nil
	}

	b, err := f(context.Background(), conf)
	if err != nil {
		return nil
	}

	return b.Describe()
}
