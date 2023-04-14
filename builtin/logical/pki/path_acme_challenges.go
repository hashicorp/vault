package pki

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeRootChallenge(b *backend) *framework.Path {
	return patternAcmeChallenge(b,
		"acme/challenge/"+framework.MatchAllRegex("auth_id")+"/"+
			framework.MatchAllRegex("challenge_type"))
}

func pathAcmeRoleChallenge(b *backend) *framework.Path {
	return patternAcmeChallenge(b,
		"roles/"+framework.GenericNameRegex("role")+"/acme/challenge/"+
			framework.MatchAllRegex("auth_id")+"/"+
			framework.MatchAllRegex("challenge_type"))
}

func pathAcmeIssuerChallenge(b *backend) *framework.Path {
	return patternAcmeChallenge(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/challenge/"+
			framework.MatchAllRegex("auth_id")+"/"+
			framework.MatchAllRegex("challenge_type"))
}

func pathAcmeIssuerAndRoleChallenge(b *backend) *framework.Path {
	return patternAcmeChallenge(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/challenge/"+
			framework.MatchAllRegex("auth_id")+"/"+
			framework.MatchAllRegex("challenge_type"))
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

func patternAcmeChallenge(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEChallenge(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeParsedWrapper(b.acmeChallengeHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (b *backend) acmeChallengeHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
	authId := fields.Get("auth_id").(string)
	challengeType := fields.Get("challenge_type").(string)

	authz, err := b.acmeState.LoadAuthorization(acmeCtx, userCtx, authId)
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

	// XXX: Prompt for challenge to be tried by the server.

	return &logical.Response{
		Data: challenge.NetworkMarshal(),
	}, nil
}
