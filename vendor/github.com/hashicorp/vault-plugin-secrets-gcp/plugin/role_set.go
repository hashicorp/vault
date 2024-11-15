// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iam/v1"
)

var serviceAccountRegex = regexp.MustCompile("[^a-zA-Z0-9-]+")

const (
	serviceAccountMaxLen             = 30
	serviceAccountEmailTemplate      = "%s@%s.iam.gserviceaccount.com"
	serviceAccountDisplayNameHashLen = 8
	serviceAccountDisplayNameMaxLen  = 100
	serviceAccountDisplayNameTmpl    = "Service account for Vault secrets backend role set %s"
)

type RoleSet struct {
	Name       string
	SecretType string

	RawBindings string
	Bindings    ResourceBindings

	AccountId *gcputil.ServiceAccountId
	TokenGen  *TokenGenerator
}

// boundResources is a helper method to get the bound gcpAccountResources
func (rs *RoleSet) boundResources() *gcpAccountResources {
	if rs.AccountId == nil {
		return nil
	}
	return &gcpAccountResources{
		accountId: *rs.AccountId,
		bindings:  rs.Bindings,
		tokenGen:  rs.TokenGen,
	}
}

// validate checks whether a RoleSet has been populated properly before saving
func (rs *RoleSet) validate() error {
	var err *multierror.Error
	if rs.Name == "" {
		err = multierror.Append(err, errors.New("role set name is empty"))
	}

	if rs.SecretType == "" {
		err = multierror.Append(err, errors.New("role set secret type is empty"))
	}

	if rs.AccountId == nil {
		err = multierror.Append(err, fmt.Errorf("role set should have account associated"))
	}

	if len(rs.Bindings) == 0 {
		err = multierror.Append(err, fmt.Errorf("role set bindings cannot be empty"))
	}

	if len(rs.RawBindings) == 0 {
		err = multierror.Append(err, fmt.Errorf("role set raw bindings cannot be empty string"))
	}

	switch rs.SecretType {
	case SecretTypeAccessToken:
		if rs.TokenGen == nil {
			err = multierror.Append(err, fmt.Errorf("access token role set should have initialized token generator"))
		} else if len(rs.TokenGen.Scopes) == 0 {
			err = multierror.Append(err, fmt.Errorf("access token role set should have defined scopes"))
		}
	case SecretTypeKey:
		break
	default:
		err = multierror.Append(err, fmt.Errorf("unknown secret type: %s", rs.SecretType))
	}
	return err.ErrorOrNil()
}

// save saves a roleset to storage
func (rs *RoleSet) save(ctx context.Context, s logical.Storage) error {
	if err := rs.validate(); err != nil {
		return err
	}

	entry, err := logical.StorageEntryJSON(fmt.Sprintf("%s/%s", rolesetStoragePrefix, rs.Name), rs)
	if err != nil {
		return err
	}

	return s.Put(ctx, entry)
}

func (rs *RoleSet) bindingHash() string {
	return getStringHash(rs.RawBindings)
}

// getServiceAccount fetches an service account from the GCP IAM admin API.
func (b *backend) getServiceAccount(iamAdmin *iam.Service, accountId *gcputil.ServiceAccountId) (*iam.ServiceAccount, error) {
	if accountId == nil {
		return nil, fmt.Errorf("cannot fetch nil service account")
	}

	account, err := iamAdmin.Projects.ServiceAccounts.Get(accountId.ResourceName()).Do()
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("could not find service account %q: {{err}}", accountId.ResourceName()), err)
	}
	return account, nil
}

