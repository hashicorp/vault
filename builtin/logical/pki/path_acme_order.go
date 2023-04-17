// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package pki

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/net/idna"
)

func pathAcmeListOrders(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeListOrders, "/orders")
}

func pathAcmeGetOrder(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeGetOrder, "/order/"+uuidNameRegex("order_id"))
}

func pathAcmeNewOrder(b *backend) []*framework.Path {
	return buildAcmeFrameworkPaths(b, patternAcmeNewOrder, "/new-order")
}

func patternAcmeNewOrder(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeNewOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeListOrders(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeListOrdersHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func patternAcmeGetOrder(b *backend, pattern string) *framework.Path {
	fields := map[string]*framework.FieldSchema{}
	addFieldsForACMEPath(fields, pattern)
	addFieldsForACMERequest(fields)
	fields["order_id"] = &framework.FieldSchema{
		Type:        framework.TypeString,
		Description: `The ACME order identifier to fetch`,
		Required:    true,
	}

	return &framework.Path{
		Pattern: pattern,
		Fields:  fields,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback:                    b.acmeAccountRequiredWrapper(b.acmeGetOrderHandler),
				ForwardPerformanceSecondary: false,
				ForwardPerformanceStandby:   true,
			},
		},

		HelpSynopsis:    "",
		HelpDescription: "",
	}
}

func (b *backend) acmeGetOrderHandler(ac *acmeContext, _ *logical.Request, fields *framework.FieldData, uc *jwsCtx, _ map[string]interface{}, acct *acmeAccount) (*logical.Response, error) {
	orderId := fields.Get("order_id").(string)

	order, err := b.acmeState.LoadOrder(ac, uc, orderId)
	if err != nil {
		return nil, err
	}

	// Per RFC 8555 -> 7.1.3.  Order Objects
	// For final orders (in the "valid" or "invalid" state), the authorizations that were completed.
	//
	// Otherwise, for "pending" orders we will return our list as it was originally saved.
	requiresFiltering := order.Status == ACMEOrderValid || order.Status == ACMEOrderInvalid
	if requiresFiltering {
		filteredAuthorizationIds := []string{}

		for _, authId := range order.AuthorizationIds {
			authorization, err := b.acmeState.LoadAuthorization(ac, uc, authId)
			if err != nil {
				return nil, err
			}

			if (order.Status == ACMEOrderInvalid || order.Status == ACMEOrderValid) &&
				authorization.Status == ACMEAuthorizationValid {
				filteredAuthorizationIds = append(filteredAuthorizationIds, authId)
			}
		}

		order.AuthorizationIds = filteredAuthorizationIds
	}

	return formatOrderResponse(ac, order), nil
}

func (b *backend) acmeListOrdersHandler(ac *acmeContext, _ *logical.Request, _ *framework.FieldData, uc *jwsCtx, _ map[string]interface{}, acct *acmeAccount) (*logical.Response, error) {
	orderIds, err := b.acmeState.ListOrderIds(ac, acct.KeyId)
	if err != nil {
		return nil, err
	}

	orderUrls := []string{}
	for _, orderId := range orderIds {
		order, err := b.acmeState.LoadOrder(ac, uc, orderId)
		if err != nil {
			return nil, err
		}

		if order.Status == ACMEOrderInvalid {
			// Per RFC8555 -> 7.1.2.1 - Orders List
			// The server SHOULD include pending orders and SHOULD NOT
			// include orders that are invalid in the array of URLs.
			continue
		}

		orderUrls = append(orderUrls, buildOrderUrl(ac, orderId))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"orders": orderUrls,
		},
	}

	return resp, nil
}

