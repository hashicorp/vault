// Copyright IBM Corp. 2016, 2025
// SPDX-License-Identifier: MPL-2.0

package cryptoutil

import "testing"

func TestBlake2b256Hash(t *testing.T) {
	hashVal := Blake2b256Hash("sampletext")

	if string(hashVal) == "" || string(hashVal) == "sampletext" {
		t.Fatalf("failed to hash the text")
	}
}
