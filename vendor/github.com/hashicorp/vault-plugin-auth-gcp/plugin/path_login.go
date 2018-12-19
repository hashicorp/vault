package gcpauth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/SermoDigital/jose/jws"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-gcp-common/gcputil"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
	"gopkg.in/square/go-jose.v2/jwt"
)

const (
	expectedJwtAudTemplate string = "vault/%s"

	clientErrorTemplate string = "backend not configured properly, could not create %s client: %v"
)

func pathLogin(b *GcpAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",
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

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

func (b *GcpAuthBackend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Validate we didn't get extraneous fields
	if err := validateFields(req, data); err != nil {
		return nil, logical.CodedError(422, err.Error())
	}

	loginInfo, err := b.parseAndValidateJwt(ctx, req, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	roleType := loginInfo.Role.RoleType
	switch roleType {
	case iamRoleType:
		return b.pathIamLogin(ctx, req, loginInfo)
	case gceRoleType:
		return b.pathGceLogin(ctx, req, loginInfo)
	default:
		return logical.ErrorResponse(fmt.Sprintf("login against role type '%s' is unsupported", roleType)), nil
	}
}

func (b *GcpAuthBackend) pathLoginRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Check role exists and allowed policies are still the same.
	roleName := req.Auth.Metadata["role"]
	if roleName == "" {
		return logical.ErrorResponse("role name metadata not associated with auth token, invalid"), nil
	}
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	} else if role == nil {
		return logical.ErrorResponse("role '%s' no longer exists"), nil
	} else if !policyutil.EquivalentPolicies(role.Policies, req.Auth.Policies) {
		return logical.ErrorResponse("policies on role '%s' have changed, cannot renew"), nil
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
		return nil, fmt.Errorf("unexpected role type '%s' for login renewal", role.RoleType)
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.Period = role.Period
	resp.Auth.TTL = role.TTL
	resp.Auth.MaxTTL = role.MaxTTL
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

func (b *GcpAuthBackend) parseAndValidateJwt(ctx context.Context, req *logical.Request, data *framework.FieldData) (*gcpLoginInfo, error) {
	loginInfo := &gcpLoginInfo{}
	var err error

	loginInfo.RoleName = data.Get("role").(string)
	if loginInfo.RoleName == "" {
		return nil, errors.New("role is required")
	}

	loginInfo.Role, err = b.role(ctx, req.Storage, loginInfo.RoleName)
	if err != nil {
		return nil, err
	}
	if loginInfo.Role == nil {
		return nil, fmt.Errorf("role '%s' not found", loginInfo.RoleName)
	}

	// Process JWT string.
	signedJwt, ok := data.GetOk("jwt")
	if !ok {
		return nil, errors.New("jwt argument is required")
	}

	// Parse 'kid' key id from headers.
	jwtVal, err := jwt.ParseSigned(signedJwt.(string))
	if err != nil {
		return nil, err
	}

	key, err := b.getSigningKey(ctx, jwtVal, signedJwt.(string), loginInfo.Role, req.Storage)
	if err != nil {
		return nil, err
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

	if len(baseClaims.Subject) == 0 {
		return nil, errors.New("expected JWT to have non-empty 'sub' claim")
	}
	loginInfo.EmailOrId = baseClaims.Subject

	if customClaims.Google != nil && customClaims.Google.Compute != nil && len(customClaims.Google.Compute.InstanceId) > 0 {
		loginInfo.GceMetadata = customClaims.Google.Compute
	}

	if loginInfo.Role.RoleType == gceRoleType && loginInfo.GceMetadata == nil {
		return nil, errors.New("expected JWT to have claims with GCE metadata")
	}

	return loginInfo, nil
}

func (b *GcpAuthBackend) getSigningKey(ctx context.Context, token *jwt.JSONWebToken, rawToken string, role *gcpRole, s logical.Storage) (interface{}, error) {
	if len(token.Headers) != 1 {
		return nil, errors.New("expected token to have exactly one header")
	}

	keyId := token.Headers[0].KeyID

	switch role.RoleType {
	case iamRoleType:
		clients, err := b.newGcpClients(ctx, s)
		if err != nil {
			return nil, err
		}
		serviceAccountId, err := parseServiceAccountFromIAMJWT(rawToken)
		if err != nil {
			return nil, err
		}

		accountKey, err := gcputil.ServiceAccountKey(clients.iam, &gcputil.ServiceAccountKeyId{
			Project:   "-",
			EmailOrId: serviceAccountId,
			Key:       keyId,
		})
		if err != nil {
			// Attempt to get a normal Google Oauth cert in case of GCE inferrence.
			key, err := gcputil.OAuth2RSAPublicKey(keyId, "")
			if err != nil {
				return nil, errwrap.Wrapf(
					fmt.Sprintf("could not find service account key or Google Oauth cert with given 'kid' id %s: {{err}}", keyId),
					err)
			}
			return key, nil
		}
		return gcputil.PublicKey(accountKey.PublicKeyData)
	case gceRoleType:
		return gcputil.OAuth2RSAPublicKey(keyId, "")
	default:
		return nil, fmt.Errorf("unexpected role type %s", role.RoleType)
	}
}

// ParseServiceAccountFromIAMJWT parses the service account from the 'sub' claim given a serialized signed JWT.
func parseServiceAccountFromIAMJWT(signedJwt string) (string, error) {
	jwtVal, err := jws.ParseJWT([]byte(signedJwt))
	if err != nil {
		return "", fmt.Errorf("could not parse service account from JWT 'sub' claim: %v", err)
	}
	accountId, ok := jwtVal.Claims().Subject()
	if !ok {
		return "", errors.New("expected 'sub' claim with service account ID or name")
	}
	return accountId, nil
}

func (b *GcpAuthBackend) getGoogleOauthCert(ctx context.Context, keyId string) (interface{}, error) {
	key, err := gcputil.OAuth2RSAPublicKey(keyId, "")
	if err != nil {
		return nil, err
	}
	return key, nil
}

func validateBaseJWTClaims(c *jwt.Claims, roleName string) error {
	exp := c.Expiry.Time()
	if exp.IsZero() || exp.Before(time.Now()) {
		return errors.New("JWT is expired or does not have proper 'exp' claim")
	} else if exp.After(time.Now().Add(time.Minute * time.Duration(maxJwtExpMaxMinutes))) {
		return fmt.Errorf("JWT must expire in %d minutes", maxJwtExpMaxMinutes)
	}

	sub := c.Subject
	if len(sub) < 0 {
		return errors.New("expected JWT to have 'sub' claim with service account id or email")
	}

	expectedAudSuffix := fmt.Sprintf(expectedJwtAudTemplate, roleName)
	for _, aud := range c.Audience {
		if !strings.HasSuffix(aud, expectedAudSuffix) {
			return fmt.Errorf("at least one of the JWT claim 'aud' must end in '%s'", expectedAudSuffix)
		}
	}

	return nil
}

// ---- IAM login domain ----
// pathIamLogin attempts a login operation using the parsed login info.
func (b *GcpAuthBackend) pathIamLogin(ctx context.Context, req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	clients, err := b.newGcpClients(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	role := loginInfo.Role
	if !role.AllowGCEInference && loginInfo.GceMetadata != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Got GCE token but IAM role '%s' does not allow GCE inference", loginInfo.RoleName)), nil
	}

	// TODO(emilymye): move to general JWT validation once custom expiry is supported for other JWT types.
	if loginInfo.JWTClaims.Expiry.Time().After(time.Now().Add(role.MaxJwtExp)) {
		return logical.ErrorResponse(fmt.Sprintf("role requires that JWTs must expire within %d seconds", int(role.MaxJwtExp/time.Second))), nil
	}

	// Get service account and make sure it still exists.
	accountId := &gcputil.ServiceAccountId{
		Project:   "-",
		EmailOrId: loginInfo.EmailOrId,
	}
	serviceAccount, err := gcputil.ServiceAccount(clients.iam, accountId)
	if err != nil {
		return nil, err
	}
	if serviceAccount == nil {
		return nil, errors.New("service account is empty")
	}

	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: serviceAccount.UniqueId,
				},
			},
		}, nil
	}

	// Validate service account can login against role.
	if err := b.authorizeIAMServiceAccount(serviceAccount, role); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period: role.Period,
			Alias: &logical.Alias{
				Name: serviceAccount.UniqueId,
			},
			Policies:    role.Policies,
			Metadata:    authMetadata(loginInfo, serviceAccount),
			DisplayName: serviceAccount.Email,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
		},
	}
	if role.AddGroupAliases {
		clients, err := b.newGcpClients(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		aliases, err := b.groupAliases(clients.resourceManager, ctx, serviceAccount.ProjectId)
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
	clients, err := b.newGcpClients(ctx, req.Storage)
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

	serviceAccount, err := gcputil.ServiceAccount(clients.iam, &gcputil.ServiceAccountId{
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
		return logical.ErrorResponse(fmt.Sprintf(
			"instance %q (project %q) not in bound projects %+v", metadata.InstanceId, metadata.ProjectId, role.BoundProjects)), nil
	}

	// Verify instance exists.
	clients, err := b.newGcpClients(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	instance, err := metadata.GetVerifiedInstance(clients.gce)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"error when attempting to find instance (project %s, zone: %s, instance: %s) :%v",
			metadata.ProjectId, metadata.Zone, metadata.InstanceName, err)), nil
	}

	if err := b.authorizeGCEInstance(ctx, metadata.ProjectId, instance, req.Storage, role, loginInfo.EmailOrId); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if req.Operation == logical.AliasLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Alias: &logical.Alias{
					Name: fmt.Sprintf("gce-%s", strconv.FormatUint(instance.Id, 10)),
				},
			},
		}, nil
	}

	serviceAccount, err := gcputil.ServiceAccount(clients.iam, &gcputil.ServiceAccountId{
		Project:   "-",
		EmailOrId: loginInfo.EmailOrId,
	})
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Could not find service account '%s' used for GCE metadata token: %s",
			loginInfo.EmailOrId, err)), nil
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			InternalData: map[string]interface{}{},
			Period:       role.Period,
			Alias: &logical.Alias{
				Name: fmt.Sprintf("gce-%s", strconv.FormatUint(instance.Id, 10)),
			},
			Policies:    role.Policies,
			Metadata:    authMetadata(loginInfo, serviceAccount),
			DisplayName: instance.Name,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
		},
	}

	if role.AddGroupAliases {
		aliases, err := b.groupAliases(clients.resourceManager, ctx, metadata.ProjectId)
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
//   "project-my-project"
//   "folder-my-subfolder"
//   "folder-my-project"
//   "organization-my-org"
// ]
func (b *GcpAuthBackend) groupAliases(crmClient *cloudresourcemanager.Service, ctx context.Context, projectId string) ([]*logical.Alias, error) {
	ancestry, err := crmClient.Projects.
		GetAncestry(projectId, &cloudresourcemanager.GetAncestryRequest{}).
		Context(ctx).
		Do()
	if err != nil {
		return nil, err
	}

	aliases := make([]*logical.Alias, len(ancestry.Ancestor)+1)
	aliases[0] = &logical.Alias{
		Name: fmt.Sprintf("project-%s", projectId),
	}
	for i, parent := range ancestry.Ancestor {
		aliases[i+1] = &logical.Alias{
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
	httpC, err := b.httpClient(ctx, req.Storage)
	if err != nil {
		return err
	}

	gceClient, err := compute.New(httpC)
	if err != nil {
		return fmt.Errorf(clientErrorTemplate, "GCE", err)
	}

	meta, err := getInstanceMetadataFromAuth(req.Auth.Metadata)
	if err != nil {
		return fmt.Errorf("invalid auth metadata: %v", err)
	}

	instance, err := meta.GetVerifiedInstance(gceClient)
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
		return nil, fmt.Errorf("expected 'project_number' value '%s' to be a int64", projectNumber)
	}

	createdAt, ok := authMetadata["instance_creation_timestamp"]
	if !ok {
		return nil, errors.New("expected 'instance_creation_timestamp' field")
	}
	meta.CreatedAt, err = strconv.ParseInt(createdAt, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("expected 'instance_creation_timestamp' value '%s' to be int64", createdAt)
	}

	return meta, nil
}

// authorizeGCEInstance returns an error if the given GCE instance is not
// authorized for the role.
func (b *GcpAuthBackend) authorizeGCEInstance(ctx context.Context, project string, instance *compute.Instance, s logical.Storage, role *gcpRole, serviceAccountId string) error {
	httpC, err := b.httpClient(ctx, s)
	if err != nil {
		return err
	}

	iamClient, err := iam.New(httpC)
	if err != nil {
		return fmt.Errorf(clientErrorTemplate, "IAM", err)
	}

	gceClient, err := compute.New(httpC)
	if err != nil {
		return fmt.Errorf(clientErrorTemplate, "GCE", err)
	}

	return AuthorizeGCE(ctx, &AuthorizeGCEInput{
		client: &gcpClient{
			computeSvc: gceClient,
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

const pathLoginHelpSyn = `Authenticates Google Cloud Platform entities with Vault.`
const pathLoginHelpDesc = `
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
