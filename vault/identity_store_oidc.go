package vault

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hashicorp/go-uuid"
	"gopkg.in/square/go-jose.v2"

	"github.com/hashicorp/vault/sdk/framework"
	"github.com/hashicorp/vault/sdk/logical"
)

type audience []string

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
	/* From [NIST_SP800-63] .

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

// oidcPaths returns the API endpoints supported to operate on OIDC tokens:
// oidc/token - To register generate a new odic token
// oidc/??? -
func oidcPaths(i *IdentityStore) []*framework.Path {
	return []*framework.Path{
		{
			Pattern: "oidc/token",
			Callbacks: map[logical.Operation]framework.OperationFunc{
				logical.UpdateOperation: i.handleOIDCGenerateIDToken(),
			},

			HelpSynopsis:    "HelpSynopsis here",
			HelpDescription: "HelpDecription here",
		},
	}
}

// handleOIDCGenerate is used to generate an OIDC token
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
		idToken := &idToken{
			Issuer:   "Issuer",
			Subject:  te.EntityID,
			Audience: []string{"client_id_of_relying_party"},
			Expiry:   now.Add(2 * time.Minute).Unix(),
			IssuedAt: now.Unix(),
			Claims:   te,
		}

		// signing
		key, _ := rsa.GenerateKey(rand.Reader, 2048)
		keyID, err := uuid.GenerateUUID()
		if err != nil {
			return nil, err
		}

		priv := &jose.JSONWebKey{
			Key:       key,
			KeyID:     keyID,
			Algorithm: "RS256",
			Use:       "sig",
		}
		pub := &jose.JSONWebKey{
			Key:       key.Public(),
			KeyID:     keyID,
			Algorithm: "RS256",
			Use:       "sig",
		}

		payload, err := json.Marshal(idToken)
		signedIDToken, err := signPayload(priv, jose.RS256, payload)

		jwks := jose.JSONWebKeySet{
			Keys: make([]jose.JSONWebKey, 1),
		}
		jwks.Keys[0] = *pub

		//data2, err := json.MarshalIndent(jwks, "", "  ")

		return &logical.Response{
			Data: map[string]interface{}{
				"token": signedIDToken,
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
