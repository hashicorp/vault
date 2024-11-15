// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/vault/sdk/helper/jsonutil"
	"math/rand"
	"time"

	"github.com/hashicorp/vault-plugin-secrets-alicloud/clients"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func (b *backend) pathCreds() *framework.Path {
	return &framework.Path{
		Pattern: "creds/" + framework.GenericNameRegex("name"),
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationVerb:   "generate",
			OperationSuffix: "credentials",
		},
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "The name of the role.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.ReadOperation: b.operationCredsRead,
		},
		HelpSynopsis:    pathCredsHelpSyn,
		HelpDescription: pathCredsHelpDesc,
	}
}

func (b *backend) operationCredsRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return nil, errors.New("name is required")
	}

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		// Attempting to read a role that doesn't exist.
		return nil, nil
	}

	creds, err := readCredentials(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if creds == nil {
		return nil, errors.New("unable to create secret because no credentials are configured")
	}

	switch role.Type() {

	case roleTypeSTS:
		client, err := clients.NewSTSClient(b.sdkConfig, creds.AccessKey, creds.SecretKey)
		if err != nil {
			return nil, err
		}
		assumeRoleResp, err := client.AssumeRole(generateRoleSessionName(req.DisplayName, roleName), role.RoleARN)
		if err != nil {
			return nil, err
		}
		// Parse the expiration into a time, so that when we return it from our API it's formatted
		// the same way as how _we_ format times, which could differ from this over time.
		expiration, err := time.Parse("2006-01-02T15:04:05Z", assumeRoleResp.Credentials.Expiration)
		if err != nil {
			return nil, err
		}
		resp := b.Secret(secretType).Response(map[string]interface{}{
			"access_key":     assumeRoleResp.Credentials.AccessKeyId,
			"secret_key":     assumeRoleResp.Credentials.AccessKeySecret,
			"security_token": assumeRoleResp.Credentials.SecurityToken,
			"expiration":     expiration,
		}, map[string]interface{}{
			"role_type": roleTypeSTS.String(),
		})

		// Set the secret TTL to appropriately match the expiration of the token.
		ttl := expiration.Sub(time.Now())
		resp.Secret.TTL = ttl
		resp.Secret.MaxTTL = ttl

		// STS credentials are purposefully short-lived and aren't renewable.
		resp.Secret.Renewable = false
		return resp, nil

	case roleTypeRAM:
		client, err := clients.NewRAMClient(b.sdkConfig, creds.AccessKey, creds.SecretKey)
		if err != nil {
			return nil, err
		}

		/*
			Now we're embarking upon a multi-step process that could fail at any time.
			If it does, let's do our best to clean up after ourselves. Success will be
			our flag at the end indicating whether we should leave things be, or clean
			things up, based on how we exit this method. Since defer statements are
			last-in-first-out, it will perfectly reverse the order of everything just
			like we need.
		*/
		success := false

		createUserResp, err := client.CreateUser(generateUsername(req.DisplayName, roleName))
		if err != nil {
			return nil, err
		}
		defer func() {
			if success {
				return
			}
			if err := client.DeleteUser(createUserResp.User.UserName); err != nil {
				if b.Logger().IsError() {
					b.Logger().Error(fmt.Sprintf("unable to delete user %s", createUserResp.User.UserName), err)
				}
			}
		}()

		// We need to gather up all the names and types of the remote policies we're
		// about to create so we can detach and delete them later.
		inlinePolicies := make([]*remotePolicy, len(role.InlinePolicies))

		for i, inlinePolicy := range role.InlinePolicies {

			// By combining the userName with the particular policy's UUID,
			// it'll be possible to figure out who this policy is for and which one
			// it is using the policy name alone. The max length of a policy name is
			// 128, but given the max lengths of our username and inline policy UUID,
			// we will always remain well under that here.
			policyName := createUserResp.User.UserName + "-" + inlinePolicy.UUID

			policyDoc, err := jsonutil.EncodeJSON(inlinePolicy.PolicyDocument)
			if err != nil {
				return nil, err
			}

			createPolicyResp, err := client.CreatePolicy(policyName, string(policyDoc))
			if err != nil {
				return nil, err
			}

			inlinePolicies[i] = &remotePolicy{
				Name: createPolicyResp.Policy.PolicyName,
				Type: createPolicyResp.Policy.PolicyType,
			}

			// This defer is in this loop on purpose.
			defer func() {
				if success {
					return
				}
				if err := client.DeletePolicy(createPolicyResp.Policy.PolicyName); err != nil {
					if b.Logger().IsError() {
						b.Logger().Error(fmt.Sprintf("unable to delete policy %s", createPolicyResp.Policy.PolicyName), err)
					}
				}
			}()

			if err := client.AttachPolicy(createUserResp.User.UserName, createPolicyResp.Policy.PolicyName, createPolicyResp.Policy.PolicyType); err != nil {
				return nil, err
			}
			// This defer is also in this loop on purpose.
			defer func() {
				if success {
					return
				}
				if err := client.DetachPolicy(createUserResp.User.UserName, createPolicyResp.Policy.PolicyName, createPolicyResp.Policy.PolicyType); err != nil {
					if b.Logger().IsError() {
						b.Logger().Error(fmt.Sprintf(
							"unable to detach policy name:%s, type:%s from user:%s", createPolicyResp.Policy.PolicyName, createPolicyResp.Policy.PolicyType, createUserResp.User.UserName))
					}
				}
			}()
		}

		for _, remotePol := range role.RemotePolicies {
			if err := client.AttachPolicy(createUserResp.User.UserName, remotePol.Name, remotePol.Type); err != nil {
				return nil, err
			}
			// This defer is also in this loop on purpose.
			// Separate these strings from the remotePol pointer so the defer statement will retain the correct values
			// due to pointer reuse for the remotePol var on each iteration of the loop.
			remotePolName := remotePol.Name
			remotePolType := remotePol.Type
			defer func() {
				if success {
					return
				}
				if err := client.DetachPolicy(createUserResp.User.UserName, remotePolName, remotePolType); err != nil {
					if b.Logger().IsError() {
						b.Logger().Error(fmt.Sprintf("unable to detach policy name:%s, type:%s from user:%s", remotePolName, remotePolType, createUserResp.User.UserName))
					}
				}
			}()
		}

		accessKeyResp, err := client.CreateAccessKey(createUserResp.User.UserName)
		if err != nil {
			return nil, err
		}
		// It's unlikely we wouldn't have success at this point because there are no further errors returned below, but
		// there could be a panic if somehow one of the objects below were missing a pointer, so let's play it safe and
		// add a defer rolling back the access key if that happens.
		defer func() {
			if success {
				return
			}
			if err := client.DeleteAccessKey(createUserResp.User.UserName, accessKeyResp.AccessKey.AccessKeyId); err != nil {
				if b.Logger().IsError() {
					b.Logger().Error(fmt.Sprintf("unable to delete access key for username:%s", createUserResp.User.UserName))
				}
			}
		}()

		resp := b.Secret(secretType).Response(map[string]interface{}{
			"access_key": accessKeyResp.AccessKey.AccessKeyId,
			"secret_key": accessKeyResp.AccessKey.AccessKeySecret,
		}, map[string]interface{}{
			"role_type":       roleTypeRAM.String(),
			"role_name":       roleName,
			"username":        createUserResp.User.UserName,
			"access_key_id":   accessKeyResp.AccessKey.AccessKeyId,
			"inline_policies": inlinePolicies,
			"remote_policies": role.RemotePolicies,
		})
		if role.TTL != 0 {
			resp.Secret.TTL = role.TTL
		}
		if role.MaxTTL != 0 {
			resp.Secret.MaxTTL = role.MaxTTL
		}

		success = true
		return resp, nil

	default:
		return nil, fmt.Errorf("unsupported role type: %s", role.Type())
	}
}

// The max length of a username per AliCloud is 64.
func generateUsername(displayName, roleName string) string {
	return generateName(displayName, roleName, 64)
}

// The max length of a role session name per AliCloud is 32.
func generateRoleSessionName(displayName, roleName string) string {
	return generateName(displayName, roleName, 32)
}

func generateName(displayName, roleName string, maxLength int) string {
	name := fmt.Sprintf("%s-%s-", displayName, roleName)

	// The time and random number take up to 15 more in length, so if the name
	// is too long we need to trim it.
	if len(name) > maxLength-15 {
		name = name[:maxLength-15]
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%s%d-%d", name, time.Now().Unix(), r.Intn(10000))
}

const pathCredsHelpSyn = `
Generate an API key or STS credential using the given role's configuration.'
`

const pathCredsHelpDesc = `
This path will generate a new API key or STS credential for
accessing AliCloud. The RAM policies used to back this key pair will be
configured on the role. For example, if this backend is mounted at "alicloud",
then "alicloud/creds/deploy" would generate access keys for the "deploy" role.

The API key or STS credential will have a ttl associated with it. API keys can
be renewed or revoked as described here: 
https://www.vaultproject.io/docs/concepts/lease.html,
but STS credentials do not support renewal or revocation.
`
