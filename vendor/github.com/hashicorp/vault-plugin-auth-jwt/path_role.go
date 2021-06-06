package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-sockaddr"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/helper/tokenutil"
	"github.com/hashicorp/vault/sdk/logical"
	"gopkg.in/square/go-jose.v2/jwt"
)

var reservedMetadata = []string{"role"}

const (
	claimDefaultLeeway    = 150
	boundClaimsTypeString = "string"
	boundClaimsTypeGlob   = "glob"
)

func pathRoleList(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "role/?",
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ListOperation: &framework.PathOperation{
				Callback:    b.pathRoleList,
				Summary:     strings.TrimSpace(roleHelp["role-list"][0]),
				Description: strings.TrimSpace(roleHelp["role-list"][1]),
			},
		},
		HelpSynopsis:    strings.TrimSpace(roleHelp["role-list"][0]),
		HelpDescription: strings.TrimSpace(roleHelp["role-list"][1]),
	}
}

// pathRole returns the path configurations for the CRUD operations on roles
func pathRole(b *jwtAuthBackend) *framework.Path {
	p := &framework.Path{
		Pattern: "role/" + framework.GenericNameRegex("name"),
		Fields: map[string]*framework.FieldSchema{
			"name": {
				Type:        framework.TypeLowerCaseString,
				Description: "Name of the role.",
			},
			"role_type": {
				Type:        framework.TypeString,
				Description: "Type of the role, either 'jwt' or 'oidc'.",
			},

			"policies": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_policies"),
				Deprecated:  true,
			},
			"num_uses": {
				Type:        framework.TypeInt,
				Description: tokenutil.DeprecationText("token_num_uses"),
				Deprecated:  true,
			},
			"ttl": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_ttl"),
				Deprecated:  true,
			},
			"max_ttl": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_max_ttl"),
				Deprecated:  true,
			},
			"period": {
				Type:        framework.TypeDurationSecond,
				Description: tokenutil.DeprecationText("token_period"),
				Deprecated:  true,
			},
			"bound_cidrs": {
				Type:        framework.TypeCommaStringSlice,
				Description: tokenutil.DeprecationText("token_bound_cidrs"),
				Deprecated:  true,
			},
			"expiration_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating expiration of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
				Default: claimDefaultLeeway,
			},
			"not_before_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating not before values of a token to account for clock skew. 
Defaults to 150 (2.5 minutes) if set to 0 and can be disabled if set to -1.`,
				Default: claimDefaultLeeway,
			},
			"clock_skew_leeway": {
				Type: framework.TypeSignedDurationSecond,
				Description: `Duration in seconds of leeway when validating all claims to account for clock skew. 
Defaults to 60 (1 minute) if set to 0 and can be disabled if set to -1.`,
				Default: jwt.DefaultLeeway,
			},
			"bound_subject": {
				Type:        framework.TypeString,
				Description: `The 'sub' claim that is valid for login. Optional.`,
			},
			"bound_audiences": {
				Type:        framework.TypeCommaStringSlice,
				Description: `Comma-separated list of 'aud' claims that are valid for login; any match is sufficient`,
			},
			"bound_claims_type": {
				Type:        framework.TypeString,
				Description: `How to interpret values in the map of claims/values (which must match for login): allowed values are 'string' or 'glob'`,
				Default:     boundClaimsTypeString,
			},
			"bound_claims": {
				Type:        framework.TypeMap,
				Description: `Map of claims/values which must match for login`,
			},
			"claim_mappings": {
				Type:        framework.TypeKVPairs,
				Description: `Mappings of claims (key) that will be copied to a metadata field (value)`,
			},
			"user_claim": {
				Type:        framework.TypeString,
				Description: `The claim to use for the Identity entity alias name`,
			},
			"groups_claim": {
				Type:        framework.TypeString,
				Description: `The claim to use for the Identity group alias names`,
			},
			"oidc_scopes": {
				Type:        framework.TypeCommaStringSlice,
				Description: `Comma-separated list of OIDC scopes`,
			},
			"allowed_redirect_uris": {
				Type:        framework.TypeCommaStringSlice,
				Description: `Comma-separated list of allowed values for redirect_uri`,
			},
			"verbose_oidc_logging": {
				Type: framework.TypeBool,
				Description: `Log received OIDC tokens and claims when debug-level logging is active. 
