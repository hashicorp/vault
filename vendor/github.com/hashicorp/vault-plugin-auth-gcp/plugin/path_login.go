// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package gcpauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/go-jose/go-jose/v4/jwt"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/helper/policyutil"
	"github.com/hashicorp/vault/sdk/logical"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
)

const (
	expectedJwtAudTemplate string = "vault/%s"
	jwtExpToleranceSec            = 60
)

var (
	allowedSignatureAlgorithms = []jose.SignatureAlgorithm{
		jose.RS256,
		jose.ES256,
		jose.HS256,
	}
)

func pathLogin(b *GcpAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixGoogleCloud,
			OperationVerb:   "login",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `Name of the role against which the login is being attempted. Required.`,
			},
			"jwt": {
				Type: framework.TypeString,
				Description: `
A signed JWT. This is either a self-signed service account JWT ('iam' roles only) or a
GCE identity metadata token ('iam', 'gce' roles).`,
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLogin,
			},
			logical.AliasLookaheadOperation: &framework.PathOperation{
				Callback: b.pathLogin,
			},
			logical.ResolveRoleOperation: &framework.PathOperation{
				Callback: b.pathResolveRole,
			},
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}
func (b *GcpAuthBackend) pathResolveRole(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	roleName := data.Get("role").(string)
	if roleName == "" {
		return logical.ErrorResponse("role is required"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q not found", roleName), nil
	}

	return logical.ResolveRoleResponse(roleName)
}

func (b *GcpAuthBackend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(http.StatusUnprocessableEntity, err.Error())
	}

	loginInfo, err := b.parseAndValidateJwt(ctx, req.Storage, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if len(loginInfo.Role.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, loginInfo.Role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	roleType := loginInfo.Role.RoleType
	switch roleType {
	case iamRoleType:
		return b.pathIamLogin(ctx, req, loginInfo)
	case gceRoleType:
		return b.pathGceLogin(ctx, req, loginInfo)
	default:
		return logical.ErrorResponse("login against role type %q is unsupported", roleType), nil
	}
}

func (b *GcpAuthBackend) pathLoginRenew(ctx context.Context, req *logical.Request, _ *framework.FieldData) (*logical.Response, error) {
	// Check role exists and allowed policies are still the same.
	roleName := req.Auth.Metadata["role"]
	if roleName == "" {
		return logical.ErrorResponse("role name metadata not associated with auth token, invalid"), nil
	}
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	} else if role == nil {
		return logical.ErrorResponse("role %q no longer exists", roleName), nil
	} else if !policyutil.EquivalentPolicies(role.TokenPolicies, req.Auth.TokenPolicies) {
		return logical.ErrorResponse("policies on role %q have changed, cannot renew", roleName), nil
	}

	switch role.RoleType {
	case iamRoleType:
		if err := b.pathIamRenew(ctx, req, roleName, role); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	case gceRoleType:
		if err := b.pathGceRenew(ctx, req, roleName, role); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	default:
		return nil, fmt.Errorf("unexpected role type %q for login renewal", role.RoleType)
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.Period = role.TokenPeriod
	resp.Auth.TTL = role.TokenTTL
	resp.Auth.MaxTTL = role.TokenMaxTTL
	return resp, nil
}

// gcpLoginInfo represents the data given to Vault for logging in using the IAM method.
type gcpLoginInfo struct {
	// Name of the role being logged in against
	RoleName string

	// Role being logged in against
	Role *gcpRole

	// ID or email of an IAM service account or that inferred for a GCE VM.
	EmailOrId string

	// Base JWT Claims (registered claims such as 'exp', 'iss', etc)
	JWTClaims *jwt.Claims

	// Metadata from a GCE instance identity token.
	GceMetadata *gcputil.GCEIdentityMetadata
}

func (b *GcpAuthBackend) parseAndValidateJwt(ctx context.Context, s logical.Storage, data *framework.FieldData) (*gcpLoginInfo, error) {
	loginInfo := &gcpLoginInfo{}
	var err error

	conf, err := b.config(ctx, s)
	if err != nil {
		return nil, errors.New("unable to retrieve GCP configuration")
	}

	roleName := data.Get("role").(string)
	if roleName == "" {
		return nil, errors.New("role is required")
	}

	role, err := b.role(ctx, s, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, fmt.Errorf("role %q not found", roleName)
	}

	loginInfo.RoleName = roleName
	loginInfo.Role = role

	// Process JWT string.
	signedJwt, ok := data.GetOk("jwt")
	if !ok {
		return nil, errors.New("jwt argument is required")
	}

	// Parse 'kid' key id from headers.
	jwtVal, err := jwt.ParseSigned(signedJwt.(string), allowedSignatureAlgorithms)
	if err != nil {
		return nil, fmt.Errorf("unable to parse signed JWT: %w", err)
	}

	key, err := b.getSigningKey(ctx, jwtVal, signedJwt.(string), loginInfo.Role, conf.APICustomEndpoint)
	if err != nil {
		return nil, fmt.Errorf("unable to get public key for signed JWT: %w", err)
	}

	// Parse claims and verify signature.
	baseClaims := &jwt.Claims{}
	customClaims := &gcputil.CustomJWTClaims{}

	if err = jwtVal.Claims(key, baseClaims, customClaims); err != nil {
		return nil, err
	}

	if err = validateBaseJWTClaims(baseClaims, loginInfo.RoleName); err != nil {
		return nil, err
	}

	loginInfo.JWTClaims = baseClaims
	loginInfo.EmailOrId = baseClaims.Subject

	if loginInfo.Role.RoleType == gceRoleType {
		if customClaims.Google != nil && customClaims.Google.Compute != nil && len(customClaims.Google.Compute.InstanceId) > 0 {
			loginInfo.GceMetadata = customClaims.Google.Compute
		}
		if loginInfo.GceMetadata == nil {
			return nil, errors.New("expected JWT to have claims with GCE metadata")
		}
	}
	return loginInfo, nil
}

func (b *GcpAuthBackend) getSigningKey(ctx context.Context, token *jwt.JSONWebToken, rawToken string, role *gcpRole, endpoint string) (interface{}, error) {
	b.Logger().Debug("Getting signing Key for JWT")

	if len(token.Headers) != 1 {
		return nil, errors.New("expected token to have exactly one header")
	}
	kid := token.Headers[0].KeyID
	b.Logger().Debug("kid found for JWT", "kid", kid)

	// Try getting Google-wide key
	k, gErr := gcputil.OAuth2RSAPublicKeyWithEndpoint(ctx, kid, endpoint)
	if gErr == nil {
		b.Logger().Debug("Found Google OAuth2 provider key", "kid", kid)
		return k, nil
	}

	saId, err := getJWTSubject(rawToken)
	if err != nil {
		return nil, err
	}

	if role.RoleType == iamRoleType {
		// If that failed, and the authentication type is IAM, try to get account-specific key
		b.Logger().Debug("Unable to get Google-wide OAuth2 Key, trying service-account public key")
		k, saErr := gcputil.ServiceAccountPublicKeyWithEndpoint(ctx, saId, kid, endpoint)
		if saErr == nil {
			return k, nil
		}
		return nil, fmt.Errorf("unable to get public key %q for JWT subject %q: %w", kid, saId, saErr)
	}
	return nil, fmt.Errorf("unable to get public key %q for JWT subject %q: no Google OAuth2 provider key found for GCE role", kid, saId)
}

// getJWTSubject grabs 'sub' claim given an unverified signed JWT.
func getJWTSubject(signedJwt string) (string, error) {
	jwtVal, err := jwt.ParseSigned(signedJwt, allowedSignatureAlgorithms)
	if err != nil {
		return "", fmt.Errorf("could not parse JWT: %v", err)
	}
	var claims jwt.Claims
	if err = jwtVal.UnsafeClaimsWithoutVerification(&claims); err != nil {
		return "", fmt.Errorf("could not parse claims from JWT: %v", err)
	}
	accountID := claims.Subject
	if accountID == "" {
		return "", errors.New("expected 'sub' claim from JWT")
	}
	return accountID, nil
}

func validateBaseJWTClaims(c *jwt.Claims, roleName string) error {
	exp := c.Expiry.Time()
	tolDelt := time.Second * jwtExpToleranceSec
	// Compare expiration to current time with tolerance
	if exp.IsZero() || exp.Before(time.Now().Add(-tolDelt)) {
		return errors.New("JWT is expired or does not have proper 'exp' claim")
	}

	// Compare expiration to max expiration with tolerance
	allowedDelta := time.Minute*time.Duration(maxJwtExpMaxMinutes) + tolDelt
	expIn := exp.Sub(time.Now())
	if expIn > allowedDelta {
		return fmt.Errorf("JWT must expire in %d minutes, expires in %v", maxJwtExpMaxMinutes, expIn)
	}

	if len(c.Subject) == 0 {
		return errors.New("expected JWT to have 'sub' claim with service account id or email")
	}

	expectedAudSuffix := fmt.Sprintf(expectedJwtAudTemplate, roleName)
	for _, aud := range c.Audience {
		if !strings.HasSuffix(aud, expectedAudSuffix) {
			return fmt.Errorf("at least one of the JWT claim 'aud' must end in %q", expectedAudSuffix)
		}
	}

	return nil
}

// ---- IAM login domain ----
// pathIamLogin attempts a login operation using the parsed login info.
func (b *GcpAuthBackend) pathIamLogin(ctx context.Context, req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	iamClient, err := b.IAMClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	role := loginInfo.Role
	if !role.AllowGCEInference && loginInfo.GceMetadata != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Got GCE token but IAM role %q does not allow GCE inference", loginInfo.RoleName)), nil
	}

	// TODO(emilymye): move to general JWT validation once custom expiry is supported for other JWT types.
	if loginInfo.GceMetadata != nil {
		b.Logger().Info("GCE Metadata found in JWT, skipping custom expiry check")
	} else if loginInfo.JWTClaims.Expiry.Time().After(time.Now().Add(role.MaxJwtExp)) {
		return logical.ErrorResponse("role requires that service account JWTs expire within %d seconds", int(role.MaxJwtExp/time.Second)), nil
	}

	// Get service account and make sure it still exists.
	accountId := &gcputil.ServiceAccountId{
		Project:   "-",
		EmailOrId: loginInfo.EmailOrId,
	}
	serviceAccount, err := gcputil.ServiceAccount(iamClient, accountId)
	if err != nil {
		return nil, err
	}
	if serviceAccount == nil {
		return nil, errors.New("service account is empty")
	}

	conf, err := b.config(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("unable to retrieve GCP configuration"), nil
	}

	alias, err := conf.getIAMAlias(role, serviceAccount)
	if err != nil {
		return logical.ErrorResponse("unable to create alias: %s", err), nil
	}

	if req.Operation == logical.AliasLookaheadOperation {
		resp := &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: alias,
				},
			},
		}
		return resp, nil
	}

	// Validate service account can login against role.
	if err := b.authorizeIAMServiceAccount(serviceAccount, role); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	auth := &logical.Auth{
		Alias: &logical.Alias{
			Name: alias,
		},
		DisplayName: serviceAccount.Email,
	}
	role.PopulateTokenAuth(auth)
	if err := conf.IAMAuthMetadata.PopulateDesiredMetadata(auth, authMetadata(loginInfo, serviceAccount)); err != nil {
		b.Logger().Warn("unable to populate iam metadata", "err", err.Error())
	}

	resp := &logical.Response{
		Auth: auth,
	}

	if role.AddGroupAliases {
		crmClient, err := b.CRMClient(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		aliases, err := b.groupAliases(crmClient, ctx, serviceAccount.ProjectId)
		if err != nil {
			return nil, err
		}
		resp.Auth.GroupAliases = aliases
	}

	return resp, nil
}

