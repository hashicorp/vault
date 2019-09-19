// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"unicode"
)

// These constants store the required http path & method information for validating the signed request
const (
	PathVersionBase = "/v1"
	PathBaseFormat  = "/auth/oci/login/%s"
	PathLoginMethod = "get"
)

// Signing Header constants
const (
	// HdrRequestTarget represents the special header name used to refer to the HTTP verb and URI in the signature.
	HdrRequestTarget = `(request-target)`
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"request_headers": {
				Type:        framework.TypeHeader,
				Description: `The signed headers of the client`,
			},
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role.",
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation: b.pathLoginUpdate,
		},

		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	// Validate the role
	role, ok := data.GetOk("role")
	if !ok {
		return logical.ErrorResponse("Role is not specified"), nil
	}
	roleName := role.(string)

	b.Logger().Trace(req.ID, "pathLoginUpdate roleName", roleName)

	// Validate that the role exists
	roleEntry, err := b.getOCIRole(ctx, req.Storage, roleName)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	if roleEntry == nil {
		return unauthorizedLogicalResponse(req, b.Logger(), fmt.Errorf("Role is not found"))
	}

	// Parse the authentication headers
	requestHeaders := data.Get("request_headers")
	if !ok {
		return logical.ErrorResponse("request_headers is not specified"), nil
	}
	authenticateRequestHeaders := requestHeaders.(http.Header)

	// Find the targetUrl and Method
	finalLoginPath := PathVersionBase + fmt.Sprintf(PathBaseFormat, roleName)
	method, targetUrl, err := requestTargetToMethodURL(authenticateRequestHeaders[HdrRequestTarget], PathLoginMethod, finalLoginPath)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}
	b.Logger().Trace(req.ID, "Method:", method, "targetUrl:", targetUrl)

	authenticateClientDetails := AuthenticateClientDetails{
		RequestHeaders: authenticateRequestHeaders,
	}

	requestMetadata := common.RequestMetadata{
		nil,
	}

	authenticateClientRequest := AuthenticateClientRequest{
		authenticateClientDetails,
		nil,
		&req.ID,
		requestMetadata,
	}

	// Authenticate the request with Identity
	if b.authenticationClient == nil && b.createAuthClient() != nil {
		return logical.RespondWithStatusCode(nil, req, http.StatusInternalServerError)
	}
	authenticateClientResponse, err := b.authenticationClient.AuthenticateClient(ctx, authenticateClientRequest)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}
	if authenticateClientResponse.Principal == nil ||
		len(authenticateClientResponse.Principal.Claims) == 0 ||
		*authenticateClientResponse.IsSuccess == false {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}
	internalClaims := FromClaims(authenticateClientResponse.Principal.Claims)
	principalType := internalClaims.GetString(ClaimPrincipalType)

	// Check the principal type
	if principalType != PrincipalTypeInstance && principalType != PrincipalTypeUser {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	b.Logger().Trace("Authentication ok", "Method:", method, "targetUrl:", targetUrl, "id", req.ID)

	// Validate the home tenancy
	err = b.validateHomeTenancy(ctx, req, *authenticateClientResponse.Principal.TenantId)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	// Find whether the entity corresponding the Principal is a part of any OCIDs allowed to take the role
	filterGroupMembershipDetails := FilterGroupMembershipDetails{
		*authenticateClientResponse.Principal,
		roleEntry.OcidList,
	}

	filterGroupMembershipRequest := FilterGroupMembershipRequest{
		filterGroupMembershipDetails,
		nil,
		&req.ID,
		requestMetadata,
	}

	filterGroupMembershipResponse, err := b.authenticationClient.FilterGroupMembership(ctx, filterGroupMembershipRequest)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}
	if filterGroupMembershipResponse.GroupIds == nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	// Validate that the filtered list contains atleast one of the OCIDs of the Role
	filteredOcidMap := sliceToMap(filterGroupMembershipResponse.GroupIds)
	found := false
	for _, item := range roleEntry.OcidList {
		_, present := filteredOcidMap[item]
		if present {
			found = true
			break
		}
	}
	if found == false {
		return unauthorizedLogicalResponse(req, b.Logger(), fmt.Errorf("Entity not a part of any of the Role OCIDs"))
	}

	b.Logger().Trace("Login ok", "Method:", method, "targetUrl:", targetUrl, "id", req.ID)

	// Return the response
	auth := &logical.Auth{
		Metadata: map[string]string{
			"role_name": roleName,
		},
		InternalData: map[string]interface{}{
			"role_name": roleName,
		},
		DisplayName: roleName,
		Alias: &logical.Alias{
			Name: "name",
		},
	}

	roleEntry.PopulateTokenAuth(auth)
	auth.Renewable = false

	resp := &logical.Response{
		Auth: auth,
	}

	return resp, nil
}

func (b *backend) validateHomeTenancy(ctx context.Context, req *logical.Request, homeTenancyId string) error {

	configEntry, err := b.getOCIConfig(ctx, req.Storage)
	if err != nil {
		return err
	}

	if configEntry == nil || configEntry.HomeTenancyId == "" {
		return fmt.Errorf("Home Tenancy is invalid")
	}

	if homeTenancyId != configEntry.HomeTenancyId {
		return fmt.Errorf("Invalid Tenancy")
	}

	return nil
}

func unauthorizedLogicalResponse(req *logical.Request, logger log.Logger, err error) (*logical.Response, error) {
	logger.Trace(req.ID, ": Failed with error:", err)
	return logical.RespondWithStatusCode(nil, req, http.StatusUnauthorized)
}

func requestTargetToMethodURL(requestTarget []string, expectedMethod string, expectedUrl string) (method string, url string, err error) {
	if len(requestTarget) == 0 {
		return "", "", errors.New("no (request-target) specified in header")
	}
	parts := strings.FieldsFunc(requestTarget[0], unicode.IsSpace)
	if len(parts) != 2 || strings.ToLower(parts[0]) != expectedMethod || strings.ToLower(parts[1]) != expectedUrl {
		return "", "", errors.New("incorrect (request-target) specified in header")
	}
	return parts[0], parts[1], nil
}

const pathLoginSyn = `
Authenticates to Vault using OCI credentials
`

const pathLoginDesc = `
Authenticates to Vault using OCI credentials such as User Api Key, Instance Principal
`
