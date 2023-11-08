// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package keymanager

import (
	"context"

	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

const (
	KeyID = "root"
)

type KeyManager interface {
	// Returns a wrapping.Wrapper which can be used to perform key-related operations.
	Wrapper() wrapping.Wrapper
	// RetrievalToken is the material returned which can be used to source back the
	// encryption key. Depending on the implementation, the token can be the
	// encryption key itself or a token/identifier used to exchange the token.
	RetrievalToken(ctx context.Context) ([]byte, error)
}
