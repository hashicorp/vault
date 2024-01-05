// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package pki

import (
	"strings"

	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
)

func (sc *storageContext) isDefaultKeySet() (bool, error) {
	config, err := sc.getKeysConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultKeyId.String()) != "", nil
}

func (sc *storageContext) isDefaultIssuerSet() (bool, error) {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return false, err
	}

	return strings.TrimSpace(config.DefaultIssuerId.String()) != "", nil
}

func (sc *storageContext) updateDefaultKeyId(id issuing.KeyID) error {
	config, err := sc.getKeysConfig()
	if err != nil {
		return err
	}

	if config.DefaultKeyId != id {
		return sc.setKeysConfig(&issuing.KeyConfigEntry{
			DefaultKeyId: id,
		})
	}

	return nil
}

func (sc *storageContext) updateDefaultIssuerId(id issuing.IssuerID) error {
	config, err := sc.getIssuersConfig()
	if err != nil {
		return err
	}

	if config.DefaultIssuerId != id {
		config.DefaultIssuerId = id
		return sc.setIssuersConfig(config)
	}

	return nil
}
