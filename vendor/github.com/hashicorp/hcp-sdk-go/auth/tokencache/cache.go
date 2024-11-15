// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokencache

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/hashicorp/hcp-sdk-go/config/files"
	"golang.org/x/oauth2"
)

// cache is used to (un-)marshal the cached tokens from/to JSON.
type cache struct {
	// Login contains the cached tokens for the interactive login session
	Login *cacheEntry `json:"login,omitempty"`

	// ServicePrincipals contains cached tokens for service principals. The key is the service principal key's client_id.
	ServicePrincipals map[string]cacheEntry `json:"service-principals,omitempty"`

	// Workloads contains cached tokens for workload identity providers. The key is the workload identity provider's
	// resource name.
	Workloads map[string]cacheEntry `json:"workloads,omitempty"`
}

// readCache will read the cached tokens from a file. If an error occurs it will return an error and an empty cache
// struct with initialized maps to store service principal and workload credentials.
func readCache(cacheFile string) (*cache, error) {
	cachedTokens := &cache{
		Login:             nil,
		ServicePrincipals: map[string]cacheEntry{},
		Workloads:         map[string]cacheEntry{},
	}

	// Read the cache information from the file, if it exists
	cacheJSON, err := os.ReadFile(cacheFile)
	if err != nil {
		// If the file does not exist just return an empty cache
		if errors.Is(err, os.ErrNotExist) {
			return cachedTokens, nil
		}

		return cachedTokens, fmt.Errorf("failed to read credentials cache from %q: %w", cacheFile, err)
	}

	// Unmarshal the cached credentials
	if err = json.Unmarshal(cacheJSON, cachedTokens); err != nil {
		return cachedTokens, fmt.Errorf("failed to unmarshal cached credentials: %w", err)
	}

	return cachedTokens, nil
}

// write will write the cache information to the cache file.
// This currently overwrites the whole file and has the risk of overwriting concurrent updates, which would result in
// additional token refreshes or oauth flows. This risk is accepted for now, but we might want to address it in the
// future (e.g. by -re-reading the content before writing, locking the file or writing cache information to multiple files).
func (c *cache) write(cacheFile string) error {
	// Marshal the new tokens
	cacheJSON, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal cached tokens: %w", err)
	}

	// Create the directory if it doesn't exist yet
	err = os.MkdirAll(path.Dir(cacheFile), files.FolderMode)
	if err != nil {
		return fmt.Errorf("failed to create credential cache directory, %w", err)
	}

	// Write the file
	err = os.WriteFile(cacheFile, cacheJSON, files.FileMode)
	if err != nil {
		return fmt.Errorf("failed to write cached credentials to file: %w", err)
	}

	return nil
}

// removeExpiredTokens will remove any expired service principal or workload tokens that don't have a refresh
// token and are expired. This function should be called to garbage collect old tokens.
//
// Expired tokens will only get removed from the in-memory version of the cache, the result still has to explicitly get
// persisted to disc by calling write().
func (c *cache) removeExpiredTokens() {
	// Remove expired service principal tokens
	for identifier, entry := range c.ServicePrincipals {
		if entry.RefreshToken == "" && entry.AccessTokenExpiry.Before(time.Now()) {
			delete(c.ServicePrincipals, identifier)
		}
	}

	// Remove expired workload tokens
	for identifier, entry := range c.Workloads {
		if entry.RefreshToken == "" && entry.AccessTokenExpiry.Before(time.Now()) {
			delete(c.Workloads, identifier)
		}
	}
}

// cacheEntry represents an individual set of cached tokens.
type cacheEntry struct {
	// AccessToken is the bearer token used to authenticate to HCP.
	AccessToken string `json:"access_token,omitempty"`

	// RefreshToken is used to get a new access token.
	RefreshToken string `json:"refresh_token,omitempty"`

	// AccessTokenExpiry is when the access token will expire.
	AccessTokenExpiry time.Time `json:"access_token_expiry,omitempty"`
}

// token will convert the cacheEntry to an oauth2.Token.
func (entry *cacheEntry) token() *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  entry.AccessToken,
		RefreshToken: entry.RefreshToken,
		Expiry:       entry.AccessTokenExpiry,
	}
}

// cacheEntryFromToken will convert an oauth2.Token to a cacheEntry
func cacheEntryFromToken(token *oauth2.Token) *cacheEntry {
	return &cacheEntry{
		AccessToken:       token.AccessToken,
		AccessTokenExpiry: token.Expiry,
		RefreshToken:      token.RefreshToken,
	}
}
