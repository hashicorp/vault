// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeChallenge(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeChallenge(b, baseUrl+
		"/challenge/"+framework.MatchAllRegex("auth_id")+"/"+framework.MatchAllRegex("challenge_type"), opts)
}

func addFieldsForACMEChallenge(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["auth_id"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME authorization identifier value",
		Required:    true,
	}

	fields["challenge_type"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME challenge type",
		Required:    true,
	}

	return fields
}

func patternAcmeChallenge(b *backend, pattern string, opts acmeWrapperOpts) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEChallenge(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(opts, b.acmeChallengeHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeHelpSync,
		HelpDescription: pathAcmeHelpDesc,
	}
}

func (b *backend) acmeChallengeHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, _ *acmeAccount) (*logical.Response, error) {
	authId := fields.Get("auth_id").(string)
	challengeType := fields.Get("challenge_type").(string)

	authz, err := b.GetAcmeState().LoadAuthorization(acmeCtx, userCtx, authId)
	if err != nil {
		return nil, fmt.Errorf("failed to load authorization: %w", err)
	}

	return b.acmeChallengeFetchHandler(acmeCtx, r, fields, userCtx, data, authz, challengeType)
}

func (b *backend) acmeChallengeFetchHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, authz *ACMEAuthorization, challengeType string) (*logical.Response, error) {
	var challenge *ACMEChallenge
	for _, c := range authz.Challenges {
		if string(c.Type) == challengeType {
			challenge = c
			break
		}
	}

	if challenge == nil {
		return nil, fmt.Errorf("unknown challenge of type '%v' in authorization: %w", challengeType, ErrMalformed)
	}

	// Per RFC 8555 Section 7.5.1. Responding to Challenges:
	//
	// > The client indicates to the server that it is ready for the challenge
	// > validation by sending an empty JSON body ("{}") carried in a POST
	// > request to the challenge URL (not the authorization URL).
	if len(data) > 0 {
		return nil, fmt.Errorf("unexpected request parameters: %w", ErrMalformed)
	}

	// If data was nil, we got a POST-as-GET request, just return current challenge without an accept,
	// otherwise we most likely got a "{}" payload which we should now accept the challenge.
	if data != nil {
		thumbprint, err := userCtx.GetKeyThumbprint()
		if err != nil {
			return nil, fmt.Errorf("failed to get thumbprint for key: %w", err)
		}

		if err := b.GetAcmeState().validator.AcceptChallenge(acmeCtx.sc, userCtx.Kid, authz, challenge, thumbprint); err != nil {
			return nil, fmt.Errorf("error submitting challenge for validation: %w", err)
		}
	}

	return &logical.Response{
		Data: challenge.NetworkMarshal(acmeCtx, authz.Id),

		// Per RFC 8555 Section 7.1. Resources:
		//
		// > The "up" link relation is used with challenge resources to indicate
		// > the authorization resource to which a challenge belongs.
		Headers: map[string][]string{
			"Link": {fmt.Sprintf("<%s>;rel=\"up\"", buildAuthorizationUrl(acmeCtx, authz.Id))},
		},
	}, nil
}
