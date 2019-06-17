package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/vault/sdk/helper/strutil"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/helper/identity"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
	"gopkg.in/square/go-jose.v2"
)

// todo fix this
var requiredClaims = []string{"iat", "aud", "exp", "iss", "sub"}

const (
	oidcTokensPrefix     = "oidc_tokens/"
	oidcConfigStorageKey = oidcTokensPrefix + "config/"
	namedKeyConfigPath   = oidcTokensPrefix + "named_keys/"
	publicKeysConfigPath = oidcTokensPrefix + "public_keys/"
	roleConfigPath       = oidcTokensPrefix + "roles/"
)

type oidcConfig struct {
	Issuer string `json:"issuer"`
}

type ExpireableKey struct {
	KeyID    string    `json:"key_id"`
	ExpireAt time.Time `json:"expire_at"`
}

type NamedKey struct {
	Name             string           `json:"name"`
	SigningAlgorithm string           `json:"signing_algorithm"`
	VerificationTTL  int              `json:"verification_ttl"`
	RotationPeriod   int              `json:"rotation_period"`
	KeyRing          []*ExpireableKey `json:"key_ring"`
	SigningKey       jose.JSONWebKey  `json:"signing_key"`
	NextRotation     time.Time        `json:"next_rotation""`
}

type Role struct {
	Name     string        `json:"name"` // TODO: do we need/want this?
	TokenTTL time.Duration `json:"token_ttl"`
	Key      string        `json:"key"`
	Template string        `json:"template"`
	RoleID   string        `json:"role_id"`
}

type discovery struct {
	Issuer        string   `json:"issuer"`
	Auth          string   `json:"authorization_endpoint"`
	Token         string   `json:"token_endpoint"`
	Keys          string   `json:"jwks_uri"`
	ResponseTypes []string `json:"response_types_supported"`
	Subjects      []string `json:"subject_types_supported"`
	IDTokenAlgs   []string `json:"id_token_signing_alg_values_supported"`
	Scopes        []string `json:"scopes_supported"`
	AuthMethods   []string `json:"token_endpoint_auth_methods_supported"`
	Claims        []string `json:"claims_supported"`
}

// oidcCache stores common response data as well as when the periodic func needs
// to run. This is conservatively managed, and most writes to the OIDC endpoints
// will invalidate the cache.
type oidcCache struct {
	lock              sync.RWMutex
	discoveryResponse []byte
	jwksResponse      []byte
	nextPeriodicRun   time.Time
}

func (c *oidcCache) purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.discoveryResponse = c.discoveryResponse[:0]
	c.jwksResponse = c.jwksResponse[:0]
	c.nextPeriodicRun = time.Time{}
}

// oidcPaths returns the API endpoints supported to operate on OIDC tokens:
// oidc/key/:key - Create a new key named key
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
					Description: "How often to generate a new keypair. Defaults to 6h",
				},

				"verification_ttl": {
					Type:        framework.TypeDurationSecond,
					Description: "Controls how long the public portion of a key will be available for verification after being rotated. Defaults to the current rotation_period, which will provide for a current key and previous key.",
				},

				"algorithm": {
					Type:        framework.TypeString,
					Description: "Signing algorithm to use. This will default to RS256, and is currently the only allowed value.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCCreateUpdateKey,
				logical.ReadOperation:   i.pathOIDCReadKey,
				logical.DeleteOperation: i.pathOIDCDeleteKey,
			},
			HelpSynopsis:    "oidc/key/:key help synopsis here",
			HelpDescription: "oidc/key/:key help description here",
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
			HelpSynopsis:    "oidc/key/:key/rotate help synopsis here",
			HelpDescription: "oidc/key/:key/rotate help description here",
		},
		{
			Pattern: "oidc/key/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListKey,
			},
			HelpSynopsis:    "list oidc/key/ help synopsis here",
			HelpDescription: "list oidc/key/ help description here",
		},
		{
			Pattern: "oidc/.well-known/openid-configuration/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCDiscovery,
			},
			HelpSynopsis:    "read oidc/.well-known/discovery/ help synopsis here",
			HelpDescription: "read oidc/.well-known/discovery/ help description here",
		},
		{
			Pattern: "oidc/.well-known/certs/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.pathOIDCReadPublicKeys,
			},
			HelpSynopsis:    "read oidc/.well-known/certs/ help synopsis here",
			HelpDescription: "read oidc/.well-known/certs/ help description here",
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
					Default:     "5m",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.pathOIDCUpdateRole,
				logical.CreateOperation: i.pathOIDCUpdateRole,
				logical.ReadOperation:   i.pathOIDCReadRole,
				logical.DeleteOperation: i.pathOIDCDeleteRole,
			},
			ExistenceCheck: i.pathOIDCRoleExistenceCheck,
		},
		{
			Pattern: "oidc/role/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.pathOIDCListRole,
			},
			HelpSynopsis:    "list oidc/role/ help synopsis here",
			HelpDescription: "list oidc/role/ help description here",
		},
	}
}

