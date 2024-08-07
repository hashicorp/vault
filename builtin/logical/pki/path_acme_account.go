// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

func uuidNameRegex(name string) string {
	return fmt.Sprintf("(?P<%s>[[:alnum:]]{8}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{4}-[[:alnum:]]{12}?)", name)
}

func pathAcmeNewAccount(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeNewAccount(b, baseUrl+"/new-account", opts)
}

func pathAcmeUpdateAccount(b *backend, baseUrl string, opts acmeWrapperOpts) *framework.Path {
	return patternAcmeNewAccount(b, baseUrl+"/account/"+uuidNameRegex("kid"), opts)
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
	if strings.Contains(pattern, framework.GenericNameRegex("policy")) {
		fields["policy"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The policy name to pass through to the CIEPS service`,
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
	if strings.Contains(pattern, uuidNameRegex("kid")) {
		fields["kid"] = &framework.FieldSchema{
			Type:        framework.TypeString,
			Description: `The key identifier provided by the CA`,
			Required:    true,
		}
	}

	return fields
}

func patternAcmeNewAccount(b *backend, pattern string, opts acmeWrapperOpts) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	addFieldsForACMEKidRequest(fields, pattern)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeParsedWrapper(opts, b.acmeNewAccountHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    pathAcmeHelpSync,
		HelpDescription: pathAcmeHelpDesc,
	}
}

func (b *backend) acmeNewAccountHandler(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
	// Parameters
	var ok bool
	var onlyReturnExisting bool
	var contacts []string
	var termsOfServiceAgreed bool
	var status string
	var eabData map[string]interface{}

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

	if eabDataRaw, ok := data["externalAccountBinding"]; ok {
		eabData, ok = eabDataRaw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: externalAccountBinding field was unparseable", ErrMalformed)
		}
	}

	// We have two paths here: search or create.
	if onlyReturnExisting {
		return b.acmeAccountSearchHandler(acmeCtx, userCtx)
	}

	// Pass through the /new-account API calls to this specific handler as its requirements are different
	// from the account update handler.
	if strings.HasSuffix(r.Path, "/new-account") {
		return b.acmeNewAccountCreateHandler(acmeCtx, userCtx, contacts, termsOfServiceAgreed, r, eabData)
	}

	return b.acmeNewAccountUpdateHandler(acmeCtx, userCtx, contacts, status, eabData)
}

func formatNewAccountResponse(acmeCtx *acmeContext, acct *acmeAccount, eabData map[string]interface{}) *logical.Response {
	resp := formatAccountResponse(acmeCtx, acct)

	// Per RFC 8555 Section 7.1.2.  Account Objects
	// Including this field in a newAccount request indicates approval by
	// the holder of an existing non-ACME account to bind that account to
	// this ACME account
	if acct.Eab != nil && len(eabData) != 0 {
		resp.Data["externalAccountBinding"] = eabData
	}

	return resp
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

func (b *backend) acmeAccountSearchHandler(acmeCtx *acmeContext, userCtx *jwsCtx) (*logical.Response, error) {
	thumbprint, err := userCtx.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint for key: %w", err)
	}

	account, err := b.GetAcmeState().LoadAccountByKey(acmeCtx, thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed to load account by thumbprint: %w", err)
	}

	if account != nil {
		if err = acmeCtx.eabPolicy.EnforceForExistingAccount(account); err != nil {
			return nil, err
		}
		return formatAccountResponse(acmeCtx, account), nil
	}

	// Per RFC 8555 Section 7.3.1. Finding an Account URL Given a Key:
	//
	// > If a client sends such a request and an account does not exist,
	// > then the server MUST return an error response with status code
	// > 400 (Bad Request) and type "urn:ietf:params:acme:error:accountDoesNotExist".
	return nil, fmt.Errorf("An account with this key does not exist: %w", ErrAccountDoesNotExist)
}

func (b *backend) acmeNewAccountCreateHandler(acmeCtx *acmeContext, userCtx *jwsCtx, contact []string, termsOfServiceAgreed bool, r *logical.Request, eabData map[string]interface{}) (*logical.Response, error) {
	if userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to newAccount with 'kid': %w", ErrMalformed)
	}

	// If the account already exists, return the existing one.
	thumbprint, err := userCtx.GetKeyThumbprint()
	if err != nil {
		return nil, fmt.Errorf("failed generating thumbprint for key: %w", err)
	}

	accountByKey, err := b.GetAcmeState().LoadAccountByKey(acmeCtx, thumbprint)
	if err != nil {
		return nil, fmt.Errorf("failed to load account by thumbprint: %w", err)
	}

	if accountByKey != nil {
		if err = acmeCtx.eabPolicy.EnforceForExistingAccount(accountByKey); err != nil {
			return nil, err
		}
		return formatAccountResponse(acmeCtx, accountByKey), nil
	}

	var eab *eabType
	if len(eabData) != 0 {
		eab, err = verifyEabPayload(b.GetAcmeState(), acmeCtx, userCtx, r.Path, eabData)
		if err != nil {
			return nil, err
		}
	}

	// Verify against our EAB policy
	if err = acmeCtx.eabPolicy.EnforceForNewAccount(eab); err != nil {
		return nil, err
	}

	// TODO: Limit this only when ToS are required or set by the operator, since we don't have a
	//       ToS URL in the directory at the moment, we can not enforce this.
	//if !termsOfServiceAgreed {
	//	return nil, fmt.Errorf("terms of service not agreed to: %w", ErrUserActionRequired)
	//}

	if eab != nil {
		// We delete the EAB to prevent future re-use after associating it with an account, worst
		// case if we fail creating the account we simply nuked the EAB which they can create another
		// and retry
		wasDeleted, err := b.GetAcmeState().DeleteEab(acmeCtx.sc, eab.KeyID)
		if err != nil {
			return nil, fmt.Errorf("failed to delete eab reference: %w", err)
		}

		if !wasDeleted {
			// Something consumed our EAB before we did bail...
			return nil, fmt.Errorf("eab was already used: %w", ErrUnauthorized)
		}
	}

	b.acmeAccountLock.RLock() // Prevents Account Creation and Tidy Interfering
	defer b.acmeAccountLock.RUnlock()

	accountByKid, err := b.GetAcmeState().CreateAccount(acmeCtx, userCtx, contact, termsOfServiceAgreed, eab)
	if err != nil {
		if eab != nil {
			return nil, fmt.Errorf("failed to create account: %w; the EAB key used for this request has been deleted as a result of this operation; fetch a new EAB key before retrying", err)
		}
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	resp := formatNewAccountResponse(acmeCtx, accountByKid, eabData)

	// Per RFC 8555 Section 7.3. Account Management:
	//
	// > The server returns this account object in a 201 (Created) response,
	// > with the account URL in a Location header field.
	resp.Data[logical.HTTPStatusCode] = http.StatusCreated
	return resp, nil
}

func (b *backend) acmeNewAccountUpdateHandler(acmeCtx *acmeContext, userCtx *jwsCtx, contact []string, status string, eabData map[string]interface{}) (*logical.Response, error) {
	if !userCtx.Existing {
		return nil, fmt.Errorf("cannot submit to account updates without a 'kid': %w", ErrMalformed)
	}

	if len(eabData) != 0 {
		return nil, fmt.Errorf("%w: not allowed to update EAB data in accounts", ErrMalformed)
	}

	account, err := b.GetAcmeState().LoadAccount(acmeCtx, userCtx.Kid)
	if err != nil {
		return nil, fmt.Errorf("error loading account: %w", err)
	}

	if err = acmeCtx.eabPolicy.EnforceForExistingAccount(account); err != nil {
		return nil, err
	}

	// Per RFC 8555 7.3.6 Account deactivation, if we were previously deactivated, we should return
	// unauthorized. There is no way to reactivate any accounts per ACME RFC.
	if account.Status != AccountStatusValid {
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
	if string(AccountStatusDeactivated) == status {
		shouldUpdate = true
		// TODO: This should cancel any ongoing operations (do not revoke certs),
		//       perhaps we should delete this account here?
		account.Status = AccountStatusDeactivated
		account.AccountRevokedDate = time.Now()
	}

	if shouldUpdate {
		err = b.GetAcmeState().UpdateAccount(acmeCtx.sc, account)
		if err != nil {
			return nil, fmt.Errorf("failed to update account: %w", err)
		}
	}

	resp := formatAccountResponse(acmeCtx, account)
	return resp, nil
}

func (b *backend) tidyAcmeAccountByThumbprint(as *acmeState, sc *storageContext, keyThumbprint string, certTidyBuffer, accountTidyBuffer time.Duration) error {
	thumbprintEntry, err := sc.Storage.Get(sc.Context, path.Join(acmeThumbprintPrefix, keyThumbprint))
	if err != nil {
		return fmt.Errorf("error retrieving thumbprint entry %v, unable to find corresponding account entry: %w", keyThumbprint, err)
	}
	if thumbprintEntry == nil {
		return fmt.Errorf("empty thumbprint entry %v, unable to find corresponding account entry", keyThumbprint)
	}

	var thumbprint acmeThumbprint
	err = thumbprintEntry.DecodeJSON(&thumbprint)
	if err != nil {
		return fmt.Errorf("unable to decode thumbprint entry %v to find account entry: %w", keyThumbprint, err)
	}

	if len(thumbprint.Kid) == 0 {
		return fmt.Errorf("unable to find account entry: empty kid within thumbprint entry: %s", keyThumbprint)
	}

	// Now Get the Account:
	accountEntry, err := sc.Storage.Get(sc.Context, acmeAccountPrefix+thumbprint.Kid)
	if err != nil {
		return err
	}
	if accountEntry == nil {
		// We delete the Thumbprint Associated with the Account, and we are done
		err = sc.Storage.Delete(sc.Context, path.Join(acmeThumbprintPrefix, keyThumbprint))
		if err != nil {
			return err
		}
		b.tidyStatusIncDeletedAcmeAccountCount()
		return nil
	}

	var account acmeAccount
	err = accountEntry.DecodeJSON(&account)
	if err != nil {
		return err
	}
	account.KeyId = thumbprint.Kid

	// Tidy Orders On the Account
	orderIds, err := as.ListOrderIds(sc, thumbprint.Kid)
	if err != nil {
		return err
	}
	allOrdersTidied := true
	maxCertExpiryUpdated := false
	for _, orderId := range orderIds {
		wasTidied, orderExpiry, err := b.acmeTidyOrder(sc, thumbprint.Kid, getOrderPath(thumbprint.Kid, orderId), certTidyBuffer)
		if err != nil {
			return err
		}
		if !wasTidied {
			allOrdersTidied = false
		}

		if !orderExpiry.IsZero() && account.MaxCertExpiry.Before(orderExpiry) {
			account.MaxCertExpiry = orderExpiry
			maxCertExpiryUpdated = true
		}
	}

	now := time.Now()
	if allOrdersTidied &&
		now.After(account.AccountCreatedDate.Add(accountTidyBuffer)) &&
		now.After(account.MaxCertExpiry.Add(accountTidyBuffer)) {
		// Tidy this account
		// If it is Revoked or Deactivated:
		if (account.Status == AccountStatusRevoked || account.Status == AccountStatusDeactivated) && now.After(account.AccountRevokedDate.Add(accountTidyBuffer)) {
			// We Delete the Account Associated with this Thumbprint:
			err = sc.Storage.Delete(sc.Context, path.Join(acmeAccountPrefix, thumbprint.Kid))
			if err != nil {
				return err
			}

			// Now we delete the Thumbprint Associated with the Account:
			err = sc.Storage.Delete(sc.Context, path.Join(acmeThumbprintPrefix, keyThumbprint))
			if err != nil {
				return err
			}
			b.tidyStatusIncDeletedAcmeAccountCount()
		} else if account.Status == AccountStatusValid {
			// Revoke This Account
			account.AccountRevokedDate = now
			account.Status = AccountStatusRevoked
			err := as.UpdateAccount(sc, &account)
			if err != nil {
				return err
			}
			b.tidyStatusIncRevAcmeAccountCount()
		}
	}

	// Only update the account if we modified the max cert expiry values and the account is still valid,
	// to prevent us from adding back a deleted account or not re-writing the revoked account that was
	// already written above.
	if maxCertExpiryUpdated && account.Status == AccountStatusValid {
		// Update our expiry time we previously setup.
		err := as.UpdateAccount(sc, &account)
		if err != nil {
			return err
		}
	}

	return nil
}
