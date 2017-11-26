package api

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/hashicorp/vault/helper/jsonutil"
)

const (
	wrappedResponseLocation = "cubbyhole/response"
)

var (
	// The default TTL that will be used with `sys/wrapping/wrap`, can be
	// changed
	DefaultWrappingTTL = "5m"

	// The default function used if no other function is set, which honors the
	// env var and wraps `sys/wrapping/wrap`
	DefaultWrappingLookupFunc = func(operation, path string) string {
		if os.Getenv(EnvVaultWrapTTL) != "" {
			return os.Getenv(EnvVaultWrapTTL)
		}

		if (operation == http.MethodPut || operation == http.MethodPost) && path == "sys/wrapping/wrap" {
			return DefaultWrappingTTL
		}

		return ""
	}
)

// Logical is used to perform logical backend operations on Vault.
type Logical struct {
	c *Client
}

// Logical is used to return the client for logical-backend API calls.
func (c *Client) Logical() *Logical {
	return &Logical{c: c}
}

// Read performs an HTTP GET operation on the given path, corresponding to a
// "read" operation.
func (c *Logical) Read(path string) (*Secret, error) {
	return c.ReadWithContext(context.Background(), path)
}

// ReadWithContext performs an HTTP GET operation on the given path with the
// given context, corresponding to a "read" operation.
func (c *Logical) ReadWithContext(ctx context.Context, path string) (*Secret, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/"+path)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

// List performs an HTTP LIST operation on the given path corresponding to a
// "list" operation.
func (c *Logical) List(path string) (*Secret, error) {
	return c.ListWithContext(context.Background(), path)
}

// ListWithContext performs an HTTP LIST operation on the given path with the
// given context, corresponding to a "list" operation.
func (c *Logical) ListWithContext(ctx context.Context, path string) (*Secret, error) {
	req := c.c.NewRequest("LIST", "/v1/"+path)
	req = req.WithContext(ctx)

	// Set this for broader compatibility, but we use LIST above to be able to
	// handle the wrapping lookup function
	req.Method = http.MethodGet
	req.Params.Set("list", "true")
	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

// Write performs an HTTP PUT operation on the given path corresponding to a
// "write" operation.
func (c *Logical) Write(path string, data map[string]interface{}) (*Secret, error) {
	return c.WriteWithContext(context.Background(), path, data)
}

// WriteWithContext performs an HTTP put operation on the given path with the
// given context, corresponding to a "write" operation.
func (c *Logical) WriteWithContext(ctx context.Context, path string, data map[string]interface{}) (*Secret, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/"+path)
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

// Delete performs and HTTP DELETE operation on the given path corresponding to
// a "delete" operation.
func (c *Logical) Delete(path string) (*Secret, error) {
	return c.DeleteWithContext(context.Background(), path)
}

// DeleteWithContext performs an HTTP DELETE operation on the given path with
// the given context, corresponding to a "delete" operation.
func (c *Logical) DeleteWithContext(ctx context.Context, path string) (*Secret, error) {
	req := c.c.NewRequest(http.MethodDelete, "/v1/"+path)
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

// Unwrap unwraps the value in the given token's cubbyhole.
func (c *Logical) Unwrap(wrappingToken string) (*Secret, error) {
	return c.UnwrapWithContext(context.Background(), wrappingToken)
}

// UnwrapWithContext unwraps the value in the given token's cubbyhole with the
// given context.
func (c *Logical) UnwrapWithContext(ctx context.Context, wrappingToken string) (*Secret, error) {
	var data map[string]interface{}
	if wrappingToken != "" {
		if c.c.Token() == "" {
			c.c.SetToken(wrappingToken)
		} else if wrappingToken != c.c.Token() {
			data = map[string]interface{}{
				"token": wrappingToken,
			}
		}
	}

	req := c.c.NewRequest(http.MethodPut, "/v1/sys/wrapping/unwrap")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if resp != nil && resp.StatusCode != 404 {
			return nil, err
		}
	}
	if resp == nil {
		return nil, nil
	}

	switch resp.StatusCode {
	case http.StatusOK: // New method is supported
		return ParseSecret(resp.Body)
	case http.StatusNotFound: // Fall back to old method
	default:
		return nil, nil
	}

	if wrappingToken != "" {
		origToken := c.c.Token()
		defer c.c.SetToken(origToken)
		c.c.SetToken(wrappingToken)
	}

	secret, err := c.Read(wrappedResponseLocation)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %s", wrappedResponseLocation, err)
	}
	if secret == nil {
		return nil, fmt.Errorf("no value found at %s", wrappedResponseLocation)
	}
	if secret.Data == nil {
		return nil, fmt.Errorf("\"data\" not found in wrapping response")
	}
	if _, ok := secret.Data["response"]; !ok {
		return nil, fmt.Errorf("\"response\" not found in wrapping response \"data\" map")
	}

	wrappedSecret := new(Secret)
	buf := bytes.NewBufferString(secret.Data["response"].(string))
	if err := jsonutil.DecodeJSONFromReader(buf, wrappedSecret); err != nil {
		return nil, fmt.Errorf("error unmarshaling wrapped secret: %s", err)
	}

	return wrappedSecret, nil
}
