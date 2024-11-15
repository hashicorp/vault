// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokencache

// ClearLoginCache will clear the login entry in the passed cache file. This can be used to "logout" principal.
func ClearLoginCache(cacheFile string) error {
	// Read cached tokens
	tokenCache, err := readCache(cacheFile)
	if err != nil {
		return err
	}

	// Unset the login entry
	tokenCache.Login = nil

	// Write cached without login entry
	if err := tokenCache.write(cacheFile); err != nil {
		return err
	}

	return nil
}
