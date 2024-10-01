// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathSign(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "sign/" + framework.GenericNameWithAtRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationVerb:   "sign",
			OperationSuffix: "certificate",
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathSign,
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `The desired role with configuration for this request.`,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `The requested Time To Live for the SSH certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
			},
			"public_key": {
				Type:        framework.TypeString,
				Description: `SSH public key that should be signed.`,
			},
			"valid_principals": {
				Type:        framework.TypeString,
				Description: `Valid principals, either usernames or hostnames, that the certificate should be signed for.`,
			},
			"cert_type": {
				Type:        framework.TypeString,
				Description: `Type of certificate to be created; either "user" or "host".`,
				Default:     "user",
			},
			"key_id": {
				Type:        framework.TypeString,
				Description: `Key id that the created certificate should have. If not specified, the display name of the token will be used.`,
			},
			"critical_options": {
				Type:        framework.TypeMap,
				Description: `Critical options that the certificate should be signed for.`,
			},
			"extensions": {
				Type:        framework.TypeMap,
				Description: `Extensions that the certificate should be signed for.`,
			},
		},

		HelpSynopsis:    `Request signing an SSH key using a certain role with the provided details.`,
		HelpDescription: `This path allows SSH keys to be signed according to the policy of the given role.`,
	}
}

func (b *backend) pathSign(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("role").(string)

	// Get the role
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("Unknown role: %s", roleName)), nil
	}

	return b.pathSignCertificate(ctx, req, data, role)
}

func (b *backend) pathSignCertificate(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole) (*logical.Response, error) {
	publicKey := data.Get("public_key").(string)
	if publicKey == "" {
		return logical.ErrorResponse("missing public_key"), nil
	}

	userPublicKey, err := parsePublicSSHKey(publicKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse public_key as SSH key: %s", err)), nil
	}

	err = b.validateSignedKeyRequirements(userPublicKey, role)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("public_key failed to meet the key requirements: %s", err)), nil
	}

	return b.pathSignIssueCertificateHelper(ctx, req, data, role, userPublicKey)
}