// saveRoleSetWithNewAccount rotates the role set service account. This includes creating a new service account with
// a new name and deleting the old service account, updating keys or bindings as required.
func (b *backend) saveRoleSetWithNewAccount(ctx context.Context, req *logical.Request, rs *RoleSet, project string, newBinds ResourceBindings, scopes []string) (warnings []string, err error) {
	b.Logger().Debug("updating roleset with new account")

	oldResources := rs.boundResources()

	// Generate name for new account
	newSaName := generateAccountNameForRoleSet(rs.Name)

	// Construct IDs for new resources.
	// The actual GCP resources are not created yet, but we need the IDs to create WAL entries.
	newResources := &gcpAccountResources{
		accountId: gcputil.ServiceAccountId{
			Project:   project,
			EmailOrId: emailForServiceAccountName(project, newSaName),
		},
		bindings: newBinds,
	}
	if len(scopes) > 0 {
		newResources.tokenGen = &TokenGenerator{Scopes: scopes}
	}

	// Add WALs for both old and new resources.
	// WAL callback checks whether resources are still being used by roleset so
	// there is no harm in adding WALs early, or adding WALs for resources that
	// will eventually get cleaned up.
	b.Logger().Debug("adding WALs for old roleset resources")
	oldWalIds, err := b.addWalsForRoleSetResources(ctx, req, rs.Name, oldResources)
	if err != nil {
		return nil, err
	}

	b.Logger().Debug("adding WALs for new roleset resources")
	newWalIds, err := b.addWalsForRoleSetResources(ctx, req, rs.Name, newResources)
	if err != nil {
		return nil, err
	}

	// Created new RoleSet resources
	// Create new service account
	sa, err := b.createServiceAccount(ctx, req, newResources.accountId.Project, newSaName, fmt.Sprintf("role set %s", rs.Name))
	if err != nil {
		return nil, err
	}

	// Create new IAM bindings.
	if err := b.createIamBindings(ctx, req, sa.Email, newResources.bindings); err != nil {
		return nil, err
	}

	// Create new token gen if a stubbed tokenGenerator (with scopes) is given.
	if newResources.tokenGen != nil && len(newResources.tokenGen.Scopes) > 0 {
		tokenGen, err := b.createNewTokenGen(ctx, req, sa.Name, newResources.tokenGen.Scopes)
		if err != nil {
			return nil, err
		}
		newResources.tokenGen = tokenGen
	}

	// Edit roleset with new resources and save to storage.
	rs.AccountId = &newResources.accountId
	rs.Bindings = newResources.bindings
	rs.TokenGen = newResources.tokenGen
	if err := rs.save(ctx, req.Storage); err != nil {
		return nil, err
	}

	// We successfully saved the new roleset with new resources, so try cleaning up WALs
	// that would rollback the roleset resources (will no-op if still in use by roleset)
	b.tryDeleteWALs(ctx, req.Storage, newWalIds...)

	return b.tryDeleteRoleSetResources(ctx, req, oldResources, oldWalIds), nil
}

// saveRoleSetWithNewTokenKey rotates the role set access_token key and saves it to storage.
func (b *backend) saveRoleSetWithNewTokenKey(ctx context.Context, req *logical.Request, rs *RoleSet, scopes []string) (warning string, err error) {
	if rs.SecretType != SecretTypeAccessToken {
		return "", fmt.Errorf("a key is not saved or used for non-access-token role set '%s'", rs.Name)
	}
	if rs.AccountId == nil {
		return "", fmt.Errorf("unable to save roleset with new key - account ID was nil")
	}

	b.Logger().Debug("updating roleset with new account key")

	var oldTokenGen *TokenGenerator
	var oldWalId string
	if rs.TokenGen != nil {
		scopes = rs.TokenGen.Scopes
		oldTokenGen = rs.TokenGen
		oldWalId, err = b.addWalRoleSetServiceAccountKey(ctx, req, rs.Name, rs.AccountId, oldTokenGen.KeyName)
		if err != nil {
			return "", err
		}
	}

	// Add WALs for new TokenGen - since we don't have a key ID yet, give an empty key name so WAL
	// will know to just clear keys that aren't being used. This also covers up cleaning up
	// the old token generator, so we don't add a separate WAL for that.
	newWalId, err := b.addWalRoleSetServiceAccountKey(ctx, req, rs.Name, rs.AccountId, "")
	if err != nil {
		return "", err
	}

	newTokenGen, err := b.createNewTokenGen(ctx, req, rs.AccountId.ResourceName(), scopes)
	if err != nil {
		return "", err
	}

	// Edit roleset with new key and save to storage.
	rs.TokenGen = newTokenGen
	if err := rs.save(ctx, req.Storage); err != nil {
		return "", err
	}

	// Try deleting the old key.
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return "", err
	}

	b.tryDeleteWALs(ctx, req.Storage, newWalId)
	if oldTokenGen != nil {
		if err := b.deleteTokenGenKey(ctx, iamAdmin, oldTokenGen); err != nil {
			return errwrap.Wrapf("roleset update succeeded but got error while trying to delete old key - will be cleaned up later by WAL: {{err}}", err).Error(), nil
		}
		b.tryDeleteWALs(ctx, req.Storage, oldWalId)
	}
	return "", nil
}

