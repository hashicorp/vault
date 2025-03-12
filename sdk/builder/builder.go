// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package builder

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type BackendBuilder[CC any, C any] struct {
	Name               string
	Version            string
	Type               logical.BackendType
	BackendHelpMessage string

	WALRollback       framework.WALRollbackFunc
	WALRollbackMinAge time.Duration

	ClientConfig *ClientConfig[CC, C]
}

type ClientConfig[CC, C any] struct {
	NewClientFunc func(*CC) (*C, error)
	Fields        map[string]*framework.FieldSchema
}

type Role[R, C any] struct {
	Fields  map[string]*framework.FieldSchema
	Secrets Secret[R, C]
}

type Secret[R, C any] struct {
	Type         string
	Fields       map[string]*framework.FieldSchema
	CreateSecret func(*R, *C) *logical.Response
	RevokeFunc   func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error)
	RenewFunc    func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error)
}

func (bb *BackendBuilder[CC, C]) build() (*GenericBackend[CC, C], error) {
	gb := &GenericBackend[CC, C]{}

	gb.newClient = bb.ClientConfig.NewClientFunc

	configPath, err := gb.pathConfig(bb.ClientConfig.Fields)
	if err != nil {
		return nil, err
	}

	gb.Backend = &framework.Backend{
		BackendType: bb.Type,
		Help:        strings.TrimSpace(bb.BackendHelpMessage),
		PathsSpecial: &logical.Paths{
			LocalStorage: []string{},
			SealWrapStorage: []string{
				"config",
				//"role/*",
			},
		},
		Paths: framework.PathAppend(
			[]*framework.Path{
				configPath,
			},
		),
		Secrets:           []*framework.Secret{},
		Invalidate:        gb.invalidate,
		WALRollback:       bb.WALRollback,
		WALRollbackMinAge: bb.WALRollbackMinAge,
		RunningVersion:    bb.Version,
	}
	return gb, nil
}
