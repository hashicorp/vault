package alicloud

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/endpoints"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/cidrutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
)

func pathLogin(b *backend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type: framework.TypeString,
				Description: `Name of the role against which the login is being attempted.
If 'role' is not specified, then the login endpoint looks for a role name in the ARN returned by 
the GetCallerIdentity request. If a matching role is not found, login fails.`,
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
			logical.UpdateOperation: b.pathLoginUpdate,
		},
		HelpSynopsis:    pathLoginSyn,
		HelpDescription: pathLoginDesc,
	}
}

func (b *backend) pathLoginUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {

	b64URL := data.Get("identity_request_url").(string)
	if b64URL == "" {
		return nil, errors.New("missing identity_request_url")
	}
	identityReqURL, err := base64.StdEncoding.DecodeString(b64URL)
	if err != nil {
		return nil, errwrap.Wrapf("failed to base64 decode identity_request_url: {{err}}", err)
	}
	if _, err := url.Parse(string(identityReqURL)); err != nil {
		return nil, errwrap.Wrapf("error parsing identity_request_url: {{err}}", err)
	}

	header := data.Get("identity_request_headers").(http.Header)
	if len(header) == 0 {
		return nil, errors.New("missing identity_request_headers")
	}

	callerIdentity, err := b.getCallerIdentity(header, string(identityReqURL))
	if err != nil {
		return nil, errwrap.Wrapf("error making upstream request: {{err}}", err)
	}

	parsedARN, err := parseARN(callerIdentity.Arn)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("unable to parse entity's arn %s due to {{err}}", callerIdentity.Arn), err)
	}
	if parsedARN.Type != arnTypeAssumedRole {
		return nil, fmt.Errorf("only %s arn types are supported at this time, but %s was provided", arnTypeAssumedRole, parsedARN.Type)
	}

	// If a role name was explicitly provided, use that, but otherwise fall back to using the role
	// in the ARN returned by the GetCallerIdentity call.
	roleName := ""
	roleNameIfc, ok := data.GetOk("role")
	if ok {
		roleName = roleNameIfc.(string)
	}
	if roleName == "" {
		roleName = parsedARN.RoleName
	}

	role, err := readRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("entry for role %s not found", parsedARN.RoleName)
	}

	if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.BoundCIDRs) {
		return nil, errors.New("login request originated from invalid CIDR")
	}
	if !parsedARN.IsMemberOf(role.ARN) {
		return nil, errors.New("the caller's arn does not match the role's arn")
	}
	return &logical.Response{
		Auth: &logical.Auth{
			Period:   role.Period,
			Policies: role.Policies,
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
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
			Alias: &logical.Alias{
				Name: callerIdentity.PrincipalId,
			},
			BoundCIDRs: role.BoundCIDRs,
		},
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
	resp.Auth.TTL = role.TTL
	resp.Auth.MaxTTL = role.MaxTTL
	resp.Auth.Period = role.Period
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
	stsEndpoint, err := getSTSEndpoint()
	if err != nil {
		return nil, err
	}
	if u.Host != stsEndpoint {
		return nil, fmt.Errorf(`expected host of "sts.aliyuncs.com" but received "%s"`, u.Host)
	}
	q := u.Query()
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
		return nil, errwrap.Wrapf("error making request: {{err}}", err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, errwrap.Wrapf("error reading response body: {{err}}", err)
		}
		return nil, fmt.Errorf("received %d checking caller identity: %s", response.StatusCode, b)
	}

	result := &sts.GetCallerIdentityResponse{}
	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return nil, errwrap.Wrapf("error decoding response: {{err}}", err)
	}
	return result, nil
}

func getSTSEndpoint() (string, error) {
	r := &endpoints.LocalGlobalResolver{}
	endpoint, supported, err := r.TryResolve(&endpoints.ResolveParam{
		Product: "sts",
	})
	if err != nil {
		return "", err
	}
	if !supported {
		return "", errors.New("sts endpoint is no longer supported")
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
