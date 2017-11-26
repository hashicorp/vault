package api

import (
	"context"
	"net/http"
)

// TokenAuth is used to perform token backend operations on Vault
type TokenAuth struct {
	c *Client
}

// Token is used to return the client for token-backend API calls
func (a *Auth) Token() *TokenAuth {
	return &TokenAuth{c: a.c}
}

// Create makes a new Vault token with the given options.
func (c *TokenAuth) Create(opts *TokenCreateRequest) (*Secret, error) {
	return c.CreateWithContext(context.Background(), opts)
}

// CreateWithContext makes a new Vault token with the given options, with a
// context.
func (c *TokenAuth) CreateWithContext(ctx context.Context, opts *TokenCreateRequest) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/create")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// CreateOrphan creates an orphan token.
func (c *TokenAuth) CreateOrphan(opts *TokenCreateRequest) (*Secret, error) {
	return c.CreateOrphanWithContext(context.Background(), opts)
}

// CreateOrphanWithContext creates an orphan token, with a context.
func (c *TokenAuth) CreateOrphanWithContext(ctx context.Context, opts *TokenCreateRequest) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/create-orphan")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// CreateWithRole creates a token with the given role.
func (c *TokenAuth) CreateWithRole(opts *TokenCreateRequest, roleName string) (*Secret, error) {
	return c.CreateWithRoleWithContext(context.Background(), opts, roleName)
}

// CreateWithRoleWithContext creates a token with the given role, with a
// context.
func (c *TokenAuth) CreateWithRoleWithContext(ctx context.Context, opts *TokenCreateRequest, roleName string) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/create/"+roleName)
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(opts); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// Lookup gets information about the given token.
func (c *TokenAuth) Lookup(token string) (*Secret, error) {
	return c.LookupWithContext(context.Background(), token)
}

// LookupWithContext gets information about the given token, with a context.
func (c *TokenAuth) LookupWithContext(ctx context.Context, token string) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/lookup")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// LookupAccessor looks up information about the token with the given accessor.
func (c *TokenAuth) LookupAccessor(accessor string) (*Secret, error) {
	return c.LookupAccessorWithContext(context.Background(), accessor)
}

// LookupAccessorWithContext looks up information about the token with the given
// accessor, with a context.
func (c *TokenAuth) LookupAccessorWithContext(ctx context.Context, accessor string) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/lookup-accessor")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"accessor": accessor,
	}); err != nil {
		return nil, err
	}
	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// LookupSelf looks up the current token.
func (c *TokenAuth) LookupSelf() (*Secret, error) {
	return c.LookupSelfWithContext(context.Background())
}

// LookupSelfWithContext looks up the current token, with a context.
func (c *TokenAuth) LookupSelfWithContext(ctx context.Context) (*Secret, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/auth/token/lookup-self")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// Renew renews the given token with the requested increment.
func (c *TokenAuth) Renew(token string, increment int) (*Secret, error) {
	return c.RenewWithContext(context.Background(), token, increment)
}

// RenewWithContext renews the given token with the requested increment, with a
// context.
func (c *TokenAuth) RenewWithContext(ctx context.Context, token string, increment int) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/renew")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"token":     token,
		"increment": increment,
	}); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RenewSelf renews the token on the client with the requested increment.
func (c *TokenAuth) RenewSelf(increment int) (*Secret, error) {
	return c.RenewSelfWithContext(context.Background(), increment)
}

// RenewSelfWithContext renews the token on the client with the requested
// increment, with a context.
func (c *TokenAuth) RenewSelfWithContext(ctx context.Context, increment int) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/renew-self")
	req = req.WithContext(ctx)

	body := map[string]interface{}{"increment": increment}
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RenewTokenAsSelf behaves like renew-self, but authenticates using a provided
// token instead of the token attached to the client.
func (c *TokenAuth) RenewTokenAsSelf(token string, increment int) (*Secret, error) {
	return c.RenewTokenAsSelfWithContext(context.Background(), token, increment)
}

