package vault

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
	"github.com/SermoDigital/jose/jwt"
	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/vault/helper/consts"
	"github.com/hashicorp/vault/helper/jsonutil"
	"github.com/hashicorp/vault/logical"
	jose "gopkg.in/square/go-jose.v2"
	squarejwt "gopkg.in/square/go-jose.v2/jwt"
)

const (
	// The location of the key used to generate response-wrapping JWTs
	coreWrappingJWTKeyPath = "core/wrapping/jwtkey"
)

func (c *Core) ensureWrappingKey(ctx context.Context) error {
	entry, err := c.barrier.Get(ctx, coreWrappingJWTKeyPath)
	if err != nil {
		return err
	}

	var keyParams clusterKeyParams

	if entry == nil {
		key, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
		if err != nil {
			return errwrap.Wrapf("failed to generate wrapping key: {{err}}", err)
		}
		keyParams.D = key.D
		keyParams.X = key.X
		keyParams.Y = key.Y
		keyParams.Type = corePrivateKeyTypeP521
		val, err := jsonutil.EncodeJSON(keyParams)
		if err != nil {
			return errwrap.Wrapf("failed to encode wrapping key: {{err}}", err)
		}
		entry = &Entry{
			Key:   coreWrappingJWTKeyPath,
			Value: val,
		}
		if err = c.barrier.Put(ctx, entry); err != nil {
			return errwrap.Wrapf("failed to store wrapping key: {{err}}", err)
		}
	}

	// Redundant if we just created it, but in this case serves as a check anyways
	if err = jsonutil.DecodeJSON(entry.Value, &keyParams); err != nil {
		return errwrap.Wrapf("failed to decode wrapping key parameters: {{err}}", err)
	}

	c.wrappingJWTKey = &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P521(),
			X:     keyParams.X,
			Y:     keyParams.Y,
		},
		D: keyParams.D,
	}

	c.logger.Info("loaded wrapping token key")

	return nil
}

