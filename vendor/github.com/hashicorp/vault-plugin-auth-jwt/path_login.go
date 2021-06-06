package jwtauth

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/cap/jwt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/logical"
	"golang.org/x/oauth2"
)

func pathLogin(b *jwtAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: `login$`,
		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeLowerCaseString,
				Description: "The role to log in against.",
			},
			"jwt": {
				Type:        framework.TypeString,
				Description: "The signed JWT to validate.",
			},
		},

		Operations: map[logical.Operation]framework.OperationHandler{
			logical.UpdateOperation: &framework.PathOperation{
				Callback: b.pathLogin,
				Summary:  pathLoginHelpSyn,
			},
			logical.AliasLookaheadOperation: &framework.PathOperation{
				Callback: b.pathLogin,
			},
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

func (b *jwtAuthBackend) pathLogin(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	config, err := b.config(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return logical.ErrorResponse("could not load configuration"), nil
	}

	roleName := d.Get("role").(string)
	if roleName == "" {
		roleName = config.DefaultRole
	}
	if roleName == "" {
		return logical.ErrorResponse("missing role"), nil
	}

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q could not be found", roleName), nil
	}

	if role.RoleType == "oidc" {
		return logical.ErrorResponse("role with oidc role_type is not allowed"), nil
	}

	token := d.Get("jwt").(string)
	if len(token) == 0 {
		return logical.ErrorResponse("missing token"), nil
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

	// Get the JWT validator based on the configured auth type
	validator, err := b.jwtValidator(config)
	if err != nil {
		return logical.ErrorResponse("error configuring token validator: %s", err.Error()), nil
	}

	// Validate JWT supported algorithms if they've been provided. Otherwise,
	// ensure that the signing algorithm is a member of the supported set.
	signingAlgorithms := toAlg(config.JWTSupportedAlgs)
	if len(signingAlgorithms) == 0 {
		signingAlgorithms = []jwt.Alg{
			jwt.RS256, jwt.RS384, jwt.RS512, jwt.ES256, jwt.ES384,
			jwt.ES512, jwt.PS256, jwt.PS384, jwt.PS512, jwt.EdDSA,
		}
	}

	// Set expected claims values to assert on the JWT
	expected := jwt.Expected{
		Issuer:            config.BoundIssuer,
		Subject:           role.BoundSubject,
		Audiences:         role.BoundAudiences,
		SigningAlgorithms: signingAlgorithms,
		NotBeforeLeeway:   role.NotBeforeLeeway,
		ExpirationLeeway:  role.ExpirationLeeway,
		ClockSkewLeeway:   role.ClockSkewLeeway,
	}

	// Validate the JWT by verifying its signature and asserting expected claims values
	allClaims, err := validator.Validate(ctx, token, expected)
	if err != nil {
		return logical.ErrorResponse("error validating token: %s", err.Error()), nil
	}

	// If there are no bound audiences for the role, then the existence of any audience
	// in the audience claim should result in an error.
	aud, ok := getClaim(b.Logger(), allClaims, "aud").([]interface{})
	if ok && len(aud) > 0 && len(role.BoundAudiences) == 0 {
		return logical.ErrorResponse("audience claim found in JWT but no audiences bound to the role"), nil
	}

	alias, groupAliases, err := b.createIdentity(ctx, allClaims, role, nil)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	if err := validateBoundClaims(b.Logger(), role.BoundClaimsType, role.BoundClaims, allClaims); err != nil {
		return logical.ErrorResponse("error validating claims: %s", err.Error()), nil
	}

	tokenMetadata := map[string]string{"role": roleName}
	for k, v := range alias.Metadata {
		tokenMetadata[k] = v
	}

	auth := &logical.Auth{
		DisplayName:  alias.Name,
		Alias:        alias,
		GroupAliases: groupAliases,
		InternalData: map[string]interface{}{
			"role": roleName,
		},
		Metadata: tokenMetadata,
	}

	role.PopulateTokenAuth(auth)

	return &logical.Response{
		Auth: auth,
	}, nil
}

func (b *jwtAuthBackend) pathLoginRenew(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName := req.Auth.InternalData["role"].(string)
	if roleName == "" {
		return nil, errors.New("failed to fetch role_name during renewal")
	}

	// Ensure that the Role still exists.
	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, errwrap.Wrapf(fmt.Sprintf("failed to validate role %s during renewal: {{err}}", roleName), err)
	}
	if role == nil {
		return nil, fmt.Errorf("role %s does not exist during renewal", roleName)
	}

	resp := &logical.Response{Auth: req.Auth}
	resp.Auth.TTL = role.TokenTTL
	resp.Auth.MaxTTL = role.TokenMaxTTL
	resp.Auth.Period = role.TokenPeriod
	return resp, nil
}

