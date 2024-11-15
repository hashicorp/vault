// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package files

import (
	"fmt"
	"os"
	"path"
)

const (
	// DefaultDirectory is the directory within the home folder in which configuration, credential and cache files will
	// be stored.
	DefaultDirectory = ".config/hcp"

	// FolderMode is the mode the configuration folder will be set to if it is automatically created.
	//
	// It is set to 755 as the .config folder might be used by other applications in the future and the names of
	// contained files are not sensitive.
	// The mode is not enforced and can be changed or set by pre-creating the folder.
	FolderMode = os.FileMode(0755)

	// FileMode is the mode the files in the configuration will be set to if they are automatically created.
	//
	// It is set to only be accessible by the user that executes the application as the file contains credentials.
	// The mode is not enforced and can be changed or set by pre-creating the file(s).
	FileMode = os.FileMode(0700)

	// TokenCacheFileName is the name of the cache file within the configuration folder.
	TokenCacheFileName = "creds-cache.json"
)

// TokenCacheFile will return the absolute path to the token cache file.
func TokenCacheFile() (string, error) {
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user's home directory path: %v", err)
	}

	cacheFile := path.Join(userHome, DefaultDirectory, TokenCacheFileName)
	return cacheFile, nil
}