func (c *Core) wrapInCubbyhole(ctx context.Context, req *logical.Request, resp *logical.Response, auth *logical.Auth) (*logical.Response, error) {
	// Before wrapping, obey special rules for listing: if no entries are
	// found, 404. This prevents unwrapping only to find empty data.
	if req.Operation == logical.ListOperation {
		if resp == nil || (len(resp.Data) == 0 && len(resp.Warnings) == 0) {
			return nil, logical.ErrUnsupportedPath
		}

		keysRaw, ok := resp.Data["keys"]
		if !ok || keysRaw == nil {
			if len(resp.Data) > 0 || len(resp.Warnings) > 0 {
				// We could be returning extra metadata on a list, or returning
				// warnings with no data, so handle these cases
				goto DONELISTHANDLING
			}
			return nil, logical.ErrUnsupportedPath
		}

		keys, ok := keysRaw.([]string)
		if !ok {
			return nil, logical.ErrUnsupportedPath
		}
		if len(keys) == 0 {
			return nil, logical.ErrUnsupportedPath
		}
	}

DONELISTHANDLING:
	var err error
	sealWrap := resp.WrapInfo.SealWrap

	// If we are wrapping, the first part (performed in this functions) happens
	// before auditing so that resp.WrapInfo.Token can contain the HMAC'd
	// wrapping token ID in the audit logs, so that it can be determined from
	// the audit logs whether the token was ever actually used.
	creationTime := time.Now()
	te := logical.TokenEntry{
		Path:           req.Path,
		Policies:       []string{"response-wrapping"},
		CreationTime:   creationTime.Unix(),
		TTL:            resp.WrapInfo.TTL,
		NumUses:        1,
		ExplicitMaxTTL: resp.WrapInfo.TTL,
	}

	if err := c.tokenStore.create(ctx, &te); err != nil {
		c.logger.Error("failed to create wrapping token", "error", err)
		return nil, ErrInternalError
	}

	resp.WrapInfo.Token = te.ID
	resp.WrapInfo.Accessor = te.Accessor
	resp.WrapInfo.CreationTime = creationTime
	// If this is not a rewrap, store the request path as creation_path
	if req.Path != "sys/wrapping/rewrap" {
		resp.WrapInfo.CreationPath = req.Path
	}

	if auth != nil && auth.EntityID != "" {
		resp.WrapInfo.WrappedEntityID = auth.EntityID
	}

	// This will only be non-nil if this response contains a token, so in that
	// case put the accessor in the wrap info.
	if resp.Auth != nil {
		resp.WrapInfo.WrappedAccessor = resp.Auth.Accessor
	}

	switch resp.WrapInfo.Format {
	case "jwt":
		// Create the JWT
		claims := squarejwt.Claims{
			// Map the JWT ID to the token ID for ease of use
			ID: te.ID,
			// Set the issue time to the creation time
			IssuedAt: squarejwt.NewNumericDate(creationTime),
			// Set the expiration to the TTL
			Expiry: squarejwt.NewNumericDate(creationTime.Add(resp.WrapInfo.TTL)),
			// Set a reasonable not-before time; since unwrapping happens on this
			// node we shouldn't have to worry much about drift
			NotBefore: squarejwt.NewNumericDate(time.Now().Add(-5 * time.Second)),
		}
		type privateClaims struct {
			Accessor string `json:"accessor"`
			Type     string `json:"type"`
			Addr     string `json:"addr"`
		}
		priClaims := &privateClaims{
			Type: "wrapping",
			Addr: c.redirectAddr,
		}
		if resp.Auth != nil {
			priClaims.Accessor = resp.Auth.Accessor
		}
		sig, err := jose.NewSigner(
			jose.SigningKey{Algorithm: jose.ES512, Key: c.wrappingJWTKey},
			(&jose.SignerOptions{}).WithType("JWT"))
		if err != nil {
			c.tokenStore.revokeOrphan(ctx, te.ID)
			c.logger.Error("failed to create JWT builder", "error", err)
			return nil, ErrInternalError
		}
		ser, err := squarejwt.Signed(sig).Claims(claims).Claims(priClaims).CompactSerialize()
		if err != nil {
			c.tokenStore.revokeOrphan(ctx, te.ID)
			c.logger.Error("failed to serialize JWT", "error", err)
			return nil, ErrInternalError
		}
		resp.WrapInfo.Token = ser
		if c.redirectAddr == "" {
			resp.AddWarning("No redirect address set in Vault so none could be encoded in the token. You may need to supply Vault's API address when unwrapping the token.")
		}
	}

	cubbyReq := &logical.Request{
		Operation:   logical.CreateOperation,
		Path:        "cubbyhole/response",
		ClientToken: te.ID,
	}
	if sealWrap {
		cubbyReq.WrapInfo = &logical.RequestWrapInfo{
			SealWrap: true,
		}
	}

	// During a rewrap, store the original response, don't wrap it again.
	if req.Path == "sys/wrapping/rewrap" {
		cubbyReq.Data = map[string]interface{}{
			"response": resp.Data["response"],
		}
	} else {
		httpResponse := logical.LogicalResponseToHTTPResponse(resp)

		// Add the unique identifier of the original request to the response
		httpResponse.RequestID = req.ID

		// Because of the way that JSON encodes (likely just in Go) we actually get
		// mixed-up values for ints if we simply put this object in the response
		// and encode the whole thing; so instead we marshal it first, then store
		// the string response. This actually ends up making it easier on the
		// client side, too, as it becomes a straight read-string-pass-to-unmarshal
		// operation.

		marshaledResponse, err := json.Marshal(httpResponse)
		if err != nil {
			c.tokenStore.revokeOrphan(ctx, te.ID)
			c.logger.Error("failed to marshal wrapped response", "error", err)
			return nil, ErrInternalError
		}

		cubbyReq.Data = map[string]interface{}{
			"response": string(marshaledResponse),
		}
	}

	cubbyResp, err := c.router.Route(ctx, cubbyReq)
	if err != nil {
		// Revoke since it's not yet being tracked for expiration
		c.tokenStore.revokeOrphan(ctx, te.ID)
		c.logger.Error("failed to store wrapped response information", "error", err)
		return nil, ErrInternalError
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		c.tokenStore.revokeOrphan(ctx, te.ID)
		c.logger.Error("failed to store wrapped response information", "error", cubbyResp.Data["error"])
		return cubbyResp, nil
	}

	// Store info for lookup
	cubbyReq.WrapInfo = nil
	cubbyReq.Path = "cubbyhole/wrapinfo"
	cubbyReq.Data = map[string]interface{}{
		"creation_ttl":  resp.WrapInfo.TTL,
		"creation_time": creationTime,
	}
	// Store creation_path if not a rewrap
	if req.Path != "sys/wrapping/rewrap" {
		cubbyReq.Data["creation_path"] = req.Path
	} else {
		cubbyReq.Data["creation_path"] = resp.WrapInfo.CreationPath
	}
	cubbyResp, err = c.router.Route(ctx, cubbyReq)
	if err != nil {
		// Revoke since it's not yet being tracked for expiration
		c.tokenStore.revokeOrphan(ctx, te.ID)
		c.logger.Error("failed to store wrapping information", "error", err)
		return nil, ErrInternalError
	}
	if cubbyResp != nil && cubbyResp.IsError() {
		c.tokenStore.revokeOrphan(ctx, te.ID)
		c.logger.Error("failed to store wrapping information", "error", cubbyResp.Data["error"])
		return cubbyResp, nil
	}

	wAuth := &logical.Auth{
		ClientToken: te.ID,
		Policies:    []string{"response-wrapping"},
		LeaseOptions: logical.LeaseOptions{
			TTL:       te.TTL,
			Renewable: false,
		},
	}

	// Register the wrapped token with the expiration manager
	if err := c.expiration.RegisterAuth(ctx, te.Path, wAuth); err != nil {
		// Revoke since it's not yet being tracked for expiration
		c.tokenStore.revokeOrphan(ctx, te.ID)
		c.logger.Error("failed to register cubbyhole wrapping token lease", "request_path", req.Path, "error", err)
		return nil, ErrInternalError
	}

	return nil, nil
}

