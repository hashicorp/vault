// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package pki

import (
	"github.com/hashicorp/vault/builtin/logical/pki/issuing"
)

func (b *backend) adjustInputBundle(input *inputBundle) {}

func entValidateRole(b *backend, entry *issuing.RoleEntry, operation string) ([]string, error) {
	return nil, nil
}
