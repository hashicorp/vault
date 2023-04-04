package pki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathAcmeRootNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "acme/new-account")
}

func pathAcmeRoleNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "roles/"+framework.GenericNameRegex("role")+"/acme/new-account")
}

func pathAcmeIssuerNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/new-account")
}

func pathAcmeIssuerAndRoleNewAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/new-account")
}

func addFieldsForACMEPath(fields map[string]*framework.FieldSchema, pattern string) map[string]*framework.FieldSchema {
	if strings.Contains(pattern, framework.GenericNameRegex("role")) {
		fields["role"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The desired role for the acme request`,
			Required:    true,
		}
	}
	if strings.Contains(pattern, framework.GenericNameRegex(issuerRefParam)) {
		fields[issuerRefParam] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `Reference to an existing issuer name or issuer id`,
			Required:    true,
		}
	}

	return fields
}

func addFieldsForACMERequest(fields map[string]*framework.FieldSchema) map[string]*framework.FieldSchema {
	fields["protected"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME request 'protected' value",
		Required:    false,
	}

	fields["payload"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME request 'payload' value",
		Required:    false,
	}

	fields["signature"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: "ACME request 'signature' value",
		Required:    false,
	}

	return fields
}

func patternAcmeNewAccount(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeParsedWrapper(b.acmeNewAccountHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathOcspHelpSyn,
		HelpDescription: pathOcspHelpDesc,
	}
}

type acmeParsedOperation func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error)

func (b *backend) acmeParsedWrapper(op acmeParsedOperation) framework.OperationFunc {
	return b.acmeWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData) (*logical.Response, error) {
		user, data, err := b.acmeState.ParseRequestParams(acmeCtx, fields)
		if err != nil {
			return nil, err
		}

		resp, err := op(acmeCtx, r, fields, user, data)

		// Our response handlers might not add the necessary headers.
		if resp != nil {
			if resp.Headers == nil {
				resp.Headers = map[string][]string{}
			}

			if _, ok := resp.Headers["Replay-Nonce"]; !ok {
				nonce, _, err := b.acmeState.GetNonce()
				if err != nil {
					return nil, err
				}

				resp.Headers["Replay-Nonce"] = []string{nonce}
			}

			if _, ok := resp.Headers["Link"]; !ok {
				resp.Headers["Link"] = genAcmeLinkHeader(acmeCtx)
			} else {
				directory := genAcmeLinkHeader(acmeCtx)[0]
				addDirectory := true
				for _, item := range resp.Headers["Link"] {
					if item == directory {
						addDirectory = false
						break
					}
				}
				if addDirectory {
					resp.Headers["Link"] = append(resp.Headers["Link"], directory)
				}
			}

			// ACME responses don't understand Vault's default encoding
			// format. Rather than expecting everything to handle creating
			// ACME-formatted responses, do the marshaling in one place.
			if _, ok := resp.Data[logical.HTTPRawBody]; !ok {
				ignored_values := map[string]bool{logical.HTTPContentType: true, logical.HTTPStatusCode: true}
				fields := map[string]interface{}{}
				body := map[string]interface{}{
					logical.HTTPContentType: "application/json",
					logical.HTTPStatusCode:  http.StatusOK,
				}

				for key, value := range resp.Data {
					if _, present := ignored_values[key]; !present {
						fields[key] = value
					} else {
						body[key] = value
					}
				}

				rawBody, err := json.Marshal(fields)
				if err != nil {
					return nil, fmt.Errorf("Error marshaling JSON body: %w", err)
				}

				body[logical.HTTPRawBody] = rawBody
				resp.Data = body
			}
		}

		return resp, err
	})
}

func (b *backend) acmeNewAccountHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
	// Parameters
	var ok bool
	var onlyReturnExisting bool
	var contacts []string
	var termsOfServiceAgreed bool

	rawContact, present := data["contact"]
	if present {
		listContact, ok := rawContact.([]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for field 'contact': %w", rawContact, ErrMalformed)
		}

		for index, singleContact := range listContact {
			contact, ok := singleContact.(string)
			if !ok {
				return nil, fmt.Errorf("invalid type (%T) for field 'contact' item %d: %w", singleContact, index, ErrMalformed)
			}

			contacts = append(contacts, contact)
		}
	}

	rawTermsOfServiceAgreed, present := data["termsOfServiceAgreed"]
	if present {
		termsOfServiceAgreed, ok = rawTermsOfServiceAgreed.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for field 'termsOfServiceAgreed': %w", rawTermsOfServiceAgreed, ErrMalformed)
		}
	}

	rawOnlyReturnExisting, present := data["onlyReturnExisting"]
	if present {
		onlyReturnExisting, ok = rawOnlyReturnExisting.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for field 'onlyReturnExisting': %w", rawOnlyReturnExisting, ErrMalformed)
		}
	}

	// We ignore the EAB parameter as it is currently not supported.

	// We have two paths here: search or create.
	if onlyReturnExisting {
		return b.acmeNewAccountSearchHandler(acmeCtx, r, fields, userCtx, data)
	}

	return b.acmeNewAccountCreateHandler(acmeCtx, r, fields, userCtx, data, contacts, termsOfServiceAgreed)
}

func formatAccountResponse(location string, acct *acmeAccount) *logical.Response {
	resp := &logical.Response{
		Data: map[string]interface{}{
			"status": acct.Status,
			"orders": location + "/orders",
		},
		Headers: map[string][]string{
			"Location": {location},
		},
	}

	if len(acct.Contact) > 0 {
		resp.Data["contact"] = acct.Contact
	}

	return resp
}

func (b *backend) acmeNewAccountSearchHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
	if userCtx.Existing || b.acmeState.DoesAccountExist(acmeCtx, userCtx.Kid) {
		// This account exists; return its details. It would be slightly
		// weird to specify a kid in the request (and not use an explicit
		// jwk here), but we might as well support it too.
		account, err := b.acmeState.LoadAccount(acmeCtx, userCtx.Kid)
		if err != nil {
			return nil, fmt.Errorf("error loading account: %w", err)
		}

		location := acmeCtx.baseUrl.String() + "account/" + userCtx.Kid
		return formatAccountResponse(location, account), nil
	}

	// Per RFC 8555 Section 7.3.1. Finding an Account URL Given a Key:
	//
	// > If a client sends such a request and an account does not exist,
	// > then the server MUST return an error response with status code
	// > 400 (Bad Request) and type "urn:ietf:params:acme:error:accountDoesNotExist".
	return nil, fmt.Errorf("An account with this key does not exist: %w", ErrAccountDoesNotExist)
}

func (b *backend) acmeNewAccountCreateHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, contact []string, termsOfServiceAgreed bool) (*logical.Response, error) {
	if userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to newAccount with 'kid': %w", ErrMalformed)
	}

	// If the account already exists, return the existing one.
	if b.acmeState.DoesAccountExist(acmeCtx, userCtx.Kid) {
		return b.acmeNewAccountSearchHandler(acmeCtx, r, fields, userCtx, data)
	}

	// TODO: Limit this only when ToS are required by the operator.
	if !termsOfServiceAgreed {
		return nil, fmt.Errorf("terms of service not agreed to: %w", ErrUserActionRequired)
	}

	account, err := b.acmeState.CreateAccount(acmeCtx, userCtx, contact, termsOfServiceAgreed)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	location := acmeCtx.baseUrl.String() + "account/" + userCtx.Kid
	resp := formatAccountResponse(location, account)

	// Per RFC 8555 Section 7.3. Account Management:
	//
	// > The server returns this account object in a 201 (Created) response,
	// > with the account URL in a Location header field.
	resp.Data[logical.HTTPStatusCode] = http.StatusCreated
	return resp, nil
}
