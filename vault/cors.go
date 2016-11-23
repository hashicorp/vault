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
	return &CORSConfig{
		enabled:        false,
		allowedOrigins: &regexp.Regexp{},
		allowedHeaders: &map[string]string{
			"Access-Control-Allow-Headers": "X-Requested-With,Content-Type,Accept,Origin,Authorization,X-Vault-Token",
			"Access-Control-Max-Age":       "1800",
			"Content-Type":                 "text/plain",
		},
		allowCredentials: true,
		allowedMethods:   &[]string{http.MethodOptions},
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
	// c.enabled = false
	// c.allowedOrigins = &regexp.Regexp{}
}

func (c *CORSConfig) AllowedMethods() []string {
	return *c.allowedMethods
}

func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) {
	// check that the origin is valid & set the header.
	origin := r.Header.Get("Origin")
	if c.validOrigin(origin) {

		methods := strings.Join(c.AllowedMethods(), ",")
		w.Header().Set("Allow", methods)

		// add the credentials header if allowed.
		if c.allowCredentials {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// apply headers for preflight requests
		if r.Method == http.MethodOptions {
			for k, v := range *c.allowedHeaders {
				w.Header().Set(k, v)
			}
		}
	}
}

func (c *CORSConfig) AllowedOrigins() *regexp.Regexp {
	return c.allowedOrigins
}

func (c *CORSConfig) validOrigin(origin string) bool {
	if len(origin) == 0 {
		return false
	}

	if c.allowedOrigins == nil {
		return false
	}

	if c.allowedOrigins.MatchString(origin) {
		return true
	}

	return false
}
