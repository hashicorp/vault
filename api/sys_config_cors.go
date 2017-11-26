package api

import (
	"context"
	"net/http"
)

// CORSStatus returns the current CORS configuration.
func (c *Sys) CORSStatus() (*CORSResponse, error) {
	return c.CORSStatusWithContext(context.Background())
}

// CORSStatusWithContext returns the current CORS configuration, with a context.
func (c *Sys) CORSStatusWithContext(ctx context.Context) (*CORSResponse, error) {
	req := c.c.NewRequest(http.MethodGet, "/v1/sys/config/cors")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CORSResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// ConfigureCORS updates the CORS configuration.
func (c *Sys) ConfigureCORS(r *CORSRequest) (*CORSResponse, error) {
	return c.ConfigureCORSWithContext(context.Background(), r)
}

// ConfigureCORSWithContext updates the CORS configuration, with a context.
func (c *Sys) ConfigureCORSWithContext(ctx context.Context, r *CORSRequest) (*CORSResponse, error) {
	req := c.c.NewRequest(http.MethodPut, "/v1/sys/config/cors")
	req = req.WithContext(ctx)

	if err := req.SetJSONBody(r); err != nil {
		return nil, err
	}

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CORSResponse
	err = resp.DecodeJSON(&result)
	return &result, err
}

// DisableCORS disables CORS.
func (c *Sys) DisableCORS() (*CORSResponse, error) {
	return c.DisableCORSWithContext(context.Background())
}

// DisableCORSWithContext disables CORS, with a context.
func (c *Sys) DisableCORSWithContext(ctx context.Context) (*CORSResponse, error) {
	req := c.c.NewRequest(http.MethodDelete, "/v1/sys/config/cors")
	req = req.WithContext(ctx)

	resp, err := c.c.RawRequest(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result CORSResponse
	err = resp.DecodeJSON(&result)
	return &result, err

}

// CORSRequest is used as input to enable CORS.
type CORSRequest struct {
	// AllowOrigins is a comma-separated list of allowed origins.
	AllowedOrigins string `json:"allowed_origins"`

	// Enabled is a boolean indicating whether CORS should be enabled.
	Enabled bool `json:"enabled"`
}

// CORSResponse is the result of a CORS request.
type CORSResponse struct {
	// AllowOrigins is a comma-separated list of allowed origins.
	AllowedOrigins string `json:"allowed_origins"`

	// Enabled is a boolean indicating whether CORS is enabled.
	Enabled bool `json:"enabled"`
}
