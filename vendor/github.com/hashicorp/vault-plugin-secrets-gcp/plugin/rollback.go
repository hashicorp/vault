// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	walTypeAccount       = "account"
	walTypeAccountKey    = "account_key"
	walTypeIamPolicy     = "iam_policy"
	walTypeIamPolicyDiff = "iam_policy_diff"
)

func (b *backend) walRollback(ctx context.Context, req *logical.Request, kind string, data interface{}) error {
	switch kind {
	case walTypeAccount:
		return b.serviceAccountRollback(ctx, req, data)
	case walTypeAccountKey:
		return b.serviceAccountKeyRollback(ctx, req, data)
	case walTypeIamPolicy:
		return b.serviceAccountPolicyRollback(ctx, req, data)
	case walTypeIamPolicyDiff:
		return b.serviceAccountPolicyDiffRollback(ctx, req, data)
	default:
		return fmt.Errorf("unknown type to rollback")
	}
}

type walAccount struct {
	RoleSet string
	Id      gcputil.ServiceAccountId
}

type walAccountKey struct {
	RoleSet            string
	StaticAccount      string
	ServiceAccountName string
	KeyName            string
}

type walIamPolicy struct {
	RoleSet   string
	AccountId gcputil.ServiceAccountId
	Resource  string
	Roles     []string
}

type walIamPolicyStaticAccount struct {
	StaticAccount string
	AccountId     gcputil.ServiceAccountId
	Resource      string
	RolesAdded    []string
	RolesRemoved  []string
}

func (b *backend) serviceAccountRollback(ctx context.Context, req *logical.Request, data interface{}) error {
	b.rolesetLock.Lock()
	defer b.rolesetLock.Unlock()

	var entry walAccount
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	rs, err := getRoleSet(entry.RoleSet, ctx, req.Storage)
	if err != nil {
		return err
	}

	// If account is still being used, WAL entry was not deleted properly after a successful operation.
	// Remove WAL entry.
	if rs != nil && entry.Id.ResourceName() == rs.AccountId.ResourceName() {
		// Still being used - don't delete this service account.
		return nil
	}

	// Delete service account.
	iamC, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return err
	}

	return b.deleteServiceAccount(ctx, iamC, entry.Id)
}

func (b *backend) serviceAccountKeyRollback(ctx context.Context, req *logical.Request, data interface{}) error {
	b.rolesetLock.Lock()
	defer b.rolesetLock.Unlock()

	var entry walAccountKey
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	b.Logger().Debug("checking parent listed in WAL generates access_token secret")
	var keyInUse string

	switch {
	case entry.RoleSet != "":
		rs, err := getRoleSet(entry.RoleSet, ctx, req.Storage)
		if err != nil {
			return err
		}

		// If roleset is not nil, get key in use.
		if rs != nil {
			if rs.SecretType != SecretTypeAccessToken {
				// Remove WAL entry - we don't clean keys if roleset generates key secrets.
				return nil
			}

			if rs.TokenGen != nil {
				keyInUse = rs.TokenGen.KeyName
			}
		}
	case entry.StaticAccount != "":
		sa, err := b.getStaticAccount(entry.StaticAccount, ctx, req.Storage)
		if err != nil {
			return err
		}

		// If roleset is not nil, get key in use.
		if sa != nil {
			if sa.SecretType != SecretTypeAccessToken {
				// Remove WAL entry - we don't clean keys if roleset generates key secrets.
				return nil
			}
			if sa.TokenGen != nil {
				keyInUse = sa.TokenGen.KeyName
			}
		}
	default:
		b.Logger().Error("removing invalid walAccountKey with empty RoleSet and empty StaticAccount, may need manual cleanup: %v", entry)
		return nil
	}

	iamC, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return err
	}

	if entry.KeyName == "" {
		// If given an empty key name, this means the WAL entry was created before the key was created.
		// We list all keys and then delete any not in use by the current roleset.
		keys, err := iamC.Projects.ServiceAccounts.Keys.List(entry.ServiceAccountName).KeyTypes("USER_MANAGED").Do()
		if err != nil {
			// If service account already deleted, no need to clean up keys.
			if isGoogleAccountNotFoundErr(err) {
				return nil
			}
			return err
		}

		for _, k := range keys.Keys {
			// Skip deleting keys still in use (empty keyInUse means no key is in use)
			if k.Name == keyInUse {
				continue
			}

			_, err = iamC.Projects.ServiceAccounts.Keys.Delete(entry.KeyName).Do()
			if err != nil && !isGoogleAccountKeyNotFoundErr(err) {
				return err
			}
		}
		return nil
	}

	// If key is still in use, don't delete (empty keyInUse means no key is in use)
	if entry.KeyName == keyInUse {
		return nil
	}

	_, err = iamC.Projects.ServiceAccounts.Keys.Delete(entry.KeyName).Do()
	if err != nil && !isGoogleAccountKeyNotFoundErr(err) {
		return err
	}
	return nil
}

