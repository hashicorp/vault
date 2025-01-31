// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	mathrand "math/rand"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-secure-stdlib/base62"
	"github.com/hashicorp/go-secure-stdlib/strutil"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/consts"
	"github.com/hashicorp/vault/sdk/helper/cryptoutil"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/exp/maps"
)

type oidcConfig struct {
	// Issuer is the scheme://host:port component of the issuer set as
	// configuration in the Vault API. It is the URL base.
	Issuer string `json:"issuer"`

	// effectiveIssuer is a calculated field and will be either Issuer (if
	// that's set) or the Vault instance's api_addr, followed by the path
	// /v1/<namespace_path>/identity/oidc.
	effectiveIssuer string
}

// fullIssuer returns the full issuer for the config, suitable for OpenID metadata and
// token claims. It takes an optional child, which must be of the value "" or "plugins".
// The child will be appended as the last path segment on the returned issuer URL.
func (c *oidcConfig) fullIssuer(child string) (string, error) {
	if !validChildIssuer(child) {
		return "", fmt.Errorf("invalid child issuer %q", child)
	}

	issuer, err := url.JoinPath(c.effectiveIssuer, child)
	if err != nil {
		return "", fmt.Errorf("failed to join issuer: %w", err)
	}

	return issuer, nil
}

type expireableKey struct {
	KeyID    string    `json:"key_id"`
	ExpireAt time.Time `json:"expire_at"`
}

type namedKey struct {
	name             string
	Algorithm        string           `json:"signing_algorithm"`
	VerificationTTL  time.Duration    `json:"verification_ttl"`
	RotationPeriod   time.Duration    `json:"rotation_period"`
	KeyRing          []*expireableKey `json:"key_ring"`
	SigningKey       *jose.JSONWebKey `json:"signing_key"`
	NextSigningKey   *jose.JSONWebKey `json:"next_signing_key"`
	NextRotation     time.Time        `json:"next_rotation"`
	AllowedClientIDs []string         `json:"allowed_client_ids"`
}

type role struct {
	TokenTTL time.Duration `json:"token_ttl"`
	Key      string        `json:"key"`
	Template string        `json:"template"`
	ClientID string        `json:"client_id"`
}

// idToken contains the required OIDC fields.
//
// Templated claims will be merged into the final output. Those claims may
// include top-level keys, but those keys may not overwrite any of the
// required OIDC fields.
type idToken struct {
	Issuer          string `json:"iss"`       // api_addr or custom Issuer
	Namespace       string `json:"namespace"` // Namespace of issuer
	Subject         string `json:"sub"`       // Entity ID
	Audience        string `json:"aud"`       // Role or client ID will be used here.
	Expiry          int64  `json:"exp"`       // Expiration, as determined by the role or client.
	IssuedAt        int64  `json:"iat"`       // Time of token creation
	Nonce           string `json:"nonce"`     // Nonce given in OIDC authentication requests
	AuthTime        int64  `json:"auth_time"` // AuthTime given in OIDC authentication requests
	AccessTokenHash string `json:"at_hash"`   // Access token hash value
	CodeHash        string `json:"c_hash"`    // Authorization code hash value
}

// discovery contains a subset of the required elements of OIDC discovery needed
// for JWT verification libraries to use the .well-known endpoint.
//
// https://openid.net/specs/openid-connect-discovery-1_0.html#ProviderMetadata
type discovery struct {
	Issuer        string   `json:"issuer"`
	Keys          string   `json:"jwks_uri"`
	ResponseTypes []string `json:"response_types_supported"`
	Subjects      []string `json:"subject_types_supported"`
	IDTokenAlgs   []string `json:"id_token_signing_alg_values_supported"`
}

// oidcCache is a thin wrapper around go-cache to partition by namespace
type oidcCache struct {
	c *cache.Cache
}

var (
	errNilNamespace = errors.New("nil namespace in oidc cache request")

	reservedClaims = []string{
		"iat", "aud", "exp", "iss",
		"sub", "namespace", "nonce",
		"auth_time", "at_hash", "c_hash",
	}
	supportedAlgs = []string{
		string(jose.RS256),
		string(jose.RS384),
		string(jose.RS512),
		string(jose.ES256),
		string(jose.ES384),
		string(jose.ES512),
		string(jose.EdDSA),
	}
)

const (
	issuerPath              = "identity/oidc"
	oidcTokensPrefix        = "oidc_tokens/"
	namedKeyCachePrefix     = "namedKeys/"
	oidcConfigStorageKey    = oidcTokensPrefix + "config/"
	namedKeyConfigPath      = oidcTokensPrefix + "named_keys/"
	publicKeysConfigPath    = oidcTokensPrefix + "public_keys/"
	roleConfigPath          = oidcTokensPrefix + "roles/"
	baseIdentityTokenIssuer = ""
	deleteKeyErrorFmt       = "unable to delete key %q because it is currently referenced by these %s: %s"
)

// optionalChildIssuerRegex is a regex for optionally accepting a field in an
// API request as a single path segment. Adapted from framework.OptionalParamRegex
// to not include additional forward slashes.
func optionalChildIssuerRegex(name string) string {
	return fmt.Sprintf(`(/(?P<%s>[^/]+))?`, name)
}

func oidcPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "oidc/config/?$",

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
			},

			Fields: map[string]*framework.FieldSchema{
				"issuer": {
					Type:        framework.TypeString,
					Description: "Issuer URL to be used in the iss claim of the token. If not set, Vault's app_addr will be used.",
				},
			},

			Operations: map[logical.Operation]framework.OperationHandler{
				logical.ReadOperation: &framework.PathOperation{
					Callback: i.pathOIDCReadConfig,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationSuffix: "configuration",
					},
				},
				logical.UpdateOperation: &framework.PathOperation{
					Callback: i.pathOIDCUpdateConfig,
					DisplayAttrs: &framework.DisplayAttributes{
						OperationVerb: "configure",
					},
				},
			},

			HelpSynopsis:    "OIDC configuration",
			HelpDescription: "Update OIDC configuration in the identity backend",
		},
		{
			Pattern: "oidc/key/" + framework.GenericNameRegex("name"),

			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "key",
			},

			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the key",
				},

				"rotation_period": {
					Type:        framework.TypeDurationSecond,
					Description: "How often to generate a new keypair.",
					Default:     "24h",
				},

				"verification_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Controls how long the public portion of a key will be available for verification after being rotated.",
					Default:     "24h",
				},

				"algorithm": {
					Type:        framework.TypeString,
					Description: "Signing algorithm to use. This will default to RS256.",
					Default:     "RS256",
				},

				"allowed_client_ids": {
					Type:        framework.TypeCommaStringSlice,
					Description: "Comma separated string or array of role client ids allowed to use this key for signing. If empty no roles are allowed. If \"*\" all roles are allowed.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: i.pathOIDCCreateUpdateKey,
				logical.UpdateOperation: i.pathOIDCCreateUpdateKey,
				logical.ReadOperation:   i.pathOIDCReadKey,
				logical.DeleteOperation: i.pathOIDCDeleteKey,
			},
			ExistenceCheck:  i.pathOIDCKeyExistenceCheck,
			HelpSynopsis:    "CRUD operations for OIDC keys.",
			HelpDescription: "Create, Read, Update, and Delete OIDC named keys.",
		},
		{
			Pattern: "oidc/key/" + framework.GenericNameRegex("name") + "/rotate/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationVerb:   "rotate",
				OperationSuffix: "key",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the key",
				},
				"verification_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Controls how long the public portion of a key will be available for verification after being rotated. Setting verification_ttl here will override the verification_ttl set on the key.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCRotateKey,
			},
			HelpSynopsis:    "Rotate a named OIDC key.",
			HelpDescription: "Manually rotate a named OIDC key. Rotating a named key will cause a new underlying signing key to be generated. The public portion of the underlying rotated signing key will continue to live for the verification_ttl duration.",
		},
		{
			Pattern: "oidc/key/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "keys",
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListKey,
			},
			HelpSynopsis:    "List OIDC keys",
			HelpDescription: "List all named OIDC keys",
		},
		{
			Pattern: "oidc" + optionalChildIssuerRegex("child") + "/\\.well-known/openid-configuration/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "open-id-configuration",
			},
			Fields: map[string]*framework.FieldSchema{
				"child": {
					Type:        framework.TypeString,
					Description: "Name of the child issuer",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCDiscovery,
			},
			HelpSynopsis:    "Query OIDC configurations",
			HelpDescription: "Query this path to retrieve the configured OIDC Issuer and Keys endpoints, response types, subject types, and signing algorithms used by the OIDC backend.",
		},
		{
			Pattern: "oidc" + optionalChildIssuerRegex("child") + "/\\.well-known/keys/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "public-keys",
			},
			Fields: map[string]*framework.FieldSchema{
				"child": {
					Type:        framework.TypeString,
					Description: "Name of the child issuer",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCReadPublicKeys,
			},
			HelpSynopsis:    "Retrieve public keys",
			HelpDescription: "Query this path to retrieve the public portion of keys used to sign OIDC tokens. Clients can use this to validate the authenticity of the OIDC token claims.",
		},
		{
			Pattern: "oidc/token/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationVerb:   "generate",
				OperationSuffix: "token",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCGenerateToken,
			},
			HelpSynopsis:    "Generate an OIDC token",
			HelpDescription: "Generate an OIDC token against a configured role. The vault token used to call this path must have a corresponding entity.",
		},
		{
			Pattern: "oidc/role/" + framework.GenericNameRegex("name"),
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "role",
			},
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role",
				},
				"key": {
					Type:        framework.TypeString,
					Description: "The OIDC key to use for generating tokens. The specified key must already exist.",
					Required:    true,
				},
				"template": {
					Type:        framework.TypeString,
					Description: "The template string to use for generating tokens. This may be in string-ified JSON or base64 format.",
				},
				"ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "TTL of the tokens generated against the role.",
					Default:     "24h",
				},
				"client_id": {
					Type:        framework.TypeString,
					Description: "Optional client_id",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCCreateUpdateRole,
				logical.CreateOperation: i.pathOIDCCreateUpdateRole,
				logical.ReadOperation:   i.pathOIDCReadRole,
				logical.DeleteOperation: i.pathOIDCDeleteRole,
			},
			ExistenceCheck:  i.pathOIDCRoleExistenceCheck,
			HelpSynopsis:    "CRUD operations on OIDC Roles",
			HelpDescription: "Create, Read, Update, and Delete OIDC Roles. OIDC tokens are generated against roles which can be configured to determine how OIDC tokens are generated.",
		},
		{
			Pattern: "oidc/role/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationSuffix: "roles",
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListRole,
			},
			HelpSynopsis:    "List configured OIDC roles",
			HelpDescription: "List all configured OIDC roles in the identity backend.",
		},
		{
			Pattern: "oidc/introspect/?$",
			DisplayAttrs: &framework.DisplayAttributes{
				OperationPrefix: "oidc",
				OperationVerb:   "introspect",
			},
			Fields: map[string]*framework.FieldSchema{
				"token": {
					Type:        framework.TypeString,
					Description: "Token to verify",
				},
				"client_id": {
					Type:        framework.TypeString,
					Description: "Optional client_id to verify",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCIntrospect,
			},
			HelpSynopsis:    "Verify the authenticity of an OIDC token",
			HelpDescription: "Use this path to verify the authenticity of an OIDC token and whether the associated entity is active and enabled.",
		},
	}
}

func (i *IdentityStore) pathOIDCReadConfig(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := i.getOIDCConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			"issuer": c.Issuer,
		},
	}

	if i.redirectAddr == "" && c.Issuer == "" {
		resp.AddWarning(`Both "issuer" and Vault's "api_addr" are empty. ` +
			`The issuer claim in generated tokens will not be network reachable.`)
	}

	return resp, nil
}

