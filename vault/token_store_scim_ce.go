// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault/sdk/logical"
)

// VerifySCIMTokenCreation is a no-op stub for CE builds.  Enterprise builds
// perform entity resolution and renewable/orphan guards when type=scim is
// requested.
func (ts *TokenStore) VerifySCIMTokenCreation(
	_ context.Context,
	_ *logical.Request,
	_ bool,
	_ bool,
) (string, time.Duration, int32, *logical.Response, error) {
	return "", 0, 0, nil, nil
}

// AcquireSCIMTokenSlot is a no-op stub for CE builds.  Enterprise builds
// acquire a per-client lock and verify the active-token cap before issuance.
func (ts *TokenStore) AcquireSCIMTokenSlot(
	_ context.Context,
	_ *logical.Request,
	_ string,
	_ int32,
) (func(), *logical.Response, error) {
	return func() {}, nil, nil
}
