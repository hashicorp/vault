package pki

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hashicorp/go-secure-stdlib/strutil"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func uuidNameRegex(name string) string {
	return fmt.Sprintf("(?P<%s>[[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}?)", name)
}

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

func pathAcmeRootUpdateAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "acme/account/"+uuidNameRegex("kid"))
}

func pathAcmeRoleUpdateAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "roles/"+framework.GenericNameRegex("role")+"/acme/account/"+uuidNameRegex("kid"))
}

func pathAcmeIssuerUpdateAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/account/"+uuidNameRegex("kid"))
}

func pathAcmeIssuerAndRoleUpdateAccount(b *backend) *framework.Path {
	return patternAcmeNewAccount(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/account/"+uuidNameRegex("kid"))
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

func addFieldsForACMEKidRequest(fields map[string]*framework.FieldSchema, pattern string) map[string]*framework.FieldSchema {
	if strings.Contains(pattern, framework.GenericNameRegex("kid")) {
		fields["kid"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The key identifier provided by the CA`,
			Required:    true,
		}
	}

	return fields
}

func patternAcmeNewAccount(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEKidRequest(fields, pattern)

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
	var status string

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

	// Per RFC 8555 7.3.6 Account deactivation, we will handle it within our update API.
	rawStatus, present := data["status"]
	if present {
		status, ok = rawStatus.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for field 'onlyReturnExisting': %w", rawOnlyReturnExisting, ErrMalformed)
		}
	}

	// We ignore the EAB parameter as it is currently not supported.

	// We have two paths here: search or create.
	if onlyReturnExisting {
		return b.acmeNewAccountSearchHandler(acmeCtx, userCtx)
	}

	// Pass through the /new-account API calls to this specific handler as its requirements are different
	// from the account update handler.
	if strings.HasSuffix(r.Path, "/new-account") {
		return b.acmeNewAccountCreateHandler(acmeCtx, userCtx, contacts, termsOfServiceAgreed)
	}

	return b.acmeNewAccountUpdateHandler(acmeCtx, userCtx, contacts, status)
}

func formatAccountResponse(acmeCtx *acmeContext, acct *acmeAccount) *logical.Response {
	location := acmeCtx.baseUrl.String() + "account/" + acct.KeyId

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

func (b *backend) acmeNewAccountSearchHandler(acmeCtx *acmeContext, userCtx *jwsCtx) (*logical.Response, error) {
	thumbprint, err := userCtx.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint for key: %w", err)
	}

	account, err := b.acmeState.LoadAccountByKey(acmeCtx, thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed to load account by thumbprint: %w", err)
	}

	if account != nil {
		return formatAccountResponse(acmeCtx, account), nil
	}

	// Per RFC 8555 Section 7.3.1. Finding an Account URL Given a Key:
	//
	// > If a client sends such a request and an account does not exist,
	// > then the server MUST return an error response with status code
	// > 400 (Bad Request) and type "urn:ietf:params:acme:error:accountDoesNotExist".
	return nil, fmt.Errorf("An account with this key does not exist: %w", ErrAccountDoesNotExist)
}

func (b *backend) acmeNewAccountCreateHandler(acmeCtx *acmeContext, userCtx *jwsCtx, contact []string, termsOfServiceAgreed bool) (*logical.Response, error) {
	if userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to newAccount with 'kid': %w", ErrMalformed)
	}

	// If the account already exists, return the existing one.
	thumbprint, err := userCtx.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint for key: %w", err)
	}

	accountByKey, err := b.acmeState.LoadAccountByKey(acmeCtx, thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed to load account by thumbprint: %w", err)
	}

	if accountByKey != nil {
		return formatAccountResponse(acmeCtx, accountByKey), nil
	}

	// TODO: Limit this only when ToS are required or set by the operator, since we don't have a
	//       ToS URL in the directory at the moment, we can not enforce this.
	//if !termsOfServiceAgreed {
	//	return nil, fmt.Errorf("terms of service not agreed to: %w", ErrUserActionRequired)
	//}

	accountByKid, err := b.acmeState.CreateAccount(acmeCtx, userCtx, contact, termsOfServiceAgreed)
	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	resp := formatAccountResponse(acmeCtx, accountByKid)

	// Per RFC 8555 Section 7.3. Account Management:
	//
	// > The server returns this account object in a 201 (Created) response,
	// > with the account URL in a Location header field.
	resp.Data[logical.HTTPStatusCode] = http.StatusCreated
	return resp, nil
}

func (b *backend) acmeNewAccountUpdateHandler(acmeCtx *acmeContext, userCtx *jwsCtx, contact []string, status string) (*logical.Response, error) {
	if !userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to account updates without a 'kid': %w", ErrMalformed)
	}

	account, err := b.acmeState.LoadAccount(acmeCtx, userCtx.Kid)
	if err != nil {
		return nil, fmt.Errorf("error loading account: %w", err)
	}

	// Per RFC 8555 7.3.6 Account deactivation, if we were previously deactivated, we should return
	// unauthorized. There is no way to reactivate any accounts per ACME RFC.
	if account.Status != StatusValid {
		// Treating "revoked" and "deactivated" as the same here.
		return nil, ErrUnauthorized
	}

	shouldUpdate := false
	// Check to see if we should update, we don't really care about ordering
	if !strutil.EquivalentSlices(account.Contact, contact) {
		shouldUpdate = true
		account.Contact = contact
	}

	// Check to process account de-activation status was requested.
	// 7.3.6. Account Deactivation
	if string(StatusDeactivated) == status {
		shouldUpdate = true
		// TODO: This should cancel any ongoing operations (do not revoke certs),
		//       perhaps we should delete this account here?
		account.Status = StatusDeactivated
	}

	if shouldUpdate {
		err = b.acmeState.UpdateAccount(acmeCtx, account)
		if err != nil {
			return nil, fmt.Errorf("failed to update account: %w", err)
		}
	}

	resp := formatAccountResponse(acmeCtx, account)
	return resp, nil
}
