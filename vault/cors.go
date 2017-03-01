package vault

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/hashicorp/vault/helper/strutil"
	"github.com/hashicorp/vault/logical"
)

var errCORSNotConfigured = errors.New("CORS is not configured")

// CORSConfig stores the state of the CORS configuration.
type CORSConfig struct {
	sync.RWMutex
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins"`
}

func (c *Core) saveCORSConfig() error {
	view := c.systemBarrierView.SubView("config/")

	entry, err := logical.StorageEntryJSON("cors", c.corsConfig)
	if err != nil {
		return fmt.Errorf("failed to create CORS config entry: %v", err)
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
func (c *CORSConfig) Enable(urls string) error {

	if strings.Contains("*", urls) && len(urls) > 1 {
		return errors.New("wildcard must be the only value")
	}

	c.Lock()
	defer c.Unlock()

	c.AllowedOrigins = strings.Split(urls, ",")
	c.Enabled = true

	return nil
}

// IsEnabled returns the value of CORSConfig.isEnabled
func (c *CORSConfig) IsEnabled() bool {
	c.RLock()
	defer c.RUnlock()

	return c.Enabled
}

// Disable sets CORS to disabled and clears the allowed origins
func (c *CORSConfig) Disable() {
	c.Lock()
	defer c.Unlock()

	c.Enabled = false
	c.AllowedOrigins = []string{}
}

// IsValidOrigin determines if the origin of the request is allowed to make
// cross-origin requests based on the CORSConfig.
func (c *CORSConfig) IsValidOrigin(origin string) bool {
	c.RLock()
	defer c.RUnlock()

	if c.AllowedOrigins == nil {
		return false
	}

	if len(c.AllowedOrigins) == 1 && (c.AllowedOrigins)[0] == "*" {
		return true
	}

	return strutil.StrListContains(c.AllowedOrigins, origin)
}
