// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

// FileCredentialSource sources credentials by reading the file at the given
// path.
type FileCredentialSource struct {
	// Path sources the external credential by reading the value from the
	// specified file path.
	Path string `json:"path,omitempty"`

	// CredentialFormat configures how the credentials are extracted from the file.
	CredentialFormat
}

// Validate validates the config.
func (fc *FileCredentialSource) Validate() error {
	if fc.Path == "" {
		return fmt.Errorf("path must be set")
	}

	return fc.CredentialFormat.Validate()
}

// token retrieves the token from the specified file
func (fc *FileCredentialSource) token() (string, error) {
	credFile, err := os.Open(fc.Path)
	if err != nil {
		return "", fmt.Errorf("failed to open credential file %q", fc.Path)
	}
	defer credFile.Close()

	// Read the file but limit the size we read to a MB
	credBytes, err := io.ReadAll(io.LimitReader(credFile, 1<<20))
	if err != nil {
		return "", fmt.Errorf("failed to read credential file: %v", err)
	}

	if len(credBytes) == 0 {
		return "", fmt.Errorf("credential file is empty")
	}

	value := bytes.TrimSpace(credBytes)
	return fc.CredentialFormat.get(value)
}
