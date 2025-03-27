// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package managed_key

import (
	"context"
	"crypto"
	"errors"
)

func GetManagedKeyPublicKey(ctx context.Context, mkv SSHManagedKeyView, keyId managedKeyId) (crypto.PublicKey, error) {
	return nil, errors.New("managed keys are supported within enterprise edition only")
}
