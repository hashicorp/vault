// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package minimal

import (
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	"github.com/hashicorp/vault/audit"
	auditFile "github.com/hashicorp/vault/builtin/audit/file"
	auditSocket "github.com/hashicorp/vault/builtin/audit/socket"
	auditSyslog "github.com/hashicorp/vault/builtin/audit/syslog"
	logicalDb "github.com/hashicorp/vault/builtin/logical/database"
	"github.com/hashicorp/vault/builtin/plugin"
	"github.com/hashicorp/vault/helper/builtinplugins"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/physical/inmem"
	"github.com/hashicorp/vault/vault"
	"github.com/mitchellh/copystructure"
	"github.com/mitchellh/go-testing-interface"
)

// NewTestSoloCluster is a simpler version of NewTestCluster that only creates
// single-node clusters.  It is intentionally minimalist, if you need something
// from vault.TestClusterOptions, use NewTestCluster instead.  It should work fine
// with a nil config argument.  There is no need to call Start or Cleanup or
// TestWaitActive on the resulting cluster.
func NewTestSoloCluster(t testing.T, config *vault.CoreConfig) *vault.TestCluster {
	logger := corehelpers.NewTestLogger(t)

	mycfg := &vault.CoreConfig{}

	if config != nil {
		// It's rude to modify an input argument as a side-effect
		copy, err := copystructure.Copy(config)
		if err != nil {
			t.Fatal(err)
		}
		mycfg = copy.(*vault.CoreConfig)
	}
	if mycfg.Physical == nil {
		// Don't use NewTransactionalInmem because that would enable replication,
		// which we don't care about in our case (use NewTestCluster for that.)
		inm, err := inmem.NewInmem(nil, logger)
		if err != nil {
			t.Fatal(err)
		}
		mycfg.Physical = inm
	}
	if mycfg.CredentialBackends == nil {
		mycfg.CredentialBackends = map[string]logical.Factory{
			"plugin": plugin.Factory,
		}
	}
	if mycfg.LogicalBackends == nil {
		mycfg.LogicalBackends = map[string]logical.Factory{
			"plugin":   plugin.Factory,
			"database": logicalDb.Factory,
			// This is also available in the plugin catalog, but is here due to the need to
			// automatically mount it.
			"kv": logicalKv.Factory,
		}
	}
	if mycfg.AuditBackends == nil {
		mycfg.AuditBackends = map[string]audit.Factory{
			"file":   auditFile.Factory,
			"socket": auditSocket.Factory,
			"syslog": auditSyslog.Factory,
		}
	}
	if mycfg.BuiltinRegistry == nil {
		mycfg.BuiltinRegistry = builtinplugins.Registry
	}

	cluster := vault.NewTestCluster(t, mycfg, &vault.TestClusterOptions{
		NumCores:    1,
		HandlerFunc: http.Handler,
		Logger:      logger,
	})
	t.Cleanup(cluster.Cleanup)
	return cluster
}
