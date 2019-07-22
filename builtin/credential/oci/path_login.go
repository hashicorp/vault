// Copyright © 2019, Oracle and/or its affiliates.
package ociauth

import (
	"context"
	"fmt"
	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"time"
	"unicode"
)

const PATH_VERSION_BASE = "/v1"
const PATH_BASE_FORMAT = "/auth/oci/login/%s"
const PATH_LOGIN_METHOD = "get"

// Header constants
const (
	// HdrRequestTarget represents the special header name used to refer to the HTTP verb and URI in the signature.
	HdrRequestTarget = `(request-target)`
	HdrHost = `host`
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login/" + framework.GenericNameRegex("role"),
		Fields: map[string]*framework.FieldSchema{
			"requestHeaders": {
				Type:        framework.TypeMap,
				Description: `The signed headers of the client`,
			},
			"role": {
				Type:        framework.TypeString,
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

	//Validate the role
	role, ok := data.GetOk("role")
	if !ok {
		return unauthorizedLogicalResponse(req, b.Logger(), fmt.Errorf("role is not specified"))
	}
	roleName := role.(string)
	if roleName != strings.ToLower(roleName) { //sanity check to prevent early exit when roleName case is mismatched, to prevent auth verification headers later-on
		return unauthorizedLogicalResponse(req, b.Logger(), fmt.Errorf("role is not in lower case"))
	}

	b.Logger().Debug(req.ID, "pathLoginUpdate roleName", roleName)

	//Parse the authentication headers
	authenticateRequestHeaders, err := deserializeRequest(data, b.Logger())
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	//Find the targetUrl and Method
	finalLoginPath := PATH_VERSION_BASE + fmt.Sprintf(PATH_BASE_FORMAT, roleName)
	method, targetUrl, err := requestTargetToMethodURL(authenticateRequestHeaders[HdrRequestTarget], PATH_LOGIN_METHOD, finalLoginPath)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}
	b.Logger().Debug(req.ID, "Method:", method, "targetUrl:", targetUrl)

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

	//Authenticate the request with Identity
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

	//Check the principal type
	if principalType != PrincipalTypeInstance && principalType != PrincipalTypeUser {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	b.Logger().Debug("Authentication ok", "Method:", method, "targetUrl:", targetUrl, "id", req.ID)

	//Validate the home tenancy
	err = b.validateHomeTenancy(ctx, req, *authenticateClientResponse.Principal.TenantId)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	//Validate that the role exists
	roleEntry, err := b.nonLockedOCIRole(ctx, req.Storage, roleName)
	if err != nil {
		return unauthorizedLogicalResponse(req, b.Logger(), err)
	}

	if roleEntry == nil {
		return unauthorizedLogicalResponse(req, b.Logger(), fmt.Errorf("Role is not found"))
	}

	//Find whether the entity corresponding the Principal is a part of any OCIDs allowed to take the role
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

	//Validate that the filtered list contains atleast one of the OCIDs of the Role
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

	b.Logger().Debug("Login ok", "Method:", method, "targetUrl:", targetUrl, "id", req.ID)

	//Return the response
	resp := &logical.Response{
		Auth: &logical.Auth{
			Period:   time.Duration(roleEntry.TTL) * time.Second,
			Policies: roleEntry.PolicyList,
			Metadata: map[string]string{
				"role_name": roleName,
			},
			InternalData: map[string]interface{}{
				"role_name": roleName,
			},
			DisplayName: roleName,
			LeaseOptions: logical.LeaseOptions{
				Renewable: false,
				TTL:       time.Duration(roleEntry.TTL) * time.Second,
			},
			Alias: &logical.Alias{
				Name: "name",
			},
		},
	}

	return resp, nil
}

func (b *backend) validateHomeTenancy(ctx context.Context, req *logical.Request, homeTenancyId string) error {

	configEntry, err := b.nonLockedOCIConfig(ctx, req.Storage, HOME_TENANCY_ID_CONFIG_NAME)
	if err != nil {
		return err
	}

	if configEntry == nil || configEntry.ConfigValue == "" {
		return fmt.Errorf("Home Tenancy is invalid")
	}
	configuredHomeTenancyId := configEntry.ConfigValue

	if homeTenancyId != configuredHomeTenancyId {
		return fmt.Errorf("Invalid Tenancy")
	}

	return nil
}

func deserializeRequest(data *framework.FieldData, logger log.Logger) (map[string][]string, error) {
	requestHeaders := data.Get("requestHeaders")
	if requestHeaders == nil {
		return nil, errors.New("Empty Authentication Request")
	}

	return convertHeaders(requestHeaders.(map[string]interface{}))
}

// Vault provides the header of type map[string]interface{}, interfaceToString
// meant to convert  map[string]interface{} to map[string][]string for later processing
func convertHeaders(header map[string]interface{}) (map[string][]string, error) {
	returnHeader := map[string][]string{}

	var err error
	for key, _ := range header {
		// HdrRequestTarget and HdrHost are expected to be lower case, anything else needs to be Title case
		if key == HdrRequestTarget || key == HdrHost {
			returnHeader[key], err = interfaceToStringSlice(header[key])
		} else {
			returnHeader[strings.Title(key)], err = interfaceToStringSlice(header[key])
		}
		if err != nil {
			return nil, err
		}
	}
	return returnHeader, nil
}

func interfaceToStringSlice(interfaceList interface{}) ([]string, error) {
	// try to convert interface{} to []interface{}
	switch interfaceList.(type) {
	// If it is already []string, assert type and return value
	case []string:
		interfaceAsList, ok := interfaceList.([]string)
		if !ok {
			return nil, errors.New("interfaceToStringSlice failure 1")
		}
		return interfaceAsList, nil

	case interface{}:
		interfaceAsList, ok := interfaceList.([]interface{})
		if !ok {
			return nil, errors.New("interfaceToStringSlice failure 2")
		}
		returnString := make([]string, len(interfaceAsList))
		// for every element in the slice, try to convert to string
		for i, value := range interfaceAsList {
			switch typedValue := value.(type) {
			case string:
				returnString[i] = typedValue
			default:
				// if not string return error
				return nil, errors.New("interfaceToStringSlice failure 3")
			}
		}
		return returnString, nil
	default:
		return nil, errors.New("interfaceToStringSlice failure 4")
	}
}

func unauthorizedLogicalResponse(req *logical.Request, logger log.Logger, err error) (*logical.Response, error) {
	logger.Debug(req.ID, ": Failed with error:", err)
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