// pathIamRenew returns an error if the service account referenced in the auth token metadata cannot renew the
// auth token for the given role.
func (b *GcpAuthBackend) pathIamRenew(ctx context.Context, req *logical.Request, roleName string, role *gcpRole) error {
	iamClient, err := b.IAMClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	serviceAccountId, ok := req.Auth.Metadata["service_account_id"]
	if !ok {
		return errors.New("service account id metadata not associated with auth token, invalid")
	}

	// This project is the service account's project.
	project, ok := req.Auth.Metadata["project_id"]
	if !ok {
		project = "-"
	}

	serviceAccount, err := gcputil.ServiceAccount(iamClient, &gcputil.ServiceAccountId{
		Project:   project,
		EmailOrId: serviceAccountId,
	})
	if err != nil {
		return fmt.Errorf("cannot find service account %s", serviceAccountId)
	}

	_, isGceInferred := req.Auth.Metadata["instance_id"]
	if isGceInferred && !role.AllowGCEInference {
		return fmt.Errorf("GCE inferrence is no longer allowed for role %s", roleName)
	}

	if err := b.authorizeIAMServiceAccount(serviceAccount, role); err != nil {
		return fmt.Errorf("service account is no longer authorized for role %s", roleName)
	}

	return nil
}

// validateAgainstIAMRole returns an error if the given IAM service account is not authorized for the role.
func (b *GcpAuthBackend) authorizeIAMServiceAccount(serviceAccount *iam.ServiceAccount, role *gcpRole) error {
	if len(role.BoundProjects) > 0 && !strutil.StrListContains(role.BoundProjects, serviceAccount.ProjectId) {
		return fmt.Errorf("service account %q not in bound projects %+v", serviceAccount.Email, role.BoundProjects)
	}

	// Check if role has the wildcard as the only service account.
	if len(role.BoundServiceAccounts) == 1 && role.BoundServiceAccounts[0] == serviceAccountsWildcard {
		return nil
	}

	// Check for service account id/email.
	if strutil.StrListContains(role.BoundServiceAccounts, serviceAccount.Email) ||
		strutil.StrListContains(role.BoundServiceAccounts, serviceAccount.UniqueId) {
		return nil
	}

	return fmt.Errorf("service account %s (id: %s) is not authorized for role",
		serviceAccount.Email, serviceAccount.UniqueId)
}

