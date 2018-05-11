package api

import (
	"io"
	"os"
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

		if (operation == "PUT" || operation == "POST") && path == "sys/wrapping/wrap" {
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

func (c *Logical) Read(path string) (*Secret, error) {
	r := c.c.NewRequest("GET", "/v1/"+path)
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) List(path string) (*Secret, error) {
	r := c.c.NewRequest("LIST", "/v1/"+path)
	// Set this for broader compatibility, but we use LIST above to be able to
	// handle the wrapping lookup function
	r.Method = "GET"
	r.Params.Set("list", "true")
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}

func (c *Logical) Write(path string, data map[string]interface{}) (*Secret, error) {
	r := c.c.NewRequest("PUT", "/v1/"+path)
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

func (c *Logical) Delete(path string) (*Secret, error) {
	r := c.c.NewRequest("DELETE", "/v1/"+path)
	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, err
		}
	}
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		return ParseSecret(resp.Body)
	}

	return nil, nil
}

func (c *Logical) Unwrap(wrappingToken string) (*Secret, error) {
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

	r := c.c.NewRequest("PUT", "/v1/sys/wrapping/unwrap")
	if err := r.SetJSONBody(data); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(r)
	if resp != nil {
		defer resp.Body.Close()
	}
	if resp != nil && resp.StatusCode == 404 {
		secret, parseErr := ParseSecret(resp.Body)
		switch parseErr {
		case nil:
		case io.EOF:
			return nil, nil
		default:
			return nil, err
		}
		if secret != nil && (len(secret.Warnings) > 0 || len(secret.Data) > 0) {
			return secret, nil
		}
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return ParseSecret(resp.Body)
}