// ValidateWrappingToken checks whether a token is a wrapping token.
func (c *Core) ValidateWrappingToken(req *logical.Request) (bool, error) {
	if req == nil {
		return false, fmt.Errorf("invalid request")
	}

	var err error

	var token string
	var thirdParty bool
	if req.Data != nil && req.Data["token"] != nil {
		thirdParty = true
		if tokenStr, ok := req.Data["token"].(string); !ok {
			return false, fmt.Errorf("could not decode token in request body")
		} else if tokenStr == "" {
			return false, fmt.Errorf("empty token in request body")
		} else {
			token = tokenStr
		}
	} else {
		token = req.ClientToken
	}

	// Check for it being a JWT. If it is, and it is valid, we extract the
	// internal client token from it and use that during lookup.
	if strings.Count(token, ".") == 2 {
		{
			// Implement the jose library way
			parsedJWT, err := squarejwt.ParseSigned(token)
			if err != nil {
				goto FALLBACK
			}
			var claims squarejwt.Claims
			var allClaims = make(map[string]interface{})
			if err = parsedJWT.Claims(&c.wrappingJWTKey.PublicKey, &claims, &allClaims); err != nil {
				goto FALLBACK
			}
			typeClaimRaw, ok := allClaims["type"]
			if !ok {
				return false, errors.New("could not validate type claim")
			}
			typeClaim, ok := typeClaimRaw.(string)
			if !ok {
				return false, errors.New("could not parse type claim")
			}
			if typeClaim != "wrapping" {
				return false, errors.New("unexpected type claim")
			}
			if !thirdParty {
				req.ClientToken = claims.ID
			} else {
				req.Data["token"] = claims.ID
			}
			goto JWTDONE
		}
	FALLBACK:
		// Eventually simply rejig the signature to be the right format. Look at asymmetric.go in the square library, function signPayload, to see what it's doing
		{
			// This is fallback because the SermoDigital library produces
			// invalid signatures, so we keep this for previously wrapped
			// tokens
			wt, err := jws.ParseJWT([]byte(token))
			// If there's an error we simply fall back to attempting to use it as a regular token
			if err == nil && wt != nil {
				validator := &jwt.Validator{}
				validator.SetClaim("type", "wrapping")
				if err = wt.Validate(&c.wrappingJWTKey.PublicKey, crypto.SigningMethodES512, []*jwt.Validator{validator}...); err != nil {
					return false, errwrap.Wrapf("wrapping token signature could not be validated: {{err}}", err)
				}
				token, _ = wt.Claims().JWTID()
				// We override the given request client token so that the rest of
				// Vault sees the real value. This also ensures audit logs are
				// consistent with the actual token that was issued.
				if !thirdParty {
					req.ClientToken = token
				} else {
					req.Data["token"] = token
				}
			}
		}
	}

JWTDONE:
	if token == "" {
		return false, fmt.Errorf("token is empty")
	}

	if c.Sealed() {
		return false, consts.ErrSealed
	}

	c.stateLock.RLock()
	defer c.stateLock.RUnlock()
	if c.standby {
		return false, consts.ErrStandby
	}

	te, err := c.tokenStore.Lookup(c.activeContext, token)
	if err != nil {
		return false, err
	}
	if te == nil {
		return false, nil
	}

	if len(te.Policies) != 1 {
		return false, nil
	}

	if te.Policies[0] != responseWrappingPolicyName && te.Policies[0] != controlGroupPolicyName {
		return false, nil
	}

	return true, nil
}
