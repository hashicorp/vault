// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpsecrets

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/iam/v1"
)

const (
	SecretTypeKey      = "service_account_key"
	keyAlgorithmRSA2k  = "KEY_ALG_RSA_2048"
	privateKeyTypeJson = "TYPE_GOOGLE_CREDENTIALS_FILE"
)

type secretKeyParams struct {
	keyType           string
	keyAlgorithm      string
	ttl               int
	extraInternalData map[string]interface{}
}

func secretServiceAccountKey(b *backend) *framework.Secret {
	return &framework.Secret{
		Type: SecretTypeKey,
		Fields: map[string]*framework.FieldSchema{
			"private_key_data": {
				Type:        framework.TypeString,
				Description: "Base-64 encoded string. Private key data for a service account key",
			},
			"key_algorithm": {
				Type:        framework.TypeString,
				Description: "Which type of key and algorithm to use for the key (defaults to 2K RSA). Valid values are GCP enum(ServiceAccountKeyAlgorithm)",
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: "Type of the private key (i.e. whether it is JSON or P12). Valid values are GCP enum(ServiceAccountPrivateKeyType)",
			},
		},
		Renew:  b.secretKeyRenew,
		Revoke: b.secretKeyRevoke,
	}
}

func (b *backend) secretKeyRenew(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	resp, err := b.verifySecretServiceKeyExists(ctx, req)
	if err != nil {
		return resp, err
	}
	if resp == nil {
		resp = &logical.Response{}
	}
	cfg, err := getConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = &config{}
	}

	resp.Secret = req.Secret
	resp.Secret.TTL = cfg.TTL
	resp.Secret.MaxTTL = cfg.MaxTTL
	return resp, nil
}

func (b *backend) verifyBindingsNotUpdatedForSecret(ctx context.Context, req *logical.Request) error {
	if v, ok := req.Secret.InternalData["role_set"]; ok {
		bindingSum, ok := req.Secret.InternalData["role_set_bindings"]
		if !ok {
			return fmt.Errorf("invalid secret, internal data is missing role set bindings checksum")
		}

		// Verify role set was not deleted.
		rs, err := getRoleSet(v.(string), ctx, req.Storage)
		if err != nil {
			return fmt.Errorf("could not find role set %q to verify secret", v)
		}

		// Verify role set bindings have not changed since secret was generated.
		if rs.bindingHash() != bindingSum.(string) {
			return fmt.Errorf("role set '%v' bindings were updated since secret was generated, cannot renew", v)
		}
	} else if v, ok := req.Secret.InternalData["static_account"]; ok {
		bindingSum, ok := req.Secret.InternalData["static_account_bindings"]
		if !ok {
			return fmt.Errorf("invalid secret, internal data is missing static account bindings checksum")
		}

		// Verify static account was not deleted.
		sa, err := b.getStaticAccount(v.(string), ctx, req.Storage)
		if err != nil {
			return fmt.Errorf("could not find static account %q to verify secret", v)
		}

		// Verify static account bindings have not changed since secret was generated.
		if sa.bindingHash() != bindingSum.(string) {
			return fmt.Errorf("static account '%v' bindings were updated since secret was generated, cannot renew", v)
		}
	} else {
		return fmt.Errorf("invalid secret, internal data is missing role set or static account name")
	}

	return nil
}

func (b *backend) verifySecretServiceKeyExists(ctx context.Context, req *logical.Request) (*logical.Response, error) {
	keyName, ok := req.Secret.InternalData["key_name"]
	if !ok {
		return nil, fmt.Errorf("invalid secret, internal data is missing key name")
	}

	if err := b.verifyBindingsNotUpdatedForSecret(ctx, req); err != nil {
		return logical.ErrorResponse(err.Error()), err
	}

	// Verify service account key still exists.
	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return logical.ErrorResponse("could not confirm key still exists in GCP"), nil
	}

	if k, err := iamAdmin.Projects.ServiceAccounts.Keys.Get(keyName.(string)).Do(); err != nil || k == nil {
		return logical.ErrorResponse("could not confirm key still exists in GCP: %v", err), nil
	}

	return nil, nil
}

func (b *backend) secretKeyRevoke(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keyNameRaw, ok := req.Secret.InternalData["key_name"]
	if !ok {
		return nil, fmt.Errorf("secret is missing key_name internal data")
	}

	iamAdmin, err := b.IAMAdminClient(req.Storage)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	_, err = iamAdmin.Projects.ServiceAccounts.Keys.Delete(keyNameRaw.(string)).Context(ctx).Do()
	if err != nil && !isGoogleAccountKeyNotFoundErr(err) {
		return logical.ErrorResponse("unable to delete service account key: %v", err), nil
	}

	return nil, nil
}

func (b *backend) createServiceAccountKeySecret(ctx context.Context, s logical.Storage, id *gcputil.ServiceAccountId, params secretKeyParams) (*logical.Response, error) {
	cfg, err := getConfig(ctx, s)
	if err != nil {
		return nil, errwrap.Wrapf("could not read backend config: {{err}}", err)
	}
	if cfg == nil {
		cfg = &config{}
	}

	iamC, err := b.IAMAdminClient(s)
	if err != nil {
		return nil, errwrap.Wrapf("could not create IAM Admin client: {{err}}", err)
	}

	key, err := iamC.Projects.ServiceAccounts.Keys.Create(
		id.ResourceName(), &iam.CreateServiceAccountKeyRequest{
			KeyAlgorithm:   params.keyAlgorithm,
			PrivateKeyType: params.keyType,
		}).Do()
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	secretD := map[string]interface{}{
		"private_key_data": key.PrivateKeyData,
		"key_algorithm":    key.KeyAlgorithm,
		"key_type":         key.PrivateKeyType,
	}
	internalD := map[string]interface{}{
		"key_name": key.Name,
	}

	for k, v := range params.extraInternalData {
		internalD[k] = v
	}

	resp := b.Secret(SecretTypeKey).Response(secretD, internalD)
	resp.Secret.Renewable = true

	resp.Secret.MaxTTL = cfg.MaxTTL
	resp.Secret.TTL = cfg.TTL

	// If the request came with a TTL value, overwrite the config default
	if params.ttl > 0 {
		resp.Secret.TTL = time.Duration(params.ttl) * time.Second
	}

	return resp, nil
}

const pathServiceAccountKeySyn = `Generate a service account private key secret.`
const pathServiceAccountKeyDesc = `
This path will generate a new service account key for accessing GCP APIs.

Either specify "roleset/my-roleset" or "static/my-account" to generate a key corresponding
to a roleset or static account respectively.

Please see backend documentation for more information:
https://www.vaultproject.io/docs/secrets/gcp/index.html
`
