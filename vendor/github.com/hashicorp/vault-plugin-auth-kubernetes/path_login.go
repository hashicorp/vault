// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kubeauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-jose/go-jose/v4"
	josejwt "github.com/go-jose/go-jose/v4/jwt"
	capjwt "github.com/hashicorp/cap/jwt"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/cidrutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/mitchellh/mapstructure"
)

const (
	metadataKeySAUID        = "service_account_uid"
	metadataKeySAName       = "service_account_name"
	metadataKeySANamespace  = "service_account_namespace"
	metadataKeySASecretName = "service_account_secret_name"
)

var reservedAliasMetadataKeys = map[string]struct{}{
	metadataKeySAUID:        {},
	metadataKeySAName:       {},
	metadataKeySANamespace:  {},
	metadataKeySASecretName: {},
}

// defaultJWTIssuer is used to verify the iss header on the JWT if the config doesn't specify an issuer.
var defaultJWTIssuer = "kubernetes/serviceaccount"

var (
	// signing algorithms supported by k8s OIDC
	// ref: https://github.com/kubernetes/kubernetes/blob/b4935d910dcf256288694391ef675acfbdb8e7a3/staging/src/k8s.io/apiserver/plugin/pkg/authenticator/token/oidc/oidc.go#L222-L233
	allowedSigningAlgs = []jose.SignatureAlgorithm{
		jose.RS256,
		jose.RS384,
		jose.RS512,
		jose.ES256,
		jose.ES384,
		jose.ES512,
		jose.PS256,
		jose.PS384,
		jose.PS512,
	}
	// allowedSigningAlgsCap is initialized with the values from allowedSigningAlgs
	allowedSigningAlgsCap = make([]capjwt.Alg, len(allowedSigningAlgs))
)

func init() {
	for idx := 0; idx < len(allowedSigningAlgs); idx++ {
		allowedSigningAlgsCap[idx] = capjwt.Alg(allowedSigningAlgs[idx])
	}
}

// pathLogin returns the path configurations for login endpoints
func pathLogin(b *kubeAuthBackend) *framework.Path {
	return &framework.Path{
		Pattern: "login$",

		DisplayAttrs: &framework.DisplayAttributes{
			OperationPrefix: operationPrefixKubernetes,
			OperationVerb:   "login",
		},

		Fields: map[string]*framework.FieldSchema{
			"role": {
				Type:        framework.TypeString,
				Description: `Name of the role against which the login is being attempted. This field is required`,
			},
			"jwt": {
				Type:        framework.TypeString,
				Description: `A signed JWT for authenticating a service account. This field is required.`,
			},
		},

		Callbacks: map[logical.Operation]framework.OperationFunc{
			logical.UpdateOperation:         b.pathLogin,
			logical.AliasLookaheadOperation: b.aliasLookahead,
			logical.ResolveRoleOperation:    b.pathResolveRole,
		},

		HelpSynopsis:    pathLoginHelpSyn,
		HelpDescription: pathLoginHelpDesc,
	}
}

// pathLogin is used to resolve the role to be used from a login request
func (b *kubeAuthBackend) pathResolveRole(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, resp := b.getFieldValueStr(data, "role")
	if resp != nil {
		return resp, nil
	}

	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role name %q", roleName)), nil
	}

	return logical.ResolveRoleResponse(roleName)
}

