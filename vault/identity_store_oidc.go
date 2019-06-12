package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"

	uuid "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/helper/parseutil"
	"github.com/hashicorp/vault/sdk/logical"
)

// types
type audience []string
type signingAlgorithm int

const (
	rs256 signingAlgorithm = iota
)

// globals - todo fix this
var publicKeys []*ExpireableKey = make([]*ExpireableKey, 0, 0)

type ExpireableKey struct {
	Key        jose.JSONWebKey `json:"key"`
	Expireable bool            `json:"expireable"`
	ExpireAt   time.Time       `json:"expire_at"`
}

type NamedKey struct {
	Name             string           `json:"name"`
	SigningAlgorithm signingAlgorithm `json:"signing_algorithm"`
	Verificationttl  string           `json:"verification_ttl"`
	RotationPeriod   string           `json:"rotation_period"`
	KeyRing          []string         `json:"key_ring"`
	SigningKey       jose.JSONWebKey  `json:"signing_key"`
}

// oidcPaths returns the API endpoints supported to operate on OIDC tokens:
// oidc/key/:key - Create a new key named key
func oidcPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		// {
		// 	Pattern: "oidc/token",
		// 	Callbacks: map[logical.Operation]framework.OperationFunc{
		// 		logical.UpdateOperation: i.handleOIDCGenerateIDToken(),
		// 	},

		// 	HelpSynopsis:    "HelpSynopsis here",
		// 	HelpDescription: "HelpDecription here",
		// },

		// {
		{
			Pattern: "oidc/key/" + framework.GenericNameRegex("name"),
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the key",
				},

				"rotation_period": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "How often to generate a new keypair. Defaults to 6h",
				},

				"verification_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Controls how long the public portion of a key will be available for verification after being rotated. Defaults to the current rotation_period, which will provide for a current key and previous key.",
				},

				"algorithm": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Signing algorithm to use. This will default to RS256, and is currently the only allowed value.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleOIDCCreateKey(),
				logical.ReadOperation:   i.handleOIDCReadKey(),
				logical.DeleteOperation: i.handleOIDCDeleteKey(),
			},
			HelpSynopsis:    "oidc/key/:key help synopsis here",
			HelpDescription: "oidc/key/:key help description here",
		},
		{
			Pattern: "oidc/key/" + framework.GenericNameRegex("name") + "/rotate/?$",
			Fields: map[string]*framework.FieldSchema{
				"name": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Name of the key",
				},
				"verification_ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Controls how long the public portion of a key will be available for verification after being rotated. Setting verification_ttl here will override the verification_ttl set on the key.",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleOIDCRotateKey(),
			},
			HelpSynopsis:    "oidc/key/:key/rotate help synopsis here",
			HelpDescription: "oidc/key/:key/rotate help description here",
		},
		{
			Pattern: "oidc/key/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ListOperation: i.handleOIDCListKeys(),
			},
			HelpSynopsis:    "list oidc/key/ help synopsis here",
			HelpDescription: "list oidc/key/ help description here",
		},
		{
			Pattern: "oidc/well-known/certs/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.handleOIDCReadCerts(),
			},
			HelpSynopsis:    "read oidc/well-known/certs/ help synopsis here",
			HelpDescription: "read oidc/well-known/certs/ help description here",
		},
		{
			Pattern: "oidc/token/generate/?$",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.ReadOperation: i.handleOIDCGenerateSignToken(),
			},
		},
	}
}

// handleOIDCCreateKey is used to create a new key
func (i *IdentityStore) handleOIDCCreateKey() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		var namedKey *NamedKey
		var update bool = false
		var algorithm signingAlgorithm

		// parse parameters
		name := d.Get("name").(string)
		rotationPeriod := d.Get("rotation_period").(string)
		verificationttl := d.Get("verification_ttl").(string)
		algorithmInput := d.Get("algorithm").(string)

		// determine if we are creating a key or updating an existing key
		entry, err := req.Storage.Get(ctx, "oidc-config/named-key/"+name)
		if err != nil {
			return nil, err
		}
		if entry != nil {
			if err := entry.DecodeJSON(&namedKey); err != nil {
				return nil, err
			}
			update = true
		}

		// set defaults if creating a new key
		if !update {
			if rotationPeriod == "" {
				rotationPeriod = "6h"
			}
			if verificationttl == "" {
				verificationttl = rotationPeriod
			}
			if algorithmInput == "" {
				algorithmInput = "RS256"
			}
		}

		// verify inputs and set values if this is an update
		if rotationPeriod != "" {
			_, err = parseutil.ParseDurationSecond(rotationPeriod)
			if err != nil {
				return nil, fmt.Errorf("unable to parse provided rotation_period of: %s", rotationPeriod)
			}
			if update {
				namedKey.RotationPeriod = rotationPeriod
			}
		}

		if verificationttl != "" {
			_, err = parseutil.ParseDurationSecond(verificationttl)
			if err != nil {
				return nil, fmt.Errorf("unable to parse provided verification_ttl of: %s", verificationttl)
			}
			if update {
				namedKey.Verificationttl = verificationttl
			}
		}

		if algorithmInput != "" {
			switch algorithmInput {
			case "RS256":
				algorithm = rs256
			default:
				return logical.ErrorResponse(fmt.Sprintf("unknown signing algorithm %q", algorithmInput)), logical.ErrInvalidRequest
			}
			if update {
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
			keyRing := make([]string, 1, 1)
			keyRing[0] = publicKey.Key.KeyID

			// create the named key
			namedKey = &NamedKey{
				Name:             name,
				SigningAlgorithm: algorithm,
				RotationPeriod:   rotationPeriod,
				Verificationttl:  verificationttl,
				KeyRing:          keyRing,
				SigningKey:       signingKey,
			}

			// add the public key to the struct containing all public keys and store it
			publicKeys = append(publicKeys, &publicKey)
			entry, err = logical.StorageEntryJSON("oidc-config/public-keys/", publicKeys)
			if err != nil {
				return nil, err
			}
			if err := req.Storage.Put(ctx, entry); err != nil {
				return nil, err
			}
		}

		// store named key (which was either just created or updated)
		entry, err = logical.StorageEntryJSON("oidc-config/named-key/"+name, namedKey)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"rotation_period":  namedKey.RotationPeriod,
				"verification_ttl": namedKey.Verificationttl,
				"algorithm":        namedKey.SigningAlgorithm.String(),
			},
		}, nil
	}
}