func (i *IdentityStore) pathOIDCUpdateConfig(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var resp *logical.Response

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	issuerRaw, ok := d.GetOk("issuer")
	if !ok {
		return nil, nil
	}

	issuer := issuerRaw.(string)

	if issuer != "" {
		// verify that issuer is the correct format:
		//   - http or https
		//   - host name
		//   - optional port
		//   - nothing more
		valid := false
		if u, err := url.Parse(issuer); err == nil {
			u2 := url.URL{
				Scheme: u.Scheme,
				Host:   u.Host,
			}
			valid = (*u == u2) &&
				(u.Scheme == "http" || u.Scheme == "https") &&
				u.Host != ""
		}

		if !valid {
			return logical.ErrorResponse(
				"invalid issuer, which must include only a scheme, host, " +
					"and optional port (e.g. https://example.com:8200)"), nil
		}

		resp = &logical.Response{
			Warnings: []string{`If "issuer" is set explicitly, all tokens must be ` +
				`validated against that address, including those issued by secondary ` +
				`clusters. Setting issuer to "" will restore the default behavior of ` +
				`using the cluster's api_addr as the issuer.`},
		}
	}

	c := oidcConfig{
		Issuer: issuer,
	}

	entry, err := logical.StorageEntryJSON(oidcConfigStorageKey, c)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *IdentityStore) getOIDCConfig(ctx context.Context, s logical.Storage) (*oidcConfig, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	v, ok, err := i.oidcCache.Get(ns, "config")
	if err != nil {
		return nil, err
	}

	if ok {
		return v.(*oidcConfig), nil
	}

	var c oidcConfig
	entry, err := s.Get(ctx, oidcConfigStorageKey)
	if err != nil {
		return nil, err
	}

	if entry != nil {
		if err := entry.DecodeJSON(&c); err != nil {
			return nil, err
		}
	}

	c.effectiveIssuer = c.Issuer
	if c.effectiveIssuer == "" {
		c.effectiveIssuer = i.redirectAddr
	}

	c.effectiveIssuer += "/v1/" + ns.Path + issuerPath

	if err := i.oidcCache.SetDefault(ns, "config", &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// handleOIDCCreateKey is used to create a new named key or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	var key namedKey
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&key); err != nil {
				return nil, err
			}
		}
	}

	if rotationPeriodRaw, ok := d.GetOk("rotation_period"); ok {
		key.RotationPeriod = time.Duration(rotationPeriodRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		key.RotationPeriod = time.Duration(d.Get("rotation_period").(int)) * time.Second
	}

	if key.RotationPeriod < 1*time.Minute {
		return logical.ErrorResponse("rotation_period must be at least one minute"), nil
	}

	if verificationTTLRaw, ok := d.GetOk("verification_ttl"); ok {
		key.VerificationTTL = time.Duration(verificationTTLRaw.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		key.VerificationTTL = time.Duration(d.Get("verification_ttl").(int)) * time.Second
	}

	if key.VerificationTTL > 10*key.RotationPeriod {
		return logical.ErrorResponse("verification_ttl cannot be longer than 10x rotation_period"), nil
	}

	if req.Operation == logical.UpdateOperation {
		// ensure any roles referencing this key do not already have a token_ttl
		// greater than the key's verification_ttl
		roles, err := i.rolesReferencingTargetKeyName(ctx, req, name)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			if role.TokenTTL > key.VerificationTTL {
				errorMessage := fmt.Sprintf(
					"unable to update key %q because it is currently referenced by one or more roles with a token ttl greater than %d seconds",
					name,
					key.VerificationTTL/time.Second,
				)
				return logical.ErrorResponse(errorMessage), nil
			}
		}

		// ensure any clients referencing this key do not already have a id_token_ttl
		// greater than the key's verification_ttl
		clients, err := i.clientsReferencingTargetKeyName(ctx, req, name)
		if err != nil {
			return nil, err
		}
		for _, client := range clients {
			if client.IDTokenTTL > key.VerificationTTL {
				errorMessage := fmt.Sprintf(
					"unable to update key %q because it is currently referenced by one or more clients with an id_token_ttl greater than %d seconds",
					name,
					key.VerificationTTL/time.Second,
				)
				return logical.ErrorResponse(errorMessage), nil
			}
		}
	}

	if allowedClientIDsRaw, ok := d.GetOk("allowed_client_ids"); ok {
		key.AllowedClientIDs = allowedClientIDsRaw.([]string)
	} else if req.Operation == logical.CreateOperation {
		key.AllowedClientIDs = d.Get("allowed_client_ids").([]string)
	}

	prevAlgorithm := key.Algorithm
	if algorithm, ok := d.GetOk("algorithm"); ok {
		key.Algorithm = algorithm.(string)
	} else if req.Operation == logical.CreateOperation {
		key.Algorithm = d.Get("algorithm").(string)
	}

	if !strutil.StrListContains(supportedAlgs, key.Algorithm) {
		return logical.ErrorResponse("unknown signing algorithm %q", key.Algorithm), nil
	}

	now := time.Now()

	// Update next rotation time if it is unset or now earlier than previously set.
	nextRotation := now.Add(key.RotationPeriod)
	if key.NextRotation.IsZero() || nextRotation.Before(key.NextRotation) {
		key.NextRotation = nextRotation
	}

	// generate current and next keys if creating a new key or changing algorithms
	if key.Algorithm != prevAlgorithm {
		err = key.generateAndSetKey(ctx, i.Logger(), req.Storage)
		if err != nil {
			return nil, err
		}

		err = key.generateAndSetNextKey(ctx, i.Logger(), req.Storage)
		if err != nil {
			return nil, err
		}
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	// store named key
	entry, err := logical.StorageEntryJSON(namedKeyConfigPath+name, key)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	return nil, nil
}

// handleOIDCReadKey is used to read an existing key
func (i *IdentityStore) pathOIDCReadKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	i.oidcLock.RLock()
	defer i.oidcLock.RUnlock()

	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}

	var storedNamedKey namedKey
	if err := entry.DecodeJSON(&storedNamedKey); err != nil {
		return nil, err
	}
	return &logical.Response{
		Data: map[string]interface{}{
			"rotation_period":    int64(storedNamedKey.RotationPeriod.Seconds()),
			"verification_ttl":   int64(storedNamedKey.VerificationTTL.Seconds()),
			"algorithm":          storedNamedKey.Algorithm,
			"allowed_client_ids": storedNamedKey.AllowedClientIDs,
		},
	}, nil
}

// keyIDsByName will return a slice of key IDs for the given key name
func (i *IdentityStore) keyIDsByName(ctx context.Context, s logical.Storage, name string) ([]string, error) {
	var keyIDs []string
	entry, err := s.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return keyIDs, err
	}
	if entry == nil {
		return keyIDs, nil
	}

	var key namedKey
	if err := entry.DecodeJSON(&key); err != nil {
		return keyIDs, err
	}

	for _, k := range key.KeyRing {
		keyIDs = append(keyIDs, k.KeyID)
	}

	return keyIDs, nil
}

// rolesReferencingTargetKeyName returns a map of role names to roles
// referencing targetKeyName.
//
// Note: this is not threadsafe. It is to be called with Lock already held.
func (i *IdentityStore) rolesReferencingTargetKeyName(ctx context.Context, req *logical.Request, targetKeyName string) (map[string]role, error) {
	roleNames, err := req.Storage.List(ctx, roleConfigPath)
	if err != nil {
		return nil, err
	}

	var tempRole role
	roles := make(map[string]role)
	for _, roleName := range roleNames {
		entry, err := req.Storage.Get(ctx, roleConfigPath+roleName)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&tempRole); err != nil {
				return nil, err
			}
			if tempRole.Key == targetKeyName {
				roles[roleName] = tempRole
			}
		}
	}

	return roles, nil
}

// roleNamesReferencingTargetKeyName returns a slice of strings of role
// names referencing targetKeyName.
//
// Note: this is not threadsafe. It is to be called with Lock already held.
func (i *IdentityStore) roleNamesReferencingTargetKeyName(ctx context.Context, req *logical.Request, targetKeyName string) ([]string, error) {
	roles, err := i.rolesReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		return nil, err
	}

	var names []string
	for key := range roles {
		names = append(names, key)
	}
	sort.Strings(names)
	return names, nil
}

// listMounts returns all mount entries in the namespace.
// Returns an error if the namespace is nil.
func (i *IdentityStore) listMounts(ns *namespace.Namespace) ([]*MountEntry, error) {
	if ns == nil {
		return nil, errors.New("namespace must not be nil")
	}

	secretMounts, err := i.mountLister.ListMounts()
	if err != nil {
		return nil, err
	}
	authMounts, err := i.mountLister.ListAuths()
	if err != nil {
		return nil, err
	}

	var allMounts []*MountEntry
	for _, mount := range append(authMounts, secretMounts...) {
		if mount.NamespaceID == ns.ID {
			allMounts = append(allMounts, mount)
		}
	}

	return allMounts, nil
}

// mountsReferencingKey returns a sorted list of all mount entry paths referencing
// the key in the namespace. Returns an error if the namespace is nil.
func (i *IdentityStore) mountsReferencingKey(ns *namespace.Namespace, key string) ([]string, error) {
	if ns == nil {
		return nil, errors.New("namespace must not be nil")
	}

	allMounts, err := i.listMounts(ns)
	if err != nil {
		return nil, err
	}

	pathsWithKey := make(map[string]struct{})
	for _, mount := range allMounts {
		if mount.Config.IdentityTokenKey == key {
			pathsWithKey[mount.Path] = struct{}{}
		}
	}

	paths := maps.Keys(pathsWithKey)
	sort.Strings(paths)
	return paths, nil
}