// readOIDCConfig returns a framework.OperationFunc for reading OIDC configuration
func (i *IdentityStore) pathOIDCReadConfig(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	c, err := i.getOIDCConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	if c == nil {
		return nil, nil
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"issuer": c.Issuer,
		},
	}, nil
}

// writeOIDCConfig returns a framework.OperationFunc for creating and updating OIDC configuration
func (i *IdentityStore) pathOIDCUpdateConfig(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	value, ok := d.GetOk("issuer")
	if !ok {
		return nil, nil
	}

	c := oidcConfig{
		Issuer: value.(string),
	}

	entry, err := logical.StorageEntryJSON(oidcConfigStorageKey, c)
	if err != nil {
		return nil, err
	}

	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	var resp logical.Response

	if c.Issuer != "" {
		resp.AddWarning(`If "issuer" is set explicitly, all tokens must be ` +
			`validated against that address, including those issued by secondary ` +
			`clusters. Setting issuer to "" will restore the default behavior of ` +
			`using the cluster's api_addr as the issuer.`)
	}

	i.oidcCache.purge()

	return &resp, nil
}

func (i *IdentityStore) getOIDCConfig(ctx context.Context, s logical.Storage) (*oidcConfig, error) {
	var c oidcConfig

	entry, err := s.Get(ctx, oidcConfigStorageKey)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	if err := entry.DecodeJSON(&c); err != nil {
		return nil, err
	}

	return &c, nil
}

// handleOIDCCreateKey is used to create a new named key or update an existing one
func (i *IdentityStore) pathOIDCCreateUpdateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var namedKey *NamedKey
	var update bool = false

	// parse parameters
	name := d.Get("name").(string)
	rotationPeriodInput, rotationPeriodInputOK := d.GetOk("rotation_period")
	verificationTTLInput, verificationTTLInputOK := d.GetOk("verification_ttl")
	algorithm := d.Get("algorithm").(string)

	// determine if we are creating a new key or updating an existing key
	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		if err := entry.DecodeJSON(&namedKey); err != nil {
			return nil, err
		}
		update = true
	}

	var rotationPeriod int
	var verificationTTL int

	// set defaults if creating a new key
	if !update {
		if rotationPeriodInputOK {
			rotationPeriod = rotationPeriodInput.(int)
		} else {
			rotationPeriod = 6 * 60 * 60 // 6h in seconds
		}

		if verificationTTLInputOK {
			verificationTTL = verificationTTLInput.(int)
		} else {
			verificationTTL = rotationPeriod
		}

		if algorithm == "" {
			algorithm = "RS256"
		}
	}

	// set values on the namedkey if they were provided and this is an update
	if update {
		if rotationPeriodInputOK {
			namedKey.RotationPeriod = rotationPeriodInput.(int)

			// TODO: this should probably be the old/new NextRotation?
			namedKey.NextRotation = time.Now().Add(time.Duration(rotationPeriod) * time.Second)
		}

		if verificationTTLInputOK {
			namedKey.VerificationTTL = verificationTTLInput.(int)
		}

		if algorithm != "" {
			if algorithm != "RS256" {
				return logical.ErrorResponse("unknown signing algorithm %q", algorithm), logical.ErrInvalidRequest
			}
			namedKey.SigningAlgorithm = algorithm
		}
	}

	// generate keys if creating a new key
	if !update {
		signingKey, publicKey, err := generateKeys(algorithm)
		if err != nil {
			return nil, err
		}

		// add public part of signing key to the key ring
		keyRing := []*ExpireableKey{
			{KeyID: publicKey.KeyID},
		}

		// create the named key
		namedKey = &NamedKey{
			Name:             name,
			SigningAlgorithm: algorithm,
			RotationPeriod:   rotationPeriod,
			VerificationTTL:  verificationTTL,
			KeyRing:          keyRing,
			SigningKey:       signingKey,
			NextRotation:     time.Now().Add(time.Duration(rotationPeriod) * time.Second),
		}

		if err := saveOIDCPublicKey(ctx, req.Storage, publicKey); err != nil {
			return nil, err
		}
	}

	// store named key (which was either just created or updated)
	entry, err = logical.StorageEntryJSON(namedKeyConfigPath+name, namedKey)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	i.oidcCache.purge()

	return &logical.Response{
		Data: map[string]interface{}{
			"rotation_period":  namedKey.RotationPeriod,
			"verification_ttl": namedKey.VerificationTTL,
			"algorithm":        namedKey.SigningAlgorithm,
		},
	}, nil
}

