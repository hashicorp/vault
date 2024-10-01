// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package ssh

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type keySpecs struct {
	Type string
	Bits int
}

func pathIssue(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "issue/" + framework.GenericNameWithAtRegex("role"),

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixSSH,
			OperationVerb:   "issue",
			OperationSuffix: "certificate",
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathIssue,
			},
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `The desired role with configuration for this request.`,
			},
			"key_type": {
				Type:        framework.TypeString,
				Description: "Specifies the desired key type; must be `rsa`, `ed25519` or `ec`",
				Default:     "rsa",
			},
			"key_bits": {
				Type:        framework.TypeInt,
				Description: "Specifies the number of bits to use for the generated keys.",
				Default:     0,
			},
			"ttl": {
				Type: framework.TypeDurationSecond,
				Description: `The requested Time To Live for the SSH certificate;
sets the expiration date. If not specified
the role default, backend default, or system
default TTL is used, in that order. Cannot
be later than the role max TTL.`,
			},
			"valid_principals": {
				Type:        framework.TypeString,
				Description: `Valid principals, either usernames or hostnames, that the certificate should be signed for.  Must be non-empty unless allow_empty_principals=true (not recommended) or a value for DefaultUser has been set in the role`,
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
		HelpSynopsis:    pathIssueHelpSyn,
		HelpDescription: pathIssueHelpDesc,
	}
}

func (b *backend) pathIssue(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Get the role
	roleName := data.Get("role").(string)
	role, err := b.getRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("unknown role: %s", roleName)), nil
	}

	if role.KeyType != "ca" {
		return logical.ErrorResponse("role key type '%s' not allowed to issue key pairs", role.KeyType), nil
	}

	// Validate and extract key specifications
	keySpecs, err := extractKeySpecs(role, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	// Issue certificate
	return b.pathIssueCertificate(ctx, req, data, role, keySpecs)
}

func (b *backend) pathIssueCertificate(ctx context.Context, req *logical.Request, data *framework.FieldData, role *sshRole, keySpecs *keySpecs) (*logical.Response, error) {
	publicKey, privateKey, err := generateSSHKeyPair(rand.Reader, keySpecs.Type, keySpecs.Bits)
	if err != nil {
		return nil, err
	}

	// Sign key
	userPublicKey, err := parsePublicSSHKey(publicKey)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("failed to parse public_key as SSH key: %s", err)), nil
	}

	response, err := b.pathSignIssueCertificateHelper(ctx, req, data, role, userPublicKey)
	if err != nil {
		return nil, err
	}
	if response.IsError() {
		return response, nil
	}

	// Additional to sign response
	response.Data["private_key"] = privateKey
	response.Data["private_key_type"] = keySpecs.Type

	return response, nil
}

func extractKeySpecs(role *sshRole, data *framework.FieldData) (*keySpecs, error) {
	keyType := data.Get("key_type").(string)
	keyBits := data.Get("key_bits").(int)
	keySpecs := keySpecs{
		Type: keyType,
		Bits: keyBits,
	}

	keyTypeToMapKey := createKeyTypeToMapKey(keyType, keyBits)

	if len(role.AllowedUserKeyTypesLengths) != 0 {
		var keyAllowed bool
		var bitsAllowed bool

	keyTypeAliasesLoop:
		for _, keyTypeAlias := range keyTypeToMapKey[keyType] {
			allowedValues, allowed := role.AllowedUserKeyTypesLengths[keyTypeAlias]
			if !allowed {
				continue
			}
			keyAllowed = true

			for _, value := range allowedValues {
				if value == keyBits {
					bitsAllowed = true
					break keyTypeAliasesLoop
				}
			}
		}

		if !keyAllowed {
			return nil, errors.New("provided key_type value not in allowed_user_key_types")
		}

		if !bitsAllowed {
			return nil, errors.New("provided key_bits value not in list of role's allowed_user_key_types")
		}
	}

	return &keySpecs, nil
}

const pathIssueHelpSyn = `
Request a certificate using a certain role with the provided details.
`

const pathIssueHelpDesc = `
This path allows requesting a certificate to be issued according to the
policy of the given role. The certificate will only be issued if the
requested details are allowed by the role policy.

This path returns a certificate and a private key. If you want a workflow
that does not expose a private key, generate a CSR locally and use the
sign path instead.
`