Not recommended in production since sensitive information may be present 
in OIDC responses.`,
			},
			"max_age": {
				Type: framework.TypeDurationSecond,
				Description: `Specifies the allowable elapsed time in seconds since the last time the 
user was actively authenticated.`,
			},
		},
		ExistenceCheck: b.pathRoleExistenceCheck,
		Operations: map[logical.Operation]framework.OperationHandler{
			logical.ReadOperation: &framework.PathOperation{
				Callback: b.pathRoleRead,
				Summary:  "Read an existing role.",
			},

			logical.UpdateOperation: &framework.PathOperation{
				Callback:    b.pathRoleCreateUpdate,
				Summary:     strings.TrimSpace(roleHelp["role"][0]),
				Description: strings.TrimSpace(roleHelp["role"][1]),
			},

			logical.CreateOperation: &framework.PathOperation{
				Callback:    b.pathRoleCreateUpdate,
				Summary:     strings.TrimSpace(roleHelp["role"][0]),
				Description: strings.TrimSpace(roleHelp["role"][1]),
			},

			logical.DeleteOperation: &framework.PathOperation{
				Callback: b.pathRoleDelete,
				Summary:  "Delete an existing role.",
			},
		},
		HelpSynopsis:    strings.TrimSpace(roleHelp["role"][0]),
		HelpDescription: strings.TrimSpace(roleHelp["role"][1]),
	}

	tokenutil.AddTokenFields(p.Fields)
	return p
}

type jwtRole struct {
	tokenutil.TokenParams

	RoleType string `json:"role_type"`

	// Duration of leeway for expiration to account for clock skew
	ExpirationLeeway time.Duration `json:"expiration_leeway"`

	// Duration of leeway for not before to account for clock skew
	NotBeforeLeeway time.Duration `json:"not_before_leeway"`

	// Duration of leeway for all claims to account for clock skew
	ClockSkewLeeway time.Duration `json:"clock_skew_leeway"`

	// Role binding properties
	BoundAudiences      []string               `json:"bound_audiences"`
	BoundSubject        string                 `json:"bound_subject"`
	BoundClaimsType     string                 `json:"bound_claims_type"`
	BoundClaims         map[string]interface{} `json:"bound_claims"`
	ClaimMappings       map[string]string      `json:"claim_mappings"`
	UserClaim           string                 `json:"user_claim"`
	GroupsClaim         string                 `json:"groups_claim"`
	OIDCScopes          []string               `json:"oidc_scopes"`
	AllowedRedirectURIs []string               `json:"allowed_redirect_uris"`
	VerboseOIDCLogging  bool                   `json:"verbose_oidc_logging"`
	MaxAge              time.Duration          `json:"max_age"`

	// Deprecated by TokenParams
	Policies   []string                      `json:"policies"`
	NumUses    int                           `json:"num_uses"`
	TTL        time.Duration                 `json:"ttl"`
	MaxTTL     time.Duration                 `json:"max_ttl"`
	Period     time.Duration                 `json:"period"`
	BoundCIDRs []*sockaddr.SockAddrMarshaler `json:"bound_cidrs"`
}

// role takes a storage backend and the name and returns the role's storage
// entry
func (b *jwtAuthBackend) role(ctx context.Context, s logical.Storage, name string) (*jwtRole, error) {
	raw, err := s.Get(ctx, rolePrefix+name)
	if err != nil {
		return nil, err
	}
	if raw == nil {
		return nil, nil
	}

	role := new(jwtRole)
	if err := raw.DecodeJSON(role); err != nil {
		return nil, err
	}

	// Report legacy roles as type "jwt"
	if role.RoleType == "" {
		role.RoleType = "jwt"
	}

	if role.BoundClaimsType == "" {
		role.BoundClaimsType = boundClaimsTypeString
	}

	if role.TokenTTL == 0 && role.TTL > 0 {
		role.TokenTTL = role.TTL
	}
	if role.TokenMaxTTL == 0 && role.MaxTTL > 0 {
		role.TokenMaxTTL = role.MaxTTL
	}
	if role.TokenPeriod == 0 && role.Period > 0 {
		role.TokenPeriod = role.Period
	}
	if role.TokenNumUses == 0 && role.NumUses > 0 {
		role.TokenNumUses = role.NumUses
	}
	if len(role.TokenPolicies) == 0 && len(role.Policies) > 0 {
		role.TokenPolicies = role.Policies
	}
	if len(role.TokenBoundCIDRs) == 0 && len(role.BoundCIDRs) > 0 {
		role.TokenBoundCIDRs = role.BoundCIDRs
	}

	return role, nil
}

// pathRoleExistenceCheck returns whether the role with the given name exists or not.
func (b *jwtAuthBackend) pathRoleExistenceCheck(ctx context.Context, req *logical.Request, data *framework.FieldData) (bool, error) {
	role, err := b.role(ctx, req.Storage, data.Get("name").(string))
	if err != nil {
		return false, err
	}
	return role != nil, nil
}

// pathRoleList is used to list all the Roles registered with the backend.
func (b *jwtAuthBackend) pathRoleList(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, rolePrefix)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

// pathRoleRead grabs a read lock and reads the options set on the role from the storage
func (b *jwtAuthBackend) pathRoleRead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing name"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	// Create a map of data to be returned
	d := map[string]interface{}{
		"role_type":             role.RoleType,
		"expiration_leeway":     int64(role.ExpirationLeeway.Seconds()),
		"not_before_leeway":     int64(role.NotBeforeLeeway.Seconds()),
		"clock_skew_leeway":     int64(role.ClockSkewLeeway.Seconds()),
		"bound_audiences":       role.BoundAudiences,
		"bound_subject":         role.BoundSubject,
		"bound_claims_type":     role.BoundClaimsType,
		"bound_claims":          role.BoundClaims,
		"claim_mappings":        role.ClaimMappings,
		"user_claim":            role.UserClaim,
		"groups_claim":          role.GroupsClaim,
		"allowed_redirect_uris": role.AllowedRedirectURIs,
		"oidc_scopes":           role.OIDCScopes,
		"verbose_oidc_logging":  role.VerboseOIDCLogging,
		"max_age":               int64(role.MaxAge.Seconds()),
	}

	role.PopulateTokenData(d)

	if len(role.Policies) > 0 {
		d["policies"] = d["token_policies"]
	}
	if len(role.BoundCIDRs) > 0 {
		d["bound_cidrs"] = d["token_bound_cidrs"]
	}
	if role.TTL > 0 {
		d["ttl"] = int64(role.TTL.Seconds())
	}
	if role.MaxTTL > 0 {
		d["max_ttl"] = int64(role.MaxTTL.Seconds())
	}
	if role.Period > 0 {
		d["period"] = int64(role.Period.Seconds())
	}
	if role.NumUses > 0 {
		d["num_uses"] = role.NumUses
	}

	return &logical.Response{
		Data: d,
	}, nil
}

// pathRoleDelete removes the role from storage
func (b *jwtAuthBackend) pathRoleDelete(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("role name required"), nil
	}

	// Delete the role itself
	if err := req.Storage.Delete(ctx, rolePrefix+roleName); err != nil {
		return nil, err
	}

	return nil, nil
}

// pathRoleCreateUpdate registers a new role with the backend or updates the options
// of an existing role
func (b *jwtAuthBackend) pathRoleCreateUpdate(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := data.Get("name").(string)
	if roleName == "" {
		return logical.ErrorResponse("missing role name"), nil
	}

	// Check if the role already exists
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}

	// Create a new entry object if this is a CreateOperation
	if role == nil {
		if req.Operation == logical.UpdateOperation {
			return nil, errors.New("role entry not found during update operation")
		}
		role = new(jwtRole)
	}

	roleType := data.Get("role_type").(string)
	if roleType == "" {
		roleType = "oidc"
	}
	if roleType != "jwt" && roleType != "oidc" {
		return logical.ErrorResponse("invalid 'role_type': %s", roleType), nil
	}
	role.RoleType = roleType

	if err := role.ParseTokenFields(req, data); err != nil {
		return logical.ErrorResponse(err.Error()), logical.ErrInvalidRequest
	}

	// Handle upgrade cases
	{
		if err := tokenutil.UpgradeValue(data, "policies", "token_policies", &role.Policies, &role.TokenPolicies); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "bound_cidrs", "token_bound_cidrs", &role.BoundCIDRs, &role.TokenBoundCIDRs); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "num_uses", "token_num_uses", &role.NumUses, &role.TokenNumUses); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "ttl", "token_ttl", &role.TTL, &role.TokenTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "max_ttl", "token_max_ttl", &role.MaxTTL, &role.TokenMaxTTL); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		if err := tokenutil.UpgradeValue(data, "period", "token_period", &role.Period, &role.TokenPeriod); err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}
	}

	if role.TokenPeriod > b.System().MaxLeaseTTL() {
		return logical.ErrorResponse(fmt.Sprintf("'period' of '%q' is greater than the backend's maximum lease TTL of '%q'", role.TokenPeriod.String(), b.System().MaxLeaseTTL().String())), nil
	}

	if tokenExpLeewayRaw, ok := data.GetOk("expiration_leeway"); ok {
		role.ExpirationLeeway = time.Duration(tokenExpLeewayRaw.(int)) * time.Second
	}

	if tokenNotBeforeLeewayRaw, ok := data.GetOk("not_before_leeway"); ok {
		role.NotBeforeLeeway = time.Duration(tokenNotBeforeLeewayRaw.(int)) * time.Second
	}

	if tokenClockSkewLeeway, ok := data.GetOk("clock_skew_leeway"); ok {
		role.ClockSkewLeeway = time.Duration(tokenClockSkewLeeway.(int)) * time.Second
	}

	if boundAudiences, ok := data.GetOk("bound_audiences"); ok {
		role.BoundAudiences = boundAudiences.([]string)
	}

	if boundSubject, ok := data.GetOk("bound_subject"); ok {
		role.BoundSubject = boundSubject.(string)
	}

	if verboseOIDCLoggingRaw, ok := data.GetOk("verbose_oidc_logging"); ok {
		role.VerboseOIDCLogging = verboseOIDCLoggingRaw.(bool)
	}

	if maxAgeRaw, ok := data.GetOk("max_age"); ok {
		role.MaxAge = time.Duration(maxAgeRaw.(int)) * time.Second
	}

	boundClaimsType := data.Get("bound_claims_type").(string)
	switch boundClaimsType {
	case boundClaimsTypeString, boundClaimsTypeGlob:
		role.BoundClaimsType = boundClaimsType
	default:
		return logical.ErrorResponse("invalid 'bound_claims_type': %s", boundClaimsType), nil
	}

	if boundClaimsRaw, ok := data.GetOk("bound_claims"); ok {
		role.BoundClaims = boundClaimsRaw.(map[string]interface{})

		if boundClaimsType == boundClaimsTypeGlob {
			// Check that the claims are all strings
			for _, claimValues := range role.BoundClaims {
				claimsValuesList, ok := normalizeList(claimValues)

				if !ok {
					return logical.ErrorResponse("claim is not a string or list: %v", claimValues), nil
				}

				for _, claimValue := range claimsValuesList {
					if _, ok := claimValue.(string); !ok {
						return logical.ErrorResponse("claim is not a string: %v", claimValue), nil
					}
				}
			}
		}
	}

	if claimMappingsRaw, ok := data.GetOk("claim_mappings"); ok {
		claimMappings := claimMappingsRaw.(map[string]string)

		// sanity check mappings for duplicates and collision with reserved names
		targets := make(map[string]bool)
		for _, metadataKey := range claimMappings {
			if strutil.StrListContains(reservedMetadata, metadataKey) {
				return logical.ErrorResponse("metadata key %q is reserved and may not be a mapping destination", metadataKey), nil
			}

			if targets[metadataKey] {
				return logical.ErrorResponse("multiple keys are mapped to metadata key %q", metadataKey), nil
			}
			targets[metadataKey] = true
		}

		role.ClaimMappings = claimMappings
	}

	if userClaim, ok := data.GetOk("user_claim"); ok {
		role.UserClaim = userClaim.(string)
	}
	if role.UserClaim == "" {
		return logical.ErrorResponse("a user claim must be defined on the role"), nil
	}

	if groupsClaim, ok := data.GetOk("groups_claim"); ok {
		role.GroupsClaim = groupsClaim.(string)
	}

	if oidcScopes, ok := data.GetOk("oidc_scopes"); ok {
		role.OIDCScopes = oidcScopes.([]string)
	}

	if allowedRedirectURIs, ok := data.GetOk("allowed_redirect_uris"); ok {
		role.AllowedRedirectURIs = allowedRedirectURIs.([]string)
	}

	if role.RoleType == "oidc" && len(role.AllowedRedirectURIs) == 0 {
		return logical.ErrorResponse(
			"'allowed_redirect_uris' must be set if 'role_type' is 'oidc' or unspecified."), nil
	}

	// OIDC verification will enforce that the audience match the configured client_id.
	// For other methods, require at least one bound constraint.
	if roleType != "oidc" {
		if len(role.BoundAudiences) == 0 &&
			len(role.TokenBoundCIDRs) == 0 &&
			role.BoundSubject == "" &&
			len(role.BoundClaims) == 0 {
			return logical.ErrorResponse("must have at least one bound constraint when creating/updating a role"), nil
		}
	}

	// Check that the TTL value provided is less than the MaxTTL.
	// Sanitizing the TTL and MaxTTL is not required now and can be performed
	// at credential issue time.
	if role.TokenMaxTTL > 0 && role.TokenTTL > role.TokenMaxTTL {
		return logical.ErrorResponse("ttl should not be greater than max ttl"), nil
	}

	resp := &logical.Response{}
	if role.TokenMaxTTL > b.System().MaxLeaseTTL() {
		resp.AddWarning("token max ttl is greater than the system or backend mount's maximum TTL value; issued tokens' max TTL value will be truncated")
	}

	if role.VerboseOIDCLogging {
		resp.AddWarning(`verbose_oidc_logging has been enabled for this role. ` +
			`This is not recommended in production since sensitive information ` +
			`may be present in OIDC responses.`)
	}

	// Store the entry.
	entry, err := logical.StorageEntryJSON(rolePrefix+roleName, role)
	if err != nil {
		return nil, err
	}
	if err = req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return resp, nil
}

// roleStorageEntry stores all the options that are set on an role
var roleHelp = map[string][2]string{
	"role-list": {
		"Lists all the roles registered with the backend.",
		"The list will contain the names of the roles.",
	},
	"role": {
		"Register an role with the backend.",
		`A role is required to authenticate with this backend. The role binds
		JWT token information with token policies and settings.
		The bindings, token polices and token settings can all be configured
		using this endpoint`,
	},
}
