package api

import "context"

func (c *Sys) GenerateRootStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRootStatusWithContext(ctx)
}

func (c *Sys) GenerateDROperationTokenStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenStatusWithContext(ctx)
}

func (c *Sys) GenerateRecoveryOperationTokenStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenStatusWithContext(ctx)
}

func (c *Sys) GenerateRootStatusWithContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonWithContext(ctx, "/v1/sys/generate-root/attempt")
}

func (c *Sys) GenerateDROperationTokenStatusWithContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonWithContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) GenerateRecoveryOperationTokenStatusWithContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonWithContext(ctx, "/v1/sys/generate-recovery-token/attempt")
}

func (c *Sys) generateRootStatusCommonWithContext(ctx context.Context, path string) (*GenerateRootStatusResponse, error) {
	r := c.c.NewRequest("GET", path)

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) GenerateRootInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRootInitWithContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateDROperationTokenInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenInitWithContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateRecoveryOperationTokenInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenInitWithContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateRootInitWithContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonWithContext(ctx, "/v1/sys/generate-root/attempt", otp, pgpKey)
}

func (c *Sys) GenerateDROperationTokenInitWithContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonWithContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt", otp, pgpKey)
}

func (c *Sys) GenerateRecoveryOperationTokenInitWithContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonWithContext(ctx, "/v1/sys/generate-recovery-token/attempt", otp, pgpKey)
}

func (c *Sys) generateRootInitCommonWithContext(ctx context.Context, path, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	body := map[string]interface{}{
		"otp":     otp,
		"pgp_key": pgpKey,
	}

	r := c.c.NewRequest("PUT", path)
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

func (c *Sys) GenerateRootCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRootCancelWithContext(ctx)
}

func (c *Sys) GenerateDROperationTokenCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenCancelWithContext(ctx)
}

func (c *Sys) GenerateRecoveryOperationTokenCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenCancelWithContext(ctx)
}

func (c *Sys) GenerateRootCancelWithContext(ctx context.Context) error {
	return c.generateRootCancelCommonWithContext(ctx, "/v1/sys/generate-root/attempt")
}

func (c *Sys) GenerateDROperationTokenCancelWithContext(ctx context.Context) error {
	return c.generateRootCancelCommonWithContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) GenerateRecoveryOperationTokenCancelWithContext(ctx context.Context) error {
	return c.generateRootCancelCommonWithContext(ctx, "/v1/sys/generate-recovery-token/attempt")
}

func (c *Sys) generateRootCancelCommonWithContext(ctx context.Context, path string) error {
	r := c.c.NewRequest("DELETE", path)

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

func (c *Sys) GenerateRootUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRootUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) GenerateDROperationTokenUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) GenerateRecoveryOperationTokenUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenUpdateWithContext(ctx, shard, nonce)
}

func (c *Sys) GenerateRootUpdateWithContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonWithContext(ctx, "/v1/sys/generate-root/update", shard, nonce)
}

func (c *Sys) GenerateDROperationTokenUpdateWithContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonWithContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/update", shard, nonce)
}

func (c *Sys) GenerateRecoveryOperationTokenUpdateWithContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonWithContext(ctx, "/v1/sys/generate-recovery-token/update", shard, nonce)
}

func (c *Sys) generateRootUpdateCommonWithContext(ctx context.Context, path, shard, nonce string) (*GenerateRootStatusResponse, error) {
	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}

	r := c.c.NewRequest("PUT", path)
	if err := r.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

type GenerateRootStatusResponse struct {
	Nonce            string `json:"nonce"`
	Started          bool   `json:"started"`
	Progress         int    `json:"progress"`
	Required         int    `json:"required"`
	Complete         bool   `json:"complete"`
	EncodedToken     string `json:"encoded_token"`
	EncodedRootToken string `json:"encoded_root_token"`
	PGPFingerprint   string `json:"pgp_fingerprint"`
	OTP              string `json:"otp"`
	OTPLength        int    `json:"otp_length"`
}
