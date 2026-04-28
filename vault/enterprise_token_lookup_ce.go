// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

import "errors"

func resolveEnterpriseTokenIDForLookup(_ string) (string, error) {
	return "", errors.New("enterprise build required")
}
