package gcpauth

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/util"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
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

func (b *GcpAuthBackend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	loginInfo, err := b.parseAndValidateJwt(req, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	roleType := loginInfo.Role.RoleType
	switch roleType {
	case iamRoleType:
		return b.pathIamLogin(req, loginInfo)
	case gceRoleType:
		return b.pathGceLogin(req, loginInfo)
	default:
		return logical.ErrorResponse(fmt.Sprintf("login against role type '%s' is unsupported", roleType)), nil
	}
}

func (b *GcpAuthBackend) pathLoginRenew(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	// Check role exists and allowed policies are still the same.
	roleName := req.Auth.Metadata["role"]
	if roleName == "" {
		return logical.ErrorResponse("role name metadata not associated with auth token, invalid"), nil
	}
	role, err := b.role(req.Storage, roleName)
	if err != nil {
		return nil, err
	} else if role == nil {
		return logical.ErrorResponse("role '%s' no longer exists"), nil
	} else if !policyutil.EquivalentPolicies(role.Policies, req.Auth.Policies) {
		return logical.ErrorResponse("policies on role '%s' have changed, cannot renew"), nil
	}

	switch role.RoleType {
	case iamRoleType:
		if err := b.pathIamRenew(req, roleName, role); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	case gceRoleType:
		if err := b.pathGceRenew(req, roleName, role); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	default:
		return nil, fmt.Errorf("unexpected role type '%s' for login renewal", role.RoleType)
	}

	// If 'Period' is set on the Role, the token should never expire.
	if role.Period > 0 {
		// Replenish the TTL with current role's Period.
		req.Auth.TTL = role.Period
		return &logical.Response{Auth: req.Auth}, nil
	} else {
		return framework.LeaseExtend(role.TTL, role.MaxTTL, b.System())(req, data)
	}
}

// gcpLoginInfo represents the data given to Vault for logging in using the IAM method.
type gcpLoginInfo struct {
	// Name of the role being logged in against
	RoleName string

	// Role being logged in against
	Role *gcpRole

	// ID or email of an IAM service account or that inferred for a GCE VM.
	ServiceAccountId string

	// Base JWT Claims (registered claims such as 'exp', 'iss', etc)
	JWTClaims *jwt.Claims

	// Metadata from a GCE instance identity token.
	GceMetadata *util.GCEIdentityMetadata
}

func (b *GcpAuthBackend) parseAndValidateJwt(req *logical.Request, data *framework.FieldData) (*gcpLoginInfo, error) {
	loginInfo := &gcpLoginInfo{}
	var err error

	loginInfo.RoleName = data.Get("role").(string)
	if loginInfo.RoleName == "" {
		return nil, errors.New("role is required")
	}

	loginInfo.Role, err = b.role(req.Storage, loginInfo.RoleName)
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

	key, err := b.getSigningKey(jwtVal, signedJwt.(string), loginInfo.Role, req.Storage)
	if err != nil {
		return nil, err
	}

	// Parse claims and verify signature.
	baseClaims := &jwt.Claims{}
	customClaims := &util.CustomJWTClaims{}

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
	loginInfo.ServiceAccountId = baseClaims.Subject

	if customClaims.Google != nil && customClaims.Google.Compute != nil && len(customClaims.Google.Compute.InstanceId) > 0 {
		loginInfo.GceMetadata = customClaims.Google.Compute
	}

	if loginInfo.Role.RoleType == gceRoleType && loginInfo.GceMetadata == nil {
		return nil, errors.New("expected JWT to have claims with GCE metadata")
	}

	return loginInfo, nil
}

func (b *GcpAuthBackend) getSigningKey(token *jwt.JSONWebToken, rawToken string, role *gcpRole, s logical.Storage) (interface{}, error) {
	if len(token.Headers) != 1 {
		return nil, errors.New("expected token to have exactly one header")
	}

	keyId := token.Headers[0].KeyID

	switch role.RoleType {
	case iamRoleType:
		iamClient, err := b.IAM(s)
		if err != nil {
			return nil, err
		}

		serviceAccountId, err := util.ParseServiceAccountFromIAMJWT(rawToken)
		if err != nil {
			return nil, err
		}

		accountKey, err := util.ServiceAccountKey(iamClient, keyId, serviceAccountId, role.ProjectId)
		if err != nil {
			return nil, err
		}

		return util.PublicKey(accountKey.PublicKeyData)
	case gceRoleType:
		var certsEndpoint string
		conf, err := b.config(s)
		if err != nil {
			return nil, fmt.Errorf("could not read config for backend: %v", err)
		}
		if conf != nil {
			certsEndpoint = conf.GoogleCertsEndpoint
		}

		key, err := util.OAuth2RSAPublicKey(keyId, certsEndpoint)
		if err != nil {
			return nil, err
		}
		return key, nil
	default:
		return nil, fmt.Errorf("unexpected role type %s", role.RoleType)
	}
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
func (b *GcpAuthBackend) pathIamLogin(req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	iamClient, err := b.IAM(req.Storage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(clientErrorTemplate, "IAM", err)), nil
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
	serviceAccount, err := util.ServiceAccount(iamClient, loginInfo.ServiceAccountId, role.ProjectId)
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
			},
		},
	}

	return resp, nil
}