func (b *backend) serviceAccountPolicyRollback(ctx context.Context, req *logical.Request, data interface{}) error {
	b.rolesetLock.Lock()
	defer b.rolesetLock.Unlock()

	var entry walIamPolicy
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	var rolesInUse util.StringSet

	// Try to verify service account not being used by roleset
	rs, err := getRoleSet(entry.RoleSet, ctx, req.Storage)
	if err != nil {
		return err
	}
	if rs != nil && rs.AccountId != nil && rs.AccountId.ResourceName() == entry.AccountId.ResourceName() {
		rolesInUse = rs.Bindings[entry.Resource]
	}

	// Take out any bindings still being used by this role set from roles being removed.
	rolesToRemove := util.ToSet(entry.Roles)
	if len(rolesInUse) > 0 {
		rolesToRemove = rolesToRemove.Sub(rolesInUse)
	}

	r, err := b.resources.Parse(entry.Resource)
	if err != nil {
		return err
	}

	httpC, err := b.HTTPClient(req.Storage)
	if err != nil {
		return err
	}

	apiHandle := iamutil.GetApiHandle(httpC, useragent.PluginString(b.pluginEnv,
		userAgentPluginName))
	p, err := r.GetIamPolicy(ctx, apiHandle)
	if err != nil {
		if isGoogleAccountNotFoundErr(err) || isGoogleAccountUnauthorizedErr(err) {
			return nil
		}
		return err
	}

	changed, newP := p.RemoveBindings(
		&iamutil.PolicyDelta{
			Email: entry.AccountId.EmailOrId,
			Roles: rolesToRemove,
		})
	if !changed {
		return nil
	}

	_, err = r.SetIamPolicy(ctx, apiHandle, newP)
	return err
}

func (b *backend) serviceAccountPolicyDiffRollback(ctx context.Context, req *logical.Request, data interface{}) error {
	b.staticAccountLock.Lock()
	defer b.staticAccountLock.Unlock()

	var entry walIamPolicyStaticAccount
	if err := mapstructure.Decode(data, &entry); err != nil {
		return err
	}

	var rolesInUse util.StringSet

	// Try to verify service account not being used by roleset
	sa, err := b.getStaticAccount(entry.StaticAccount, ctx, req.Storage)
	if err != nil {
		return err
	}
	if sa == nil {
		b.Logger().Warn("static account %s not found, dropping WAL entry", entry.StaticAccount)
		return nil
	}
	if sa.ResourceName() == entry.AccountId.ResourceName() {
		rolesInUse = sa.Bindings[entry.Resource]
	}

	// We added roles that are not actually in use
	addedRolesToRemove := util.ToSet(entry.RolesAdded).Sub(rolesInUse)

	// We removed roles that are still saved
	removedRolesToAdd := util.ToSet(entry.RolesRemoved).Intersection(rolesInUse)

	r, err := b.resources.Parse(entry.Resource)
	if err != nil {
		return err
	}

	httpC, err := b.HTTPClient(req.Storage)
	if err != nil {
		return err
	}

	apiHandle := iamutil.GetApiHandle(httpC, useragent.PluginString(b.pluginEnv,
		userAgentPluginName))
	p, err := r.GetIamPolicy(ctx, apiHandle)
	if err != nil {
		return err
	}

	changed, newP := p.ChangeBindings(
		// toAdd
		&iamutil.PolicyDelta{
			Email: entry.AccountId.EmailOrId,
			Roles: removedRolesToAdd,
		},
		// toRemove
		&iamutil.PolicyDelta{
			Email: entry.AccountId.EmailOrId,
			Roles: addedRolesToRemove,
		})
	if !changed {
		return nil
	}

	_, err = r.SetIamPolicy(ctx, apiHandle, newP)
	return err
}

// This tries to clean up WALs that are no longer needed.
// We can ignore errors if deletion fails as WAL rollback will no-op if the object is still in use or no longer exists.
// This simply attempts to reduce the number of GCP calls we will trigger in rollbacks.
func (b *backend) tryDeleteWALs(ctx context.Context, s logical.Storage, walIds ...string) {
	for _, walId := range walIds {
		// ignore errors, WALs that are not needed will just no-op
		err := framework.DeleteWAL(ctx, s, walId)
		if err != nil {
			b.Logger().Error("unable to delete unneeded WAL %s, ignoring error since WAL will no-op: %v", walId, err)
		}
	}
}
