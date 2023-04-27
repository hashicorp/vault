// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cryptoutil

import "golang.org/x/crypto/blake2b"

func Blake2b256Hash(key string) []byte {
	hf, _ := blake2b.New256(nil)

	hf.Write([]byte(key))

	return hf.Sum(nil)
}
