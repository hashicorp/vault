// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cliconfig

import (
	"github.com/hashicorp/vault/api/tokenhelper"
)

// DefaultTokenHelper returns the token helper that is configured for Vault.
// This helper should only be used for non-server CLI commands.
func DefaultTokenHelper() (tokenhelper.TokenHelper, error) {
	config, _, err := DefaultTokenHelperCheckDuplicates()
	return config, err
}

// TODO (HCL_DUP_KEYS_DEPRECATION): eventually make this consider duplicates an error. Ideally we should remove it but
// maybe we can't since it's become part of the API pkg.
func DefaultTokenHelperCheckDuplicates() (helper tokenhelper.TokenHelper, duplicate bool, err error) {
	config, duplicate, err := loadConfig("")
	if err != nil {
		return nil, duplicate, err
	}

	path := config.TokenHelper
	if path == "" {
		helper, err = tokenhelper.NewInternalTokenHelper()
		return helper, duplicate, err
	}

	path, err = tokenhelper.ExternalTokenHelperPath(path)
	if err != nil {
		return nil, duplicate, err
	}
	return &tokenhelper.ExternalTokenHelper{BinaryPath: path}, duplicate, nil
}
