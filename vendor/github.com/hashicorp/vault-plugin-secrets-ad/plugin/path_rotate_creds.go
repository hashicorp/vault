// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package plugin

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	rotateRolePath = "rotate-role/"
)

func (b *backend) pathRotateCredentials() *framework.Path {
	return &framework.Path{
		Pattern: rotateRolePath + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the static role",
			},
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathRotateCredentialsUpdate,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},

		HelpSynopsis:    pathRotateCredentialsUpdateHelpSyn,
		HelpDescription: pathRotateCredentialsUpdateHelpDesc,
	}
}

func (b *backend) pathRotateCredentialsUpdate(ctx context.Context, req *logical.Request, fieldData *framework.FieldData) (*logical.Response, error) {
	cred := make(map[string]interface{})

	config, err := readConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("the config is currently unset")
	}

	roleName := fieldData.Get("name").(string)

	b.credLock.Lock()
	defer b.credLock.Unlock()

	role, err := b.readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, fmt.Errorf("role %s does not exist", roleName)
	}

	if !role.LastVaultRotation.IsZero() {
		credIfc, found := b.credCache.Get(roleName)

		if found {
			b.Logger().Debug("checking cached credential")
			cred = credIfc.(map[string]interface{})
		} else {
			b.Logger().Debug("checking stored credential")
			entry, err := req.Storage.Get(ctx, storageKey+"/"+roleName)
			if err != nil {
				return nil, err
			}

			// If the creds aren't in storage, but roles are and we've created creds before,
			// this is an unexpected state and something has gone wrong.
			// Let's be explicit and error about this.
			if entry == nil {
				b.Logger().Warn("should have the creds for %+v but they're not found", role)
			} else {
				if err := entry.DecodeJSON(&cred); err != nil {
					return nil, err
				}
				b.credCache.SetDefault(roleName, cred)
			}
		}
	}

	_, err = b.generateAndReturnCreds(ctx, config, req.Storage, roleName, role, cred)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

const pathRotateCredentialsUpdateHelpSyn = `
Request to rotate the role's credentials.
`

const pathRotateCredentialsUpdateHelpDesc = `
This path attempts to rotate the role's credentials. 
`