// handleOIDCReadKey is used to read an existing key
func (i *IdentityStore) pathOIDCReadKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		var storedNamedKey NamedKey
		if err := entry.DecodeJSON(&storedNamedKey); err != nil {
			return nil, err
		}
		return &logical.Response{
			Data: map[string]interface{}{
				"rotation_period":  storedNamedKey.RotationPeriod,
				"verification_ttl": storedNamedKey.VerificationTTL,
				"algorithm":        storedNamedKey.SigningAlgorithm,
			},
		}, nil
	}
	return logical.ErrorResponse("no named key found at %q", name), logical.ErrInvalidRequest
}

// handleOIDCDeleteKey is used to delete a key
func (i *IdentityStore) pathOIDCDeleteKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	targetKeyName := d.Get("name").(string)

	// it is an error to delete a key that is actively referenced by a role
	roleNames, err := req.Storage.List(ctx, roleConfigPath)
	if err != nil {
		return nil, err
	}

	var role *Role
	rolesReferencingTargetKeyName := make([]string, 0)
	for _, roleName := range roleNames {
		entry, err := req.Storage.Get(ctx, roleConfigPath+roleName)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&role); err != nil {
				return nil, err
			}
			if role.Key == targetKeyName {
				rolesReferencingTargetKeyName = append(rolesReferencingTargetKeyName, role.Name)
			}
		}
	}

	if len(rolesReferencingTargetKeyName) > 0 {
		errorMessage := fmt.Sprintf("Unable to delete the key: %q because it is currently referenced by these roles: %s",
			targetKeyName, strings.Join(rolesReferencingTargetKeyName, ", "))
		return logical.ErrorResponse(errorMessage), logical.ErrInvalidRequest
	}

	// key can safely be deleted now
	err = req.Storage.Delete(ctx, namedKeyConfigPath+targetKeyName)
	if err != nil {
		return nil, err
	}

	_, err = i.expireOIDCPublicKeys(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	i.oidcCache.purge()

	return nil, nil
}

// handleOIDCListKey is used to list named keys
func (i *IdentityStore) pathOIDCListKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	keys, err := req.Storage.List(ctx, namedKeyConfigPath)
	if err != nil {
		return nil, err
	}
	return logical.ListResponse(keys), nil
}

// pathOIDCRotateKey is used to manually trigger a rotation on the named key
func (i *IdentityStore) pathOIDCRotateKey(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)
	verificationTTLOverrideInput, verificationTTLOverrideInputOK := d.GetOk("verification_ttl")

	// load the named key and perform a rotation
	entry, err := req.Storage.Get(ctx, namedKeyConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("no named key found at %q", name), logical.ErrInvalidRequest
	}

	var storedNamedKey NamedKey
	if err := entry.DecodeJSON(&storedNamedKey); err != nil {
		return nil, err
	}

	// call rotate with an appropriate overrideTTL where -1 means no override
	var verificationTTLOverride int
	if verificationTTLOverrideInputOK {
		verificationTTLOverride = verificationTTLOverrideInput.(int)
	} else {
		verificationTTLOverride = -1
	}
	verificationTTLUsed, err := storedNamedKey.rotate(ctx, req.Storage, verificationTTLOverride)
	if err != nil {
		return nil, err
	}

	i.oidcCache.purge()

	// prepare response
	return &logical.Response{
		Data: map[string]interface{}{
			"rotation_period":  storedNamedKey.RotationPeriod,
			"verification_ttl": verificationTTLUsed,
			"algorithm":        storedNamedKey.SigningAlgorithm,
		},
	}, nil
}