// handleOIDCDeleteKey is used to delete a key
func (i *IdentityStore) pathOIDCDeleteKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	targetKeyName := d.Get("name").(string)

	if targetKeyName == defaultKeyName {
		return logical.ErrorResponse("deletion of key %q not allowed",
			defaultKeyName), nil
	}

	i.oidcLock.Lock()

	roleNames, err := i.roleNamesReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		i.oidcLock.Unlock()
		return nil, err
	}

	if len(roleNames) > 0 {
		errorMessage := fmt.Sprintf(deleteKeyErrorFmt,
			targetKeyName, "roles", strings.Join(roleNames, ", "))
		i.oidcLock.Unlock()
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}

	clientNames, err := i.clientNamesReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		i.oidcLock.Unlock()
		return nil, err
	}

	if len(clientNames) > 0 {
		errorMessage := fmt.Sprintf(deleteKeyErrorFmt,
			targetKeyName, "clients", strings.Join(clientNames, ", "))
		i.oidcLock.Unlock()
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}

	mounts, err := i.mountsReferencingKey(ns, targetKeyName)
	if err != nil {
		i.oidcLock.Unlock()
		return nil, err
	}
	if len(mounts) > 0 {
		errorMessage := fmt.Sprintf(deleteKeyErrorFmt,
			targetKeyName, "mounts", strings.Join(mounts, ", "))
		i.oidcLock.Unlock()
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}

	// key can safely be deleted now
	err = req.Storage.Delete(ctx, namedKeyConfigPath+targetKeyName)
	if err != nil {
		i.oidcLock.Unlock()
		return nil, err
	}

	i.oidcLock.Unlock()

	_, err = i.expireOIDCPublicKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	return nil, nil
}

// handleOIDCListKey is used to list named keys
func (i *IdentityStore) pathOIDCListKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	i.oidcLock.RLock()
	defer i.oidcLock.RUnlock()

	keys, err := req.Storage.List(ctx, namedKeyConfigPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

// pathOIDCRotateKey is used to manually trigger a rotation on the named key
func (i *IdentityStore) pathOIDCRotateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	// load the named key and perform a rotation
	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("no named key found at %q", name), logical.ErrInvalidRequest
	}

	var storedNamedKey namedKey
	if err := entry.DecodeJSON(&storedNamedKey); err != nil {
		return nil, err
	}
	storedNamedKey.name = name

	// call rotate with an appropriate overrideTTL where < 0 means no override
	verificationTTLOverride := -1 * time.Second

	if ttlRaw, ok := d.GetOk("verification_ttl"); ok {
		verificationTTLOverride = time.Duration(ttlRaw.(int)) * time.Second
	}

	if err := storedNamedKey.rotate(ctx, i.Logger(), req.Storage, verificationTTLOverride); err != nil {
		return nil, err
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	return nil, nil
}

func (i *IdentityStore) pathOIDCKeyExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	name := d.Get("name").(string)

	i.oidcLock.RLock()
	defer i.oidcLock.RUnlock()

	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return false, err
	}

	return entry != nil, nil
}

// handleOIDCGenerateSignToken generates and signs an OIDC token
func (i *IdentityStore) pathOIDCGenerateToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	roleName := d.Get("name").(string)

	role, err := i.getOIDCRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q not found", roleName), nil
	}

	key, err := i.getNamedKey(ctx, req.Storage, role.Key)
	if err != nil {
		return nil, err
	}
	if key == nil {
		return logical.ErrorResponse("key %q not found", role.Key), nil
	}

	// Validate that the role is allowed to sign with its key (the key could have been updated)
	if !strutil.StrListContains(key.AllowedClientIDs, "*") && !strutil.StrListContains(key.AllowedClientIDs, role.ClientID) {
		return logical.ErrorResponse("the key %q does not list the client ID of the role %q as an allowed client ID", role.Key, roleName), nil
	}

	// generate an OIDC token from entity data
	if req.EntityID == "" {
		return logical.ErrorResponse("no entity associated with the request's token"), nil
	}

	config, err := i.getOIDCConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	retResp := &logical.Response{}
	expiry := role.TokenTTL
	if expiry > key.VerificationTTL {
		expiry = key.VerificationTTL
		retResp.AddWarning(fmt.Sprintf("a role's token ttl cannot be longer "+
			"than the verification_ttl of the key it references, setting token ttl to %d", expiry))
	}

	issuer, err := config.fullIssuer(baseIdentityTokenIssuer)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	now := time.Now()
	idToken := idToken{
		Issuer:    issuer,
		Namespace: ns.ID,
		Subject:   req.EntityID,
		Audience:  role.ClientID,
		Expiry:    now.Add(expiry).Unix(),
		IssuedAt:  now.Unix(),
	}

	e, err := i.MemDBEntityByID(req.EntityID, true)
	if err != nil {
		return nil, err
	}
	if e == nil {
		return nil, fmt.Errorf("error loading entity ID %q", req.EntityID)
	}

	groups, inheritedGroups, err := i.groupsByEntityID(e.ID)
	if err != nil {
		return nil, err
	}

	groups = append(groups, inheritedGroups...)

	// Parse and integrate the populated template. Structural errors with the template _should_
	// be caught during configuration. Error found during runtime will be logged, but they will
	// not block generation of the basic ID token. They should not be returned to the requester.
	_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
		Mode:        identitytpl.JSONTemplating,
		String:      role.Template,
		Entity:      identity.ToSDKEntity(e),
		Groups:      identity.ToSDKGroups(groups),
		NamespaceID: ns.ID,
	})
	if err != nil {
		i.Logger().Warn("error populating OIDC token template", "template", role.Template, "error", err)
	}

	payload, err := idToken.generatePayload(i.Logger(), populatedTemplate)
	if err != nil {
		i.Logger().Warn("error populating OIDC token template", "error", err)
	}

	signedIdToken, err := key.signPayload(payload)
	if err != nil {
		return nil, fmt.Errorf("error signing OIDC token: %w", err)
	}

	retResp.Data = map[string]interface{}{
		"token":     signedIdToken,
		"client_id": role.ClientID,
		"ttl":       int64(role.TokenTTL.Seconds()),
	}
	return retResp, nil
}

func (i *IdentityStore) getNamedKey(ctx context.Context, s logical.Storage, name string) (*namedKey, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Attempt to get the key from the cache
	keyRaw, found, err := i.oidcCache.Get(ns, namedKeyCachePrefix+name)
	if err != nil {
		return nil, err
	}
	if key, ok := keyRaw.(*namedKey); ok && found {
		return key, nil
	}

	// Fall back to reading the key from storage
	entry, err := s.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, nil
	}
	var key namedKey
	if err := entry.DecodeJSON(&key); err != nil {
		return nil, err
	}

	// Cache the key
	if err := i.oidcCache.SetDefault(ns, namedKeyCachePrefix+name, &key); err != nil {
		i.logger.Warn("failed to cache key", "error", err)
	}

	return &key, nil
}