// pathLogin is used to authenticate to this backend
func (b *kubeAuthBackend) pathLogin(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, resp := b.getFieldValueStr(data, "role")
	if resp != nil {
		return resp, nil
	}

	jwtStr, resp := b.getFieldValueStr(data, "jwt")
	if resp != nil {
		return resp, nil
	}

	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role name %q", roleName)), nil
	}

	// Check for a CIDR match.
	if len(role.TokenBoundCIDRs) > 0 {
		if req.Connection == nil {
			b.Logger().Warn("token bound CIDRs found but no connection information available for validation")
			return nil, logical.ErrPermissionDenied
		}
		if !cidrutil.RemoteAddrIsOk(req.Connection.RemoteAddr, role.TokenBoundCIDRs) {
			return nil, logical.ErrPermissionDenied
		}
	}

	config, err := b.loadConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("could not load backend configuration")
	}

	client, err := b.getHTTPClient()
	if err != nil {
		b.Logger().Error("Failed to get the HTTP client", "err", err)
		return nil, logical.ErrUnrecoverable
	}

	sa, err := b.parseAndValidateJWT(ctx, client, jwtStr, role, config)
	if err != nil {
		if err == jose.ErrCryptoFailure || strings.Contains(err.Error(), "verifying token signature") {
			b.Logger().Debug(`login unauthorized`, "err", err)
			return nil, logical.ErrPermissionDenied
		}
		return nil, err
	}

	aliasName, err := b.getAliasName(role, sa)
	if err != nil {
		return nil, err
	}

	// look up the JWT token in the kubernetes API
	err = sa.lookup(ctx, client, jwtStr, role.Audience, b.reviewFactory(config))
	if err != nil {
		b.Logger().Debug(`login unauthorized`, "err", err)
		return nil, logical.ErrPermissionDenied
	}

	annotations := map[string]string{}
	if config.UseAnnotationsAsAliasMetadata {
		annotations, err = b.serviceAccountGetterFactory(config).annotations(ctx, client, jwtStr, sa.namespace(), sa.name())
		if err != nil {
			if errors.Is(err, errAliasMetadataReservedKeysFound) {
				return logical.ErrorResponse(err.Error()), nil
			}

			return nil, err
		}
	}

	uid, err := sa.uid()
	if err != nil {
		return nil, err
	}

	metadata := annotations
	metadata[metadataKeySAUID] = uid
	metadata[metadataKeySAName] = sa.name()
	metadata[metadataKeySANamespace] = sa.namespace()
	metadata[metadataKeySASecretName] = sa.SecretName

	auth := &logical.Auth{
		Alias: &logical.Alias{
			Name:     aliasName,
			Metadata: metadata,
		},
		InternalData: map[string]interface{}{
			"role": roleName,
		},
		Metadata: map[string]string{
			metadataKeySAUID:        uid,
			metadataKeySAName:       sa.name(),
			metadataKeySANamespace:  sa.namespace(),
			metadataKeySASecretName: sa.SecretName,
			"role":                  roleName,
		},
		DisplayName: fmt.Sprintf("%s-%s", sa.namespace(), sa.name()),
	}

	role.PopulateTokenAuth(auth)

	return &logical.Response{
		Auth: auth,
	}, nil
}

func (b *kubeAuthBackend) getFieldValueStr(data *framework.FieldData, param string) (string, *logical.Response) {
	val := data.Get(param).(string)
	if len(val) == 0 {
		return "", logical.ErrorResponse("missing %s", param)
	}
	return val, nil
}

func (b *kubeAuthBackend) getAliasName(role *roleStorageEntry, serviceAccount *serviceAccount) (string, error) {
	switch role.AliasNameSource {
	case aliasNameSourceSAUid, aliasNameSourceUnset:
		uid, err := serviceAccount.uid()
		if err != nil {
			return "", err
		}
		return uid, nil
	case aliasNameSourceSAName:
		ns, name := serviceAccount.namespace(), serviceAccount.name()
		if ns == "" || name == "" {
			return "", fmt.Errorf("service account namespace and name must be set")
		}
		return fmt.Sprintf("%s/%s", ns, name), nil
	default:
		return "", fmt.Errorf("unknown alias_name_source %q", role.AliasNameSource)
	}
}

// aliasLookahead returns the alias object with the SA UID from the JWT
// Claims.
// Only JWTs matching the specified role's configuration will be accepted as valid.
func (b *kubeAuthBackend) aliasLookahead(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
	roleName, resp := b.getFieldValueStr(data, "role")
	if resp != nil {
		return resp, nil
	}

	jwtStr, resp := b.getFieldValueStr(data, "jwt")
	if resp != nil {
		return resp, nil
	}

	b.l.RLock()
	defer b.l.RUnlock()

	role, err := b.role(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse(fmt.Sprintf("invalid role name %q", roleName)), nil
	}

	config, err := b.loadConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return nil, errors.New("could not load backend configuration")
	}
	// validation of the JWT against the provided role ensures alias look ahead requests
	// are authentic.
	client, err := b.getHTTPClient()
	if err != nil {
		b.Logger().Error("Failed to get the HTTP client", "err", err)
		return nil, logical.ErrUnrecoverable
	}

	sa, err := b.parseAndValidateJWT(ctx, client, jwtStr, role, config)
	if err != nil {
		return nil, err
	}

	aliasName, err := b.getAliasName(role, sa)
	if err != nil {
		return nil, err
	}

	return &logical.Response{
		Auth: &logical.Auth{
			Alias: &logical.Alias{
				Name: aliasName,
			},
		},
	}, nil
}