// pathIamRenew returns an error if the service account referenced in the auth token metadata cannot renew the
// auth token for the given role.
func (b *GcpAuthBackend) pathIamRenew(req *logical.Request, roleName string, role *gcpRole) error {
	iamClient, err := b.IAM(req.Storage)
	if err != nil {
		return fmt.Errorf(clientErrorTemplate, "IAM", err)
	}

	serviceAccountId, ok := req.Auth.Metadata["service_account_id"]
	if !ok {
		return errors.New("service account id metadata not associated with auth token, invalid")
	}

	serviceAccount, err := util.ServiceAccount(iamClient, serviceAccountId, role.ProjectId)
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
	// This is just in case - project should already be used to retrieve service account.
	if role.ProjectId != serviceAccount.ProjectId {
		return fmt.Errorf("service account %s does not belong to project %s", serviceAccount.Email, role.ProjectId)
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
func (b *GcpAuthBackend) pathGceLogin(req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	role := loginInfo.Role
	metadata := loginInfo.GceMetadata
	if metadata == nil {
		return logical.ErrorResponse("could not get GCE metadata from given JWT"), nil
	}

	if role.ProjectId != metadata.ProjectId {
		return logical.ErrorResponse(fmt.Sprintf(
			"GCE instance must belong to project %s; metadata given has project %s",
			role.ProjectId, metadata.ProjectId)), nil
	}

	// Verify instance exists.
	gceClient, err := b.GCE(req.Storage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(clientErrorTemplate, "GCE", err)), nil
	}

	instance, err := metadata.GetVerifiedInstance(gceClient)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"error when attempting to find instance (project %s, zone: %s, instance: %s) :%v",
			metadata.ProjectId, metadata.Zone, metadata.InstanceName, err)), nil
	}

	if err := b.authorizeGCEInstance(instance, req.Storage, role, metadata.Zone, loginInfo.ServiceAccountId); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	iamClient, err := b.IAM(req.Storage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(clientErrorTemplate, "IAM", err)), nil
	}

	serviceAccount, err := util.ServiceAccount(iamClient, loginInfo.ServiceAccountId, loginInfo.Role.ProjectId)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(
			"Could not find service account '%s' used for GCE metadata token: %s",
			loginInfo.ServiceAccountId, err)), nil
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
			},
		},
	}

	return resp, nil
}

