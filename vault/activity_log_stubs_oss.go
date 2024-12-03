// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"errors"
)

//go:generate go run github.com/hashicorp/vault/tools/stubmaker

// ErrEmptyResponse error is used to avoid returning "nil, nil" from a function
var ErrEmptyResponse = errors.New("empty response; the system encountered a statement that exclusively returns nil values")

// sendGlobalClients is a no-op on CE
func (a *ActivityLog) sendGlobalClients(ctx context.Context) error {
	return nil
}

// waitForSecondaryGlobalClients is a no-op on CE
func (a *ActivityLog) waitForSecondaryGlobalClients(ctx context.Context) error {
	return nil
}
