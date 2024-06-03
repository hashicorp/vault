// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package builtinplugins

import (
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	logicalPki "github.com/hashicorp/vault/builtin/logical/pki"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
)

func newRegistry() *registry {
	reg := &registry{
		credentialBackends: map[string]credentialBackend{},
		databasePlugins:    map[string]databasePlugin{},
		logicalBackends: map[string]logicalBackend{
			"kv":      {Factory: logicalKv.Factory},
			"pki":     {Factory: logicalPki.Factory},
			"transit": {Factory: logicalTransit.Factory},
		},
	}

	entAddExtPlugins(reg)

	return reg
}
