// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/vault-plugin-secrets-alicloud/clients"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"github.com/hashicorp/vault/sdk/logical"
)

const secretType = "alicloud"

func (b *backend) pathSecrets() *framework.Secret {
	return &framework.Secret{
		Type: secretType,
		Fields: map[string]*framework.FieldSchema{
			"access_key": {
				Type:        framework.TypeString,
				Description: "Access Key",
			},
			"secret_key": {
				Type:        framework.TypeString,
				Description: "Secret Key",
			},
		},
		Renew:  b.operationRenew,
		Revoke: b.operationRevoke,
	}
}

func (b *backend) operationRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleTypeRaw, ok := req.Secret.InternalData["role_type"]
	if !ok {
		return nil, errors.New("role_type missing from secret")
	}
	nameOfRoleType, ok := roleTypeRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unable to read role_type: %+v", roleTypeRaw)
	}
	rType, err := parseRoleType(nameOfRoleType)
	if err != nil {
		return nil, err
	}

	switch rType {

	case roleTypeSTS:
		// STS already has a lifetime, and we don'nameOfRoleType support renewing it.
		return nil, nil

	case roleTypeRAM:
		roleName, err := getStringValue(req.Secret.InternalData, "role_name")
		if err != nil {
			return nil, err
		}

		role, err := readRole(ctx, req.Storage, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			// The role has been deleted since the secret was issued or last renewed.
			// The user's expectation is probably that the caller won'nameOfRoleType continue being
			// able to perform renewals.
			return nil, fmt.Errorf("role %s has been deleted so no further renewals are allowed", roleName)
		}

		resp := &logical.Response{Secret: req.Secret}
		if role.TTL != 0 {
			resp.Secret.TTL = role.TTL
		}
		if role.MaxTTL != 0 {
			resp.Secret.MaxTTL = role.MaxTTL
		}
		return resp, nil

	default:
		return nil, fmt.Errorf("unrecognized role_type: %s", nameOfRoleType)
	}
}

func (b *backend) operationRevoke(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	roleTypeRaw, ok := req.Secret.InternalData["role_type"]
	if !ok {
		return nil, errors.New("role_type missing from secret")
	}
	nameOfRoleType, ok := roleTypeRaw.(string)
	if !ok {
		return nil, fmt.Errorf("unable to read role_type: %+v", roleTypeRaw)
	}
	rType, err := parseRoleType(nameOfRoleType)
	if err != nil {
		return nil, err
	}

	switch rType {
	case roleTypeSTS:
		// STS cleans up after itself so we can skip this if is_sts internal data
		// element set to true.
		return nil, nil

	case roleTypeRAM:
		creds, err := readCredentials(ctx, req.Storage)
		if err != nil {
			return nil, err
		}
		if creds == nil {
			return nil, errors.New("unable to delete access key because no credentials are configured")
		}
		client, err := clients.NewRAMClient(b.sdkConfig, creds.AccessKey, creds.SecretKey)
		if err != nil {
			return nil, err
		}

		userName, err := getStringValue(req.Secret.InternalData, "username")
		if err != nil {
			return nil, err
		}

		accessKeyID, err := getStringValue(req.Secret.InternalData, "access_key_id")
		if err != nil {
			return nil, err
		}

		// On the first pass here, if we successfully delete an access key but fail later and return an
		// error, we'll never again be able to progress past deleting the access key because it'll be
		// gone, leaving dangling objects. So, even if we get an error, we still try to do the remaining
		// things. Multierror flags whether to return an error or nil, and if it return any error we'll
		// try again, but if it doesn't we know we're done.
		apiErrs := &multierror.Error{}

		// Delete the access key first so if all else fails, the access key is revoked.
		if err := client.DeleteAccessKey(userName, accessKeyID); err != nil {
			apiErrs = multierror.Append(apiErrs, err)
		}

		// Inline policies are currently stored as remote policies, because they have been
		// instantiated remotely and we need their name and type to now detach and delete them.
		inlinePolicies, err := getRemotePolicies(req.Secret.InternalData, "inline_policies")
		if err != nil {
			// This shouldn't be part of the multierror because if it returns empty inline policies,
			// then we won't go through the inlinePolicies loop and we'll think we're successful
			// when we actually didn't delete the inlinePolicies we need to.
			return nil, err
		}
		for _, inlinePolicy := range inlinePolicies {
			if err := client.DetachPolicy(userName, inlinePolicy.Name, inlinePolicy.Type); err != nil {
				apiErrs = multierror.Append(apiErrs, err)
			}
			if err := client.DeletePolicy(inlinePolicy.Name); err != nil {
				apiErrs = multierror.Append(apiErrs, err)
			}
		}

		// These just need to be detached, but we're not going to delete them because they're
		// supposed to be longstanding.
		remotePolicies, err := getRemotePolicies(req.Secret.InternalData, "remote_policies")
		if err != nil {
			// This shouldn't be part of the multierror because if it returns empty remote policies,
			// then we won't go through the remotePolicies loop and we'll think we're successful
			// when we actually didn't delete the remotePolicies we need to.
			return nil, err
		}
		for _, remotePolicy := range remotePolicies {
			if err := client.DetachPolicy(userName, remotePolicy.Name, remotePolicy.Type); err != nil {
				apiErrs = multierror.Append(apiErrs, err)
			}
		}

		// Finally, delete the user. Note: this will fail if any other new associations have been
		// created with the user out of band from Vault. For example, if a new API key had been
		// manually created for them in their console that Vault didn't know about, or some other
		// thing had been created. Luckily the err returned is pretty explanatory so that will
		// help with debugging.
		if err := client.DeleteUser(userName); err != nil {
			apiErrs = multierror.Append(apiErrs, err)
		}
		return nil, apiErrs.ErrorOrNil()

	default:
		return nil, fmt.Errorf("unrecognized role_type: %s", nameOfRoleType)
	}
}

func getStringValue(internalData map[string]interface{}, key string) (string, error) {
	valueRaw, ok := internalData[key]
	if !ok {
		return "", fmt.Errorf("secret is missing %s internal data", key)
	}
	value, ok := valueRaw.(string)
	if !ok {
		return "", fmt.Errorf("secret is missing %s internal data", key)
	}
	return value, nil
}

func getRemotePolicies(internalData map[string]interface{}, key string) ([]*remotePolicy, error) {
	valuesRaw, ok := internalData[key]
	if !ok {
		return nil, fmt.Errorf("secret is missing %s internal data", key)
	}

	valuesJSON, err := jsonutil.EncodeJSON(valuesRaw)
	if err != nil {
		return nil, fmt.Errorf("malformed %s internal data", key)
	}

	policies := []*remotePolicy{}
	if err := jsonutil.DecodeJSON(valuesJSON, &policies); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s internal data as remotePolicy", key)
	}
	return policies, nil
}
