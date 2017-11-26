package api

import (
	"context"
	"fmt"
	"net/http"
)

// CapabilitiesSelf returns the list of capabilities for the current token on
// the given path.
func (c *Sys) CapabilitiesSelf(path string) ([]string, error) {
	return c.CapabilitiesSelfWithContext(context.Background(), path)
}

// CapabilitiesSelfWithContext returns the list of capabilities for the current
// token on the given path, with a context.
func (c *Sys) CapabilitiesSelfWithContext(ctx context.Context, path string) ([]string, error) {
	return c.CapabilitiesWithContext(ctx, c.c.Token(), path)
}

// Capabilities returns the list of capabilities for the given token on the
// given path.
func (c *Sys) Capabilities(token, path string) ([]string, error) {
	return c.CapabilitiesWithContext(context.Background(), token, path)
}

// CapabilitiesWithContext returns the list of capabilities for the given token
// on the given path, with a context.
func (c *Sys) CapabilitiesWithContext(ctx context.Context, token, path string) ([]string, error) {
	reqPath := "/v1/sys/capabilities"
	if token == c.c.Token() {
		reqPath = fmt.Sprintf("%s-self", reqPath)
	}

	req := c.c.NewRequest(http.MethodPost, reqPath)
	req = req.WithContext(ctx)

	body := map[string]string{
		"token": token,
		"path":  path,
	}
	if err := req.SetJSONBody(body); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = resp.DecodeJSON(&result)
	if err != nil {
		return nil, err
	}

	var capabilities []string
	capabilitiesRaw := result["capabilities"].([]interface{})
	for _, capability := range capabilitiesRaw {
		capabilities = append(capabilities, capability.(string))
	}
	return capabilities, nil
}