func (i *IdentityStore) pathOIDCDiscovery(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data []byte

	i.oidcCache.lock.RLock()
	if len(i.oidcCache.discoveryResponse) > 0 {
		data = i.oidcCache.discoveryResponse
		i.oidcCache.lock.RUnlock()
	} else {
		i.oidcCache.lock.RUnlock()
		issuer := "http://127.0.0.1:8200"

		disc := discovery{
			Issuer: issuer + "/v1/identity/oidc",
			//Auth:        s.absURL("/auth"),
			//Token:       s.absURL("/token"),
			Keys:        issuer + "/v1/identity/oidc/.well-known/certs",
			Subjects:    []string{"public"},
			IDTokenAlgs: []string{string(jose.RS256)},
			//Scopes:      []string{"openid", "email", "groups", "profile", "offline_access"},
			//AuthMethods: []string{"client_secret_basic"},
			//Claims: []string{
			//	"aud", "email", "email_verified", "exp",
			//	"iat", "iss", "locale", "name", "sub",
			//},
		}

		var err error
		data, err = json.Marshal(disc)
		if err != nil {
			return nil, err
		}

		i.oidcCache.lock.Lock()
		i.oidcCache.discoveryResponse = data
		i.oidcCache.lock.Unlock()
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

// pathOIDCReadPublicKeys is used to retrieve all public keys so that clients can
// verify the validity of a signed OIDC token.
func (i *IdentityStore) pathOIDCReadPublicKeys(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	var data []byte

	i.oidcCache.lock.RLock()
	if len(i.oidcCache.jwksResponse) > 0 {
		data = i.oidcCache.jwksResponse
		i.oidcCache.lock.RUnlock()
	} else {
		i.oidcCache.lock.RUnlock()

		keyIDs, err := listOIDCPublicKeys(ctx, req.Storage)
		if err != nil {
			return nil, err
		}

		jwks := jose.JSONWebKeySet{
			Keys: make([]jose.JSONWebKey, 0, len(keyIDs)),
		}

		for _, keyID := range keyIDs {
			key, err := loadOIDCPublicKey(ctx, req.Storage, keyID)
			if err != nil {
				return nil, err
			}
			jwks.Keys = append(jwks.Keys, *key)
		}

		data, err = json.Marshal(jwks)
		if err != nil {
			return nil, err
		}

		i.oidcCache.lock.Lock()
		i.oidcCache.jwksResponse = data
		i.oidcCache.lock.Unlock()
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

// handleOIDCGenerateSignToken generates and signs an OIDC token
func (i *IdentityStore) pathOIDCGenerateToken(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	roleName := d.Get("name").(string)

	role, err := i.getOIDCRole(ctx, req.Storage, roleName)
	if err != nil {
		return nil, err
	}
	if role == nil {
		return logical.ErrorResponse("role %q not found", roleName), nil

	}

	entry, _ := req.Storage.Get(ctx, namedKeyConfigPath+role.Key)
	if entry == nil {
		return logical.ErrorResponse("key %q not found", role.Key), nil
	}

	var namedKey NamedKey
	if err := entry.DecodeJSON(&namedKey); err != nil {
		return nil, err
	}

	// generate an OIDC token from entity data
	accessorEntry, err := i.core.tokenStore.lookupByAccessor(ctx, req.ClientTokenAccessor, false, false)
	if err != nil {
		return nil, err
	}

	te, err := i.core.LookupToken(ctx, accessorEntry.TokenID)
	if err != nil {
		return nil, err
	}

	if te == nil || te.EntityID == "" {
		return logical.ErrorResponse("No entity associated with this Vault token"), nil
	}

	config, err := i.getOIDCConfig(ctx, req.Storage)
	if err != nil {
		return nil, err
	}

	issuer := i.core.redirectAddr
	if config != nil && config.Issuer != "" {
		issuer = config.Issuer
	}

	now := time.Now()
	idToken := idToken{
		Issuer:   issuer + "/v1/identity/oidc",
		Subject:  te.EntityID,
		Audience: role.RoleID,
		Expiry:   now.Add(role.TokenTTL).Unix(),
		IssuedAt: now.Unix(),
		Claims:   te,
	}

	e, err := i.MemDBEntityByID(te.EntityID, true)
	if err != nil {
		return nil, err
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

	signedIdToken, err := namedKey.signPayload(payload)
	if err != nil {
		return nil, errwrap.Wrapf("error signing OIDC token: {{err}}", err)
	}

	return &logical.Response{
		Data: map[string]interface{}{
			"token":     signedIdToken,
			"client_id": role.RoleID,
			"ttl":       int64(role.TokenTTL.Seconds()),
		},
	}, nil
}

func (key *NamedKey) signPayload(payload []byte) (string, error) {
	signingKey := jose.SigningKey{Key: key.SigningKey, Algorithm: jose.SignatureAlgorithm(key.SigningAlgorithm)}
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

// handleOIDCCreateRole is used to create a new role or update an existing one
func (i *IdentityStore) pathOIDCUpdateRole(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
	name := d.Get("name").(string)

	var role Role
	entry, err := req.Storage.Get(ctx, roleConfigPath+name)
	if err != nil {
		return nil, err
	}
	if entry != nil {
		if err := entry.DecodeJSON(&role); err != nil {
			return nil, err
		}
	}

	if key, ok := d.GetOk("key"); ok {
		role.Key = key.(string)
	} else if req.Operation == logical.CreateOperation {
		role.Key = d.Get("key").(string)
	}

	if role.Key == "" {
		return logical.ErrorResponse("key must be provided"), nil
	}

	// validate that key exists
	entry, err = req.Storage.Get(ctx, namedKeyConfigPath+role.Key)
	if err != nil {
		return nil, err
	}
	if entry == nil {
		return logical.ErrorResponse("Key %q does not exist", role.Key), nil
	}

	if template, ok := d.GetOk("template"); ok {
		role.Template = template.(string)
	} else if req.Operation == logical.CreateOperation {
		role.Template = d.Get("template").(string)
	}

	// Validate that template can be parsed and results in valid JSON
	if role.Template != "" {
		_, populatedTemplate, err := identity.PopulateString(&identity.PopulateStringInput{
			Mode:   identity.JSONTemplating,
			String: role.Template,
			Entity: new(identity.Entity),
			Groups: make([]*identity.Group, 0),
			// namespace?
		})

		if err != nil {
			return logical.ErrorResponse("Error parsing template: %s", err.Error()), nil
		}

		var tmp map[string]interface{}
		if err := json.Unmarshal([]byte(populatedTemplate), &tmp); err != nil {
			return logical.ErrorResponse("Error parsing template JSON: %s", err.Error()), nil
		}

		for key, _ := range tmp {
			if strutil.StrListContains(requiredClaims, key) {
				return logical.ErrorResponse(`Top level key %q not allowed. Restricted keys: %s`,
					key, strings.Join(requiredClaims, ", ")), nil
			}
		}
	}

	if ttl, ok := d.GetOk("ttl"); ok {
		role.TokenTTL = time.Duration(ttl.(int)) * time.Second
	} else if req.Operation == logical.CreateOperation {
		role.TokenTTL = time.Duration(d.Get("ttl").(int)) * time.Second
	}

	// create role path
	if role.RoleID == "" {
		roleID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}
		role.RoleID = roleID
	}

	role.Name = name // TODO: needed???

	// store role (which was either just created or updated)
	entry, err = logical.StorageEntryJSON(roleConfigPath+name, role)
	if err != nil {
		return nil, err
	}
	if err := req.Storage.Put(ctx, entry); err != nil {
		return nil, err
	}

	i.oidcCache.purge()

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
			"client_id": role.RoleID,
			"key":       role.Key,
			"template":  role.Template,
			"ttl":       int64(role.TokenTTL.Seconds()),
		},
	}, nil
}

func (i *IdentityStore) getOIDCRole(ctx context.Context, s logical.Storage, roleName string) (*Role, error) {
	entry, err := s.Get(ctx, roleConfigPath+roleName)
	if err != nil {
		return nil, err
	}

	if entry == nil {
		return nil, nil
	}

	var role Role
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

func (tok *idToken) generatePayload(logger hclog.Logger, template string, entity *identity.Entity, groups []*identity.Group) ([]byte, error) {
	output := map[string]interface{}{
		"iss": tok.Issuer,
		"sub": tok.Subject,
		"aud": tok.Audience,
		"exp": tok.Expiry,
		"iat": tok.IssuedAt,
	}

	// Parse and integrate the populated role template. Structural errors with the template _should_
	// be caught during role configuration. Error found during runtime will be logged, but they will
	// not block generation of the basic ID token. They should not be returned to the requester.
	_, populatedTemplate, err := identity.PopulateString(&identity.PopulateStringInput{
		Mode:   identity.JSONTemplating,
		String: template,
		Entity: entity,
		Groups: groups,
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

// --- some helper methods ---

// NamedKey.rotate(overrides) performs a key rotation on a NamedKey and returns the
// verification_ttl that was applied. verification_ttl can be overriden with an
// overrideVerificationTTL value >= 0
func (namedKey *NamedKey) rotate(ctx context.Context, s logical.Storage, overrideVerificationTTL int) (int, error) {
	verificationTTL := namedKey.VerificationTTL

	if overrideVerificationTTL >= 0 {
		verificationTTL = overrideVerificationTTL
	}

	// determine verificationTTL duration
	verificationTTLDuration := time.Duration(verificationTTL) * time.Second

	// generate new key
	signingKey, publicKey, err := generateKeys(namedKey.SigningAlgorithm)
	if err != nil {
		return -1, err
	}
	if err := saveOIDCPublicKey(ctx, s, publicKey); err != nil {
		return -1, err
	}

	// rotation involves overwriting current signing key, updating key ring, and updating global
	// public keys to expire the signing key that was just rotated
	rotateID := namedKey.SigningKey.KeyID
	namedKey.SigningKey = signingKey
	namedKey.KeyRing = append(namedKey.KeyRing, &ExpireableKey{KeyID: publicKey.KeyID})

	// set the previous public key's  expiry time
	for _, key := range namedKey.KeyRing {
		if key.KeyID == rotateID {
			key.ExpireAt = time.Now().Add(verificationTTLDuration)
			break
		}
	}

	namedKey.NextRotation = time.Now().Add(time.Duration(namedKey.RotationPeriod) * time.Second)

	// store named key (it was modified when rotate was called on it)
	entry, err := logical.StorageEntryJSON(namedKeyConfigPath+namedKey.Name, namedKey)
	if err != nil {
		return -1, err
	}
	if err := s.Put(ctx, entry); err != nil {
		return -1, err
	}

	return verificationTTL, nil
}

// generateKeys returns a signingKey and publicKey pair
func generateKeys(algorithm string) (signingKey jose.JSONWebKey, publicKey jose.JSONWebKey, err error) {
	// 2048 is recommended by RSA Laboratories as a MINIMUM post 2015, 3072 bits
	// is also seen in the wild, this could be configurable if need be
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return
	}

	signingKey = jose.JSONWebKey{
		Key:       key,
		KeyID:     id,
		Algorithm: algorithm,
		Use:       "sig",
	}

	publicKey = jose.JSONWebKey{
		Key:       &key.PublicKey,
		KeyID:     id,
		Algorithm: algorithm,
		Use:       "sig",
	}

	return
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

func (i *IdentityStore) expireOIDCPublicKeys(ctx context.Context, s logical.Storage) (time.Time, error) {
	var didUpdate bool

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

		var key NamedKey
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
				i.Logger().Error("error updating key", "key", key.Name, "error", err)
			}

			if err := s.Put(ctx, entry); err != nil {
				i.Logger().Error("error saving key", "key", key.Name, "error", err)

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
		i.oidcCache.purge()
	}

	return nextExpiration, nil
}

func (i *IdentityStore) oidcKeyRotation(ctx context.Context, s logical.Storage) (time.Time, error) {
	// soonestRotation will be the soonest rotation time of all keys. Initialize
	// here to a relatively distant time.
	soonestRotation := time.Now().Add(24 * time.Hour)
	now := time.Now()

	keys, err := s.List(ctx, namedKeyConfigPath)
	if err != nil {
		return now, err
	}

	for _, k := range keys {
		entry, err := s.Get(ctx, namedKeyConfigPath+k)
		if err != nil {
			return now, err
		}

		var key NamedKey
		if err := entry.DecodeJSON(&key); err != nil {
			return now, err
		}

		// Future key rotation that is the earliest we've seen.
		if now.Before(key.NextRotation) && key.NextRotation.Before(soonestRotation) {
			soonestRotation = key.NextRotation
		}

		// Key that is due to be rotated.
		if now.After(key.NextRotation) {
			i.Logger().Debug("rotating OIDC key", "key", key.Name)
			if _, err := key.rotate(ctx, s, -1); err != nil {
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
func (i *IdentityStore) oidcPeriodicFunc(ctx context.Context, s logical.Storage) {
	i.oidcCache.lock.RLock()
	nextRun := i.oidcCache.nextPeriodicRun
	i.oidcCache.lock.RUnlock()

	// The condition here is for performance, not precise timing. The actions can
	// be run at any time safely, but there is no need to invoke them (which
	// might be somewhat expensive if there are many roles/keys) if we're not
	// past any rotation/expiration TTLs.
	if time.Now().After(nextRun) {
		nextRotation, err := i.oidcKeyRotation(ctx, s)
		if err != nil {
			i.Logger().Warn("error rotating OIDC keys", "err", err)
		}

		nextExpiration, err := i.expireOIDCPublicKeys(ctx, s)
		if err != nil {
			i.Logger().Warn("error expiring OIDC public keys", "err", err)
		}

		i.oidcCache.purge()

		// re-run at the soonest expiration or rotation time
		i.oidcCache.lock.Lock()
		if nextRotation.Before(nextExpiration) {
			i.oidcCache.nextPeriodicRun = nextRotation
		} else {
			i.oidcCache.nextPeriodicRun = nextExpiration
		}
		i.oidcCache.lock.Unlock()
	}
}

type idToken struct {
	// ---- OIDC CLAIMS WITH NOTES FROM SPEC ----
	// required fields
	Issuer   string `json:"iss"` // Vault server address?
	Subject  string `json:"sub"`
	Audience string `json:"aud"`
	// Audience _should_ contain the OAuth 2.0 client_id of the Relying Party.
	// Not sure how/if we will leverage this

	Expiry   int64 `json:"exp"`
	IssuedAt int64 `json:"iat"`

	AuthTime int64 `json:"auth_time"`
	// required if max_age is specified in the Authentication request (which we aren't doing) or auth_time is identified by the client as an "Essential Claim"
	// we could return the time that the token was created at

	// optional fields

	// Nonce                               string `json:"nonce,omitempty"`
	// I don't think that Nonce will apply because we will not have any concept of "Authentication request".
	// From spec:
	// String value used to associate a Client session with an ID Token, and to mitigate replay attacks. The value is passed through unmodified from the Authentication Request to the ID Token. If present in the ID Token, Clients MUST verify that the nonce Claim Value is equal to the value of the nonce parameter sent in the Authentication Request. If present in the Authentication Request, Authorization Servers MUST include a nonce Claim in the ID Token with the Claim Value being the nonce value sent in the Authentication Request. Authorization Servers SHOULD perform no other processing on nonce values used. The nonce value is a case sensitive string.
	// where Authentication Request means:
	// OAuth 2.0 Authorization Request using extension parameters and scopes defined by OpenID Connect to request that the End-User be authenticated by the Authorization Server, which is an OpenID Connect Provider, to the Client, which is an OpenID Connect Relying Party.

	AuthenticationContextClassReference string `json:"acr,omitempty"`
	// Optional, very up to the implementation to decide on details.
	// from the spec:
	// Parties using this claim will need to agree upon the meanings of the values used, which may be context-specific.

	// maybe userpass auth is a lower level than approle or userpass ent with mfa enabled...
	// here is one spec...
	/*
		+--------------------------+---------+---------+---------+---------+
		| Token Type               | Level 1 | Level 2 | Level 3 | Level 4 |
		+--------------------------+---------+---------+---------+---------+
		| Hard crypto token        | X       | X       | X       | X       |
		|                          |         |         |         |         |
		| One-time password device | X       | X       | X       |         |
		|                          |         |         |         |         |
		| Soft crypto token        | X       | X       | X       |         |
		|                          |         |         |         |         |
		| Passwords & PINs         | X       | X       |         |         |
		+--------------------------+---------+---------+---------+---------+

		 +------------------------+---------+---------+---------+---------+
		 | Protect Against        | Level 1 | Level 2 | Level 3 | Level 4 |
		 +------------------------+---------+---------+---------+---------+
		 | On-line guessing       | X       | X       | X       | X       |
		 |                        |         |         |         |         |
		 | Replay                 | X       | X       | X       | X       |
		 |                        |         |         |         |         |
		 | Eavesdropper           |         | X       | X       | X       |
		 |                        |         |         |         |         |
		 | Verifier impersonation |         |         | X       | X       |
		 |                        |         |         |         |         |
		 | Man-in-the-middle      |         |         | X       | X       |
		 |                        |         |         |         |         |
		 | Session hijacking      |         |         |         | X       |
		 +------------------------+---------+---------+---------+---------+
	*/
	AuthenticationMethodsReference string `json:"amr,omitempty"`
	// I think this is only useful if downstream services will be making decisions based on what auth method was used to acquire a Vault token
	// which is something that we are trying to abstract away in using entityID as our sub. Think we can remove this.

	AuthorizingParty string `json:"azp,omitempty"`
	// I don't think we should use this for same, reasoning builds on not leveraging "aud" - checkout: thhttps://bitbucket.org/openid/connect/issues/973/

	// AccessTokenHash string `json:"at_hash,omitempty"`
	// I don't think that at_hash will apply because we are not creating any kind of access token (maybe the Vault Token is like an access token but how it was acquired is different from a typical oauth access token)
	// From the spec:
	// The contents of the ID Token are as described in Section 2. When using the Authorization Code Flow, these additional requirements for the following ID Token Claims apply:
	// at_hash
	// OPTIONAL. Access Token hash value. Its value is the base64url encoding of the left-most half of the hash of the octets of the ASCII representation of the access_token value, where the hash algorithm used is the hash algorithm used in the alg Header Parameter of the ID Token's JOSE Header. For instance, if the alg is RS256, hash the access_token value with SHA-256, then take the left-most 128 bits and base64url encode them. The at_hash value is a case sensitive string.

	// Email         string `json:"email,omitempty"`
	// EmailVerified *bool  `json:"email_verified,omitempty"`
	// Groups []string `json:"groups,omitempty"`
	// Name   string      `json:"name,omitempty"`
	Claims interface{} `json:"claims",omitempty`
	//FederatedIDClaims *federatedIDClaims `json:"federated_claims,omitempty"`
}

/*
// oidcPaths returns the API endpoints supported to operate on OIDC tokens:
// oidc/token - To register generate a new odic token
// oidc/key/:key - Create a new keyring
		oidcPathConfig(i),

//handleOIDCGenerate is used to generate an OIDC token
func (i *IdentityStore) handleOIDCGenerateIDToken() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		// Get entity linked to the requesting token
		// te could be nil if it is a root token
		// te could not have an entity if it was created from the token backend

		accessorEntry, err := i.core.tokenStore.lookupByAccessor(ctx, req.ClientTokenAccessor, false, false)
		if err != nil {
			return nil, err
		}

		te, err := i.core.LookupToken(ctx, accessorEntry.TokenID)
		if te == nil {
			return nil, errors.New("No token entry for this token")
		}
		fmt.Printf("-- -- --\nreq:\n%#v\n", req)
		fmt.Printf("-- -- --\nte:\n%#v\n", te)
		if err != nil {
			return nil, err
		}
		if te.EntityID == "" {
			return nil, errors.New("No EntityID associated with this request's Vault token")
		}

		now := time.Now()
		idToken := idToken{
			Issuer:   "Issuer",
			Subject:  te.EntityID,
			Audience: []string{"client_id_of_relying_party"},
			Expiry:   now.Add(2 * time.Minute).Unix(),
			IssuedAt: now.Unix(),
			Claims:   te,
		}

		// signing
		keyRing, _ := i.createKeyRing("foo")
		privWebKey, pubWebKey := keyRing.GenerateWebKeys()
		signedIdToken, _ := keyRing.SignIdToken(privWebKey, idToken)

		jwks := jose.JSONWebKeySet{
			Keys: make([]jose.JSONWebKey, 1),
		}
		jwks.Keys[0] = *pubWebKey

		return &logical.Response{
			Data: map[string]interface{}{
				"token": signedIdToken,
				"pub":   jwks,
			},
		}, nil
	}
}

func signPayload(key *jose.JSONWebKey, alg jose.SignatureAlgorithm, payload []byte) (jws string, err error) {
	signingKey := jose.SigningKey{Key: key, Algorithm: alg}

	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", fmt.Errorf("new signier: %v", err)
	}
	signature, err := signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("signing payload: %v", err)
	}
	return signature.CompactSerialize()
}


// --- --- KEY SIGNING FUNCTIONALITY --- ---

// type keyRing struct {
// 	mostRecentKeyAt int // locates the most recent key within keys, -1 means that there is no key in the key ring
// 	name            string
// 	numberOfKeys    int
// 	keyTTL          time.Duration
// 	keys            []keyRingKey
// }

type keyRingKey struct {
	createdAt time.Time
	key       *rsa.PrivateKey
	id        string
}

// TODO
// - USE A REAL CONFIG
// - STORE AND CACHE (look at upsertEntity)
// - LOCKS AROUND ROTATING

// populates an empty keyring from defaults or config
func (i *IdentityStore) emptyKeyRing() *keyRing {
	// retrieve config values if they exist
	numberOfKeys := 4
	keyTTL := 6 * time.Hour
	return &keyRing{
		mostRecentKeyAt: -1,
		numberOfKeys:    numberOfKeys,
		keyTTL:          keyTTL,
		keys:            make([]keyRingKey, numberOfKeys, numberOfKeys),
	}
}

// Create a keyRing
func (i *IdentityStore) createKeyRing(name string) (*keyRing, error) {
	// err if name already exist
	// retrieve configurations - hardcoded for now
	kr := i.emptyKeyRing()
	kr.name = name
	kr.Rotate()
	// store keyring
	return kr, nil
}

// RotateIfRequired performs a rotate if the current key is outdated
func (kr *keyRing) RotateIfRequired() error {
	expireAt := kr.keys[kr.mostRecentKeyAt].createdAt.Add(kr.keyTTL)
	now := time.Now().UTC().Round(time.Millisecond)
	if now.After(expireAt) {
		err := kr.Rotate()
		if err != nil {
			return err
		}
	}
	return nil
}

// Rotate adds a new key to a keyRing which may override existing entries
func (kr *keyRing) Rotate() error {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	var insertInKeyRingAt int
	switch {
	case kr.mostRecentKeyAt < 0:
		insertInKeyRingAt = 0
	case kr.mostRecentKeyAt >= 0:
		insertInKeyRingAt = (kr.mostRecentKeyAt + 1) % len(kr.keys)
	}

	kr.keys[insertInKeyRingAt].key = key
	kr.keys[insertInKeyRingAt].createdAt = time.Now().UTC().Round(time.Millisecond)
	kr.keys[insertInKeyRingAt].id = id
	kr.mostRecentKeyAt = insertInKeyRingAt
	return nil
}

//Sign a payload with a keyRing
func (kr *keyRing) SignIdToken(webKey *jose.JSONWebKey, token idToken) (string, error) {
	err := kr.RotateIfRequired()
	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(token)
	if err != nil {
		return "", err
	}

	signingKey := jose.SigningKey{Key: webKey, Algorithm: jose.SignatureAlgorithm(webKey.Algorithm)}
	signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
	if err != nil {
		return "", fmt.Errorf("new signier: %v", err)
	}
	signature, err := signer.Sign(payload)
	if err != nil {
		return "", fmt.Errorf("signing payload: %v", err)
	}
	return signature.CompactSerialize()
}

func (kr *keyRing) GenerateWebKeys() (priv *jose.JSONWebKey, pub *jose.JSONWebKey) {
	kr.RotateIfRequired()
	keyRingKey := kr.keys[kr.mostRecentKeyAt]

	priv = &jose.JSONWebKey{
		Key:       keyRingKey.key,
		KeyID:     keyRingKey.id,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	pub = &jose.JSONWebKey{
		Key:       keyRingKey.key.Public(),
		KeyID:     keyRingKey.id,
		Algorithm: string(jose.RS256),
		Use:       "sig",
	}
	return
}
*/
