// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/hcp-sdk-go/auth/workload"
	"github.com/hashicorp/hcp-sdk-go/config/files"
)

const (
	// EnvHCPCredFile is the environment variable that sets the HCP Credential
	// File location.
	EnvHCPCredFile = "HCP_CRED_FILE"

	// CredentialFileName is the file name for the HCP credential file.
	CredentialFileName = "cred_file.json"

	// CredentialFileSchemeServicePrincipal is the credential file scheme value
	// that indicates service principal credentials should be used to
	// authenticate to HCP.
	CredentialFileSchemeServicePrincipal = "service_principal_creds"

	// CredentialFileSchemeWorkload is the credential file scheme value
	// that indicates workload identity credentials should be used to
	// authenticate to HCP.
	CredentialFileSchemeWorkload = "workload"
)

var (
	// testDefaultHCPCredFilePath is the default HCP Credential File location during
	// tests. The test should set its value.
	testDefaultHCPCredFilePath = ""
)

// CredentialFile stores information required to authenticate to HCP APIs. It
// supports various authentication schemes, such as service principal
type CredentialFile struct {
	// ProjectID captures the project ID of the service principal. It may be blank.
	ProjectID string `json:"project_id,omitempty"`

	// Scheme is the authentication scheme. It may be one of: service_principal_creds, workload.
	Scheme string `json:"scheme,omitempty"`

	// Workload configures the workload identity provider to exchange tokens
	// with.
	Workload *workload.IdentityProviderConfig `json:"workload,omitempty"`

	// Oauth configures authentication via Oauth.
	Oauth *OauthConfig `json:"oauth,omitempty"`
}

// OauthConfig configures authentication based on OAuth credentials.
type OauthConfig struct {
	// ClientID is the client id of an HCP Service Principal
	ClientID string `json:"client_id,omitempty"`

	// ClientSecret is the client secret of an HCP Service Principal
	ClientSecret string `json:"client_secret,omitempty"`
}

// Validate validates the CredentialFile
func (c *CredentialFile) Validate() error {
	if c == nil {
		return nil
	}

	if c.Scheme == CredentialFileSchemeServicePrincipal {
		if c.Oauth == nil {
			return fmt.Errorf("oauth config must be set when scheme is %q", CredentialFileSchemeServicePrincipal)
		}

		if err := c.Oauth.Validate(); err != nil {
			return fmt.Errorf("oauth: %v", err)
		}
	} else if c.Scheme == CredentialFileSchemeWorkload {
		if c.Workload == nil {
			return fmt.Errorf("workload config must be set when scheme is %q", CredentialFileSchemeWorkload)
		}

		if err := c.Workload.Validate(); err != nil {
			return fmt.Errorf("workload: %v", err)
		}
	} else {
		return fmt.Errorf("scheme must be one of: %q, %q", CredentialFileSchemeServicePrincipal, CredentialFileSchemeWorkload)
	}

	if c.Workload != nil && c.Oauth != nil {
		return fmt.Errorf("only one of oauth or workload may be set")
	}

	return nil
}

// Validate validates the OauthConfig
func (o *OauthConfig) Validate() error {
	if o.ClientID == "" || o.ClientSecret == "" {
		return fmt.Errorf("both client_id and client_secret must be set")
	}

	return nil
}

// ReadCredentialFile returns the credential file at the given path.
func ReadCredentialFile(path string) (*CredentialFile, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read credential file: %w", err)
	}

	var f CredentialFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("failed to unmarshal credential file: %v", err)
	}

	return &f, f.Validate()

}

// GetDefaultCredentialFile returns the credential file by searching the default
// credential file location or by using the credential file environment variable
// to look for an override. If no credential file is found, a nil value will be
// returned with no error set.
func GetDefaultCredentialFile() (*CredentialFile, error) {
	p, err := GetCredentialFilePath()
	if err != nil {
		return nil, fmt.Errorf("failed to find credential file: %v", err)
	}

	// Read the credential file, but if no credential file is found, suppress
	// the erorr.
	cf, err := ReadCredentialFile(p)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}

	return cf, err
}

// GetCredentialFilePath returns the credential file path, first looking for an
// overriding environment variable and then falling back to the default file
// location.
func GetCredentialFilePath() (string, error) {
	if testDefaultHCPCredFilePath != "" {
		return testDefaultHCPCredFilePath, nil
	}

	if p, ok := os.LookupEnv(EnvHCPCredFile); ok {
		return p, nil
	}

	// Get the user's home directory.
	userHome, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve user's home directory path: %v", err)
	}

	p := filepath.Join(userHome, files.DefaultDirectory, CredentialFileName)
	return p, nil
}

// WriteDefaultCredentialFile writes the credential file to the default
// credential file location or to the value of EnvHCPCredFile if set.
func WriteDefaultCredentialFile(cf *CredentialFile) error {
	p, err := GetCredentialFilePath()
	if err != nil {
		return err
	}

	return WriteCredentialFile(p, cf)
}

// WriteCredentialFile writes the given credential file to the path.
func WriteCredentialFile(path string, cf *CredentialFile) error {
	data, err := json.MarshalIndent(cf, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, files.FileMode)
}
