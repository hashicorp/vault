// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package builder

import (
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type BackendBuilder[CC, C, R any] struct {
	Name               string
	Version            string
	Type               logical.BackendType
	BackendHelpMessage string

	WALRollback       framework.WALRollbackFunc
	WALRollbackMinAge time.Duration

	ClientConfig *ClientConfig[CC, C, R]
	Role         *Role[R, C]
}

type ClientConfig[CC, C, R any] struct {
	NewClientFunc func(*CC) (*C, error)
	Fields        map[string]*framework.FieldSchema
	ValidateFunc  func(*CC) error
}

type Role[R, C any] struct {
	Fields       map[string]*framework.FieldSchema
	Secret       *Secret[R, C]
	ValidateFunc func(*R) error
}

type Secret[R, C any] struct {
	Type            string
	Fields          map[string]*framework.FieldSchema
	FetchSecretFunc func(req *logical.Request, d *framework.FieldData, client *C, role *R, resp *framework.Secret) (*logical.Response, error)
	RevokeFunc      func(req *logical.Request, d *framework.FieldData, client *C, role *R) (*logical.Response, error)
	RenewFunc       func(req *logical.Request, d *framework.FieldData, client *C, role *R) (*logical.Response, error)
}

func (bb *BackendBuilder[CC, C, R]) build() (*GenericBackend[CC, C, R], error) {
	gb := &GenericBackend[CC, C, R]{}

	gb.newClient = bb.ClientConfig.NewClientFunc
	gb.validateConfig = bb.ClientConfig.ValidateFunc
	gb.validateRole = bb.Role.ValidateFunc

	configPath := gb.pathConfig(bb.ClientConfig)
	rolePaths := gb.pathRole(bb.Role)

	gb.Backend = &framework.Backend{
		BackendType: bb.Type,
		Help:        strings.TrimSpace(bb.BackendHelpMessage),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				configPath.Pattern,
				"role/*",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				configPath,
			},
			rolePaths,
		),
		Secrets:           []*framework.Secret{},
		Invalidate:        gb.invalidate,
		WALRollback:       bb.WALRollback,
		WALRollbackMinAge: bb.WALRollbackMinAge,
		RunningVersion:    bb.Version,
	}
	return gb, nil
}
