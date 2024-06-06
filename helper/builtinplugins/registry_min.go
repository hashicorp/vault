// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package builtinplugins

import (
	credJWT "github.com/hashicorp/vault-plugin-auth-jwt"
	logicalKv "github.com/hashicorp/vault-plugin-secrets-kv"
	credAppRole "github.com/hashicorp/vault/builtin/credential/approle"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	logicalPki "github.com/hashicorp/vault/builtin/logical/pki"
	logicalSsh "github.com/hashicorp/vault/builtin/logical/ssh"
	logicalTransit "github.com/hashicorp/vault/builtin/logical/transit"
)

func newRegistry() *registry {
	reg := &registry{
		credentialBackends: map[string]credentialBackend{
			"approle":  {Factory: credAppRole.Factory},
			"cert":     {Factory: credCert.Factory},
			"jwt":      {Factory: credJWT.Factory},
			"oidc":     {Factory: credJWT.Factory},
			"userpass": {Factory: credUserpass.Factory},
		},
		databasePlugins: map[string]databasePlugin{},
		logicalBackends: map[string]logicalBackend{
			"kv":      {Factory: logicalKv.Factory},
			"pki":     {Factory: logicalPki.Factory},
			"ssh":     {Factory: logicalSsh.Factory},
			"transit": {Factory: logicalTransit.Factory},
		},
	}

	entAddExtPlugins(reg)

	return reg
}
