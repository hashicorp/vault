// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getImpersonatedAccount(name string, ctx context.Context, s logical.Storage) (*ImpersonatedAccount, error) {
	b.Logger().Debug("getting impersonated account from storage", "impersonated_account_name", name)
	entry, err := s.Get(ctx, fmt.Sprintf("%s/%s", impersonatedAccountStoragePrefix, name))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	a := &ImpersonatedAccount{}
	if err := entry.DecodeJSON(a); err != nil {
		return nil, err
	}
	return a, nil
}

type ImpersonatedAccount struct {
	Name string
	gcputil.ServiceAccountId

	TokenScopes []string
	Ttl         int
}

func (a *ImpersonatedAccount) validate() error {
	err := &multierror.Error{}
	if a.Name == "" {
		err = multierror.Append(err, errors.New("impersonated account name is empty"))
	}

	if a.EmailOrId == "" {
		err = multierror.Append(err, fmt.Errorf("impersonated account must have service account email"))
	}

	if len(a.TokenScopes) == 0 {
		err = multierror.Append(err, fmt.Errorf("access token impersonated account should have defined scopes"))
	}

	return err.ErrorOrNil()
}

// parseOkInputServiceAccountEmail checks that when creating a static account, a service account
// email is provided. A service account email can be provided while updating the static account
// but it must be the same as the one in the static account and cannot be updated.
func (a *ImpersonatedAccount) parseOkInputServiceAccountEmail(d *framework.FieldData) (warnings []string, err error) {
	email := d.Get("service_account_email").(string)
	if email == "" && a.EmailOrId == "" {
		return nil, fmt.Errorf("email is required")
	}
	if a.EmailOrId != "" && email != "" && a.EmailOrId != email {
		return nil, fmt.Errorf("cannot update email")
	}

	a.EmailOrId = email
	return nil, nil
}

func (a *ImpersonatedAccount) parseOkInputTokenScopes(d *framework.FieldData) (warnings []string, err error) {
	v, ok := d.GetOk("token_scopes")
	if ok {
		scopes, castOk := v.([]string)
		if !castOk {
			return nil, fmt.Errorf("scopes unexpected type %T, expected []string", v)
		}
		a.TokenScopes = scopes
	}

	if len(a.TokenScopes) == 0 {
		return nil, fmt.Errorf("non-empty token_scopes must be provided for generating secrets")
	}

	return
}

func (a *ImpersonatedAccount) save(ctx context.Context, s logical.Storage) error {
	if err := a.validate(); err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", impersonatedAccountStoragePrefix, a.Name), a)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func (b *backend) tryDeleteImpersonatedAccountResources(ctx context.Context, req *logical.Request, boundResources *gcpAccountResources, walIds []string) []string {
	return b.tryDeleteGcpAccountResources(ctx, req, boundResources, flagMustKeepServiceAccount, walIds)
}

func (b *backend) createImpersonatedAccount(ctx context.Context, req *logical.Request, input *ImpersonatedAccount) (err error) {
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return err
	}

	gcpAcct, err := b.getServiceAccount(iamAdmin, &gcputil.ServiceAccountId{
		Project:   gcpServiceAccountInferredProject,
		EmailOrId: input.EmailOrId,
	})
	if err != nil {
		if isGoogleAccountNotFoundErr(err) {
			return fmt.Errorf("unable to create impersonated account, service account %q should exist", input.EmailOrId)
		}
		return fmt.Errorf("unable to create impersonated account, could not confirm service account %q exists: %w", input.EmailOrId, err)
	}

	acctId := gcputil.ServiceAccountId{
		Project:   gcpAcct.ProjectId,
		EmailOrId: gcpAcct.Email,
	}

	// Construct new impersonated account
	a := &ImpersonatedAccount{
		Name:             input.Name,
		ServiceAccountId: acctId,
		TokenScopes:      input.TokenScopes,
		Ttl:              input.Ttl,
	}

	// Save to storage.
	if err := a.save(ctx, req.Storage); err != nil {
		return err
	}

	return err
}

func (b *backend) updateImpersonatedAccount(ctx context.Context, req *logical.Request, a *ImpersonatedAccount, updateInput *ImpersonatedAccount) (warnings []string, err error) {
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, err
	}

	_, err = b.getServiceAccount(iamAdmin, &a.ServiceAccountId)
	if err != nil {
		if isGoogleAccountNotFoundErr(err) {
			return nil, fmt.Errorf("unable to update impersonated account, could not find service account %q", a.ResourceName())
		}
		return nil, fmt.Errorf("unable to create impersonated account, could not confirm service account %q exists: %w", a.ResourceName(), err)
	}

	madeChange := false
	if !strutil.EquivalentSlices(updateInput.TokenScopes, a.TokenScopes) {
		b.Logger().Debug("detected scopes change, updating scopes for impersonated account")
		a.TokenScopes = updateInput.TokenScopes
		madeChange = true
	}

	if updateInput.Ttl != a.Ttl {
		b.Logger().Debug("detected ttl change, updating ttl for impersonated account")
		a.Ttl = updateInput.Ttl
		madeChange = true
	}

	if !madeChange {
		return nil, nil
	}

	if err := a.save(ctx, req.Storage); err != nil {
		return nil, err
	}

	return
}
