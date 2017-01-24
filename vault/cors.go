package vault

import (
	"errors"
	"net/http"
	"strings"
)

var errCORSNotConfigured = errors.New("CORS is not configured")

var allowedHeaders = map[string]string{
	"Access-Control-Allow-Headers": "origin,content-type,cache-control,accept,options,authorization,x-requested-with,x-vault-token",
	"Access-Control-Max-Age":       "1800",
}

var allowedMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
}

type CORSConfig struct {
	enabled        bool
	allowedOrigins []string
}

func newCORSConfig() *CORSConfig {
	return &CORSConfig{
		enabled: false,
	}
}

func (c *CORSConfig) Enabled() bool {
	return c.enabled
}

func (c *CORSConfig) Enable(s string) error {
	if strings.Contains("*", s) && len(s) > 1 {
		return errors.New("wildcard must be the only value")
	}

	allowedOrigins := strings.Split(s, " ")

	c.allowedOrigins = allowedOrigins
	c.enabled = true

	return nil
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.enabled = false
	c.allowedOrigins = []string{}
}

// ApplyHeaders examines the CORS configuration and the request to determine
// if the CORS headers should be returned with the response.
func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) int {
	origin := r.Header.Get("Origin")

	// If CORS is not enabled or if no Origin header is present (i.e. the request
	// is from the Vault CLI. A browser will always send an Origin header), then
	// just return a 200.
	if !c.enabled || origin == "" {
		return http.StatusOK
	}

	// Return a 403 if the origin is not
	// allowed to make cross-origin requests.
	if !c.validOrigin(origin) {
		return http.StatusForbidden
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")

	// apply headers for preflight requests
	if r.Method == http.MethodOptions {
		methodAllowed := false
		requestedMethod := r.Header.Get("Access-Control-Request-Method")
		for _, method := range allowedMethods {
			if method == requestedMethod {
				methodAllowed = true
				continue
			}
		}

		if !methodAllowed {
			return http.StatusMethodNotAllowed
		}

		methods := strings.Join(allowedMethods, ",")
		w.Header().Set("Access-Control-Allow-Methods", methods)

		for k, v := range allowedHeaders {
			w.Header().Set(k, v)
		}
	}

	return http.StatusOK
}

// AllowedOrigins returns a space-separated list of origins which can make
// cross-origin requests.
func (c *CORSConfig) AllowedOrigins() string {
	return strings.Join(c.allowedOrigins, " ")
}

func (c *CORSConfig) validOrigin(origin string) bool {
	if c.allowedOrigins == nil {
		return false
	}

	if len(c.allowedOrigins) == 1 && (c.allowedOrigins)[0] == "*" {
		return true
	}

	for _, allowedOrigin := range c.allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}
