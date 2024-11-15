// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package alicloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/logical"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixAliCloud,
			OperationVerb:   "login",
		},
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type: framework.TypeString,
				Description: `Name of the role against which the login is being attempted.
If a matching role is not found, login fails.`,
				Required: true,
			},
			"identity_request_url": {
				Type:        framework.TypeString,
				Description: "Base64-encoded full URL against which to make the AliCloud request.",
			},
			"identity_request_headers": {
				Type: framework.TypeHeader,
				Description: `The request headers. This must include the headers over which AliCloud
has included a signature.`,
			},
		},
		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:      b.pathLoginUpdate,
			logical.ResolveRoleOperation: b.pathLoginResolveRole,
		},
		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

// pathLoginResolveRole will identify the role that pathLoginUpdate will use to log-in
// Note: Most of this function is duplicated logic. The reason for this is so that callers
// to this function receive logical errors instead of internal server errors where appropriate
// logic updates relating to role determination should be kept consistent between the two.
func (b *backend) pathLoginResolveRole(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	var err error

	roleNameIfc, ok := data.GetOk("role")
	if !ok {
		return logical.ErrorResponse("missing role"), nil
	}
	roleName := roleNameIfc.(string)

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("entry for role %s not found", roleName), nil
	}

	return logical.ResolveRoleResponse(roleName)
}

func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	b64URL := data.Get("identity_request_url").(string)
	if b64URL == "" {
		return nil, errors.New("missing identity_request_url")
	}
	identityReqURL, err := base64.StdEncoding.DecodeString(b64URL)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 decode identity_request_url: %w", err)
	}
	if _, err := url.Parse(string(identityReqURL)); err != nil {
		return nil, fmt.Errorf("error parsing identity_request_url: %w", err)
	}
	header := data.Get("identity_request_headers").(http.Header)
	if len(header) == 0 {
		return nil, errors.New("missing identity_request_headers")
	}

	callerIdentity, err := b.getCallerIdentity(header, string(identityReqURL))
	if err != nil {
		return nil, fmt.Errorf("error making upstream request: %w", err)
	}

	parsedARN, err := parseARN(callerIdentity.Arn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse entity's arn %s due to %w", callerIdentity.Arn, err)
	}
	if parsedARN.Type != arnTypeAssumedRole {
		return nil, fmt.Errorf("only %s arn types are supported at this time, but %s was provided", arnTypeAssumedRole, parsedARN.Type)
	}

	roleNameIfc, ok := data.GetOk("role")
	if !ok {
		return nil, errors.New("missing role")
	}
	roleName := roleNameIfc.(string)

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("entry for role %s not found", roleName)
	}

	if len(role.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	if !parsedARN.IsMemberOf(role.ARN) {
		return nil, errors.New("the caller's arn does not match the role's arn")
	}

	auth := &logical.Auth{
		Metadata: map[string]string{
			"account_id":    callerIdentity.AccountId,
			"user_id":       callerIdentity.UserId,
			"role_id":       callerIdentity.RoleId,
			"arn":           callerIdentity.Arn,
			"identity_type": callerIdentity.IdentityType,
			"principal_id":  callerIdentity.PrincipalId,
			"request_id":    callerIdentity.RequestId,
			"role_name":     roleName,
		},
		DisplayName: callerIdentity.PrincipalId,
		Alias: &logical.Alias{
			Name: callerIdentity.PrincipalId,
		},
	}

	role.PopulateTokenAuth(auth)

	return &logical.Response{
		Auth: auth,
	}, nil
}

func (b *backend) pathLoginRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// The arn set in metadata earlier is the assumed-role arn.
	arn := req.Auth.Metadata["arn"]
	if arn == "" {
		return nil, errors.New("unable to retrieve arn from metadata during renewal")
	}
	parsedARN, err := parseARN(arn)
	if err != nil {
		return nil, err
	}

	roleName, ok := req.Auth.Metadata["role_name"]
	if !ok {
		return nil, errors.New("error retrieving role_name during renewal")
	}

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, errors.New("role entry not found")
	}

	if !parsedARN.IsMemberOf(role.ARN) {
		return nil, errors.New("the caller's arn does not match the role's arn")
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = role.TokenTTL
	resp.Auth.MaxTTL = role.TokenMaxTTL
	resp.Auth.Period = role.TokenPeriod
	return resp, nil
}

func (b *backend) getCallerIdentity(header http.Header, rawURL string) (*sts.GetCallerIdentityResponse, error) {
	/*
		Here we need to ensure we're actually hitting the AliCloud service, and that the caller didn't
		inject a URL to their own service that will respond as desired.
	*/
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	if u.Scheme != "https" {
		return nil, fmt.Errorf(`expected "https" url scheme but received "%s"`, u.Scheme)
	}
	q := u.Query()
	regionID := q.Get("RegionId")
	if regionID == "" {
		return nil, fmt.Errorf("query RegionId must not be empty")
	}

	stsEndpoint, err := getSTSEndpoint(regionID)
	if err != nil {
		return nil, err
	}
	if u.Host != stsEndpoint {
		return nil, fmt.Errorf(`expected host of "%s" but received "%s"`, stsEndpoint, u.Host)
	}
	if q.Get("Format") != "JSON" {
		return nil, fmt.Errorf("query Format must be JSON but received %s", q.Get("Format"))
	}
	if q.Get("Action") != "GetCallerIdentity" {
		return nil, fmt.Errorf("query Action must be GetCallerIdentity but received %s", q.Get("Action"))
	}

	request, err := http.NewRequest(http.MethodPost, rawURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header = header

	response, err := b.identityClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("error reading response body: %w", err)
		}
		return nil, fmt.Errorf("received %d checking caller identity: %s", response.StatusCode, b)
	}

	result := &sts.GetCallerIdentityResponse{}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}
	return result, nil
}

// getSTSEndpoint will build endpoints from the given region using the sts
// GetEndpointRules API method. See the endpoint docs at:
// https://www.alibabacloud.com/help/en/resource-access-management/latest/api-doc-sts-2015-04-01-endpoint
//
// Alicloud Support said that there is not an API to fetch all sts endpoints.
// See the github ticket at: https://github.com/aliyun/alibaba-cloud-sdk-go/issues/577
func getSTSEndpoint(regionID string) (string, error) {
	config := sdk.NewConfig()
	config.Scheme = "https"

	// we don't need real creds because we only need the client to build the
	// endpoint for the given region
	creds := credentials.NewAccessKeyCredential("", "")
	client, err := sts.NewClientWithOptions(regionID, config, creds)
	if err != nil {
		return "", err
	}
	endpoint, err := client.GetEndpointRules(regionID, "sts")
	if err != nil {
		return "", err
	}
	if endpoint == "" {
		return "", errors.New("got an empty endpoint")
	}

	return endpoint, nil
}

const pathLoginSyn = `
Authenticates an RAM entity with Vault.
`

const pathLoginDesc = `
Authenticate AliCloud entities using an arbitrary RAM principal.

RAM principals are authenticated by processing a signed sts:GetCallerIdentity
request and then parsing the response to see who signed the request.
`