func (b *backend) acmeNewOrderHandler(ac *acmeContext, r *logical.Request, _ *framework.FieldData, _ *jwsCtx, data map[string]interface{}, account *acmeAccount) (*logical.Response, error) {
	identifiers, err := parseOrderIdentifiers(data)
	if err != nil {
		return nil, err
	}

	notBefore, err := parseOptRFC3339Field(data, "notBefore")
	if err != nil {
		return nil, err
	}

	notAfter, err := parseOptRFC3339Field(data, "notAfter")
	if err != nil {
		return nil, err
	}

	err = validateAcmeProvidedOrderDates(notBefore, notAfter)
	if err != nil {
		return nil, err
	}

	// TODO: Implement checks against role here.

	// Per RFC 8555 -> 7.1.3. Order Objects
	// For pending orders, the authorizations that the client needs to complete before the
	// requested certificate can be issued (see Section 7.5), including
	// unexpired authorizations that the client has completed in the past
	// for identifiers specified in the order.
	//
	// Since we are generating all authorizations here, there is no need to filter them out
	// IF/WHEN we support pre-authz workflows and associate existing authorizations to this
	// order they will need filtering.
	var authorizations []*ACMEAuthorization
	var authorizationIds []string
	for _, identifier := range identifiers {
		authz, err := generateAuthorization(ac, account, identifier)
		if err != nil {
			return nil, fmt.Errorf("error generating authorizations: %w", err)
		}
		authorizations = append(authorizations, authz)

		err = b.acmeState.SaveAuthorization(ac, authz)
		if err != nil {
			return nil, fmt.Errorf("failed storing authorization: %w", err)
		}

		authorizationIds = append(authorizationIds, authz.Id)
	}

	order := &acmeOrder{
		OrderId:          genUuid(),
		AccountId:        account.KeyId,
		Status:           ACMEOrderPending,
		Expires:          time.Now().Add(1 * time.Hour).Format(time.RFC3339),
		Identifiers:      identifiers,
		AuthorizationIds: authorizationIds,
	}

	if !notBefore.IsZero() {
		order.NotBefore = notBefore.Format(time.RFC3339)
	}

	if !notAfter.IsZero() {
		order.NotAfter = notAfter.Format(time.RFC3339)
	}

	err = b.acmeState.SaveOrder(ac, order)
	if err != nil {
		return nil, fmt.Errorf("failed storing order: %w", err)
	}

	resp := formatOrderResponse(ac, order)

	// Per RFC 8555 Section 7.4. Applying for Certificate Issuance:
	//
	// > If the server is willing to issue the requested certificate, it
	// > responds with a 201 (Created) response.
	resp.Data[logical.HTTPStatusCode] = http.StatusCreated
	return resp, nil
}

func validateAcmeProvidedOrderDates(notBefore time.Time, notAfter time.Time) error {
	if !notBefore.IsZero() && !notAfter.IsZero() {
		if notBefore.Equal(notAfter) {
			return fmt.Errorf("%w: provided notBefore and notAfter dates can not be equal", ErrMalformed)
		}

		if notBefore.After(notAfter) {
			return fmt.Errorf("%w: provided notBefore can not be greater than notAfter", ErrMalformed)
		}
	}

	if !notAfter.IsZero() {
		if time.Now().After(notAfter) {
			return fmt.Errorf("%w: provided notAfter can not be in the past", ErrMalformed)
		}
	}

	return nil
}

func formatOrderResponse(acmeCtx *acmeContext, order *acmeOrder) *logical.Response {
	baseOrderUrl := buildOrderUrl(acmeCtx, order.OrderId)

	var authorizationUrls []string
	for _, authId := range order.AuthorizationIds {
		authorizationUrls = append(authorizationUrls, buildAuthorizationUrl(acmeCtx, authId))
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"status":         ACMEOrderPending,
			"expires":        order.Expires,
			"identifiers":    order.Identifiers,
			"authorizations": authorizationUrls,
			"finalize":       baseOrderUrl + "/finalize",
		},
		Headers: map[string][]string{
			"Location": {baseOrderUrl},
		},
	}

	// Only reply with the certificate URL if we are in a valid order state.
	if order.Status == ACMEOrderValid {
		resp.Data["certificate"] = baseOrderUrl + "/cert"
	}

	return resp
}

func buildAuthorizationUrl(acmeCtx *acmeContext, authId string) string {
	return acmeCtx.baseUrl.JoinPath("authorization", authId).String()
}

func buildOrderUrl(acmeCtx *acmeContext, orderId string) string {
	return acmeCtx.baseUrl.JoinPath("order", orderId).String()
}

