// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/iam/v1"
)

const (
	flagCanDeleteServiceAccount = true
	flagMustKeepServiceAccount  = false
)

type (
	// gcpAccountResources is a wrapper around the GCP resources Vault creates to generate credentials.
	// This includes a Vault-managed GCP service account (required), IAM bindings, and/or key via TokenGenerator
	// (for generating access tokens).
	gcpAccountResources struct {
		accountId gcputil.ServiceAccountId
		bindings  ResourceBindings
		tokenGen  *TokenGenerator
	}

	// ResourceBindings represent a map of GCP resource name to IAM roles to be bound on that resource.
	ResourceBindings map[string]util.StringSet

	// TokenGenerator wraps the service account key and params required to create access tokens.
	TokenGenerator struct {
		KeyName    string
		B64KeyJSON string
		Scopes     []string
	}
)

func (rb ResourceBindings) asOutput() map[string][]string {
	out := make(map[string][]string, len(rb))
	for k, v := range rb {
		out[k] = v.ToSlice()
	}
	return out
}

func (rb ResourceBindings) sub(toRemove ResourceBindings) ResourceBindings {
	subbed := make(ResourceBindings)
	for r, iamRoles := range rb {
		toRemoveIamRoles, ok := toRemove[r]
		if ok {
			iamRoles = iamRoles.Sub(toRemoveIamRoles)
		}
		subbed[r] = iamRoles
	}
	return subbed
}

func getStringHash(bindingsRaw string) string {
	ssum := sha256.Sum256([]byte(bindingsRaw))
	return base64.StdEncoding.EncodeToString(ssum[:])
}

func (b *backend) createNewTokenGen(ctx context.Context, req *logical.Request, parent string, scopes []string) (*TokenGenerator, error) {
	b.Logger().Debug("creating new TokenGenerator (service account key)", "account", parent, "scopes", scopes)

	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, err
	}

	key, err := iamAdmin.Projects.ServiceAccounts.Keys.Create(
		parent,
		&iam.CreateServiceAccountKeyRequest{
			PrivateKeyType: privateKeyTypeJson,
		}).Context(ctx).Do()
	if err != nil {
		return nil, err
	}
	return &TokenGenerator{
		KeyName:    key.Name,
		B64KeyJSON: key.PrivateKeyData,
		Scopes:     scopes,
	}, nil
}

func (b *backend) createIamBindings(ctx context.Context, req *logical.Request, saEmail string, binds ResourceBindings) error {
	b.Logger().Debug("creating IAM bindings", "account_email", saEmail, "bindings", binds)
	httpC, err := b.HTTPClient(req.Storage)
	if err != nil {
		return err
	}
	apiHandle := iamutil.GetApiHandle(httpC, useragent.PluginString(b.pluginEnv,
		userAgentPluginName))

	for resourceName, roles := range binds {
		b.Logger().Debug("setting IAM binding", "resource", resourceName, "roles", roles)
		resource, err := b.resources.Parse(resourceName)
		if err != nil {
			return err
		}

		b.Logger().Debug("getting IAM policy for resource name", "name", resourceName)
		p, err := resource.GetIamPolicy(ctx, apiHandle)
		if err != nil {
			return err
		}

		b.Logger().Debug("got IAM policy for resource name", "name", resourceName)
		changed, newP := p.AddBindings(&iamutil.PolicyDelta{
			Roles: roles,
			Email: saEmail,
		})
		if !changed || newP == nil {
			continue
		}

		b.Logger().Debug("setting IAM policy for resource name", "name", resourceName)
		if _, err := resource.SetIamPolicy(ctx, apiHandle, newP); err != nil {
			return errwrap.Wrapf(fmt.Sprintf("unable to set IAM policy for resource %q: {{err}}", resourceName), err)
		}
	}

	return nil
}

func (b *backend) createServiceAccount(ctx context.Context, req *logical.Request, project, saName, descriptor string) (*iam.ServiceAccount, error) {
	createSaReq := &iam.CreateServiceAccountRequest{
		AccountId: saName,
		ServiceAccount: &iam.ServiceAccount{
			DisplayName: roleSetServiceAccountDisplayName(descriptor),
		},
	}

	b.Logger().Debug("creating service account",
		"project", project,
		"request", createSaReq)

	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return nil, err
	}

	return iamAdmin.Projects.ServiceAccounts.Create(fmt.Sprintf("projects/%s", project), createSaReq).Context(ctx).Do()
}

