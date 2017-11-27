package api

import (
	"context"
	"net/http"
)

// Renew renews the given lease id. The optional increment parameter can be used
// to request additional time on the lease. Vault does not have to honor this
// request.
func (c *Sys) Renew(id string, increment int) (*Secret, error) {
	return c.RenewWithContext(context.Background(), id, increment)
}

// RenewWithContext renews the given lease id with a context. The optional increment
// parameter can be used to request additional time on the lease. Vault does not
// have to honor this request.
func (c *Sys) RenewWithContext(ctx context.Context, id string, increment int) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/renew")
	req = req.WithContext(ctx)

	body := map[string]interface{}{
		"increment": increment,
		"lease_id":  id,
	}
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

// Revoke revokes the given lease id.
func (c *Sys) Revoke(id string) error {
	return c.RevokeWithContext(context.Background(), id)
}

// RevokeWithContext revokes the given lease id.
func (c *Sys) RevokeWithContext(ctx context.Context, id string) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke/"+id)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// RevokePrefix revokes all leases which begin with the given prefix.
func (c *Sys) RevokePrefix(id string) error {
	return c.RevokePrefixWithContext(context.Background(), id)
}

// RevokePrefixWithContext revokes all leases which begin with the given prefix,
// with a context.
func (c *Sys) RevokePrefixWithContext(ctx context.Context, id string) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke-prefix/"+id)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// RevokeForce force revokes the given leases.
func (c *Sys) RevokeForce(id string) error {
	return c.RevokeForceWithContext(context.Background(), id)
}

// RevokeForceWithContext force revokes the given leases, with a context.
func (c *Sys) RevokeForceWithContext(ctx context.Context, id string) error {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke-force/"+id)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}
