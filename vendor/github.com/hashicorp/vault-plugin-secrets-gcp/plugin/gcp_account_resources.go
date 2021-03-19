package gcpsecrets

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/iamutil"
	"github.com/hashicorp/vault-plugin-secrets-gcp/plugin/util"
	"github.com/hashicorp/vault/sdk/helper/useragent"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iam/v1"
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
		}).Do()
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
	b.Logger().Debug("creating IAM bindings", "account_email", saEmail)
	httpC, err := b.HTTPClient(req.Storage)
	if err != nil {
		return err
	}
	apiHandle := iamutil.GetApiHandle(httpC, useragent.String())

	for resourceName, roles := range binds {
		b.Logger().Debug("setting IAM binding", "resource", resourceName, "roles", roles)
		resource, err := b.resources.Parse(resourceName)
		if err != nil {
			return err
		}

		b.Logger().Debug("getting IAM policy for resource name", "name", resourceName)
		p, err := resource.GetIamPolicy(ctx, apiHandle)
		if err != nil {
			return nil
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

	return iamAdmin.Projects.ServiceAccounts.Create(fmt.Sprintf("projects/%s", project), createSaReq).Do()
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
