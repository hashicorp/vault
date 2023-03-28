package acme

import (
	"encoding/json"
	"fmt"

	jose "gopkg.in/square/go-jose.v2"
)

// This wraps a JWS message structure.
type JWSCtx struct {
	Algo  string          `json:"alg"`
	Kid   string          `json:"kid"`
	jwk   json.RawMessage `json:"jwk"`
	Nonce string          `json:"nonce"`
	Url   string          `json:"url"`
	key   jose.JSONWebKey `json:"-"`
}

func (c *JWSCtx) UnmarshalJSON(a *ACMEState, jws []byte) error {
	var err error
	if err = json.Unmarshal(jws, c); err != nil {
		return err
	}

	if c.Kid != "" && len(c.jwk) > 0 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The "jwk" and "kid" fields are mutually exclusive.  Servers MUST
		// > reject requests that contain both.
		return fmt.Errorf("invalid header: got both account 'kid' and 'jwk' in the same message; expected only one")
	}

	if c.Kid == "" && len(c.jwk) == 0 {
		// See RFC 8555 Section 6.2 Request Authorization:
		//
		// > Either "jwk" (JSON Web Key) or "kid" (Key ID) as specified
		// > below
		return fmt.Errorf("invalid header: got neither required fields of 'kid' nor 'jwk'")
	}

	if c.Kid != "" {
		// Load KID from storage first.
		c.jwk, err = a.LoadJWK(c.Kid)
		if err != nil {
			return err
		}
	}

	if err = c.key.UnmarshalJSON(c.jwk); err != nil {
		return err
	}

	if !c.key.Valid() {
		return fmt.Errorf("received invalid jwk")
	}

	return nil
}

func hasValues(h jose.Header) bool {
	return h.KeyID != "" || h.JSONWebKey != nil || h.Algorithm != "" || h.Nonce != "" || len(h.ExtraHeaders) > 0
}

func (c *JWSCtx) VerifyJWS(signature string) (map[string]interface{}, error) {
	sig, err := jose.ParseSigned(signature)
	if err != nil {
		return nil, fmt.Errorf("error parsing signature: %w", err)
	}

	if len(sig.Signatures) > 1 {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS MUST NOT have multiple signatures
		return nil, fmt.Errorf("request had multiple signatures")
	}

	if hasValues(sig.Signatures[0].Unprotected) {
		// See RFC 8555 Section 6.2. Request Authentication:
		//
		// > The JWS Unprotected Header [RFC7515] MUST NOT be used
		return nil, fmt.Errorf("request had unprotected headers")
	}

	payload, err := sig.Verify(c.key)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	if err := json.Unmarshal(payload, &m); err != nil {
		return nil, fmt.Errorf("failed to json unmarshal 'payload': %w", err)
	}

	return m, nil
}
