// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"bytes"
	"crypto"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-jose/go-jose/v3"
)

var AllowedOuterJWSTypes = map[string]interface{}{
	"RS256":  true,
	"RS384":  true,
	"RS512":  true,
	"PS256":  true,
	"PS384":  true,
	"PS512":  true,
	"ES256":  true,
	"ES384":  true,
	"ES512":  true,
	"EdDSA2": true,
}

var AllowedEabJWSTypes = map[string]interface{}{
	"HS256": true,
	"HS384": true,
	"HS512": true,
}

// This wraps a JWS message structure.
type jwsCtx struct {
	Algo     string          `json:"alg"`
	Kid      string          `json:"kid"`
	Jwk      json.RawMessage `json:"jwk"`
	Nonce    string          `json:"nonce"`
	Url      string          `json:"url"`
	Key      jose.JSONWebKey `json:"-"`
	Existing bool            `json:"-"`
}

func (c *jwsCtx) GetKeyThumbprint() (string, error) {
	keyThumbprint, err := c.Key.Thumbprint(crypto.SHA256)
	if err != nil {
		return "", fmt.Errorf("failed creating thumbprint: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(keyThumbprint), nil
}

func UnmarshalEabJwsJson(eabBytes []byte) (*jwsCtx, error) {
	var eabJws jwsCtx
	var err error
	if err = json.Unmarshal(eabBytes, &eabJws); err != nil {
		return nil, err
	}

	if eabJws.Kid == "" {
		return nil, fmt.Errorf("invalid header: got missing required field 'kid': %w", ErrMalformed)
	}

	if _, present := AllowedEabJWSTypes[eabJws.Algo]; !present {
		return nil, fmt.Errorf("invalid header: unexpected value for 'algo': %w", ErrMalformed)
	}

	return &eabJws, nil
}

func (c *jwsCtx) UnmarshalOuterJwsJson(a *acmeState, ac *acmeContext, jws []byte) error {
	var err error
	if err = json.Unmarshal(jws, c); err != nil {
		return err
	}

	if c.Kid != "" && len(c.Jwk) > 0 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The "jwk" and "kid" fields are mutually exclusive.  Servers MUST
		// > reject requests that contain both.
		return fmt.Errorf("invalid header: got both account 'kid' and 'jwk' in the same message; expected only one: %w", ErrMalformed)
	}

	if c.Kid == "" && len(c.Jwk) == 0 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > Either "jwk" (JSON Web Key) or "kid" (Key ID) as specified
		// > below
		return fmt.Errorf("invalid header: got neither required fields of 'kid' nor 'jwk': %w", ErrMalformed)
	}

	if _, present := AllowedOuterJWSTypes[c.Algo]; !present {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS Protected Header MUST include the following fields:
		// >
		// > - "alg" (Algorithm)
		// >
		// >   * This field MUST NOT contain "none" or a Message
		// >     Authentication Code (MAC) algorithm (e.g. one in which the
		// >     algorithm registry description mentions MAC/HMAC).
		return fmt.Errorf("invalid header: unexpected value for 'algo': %w", ErrMalformed)
	}

	if c.Kid != "" {
		// Load KID from storage first.
		kid := getKeyIdFromAccountUrl(c.Kid)
		c.Jwk, err = a.LoadJWK(ac, kid)
		if err != nil {
			return err
		}
		c.Kid = kid // Use the uuid itself, not the full account url that was originally provided to us.
		c.Existing = true
	}

	if err = c.Key.UnmarshalJSON(c.Jwk); err != nil {
		return err
	}

	if !c.Key.Valid() {
		return fmt.Errorf("received invalid jwk: %w", ErrMalformed)
	}

	if c.Kid == "" {
		c.Kid = genUuid()
		c.Existing = false
	}

	return nil
}

func getKeyIdFromAccountUrl(accountUrl string) string {
	pieces := strings.Split(accountUrl, "/")
	return pieces[len(pieces)-1]
}

func hasValues(h jose.Header) bool {
	return h.KeyID != "" || h.JSONWebKey != nil || h.Algorithm != "" || h.Nonce != "" || len(h.ExtraHeaders) > 0
}