// tryDeleteGcpAccountResources creates WALs to clean up a service account's
// bindings, key, and account (if removeServiceAccount is true)
func (b *backend) tryDeleteGcpAccountResources(ctx context.Context, req *logical.Request, boundResources *gcpAccountResources, removeServiceAccount bool, walIds []string) []string {
	if boundResources == nil {
		b.Logger().Debug("skip deletion for nil roleset resources")
		return nil
	}

	b.Logger().Debug("try to delete GCP account resources", "bound_resources", boundResources, "remove_service_account", removeServiceAccount)

	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return []string{err.Error()}
	}

	warnings := make([]string, 0)
	if boundResources.tokenGen != nil {
		if err := b.deleteTokenGenKey(ctx, iamAdmin, boundResources.tokenGen); err != nil {
			w := fmt.Sprintf("unable to delete key under service account %q (WAL entry to clean-up later has been added): %v", boundResources.accountId.ResourceName(), err)
			warnings = append(warnings, w)
		}
	}

	if merr := b.removeBindings(ctx, req, boundResources.accountId.EmailOrId, boundResources.bindings); merr != nil {
		for _, err := range merr.Errors {
			w := fmt.Sprintf("unable to delete IAM policy bindings for service account %q (WAL entry to clean-up later has been added): %v", boundResources.accountId.EmailOrId, err)
			warnings = append(warnings, w)
		}
	}

	if removeServiceAccount {
		if err := b.deleteServiceAccount(ctx, iamAdmin, boundResources.accountId); err != nil {
			w := fmt.Sprintf("unable to delete service account %q (WAL entry to clean-up later has been added): %v", boundResources.accountId.ResourceName(), err)
			warnings = append(warnings, w)
		}
	}

	// If resources were deleted, we don't need the WAL rollbacks we created for these resources.
	if len(warnings) == 0 {
		b.tryDeleteWALs(ctx, req.Storage, walIds...)
	}

	return nil
}

func (b *backend) deleteTokenGenKey(ctx context.Context, iamAdmin *iam.Service, tgen *TokenGenerator) error {
	if tgen == nil || tgen.KeyName == "" {
		return nil
	}

	_, err := iamAdmin.Projects.ServiceAccounts.Keys.Delete(tgen.KeyName).Context(ctx).Do()
	if err != nil && !isGoogleAccountKeyNotFoundErr(err) {
		return errwrap.Wrapf("unable to delete service account key: {{err}}", err)
	}
	return nil
}

func (b *backend) removeBindings(ctx context.Context, req *logical.Request, email string, bindings ResourceBindings) (allErr *multierror.Error) {
	httpC, err := b.HTTPClient(req.Storage)
	if err != nil {
		return &multierror.Error{Errors: []error{err}}
	}

	apiHandle := iamutil.GetApiHandle(httpC, useragent.PluginString(b.pluginEnv,
		userAgentPluginName))

	for resName, roles := range bindings {
		resource, err := b.resources.Parse(resName)
		if err != nil {
			allErr = multierror.Append(allErr, errwrap.Wrapf(fmt.Sprintf("unable to delete role binding for resource '%s': {{err}}", resName), err))
			continue
		}

		p, err := resource.GetIamPolicy(ctx, apiHandle)
		if err != nil {
			allErr = multierror.Append(allErr, errwrap.Wrapf(fmt.Sprintf("unable to delete role binding for resource '%s': {{err}}", resName), err))
			continue
		}

		changed, newP := p.RemoveBindings(&iamutil.PolicyDelta{
			Email: email,
			Roles: roles,
		})
		if !changed {
			continue
		}
		if _, err = resource.SetIamPolicy(ctx, apiHandle, newP); err != nil {
			allErr = multierror.Append(allErr, errwrap.Wrapf(fmt.Sprintf("unable to delete role binding for resource '%s': {{err}}", resName), err))
			continue
		}
	}
	return
}

func (b *backend) deleteServiceAccount(ctx context.Context, iamAdmin *iam.Service, account gcputil.ServiceAccountId) error {
	if account.EmailOrId == "" {
		return nil
	}

	_, err := iamAdmin.Projects.ServiceAccounts.Delete(account.ResourceName()).Context(ctx).Do()
	if err != nil && !isGoogleAccountNotFoundErr(err) {
		return errwrap.Wrapf("unable to delete service account: {{err}}", err)
	}
	return nil
}

func isGoogleAccountNotFoundErr(err error) bool {
	return isGoogleApiErrorWithCodes(err, 404)
}

func isGoogleAccountKeyNotFoundErr(err error) bool {
	return isGoogleApiErrorWithCodes(err, 403, 404)
}

func isGoogleAccountUnauthorizedErr(err error) bool {
	return isGoogleApiErrorWithCodes(err, 403)
}

func isGoogleApiErrorWithCodes(err error, validErrCodes ...int) bool {
	if err == nil {
		return false
	}

	gErr, ok := err.(*googleapi.Error)
	if !ok {
		wrapErrV := errwrap.GetType(err, &googleapi.Error{})
		if wrapErrV == nil {
			return false
		}
		gErr = wrapErrV.(*googleapi.Error)
	}

	for _, code := range validErrCodes {
		if gErr.Code == code {
			return true
		}
	}

	return false
}

func emailForServiceAccountName(project, accountName string) string {
	return fmt.Sprintf(serviceAccountEmailTemplate, accountName, project)
}

func roleSetServiceAccountDisplayName(name string) string {
	fullDisplayName := fmt.Sprintf(serviceAccountDisplayNameTmpl, name)
	displayName := fullDisplayName
	if len(fullDisplayName) > serviceAccountDisplayNameMaxLen {
		truncIndex := serviceAccountDisplayNameMaxLen - serviceAccountDisplayNameHashLen
		h := fmt.Sprintf("%x", sha256.Sum256([]byte(fullDisplayName[truncIndex:])))
		displayName = fullDisplayName[:truncIndex] + h[:serviceAccountDisplayNameHashLen]
	}
	return displayName
}
