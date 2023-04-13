package pki

import (
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/net/idna"
)

func pathAcmeRootListOrders(b *backend) *framework.Path {
	return patternAcmeListOrders(b, "acme/orders")
}

func pathAcmeRoleListOrders(b *backend) *framework.Path {
	return patternAcmeListOrders(b, "roles/"+framework.GenericNameRegex("role")+"/acme/orders")
}

func pathAcmeIssuerListOrders(b *backend) *framework.Path {
	return patternAcmeListOrders(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/orders")
}

func pathAcmeIssuerAndRoleListOrders(b *backend) *framework.Path {
	return patternAcmeListOrders(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/orders")
}

func pathAcmeRootNewOrder(b *backend) *framework.Path {
	return patternAcmeNewOrder(b, "acme/new-order")
}

func pathAcmeRoleNewOrder(b *backend) *framework.Path {
	return patternAcmeNewOrder(b, "roles/"+framework.GenericNameRegex("role")+"/acme/new-order")
}

func pathAcmeIssuerNewOrder(b *backend) *framework.Path {
	return patternAcmeNewOrder(b, "issuer/"+framework.GenericNameRegex(issuerRefParam)+"/acme/new-order")
}

func pathAcmeIssuerAndRoleNewOrder(b *backend) *framework.Path {
	return patternAcmeNewOrder(b,
		"issuer/"+framework.GenericNameRegex(issuerRefParam)+
			"/roles/"+framework.GenericNameRegex("role")+"/acme/new-order")
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

type acmeAccountRequiredOperation func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, userCtx *jwsCtx, data map[string]interface{}, acct *acmeAccount) (*logical.Response, error)

func (b *backend) acmeAccountRequiredWrapper(op acmeAccountRequiredOperation) framework.OperationFunc {
	return b.acmeParsedWrapper(func(acmeCtx *acmeContext, r *logical.Request, fields *framework.FieldData, uc *jwsCtx, data map[string]interface{}) (*logical.Response, error) {
		if !uc.Existing {
			return nil, fmt.Errorf("cannot process request without a 'kid': %w", ErrMalformed)
		}

		account, err := b.acmeState.LoadAccount(acmeCtx, uc.Kid)
		if err != nil {
			return nil, fmt.Errorf("error loading account: %w", err)
		}

		if account.Status != StatusValid {
			// Treating "revoked" and "deactivated" as the same here.
			return nil, ErrUnauthorized
		}

		return op(acmeCtx, r, fields, uc, data, account)
	})
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

func (b *backend) acmeNewOrderHandler(ac *acmeContext, _ *logical.Request, _ *framework.FieldData, _ *jwsCtx, data map[string]interface{}, account *acmeAccount) (*logical.Response, error) {
	identifiers, err := parseOrderIdentifiers(data)
	if err != nil {
		return nil, err
	}

	notBefore, _, err := parseOptRFC3339Field(data, "notBefore")
	if err != nil {
		return nil, err
	}

	notAfter, _, err := parseOptRFC3339Field(data, "notAfter")
	if err != nil {
		return nil, err
	}

	// TODO DO more validations here, role can issue, notBefore/notAfter makes sense...

	var authorizations []*ACMEAuthorization
	var authorizationIds []string
	for _, identifier := range identifiers {
		authz := generateAuthorization(account, identifier)
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

func formatOrderResponse(acmeCtx *acmeContext, order *acmeOrder) *logical.Response {
	location := buildOrderUrl(acmeCtx, order.OrderId)

	var authorizationUrls []string
	for _, authId := range order.AuthorizationIds {
		authorizationUrls = append(authorizationUrls, acmeCtx.baseUrl.String()+"authz/"+authId)
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"status":         ACMEOrderPending,
			"expires":        order.Expires,
			"identifiers":    order.Identifiers,
			"authorizations": authorizationUrls,
			"finalize":       location + "/finalize",
		},
		Headers: map[string][]string{
			"Location": {location},
		},
	}

	return resp
}

func buildOrderUrl(acmeCtx *acmeContext, orderId string) string {
	return acmeCtx.baseUrl.String() + "order/" + orderId
}

func generateAuthorization(acct *acmeAccount, identifier *ACMEIdentifier) *ACMEAuthorization {
	challenges := []*ACMEChallenge{
		{
			Type:            ACMEHTTPChallenge,
			URL:             genUuid(),
			Status:          ACMEChallengePending,
			ChallengeFields: map[string]interface{}{}, // TODO fill this in properly
		},
	}

	return &ACMEAuthorization{
		Id:         genUuid(),
		AccountId:  acct.KeyId,
		Identifier: identifier,
		Status:     ACMEAuthorizationPending,
		Expires:    "", // only populated when it switches to valid.
		Challenges: challenges,
		Wildcard:   false,
	}
}

func parseOptRFC3339Field(data map[string]interface{}, keyName string) (time.Time, bool, error) {
	var timeVal time.Time
	var err error

	rawBefore, present := data[keyName]
	if present {
		beforeStr, ok := rawBefore.(string)
		if !ok {
			return timeVal, false, fmt.Errorf("invalid type (%T) for field '%s': %w", rawBefore, keyName, ErrMalformed)
		}
		timeVal, err = time.Parse(time.RFC3339, beforeStr)
		if err != nil {
			return timeVal, false, fmt.Errorf("failed parsing field '%s' (%s): %s: %w", keyName, rawBefore, err.Error(), ErrMalformed)
		}

		return timeVal, true, nil
	}

	return timeVal, false, nil
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
