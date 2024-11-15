// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathStaticAccountRotateKey(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s/rotate-key", staticAccountPathPrefix, framework.GenericNameRegex("name")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "rotate",
			OperationSuffix: "static-account-key",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Name of the account.",
			},
		},
		ExistenceCheck: b.pathStaticAccountExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.pathStaticAccountRotateKey,
				ForwardPerformanceStandby:   true,
				ForwardPerformanceSecondary: true,
			},
		},
		HelpSynopsis:    pathStaticAccountRotateKeyHelpSyn,
		HelpDescription: pathStaticAccountRotateKeyHelpDesc,
	}
}

func (b *backend) pathStaticAccountRotateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}
	name := nameRaw.(string)

	b.staticAccountLock.Lock()
	defer b.staticAccountLock.Unlock()

	acct, err := b.getStaticAccount(name, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return logical.ErrorResponse("account '%s' not found", name), nil
	}

	if acct.SecretType != SecretTypeAccessToken {
		return logical.ErrorResponse("cannot rotate key for non-access-token static account"), nil
	}

	if acct.TokenGen == nil {
		return nil, fmt.Errorf("unexpected invalid account has no TokenGen")
	}

	scopes := acct.TokenGen.Scopes
	oldTokenGen := acct.TokenGen
	oldWalId, err := b.addWalRoleSetServiceAccountKey(ctx, req, acct.Name, &acct.ServiceAccountId, oldTokenGen.KeyName)
	if err != nil {
		return nil, err
	}

	// Add WALs for new TokenGen - since we don't have a key ID yet, give an empty key name so WAL
	// will know to just clear keys that aren't being used. This also covers up cleaning up
	// the old token generator, so we don't add a separate WAL for that.
	newWalId, err := b.addWalRoleSetServiceAccountKey(ctx, req, acct.Name, &acct.ServiceAccountId, "")
	if err != nil {
		return nil, err
	}

	newTokenGen, err := b.createNewTokenGen(ctx, req, acct.ResourceName(), scopes)
	if err != nil {
		return nil, err
	}

	// Edit roleset with new key and save to storage.
	acct.TokenGen = newTokenGen
	if err := acct.save(ctx, req.Storage); err != nil {
		return nil, err
	}

	// Try deleting the old key.
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, err
	}

	b.tryDeleteWALs(ctx, req.Storage, newWalId)

	if oldTokenGen != nil {
		if err := b.deleteTokenGenKey(ctx, iamAdmin, oldTokenGen); err != nil {
			return &logical.Response{
				Warnings: []string{
					fmt.Sprintf("saved static account with new token generator service account key but failed to delete old key (covered by WAL): %v", err),
				},
			}, nil
		}
		b.tryDeleteWALs(ctx, req.Storage, oldWalId)
	}
	return nil, nil
}

const pathStaticAccountRotateKeyHelpSyn = `Rotate the key used to generate access tokens for a static account`
const pathStaticAccountRotateKeyHelpDesc = `
This path allows you to manually rotate the service account key
created by Vault for a static account that generates access tokens secrets.
This path only applies to static accounts that generate access tokens. 
It will not delete the associated service account or change bindings.

Note that this will not invalidate access tokens created with the old key.
The only way to do so is to delete the service account.
`