func (tok *idToken) generatePayload(logger hclog.Logger, templates ...string) ([]byte, error) {
	output := map[string]interface{}{
		"iss":       tok.Issuer,
		"namespace": tok.Namespace,
		"sub":       tok.Subject,
		"aud":       tok.Audience,
		"exp":       tok.Expiry,
		"iat":       tok.IssuedAt,
	}

	// Copy optional claims into output
	if len(tok.Nonce) > 0 {
		output["nonce"] = tok.Nonce
	}
	if tok.AuthTime > 0 {
		output["auth_time"] = tok.AuthTime
	}
	if len(tok.AccessTokenHash) > 0 {
		output["at_hash"] = tok.AccessTokenHash
	}
	if len(tok.CodeHash) > 0 {
		output["c_hash"] = tok.CodeHash
	}

	// Merge each of the populated JSON templates into output
	err := mergeJSONTemplates(logger, output, templates...)
	if err != nil {
		logger.Error("failed to populate templates for ID token generation", "error", err)
		return nil, err
	}

	payload, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

// mergeJSONTemplates will merge each of the given JSON templates into the given
// output map. It will simply merge the top-level keys of the unmarshalled JSON
// templates into output, which means that any conflicting keys will be overwritten.
func mergeJSONTemplates(logger hclog.Logger, output map[string]interface{}, templates ...string) error {
	for _, template := range templates {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(template), &parsed); err != nil {
			logger.Warn("error parsing OIDC template", "template", template, "err", err)
		}

		for k, v := range parsed {
			if !strutil.StrListContains(reservedClaims, k) {
				output[k] = v
			} else {
				logger.Warn("invalid top level OIDC template key", "template", template, "key", k)
			}
		}
	}

	return nil
}

// generateAndSetKey will generate new signing and public key pairs and set
// them as the SigningKey.
func (k *namedKey) generateAndSetKey(ctx context.Context, logger hclog.Logger, s logical.Storage) error {
	signingKey, err := generateKeys(k.Algorithm)
	if err != nil {
		return err
	}

	k.SigningKey = signingKey
	k.KeyRing = append(k.KeyRing, &expireableKey{KeyID: signingKey.Public().KeyID})

	if err := saveOIDCPublicKey(ctx, s, signingKey.Public()); err != nil {
		return err
	}
	logger.Debug("generated OIDC public key to sign JWTs", "key_id", signingKey.Public().KeyID)
	return nil
}

// generateAndSetNextKey will generate new signing and public key pairs and set
// them as the NextSigningKey.
func (k *namedKey) generateAndSetNextKey(ctx context.Context, logger hclog.Logger, s logical.Storage) error {
	signingKey, err := generateKeys(k.Algorithm)
	if err != nil {
		return err
	}

	k.NextSigningKey = signingKey
	k.KeyRing = append(k.KeyRing, &expireableKey{KeyID: signingKey.Public().KeyID})

	if err := saveOIDCPublicKey(ctx, s, signingKey.Public()); err != nil {
		return err
	}
	logger.Debug("generated OIDC public key for future use", "key_id", signingKey.Public().KeyID)
	return nil
}

func (k *namedKey) signPayload(payload []byte) (string, error) {
	if k.SigningKey == nil {
		return "", fmt.Errorf("signing key is nil; rotate the key and try again")
	}
	signingKey := jose.SigningKey{Key: k.SigningKey, Algorithm: jose.SignatureAlgorithm(k.Algorithm)}
	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", err
	}

	signature, err := signer.Sign(payload)
	if err != nil {
		return "", err
	}

	signedIdToken, err := signature.CompactSerialize()
	if err != nil {
		return "", err
	}

	return signedIdToken, nil
}

func (i *IdentityStore) pathOIDCRoleExistenceCheck(ctx context.Context, req *logical.Request, d *framework.FieldData) (bool, error) {
	role, err := i.getOIDCRole(ctx, req.Storage, d.Get("name").(string))
	if err != nil {
		return false, err
	}

	return role != nil, nil
}

// pathOIDCCreateUpdateRole is used to create a new role or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	name := d.Get("name").(string)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	var role role
	if req.Operation == logical.UpdateOperation {
		entry, err := req.Storage.Get(ctx, roleConfigPath+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&role); err != nil {
				return nil, err
			}
		}
	}

	if key, ok := d.GetOk("key"); ok {
		role.Key = key.(string)
	} else if req.Operation == logical.CreateOperation {
		role.Key = d.Get("key").(string)
	}

	if role.Key == "" {
		return logical.ErrorResponse("the key parameter is required"), nil
	}

	if role.Key == defaultKeyName {
		if err := i.lazyGenerateDefaultKey(ctx, req.Storage); err != nil {
			return nil, fmt.Errorf("failed to generate default key: %w", err)
		}
	}

	if template, ok := d.GetOk("template"); ok {
		role.Template = template.(string)
	} else if req.Operation == logical.CreateOperation {
		role.Template = d.Get("template").(string)
	}

	// Attempt to decode as base64 and use that if it works
	if decoded, err := base64.StdEncoding.DecodeString(role.Template); err == nil {
		role.Template = string(decoded)
	}

	// Validate that template can be parsed and results in valid JSON
	if role.Template != "" {
		_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
			Mode:   identitytpl.JSONTemplating,
			String: role.Template,
			Entity: new(logical.Entity),
			Groups: make([]*logical.Group, 0),
			// namespace?
		})
		if err != nil {
			return logical.ErrorResponse("error parsing template: %s", err.Error()), nil
		}

		var tmp map[string]interface{}
		if err := json.Unmarshal([]byte(populatedTemplate), &tmp); err != nil {
			return logical.ErrorResponse("error parsing template JSON: %s", err.Error()), nil
		}

		for key := range tmp {
			if strutil.StrListContains(reservedClaims, key) {
				return logical.ErrorResponse(`top level key %q not allowed. Restricted keys: %s`,
					key, strings.Join(reservedClaims, ", ")), nil
			}
		}
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		role.TokenTTL = time.Duration(ttl.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.TokenTTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	// get the key referenced by this role if it exists
	var key namedKey
	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+role.Key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("cannot find key %q", role.Key), nil
	}

	if err := entry.DecodeJSON(&key); err != nil {
		return nil, err
	}
	if role.TokenTTL > key.VerificationTTL {
		return logical.ErrorResponse("a role's token ttl cannot be longer than the verification_ttl of the key it references"), nil
	}

	if clientID, ok := d.GetOk("client_id"); ok {
		role.ClientID = clientID.(string)
	}

	// create role path
	if role.ClientID == "" {
		clientID, err := base62.Random(26)
		if err != nil {
			return nil, err
		}
		role.ClientID = clientID
	}

	// store role (which was either just created or updated)
	entry, err = logical.StorageEntryJSON(roleConfigPath+name, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		return nil, err
	}

	return nil, nil
}

// handleOIDCReadRole is used to read an existing role
func (i *IdentityStore) pathOIDCReadRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	role, err := i.getOIDCRole(ctx, req.Storage, name)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"client_id": role.ClientID,
			"key":       role.Key,
			"template":  role.Template,
			"ttl":       int64(role.TokenTTL.Seconds()),
		},
	}, nil
}

func (i *IdentityStore) getOIDCRole(ctx context.Context, s logical.Storage, roleName string) (*role, error) {
	entry, err := s.Get(ctx, roleConfigPath+roleName)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var role role
	if err := entry.DecodeJSON(&role); err != nil {
		return nil, err
	}

	return &role, nil
}

