// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

const (
	impersonatedAccountStoragePrefix = "impersonated-account"
	impersonatedAccountPathPrefix    = "impersonated-account"
)

func pathImpersonatedAccount(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: fmt.Sprintf("%s/%s", impersonatedAccountPathPrefix, framework.GenericNameRegex("name")),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationSuffix: "impersonated-account",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeString,
				Description: "Required. Name to refer to this impersonated account in Vault. Cannot be updated.",
			},
			"service_account_email": {
				Type:        framework.TypeString,
				Description: "Required. Email of the GCP service account to manage. Cannot be updated.",
			},
			"token_scopes": {
				Type:        framework.TypeCommaStringSlice,
				Description: "List of OAuth scopes to assign to access tokens generated under this account.",
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: "Lifetime of the token for the impersonated account.",
			},
		},
		ExistenceCheck: b.pathImpersonatedAccountExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountDelete,
			},
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountRead,
			},
			logical.CreateOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountCreate,
			},
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountUpdate,
			},
		},
		HelpSynopsis:    pathImpersonatedAccountHelpSyn,
		HelpDescription: pathImpersonatedAccountHelpDesc,
	}
}

func pathImpersonatedAccountList(b *backend) *framework.Path {
	// Paths for listing impersonated accounts
	return &framework.Path{
		Pattern: fmt.Sprintf("%ss?/?", impersonatedAccountPathPrefix),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "list",
			OperationSuffix: "impersonated-accounts|impersonated-accounts2",
		},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback: b.pathImpersonatedAccountList,
			},
		},
		HelpSynopsis:    pathListImpersonatedAccountHelpSyn,
		HelpDescription: pathListImpersonatedAccountHelpDesc,
	}
}

func (b *backend) pathImpersonatedAccountExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return false, errors.New("impersonated account name is required")
	}

	acct, err := b.getImpersonatedAccount(nameRaw.(string), ctx, req.Storage)
	if err != nil {
		return false, err
	}

	return acct != nil, nil
}

func (b *backend) pathImpersonatedAccountRead(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}

	acct, err := b.getImpersonatedAccount(nameRaw.(string), ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return nil, nil
	}

	data := map[string]interface{}{
		"service_account_project": acct.Project,
		"service_account_email":   acct.EmailOrId,
		"token_scopes":            acct.TokenScopes,
		"ttl":                     acct.Ttl,
	}

	return &logical.Response{
		Data: data,
	}, nil
}

func (b *backend) pathImpersonatedAccountDelete(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}
	name := nameRaw.(string)

	b.impersonatedAccountLock.Lock()
	defer b.impersonatedAccountLock.Unlock()

	// Delete impersonated account
	b.Logger().Debug("deleting impersonated account from storage", "name", name)
	if err := req.Storage.Delete(ctx, fmt.Sprintf("%s/%s", impersonatedAccountStoragePrefix, name)); err != nil {
		return nil, err
	}

	b.Logger().Debug("finished deleting impersonated account from storage", "name", name)
	return nil, nil
}

func (b *backend) pathImpersonatedAccountCreate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	input, warnings, err := b.parseImpersonateInformation(ImpersonatedAccount{}, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if input == nil {
		return nil, fmt.Errorf("plugin error - parse returned unexpected nil input")
	}

	b.impersonatedAccountLock.Lock()
	defer b.impersonatedAccountLock.Unlock()

	// Create and save impersonated account with new resources.
	if err := b.createImpersonatedAccount(ctx, req, input); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if len(warnings) > 0 {
		return &logical.Response{Warnings: warnings}, nil
	}
	return nil, nil
}

func (b *backend) pathImpersonatedAccountUpdate(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	nameRaw, ok := d.GetOk("name")
	if !ok {
		return logical.ErrorResponse("name is required"), nil
	}
	name := nameRaw.(string)

	b.impersonatedAccountLock.Lock()
	defer b.impersonatedAccountLock.Unlock()

	acct, err := b.getImpersonatedAccount(name, ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if acct == nil {
		return nil, fmt.Errorf("unable to find impersonated account %s to update", name)
	}

	updateInput, warnings, err := b.parseImpersonateInformation(*acct, d)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}
	if updateInput == nil {
		return nil, fmt.Errorf("plugin error - parse returned unexpected nil input")
	}

	updateWarns, err := b.updateImpersonatedAccount(ctx, req, acct, updateInput)
	if err != nil {
		return logical.ErrorResponse("unable to update: %s", err), nil
	}
	warnings = append(warnings, updateWarns...)
	if len(warnings) > 0 {
		return &logical.Response{Warnings: warnings}, nil
	}
	return nil, nil
}

func (b *backend) pathImpersonatedAccountList(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	accounts, err := req.Storage.List(ctx, fmt.Sprintf("%s/", impersonatedAccountStoragePrefix))
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(accounts), nil
}

func (b *backend) parseImpersonateInformation(prevValues ImpersonatedAccount, d *framework.FieldData) (*ImpersonatedAccount, []string, error) {
	var warnings []string

	nameRaw, ok := d.GetOk("name")
	if !ok {
		return nil, nil, fmt.Errorf("name is required")
	}
	prevValues.Name = nameRaw.(string)

	ws, err := prevValues.parseOkInputServiceAccountEmail(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	ws, err = prevValues.parseOkInputTokenScopes(d)
	if err != nil {
		return nil, nil, err
	} else if len(ws) > 0 {
		warnings = append(warnings, ws...)
	}

	ttl, ok := d.GetOk("ttl")
	if ok {
		prevValues.Ttl = ttl.(int)
	}

	return &prevValues, warnings, nil
}

const pathImpersonatedAccountHelpSyn = `Register and manage a GCP service account to generate credentials under`
const pathImpersonatedAccountHelpDesc = `
This path allows you to register an impersonated GCP service account that you want to generate secrets against.
Secrets (i.e.access tokens) are generated under this account. The account must exist at creation of impersonated
account creation.`

const pathListImpersonatedAccountHelpSyn = `List created impersonated accounts.`
const pathListImpersonatedAccountHelpDesc = `List created impersonated accounts.`
