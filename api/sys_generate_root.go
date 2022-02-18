package api

import "context"

func (c *Sys) GenerateRootStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRootStatusContext(ctx)
}

func (c *Sys) GenerateDROperationTokenStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenStatusContext(ctx)
}

func (c *Sys) GenerateRecoveryOperationTokenStatus() (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenStatusContext(ctx)
}

func (c *Sys) GenerateRootStatusContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonContext(ctx, "/v1/sys/generate-root/attempt")
}

func (c *Sys) GenerateDROperationTokenStatusContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) GenerateRecoveryOperationTokenStatusContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommonContext(ctx, "/v1/sys/generate-recovery-token/attempt")
}

func (c *Sys) generateRootStatusCommonContext(ctx context.Context, path string) (*GenerateRootStatusResponse, error) {
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
	return c.GenerateRootInitContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateDROperationTokenInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenInitContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateRecoveryOperationTokenInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenInitContext(ctx, otp, pgpKey)
}

func (c *Sys) GenerateRootInitContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonContext(ctx, "/v1/sys/generate-root/attempt", otp, pgpKey)
}

func (c *Sys) GenerateDROperationTokenInitContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt", otp, pgpKey)
}

func (c *Sys) GenerateRecoveryOperationTokenInitContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommonContext(ctx, "/v1/sys/generate-recovery-token/attempt", otp, pgpKey)
}

func (c *Sys) generateRootInitCommonContext(ctx context.Context, path, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
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
	return c.GenerateRootCancelContext(ctx)
}

func (c *Sys) GenerateDROperationTokenCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenCancelContext(ctx)
}

func (c *Sys) GenerateRecoveryOperationTokenCancel() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenCancelContext(ctx)
}

func (c *Sys) GenerateRootCancelContext(ctx context.Context) error {
	return c.generateRootCancelCommonContext(ctx, "/v1/sys/generate-root/attempt")
}

func (c *Sys) GenerateDROperationTokenCancelContext(ctx context.Context) error {
	return c.generateRootCancelCommonContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) GenerateRecoveryOperationTokenCancelContext(ctx context.Context) error {
	return c.generateRootCancelCommonContext(ctx, "/v1/sys/generate-recovery-token/attempt")
}

func (c *Sys) generateRootCancelCommonContext(ctx context.Context, path string) error {
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
	return c.GenerateRootUpdateContext(ctx, shard, nonce)
}

func (c *Sys) GenerateDROperationTokenUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateDROperationTokenUpdateContext(ctx, shard, nonce)
}

func (c *Sys) GenerateRecoveryOperationTokenUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	return c.GenerateRecoveryOperationTokenUpdateContext(ctx, shard, nonce)
}

func (c *Sys) GenerateRootUpdateContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonContext(ctx, "/v1/sys/generate-root/update", shard, nonce)
}

func (c *Sys) GenerateDROperationTokenUpdateContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonContext(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/update", shard, nonce)
}

func (c *Sys) GenerateRecoveryOperationTokenUpdateContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommonContext(ctx, "/v1/sys/generate-recovery-token/update", shard, nonce)
}

func (c *Sys) generateRootUpdateCommonContext(ctx context.Context, path, shard, nonce string) (*GenerateRootStatusResponse, error) {
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