// handleOIDCDeleteRole is used to delete a role if it exists
func (i *IdentityStore) pathOIDCDeleteRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	err := req.Storage.Delete(ctx, roleConfigPath+name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// handleOIDCListRole is used to list stored a roles
func (i *IdentityStore) pathOIDCListRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roles, err := req.Storage.List(ctx, roleConfigPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(roles), nil
}

func (i *IdentityStore) pathOIDCDiscovery(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data []byte

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var child string
	if childRaw, ok := d.GetOk("child"); ok {
		child = childRaw.(string)
	}

	cacheKey := fmt.Sprintf("%s/discoveryResponse", child)
	v, ok, err := i.oidcCache.Get(ns, cacheKey)
	if err != nil {
		return nil, err
	}

	if ok {
		data = v.([]byte)
	} else {
		c, err := i.getOIDCConfig(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		issuer, err := c.fullIssuer(child)
		if err != nil {
			return logical.ErrorResponse(err.Error()), nil
		}

		disc := discovery{
			Issuer:        issuer,
			Keys:          issuer + "/.well-known/keys",
			ResponseTypes: []string{"id_token"},
			Subjects:      []string{"public"},
			IDTokenAlgs:   supportedAlgs,
		}

		data, err = json.Marshal(disc)
		if err != nil {
			return nil, err
		}

		if err := i.oidcCache.SetDefault(ns, cacheKey, data); err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:         200,
			logical.HTTPRawBody:            data,
			logical.HTTPContentType:        "application/json",
			logical.HTTPCacheControlHeader: "max-age=3600",
		},
	}

	return resp, nil
}

// getKeysCacheControlHeader returns the cache control header for all public
// keys at the .well-known/keys endpoint
func (i *IdentityStore) getKeysCacheControlHeader(ns *namespace.Namespace) (string, error) {
	// if jwksCacheControlMaxAge is set use that, otherwise fall back on the
	// more conservative nextRun values
	jwksCacheControlMaxAge, ok, err := i.oidcCache.Get(ns, "jwksCacheControlMaxAge")
	if err != nil {
		return "", err
	}

	if ok {
		maxDuration := int64(jwksCacheControlMaxAge.(time.Duration))
		randDuration := mathrand.Int63n(maxDuration)
		durationInSeconds := time.Duration(randDuration).Seconds()
		return fmt.Sprintf("max-age=%.0f", durationInSeconds), nil
	}

	nextRun, ok, err := i.oidcCache.Get(ns, "nextRun")
	if err != nil {
		return "", err
	}

	if ok {
		now := time.Now()
		expireAt := nextRun.(time.Time)
		if expireAt.After(now) {
			i.Logger().Debug("use nextRun value for Cache Control header", "nextRun", nextRun)
			expireInSeconds := expireAt.Sub(time.Now()).Seconds()
			return fmt.Sprintf("max-age=%.0f", expireInSeconds), nil
		}
	}
	return "", nil
}

// pathOIDCReadPublicKeys is used to retrieve all public keys so that clients can
// verify the validity of a signed OIDC token.
func (i *IdentityStore) pathOIDCReadPublicKeys(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data []byte

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	var child string
	if childRaw, ok := d.GetOk("child"); ok {
		child = childRaw.(string)
	}
	if !validChildIssuer(child) {
		return logical.ErrorResponse("invalid child issuer %q", child), nil
	}

	v, ok, err := i.oidcCache.Get(ns, "jwksResponse")
	if err != nil {
		return nil, err
	}

	if ok {
		data = v.([]byte)
	} else {
		jwks, err := i.generatePublicJWKS(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		data, err = json.Marshal(jwks)
		if err != nil {
			return nil, err
		}

		if err := i.oidcCache.SetDefault(ns, "jwksResponse", data); err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:  200,
			logical.HTTPRawBody:     data,
			logical.HTTPContentType: "application/json",
		},
	}

	// set a Cache-Control header only if there are keys
	keys, err := listOIDCPublicKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if len(keys) > 0 {
		header, err := i.getKeysCacheControlHeader(ns)
		if err != nil {
			return nil, err
		}

		if header != "" {
			resp.Data[logical.HTTPCacheControlHeader] = header
		}
	}

	return resp, nil
}

func (i *IdentityStore) pathOIDCIntrospect(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var claims jwt.Claims

	// helper for preparing the non-standard introspection response
	introspectionResp := func(errorMsg string) (*logical.Response, error) {
		response := map[string]interface{}{
			"active": true,
		}

		if errorMsg != "" {
			response["active"] = false
			response["error"] = errorMsg
		}

		data, err := json.Marshal(response)
		if err != nil {
			return nil, err
		}

		resp := &logical.Response{
			Data: map[string]interface{}{
				logical.HTTPStatusCode:  200,
				logical.HTTPRawBody:     data,
				logical.HTTPContentType: "application/json",
			},
		}

		return resp, nil
	}

	rawIDToken := d.Get("token").(string)
	clientID := d.Get("client_id").(string)

	// validate basic JWT structure
	parsedJWT, err := jwt.ParseSigned(rawIDToken)
	if err != nil {
		return introspectionResp(fmt.Sprintf("error parsing token: %s", err.Error()))
	}

	// validate signature
	jwks, err := i.generatePublicJWKS(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	var valid bool
	for _, key := range jwks.Keys {
		if err := parsedJWT.Claims(key, &claims); err == nil {
			valid = true
			break
		}
	}

	if !valid {
		return introspectionResp("unable to validate the token signature")
	}

	// validate claims
	c, err := i.getOIDCConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	issuer, err := c.fullIssuer(baseIdentityTokenIssuer)
	if err != nil {
		return logical.ErrorResponse(err.Error()), nil
	}

	expected := jwt.Expected{
		Issuer: issuer,
		Time:   time.Now(),
	}

	if clientID != "" {
		expected.Audience = []string{clientID}
	}

	if claimsErr := claims.Validate(expected); claimsErr != nil {
		return introspectionResp(fmt.Sprintf("error validating claims: %s", claimsErr.Error()))
	}

	// validate entity exists and is active
	entity, err := i.MemDBEntityByID(claims.Subject, true)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return introspectionResp("entity was not found")
	} else if entity.Disabled {
		return introspectionResp("entity is disabled")
	}

	return introspectionResp("")
}

// namedKey.rotate(overrides) performs a key rotation on a namedKey.
// verification_ttl can be overridden with an overrideVerificationTTL value >= 0
func (k *namedKey) rotate(ctx context.Context, logger hclog.Logger, s logical.Storage, overrideVerificationTTL time.Duration) error {
	verificationTTL := k.VerificationTTL
	if overrideVerificationTTL >= 0 {
		verificationTTL = overrideVerificationTTL
	}

	now := time.Now()
	if k.SigningKey != nil {
		// set the previous public key's expiry time
		for _, key := range k.KeyRing {
			if key.KeyID == k.SigningKey.KeyID {
				key.ExpireAt = now.Add(verificationTTL)
				break
			}
		}
	} else {
		// this can occur for keys generated before vault 1.9.0 but rotated on
		// vault 1.9.0
		logger.Debug("nil signing key detected on rotation")
	}

	if k.NextSigningKey == nil {
		logger.Debug("nil next signing key detected on rotation")
		// keys will not have a NextSigningKey if they were generated before
		// vault 1.9
		err := k.generateAndSetNextKey(ctx, logger, s)
		if err != nil {
			return err
		}
	}

	// do the rotation
	k.SigningKey = k.NextSigningKey
	k.NextRotation = now.Add(k.RotationPeriod)

	// now that we have rotated, generate a new NextSigningKey
	err := k.generateAndSetNextKey(ctx, logger, s)
	if err != nil {
		return err
	}

	// store named key (it was modified when rotate was called on it)
	entry, err := logical.StorageEntryJSON(namedKeyConfigPath+k.name, k)
	if err != nil {
		return err
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	logger.Debug("rotated OIDC public key, now using", "key_id", k.SigningKey.Public().KeyID)
	return nil
}

// generateKeys returns a signingKey and publicKey pair
func generateKeys(algorithm string) (*jose.JSONWebKey, error) {
	var key interface{}
	var err error

	switch algorithm {
	case "RS256", "RS384", "RS512":
		// 2048 bits is recommended by RSA Laboratories as a minimum post 2015
		if key, err = cryptoutil.GenerateRSAKey(rand.Reader, 2048); err != nil {
			return nil, err
		}
	case "ES256", "ES384", "ES512":
		var curve elliptic.Curve

		switch algorithm {
		case "ES256":
			curve = elliptic.P256()
		case "ES384":
			curve = elliptic.P384()
		case "ES512":
			curve = elliptic.P521()
		}

		if key, err = ecdsa.GenerateKey(curve, rand.Reader); err != nil {
			return nil, err
		}
	case "EdDSA":
		_, key, err = ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown algorithm %q", algorithm)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	jwk := &jose.JSONWebKey{
		Key:       key,
		KeyID:     id,
		Algorithm: algorithm,
		Use:       "sig",
	}

	return jwk, nil
}

func saveOIDCPublicKey(ctx context.Context, s logical.Storage, key jose.JSONWebKey) error {
	entry, err := logical.StorageEntryJSON(publicKeysConfigPath+key.KeyID, key)
	if err != nil {
		return err
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

func loadOIDCPublicKey(ctx context.Context, s logical.Storage, keyID string) (*jose.JSONWebKey, error) {
	entry, err := s.Get(ctx, publicKeysConfigPath+keyID)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return nil, fmt.Errorf("could not find key with ID %s", keyID)
	}

	var key jose.JSONWebKey
	if err := entry.DecodeJSON(&key); err != nil {
		return nil, err
	}

	return &key, nil
}

func listOIDCPublicKeys(ctx context.Context, s logical.Storage) ([]string, error) {
	keys, err := s.List(ctx, publicKeysConfigPath)
	if err != nil {
		return nil, err
	}

	return keys, nil
}

func (i *IdentityStore) generatePublicJWKS(ctx context.Context, s logical.Storage) (*jose.JSONWebKeySet, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	jwksRaw, ok, err := i.oidcCache.Get(ns, "jwks")
	if err != nil {
		return nil, err
	}

	if ok {
		return jwksRaw.(*jose.JSONWebKeySet), nil
	}

	i.generateJWKSLock.Lock()
	defer i.generateJWKSLock.Unlock()

	// Check the cache again incase another requset acquired the lock
	// before this request.
	jwksRaw, ok, err = i.oidcCache.Get(ns, "jwks")
	if err != nil {
		return nil, err
	}

	if ok {
		return jwksRaw.(*jose.JSONWebKeySet), nil
	}

	if _, err := i.expireOIDCPublicKeys(ctx, s); err != nil {
		return nil, err
	}

	// Only return keys that are associated with a role or plugin mount
	// by collecting and de-duplicating keys and key IDs for each
	keyNames := make(map[string]struct{})
	keyIDs := make(map[string]struct{})

	// First collect the set of unique key names
	roleNames, err := s.List(ctx, roleConfigPath)
	if err != nil {
		return nil, err
	}
	for _, roleName := range roleNames {
		role, err := i.getOIDCRole(ctx, s, roleName)
		if err != nil {
			return nil, err
		}
		if role == nil {
			continue
		}

		keyNames[role.Key] = struct{}{}
	}
	mounts, err := i.listMounts(ns)
	if err != nil {
		return nil, err
	}
	for _, me := range mounts {
		key := defaultKeyName
		if me.Config.IdentityTokenKey != "" {
			key = me.Config.IdentityTokenKey
		}

		keyNames[key] = struct{}{}
	}

	// Second collect the set of unique key IDs for each key name
	for name := range keyNames {
		ids, err := i.keyIDsByName(ctx, s, name)
		if err != nil {
			return nil, err
		}

		for _, id := range ids {
			keyIDs[id] = struct{}{}
		}
	}

	jwks := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0, len(keyIDs)),
	}

	// load the JSON web key for each key ID
	for keyID := range keyIDs {
		key, err := loadOIDCPublicKey(ctx, s, keyID)
		if err != nil {
			return nil, err
		}
		jwks.Keys = append(jwks.Keys, *key)
	}

	if err := i.oidcCache.SetDefault(ns, "jwks", jwks); err != nil {
		return nil, err
	}

	return jwks, nil
}

func (i *IdentityStore) expireOIDCPublicKeys(ctx context.Context, s logical.Storage) (time.Time, error) {
	var didUpdate bool

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return time.Time{}, err
	}

	// nextExpiration will be the soonest expiration time of all keys. Initialize
	// here to a relatively distant time.
	nextExpiration := time.Now().Add(24 * time.Hour)
	now := time.Now()

	publicKeyIDs, err := listOIDCPublicKeys(ctx, s)
	if err != nil {
		return now, err
	}

	keyNames, err := s.List(ctx, namedKeyConfigPath)
	if err != nil {
		return now, err
	}

	usedKeys := make([]string, 0)

	for _, keyName := range keyNames {
		entry, err := s.Get(ctx, namedKeyConfigPath+keyName)
		if err != nil {
			return now, err
		}

		if entry == nil {
			i.Logger().Warn("could not find key to update", "key", keyName)
			continue
		}

		var key namedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return now, err
		}

		// Remove any expired keys from the keyring.
		keyRing := key.KeyRing
		var keyringUpdated bool

		for j := 0; j < len(keyRing); j++ {
			k := keyRing[j]
			if !k.ExpireAt.IsZero() && k.ExpireAt.Before(now) {
				keyRing[j] = keyRing[len(keyRing)-1]
				keyRing = keyRing[:len(keyRing)-1]

				keyringUpdated = true
				j--
				continue
			}

			// Save a remaining key's next expiration if it is the earliest we've
			// seen (for use by the periodicFunc for scheduling).
			if !k.ExpireAt.IsZero() && k.ExpireAt.Before(nextExpiration) {
				nextExpiration = k.ExpireAt
			}

			// Mark the KeyId as in use so it doesn't get deleted in the next step
			usedKeys = append(usedKeys, k.KeyID)
		}

		// Persist any keyring updates if necessary
		if keyringUpdated {
			key.KeyRing = keyRing
			entry, err := logical.StorageEntryJSON(entry.Key, key)
			if err != nil {
				i.Logger().Error("error creating storage entry", "key", key.name, "error", err)
				continue
			}

			if err := s.Put(ctx, entry); err != nil {
				i.Logger().Error("error writing key", "key", key.name, "error", err)
				continue
			}
			didUpdate = true
		}
	}

	// Delete all public keys that were not determined to be not expired and in
	// use by some role.
	for _, keyID := range publicKeyIDs {
		if !strutil.StrListContains(usedKeys, keyID) {
			if err := s.Delete(ctx, publicKeysConfigPath+keyID); err != nil {
				i.Logger().Error("error deleting OIDC public key", "key_id", keyID, "error", err)
				nextExpiration = now
				continue
			}
			didUpdate = true
			i.Logger().Debug("deleted OIDC public key", "key_id", keyID)
		}
	}

	if didUpdate {
		if err := i.oidcCache.Flush(ns); err != nil {
			i.Logger().Error("error flushing oidc cache", "error", err)
		}
	}

	return nextExpiration, nil
}

