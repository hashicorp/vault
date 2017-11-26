package api

import (
	"context"
	"net/http"
)

// GenerateRootStatus gets the current root generation status.
func (c *Sys) GenerateRootStatus() (*GenerateRootStatusResponse, error) {
	return c.GenerateRootStatusWithContext(context.Background())
}

// GenerateRootStatusWithContext gets the current root generation status, with a
// context.
func (c *Sys) GenerateRootStatusWithContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommon(ctx, "/v1/sys/generate-root/attempt")
}

// GenerateDROperationTokenStatus gets the current DR root token generation
// status.
func (c *Sys) GenerateDROperationTokenStatus() (*GenerateRootStatusResponse, error) {
	return c.GenerateDROperationTokenStatusWithContext(context.Background())
}

// GenerateDROperationTokenStatusWithContext gets the current DR root token generation
// status, with a context.
func (c *Sys) GenerateDROperationTokenStatusWithContext(ctx context.Context) (*GenerateRootStatusResponse, error) {
	return c.generateRootStatusCommon(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) generateRootStatusCommon(ctx context.Context, path string) (*GenerateRootStatusResponse, error) {
	req := c.c.NewRequest(http.MethodGet, path)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// GenerateRootInit initializes a new root token generation.
func (c *Sys) GenerateRootInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.GenerateRootInitWithContext(context.Background(), otp, pgpKey)
}

// GenerateRootInitWithContext initializes a new root token generation, with a
// context.
func (c *Sys) GenerateRootInitWithContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommon(ctx, "/v1/sys/generate-root/attempt", otp, pgpKey)
}

// GenerateDROperationTokenInit initializes a new DR root token generation.
func (c *Sys) GenerateDROperationTokenInit(otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.GenerateDROperationTokenInitWithContext(context.Background(), otp, pgpKey)
}

// GenerateDROperationTokenInitWithContext initializes a new DR root token
// generation, with a context.
func (c *Sys) GenerateDROperationTokenInitWithContext(ctx context.Context, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	return c.generateRootInitCommon(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt", otp, pgpKey)
}

func (c *Sys) generateRootInitCommon(ctx context.Context, path, otp, pgpKey string) (*GenerateRootStatusResponse, error) {
	req := c.c.NewRequest(http.MethodPut, path)
	req = req.WithContext(ctx)

	body := map[string]interface{}{
		"otp":     otp,
		"pgp_key": pgpKey,
	}
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// GenerateRootCancel cancels an existing root generation.
func (c *Sys) GenerateRootCancel() error {
	return c.GenerateRootCancelWithContext(context.Background())
}

// GenerateRootCancelWithContext cancels an existing root generation, with a
// context.
func (c *Sys) GenerateRootCancelWithContext(ctx context.Context) error {
	return c.generateRootCancelCommon(ctx, "/v1/sys/generate-root/attempt")
}

// GenerateDROperationTokenCancel cancles an existing DR root generation.
func (c *Sys) GenerateDROperationTokenCancel() error {
	return c.GenerateDROperationTokenCancelWithContext(context.Background())
}

// GenerateDROperationTokenCancelWithContext cancles an existing DR root
// generation, with a context.
func (c *Sys) GenerateDROperationTokenCancelWithContext(ctx context.Context) error {
	return c.generateRootCancelCommon(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/attempt")
}

func (c *Sys) generateRootCancelCommon(ctx context.Context, path string) error {
	req := c.c.NewRequest(http.MethodDelete, path)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err == nil {
		defer resp.Body.Close()
	}
	return err
}

// GenerateRootUpdate submits a shard to the root generation process.
func (c *Sys) GenerateRootUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.GenerateRootUpdateWithContext(context.Background(), shard, nonce)
}

// GenerateRootUpdateWithContext submits a shard to the root generation process,
// with a context.
func (c *Sys) GenerateRootUpdateWithContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommon(ctx, "/v1/sys/generate-root/update", shard, nonce)
}

// GenerateDROperationTokenUpdate submits a shard to the DR token generation
// process.
func (c *Sys) GenerateDROperationTokenUpdate(shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.GenerateDROperationTokenUpdateWithContext(context.Background(), shard, nonce)
}

// GenerateDROperationTokenUpdateWithContext submits a shard to the DR token
// generation process, with a context.
func (c *Sys) GenerateDROperationTokenUpdateWithContext(ctx context.Context, shard, nonce string) (*GenerateRootStatusResponse, error) {
	return c.generateRootUpdateCommon(ctx, "/v1/sys/replication/dr/secondary/generate-operation-token/update", shard, nonce)
}

func (c *Sys) generateRootUpdateCommon(ctx context.Context, path, shard, nonce string) (*GenerateRootStatusResponse, error) {
	req := c.c.NewRequest(http.MethodPut, path)
	req = req.WithContext(ctx)

	body := map[string]interface{}{
		"key":   shard,
		"nonce": nonce,
	}
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result GenerateRootStatusResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// GenerateRootStatusResponse is the response from the generate root operation.
type GenerateRootStatusResponse struct {
	// Nonce is the operation nonce.
	Nonce string

	// Started is a boolean indicating whether the operation is started.
	Started bool

	// Progress is the current root token generation progress.
	Progress int

	// Required is the number of required keys.
	Required int

	// Complete is a boolean indicating whether the operation is complete.
	Complete bool

	// EncodedToken is the encoded token.
	EncodedToken string `json:"encoded_token"`

	// EncodedRootToken is the encoded root token.
	EncodedRootToken string `json:"encoded_root_token"`

	// PGPFingerprint is the PGP fingerprint, if given.
	PGPFingerprint string `json:"pgp_fingerprint"`
}
