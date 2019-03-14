package jwtauth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/cidrutil"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"gopkg.in/square/go-jose.v2/jwt"
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

	token := d.Get("jwt").(string)
	if len(token) == 0 {
		return logical.ErrorResponse("missing token"), nil
	}

	if req.Connection != nil && !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.BoundCIDRs) {
		return logical.ErrorResponse("request originated from invalid CIDR"), nil
	}

	// Here is where things diverge. If it is using OIDC Discovery, validate
	// that way; otherwise validate against the locally configured keys. Once
	// things are validated, we re-unify the request path when evaluating the
	// claims.
	allClaims := map[string]interface{}{}
	switch {
	case len(config.ParsedJWTPubKeys) != 0:
		parsedJWT, err := jwt.ParseSigned(token)
		if err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("error parsing token: {{err}}", err).Error()), nil
		}

		claims := jwt.Claims{}

		var valid bool
		for _, key := range config.ParsedJWTPubKeys {
			if err := parsedJWT.Claims(key, &claims, &allClaims); err == nil {
				valid = true
				break
			}
		}
		if !valid {
			return logical.ErrorResponse("no known key successfully validated the token signature"), nil
		}

		// We require notbefore or expiry; if only one is provided, we allow 5 minutes of leeway.
		if claims.IssuedAt == 0 && claims.Expiry == 0 && claims.NotBefore == 0 {
			return logical.ErrorResponse("no issue time, notbefore, or expiration time encoded in token"), nil
		}
		if claims.Expiry == 0 {
			latestStart := claims.IssuedAt
			if claims.NotBefore > claims.IssuedAt {
				latestStart = claims.NotBefore
			}
			claims.Expiry = latestStart + 300
		}
		if claims.NotBefore == 0 {
			if claims.IssuedAt != 0 {
				claims.NotBefore = claims.IssuedAt
			} else {
				claims.NotBefore = claims.Expiry - 300
			}
		}

		if len(claims.Audience) > 0 && len(role.BoundAudiences) == 0 {
			return logical.ErrorResponse("audience claim found in JWT but no audiences bound to the role"), nil
		}

		expected := jwt.Expected{
			Issuer:  config.BoundIssuer,
			Subject: role.BoundSubject,
			Time:    time.Now(),
		}

		if err := claims.Validate(expected); err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("error validating claims: {{err}}", err).Error()), nil
		}

		if err := validateAudience(role.BoundAudiences, claims.Audience, true); err != nil {
			return logical.ErrorResponse(errwrap.Wrapf("error validating claims: {{err}}", err).Error()), nil
		}

	case config.OIDCDiscoveryURL != "":
		allClaims, err = b.verifyOIDCToken(ctx, config, role, token)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

	default:
		return nil, errors.New("unhandled case during login")
	}

	if err := validateBoundClaims(b.Logger(), role.BoundClaims, allClaims); err != nil {
		return logical.ErrorResponse("error validating claims: %s", err.Error()), nil
	}

	alias, groupAliases, err := b.createIdentity(allClaims, role)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	tokenMetadata := map[string]string{"role": roleName}
	for k, v := range alias.Metadata {
		tokenMetadata[k] = v
	}

	resp := &logical.Response{
		Auth: &logical.Auth{
			Policies:     role.Policies,
			DisplayName:  alias.Name,
			Period:       role.Period,
			NumUses:      role.NumUses,
			Alias:        alias,
			GroupAliases: groupAliases,
			InternalData: map[string]interface{}{
				"role": roleName,
			},
			Metadata: tokenMetadata,
			LeaseOptions: logical.LeaseOptions{
				Renewable: true,
				TTL:       role.TTL,
				MaxTTL:    role.MaxTTL,
			},
			BoundCIDRs: role.BoundCIDRs,
		},
	}

	return resp, nil
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
	resp.Auth.TTL = role.TTL
	resp.Auth.MaxTTL = role.MaxTTL
	resp.Auth.Period = role.Period
	return resp, nil
}

func (b *jwtAuthBackend) verifyOIDCToken(ctx context.Context, config *jwtConfig, role *jwtRole, rawToken string) (map[string]interface{}, error) {
	allClaims := make(map[string]interface{})

	provider, err := b.getProvider(ctx, config)
	if err != nil {
		return nil, errwrap.Wrapf("error getting provider for login operation: {{err}}", err)
	}

	oidcConfig := &oidc.Config{
		SupportedSigningAlgs: config.JWTSupportedAlgs,
	}

	if role.RoleType == "oidc" {
		oidcConfig.ClientID = config.OIDCClientID
	} else {
		oidcConfig.SkipClientIDCheck = true
	}
	verifier := provider.Verifier(oidcConfig)

	idToken, err := verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, errwrap.Wrapf("error validating signature: {{err}}", err)
	}

	if err := idToken.Claims(&allClaims); err != nil {
		return nil, errwrap.Wrapf("unable to successfully parse all claims from token: {{err}}", err)
	}

	if role.BoundSubject != "" && role.BoundSubject != idToken.Subject {
		return nil, errors.New("sub claim does not match bound subject")
	}

	if err := validateAudience(role.BoundAudiences, idToken.Audience, false); err != nil {
		return nil, errwrap.Wrapf("error validating claims: {{err}}", err)
	}

	return allClaims, nil
}

// createIdentity creates an alias and set of groups aliases based on the role
// definition and received claims.
func (b *jwtAuthBackend) createIdentity(allClaims map[string]interface{}, role *jwtRole) (*logical.Alias, []*logical.Alias, error) {
	userClaimRaw, ok := allClaims[role.UserClaim]
	if !ok {
		return nil, nil, fmt.Errorf("claim %q not found in token", role.UserClaim)
	}
	userName, ok := userClaimRaw.(string)
	if !ok {
		return nil, nil, fmt.Errorf("claim %q could not be converted to string", role.UserClaim)
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

	groupsClaimRaw := getClaim(b.Logger(), allClaims, role.GroupsClaim)

	if groupsClaimRaw == nil {
		return nil, nil, fmt.Errorf("%q claim not found in token", role.GroupsClaim)
	}
	groups, ok := groupsClaimRaw.([]interface{})

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

const (
	pathLoginHelpSyn = `
	Authenticates to Vault using a JWT (or OIDC) token.
	`
	pathLoginHelpDesc = `
Authenticates JWTs.
`
)
