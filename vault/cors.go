package vault

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

var errCORSNotConfigured = errors.New("CORS is not configured")

type CORSConfig struct {
	enabled          bool
	allowedOrigins   *regexp.Regexp
	allowedHeaders   *map[string]string
	allowedMethods   *[]string
	allowCredentials bool
}

func newCORSConfig() *CORSConfig {
	// defaultOrigins, err := regexp.Compile(expr)
	defaultOrigins := &regexp.Regexp{}
	return &CORSConfig{
		enabled:        false,
		allowedOrigins: defaultOrigins,
		allowedHeaders: &map[string]string{
			"Access-Control-Allow-Headers": "origin,content-type,cache-control,accept,options,authorization,x-requested-with,x-vault-token",
			"Access-Control-Max-Age":       "1800",
			"Content-Type":                 "text/plain",
		},
		allowCredentials: true,
		allowedMethods: &[]string{
			http.MethodDelete,
			http.MethodGet,
			http.MethodOptions,
			http.MethodPost,
			http.MethodPut,
		},
	}
}

func (c *CORSConfig) Enabled() bool {
	return c.enabled
}

func (c *CORSConfig) Enable(s string) error {
	if s == "" {
		return errors.New("regexp cannot be an empty string")
	}

	allowedOrigins, err := regexp.Compile(s)
	if err != nil {
		return err
	}

	c.allowedOrigins = allowedOrigins
	c.enabled = true

	return nil
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c = nil
}

func (c *CORSConfig) AllowedMethods() []string {
	return *c.allowedMethods
}

// ApplyHeaders examines the CORS configuration and the request to determine
// if the CORS headers should be returned with the response.
func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) int {
	// If CORS is not enabled, just return a 200
	if !c.enabled {
		return http.StatusOK
	}

	// Return a 403 if the origin is not
	// allowed to make cross-origin requests.
	origin := r.Header.Get("Origin")
	if !c.validOrigin(origin) {
		return http.StatusForbidden
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)

	// apply headers for preflight requests
	if r.Method == http.MethodOptions {
		methodAllowed := false
		requestedMethod := r.Header.Get("Access-Control-Request-Method")
		for _, method := range c.AllowedMethods() {
			if method == requestedMethod {
				methodAllowed = true
				continue
			}
		}

		if !methodAllowed {
			return http.StatusMethodNotAllowed
		}

		methods := strings.Join(c.AllowedMethods(), ",")
		w.Header().Set("Access-Control-Allow-Methods", methods)

		// add the credentials header if allowed.
		if c.allowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		for k, v := range *c.allowedHeaders {
			w.Header().Set(k, v)
		}
	}

	return http.StatusOK
}

func (c *CORSConfig) AllowedOrigins() *regexp.Regexp {
	return c.allowedOrigins
}

func (c *CORSConfig) validOrigin(origin string) bool {
	if c.allowedOrigins == nil {
		return false
	}

	if c.allowedOrigins.MatchString(origin) {
		return true
	}

	return false
}
