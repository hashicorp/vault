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
	Enabled        bool
	AllowedOrigins []string
}

func newCORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled: false,
	}
}

func (c *CORSConfig) Enable(s string) error {
	if strings.Contains("*", s) && len(s) > 1 {
		return errors.New("wildcard must be the only value")
	}

	allowedOrigins := strings.Split(s, " ")

	c.AllowedOrigins = allowedOrigins
	c.Enabled = true

	return nil
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.Enabled = false
	c.AllowedOrigins = []string{}
}

// ApplyHeaders examines the CORS configuration and the request to determine
// if the CORS headers should be returned with the response.
func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) int {
	origin := r.Header.Get("Origin")

	// If CORS is not enabled or if no Origin header is present (i.e. the request
	// is from the Vault CLI. A browser will always send an Origin header), then
	// just return a 200.
	if !c.Enabled || origin == "" {
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

func (c *CORSConfig) validOrigin(origin string) bool {
	if c.AllowedOrigins == nil {
		return false
	}

	if len(c.AllowedOrigins) == 1 && (c.AllowedOrigins)[0] == "*" {
		return true
	}

	for _, allowedOrigin := range c.AllowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}

	return false
}
