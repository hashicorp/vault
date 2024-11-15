// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mongodbatlas

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"go.mongodb.org/atlas/mongodbatlas"
)

// operationPrefixMongoDBAtlas is used as a prefix for OpenAPI operation id's.
const operationPrefixMongoDBAtlas = "mongo-db-atlas"

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	b := NewBackend(conf.System)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

func NewBackend(system logical.SystemView) *Backend {
	var b Backend
	b.Backend = &framework.Backend{
		Help: strings.TrimSpace(backendHelp),

		PathsSpecial: &logical.Paths{
			LocalStorage: []string{
				framework.WALPrefix,
			},
			SealWrapStorage: []string{
				"config",
			},
		},

		Paths: []*framework.Path{
			b.pathRolesList(),
			b.pathRoles(),
			b.pathConfig(),
			b.pathCredentials(),
		},

		Secrets: []*framework.Secret{
			b.programmaticAPIKeys(),
		},

		WALRollback:       b.pathProgrammaticAPIKeyRollback,
		WALRollbackMinAge: minUserRollbackAge,
		BackendType:       logical.TypeLogical,
	}
	b.system = system
	return &b
}

type Backend struct {
	*framework.Backend

	credentialMutex sync.RWMutex
	clientMutex     sync.RWMutex

	client *mongodbatlas.Client

	system logical.SystemView
}

const backendHelp = `
The MongoDB Atlas backend dynamically generates API keys for a set of 
Organization or Project roles. The API keys have a configurable lease 
set and are automatically revoked at the end of the lease.

After mounting this backend, the Public and Private keys to generate 
API keys must be configured with the "config" path and roles must be 
written  using the "roles/" endpoints before any API keys can be generated.

`
const minUserRollbackAge = 5 * time.Minute