// RenewTokenAsSelfWithConext behaves like renew-self, but authenticates using a
// provided token instead of the token attached to the client, with a context.
func (c *TokenAuth) RenewTokenAsSelfWithContext(ctx context.Context, token string, increment int) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/renew-self")
	req = req.WithContext(ctx)
	req.ClientToken = token

	body := map[string]interface{}{"increment": increment}
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RevokeAccessor revokes a token associated with the given accessor along with
// all the child tokens.
func (c *TokenAuth) RevokeAccessor(accessor string) error {
	return c.RevokeAccessorWithContext(context.Background(), accessor)
}

// RevokeAccessorWithContext revokes a token associated with the given accessor
// along with all the child tokens, with a context.
func (c *TokenAuth) RevokeAccessorWithContext(ctx context.Context, accessor string) error {
	req := c.c.NewRequest(http.MethodPost, "/v1/auth/token/revoke-accessor")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"accessor": accessor,
	}); err != nil {
		return err
	}
	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeOrphan revokes a token without revoking the tree underneath it (so
// child tokens are orphaned rather than revoked).
func (c *TokenAuth) RevokeOrphan(token string) error {
	return c.RevokeOrphanWithContext(context.Background(), token)
}

// RevokeOrphanWithContext revokes a token without revoking the tree underneath
// it (so child tokens are orphaned rather than revoked), with a context.
func (c *TokenAuth) RevokeOrphanWithContext(ctx context.Context, token string) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/revoke-orphan")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeSelf revokes the token making the call. The `token` parameter is kept
// for backwards compatibility but is ignored; only the client's set token has
// an effect.
func (c *TokenAuth) RevokeSelf(token string) error {
	return c.RevokeSelfWithContext(context.Background())
}

// RevokeSelfWithContext revokes the token making the call, with a context.
func (c *TokenAuth) RevokeSelfWithContext(ctx context.Context) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/revoke-self")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// RevokeTree is the "normal" revoke operation that revokes the given token and
// the entire tree underneath -- all of its child tokens, their child tokens,
// etc.
func (c *TokenAuth) RevokeTree(token string) error {
	return c.RevokeTreeWithContext(context.Background(), token)
}

// RevokeTreeWithContext is the "normal" revoke operation that revokes the given token and
// the entire tree underneath -- all of its child tokens, their child tokens,
// etc, with a context.
func (c *TokenAuth) RevokeTreeWithContext(ctx context.Context, token string) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/auth/token/revoke")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(map[string]interface{}{
		"token": token,
	}); err != nil {
		return err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

// TokenCreateRequest is the options structure for creating a token.
type TokenCreateRequest struct {
	// ID is the token - only root tokens can request a token with a specified ID.
	ID string `json:"id,omitempty"`

	// Policies is the list of policies to attach to the token.
	Policies []string `json:"policies,omitempty"`

	// Metadata is arbitrary key-value metadata for the token.
	Metadata map[string]string `json:"meta,omitempty"`

	// Lease is the lease.
	//
	// Deprecated: Lease is deprecated. Use TTL instead.
	Lease string `json:"lease,omitempty"`

	// TTL is the minimum TTL for the token.
	TTL string `json:"ttl,omitempty"`

	// ExplicitMaxTTL is the maximum TTL for the token.
	ExplicitMaxTTL string `json:"explicit_max_ttl,omitempty"`

	// Period is period for the token, if a periodic token. Only root tokens can
	// create periodic tokens.
	Period string `json:"period,omitempty"`

	// NoParent detaches the token from its parent. Only root tokens can create
	// orphan tokens.
	NoParent bool `json:"no_parent,omitempty"`

	// NoDefaultPolicy removes the default policy from the token.
	NoDefaultPolicy bool `json:"no_default_policy,omitempty"`

	// DisplayName is the human name of the token.
	DisplayName string `json:"display_name"`

	// NumUses is the number of uses before the token expires, regardless of
	// remaining TTLs.
	NumUses int `json:"num_uses"`

	// Renewable is a pointer to a bool indicating whether the token is renewable.
	Renewable *bool `json:"renewable,omitempty"`
}
