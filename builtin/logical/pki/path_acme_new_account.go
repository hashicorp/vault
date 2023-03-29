package pki

import (
	"fmt"

	"github.com/hashicorp/vault/builtin/logical/pki/acme"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeRootNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "acme/directory")
}

func pathAcmeRoleNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "roles/"+framework.GenericNameRegex("role")+"/acme/directory")
}

func pathAcmeIssuerNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/directory")
}

func pathAcmeIssuerAndRoleNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/roles/"+framework.GenericNameRegex("role")+"/acme/directory")
}

func patternAcmeNewAccount(b *backend, pattern string) *framework.Path {
	return &framework.Path{
		Pattern: pattern,
		Fields:  map[string]*framework.FieldSchema{},
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback:                    b.acmeParsedWrapper(b.acmeNewAccountHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathOcspHelpSyn,
		HelpDescription: pathOcspHelpDesc,
	}
}

type acmeParsedOperation func(acmeCtx acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *acme.JWSCtx, data map[string]interface{}) (*logical.Response, error)

func (b *backend) acmeParsedWrapper(op acmeParsedOperation) framework.OperationFunc {
	return b.acmeWrapper(func(acmeCtx acmeContext, r *logical.Request, fields *framework.FieldData) (*logical.Response, error) {
		user, data, err := b.acme.ParseRequestParams(fields)
		if err != nil {
			return nil, err
		}

		return op(acmeCtx, r, fields, user, data)
	})
}

func (b *backend) acmeNewAccountHandler(acmeCtx acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *acme.JWSCtx, data map[string]interface{}) (*logical.Response, error) {
	// Parameters
	var ok bool
	var onlyReturnExisting bool
	var contact []string
	var termsOfServiceAgreed bool

	rawContact, present := data["contact"]
	if present {
		contact, ok = rawContact.([]string)
		if !ok {
			return nil, fmt.Errorf("invalid type for field 'contact': %w", acme.ErrMalformed)
		}
	}

	rawTermsOfServiceAgreed, present := data["termsOfServiceAgreed"]
	if present {
		termsOfServiceAgreed, ok = rawTermsOfServiceAgreed.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type for field 'termsOfServiceAgreed': %w", acme.ErrMalformed)
		}
	}

	rawOnlyReturnExisting, present := data["onlyReturnExisting"]
	if present {
		onlyReturnExisting, ok = rawOnlyReturnExisting.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type for field 'onlyReturnExisting': %w", acme.ErrMalformed)
		}
	}

	// We ignore the EAB parameter as it is currently not supported.

	// We have two paths here: search or create.
	if onlyReturnExisting {
		return b.acmeNewAccountSearchHandler(acmeCtx, r, fields, userCtx, data)
	}

	return b.acmeNewAccountCreateHandler(acmeCtx, r, fields, userCtx, data, contact, termsOfServiceAgreed)
}

func formatAccountResponse(location string, status string, contact []string) *logical.Response {
	resp := &logical.Response{
		Data: map[string]interface{}{
			"status": status,
			"orders": location + "/orders",
		},
	}

	if len(contact) > 0 {
		resp.Data["contact"] = contact
	}

	resp.Headers["Location"] = []string{location}

	return resp
}

func (b *backend) acmeNewAccountSearchHandler(acmeCtx acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *acme.JWSCtx, data map[string]interface{}) (*logical.Response, error) {
	if userCtx.Existing || b.acme.DoesAccountExist(userCtx.Kid) {
		// This account exists; return its details. It would be slightly
		// weird to specify a kid in the request (and not use an explicit
		// jwk here), but we might as well support it too.
		account, err := b.acme.LoadAccount(userCtx.Kid)
		if err != nil {
			return nil, fmt.Errorf("error loading account: %w", err)
		}

		location := acmeCtx.baseUrl.String() + "/acme/account/" + userCtx.Kid
		return formatAccountResponse(location, account["status"].(string), account["contact"].([]string)), nil
	}

	// Per RFC 8555 Section 7.3.1. Finding an Account URL Given a Key:
	//
	// > If a client sends such a request and an account does not exist,
	// > then the server MUST return an error response with status code
	// > 400 (Bad Request) and type "urn:ietf:params:acme:error:accountDoesNotExist".
	return nil, fmt.Errorf("An account with this key does not exist: %w", acme.ErrAccountDoesNotExist)
}

func (b *backend) acmeNewAccountCreateHandler(acmeCtx acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *acme.JWSCtx, data map[string]interface{}, contact []string, termsOfServiceAgreed bool) (*logical.Response, error) {
	if userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to newAccount with 'kid': %w", acme.ErrMalformed)
	}

	// If the account already exists, return the existing one.
	if b.acme.DoesAccountExist(userCtx.Kid) {
		return b.acmeNewAccountSearchHandler(acmeCtx, r, fields, userCtx, data)
	}

	// TODO: Limit this only when ToS are required by the operator.
	if !termsOfServiceAgreed {
		return nil, fmt.Errorf("terms of service not agreed to: %w", acme.ErrUserActionRequired)
	}

	account, err := b.acme.CreateAccount(userCtx, contact, termsOfServiceAgreed)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	location := acmeCtx.baseUrl.String() + "/acme/account/" + userCtx.Kid
	return formatAccountResponse(location, account["status"].(string), account["contact"].([]string)), nil
}
