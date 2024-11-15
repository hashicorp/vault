// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package workload

import (
	"fmt"
	"os"
)

// EnvironmentVariableCredentialSource sources credentials by reading the
// specified environment variable.
type EnvironmentVariableCredentialSource struct {
	// Var sources the external credential value from the given environment variable.
	Var string `json:"var,omitempty"`

	// CredentialFormat configures how the credentials are extracted from the environment
	// variable value.
	CredentialFormat
}

// Validate validates the config.
func (ec *EnvironmentVariableCredentialSource) Validate() error {
	if ec.Var == "" {
		return fmt.Errorf("environment variable must be specified")
	}

	return ec.CredentialFormat.Validate()
}

// token retrieves the token from the environment variable
func (ec *EnvironmentVariableCredentialSource) token() (string, error) {
	value, ok := os.LookupEnv(ec.Var)
	if !ok {
		return "", fmt.Errorf("environment variable not found")
	}
	if value == "" {
		return "", fmt.Errorf("environment variable value is empty")
	}

	return ec.CredentialFormat.get([]byte(value))
}