// ---- GCE login domain ----
// pathGceLogin attempts a login operation using the parsed login info.
func (b *GcpAuthBackend) pathGceLogin(ctx context.Context, req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	role := loginInfo.Role
	metadata := loginInfo.GceMetadata
	if metadata == nil {
		return logical.ErrorResponse("could not get GCE metadata from given JWT"), nil
	}

	if len(role.BoundProjects) > 0 && !strutil.StrListContains(role.BoundProjects, metadata.ProjectId) {
		return logical.ErrorResponse("instance %q (project %q) not in bound projects %+v", metadata.InstanceId, metadata.ProjectId, role.BoundProjects), nil
	}

	// Verify instance exists.
	computeClient, err := b.ComputeClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	instance, err := metadata.GetVerifiedInstance(computeClient)
	if err != nil {
		return logical.ErrorResponse("error when attempting to find instance (project %s, zone: %s, instance: %s) :%v",
			metadata.ProjectId, metadata.Zone, metadata.InstanceName, err), nil
	}

	if err := b.authorizeGCEInstance(ctx, metadata.ProjectId, instance, req.Storage, role, loginInfo.EmailOrId); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	conf, err := b.config(ctx, req.Storage)
	if err != nil {
		return logical.ErrorResponse("unable to retrieve GCP configuration"), nil
	}

	alias, err := conf.getGCEAlias(role, instance)
	if err != nil {
		return logical.ErrorResponse("unable to create alias: %s", err), nil
	}

	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: alias,
				},
			},
		}, nil
	}

	iamClient, err := b.IAMClient(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	serviceAccount, err := gcputil.ServiceAccount(iamClient, &gcputil.ServiceAccountId{
		Project:   "-",
		EmailOrId: loginInfo.EmailOrId,
	})
	if err != nil {
		return logical.ErrorResponse("Could not find service account %q used for GCE metadata token: %s", loginInfo.EmailOrId, err), nil
	}

	auth := &logical.Auth{
		InternalData: map[string]interface{}{},
		Alias: &logical.Alias{
			Name: alias,
		},
		DisplayName: instance.Name,
	}
	role.PopulateTokenAuth(auth)
	if err := conf.GCEAuthMetadata.PopulateDesiredMetadata(auth, authMetadata(loginInfo, serviceAccount)); err != nil {
		b.Logger().Warn("unable to populate gce metadata", "err", err.Error())
	}

	resp := &logical.Response{
		Auth: auth,
	}

	if role.AddGroupAliases {
		crmClient, err := b.CRMClient(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		aliases, err := b.groupAliases(crmClient, ctx, metadata.ProjectId)
		if err != nil {
			return nil, err
		}
		resp.Auth.GroupAliases = aliases
	}
	return resp, nil
}

// groupAliases will add group aliases for an authenticating GCP entity
// starting at project-level and going up the Cloud Resource Manager
// hierarchy
//
// For example, given a project hierarchy of
// "my-org" --> "my-folder" --> "my-subfolder" --> "my-project",
// this returns the following group aliases:
// [
//
//	"project-my-project"
//	"folder-my-subfolder"
//	"folder-my-project"
//	"organization-my-org"
//
// ]
func (b *GcpAuthBackend) groupAliases(crmClient *cloudresourcemanager.Service, ctx context.Context, projectId string) ([]*logical.Alias, error) {
	ancestry, err := crmClient.Projects.
		GetAncestry(projectId, &cloudresourcemanager.GetAncestryRequest{}).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	aliases := make([]*logical.Alias, len(ancestry.Ancestor))
	for i, parent := range ancestry.Ancestor {
		aliases[i] = &logical.Alias{
			Name: fmt.Sprintf("%s-%s", parent.ResourceId.Type, parent.ResourceId.Id),
		}
	}
	return aliases, nil
}

func authMetadata(loginInfo *gcpLoginInfo, serviceAccount *iam.ServiceAccount) map[string]string {
	metadata := map[string]string{
		"role":                  loginInfo.RoleName,
		"service_account_id":    serviceAccount.UniqueId,
		"service_account_email": serviceAccount.Email,
		"project_id":            serviceAccount.ProjectId,
	}

	if loginInfo.GceMetadata != nil {
		gceMetadata := loginInfo.GceMetadata
		metadata["project_id"] = gceMetadata.ProjectId
		metadata["project_number"] = strconv.FormatInt(gceMetadata.ProjectNumber, 10)
		metadata["zone"] = gceMetadata.Zone
		metadata["instance_id"] = gceMetadata.InstanceId
		metadata["instance_name"] = gceMetadata.InstanceName
		metadata["instance_creation_timestamp"] = strconv.FormatInt(gceMetadata.CreatedAt, 10)
	}
	return metadata
}

// pathGceRenew returns an error if the instance referenced in the auth token metadata cannot renew the
// auth token for the given role.
func (b *GcpAuthBackend) pathGceRenew(ctx context.Context, req *logical.Request, roleName string, role *gcpRole) error {
	computeClient, err := b.ComputeClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	meta, err := getInstanceMetadataFromAuth(req.Auth.Metadata)
	if err != nil {
		return fmt.Errorf("invalid auth metadata: %v", err)
	}

	instance, err := meta.GetVerifiedInstance(computeClient)
	if err != nil {
		return err
	}

	serviceAccountId, ok := req.Auth.Metadata["service_account_id"]
	if !ok {
		return errors.New("invalid auth metadata: service_account_id not found")
	}
	if err := b.authorizeGCEInstance(ctx, meta.ProjectId, instance, req.Storage, role, serviceAccountId); err != nil {
		return fmt.Errorf("could not renew token for role %s: %v", roleName, err)
	}

	return nil
}

func getInstanceMetadataFromAuth(authMetadata map[string]string) (*gcputil.GCEIdentityMetadata, error) {
	meta := &gcputil.GCEIdentityMetadata{}
	var ok bool
	var err error

	meta.ProjectId, ok = authMetadata["project_id"]
	if !ok {
		return nil, errors.New("expected 'project_id' field")
	}

	meta.Zone, ok = authMetadata["zone"]
	if !ok {
		return nil, errors.New("expected 'zone' field")
	}

	meta.InstanceId, ok = authMetadata["instance_id"]
	if !ok {
		return nil, errors.New("expected 'instance_id' field")
	}

	meta.InstanceName, ok = authMetadata["instance_name"]
	if !ok {
		return nil, errors.New("expected 'instance_name' field")
	}

	// Parse numbers back into int values.
	projectNumber, ok := authMetadata["project_number"]
	if !ok {
		return nil, errors.New("expected 'project_number' field, got %v")
	}
	meta.ProjectNumber, err = strconv.ParseInt(projectNumber, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("expected 'project_number' value %q to be a int64", projectNumber)
	}

	createdAt, ok := authMetadata["instance_creation_timestamp"]
	if !ok {
		return nil, errors.New("expected 'instance_creation_timestamp' field")
	}
	meta.CreatedAt, err = strconv.ParseInt(createdAt, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("expected 'instance_creation_timestamp' value %q to be int64", createdAt)
	}

	return meta, nil
}

// authorizeGCEInstance returns an error if the given GCE instance is not
// authorized for the role.
func (b *GcpAuthBackend) authorizeGCEInstance(ctx context.Context, project string, instance *compute.Instance, s logical.Storage, role *gcpRole, serviceAccountId string) error {
	iamClient, err := b.IAMClient(ctx, s)
	if err != nil {
		return err
	}

	computeClient, err := b.ComputeClient(ctx, s)
	if err != nil {
		return nil
	}

	return AuthorizeGCE(ctx, &AuthorizeGCEInput{
		client: &gcpClient{
			logger:     b.Logger(),
			computeSvc: computeClient,
			iamSvc:     iamClient,
		},
		serviceAccount:   serviceAccountId,
		project:          project,
		instanceLabels:   instance.Labels,
		instanceSelfLink: instance.SelfLink,
		instanceZone:     instance.Zone,

		boundLabels:  role.BoundLabels,
		boundRegions: role.BoundRegions,
		boundZones:   role.BoundZones,

		boundInstanceGroups:  role.BoundInstanceGroups,
		boundServiceAccounts: role.BoundServiceAccounts,
	})
}

const (
	pathLoginHelpSyn  = `Authenticates Google Cloud Platform entities with Vault.`
	pathLoginHelpDesc = `
Authenticate Google Cloud Platform (GCP) entities.

Currently supports authentication for:

IAM service accounts
=====================
IAM service accounts can use GCP APIs or tools to sign a JSON Web Token (JWT).
This JWT should contain the id (expected field 'client_id') or email
(expected field 'client_email') of the authenticating service account in its claims.
Vault verifies the signed JWT and parses the identity of the account.

Renewal is rejected if the role, service account, or original signing key no longer exists.
`
)
