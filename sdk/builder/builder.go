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
	version            string
	backendType        logical.BackendType
	backendHelpMessage string

	walRollback       framework.WALRollbackFunc
	walRollbackMinAge time.Duration

	clientConfig *ClientConfig[CC, C, R]
	role         *Role[R, C]
}

type ClientConfig[CC, C, R any] struct {
	NewClientFunc   func(*CC) (*C, error)
	Fields          map[string]*framework.FieldSchema
	ValidateFunc    func(*CC) error
	HelpSynopsis    string
	HelpDescription string
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

	gb.newClient = bb.clientConfig.NewClientFunc
	gb.validateConfig = bb.clientConfig.ValidateFunc
	gb.validateRole = bb.role.ValidateFunc

	gb.role = bb.role

	configPath := gb.pathConfig(bb.clientConfig)
	rolePaths := gb.pathRole(bb.role)
	credsPath := gb.pathCredentials()

	gb.Backend = &framework.Backend{
		BackendType: bb.backendType,
		Help:        strings.TrimSpace(bb.backendHelpMessage),
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
				credsPath,
			},
			rolePaths,
		),
		Secrets: []*framework.Secret{
			gb.secret(bb.role.Secret),
		},
		Invalidate:        gb.invalidate,
		WALRollback:       bb.walRollback,
		WALRollbackMinAge: bb.walRollbackMinAge,
		RunningVersion:    bb.version,
	}
	return gb, nil
}

func (bb *BackendBuilder[CC, C, R]) WithVersion(version string) *BackendBuilder[CC, C, R] {
	bb.version = version
	return bb
}

func (bb *BackendBuilder[CC, C, R]) WithBackendType(backendType logical.BackendType) *BackendBuilder[CC, C, R] {
	bb.backendType = backendType
	return bb
}

func (bb *BackendBuilder[CC, C, R]) WithBackendHelpMessage(backendHelpMessage string) *BackendBuilder[CC, C, R] {
	bb.backendHelpMessage = backendHelpMessage
	return bb
}

func (bb *BackendBuilder[CC, C, R]) WithWalRollbackFunc(wallRollBackFunc framework.WALRollbackFunc) *BackendBuilder[CC, C, R] {
	bb.walRollback = wallRollBackFunc
	return bb
}

func (bb *BackendBuilder[CC, C, R]) WithClientConfig(clientConfig *ClientConfig[CC, C, R]) *BackendBuilder[CC, C, R] {
	bb.clientConfig = clientConfig
	return bb
}

func (bb *BackendBuilder[CC, C, R]) WithRole(role *Role[R, C]) *BackendBuilder[CC, C, R] {
	bb.role = role
	return bb
}