// oidcKeyRotation will rotate any keys that are due to be rotated.
//
// It will return the time of the soonest rotation and the minimum
// verificationTTL or minimum rotationPeriod out of all the current keys.
func (i *IdentityStore) oidcKeyRotation(ctx context.Context, s logical.Storage) (time.Time, time.Duration, error) {
	// soonestRotation will be the soonest rotation time of all keys. Initialize
	// here to a relatively distant time.
	now := time.Now()
	soonestRotation := now.Add(24 * time.Hour)

	jwksClientCacheDuration := time.Duration(math.MaxInt64)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	keys, err := s.List(ctx, namedKeyConfigPath)
	if err != nil {
		return now, jwksClientCacheDuration, err
	}

	for _, k := range keys {
		entry, err := s.Get(ctx, namedKeyConfigPath+k)
		if err != nil {
			return now, jwksClientCacheDuration, err
		}

		if entry == nil {
			continue
		}

		var key namedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return now, jwksClientCacheDuration, err
		}
		key.name = k

		if key.VerificationTTL < jwksClientCacheDuration {
			jwksClientCacheDuration = key.VerificationTTL
		}

		if key.RotationPeriod < jwksClientCacheDuration {
			jwksClientCacheDuration = key.RotationPeriod
		}

		// Future key rotation that is the earliest we've seen.
		if now.Before(key.NextRotation) && key.NextRotation.Before(soonestRotation) {
			soonestRotation = key.NextRotation
		}

		// Key that is due to be rotated.
		if now.After(key.NextRotation) {
			i.Logger().Debug("rotating OIDC key", "key", key.name)
			if err := key.rotate(ctx, i.Logger(), s, -1); err != nil {
				return now, jwksClientCacheDuration, err
			}

			// Possibly save the new rotation time
			if key.NextRotation.Before(soonestRotation) {
				soonestRotation = key.NextRotation
			}
		}
	}

	return soonestRotation, jwksClientCacheDuration, nil
}

