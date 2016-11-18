package vault

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
)

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

func (c *CORSConfig) AddMethod(newMethod string) *[]string {
	// Do not duplicate methods.
	for _, method := range *c.allowedMethods {
		if method == newMethod {
			return c.allowedMethods
		}
	}

	*c.allowedMethods = append(*c.allowedMethods, newMethod)
	return c.allowedMethods
}

func (c *CORSConfig) AllowedMethods() []string {
	return *c.allowedMethods
}

func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) error {
	// check that the origin is valid & set the header.
	origin := r.Header.Get("Origin")
	validOrigin, err := c.validOrigin(origin)
	if err != nil {
		return err
	}
	if !validOrigin {
		return nil
	}

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

	return nil
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.enabled = false
	c.allowedOrigins = &regexp.Regexp{}
}

func (c *CORSConfig) AllowedOrigins() *regexp.Regexp {
	return c.allowedOrigins
}

func (c *CORSConfig) EnableCORS(s string) (*regexp.Regexp, error) {
	if s == "" {
		return nil, errors.New("regexp cannot be an empty string")
	}

	allowedOrigins, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}

	c.allowedOrigins = allowedOrigins
	c.enabled = true

	return c.allowedOrigins, nil
}

func (c *CORSConfig) validOrigin(origin string) (bool, error) {
	if len(origin) == 0 {
		return false, errors.New("origin is empty")
	}

	if c.allowedOrigins == nil {
		return false, errors.New("no origins are allowed")
	}

	if c.allowedOrigins.MatchString(origin) {
		return true, nil
	}

	return false, nil
}
