package vault

import (
	"errors"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/strutil"
)

var errCORSNotConfigured = errors.New("CORS is not configured")

var preflightHeaders = map[string]string{
	"Access-Control-Allow-Headers":     "*",
	"Access-Control-Max-Age":           "1800",
	"Access-Control-Allow-Credentials": "true",
}

var allowedMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodOptions,
	http.MethodPost,
	http.MethodPut,
	"LIST", // LIST is not an official HTTP method, but Vault supports it.
}

// CORSConfig stores the state of the CORS configuration.
type CORSConfig struct {
	isEnabled      bool
	allowedOrigins []string
	mutex          *sync.RWMutex
}

// Enable takes either a '*' or a comma-seprated list of URLs that can make
// cross-origin requests to Vault.
func (c *CORSConfig) Enable(s string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if strings.Contains("*", s) && len(s) > 1 {
		return errors.New("wildcard must be the only value")
	}

	c.allowedOrigins = strings.Split(s, ",")
	c.isEnabled = true

	return nil
}

// Get returns the state of the CORS configuration.
func (c *CORSConfig) Get() *CORSConfig {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c
}

// Enabled returns the value of CORSConfig.isEnabled
func (c *CORSConfig) Enabled() bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.isEnabled
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.isEnabled = false
	c.allowedOrigins = []string{}
}

// ApplyHeaders examines the CORS configuration and the request to determine
// if the CORS headers should be returned with the response.
func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	origin := r.Header.Get("Origin")

	// Return a 403 if the origin is not
	// allowed to make cross-origin requests.
	if !c.validOrigin(origin) {
		return http.StatusForbidden
	}

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")

	// apply headers for preflight requests
	if r.Method == http.MethodOptions {
		requestedMethod := r.Header.Get("Access-Control-Request-Method")

		if !strutil.StrListContains(allowedMethods, requestedMethod) {
			return http.StatusMethodNotAllowed
		}

		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))

		for k, v := range preflightHeaders {
			w.Header().Set(k, v)
		}
	}

	return http.StatusNoContent
}

func (c *CORSConfig) validOrigin(origin string) bool {
	if c.allowedOrigins == nil {
		return false
	}

	if len(c.allowedOrigins) == 1 && (c.allowedOrigins)[0] == "*" {
		return true
	}

	return strutil.StrListContains(c.allowedOrigins, origin)
}
