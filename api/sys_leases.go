package api

import (
	"context"
	"errors"
	"net/http"
)

func (c *Sys) Renew(id string, increment int) (*Secret, error) {
	return c.RenewWithContext(context.Background(), id, increment)
}

// RenewResponse returns the full renew response instead of just the secret.
// Returning the full response allows access to the Age header which is
// vital for cached requests to calculate correct expiration times.
func (c *Sys) RenewResponse(id string, increment int) (*Response, error) {
	return c.RenewResponseWithContext(context.Background(), id, increment)
}

func (c *Sys) RenewWithContext(ctx context.Context, id string, increment int) (*Secret, error) {
	resp, err := c.RenewResponseWithContext(ctx, id, increment)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

// RenewResponseWithContext returns the full renew response instead of just the secret.
// Returning the full response allows access to the Age header which is
// vital for cached requests to calculate correct expiration times.
func (c *Sys) RenewResponseWithContext(ctx context.Context, id string, increment int) (*Response, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/renew")

	body := map[string]interface{}{
		"increment": increment,
		"lease_id":  id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (c *Sys) Lookup(id string) (*Secret, error) {
	return c.LookupWithContext(context.Background(), id)
}

func (c *Sys) LookupWithContext(ctx context.Context, id string) (*Secret, error) {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/lookup")

	body := map[string]interface{}{
		"lease_id": id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *Sys) Revoke(id string) error {
	return c.RevokeWithContext(context.Background(), id)
}

func (c *Sys) RevokeWithContext(ctx context.Context, id string) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke")
	body := map[string]interface{}{
		"lease_id": id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokePrefix(id string) error {
	return c.RevokePrefixWithContext(context.Background(), id)
}

func (c *Sys) RevokePrefixWithContext(ctx context.Context, id string) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke-prefix/"+id)

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokeForce(id string) error {
	return c.RevokeForceWithContext(context.Background(), id)
}

func (c *Sys) RevokeForceWithContext(ctx context.Context, id string) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	r := c.c.NewRequest(http.MethodPut, "/v1/sys/leases/revoke-force/"+id)

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokeWithOptions(opts *RevokeOptions) error {
	return c.RevokeWithOptionsWithContext(context.Background(), opts)
}

func (c *Sys) RevokeWithOptionsWithContext(ctx context.Context, opts *RevokeOptions) error {
	ctx, cancelFunc := c.c.withConfiguredTimeout(ctx)
	defer cancelFunc()

	if opts == nil {
		return errors.New("nil options provided")
	}

	// Construct path
	path := "/v1/sys/leases/revoke/"
	switch {
	case opts.Force:
		path = "/v1/sys/leases/revoke-force/"
	case opts.Prefix:
		path = "/v1/sys/leases/revoke-prefix/"
	}
	path += opts.LeaseID

	r := c.c.NewRequest(http.MethodPut, path)
	if !opts.Force {
		body := map[string]interface{}{
			"sync": opts.Sync,
		}
		if err := r.SetJSONBody(body); err != nil {
			return err
		}
	}

	resp, err := c.c.rawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

type RevokeOptions struct {
	LeaseID string
	Force   bool
	Prefix  bool
	Sync    bool
}