func generateAuthorization(acmeCtx *acmeContext, acct *acmeAccount, identifier *ACMEIdentifier) (*ACMEAuthorization, error) {
	authId := genUuid()
	var challenges []*ACMEChallenge
	for _, challengeType := range []ACMEChallengeType{ACMEHTTPChallenge} {
		token, err := getACMEToken()
		if err != nil {
			return nil, err
		}

		challenge := &ACMEChallenge{
			Type:   challengeType,
			Status: ACMEChallengePending,
			ChallengeFields: map[string]interface{}{
				"token": token,
			},
		}

		challenges = append(challenges, challenge)
	}

	return &ACMEAuthorization{
		Id:         authId,
		AccountId:  acct.KeyId,
		Identifier: identifier,
		Status:     ACMEAuthorizationPending,
		Expires:    "", // only populated when it switches to valid.
		Challenges: challenges,
		Wildcard:   strings.HasPrefix(identifier.Value, "*."),
	}, nil
}

func parseOptRFC3339Field(data map[string]interface{}, keyName string) (time.Time, error) {
	var timeVal time.Time
	var err error

	rawBefore, present := data[keyName]
	if present {
		beforeStr, ok := rawBefore.(string)
		if !ok {
			return timeVal, fmt.Errorf("invalid type (%T) for field '%s': %w", rawBefore, keyName, ErrMalformed)
		}
		timeVal, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return timeVal, fmt.Errorf("failed parsing field '%s' (%s): %s: %w", keyName, rawBefore, err.Error(), ErrMalformed)
		}

		if timeVal.IsZero() {
			return timeVal, fmt.Errorf("provided time value is invalid '%s' (%s): %w", keyName, rawBefore, ErrMalformed)
		}
	}

	return timeVal, nil
}

func parseOrderIdentifiers(data map[string]interface{}) ([]*ACMEIdentifier, error) {
	rawIdentifiers, present := data["identifiers"]
	if !present {
		return nil, fmt.Errorf("missing required identifiers argument: %w", ErrMalformed)
	}

	listIdentifiers, ok := rawIdentifiers.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid type (%T) for field 'identifiers': %w", rawIdentifiers, ErrMalformed)
	}

	var identifiers []*ACMEIdentifier
	for _, rawIdentifier := range listIdentifiers {
		mapIdentifier, ok := rawIdentifier.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid type (%T) for value in 'identifiers': %w", rawIdentifier, ErrMalformed)
		}

		typeVal, present := mapIdentifier["type"]
		if !present {
			return nil, fmt.Errorf("missing type argument for value in 'identifiers': %w", ErrMalformed)
		}
		typeStr, ok := typeVal.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for type argument (%T) for value in 'identifiers': %w", typeStr, ErrMalformed)
		}

		var acmeIdentifierType ACMEIdentifierType
		switch typeStr {
		// TODO: No support for this yet.
		// case string(ACMEIPIdentifier):
		//	acmeIdentifierType = ACMEIPIdentifier
		case string(ACMEDNSIdentifier):
			acmeIdentifierType = ACMEDNSIdentifier
		default:
			return nil, fmt.Errorf("unsupported identifier type %s: %w", typeStr, ErrUnsupportedIdentifier)
		}

		valueVal, present := mapIdentifier["value"]
		if !present {
			return nil, fmt.Errorf("missing value argument for value in 'identifiers': %w", ErrMalformed)
		}
		valueStr, ok := valueVal.(string)
		if !ok {
			return nil, fmt.Errorf("invalid type for value argument (%T) for value in 'identifiers': %w", valueStr, ErrMalformed)
		}

		if len(valueStr) == 0 {
			return nil, fmt.Errorf("value argument for value in 'identifiers' can not be blank: %w", ErrMalformed)
		}

		p := idna.New(
			idna.StrictDomainName(true),
			idna.VerifyDNSLength(true),
		)
		converted, err := p.ToASCII(valueStr)
		if err != nil {
			return nil, fmt.Errorf("value argument (%s) failed validation: %s: %w", valueStr, err.Error(), ErrMalformed)
		}
		if !hostnameRegex.MatchString(converted) {
			return nil, fmt.Errorf("value argument (%s) failed validation: %w", valueStr, ErrMalformed)
		}

		identifiers = append(identifiers, &ACMEIdentifier{
			Type:  acmeIdentifierType,
			Value: converted,
		})
	}

	return identifiers, nil
}