// addWalsForRoleSetResources creates WALs to clean up a roleset's service account, bindings, and a key if needed.
func (b *backend) addWalsForRoleSetResources(ctx context.Context, req *logical.Request, rolesetName string, boundResources *gcpAccountResources) (walIds []string, err error) {
	if boundResources == nil {
		b.Logger().Debug("skip WALs for nil roleset resources")
		return nil, nil
	}

	walIds = make([]string, 0, len(boundResources.bindings)+2)
	walId, err := framework.PutWAL(ctx, req.Storage, walTypeAccount, &walAccount{
		RoleSet: rolesetName,
		Id:      boundResources.accountId,
	})
	if err != nil {
		return walIds, errwrap.Wrapf("unable to create WAL entry to clean up service account: {{err}}", err)
	}
	walIds = append(walIds, walId)

	for resource, roles := range boundResources.bindings {
		walId, err := framework.PutWAL(ctx, req.Storage, walTypeIamPolicy, &walIamPolicy{
			RoleSet:   rolesetName,
			AccountId: boundResources.accountId,
			Resource:  resource,
			Roles:     roles.ToSlice(),
		})
		if err != nil {
			return walIds, errwrap.Wrapf("unable to create WAL entry to clean up service account bindings: {{err}}", err)
		}
		walIds = append(walIds, walId)
	}

	if boundResources.tokenGen != nil {
		walId, err := b.addWalRoleSetServiceAccountKey(ctx, req, rolesetName, &boundResources.accountId, boundResources.tokenGen.KeyName)
		if err != nil {
			return nil, err
		}
		walIds = append(walIds, walId)
	}
	return walIds, nil
}

// addWalRoleSetServiceAccountKey creates WAL to clean up a service account key (for access tokens) if needed.
func (b *backend) addWalRoleSetServiceAccountKey(ctx context.Context, req *logical.Request, roleset string, accountId *gcputil.ServiceAccountId, keyName string) (string, error) {
	if accountId == nil {
		return "", fmt.Errorf("given nil account ID for WAL for roleset service account key")
	}

	b.Logger().Debug("add WAL for service account key", "account", accountId.ResourceName(), "keyName", keyName)

	walId, err := framework.PutWAL(ctx, req.Storage, walTypeAccount, &walAccountKey{
		RoleSet:            roleset,
		ServiceAccountName: accountId.ResourceName(),
		KeyName:            keyName,
	})
	if err != nil {
		return "", errwrap.Wrapf("unable to create WAL entry to clean up service account key: {{err}}", err)
	}
	return walId, nil
}

// tryDeleteRoleSetResources tries to delete GCP resources previously managed by a roleset.
// This assumes that deletion of these resources will already be guaranteed by WAL rollbacks (referred to by the walIds)
// and will return errors as a list of warnings instead.
func (b *backend) tryDeleteRoleSetResources(ctx context.Context, req *logical.Request, boundResources *gcpAccountResources, walIds []string) []string {
	return b.tryDeleteGcpAccountResources(ctx, req, boundResources, flagCanDeleteServiceAccount, walIds)
}

// generateAccountNameForRoleSet returns a new random name for a Vault service account based off roleset name and time.
// Note this is the name rather than the full email (i.e. string before @)
//
// As an example, for roleset "my-role" this returns `vaultmy-role-1234613`
func generateAccountNameForRoleSet(rsName string) (name string) {
	// Sanitize role name
	rsName = serviceAccountRegex.ReplaceAllString(rsName, "-")

	intSuffix := fmt.Sprintf("%d", time.Now().Unix())
	fullName := fmt.Sprintf("vault%s-%s", rsName, intSuffix)
	name = fullName
	if len(fullName) > serviceAccountMaxLen {
		toTrunc := len(fullName) - serviceAccountMaxLen
		name = fmt.Sprintf("vault%s-%s", rsName[:len(rsName)-toTrunc], intSuffix)
	}
	return name
}