func (c *jwsCtx) VerifyJWS(signature string) (map[string]interface{}, error) {
	// See RFC 8555 Section 6.2. Request Authentication:
	//
	// > The JWS Unencoded Payload Option [RFC7797] MUST NOT be used
	//
	// This is validated by go-jose.
	sig, err := jose.ParseSigned(signature)
	if err != nil {
		return nil, fmt.Errorf("error parsing signature: %s: %w", err, ErrMalformed)
	}

	if len(sig.Signatures) > 1 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS MUST NOT have multiple signatures
		return nil, fmt.Errorf("request had multiple signatures: %w", ErrMalformed)
	}

	if hasValues(sig.Signatures[0].Unprotected) {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS Unprotected Header [RFC7515] MUST NOT be used
		return nil, fmt.Errorf("request had unprotected headers: %w", ErrMalformed)
	}

	payload, err := sig.Verify(c.Key)
	if err != nil {
		return nil, err
	}

	if len(payload) == 0 {
		// Distinguish POST-AS-GET from POST-with-an-empty-body.
		return nil, nil
	}

	var m map[string]interface{}
	if err := json.Unmarshal(payload, &m); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal 'payload': %s: %w", err, ErrMalformed)
	}

	return m, nil
}

func verifyEabPayload(acmeState *acmeState, ac *acmeContext, outer *jwsCtx, expectedPath string, payload map[string]interface{}) (*eabType, error) {
	// Parse the key out.
	rawProtectedBase64, ok := payload["protected"]
	if !ok {
		return nil, fmt.Errorf("missing required field 'protected': %w", ErrMalformed)
	}
	jwkBase64, ok := rawProtectedBase64.(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse 'protected' field: %w", ErrMalformed)
	}

	jwkBytes, err := base64.RawURLEncoding.DecodeString(jwkBase64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64 parse eab 'protected': %s: %w", err, ErrMalformed)
	}

	eabJws, err := UnmarshalEabJwsJson(jwkBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to json unmarshal eab 'protected': %w", err)
	}

	if len(eabJws.Url) == 0 {
		return nil, fmt.Errorf("missing required parameter 'url' in eab 'protected': %w", ErrMalformed)
	}
	expectedUrl := ac.clusterUrl.JoinPath(expectedPath).String()
	if expectedUrl != eabJws.Url {
		return nil, fmt.Errorf("invalid value for 'url' in eab 'protected': got '%v' expected '%v': %w", eabJws.Url, expectedUrl, ErrUnauthorized)
	}

	if len(eabJws.Nonce) != 0 {
		return nil, fmt.Errorf("nonce should not be provided in eab 'protected': %w", ErrMalformed)
	}

	rawPayloadBase64, ok := payload["payload"]
	if !ok {
		return nil, fmt.Errorf("missing required field eab 'payload': %w", ErrMalformed)
	}
	payloadBase64, ok := rawPayloadBase64.(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse 'payload' field: %w", ErrMalformed)
	}

	rawSignatureBase64, ok := payload["signature"]
	if !ok {
		return nil, fmt.Errorf("missing required field 'signature': %w", ErrMalformed)
	}
	signatureBase64, ok := rawSignatureBase64.(string)
	if !ok {
		return nil, fmt.Errorf("failed to parse 'signature' field: %w", ErrMalformed)
	}

	// go-jose only seems to support compact signature encodings.
	compactSig := fmt.Sprintf("%v.%v.%v", jwkBase64, payloadBase64, signatureBase64)
	sig, err := jose.ParseSigned(compactSig)
	if err != nil {
		return nil, fmt.Errorf("error parsing eab signature: %s: %w", err, ErrMalformed)
	}

	if len(sig.Signatures) > 1 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS MUST NOT have multiple signatures
		return nil, fmt.Errorf("eab had multiple signatures: %w", ErrMalformed)
	}

	if hasValues(sig.Signatures[0].Unprotected) {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS Unprotected Header [RFC7515] MUST NOT be used
		return nil, fmt.Errorf("eab had unprotected headers: %w", ErrMalformed)
	}

	// Load the EAB to validate the signature against
	eabEntry, err := acmeState.LoadEab(ac.sc, eabJws.Kid)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to verify eab", ErrUnauthorized)
	}

	verifiedPayload, err := sig.Verify(eabEntry.PrivateBytes)
	if err != nil {
		return nil, err
	}

	// Make sure how eab payload matches the outer JWK key value
	if !bytes.Equal(outer.Jwk, verifiedPayload) {
		return nil, fmt.Errorf("eab payload does not match outer JWK key: %w", ErrMalformed)
	}

	if eabEntry.AcmeDirectory != ac.acmeDirectory {
		// This EAB was not created for this specific ACME directory, reject it
		return nil, fmt.Errorf("%w: failed to verify eab", ErrUnauthorized)
	}

	return eabEntry, nil
}