// handleOIDCReadKey is used to read an existing key
func (i *IdentityStore) handleOIDCReadKey() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		name := d.Get("name").(string)

		entry, err := req.Storage.Get(ctx, "oidc-config/named-key/"+name)
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
					"verification_ttl": storedNamedKey.Verificationttl,
					"algorithm":        storedNamedKey.SigningAlgorithm.String(),
				},
			}, nil
		}
		return logical.ErrorResponse(fmt.Sprintf("no named key was stored at %q", name)), logical.ErrInvalidRequest
	}
}

// handleOIDCDeleteKey is used to delete a key
func (i *IdentityStore) handleOIDCDeleteKey() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		name := d.Get("name").(string)

		err := req.Storage.Delete(ctx, "oidc-config/named-key/"+name)
		if err != nil {
			return nil, err
		}
		//todo this is supposed to return  204-no content
		return nil, nil
	}
}

// handleOIDCListKey is used to list named keys
func (i *IdentityStore) handleOIDCListKeys() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		keys, err := req.Storage.List(ctx, "oidc-config/named-key/")
		if err != nil {
			return nil, err
		}
		return logical.ListResponse(keys), nil
	}
}

// handleOIDCRotateKey is used to manually trigger a rotation on the named key
func (i *IdentityStore) handleOIDCRotateKey() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		name := d.Get("name").(string)
		verificationttlInput := d.Get("verification_ttl").(string)
		var responseVerificationttl string

		// validate inputs and update rotation overrides
		overrides := make(map[string]interface{})
		if verificationttlInput != "" {
			verificationttlDuration, err := parseutil.ParseDurationSecond(verificationttlInput)
			if err != nil {
				return nil, fmt.Errorf("unable to parse provided verification_ttl of: %s", verificationttlInput)
			}
			overrides["verificationttlDuration"] = verificationttlDuration
		}

		// load the named key and perform a rotation
		entry, err := req.Storage.Get(ctx, "oidc-config/named-key/"+name)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			return logical.ErrorResponse(fmt.Sprintf("no named key found at %q", name)), logical.ErrInvalidRequest
		}

		var storedNamedKey NamedKey
		if err := entry.DecodeJSON(&storedNamedKey); err != nil {
			return nil, err
		}
		err = storedNamedKey.rotate(overrides)
		if err != nil {
			return nil, err
		}

		// store named key (it was modified when rotate was called on it)
		entry, err = logical.StorageEntryJSON("oidc-config/named-key/"+name, storedNamedKey)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		// prepare response
		if verificationttlInput != "" {
			responseVerificationttl = verificationttlInput
		} else {
			responseVerificationttl = storedNamedKey.Verificationttl
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"rotation_period":  storedNamedKey.RotationPeriod,
				"verification_ttl": responseVerificationttl,
				"algorithm":        storedNamedKey.SigningAlgorithm.String(),
			},
		}, nil
	}
}

// handleOIDCReadCerts is used to retrieve all certs from all keys so that clients can
// verify the validity of a signed OIDC token
func (i *IdentityStore) handleOIDCReadCerts() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		// fmt.Printf("\npublicKeys before expiry:\n%#v", publicKeys)
		publicKeys = expire(publicKeys)
		// fmt.Printf("\npublicKeys after expiry:\n%#v", publicKeys)

		jwks := jose.JSONWebKeySet{
			Keys: make([]jose.JSONWebKey, len(publicKeys)),
		}

		for i, expireableKey := range publicKeys {
			jwks.Keys[i] = expireableKey.Key
		}

		return &logical.Response{
			Data: map[string]interface{}{
				"keys": jwks,
			},
		}, nil
	}
}

