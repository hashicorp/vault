// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build minimal

package command

import (
	credOIDC "github.com/hashicorp/vault-plugin-auth-jwt"
	credCert "github.com/hashicorp/vault/builtin/credential/cert"
	credToken "github.com/hashicorp/vault/builtin/credential/token"
	credUserpass "github.com/hashicorp/vault/builtin/credential/userpass"
	_ "github.com/hashicorp/vault/helper/builtinplugins"
	physRaft "github.com/hashicorp/vault/physical/raft"
	"github.com/hashicorp/vault/sdk/physical"
	physInmem "github.com/hashicorp/vault/sdk/physical/inmem"
)

var (
	physicalBackends = map[string]physical.Factory{
		"inmem_ha":               physInmem.NewInmemHA,
		"inmem_transactional_ha": physInmem.NewTransactionalInmemHA,
		"inmem_transactional":    physInmem.NewTransactionalInmem,
		"inmem":                  physInmem.NewInmem,
		"raft":                   physRaft.NewRaftBackend,
	}

	loginHandlers = map[string]LoginHandler{
		"cert":  &credCert.CLIHandler{},
		"oidc":  &credOIDC.CLIHandler{},
		"token": &credToken.CLIHandler{},
		"userpass": &credUserpass.CLIHandler{
			DefaultMount: "userpass",
		},
	}
)