type DontVerifySignature struct{}

func (keySet DontVerifySignature) VerifySignature(_ context.Context, token string) (map[string]interface{}, error) {
	parsed, err := josejwt.ParseSigned(token, allowedSigningAlgs)
	if err != nil {
		return nil, err
	}
	claims := map[string]interface{}{}
	err = parsed.UnsafeClaimsWithoutVerification(&claims)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

// parseAndValidateJWT is used to parse, validate and lookup the JWT token.
func (b *kubeAuthBackend) parseAndValidateJWT(ctx context.Context, client *http.Client, jwtStr string,
	role *roleStorageEntry, config *kubeConfig,
) (*serviceAccount, error) {
	expected := capjwt.Expected{
		SigningAlgorithms: allowedSigningAlgsCap,
	}

	// perform ISS Claim validation if configured
	if !config.DisableISSValidation {
		// set the expected issuer to the default kubernetes issuer if the config doesn't specify it
		if config.Issuer != "" {
			expected.Issuer = config.Issuer
		} else {
			config.Issuer = defaultJWTIssuer
		}
	}

	// validate the audience if the role expects it
	if role.Audience != "" {
		expected.Audiences = []string{role.Audience}
	}

	// Parse into JWT
	var err error
	var keySet capjwt.KeySet
	if len(config.PublicKeys) == 0 {
		// we don't verify the signature if we aren't configured with public keys
		keySet = DontVerifySignature{}
	} else {
		keySet, err = capjwt.NewStaticKeySet(config.PublicKeys)
		if err != nil {
			return nil, err
		}
	}
	validator, err := capjwt.NewValidator(keySet)
	if err != nil {
		return nil, err
	}

	claims, err := validator.ValidateAllowMissingIatNbfExp(nil, jwtStr, expected)
	if err != nil {
		return nil, logical.CodedError(http.StatusForbidden, err.Error())
	}

	sa := &serviceAccount{}

	// Decode claims into a service account object
	err = mapstructure.Decode(claims, sa)
	if err != nil {
		return nil, err
	}

	// verify the service account name is allowed
	if len(role.ServiceAccountNames) > 1 || role.ServiceAccountNames[0] != "*" {
		if !strutil.StrListContainsGlob(role.ServiceAccountNames, sa.name()) {
			return nil, logical.CodedError(http.StatusForbidden,
				fmt.Sprintf("service account name not authorized"))
		}
	}

	// verify the namespace is allowed
	var allowed bool
	if len(role.ServiceAccountNamespaces) > 0 {
		if role.ServiceAccountNamespaces[0] == "*" || strutil.StrListContainsGlob(role.ServiceAccountNamespaces, sa.namespace()) {
			allowed = true
		}
	}

	// verify the namespace selector matches the namespace
	if !allowed && role.ServiceAccountNamespaceSelector != "" {
		allowed, err = b.namespaceValidatorFactory(config).validateLabels(ctx,
			client, sa.namespace(), role.ServiceAccountNamespaceSelector)
	}

	if !allowed {
		errMsg := "namespace not authorized"
		if err != nil {
			errMsg = fmt.Sprintf("%s err=%s", errMsg, err)
		}
		codedErr := logical.CodedError(http.StatusForbidden, errMsg)
		return nil, codedErr
	}

	// If we don't have any public keys to verify, return the sa and end early.
	if len(config.PublicKeys) == 0 {
		return sa, nil
	}

	return sa, nil
}

// serviceAccount holds the metadata from the JWT token and is used to lookup
// the JWT in the kubernetes API and compare the results.
type serviceAccount struct {
	Name       string   `mapstructure:"kubernetes.io/serviceaccount/service-account.name"`
	UID        string   `mapstructure:"kubernetes.io/serviceaccount/service-account.uid"`
	SecretName string   `mapstructure:"kubernetes.io/serviceaccount/secret.name"`
	Namespace  string   `mapstructure:"kubernetes.io/serviceaccount/namespace"`
	Audience   []string `mapstructure:"aud"`

	// the JSON returned from reviewing a Projected Service account has a
	// different structure, where the information is in a sub-structure instead of
	// at the top level
	Kubernetes *projectedServiceToken `mapstructure:"kubernetes.io"`
	Expiration int64                  `mapstructure:"exp"`
	IssuedAt   int64                  `mapstructure:"iat"`
}

// uid returns the UID for the service account, preferring the projected service
// account value if found
// return an error when the UID is empty.
func (s *serviceAccount) uid() (string, error) {
	uid := s.UID
	if s.Kubernetes != nil && s.Kubernetes.ServiceAccount != nil {
		uid = s.Kubernetes.ServiceAccount.UID
	}

	if uid == "" {
		return "", errors.New("could not parse UID from claims")
	}
	return uid, nil
}

// name returns the name for the service account, preferring the projected
// service account value if found. This is "default" for projected service
// accounts
func (s *serviceAccount) name() string {
	if s.Kubernetes != nil && s.Kubernetes.ServiceAccount != nil {
		return s.Kubernetes.ServiceAccount.Name
	}
	return s.Name
}

// namespace returns the namespace for the service account, preferring the
// projected service account value if found
func (s *serviceAccount) namespace() string {
	if s.Kubernetes != nil {
		return s.Kubernetes.Namespace
	}
	return s.Namespace
}

type projectedServiceToken struct {
	Namespace      string        `mapstructure:"namespace"`
	Pod            *k8sObjectRef `mapstructure:"pod"`
	ServiceAccount *k8sObjectRef `mapstructure:"serviceaccount"`
}

type k8sObjectRef struct {
	Name string `mapstructure:"name"`
	UID  string `mapstructure:"uid"`
}

// lookup calls the TokenReview API in kubernetes to verify the token and secret
// still exist.
func (s *serviceAccount) lookup(ctx context.Context, client *http.Client, jwtStr string, aud string, tr tokenReviewer) error {
	// This is somewhat redundant, as we are asking k8s if the token's audiences
	// overlap with wantAud, but setting wantAud to the token's own audiences by
	// default. In the case that `audience` was set on the Vault role, we have
	// already validated that it matches one of the token's audiences by this
	// point, so we are essentially asking k8s to either ignore audience
	// validation or check our homework.
	wantAud := s.Audience
	if aud != "" {
		wantAud = []string{aud}
	}
	r, err := tr.Review(ctx, client, jwtStr, wantAud)
	if err != nil {
		return err
	}

	// Verify the returned metadata matches the expected data from the service
	// account.
	if s.name() != r.Name {
		return errors.New("JWT names did not match")
	}
	uid, err := s.uid()
	if err != nil {
		return err
	}
	if uid != r.UID {
		return errors.New("JWT UIDs did not match")
	}
	if s.namespace() != r.Namespace {
		return errors.New("JWT namepaces did not match")
	}

	return nil
}

// Invoked when the token issued by this backend is attempting a renewal.
func (b *kubeAuthBackend) pathLoginRenew() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, data *framework.FieldData) (*logical.Response, error) {
		roleName := req.Auth.InternalData["role"].(string)
		if roleName == "" {
			return nil, fmt.Errorf("failed to fetch role_name during renewal")
		}

		b.l.RLock()
		defer b.l.RUnlock()

		// Ensure that the Role still exists.
		role, err := b.role(ctx, req.Storage, roleName)
		if err != nil {
			return nil, fmt.Errorf("failed to validate role %s during renewal:%s", roleName, err)
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
}

const (
	pathLoginHelpSyn  = `Authenticates Kubernetes service accounts with Vault.`
	pathLoginHelpDesc = `
Authenticate Kubernetes service accounts.
`
)
