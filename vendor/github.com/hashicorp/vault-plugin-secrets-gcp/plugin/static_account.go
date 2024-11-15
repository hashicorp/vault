// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) getStaticAccount(name string, ctx context.Context, s logical.Storage) (*StaticAccount, error) {
	b.Logger().Debug("getting static account from storage", "static_account_name", name)
	entry, err := s.Get(ctx, fmt.Sprintf("%s/%s", staticAccountStoragePrefix, name))
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	a := &StaticAccount{}
	if err := entry.DecodeJSON(a); err != nil {
		return nil, err
	}
	return a, nil
}

type StaticAccount struct {
	Name        string
	SecretType  string
	RawBindings string
	Bindings    ResourceBindings
	gcputil.ServiceAccountId

	TokenGen *TokenGenerator
}

func (a *StaticAccount) boundResources() *gcpAccountResources {
	return &gcpAccountResources{
		accountId: a.ServiceAccountId,
		bindings:  a.Bindings,
		tokenGen:  a.TokenGen,
	}
}

func (a *StaticAccount) bindingHash() string {
	return getStringHash(a.RawBindings)
}

func (a *StaticAccount) validate() error {
	err := &multierror.Error{}
	if a.Name == "" {
		err = multierror.Append(err, errors.New("static account name is empty"))
	}

	if a.SecretType == "" {
		err = multierror.Append(err, errors.New("static account secret type is empty"))
	}

	if a.EmailOrId == "" {
		err = multierror.Append(err, fmt.Errorf("static account must have service account email"))
	}

	switch a.SecretType {
	case SecretTypeAccessToken:
		if a.TokenGen == nil {
			err = multierror.Append(err, fmt.Errorf("access token static account should have initialized token generator"))
		} else if len(a.TokenGen.Scopes) == 0 {
			err = multierror.Append(err, fmt.Errorf("access token static account should have defined scopes"))
		}
	case SecretTypeKey:
		break
	default:
		err = multierror.Append(err, fmt.Errorf("unknown secret type: %s", a.SecretType))
	}
	return err.ErrorOrNil()
}

func (a *StaticAccount) save(ctx context.Context, s logical.Storage) error {
	if err := a.validate(); err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", staticAccountStoragePrefix, a.Name), a)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func (b *backend) tryDeleteStaticAccountResources(ctx context.Context, req *logical.Request, boundResources *gcpAccountResources, walIds []string) []string {
	return b.tryDeleteGcpAccountResources(ctx, req, boundResources, flagMustKeepServiceAccount, walIds)
}

func (b *backend) createStaticAccount(ctx context.Context, req *logical.Request, input *inputParams) (err error) {
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return err
	}

	gcpAcct, err := b.getServiceAccount(iamAdmin, &gcputil.ServiceAccountId{
		Project:   gcpServiceAccountInferredProject,
		EmailOrId: input.serviceAccountEmail,
	})
	if err != nil {
		if isGoogleAccountNotFoundErr(err) {
			return fmt.Errorf("unable to create static account, service account %q should exist", input.serviceAccountEmail)
		}
		return errwrap.Wrapf(fmt.Sprintf("unable to create static account, could not confirm service account %q exists: {{err}}", input.serviceAccountEmail), err)
	}

	acctId := gcputil.ServiceAccountId{
		Project:   gcpAcct.ProjectId,
		EmailOrId: gcpAcct.Email,
	}

	// Construct gcpAccountResources references. Note bindings/key are yet to be created.
	newResources := &gcpAccountResources{
		accountId: acctId,
		bindings:  input.bindings,
	}
	if input.secretType == SecretTypeAccessToken {
		newResources.tokenGen = &TokenGenerator{
			Scopes: input.scopes,
		}
	}

	// add WALs for static account resources
	newWalIds, err := b.addWalsForStaticAccountResources(ctx, req, input.name, newResources)
	if err != nil {
		b.tryDeleteWALs(ctx, req.Storage, newWalIds...)
		return err
	}

	// Create new IAM bindings.
	if err := b.createIamBindings(ctx, req, gcpAcct.Email, newResources.bindings); err != nil {
		return err
	}

	// Create new token gen if a stubbed tokenGenerator (with scopes) is given.
	if newResources.tokenGen != nil && len(newResources.tokenGen.Scopes) > 0 {
		tokenGen, err := b.createNewTokenGen(ctx, req, gcpAcct.Name, newResources.tokenGen.Scopes)
		if err != nil {
			return err
		}
		newResources.tokenGen = tokenGen
	}

	// Construct new static account
	a := &StaticAccount{
		Name:             input.name,
		SecretType:       input.secretType,
		RawBindings:      input.rawBindings,
		Bindings:         input.bindings,
		ServiceAccountId: acctId,
		TokenGen:         newResources.tokenGen,
	}

	// Save to storage.
	if err := a.save(ctx, req.Storage); err != nil {
		return err
	}

	// We successfully saved the new static account so try cleaning up WALs
	// that would rollback the GCP resources (will no-op if still in use)
	b.tryDeleteWALs(ctx, req.Storage, newWalIds...)
	return err
}

