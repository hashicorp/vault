// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package managed_key

import (
	"context"
	"crypto"
	"errors"
	"io"
)

func GetManagedKeyInfo(ctx context.Context, mkv SSHManagedKeyView, keyId managedKeyId) (*ManagedKeyInfo, error) {
	return nil, errors.New("managed keys are supported within enterprise edition only")
}

func Sign(ctx context.Context, mkv SSHManagedKeyView, keyId managedKeyId, rand io.Reader, msg []byte, opts crypto.SignerOpts) ([]byte, error) {
	return nil, errors.New("managed keys are supported within enterprise edition only")
}
