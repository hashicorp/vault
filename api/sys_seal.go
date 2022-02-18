package api

import "context"

func (c *Sys) SealStatus() (*SealStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.SealStatusContext(ctx)
}

func (c *Sys) SealStatusContext(ctx context.Context) (*SealStatusResponse, error) {
	r := c.c.NewRequest("GET", "/v1/sys/seal-status")
	return sealStatusRequestContext(ctx, c, r)
}

func (c *Sys) Seal() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.SealContext(ctx)
}

func (c *Sys) SealContext(ctx context.Context) error {
	r := c.c.NewRequest("PUT", "/v1/sys/seal")

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) ResetUnsealProcess() (*SealStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.ResetUnsealProcessContext(ctx)
}

func (c *Sys) ResetUnsealProcessContext(ctx context.Context) (*SealStatusResponse, error) {
	body := map[string]interface{}{"reset": true}

	r := c.c.NewRequest("PUT", "/v1/sys/unseal")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	return sealStatusRequestContext(ctx, c, r)
}

func (c *Sys) Unseal(shard string) (*SealStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.UnsealContext(ctx, shard)
}

func (c *Sys) UnsealContext(ctx context.Context, shard string) (*SealStatusResponse, error) {
	body := map[string]interface{}{"key": shard}

	r := c.c.NewRequest("PUT", "/v1/sys/unseal")
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	return sealStatusRequestContext(ctx, c, r)
}

func (c *Sys) UnsealWithOptions(opts *UnsealOpts) (*SealStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.UnsealWithOptionsContext(ctx, opts)
}

func (c *Sys) UnsealWithOptionsContext(ctx context.Context, opts *UnsealOpts) (*SealStatusResponse, error) {
	r := c.c.NewRequest("PUT", "/v1/sys/unseal")
	if err := r.SetJSONBody(opts); err != nil {
		return nil, err
	}

	return sealStatusRequestContext(ctx, c, r)
}

func sealStatusRequestContext(ctx context.Context, c *Sys, r *Request) (*SealStatusResponse, error) {
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result SealStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type SealStatusResponse struct {
	Type         string `json:"type"`
	Initialized  bool   `json:"initialized"`
	Sealed       bool   `json:"sealed"`
	T            int    `json:"t"`
	N            int    `json:"n"`
	Progress     int    `json:"progress"`
	Nonce        string `json:"nonce"`
	Version      string `json:"version"`
	Migration    bool   `json:"migration"`
	ClusterName  string `json:"cluster_name,omitempty"`
	ClusterID    string `json:"cluster_id,omitempty"`
	RecoverySeal bool   `json:"recovery_seal"`
	StorageType  string `json:"storage_type,omitempty"`
}

type UnsealOpts struct {
	Key     string `json:"key"`
	Reset   bool   `json:"reset"`
	Migrate bool   `json:"migrate"`
}
