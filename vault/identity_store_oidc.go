package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/helper/namespace"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/base62"
	"github.com/hashicorp/vault/sdk/helper/identitytpl"
	"github.com/hashicorp/vault/sdk/helper/strutil"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/patrickmn/go-cache"
	"golang.org/x/crypto/ed25519"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type oidcConfig struct {
	Issuer string `json:"issuer"`

	// effectiveIssuer is a calculated field and will be either Issuer (if
	// that's set) or the Vault instance's api_addr.
	effectiveIssuer string
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
	Issuer    string `json:"iss"`       // api_addr or custom Issuer
	Namespace string `json:"namespace"` // Namespace of issuer
	Subject   string `json:"sub"`       // Entity ID
	Audience  string `json:"aud"`       // role ID will be used here.
	Expiry    int64  `json:"exp"`       // Expiration, as determined by the role.
	IssuedAt  int64  `json:"iat"`       // Time of token creation
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

var errNilNamespace = errors.New("nil namespace in oidc cache request")

const (
	issuerPath           = "identity/oidc"
	oidcTokensPrefix     = "oidc_tokens/"
	oidcConfigStorageKey = oidcTokensPrefix + "config/"
	namedKeyConfigPath   = oidcTokensPrefix + "named_keys/"
	publicKeysConfigPath = oidcTokensPrefix + "public_keys/"
	roleConfigPath       = oidcTokensPrefix + "roles/"
)

var (
	requiredClaims = []string{"iat", "aud", "exp", "iss", "sub", "namespace"}
	supportedAlgs  = []string{
		string(jose.RS256),
		string(jose.RS384),
		string(jose.RS512),
		string(jose.ES256),
		string(jose.ES384),
		string(jose.ES512),
		string(jose.EdDSA),
	}
)

// pseudo-namespace for cache items that don't belong to any real namespace.
var noNamespace = &namespace.Namespace{ID: "__NO_NAMESPACE"}

func oidcPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "oidc/config/?$",
			Fields: map[string]*framework.FieldSchema{
				"issuer": {
					Type:        framework.TypeString,
					Description: "Issuer URL to be used in the iss claim of the token. If not set, Vault's app_addr will be used.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation:   i.pathOIDCReadConfig,
				logical.UpdateOperation: i.pathOIDCUpdateConfig,
			},
			HelpSynopsis:    "OIDC configuration",
			HelpDescription: "Update OIDC configuration in the identity backend",
		},
		{
			Pattern: "oidc/key/" + framework.GenericNameRegex("name"),
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
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListKey,
			},
			HelpSynopsis:    "List OIDC keys",
			HelpDescription: "List all named OIDC keys",
		},
		{
			Pattern: "oidc/.well-known/openid-configuration/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCDiscovery,
			},
			HelpSynopsis:    "Query OIDC configurations",
			HelpDescription: "Query this path to retrieve the configured OIDC Issuer and Keys endpoints, response types, subject types, and signing algorithms used by the OIDC backend.",
		},
		{
			Pattern: "oidc/.well-known/keys/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCReadPublicKeys,
			},
			HelpSynopsis:    "Retrieve public keys",
			HelpDescription: "Query this path to retrieve the public portion of keys used to sign OIDC tokens. Clients can use this to validate the authenticity of the OIDC token claims.",
		},
		{
			Pattern: "oidc/token/" + framework.GenericNameRegex("name"),
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
			Fields: map[string]*framework.FieldSchema{
				"name": {
					Type:        framework.TypeString,
					Description: "Name of the role",
				},
				"key": {
					Type:        framework.TypeString,
					Description: "The OIDC key to use for generating tokens. The specified key must already exist.",
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
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListRole,
			},
			HelpSynopsis:    "List configured OIDC roles",
			HelpDescription: "List all configured OIDC roles in the identity backend.",
		},
		{
			Pattern: "oidc/introspect/?$",
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

	if i.core.redirectAddr == "" && c.Issuer == "" {
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
		c.effectiveIssuer = i.core.redirectAddr
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

	// Update next rotation time if it is unset or now earlier than previously set.
	nextRotation := time.Now().Add(key.RotationPeriod)
	if key.NextRotation.IsZero() || nextRotation.Before(key.NextRotation) {
		key.NextRotation = nextRotation
	}

	// generate keys if creating a new key or changing algorithms
	if key.Algorithm != prevAlgorithm {
		signingKey, err := generateKeys(key.Algorithm)
		if err != nil {
			return nil, err
		}

		key.SigningKey = signingKey
		key.KeyRing = append(key.KeyRing, &expireableKey{KeyID: signingKey.Public().KeyID})

		if err := saveOIDCPublicKey(ctx, req.Storage, signingKey.Public()); err != nil {
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

// rolesReferencingTargetKeyName returns a map of role names to roles referenced by targetKeyName.
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
// names referenced by targetKeyName.
// Note: this is not threadsafe. It is to be called with Lock already held.
func (i *IdentityStore) roleNamesReferencingTargetKeyName(ctx context.Context, req *logical.Request, targetKeyName string) ([]string, error) {
	roles, err := i.rolesReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		return nil, err
	}

	var names []string
	for key, _ := range roles {
		names = append(names, key)
	}
	sort.Strings(names)
	return names, nil
}

// handleOIDCDeleteKey is used to delete a key
func (i *IdentityStore) pathOIDCDeleteKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
	}

	targetKeyName := d.Get("name").(string)

	i.oidcLock.Lock()

	roleNames, err := i.roleNamesReferencingTargetKeyName(ctx, req, targetKeyName)
	if err != nil {
		return nil, err
	}

	if len(roleNames) > 0 {
		errorMessage := fmt.Sprintf("unable to delete key %q because it is currently referenced by these roles: %s",
			targetKeyName, strings.Join(roleNames, ", "))
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

	if err := storedNamedKey.rotate(ctx, req.Storage, verificationTTLOverride); err != nil {
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

	var key *namedKey

	keyRaw, found, err := i.oidcCache.Get(ns, "namedKeys/"+role.Key)
	if err != nil {
		return nil, err
	}

	if found {
		key = keyRaw.(*namedKey)
	} else {
		entry, _ := req.Storage.Get(ctx, namedKeyConfigPath+role.Key)
		if entry == nil {
			return logical.ErrorResponse("key %q not found", role.Key), nil
		}

		if err := entry.DecodeJSON(&key); err != nil {
			return nil, err
		}

		if err := i.oidcCache.SetDefault(ns, "namedKeys/"+role.Key, key); err != nil {
			return nil, err
		}
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

	now := time.Now()
	idToken := idToken{
		Issuer:    config.effectiveIssuer,
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

	payload, err := idToken.generatePayload(i.Logger(), role.Template, e, groups)
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

func (tok *idToken) generatePayload(logger hclog.Logger, template string, entity *identity.Entity, groups []*identity.Group) ([]byte, error) {
	output := map[string]interface{}{
		"iss":       tok.Issuer,
		"namespace": tok.Namespace,
		"sub":       tok.Subject,
		"aud":       tok.Audience,
		"exp":       tok.Expiry,
		"iat":       tok.IssuedAt,
	}

	// Parse and integrate the populated role template. Structural errors with the template _should_
	// be caught during role configuration. Error found during runtime will be logged, but they will
	// not block generation of the basic ID token. They should not be returned to the requester.
	_, populatedTemplate, err := identitytpl.PopulateString(identitytpl.PopulateStringInput{
		Mode:   identitytpl.JSONTemplating,
		String: template,
		Entity: identity.ToSDKEntity(entity),
		Groups: identity.ToSDKGroups(groups),
		// namespace?
	})
	if err != nil {
		logger.Warn("error populating OIDC token template", "template", template, "error", err)
	}

	if populatedTemplate != "" {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(populatedTemplate), &parsed); err != nil {
			logger.Warn("error parsing OIDC template", "template", template, "err", err)
		}

		for k, v := range parsed {
			if !strutil.StrListContains(requiredClaims, k) {
				output[k] = v
			} else {
				logger.Warn("invalid top level OIDC template key", "template", template, "key", k)
			}
		}
	}

	payload, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (k *namedKey) signPayload(payload []byte) (string, error) {
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
			if strutil.StrListContains(requiredClaims, key) {
				return logical.ErrorResponse(`top level key %q not allowed. Restricted keys: %s`,
					key, strings.Join(requiredClaims, ", ")), nil
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
	if entry != nil {
		if err := entry.DecodeJSON(&key); err != nil {
			return nil, err
		}

		if role.TokenTTL > key.VerificationTTL {
			return logical.ErrorResponse("a role's token ttl cannot be longer than the verification_ttl of the key it references"), nil
		}
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

	v, ok, err := i.oidcCache.Get(ns, "discoveryResponse")
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

		disc := discovery{
			Issuer:        c.effectiveIssuer,
			Keys:          c.effectiveIssuer + "/.well-known/keys",
			ResponseTypes: []string{"id_token"},
			Subjects:      []string{"public"},
			IDTokenAlgs:   supportedAlgs,
		}

		data, err = json.Marshal(disc)
		if err != nil {
			return nil, err
		}

		if err := i.oidcCache.SetDefault(ns, "discoveryResponse", data); err != nil {
			return nil, err
		}
	}

	resp := &logical.Response{
		Data: map[string]interface{}{
			logical.HTTPStatusCode:      200,
			logical.HTTPRawBody:         data,
			logical.HTTPContentType:     "application/json",
			logical.HTTPRawCacheControl: "max-age=3600",
		},
	}

	return resp, nil
}

// pathOIDCReadPublicKeys is used to retrieve all public keys so that clients can
// verify the validity of a signed OIDC token.
func (i *IdentityStore) pathOIDCReadPublicKeys(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data []byte

	ns, err := namespace.FromContext(ctx)
	if err != nil {
		return nil, err
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

	// set a Cache-Control header only if there are keys, if there aren't keys
	// then nextRun should not be used to set Cache-Control header because it chooses
	// a time in the future that isn't based on key rotation/expiration values
	keys, err := listOIDCPublicKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}
	if len(keys) > 0 {
		v, ok, err := i.oidcCache.Get(noNamespace, "nextRun")
		if err != nil {
			return nil, err
		}

		if ok {
			now := time.Now()
			expireAt := v.(time.Time)
			if expireAt.After(now) {
				expireInSeconds := expireAt.Sub(time.Now()).Seconds()
				expireInString := fmt.Sprintf("max-age=%.0f", expireInSeconds)
				resp.Data[logical.HTTPRawCacheControl] = expireInString
			}
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

	expected := jwt.Expected{
		Issuer: c.effectiveIssuer,
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

// namedKey.rotate(overrides) performs a key rotation on a namedKey and returns the
// verification_ttl that was applied. verification_ttl can be overridden with an
// overrideVerificationTTL value >= 0
func (k *namedKey) rotate(ctx context.Context, s logical.Storage, overrideVerificationTTL time.Duration) error {
	verificationTTL := k.VerificationTTL

	if overrideVerificationTTL >= 0 {
		verificationTTL = overrideVerificationTTL
	}

	// generate new key
	signingKey, err := generateKeys(k.Algorithm)
	if err != nil {
		return err
	}
	if err := saveOIDCPublicKey(ctx, s, signingKey.Public()); err != nil {
		return err
	}

	now := time.Now()

	// set the previous public key's expiry time
	for _, key := range k.KeyRing {
		if key.KeyID == k.SigningKey.KeyID {
			key.ExpireAt = now.Add(verificationTTL)
			break
		}
	}
	k.SigningKey = signingKey
	k.KeyRing = append(k.KeyRing, &expireableKey{KeyID: signingKey.KeyID})
	k.NextRotation = now.Add(k.RotationPeriod)

	// store named key (it was modified when rotate was called on it)
	entry, err := logical.StorageEntryJSON(namedKeyConfigPath+k.name, k)
	if err != nil {
		return err
	}
	if err := s.Put(ctx, entry); err != nil {
		return err
	}

	return nil
}

// generateKeys returns a signingKey and publicKey pair
func generateKeys(algorithm string) (*jose.JSONWebKey, error) {
	var key interface{}
	var err error

	switch algorithm {
	case "RS256", "RS384", "RS512":
		// 2048 bits is recommended by RSA Laboratories as a minimum post 2015
		if key, err = rsa.GenerateKey(rand.Reader, 2048); err != nil {
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

	if _, err := i.expireOIDCPublicKeys(ctx, s); err != nil {
		return nil, err
	}

	keyIDs, err := listOIDCPublicKeys(ctx, s)
	if err != nil {
		return nil, err
	}

	jwks := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 0, len(keyIDs)),
	}

	for _, keyID := range keyIDs {
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

	namedKeys, err := s.List(ctx, namedKeyConfigPath)
	if err != nil {
		return now, err
	}

	usedKeys := make([]string, 0, 2*len(namedKeys))

	for _, k := range namedKeys {
		entry, err := s.Get(ctx, namedKeyConfigPath+k)
		if err != nil {
			return now, err
		}

		var key namedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return now, err
		}

		// Remove any expired keys from the keyring.
		keyRing := key.KeyRing
		var keyringUpdated bool

		for i := 0; i < len(keyRing); i++ {
			k := keyRing[i]
			if !k.ExpireAt.IsZero() && k.ExpireAt.Before(now) {
				keyRing[i] = keyRing[len(keyRing)-1]
				keyRing = keyRing[:len(keyRing)-1]

				keyringUpdated = true
				i--
				continue
			}

			// Save a remaining key's next expiration if it is the earliest we've
			// seen (for use by the periodicFunc for scheduling).
			if !k.ExpireAt.IsZero() && k.ExpireAt.Before(nextExpiration) {
				nextExpiration = k.ExpireAt
			}

			// Mark the KeyID as in use so it doesn't get deleted in the next step
			usedKeys = append(usedKeys, k.KeyID)
		}

		// Persist any keyring updates if necessary
		if keyringUpdated {
			key.KeyRing = keyRing
			entry, err := logical.StorageEntryJSON(entry.Key, key)
			if err != nil {
				i.Logger().Error("error updating key", "key", key.name, "error", err)
			}

			if err := s.Put(ctx, entry); err != nil {
				i.Logger().Error("error saving key", "key", key.name, "error", err)
			}
			didUpdate = true
		}
	}

	// Delete all public keys that were not determined to be not expired and in
	// use by some role.
	for _, keyID := range publicKeyIDs {
		if !strutil.StrListContains(usedKeys, keyID) {
			didUpdate = true
			if err := s.Delete(ctx, publicKeysConfigPath+keyID); err != nil {
				i.Logger().Error("error deleting OIDC public key", "key_id", keyID, "error", err)
				nextExpiration = now
			}
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

func (i *IdentityStore) oidcKeyRotation(ctx context.Context, s logical.Storage) (time.Time, error) {
	// soonestRotation will be the soonest rotation time of all keys. Initialize
	// here to a relatively distant time.
	now := time.Now()
	soonestRotation := now.Add(24 * time.Hour)

	i.oidcLock.Lock()
	defer i.oidcLock.Unlock()

	keys, err := s.List(ctx, namedKeyConfigPath)
	if err != nil {
		return now, err
	}

	for _, k := range keys {
		entry, err := s.Get(ctx, namedKeyConfigPath+k)
		if err != nil {
			return now, err
		}

		if entry == nil {
			continue
		}

		var key namedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return now, err
		}
		key.name = k

		// Future key rotation that is the earliest we've seen.
		if now.Before(key.NextRotation) && key.NextRotation.Before(soonestRotation) {
			soonestRotation = key.NextRotation
		}

		// Key that is due to be rotated.
		if now.After(key.NextRotation) {
			i.Logger().Debug("rotating OIDC key", "key", key.name)
			if err := key.rotate(ctx, s, -1); err != nil {
				return now, err
			}

			// Possibly save the new rotation time
			if key.NextRotation.Before(soonestRotation) {
				soonestRotation = key.NextRotation
			}
		}
	}

	return soonestRotation, nil
}

// oidcPeriodFunc is invoked by the backend's periodFunc and runs regular key
// rotations and expiration actions.
func (i *IdentityStore) oidcPeriodicFunc(ctx context.Context) {
	var nextRun time.Time
	now := time.Now()

	v, ok, err := i.oidcCache.Get(noNamespace, "nextRun")
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
	if now.After(nextRun) {
		// Initialize to a fairly distant next run time. This will be brought in
		// based on key rotation times.
		nextRun = now.Add(24 * time.Hour)

		for _, ns := range i.listNamespaces() {
			nsPath := ns.Path

			s := i.core.router.MatchingStorageByAPIPath(ctx, nsPath+"identity/oidc")

			if s == nil {
				continue
			}

			nextRotation, err := i.oidcKeyRotation(ctx, s)
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

			// re-run at the soonest expiration or rotation time
			if nextRotation.Before(nextRun) {
				nextRun = nextRotation
			}

			if nextExpiration.Before(nextRun) {
				nextRun = nextExpiration
			}
		}
		if err := i.oidcCache.SetDefault(noNamespace, "nextRun", nextRun); err != nil {
			i.Logger().Error("error setting oidc cache", "err", err)
		}
	}
}

func newOIDCCache() *oidcCache {
	return &oidcCache{
		c: cache.New(cache.NoExpiration, cache.NoExpiration),
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

func (c *oidcCache) Flush(ns *namespace.Namespace) error {
	if ns == nil {
		return errNilNamespace
	}

	// Remove all items from the provided namespace as well as the shared, "no namespace" section.
	for itemKey := range c.c.Items() {
		if isTargetNamespacedKey(itemKey, []string{noNamespace.ID, ns.ID}) {
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
