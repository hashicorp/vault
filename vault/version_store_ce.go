// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package vault

func IsOAuthJwt(token string) bool {
	return false
}

func IsOAuthJwtId(tokenID string) bool {
	return false
}