func (b *backend) updateStaticAccount(ctx context.Context, req *logical.Request, a *StaticAccount, updateInput *inputParams) (warnings []string, err error) {
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, err
	}

	_, err = b.getServiceAccount(iamAdmin, &a.ServiceAccountId)
	if err != nil {
		if isGoogleAccountNotFoundErr(err) {
			return nil, fmt.Errorf("unable to update static account, could not find service account %q", a.ResourceName())
		}
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to create static account, could not confirm service account %q exists: {{err}}", a.ResourceName()), err)
	}

	var walIds []string
	madeChange := false

	if updateInput.hasBindings && a.bindingHash() != updateInput.rawBindings {
		b.Logger().Debug("detected bindings change, updating bindings for static account")
		newBindings := updateInput.bindings

		bindingWals, err := b.updateBindingsForStaticAccount(ctx, req, a, newBindings)
		if err != nil {
			return nil, err
		}
		walIds = append(walIds, bindingWals...)

		a.RawBindings = updateInput.rawBindings
		a.Bindings = newBindings
		madeChange = true
	}

	if a.SecretType == "access_token" {
		if a.TokenGen == nil {
			return nil, fmt.Errorf("unexpected invalid access_token static account has no TokenGen")
		}
		if !strutil.EquivalentSlices(updateInput.scopes, a.TokenGen.Scopes) {
			b.Logger().Debug("detected scopes change, updating scopes for static account")
			a.TokenGen.Scopes = updateInput.scopes
			madeChange = true
		}
	}

	if !madeChange {
		return nil, nil
	}

	if err := a.save(ctx, req.Storage); err != nil {
		return nil, err
	}

	b.tryDeleteWALs(ctx, req.Storage, walIds...)
	return
}

func (b *backend) updateBindingsForStaticAccount(ctx context.Context, req *logical.Request, a *StaticAccount, newBindings ResourceBindings) ([]string, error) {
	oldBindings := a.Bindings
	bindsToAdd := newBindings.sub(oldBindings)
	bindsToRemove := oldBindings.sub(newBindings)

	b.Logger().Debug("updating bindings for static account")
	walIds, err := b.addWalsForStaticAccountBindings(ctx, req, a, bindsToAdd, bindsToRemove)
	if err != nil {
		return nil, err
	}

	if err := b.createIamBindings(ctx, req, a.EmailOrId, bindsToAdd); err != nil {
		return nil, err
	}

	if err := b.removeBindings(ctx, req, a.EmailOrId, bindsToRemove); err != nil {
		return nil, err
	}

	return walIds, nil
}

