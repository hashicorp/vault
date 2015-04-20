package api

// TokenAuth is used to perform token backend operations on Vault.
type TokenAuth struct {
	c *Client
}

// TokenAuth is used to return the client for logical-backend API calls.
func (a *Auth) Token() *TokenAuth {
	return &TokenAuth{c: a.c}
}

func (c *TokenAuth) Create(opts *TokenCreateRequest) (*Secret, error) {
	r := c.c.NewRequest("POST", "/v1/auth/token/create")
	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) Renew(token string, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/auth/token/renew/"+token)

	body := map[string]interface{}{"increment": increment}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *TokenAuth) RevokeOrphan(token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke-orphan/"+token)
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *TokenAuth) RevokePrefix(token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke-prefix/"+token)
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *TokenAuth) RevokeTree(token string) error {
	r := c.c.NewRequest("PUT", "/v1/auth/token/revoke/"+token)
	resp, err := c.c.RawRequest(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// TokenCreateRequest is the options structure for creating a token.
type TokenCreateRequest struct {
	ID          string            `json:"id,omitempty"`
	Policies    []string          `json:"policies,omitempty"`
	Metadata    map[string]string `json:"meta,omitempty"`
	Lease       string            `json:"lease,omitempty"`
	NoParent    bool              `json:"no_parent,omitempty"`
	DisplayName string            `json:"display_name"`
	NumUses     int               `json:"num_uses"`
}
