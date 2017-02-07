package vault

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
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
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins"`
	sync.RWMutex
}

func (c *Core) saveCORSConfig() error {
	view := c.systemBarrierView.SubView("config/")

	entry, err := logical.StorageEntryJSON("cors", c.corsConfig)
	if err != nil {
		return fmt.Errorf("failed to create CORS confif entry: %v", err)
	}

	if err := view.Put(entry); err != nil {
		return fmt.Errorf("failed to save CORS config: %v", err)
	}

	return nil
}

func (c *Core) loadCORSConfig() error {
	view := c.systemBarrierView.SubView("config/")

	// Load the config in
	out, err := view.Get("cors")
	if err != nil {
		return fmt.Errorf("failed to read CORS config: %v", err)
	}
	if out == nil {
		return nil
	}

	config := new(CORSConfig)
	err = out.DecodeJSON(config)
	if err != nil {
		return err
	}

	c.corsConfig = config

	return nil
}

// Enable takes either a '*' or a comma-seprated list of URLs that can make
// cross-origin requests to Vault.
func (c *CORSConfig) Enable(s string) error {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	if strings.Contains("*", s) && len(s) > 1 {
		return errors.New("wildcard must be the only value")
	}

	c.AllowedOrigins = strings.Split(s, ",")
	c.Enabled = true

	return nil
}

// Get returns the state of the CORS configuration.
func (c *CORSConfig) Get() *CORSConfig {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	return c
}

// IsEnabled returns the value of CORSConfig.isEnabled
func (c *CORSConfig) IsEnabled() bool {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	return c.Enabled
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	c.Enabled = false
	c.AllowedOrigins = []string{}
}

// ApplyHeaders examines the CORS configuration and the request to determine
// if the CORS headers should be returned with the response.
func (c *CORSConfig) ApplyHeaders(w http.ResponseWriter, r *http.Request) int {
	c.RWMutex.Lock()
	defer c.RWMutex.Unlock()

	origin := r.Header.Get("Origin")

	w.Header().Set("Access-Control-Allow-Origin", origin)
	w.Header().Set("Vary", "Origin")

	// apply headers for preflight requests
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))

		for k, v := range preflightHeaders {
			w.Header().Set(k, v)
		}
	}

	return http.StatusNoContent
}

// IsValidOrigin determines if the origin of the request is allowed to make
// cross-origin requests based on the CORSConfig.
func (c *CORSConfig) IsValidOrigin(origin string) bool {
	if c.AllowedOrigins == nil {
		return false
	}

	if len(c.AllowedOrigins) == 1 && (c.AllowedOrigins)[0] == "*" {
		return true
	}

	return strutil.StrListContains(c.AllowedOrigins, origin)
}

// IsValidMethod determines if the verb of the HTTP request is allowed.
func (c *CORSConfig) IsValidMethod(method string) bool {
	if method == "" {
		return false
	}

	return strutil.StrListContains(allowedMethods, method)
}
