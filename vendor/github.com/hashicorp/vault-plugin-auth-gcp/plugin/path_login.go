package gcpauth

import (
	"errors"
	"fmt"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/hashicorp/vault-plugin-auth-gcp/plugin/util"
	"github.com/hashicorp/vault/helper/policyutil"
	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"google.golang.org/api/iam/v1"
)

const (
	expectedJwtAudTemplate string = "vault/%s"

	// Default duration that JWT tokens must expire within to be accepted
	defaultMaxJwtExpMin int = 15

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
				Type:        framework.TypeString,
				Description: `A signed JWT for authenticating a service account.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:           b.pathLogin,
			logical.PersonaLookaheadOperation: b.pathLogin,
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

func (b *GcpAuthBackend) pathLogin(req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	loginInfo, err := b.parseInfoFromJwt(req, data)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	roleType := loginInfo.Role.RoleType
	switch roleType {
	case iamRoleType:
		return b.pathIamLogin(req, loginInfo)
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
		if err := b.pathIamRenew(req, role); err != nil {
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

	// ID of the public key to verify the signed JWT.
	KeyId string

	// Signed JWT
	JWT jwt.JWT
}

func (b *GcpAuthBackend) parseInfoFromJwt(req *logical.Request, data *framework.FieldData) (*gcpLoginInfo, error) {
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

	signedJwt, ok := data.GetOk("jwt")
	if !ok {
		return nil, errors.New("jwt argument is required")
	}
	signedJwtBytes := []byte(signedJwt.(string))

	// Parse into JWS to get header values.
	jwsVal, err := jws.Parse(signedJwtBytes)
	if err != nil {
		return nil, err
	}
	headerVal := jwsVal.Protected()

	if headerVal.Has("kid") {
		loginInfo.KeyId = jwsVal.Protected().Get("kid").(string)
	} else {
		return nil, errors.New("provided JWT must have 'kid' header value")
	}

	// Parse claims
	loginInfo.JWT, err = jws.ParseJWT(signedJwtBytes)
	if err != nil {
		return nil, err
	}

	sub, ok := loginInfo.JWT.Claims().Subject()
	if !ok {
		return nil, errors.New("expected JWT to have 'sub' claim with service account id or email")
	}
	loginInfo.ServiceAccountId = sub

	return loginInfo, nil
}

func (info *gcpLoginInfo) validateJWT(keyPEM string, loginInfo *gcpLoginInfo) error {
	pubKey, err := util.PublicKey(keyPEM)
	if err != nil {
		return err
	}

	validator := &jwt.Validator{
		Expected: jwt.Claims{
			"aud": fmt.Sprintf(expectedJwtAudTemplate, loginInfo.RoleName),
		},
		Fn: func(c jwt.Claims) error {
			exp, ok := c.Expiration()
			if !ok {
				return errors.New("JWT claim 'exp' is required")
			}
			if exp.After(time.Now().Add(loginInfo.Role.MaxJwtExp)) {
				return fmt.Errorf("JWT expires in %v minutes but must expire within %v for this role. Please generate a new token with a valid expiration.",
					int(exp.Sub(time.Now())/time.Minute), loginInfo.Role.MaxJwtExp)
			}

			return nil
		},
	}

	if err := info.JWT.Validate(pubKey, crypto.SigningMethodRS256, validator); err != nil {
		return fmt.Errorf("invalid JWT: %v", err)
	}

	return nil
}

// ---- IAM login domain ----

func (b *GcpAuthBackend) pathIamLogin(req *logical.Request, loginInfo *gcpLoginInfo) (*logical.Response, error) {
	iamClient, err := b.IAM(req.Storage)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf(clientErrorTemplate, "IAM", err)), nil
	}

	role := loginInfo.Role

	// Verify and get service account from signed JWT.
	key, err := util.ServiceAccountKey(iamClient, loginInfo.KeyId, loginInfo.ServiceAccountId, role.ProjectId)
	if err != nil {
		return logical.ErrorResponse(fmt.Sprintf("service account %s has no key with id %s", loginInfo.ServiceAccountId, loginInfo.KeyId)), nil
	}

	if err := loginInfo.validateJWT(key.PublicKeyData, loginInfo); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	serviceAccount, err := util.ServiceAccount(iamClient, loginInfo.ServiceAccountId, role.ProjectId)
	if err != nil {
		return nil, err
	}
	if serviceAccount == nil {
		return nil, errors.New("service account is empty")
	}

	if req.Operation == logical.PersonaLookaheadOperation {
		return &logical.Response{
			Auth: &logical.Auth{
				Persona: &logical.Persona{
					Name: serviceAccount.UniqueId,
				},
			},
		}, nil
	}

	// Validate service account can login against role.
	if err := b.validateAgainstIAMRole(serviceAccount, role); err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Period: role.Period,
			Persona: &logical.Persona{
				Name: serviceAccount.UniqueId,
			},
			Policies: role.Policies,
			Metadata: map[string]string{
				"service_account_id":    serviceAccount.UniqueId,
				"service_account_email": serviceAccount.Email,
				"role":                  loginInfo.RoleName,
			},
			DisplayName: serviceAccount.Email,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
			},
		},
	}

	return resp, nil
}

func (b *GcpAuthBackend) pathIamRenew(req *logical.Request, role *gcpRole) error {
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

	if err := b.validateAgainstIAMRole(serviceAccount, role); err != nil {
		return errors.New("service account is no longer authorized for role")
	}

	return nil
}

// validateAgainstIAMRole returns an error if the given IAM service account is not authorized for the role.
func (b *GcpAuthBackend) validateAgainstIAMRole(serviceAccount *iam.ServiceAccount, role *gcpRole) error {
	// This is just in case - project should already be used to retrieve service account.
	if role.ProjectId != serviceAccount.ProjectId {
		return fmt.Errorf("service account %s does not belong to project %s", serviceAccount.Email, role.ProjectId)
	}

	// Check if role has the wildcard as the only service account.
	if len(role.ServiceAccounts) == 1 && role.ServiceAccounts[0] == serviceAccountWildcard {
		return nil
	}

	// Check for service account id/email.
	if strutil.StrListContains(role.ServiceAccounts, serviceAccount.Email) ||
		strutil.StrListContains(role.ServiceAccounts, serviceAccount.UniqueId) {
		return nil
	}

	return fmt.Errorf("service account %s (id: %s) is not authorized for role",
		serviceAccount.Email, serviceAccount.UniqueId)
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