// oidcPeriodFunc is invoked by the backend's periodFunc and runs regular key
// rotations and expiration actions.
func (i *IdentityStore) oidcPeriodicFunc(ctx context.Context, s logical.Storage) {
	// Key rotations write to storage, so only run this on the primary cluster.
	// The periodic func does not run on perf standbys or DR secondaries.
	if i.System().ReplicationState().HasState(consts.ReplicationPerformanceSecondary) {
		return
	}

	now := time.Now()

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		i.Logger().Error("error getting namespace from context", "err", err)
		return
	}

	var nextRun time.Time
	v, ok, err := i.oidcCache.Get(ns, "nextRun")
	if err != nil {
		i.Logger().Error("error reading oidc cache", "err", err)
		return
	}

	if ok {
		nextRun = v.(time.Time)
	}

	// The condition here is for performance, not precise timing. The actions can
	// be run at any time safely, but there is no need to invoke them (which
	// might be somewhat expensive if there are many roles/keys) if we're not
	// past any rotation/expiration TTLs.
	if now.Before(nextRun) {
		return
	}

	nextRotation, jwksClientCacheDuration, err := i.oidcKeyRotation(ctx, s)
	if err != nil {
		i.Logger().Warn("error rotating OIDC keys", "err", err)
	}

	nextExpiration, err := i.expireOIDCPublicKeys(ctx, s)
	if err != nil {
		i.Logger().Warn("error expiring OIDC public keys", "err", err)
	}

	if err := i.oidcCache.Flush(ns); err != nil {
		i.Logger().Error("error flushing oidc cache", "err", err)
	}

	// use the soonest time between nextRotation and nextExpiration for the next run.
	// Allow at most 24 hours though, keeping the legacy behavior from the original
	// introduction of namespaces (unclear if necessary but safer to keep for now).
	nextRun = now.Add(24 * time.Hour)
	if nextRotation.Before(nextRun) {
		nextRun = nextRotation
	}
	if nextExpiration.Before(nextRun) {
		nextRun = nextExpiration
	}

	if err := i.oidcCache.SetDefault(ns, "nextRun", nextRun); err != nil {
		i.Logger().Error("error setting oidc cache", "err", err)
	}

	if jwksClientCacheDuration < math.MaxInt64 {
		// the OIDC JWKS endpoint returns a Cache-Control HTTP header time between
		// 0 and the minimum verificationTTL or minimum rotationPeriod out of all
		// keys, whichever value is lower.
		//
		// This smooths calls from services validating JWTs to Vault, while
		// ensuring that operators can assert that servers honoring the
		// Cache-Control header will always have a superset of all valid keys, and
		// not trust any keys longer than a jwksCacheControlMaxAge duration after a
		// key is rotated out of signing use
		if err := i.oidcCache.SetDefault(ns, "jwksCacheControlMaxAge", jwksClientCacheDuration); err != nil {
			i.Logger().Error("error setting jwksCacheControlMaxAge in oidc cache", "err", err)
		}
	}
}

func newOIDCCache(defaultExpiration, cleanupInterval time.Duration) *oidcCache {
	return &oidcCache{
		c: cache.New(defaultExpiration, cleanupInterval),
	}
}

func (c *oidcCache) nskey(ns *namespace.Namespace, key string) string {
	return fmt.Sprintf("v0:%s:%s", ns.ID, key)
}

func (c *oidcCache) Get(ns *namespace.Namespace, key string) (interface{}, bool, error) {
	if ns == nil {
		return nil, false, errNilNamespace
	}
	v, found := c.c.Get(c.nskey(ns, key))
	return v, found, nil
}

func (c *oidcCache) SetDefault(ns *namespace.Namespace, key string, obj interface{}) error {
	if ns == nil {
		return errNilNamespace
	}
	c.c.SetDefault(c.nskey(ns, key), obj)

	return nil
}

func (c *oidcCache) Delete(ns *namespace.Namespace, key string) error {
	if ns == nil {
		return errNilNamespace
	}
	c.c.Delete(c.nskey(ns, key))

	return nil
}

func (c *oidcCache) Flush(ns *namespace.Namespace) error {
	if ns == nil {
		return errNilNamespace
	}

	// Remove all items from the provided namespace
	for itemKey := range c.c.Items() {
		if isTargetNamespacedKey(itemKey, []string{ns.ID}) {
			c.c.Delete(itemKey)
		}
	}

	return nil
}

// isTargetNamespacedKey returns true for a properly constructed namespaced key (<version>:<nsID>:<key>)
// where <nsID> matches any targeted nsID
func isTargetNamespacedKey(nskey string, nsTargets []string) bool {
	split := strings.Split(nskey, ":")
	return len(split) >= 3 && strutil.StrListContains(nsTargets, split[1])
}