func authMetadata(loginInfo *gcpLoginInfo, serviceAccount *iam.ServiceAccount) map[string]string {
	metadata := map[string]string{
		"role":                  loginInfo.RoleName,
		"service_account_id":    serviceAccount.UniqueId,
		"service_account_email": serviceAccount.Email,
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
func (b *GcpAuthBackend) pathGceRenew(req *logical.Request, roleName string, role *gcpRole) error {
	gceClient, err := b.GCE(req.Storage)
	if err != nil {
		return fmt.Errorf(clientErrorTemplate, "GCE", err)
	}

	meta, err := util.GetInstanceMetadataFromAuth(req.Auth.Metadata)
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
	if err := b.authorizeGCEInstance(instance, req.Storage, role, meta.Zone, serviceAccountId); err != nil {
		return fmt.Errorf("could not renew token for role %s: %v", roleName, err)
	}

	return nil
}

// validateGCEInstance returns an error if the given GCE instance is not authorized for the role.
func (b *GcpAuthBackend) authorizeGCEInstance(instance *compute.Instance, s logical.Storage, role *gcpRole, zone, serviceAccountId string) error {
	gceClient, err := b.GCE(s)
	if err != nil {
		return err
	}

	// Verify instance has role labels if labels were set on role.
	for k, expectedV := range role.BoundLabels {
		actualV, ok := instance.Labels[k]
		if !ok || actualV != expectedV {
			return fmt.Errorf("role label '%s:%s' not found on GCE instance", k, expectedV)
		}
	}

	// Verify that instance is in zone or region if given.
	if len(role.BoundZone) > 0 {
		var zone string
		idx := strings.LastIndex(instance.Zone, "zones/")
		if idx > 0 {
			// Parse zone name from full zone self-link URL.
			idx += len("zones/")
			zone = instance.Zone[idx:len(instance.Zone)]
		} else {
			// Expect full zone name to be set as instance zone.
			zone = instance.Zone
		}

		if zone != role.BoundZone {
			return fmt.Errorf("instance is not in role zone '%s'", role.BoundZone)
		}
	} else if len(role.BoundRegion) > 0 {
		zone, err := gceClient.Zones.Get(role.ProjectId, zone).Do()
		if err != nil {
			return fmt.Errorf("could not verify instance zone '%s' is available for project '%s': %v", role.ProjectId, zone, err)
		}
		if zone.Region != role.BoundRegion {
			return fmt.Errorf("zone '%s' is not in region '%s'", zone.Name, zone.Region)
		}
	}

	// If instance group is given, verify group exists and that instance is in group.
	if len(role.BoundInstanceGroup) > 0 {
		var group *compute.InstanceGroup
		var err error

		// Check if group should be zonal or regional.
		if len(role.BoundZone) > 0 {
			group, err = gceClient.InstanceGroups.Get(role.ProjectId, role.BoundZone, role.BoundInstanceGroup).Do()
			if err != nil {
				return fmt.Errorf("could not find role instance group %s (project %s, zone %s)", role.BoundInstanceGroup, role.ProjectId, role.BoundZone)
			}
		} else if len(role.BoundRegion) > 0 {
			group, err = gceClient.RegionInstanceGroups.Get(role.ProjectId, role.BoundRegion, role.BoundInstanceGroup).Do()
			if err != nil {
				return fmt.Errorf("could not find role instance group %s (project %s, region %s)", role.BoundInstanceGroup, role.ProjectId, role.BoundRegion)
			}
		} else {
			return errors.New("expected zone or region to be set for GCE role '%s' with instance group")
		}

		// Verify instance group contains authenticating instance.
		instanceIdFilter := fmt.Sprintf("instance eq %s", instance.SelfLink)
		listInstanceReq := &compute.InstanceGroupsListInstancesRequest{}
		listResp, err := gceClient.InstanceGroups.ListInstances(role.ProjectId, role.BoundZone, group.Name, listInstanceReq).Filter(instanceIdFilter).Do()
		if err != nil {
			return fmt.Errorf("could not confirm instance %s is part of instance group %s: %s", instance.Name, role.BoundInstanceGroup, err)
		}

		if len(listResp.Items) == 0 {
			return fmt.Errorf("instance %s is not part of instance group %s", instance.Name, role.BoundInstanceGroup)
		}

	}

	// Verify instance is running under one of the allowed service accounts.
	if len(role.BoundServiceAccounts) > 0 {
		iamClient, err := b.IAM(s)
		if err != nil {
			return err
		}

		serviceAccount, err := util.ServiceAccount(iamClient, serviceAccountId, role.ProjectId)
		if err != nil {
			return fmt.Errorf("could not find service account with id '%s': %v", serviceAccountId, err)
		}

		if !(strutil.StrListContains(role.BoundServiceAccounts, serviceAccount.Email) ||
			strutil.StrListContains(role.BoundServiceAccounts, serviceAccount.UniqueId)) {
			return fmt.Errorf("GCE instance's service account email (%s) or id (%s) not found in role service accounts: %v",
				serviceAccount.Email, serviceAccount.UniqueId, role.BoundServiceAccounts)
		}
	}

	return nil
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