// addWalsForStaticAccountResources creates WALs to clean up a roleset's service account, bindings, and a key if needed.
func (b *backend) addWalsForStaticAccountResources(ctx context.Context, req *logical.Request, staticAcctName string, boundResources *gcpAccountResources) (walIds []string, err error) {
	if boundResources == nil {
		b.Logger().Debug("skip WALs for nil GCP account resources")
		return nil, nil
	}

	walIds = make([]string, 0, len(boundResources.bindings)+1)
	for resName, roles := range boundResources.bindings {
		walId, err := framework.PutWAL(ctx, req.Storage, walTypeIamPolicyDiff, &walIamPolicyStaticAccount{
			StaticAccount: staticAcctName,
			AccountId:     boundResources.accountId,
			Resource:      resName,
			RolesAdded:    roles.ToSlice(),
		})
		if err != nil {
			return walIds, errwrap.Wrapf("unable to create WAL entry to clean up service account bindings: {{err}}", err)
		}
		walIds = append(walIds, walId)
	}

	if boundResources.tokenGen != nil {
		walId, err := b.addWalStaticAccountServiceAccountKey(ctx, req, staticAcctName, &boundResources.accountId, boundResources.tokenGen.KeyName)
		if err != nil {
			return walIds, err
		}
		walIds = append(walIds, walId)
	}
	return walIds, nil
}

func (b *backend) addWalsForStaticAccountBindings(ctx context.Context, req *logical.Request, a *StaticAccount, removed, added ResourceBindings) (walIds []string, err error) {
	walIds = make([]string, 0, len(removed)+len(added))

	// Add WALs for resources in added bindings
	for resource, rolesAdded := range added {
		walEntry := &walIamPolicyStaticAccount{
			StaticAccount: a.Name,
			AccountId:     a.ServiceAccountId,
			Resource:      resource,
			RolesAdded:    rolesAdded.ToSlice(),
		}
		if rolesRemoved, ok := removed[resource]; ok {
			walEntry.RolesRemoved = rolesRemoved.ToSlice()
		}
		walId, err := framework.PutWAL(ctx, req.Storage, walTypeIamPolicyDiff, walEntry)
		if err != nil {
			return walIds, errwrap.Wrapf("unable to create WAL entry to clean up service account bindings: {{err}}", err)
		}
		walIds = append(walIds, walId)
	}

	// Add WALs for resources in removed bindings not in added bindings.
	for resource, rolesRemoved := range removed {
		if _, ok := added[resource]; ok {
			continue
		}
		walEntry := &walIamPolicyStaticAccount{
			StaticAccount: a.Name,
			AccountId:     a.ServiceAccountId,
			Resource:      resource,
			RolesRemoved:  rolesRemoved.ToSlice(),
		}

		walId, err := framework.PutWAL(ctx, req.Storage, walTypeIamPolicyDiff, walEntry)
		if err != nil {
			return walIds, errwrap.Wrapf("unable to create WAL entry to clean up service account bindings: {{err}}", err)
		}
		walIds = append(walIds, walId)
	}

	return walIds, nil
}

// addWalStaticAccountServiceAccountKey creates WAL to clean up a static account's service account key (for access tokens) if needed.
func (b *backend) addWalStaticAccountServiceAccountKey(ctx context.Context, req *logical.Request, acct string, accountId *gcputil.ServiceAccountId, keyName string) (string, error) {
	if accountId == nil {
		return "", fmt.Errorf("given nil account ID for WAL for roleset service account key")
	}

	b.Logger().Debug("add WAL for service account key", "account", accountId.ResourceName(), "keyName", keyName)

	walId, err := framework.PutWAL(ctx, req.Storage, walTypeAccount, &walAccountKey{
		StaticAccount:      acct,
		ServiceAccountName: accountId.ResourceName(),
		KeyName:            keyName,
	})
	if err != nil {
		return "", errwrap.Wrapf("unable to create WAL entry to clean up service account key: {{err}}", err)
	}
	return walId, nil
}
