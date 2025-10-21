// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import (
	"context"
)

// possiblyForwardPendingLoginMFASecretWrite attempts to persist the pending MFA secret and key to storage.
func possiblyForwardPendingLoginMFASecretWrite(ctx context.Context, c *Core, entityID string, pendingSecret *selfEnrollmentPendingMFASecret) error {
	return c.writeTOTPMFASecretAndKey(ctx, entityID, pendingSecret)
}
