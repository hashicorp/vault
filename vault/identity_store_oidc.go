package vault

import (
	"context"
	"fmt"
	"time"

	"gopkg.in/square/go-jose.v2"

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

type key struct {
	name             string
	signingAlgorithm signingAlgorithm
	ttl              time.Duration
	rotationPeriod   time.Duration
	keyRing          jose.JSONWebKeySet
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
					Default:     "6h",
				},

				"ttl": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "TTL of generated keys. After expiration they will be deleted from the keyring and tokens generated with them will be invalid. Defaults to 24h.",
					Default:     "24h",
				},

				"algorithm": &framework.FieldSchema{
					Type:        framework.TypeString,
					Description: "Signing algorithm to use. This will default to RS256, and is currently the only allowed value.",
					Default:     "RS256",
				},
			},
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.CreateOperation: i.handleOIDCCreateKey(),
			},

			HelpSynopsis:    "oidc/key/:key help synopsis here",
			HelpDescription: "oidc/key/:key help description here",
		},
	}
}

//handleOIDCCreateKey is used to create a new key
func (i *IdentityStore) handleOIDCCreateKey() framework.OperationFunc {
	return func(ctx context.Context, req *logical.Request, d *framework.FieldData) (*logical.Response, error) {

		// parse parameters
		nameInput := d.Get("name").(string)
		rotationPeriodInput := d.Get("rotation_period").(string)
		ttlInput := d.Get("ttl").(string)
		algorithmInput := d.Get("algorithm").(string)
		// check inputs
		// check that this named key doesn't already exist
		rotationPeriod, err := parseutil.ParseDurationSecond(rotationPeriodInput)
		if err != nil {
			return nil, fmt.Errorf("unable to parse provided rotation_period of: %s", rotationPeriodInput)
		}

		ttl, err := parseutil.ParseDurationSecond(ttlInput)
		if err != nil {
			return nil, fmt.Errorf("unable to parse provided ttl of: %s", ttlInput)
		}

		var algorithm signingAlgorithm
		switch algorithmInput {
		case "RS256":
			algorithm = rs256
		default:
			return logical.ErrorResponse(fmt.Sprintf("unknown signing algorithm %q", algorithmInput)), logical.ErrInvalidRequest
		}

		// create key
		key := &key{
			name:             nameInput,
			signingAlgorithm: algorithm,
			rotationPeriod:   rotationPeriod,
			ttl:              ttl,
			keyRing:          jose.JSONWebKeySet{make([]jose.JSONWebKey, 0)},
		}

		// store key
		entry, err := logical.StorageEntryJSON("config/key/"+name, key)
		if err != nil {
			return nil, err
		}
		if err := req.Storage.Put(ctx, entry); err != nil {
			return nil, err
		}

		return nil, nil
	}
}

/*
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