// handleOIDCGenerateSignToken generates and signs an OIDC token
func (i *IdentityStore) handleOIDCGenerateSignToken() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {
		// load test key yolo style
		entry, _ := req.Storage.Get(ctx, "oidc-config/named-key/test")
		var namedKey NamedKey
		entry.DecodeJSON(&namedKey)

		// generate an OIDC token from entity data
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
		// keyRing, _ := i.createKeyRing("foo")
		// privWebKey, pubWebKey := keyRing.GenerateWebKeys()
		// signedIdToken, _ := keyRing.SignIdToken(privWebKey, idToken)
		payload, err := json.Marshal(idToken)
		if err != nil {
			return nil, err
		}

		signingKey := jose.SigningKey{Key: namedKey.SigningKey, Algorithm: jose.SignatureAlgorithm("RS256")}
		signer, err := jose.NewSigner(signingKey, &jose.SignerOptions{})
		if err != nil {
			return nil, fmt.Errorf("new signier: %v", err)
		}
		signature, err := signer.Sign(payload)
		if err != nil {
			return nil, fmt.Errorf("signing payload: %v", err)
		}
		signedIdToken, _ := signature.CompactSerialize()

		jwks := jose.JSONWebKeySet{
			Keys: make([]jose.JSONWebKey, 1),
		}
		jwks.Keys[0] = namedKey.SigningKey.Public()

		return &logical.Response{
			Data: map[string]interface{}{
				"token": signedIdToken,
				"pub":   jwks,
			},
		}, nil
	}
}

// --- some helper methods ---

// expire returns a lice of ExpireableKey with the expired keys removed
func expire(keys []*ExpireableKey) []*ExpireableKey {
	activeKeys := make([]*ExpireableKey, len(keys))
	now := time.Now()

	// a key is active if it is not yet expireable (because it is the signing key of some named
	// keyring) or if it is expireable but the expireAt time is in the future
	insertIndex := 0
	fmt.Printf("\nkeys:\n%#v", keys)
	for i := range keys {
		fmt.Printf("\nkey:\n%#v", keys[i])
		switch {
		// expire case
		case keys[i].Expireable && now.After(keys[i].ExpireAt):
		default:
			activeKeys[insertIndex] = keys[i]
			insertIndex = insertIndex + 1
		}
	}

	return activeKeys[:insertIndex]
}

// NamedKey.rotate(overrides) performs a key rotation on a NamedKey
// ASSUMPTION: overrides make sense and won't cause type conversion errors...
func (namedKey *NamedKey) rotate(overrides map[string]interface{}) error {
	// fail early if the verification ttl of the named key is not parse-able
	var verificationttlDuration time.Duration
	verificationttlDurationInterface, ok := overrides["verificationttlDuration"]
	if ok {
		verificationttlDuration = verificationttlDurationInterface.(time.Duration)
	} else {
		var err error
		verificationttlDuration, err = parseutil.ParseDurationSecond(namedKey.Verificationttl)
		if err != nil {
			return err
		}
	}

	// generate new key
	signingKey, publicKey, err := generateKeys(namedKey.SigningAlgorithm)
	if err != nil {
		return err
	}

	// rotation involves overwritting current signing key, updating key ring, and updating global
	// public keys to expire the signing key that was just rotated
	rotateID := namedKey.SigningKey.KeyID
	namedKey.SigningKey = signingKey
	publicKeys = append(publicKeys, &publicKey)
	namedKey.KeyRing = append(namedKey.KeyRing, publicKey.Key.KeyID)

	// give current signing key's public portion an expiry time
	fmt.Printf("\n--- looping through public keys")
	for i := range publicKeys {
		fmt.Printf("\n--- current key[%d]:\n%#v", i, publicKeys[i])
		if publicKeys[i].Key.KeyID == rotateID {
			fmt.Printf("\n--- rotateID match found")
			publicKeys[i].Expireable = true
			publicKeys[i].ExpireAt = time.Now().Add(verificationttlDuration)
			fmt.Printf("\n--- current key[%d]:\n%#v", i, publicKeys[i])
			break
		}
	}

	return nil
}

// generateKeys returns a signingKey and publicKey pair
func generateKeys(algorithm signingAlgorithm) (signingKey jose.JSONWebKey, publicKey ExpireableKey, err error) {
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
		Algorithm: algorithm.String(),
		Use:       "sig",
	}

	publicKey = ExpireableKey{
		Key: jose.JSONWebKey{
			Key:       &key.PublicKey,
			KeyID:     id,
			Algorithm: algorithm.String(),
			Use:       "sig",
		},
		Expireable: false,
		ExpireAt:   time.Time{},
	}
	return
}

// signingAlgorithm.String takes a signingAlgorithm and returns the string representation of that algorithm
func (a signingAlgorithm) String() string {
	switch a {
	case rs256:
		return "RS256"
	default:
		return "unknown"
	}
}

type idToken struct {
	// ---- OIDC CLAIMS WITH NOTES FROM SPEC ----
	// required fields
	Issuer   string   `json:"iss"` // Vault server address?
	Subject  string   `json:"sub"`
	Audience audience `json:"aud"`
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