// createIdentity creates an alias and set of groups aliases based on the role
// definition and received claims.
func (b *jwtAuthBackend) createIdentity(ctx context.Context, allClaims map[string]interface{}, role *jwtRole, tokenSource oauth2.TokenSource) (*logical.Alias, []*logical.Alias, error) {
	userClaimRaw, ok := allClaims[role.UserClaim]
	if !ok {
		return nil, nil, fmt.Errorf("claim %q not found in token", role.UserClaim)
	}
	userName, ok := userClaimRaw.(string)
	if !ok {
		return nil, nil, fmt.Errorf("claim %q could not be converted to string", role.UserClaim)
	}

	pConfig, err := NewProviderConfig(ctx, b.cachedConfig, ProviderMap())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load custom provider config: %s", err)
	}

	if err := b.fetchUserInfo(ctx, pConfig, allClaims, role); err != nil {
		return nil, nil, err
	}

	metadata, err := extractMetadata(b.Logger(), allClaims, role.ClaimMappings)
	if err != nil {
		return nil, nil, err
	}

	alias := &logical.Alias{
		Name:     userName,
		Metadata: metadata,
	}

	var groupAliases []*logical.Alias

	if role.GroupsClaim == "" {
		return alias, groupAliases, nil
	}

	groupsClaimRaw, err := b.fetchGroups(ctx, pConfig, allClaims, role, tokenSource)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to fetch groups: %s", err)
	}

	groups, ok := normalizeList(groupsClaimRaw)

	if !ok {
		return nil, nil, fmt.Errorf("%q claim could not be converted to string list", role.GroupsClaim)
	}
	for _, groupRaw := range groups {
		group, ok := groupRaw.(string)
		if !ok {
			return nil, nil, fmt.Errorf("value %v in groups claim could not be parsed as string", groupRaw)
		}
		if group == "" {
			continue
		}
		groupAliases = append(groupAliases, &logical.Alias{
			Name: group,
		})
	}

	return alias, groupAliases, nil
}

// Checks if there's a custom provider_config and calls FetchUserInfo() if implemented.
func (b *jwtAuthBackend) fetchUserInfo(ctx context.Context, pConfig CustomProvider, allClaims map[string]interface{}, role *jwtRole) error {
	// Fetch user info from custom provider if it's implemented
	if pConfig != nil {
		if uif, ok := pConfig.(UserInfoFetcher); ok {
			return uif.FetchUserInfo(ctx, b, allClaims, role)
		}
	}

	return nil
}

// Checks if there's a custom provider_config and calls FetchGroups() if implemented
func (b *jwtAuthBackend) fetchGroups(ctx context.Context, pConfig CustomProvider, allClaims map[string]interface{}, role *jwtRole, tokenSource oauth2.TokenSource) (interface{}, error) {
	// If the custom provider implements interface GroupsFetcher, call it,
	// otherwise fall through to the default method
	if pConfig != nil {
		if gf, ok := pConfig.(GroupsFetcher); ok {
			groupsRaw, err := gf.FetchGroups(ctx, b, allClaims, role, tokenSource)
			if err != nil {
				return nil, err
			}

			// Add groups obtained by provider-specific fetching to the claims
			// so that they can be used for bound_claims validation on the role.
			allClaims["groups"] = groupsRaw
		}
	}
	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw == nil {
		return nil, fmt.Errorf("%q claim not found in token", role.GroupsClaim)
	}

	return groupsClaimRaw, nil
}

func toAlg(a []string) []jwt.Alg {
	alg := make([]jwt.Alg, len(a))
	for i, e := range a {
		alg[i] = jwt.Alg(e)
	}
	return alg
}

const (
	pathLoginHelpSyn = `
	Authenticates to Vault using a JWT (or OIDC) token.
	`
	pathLoginHelpDesc = `
Authenticates JWTs.
`
)
