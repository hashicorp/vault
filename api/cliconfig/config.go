// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cliconfig

import (
	"fmt"
	"os"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl"
	"github.com/hashicorp/hcl/hcl/ast"
	"github.com/mitchellh/go-homedir"
)

const (
	// defaultConfigPath is the default path to the configuration file
	defaultConfigPath = "~/.vault"

	// configPathEnv is the environment variable that can be used to
	// override where the Vault configuration is.
	configPathEnv = "VAULT_CONFIG_PATH"
)

// Config is the CLI configuration for Vault that can be specified via
// a `$HOME/.vault` file which is HCL-formatted (therefore HCL or JSON).
type defaultConfig struct {
	// TokenHelper is the executable/command that is executed for storing
	// and retrieving the authentication token for the Vault CLI. If this
	// is not specified, then vault's internal token store will be used, which
	// stores the token on disk unencrypted.
	TokenHelper string `hcl:"token_helper"`
}

// loadConfig reads the configuration from the given path. If path is
// empty, then the default path will be used, or the environment variable
// if set.
func loadConfig(path string) (config *defaultConfig, duplicate bool, err error) {
	if path == "" {
		path = defaultConfigPath
	}
	if v := os.Getenv(configPathEnv); v != "" {
		path = v
	}

	// NOTE: requires HOME env var to be set
	path, err = homedir.Expand(path)
	if err != nil {
		return nil, false, fmt.Errorf("error expanding config path %q: %w", path, err)
	}

	contents, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, false, err
	}

	conf, duplicate, err := parseConfig(string(contents))
	if err != nil {
		return nil, duplicate, fmt.Errorf("error parsing config file at %q: %w; ensure that the file is valid; Ansible Vault is known to conflict with it", path, err)
	}

	return conf, duplicate, nil
}

// parseConfig parses the given configuration as a string.
func parseConfig(contents string) (config *defaultConfig, duplicate bool, err error) {
	// TODO (HCL_DUP_KEYS_DEPRECATION): on removal stage change this to a simple hcl.Parse, effectively treating
	// duplicate keys as an error. Also get rid of all of these "duplicate" named return values
	root, duplicate, err := parseAndCheckForDuplicateHclAttributes(contents)
	if err != nil {
		return nil, duplicate, err
	}

	// Top-level item should be the object list
	list, ok := root.Node.(*ast.ObjectList)
	if !ok {
		return nil, duplicate, fmt.Errorf("failed to parse config; does not contain a root object")
	}

	valid := map[string]struct{}{
		"token_helper": {},
	}

	var validationErrors error
	for _, item := range list.Items {
		key := item.Keys[0].Token.Value().(string)
		if _, ok := valid[key]; !ok {
			validationErrors = multierror.Append(validationErrors, fmt.Errorf("invalid key %q on line %d", key, item.Assign.Line))
		}
	}

	if validationErrors != nil {
		return nil, duplicate, validationErrors
	}

	var c defaultConfig
	if err := hcl.DecodeObject(&c, list); err != nil {
		return nil, duplicate, err
	}
	return &c, duplicate, nil
}
