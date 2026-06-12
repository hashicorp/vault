// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: BUSL-1.1

//go:build !enterprise

package transit

import "github.com/hashicorp/vault/sdk/helper/keysutil"

func getFormattedPQCPublicKey(keyType keysutil.KeyType, entry keysutil.KeyEntry) string {
	return ""
}
