package api

import (
	"context"
	"errors"
)

func (c *Sys) Renew(id string, increment int) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RenewContext(ctx, id, increment)
}

func (c *Sys) RenewContext(ctx context.Context, id string, increment int) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/leases/renew")

	body := map[string]interface{}{
		"increment": increment,
		"lease_id":  id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *Sys) Lookup(id string) (*Secret, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.LookupContext(ctx, id)
}

func (c *Sys) LookupContext(ctx context.Context, id string) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/leases/lookup")

	body := map[string]interface{}{
		"lease_id": id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ParseSecret(resp.Body)
}

func (c *Sys) Revoke(id string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeContext(ctx, id)
}

func (c *Sys) RevokeContext(ctx context.Context, id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/leases/revoke")
	body := map[string]interface{}{
		"lease_id": id,
	}
	if err := r.SetJSONBody(body); err != nil {
		return err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokePrefix(id string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokePrefixContext(ctx, id)
}

func (c *Sys) RevokePrefixContext(ctx context.Context, id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/leases/revoke-prefix/"+id)

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokeForce(id string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeForceContext(ctx, id)
}

func (c *Sys) RevokeForceContext(ctx context.Context, id string) error {
	r := c.c.NewRequest("PUT", "/v1/sys/leases/revoke-force/"+id)

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) RevokeWithOptions(opts *RevokeOptions) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.RevokeWithOptionsContext(ctx, opts)
}

func (c *Sys) RevokeWithOptionsContext(ctx context.Context, opts *RevokeOptions) error {
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

	r := c.c.NewRequest("PUT", path)
	if !opts.Force {
		body := map[string]interface{}{
			"sync": opts.Sync,
		}
		if err := r.SetJSONBody(body); err != nil {
			return err
		}
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
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
