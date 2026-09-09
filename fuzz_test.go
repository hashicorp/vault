//go:build go1.18
// +build go1.18

// Copyright IBM Corp. 2016, 2026
// SPDX-License-Identifier: MPL-2.0

package vault_test

import (
	"bytes"
	"testing"

	"github.com/hashicorp/vault/shamir"
)

// FuzzShamirCombineFlat tests Shamir Combine with arbitrary
// byte data interpreted as two shares of equal length.
//
// This is the pre-auth boundary for Vault's unseal process.
// The Combine function reconstructs a secret from shares
// provided by operators. Malformed shares must not panic.
//
// 53 CVEs exist for HashiCorp Vault — this code is part of
// the core unseal mechanism that protects every secret.
func FuzzShamirCombineFlat(f *testing.F) {
	// Seed: valid 2-of-3 shares
	secret := []byte("vault-secret-data-for-fuzzing")
	shares, err := shamir.Split(secret, 3, 2)
	if err == nil {
		f.Add(shares[0], shares[1])
	}

	f.Add([]byte{}, []byte{})
	f.Add([]byte{0x00}, []byte{0x00})
	f.Add([]byte{0x00, 0x01}, []byte{0x01, 0x02})
	f.Add([]byte{0xFF}, []byte{0xFF})

	f.Fuzz(func(t *testing.T, part1, part2 []byte) {
		if len(part1) > 10000 || len(part2) > 10000 {
			return
		}

		// Pad shorter to match longer
		maxLen := len(part1)
		if len(part2) > maxLen {
			maxLen = len(part2)
		}
		p1 := make([]byte, maxLen)
		p2 := make([]byte, maxLen)
		copy(p1, part1)
		copy(p2, part2)

		// Combine must never panic on any input
		_, _ = shamir.Combine([][]byte{p1, p2})
	})
}

// FuzzShamirCombineTriple tests Combine with three shares.
func FuzzShamirCombineTriple(f *testing.F) {
	secret := []byte("vault-secret-data")
	shares, err := shamir.Split(secret, 5, 3)
	if err == nil {
		f.Add(shares[0], shares[1], shares[2])
	}

	f.Add([]byte{0x00}, []byte{0x01}, []byte{0x02})
	f.Add([]byte{}, []byte{}, []byte{})

	f.Fuzz(func(t *testing.T, p1, p2, p3 []byte) {
		if len(p1) > 10000 || len(p2) > 10000 || len(p3) > 10000 {
			return
		}
		maxLen := len(p1)
		for _, l := range []int{len(p2), len(p3)} {
			if l > maxLen {
				maxLen = l
			}
		}
		a := make([]byte, maxLen)
		b := make([]byte, maxLen)
		c := make([]byte, maxLen)
		copy(a, p1)
		copy(b, p2)
		copy(c, p3)

		_, _ = shamir.Combine([][]byte{a, b, c})
	})
}

// FuzzShamirSplitRoundTrip tests Split→Combine round-trip
// with arbitrary secrets and valid parameters.
func FuzzShamirSplitRoundTrip(f *testing.F) {
	f.Add([]byte("test-secret"), 5, 3)
	f.Add([]byte{0x00}, 2, 2)
	f.Add(make([]byte, 100), 10, 5)

	f.Fuzz(func(t *testing.T, secret []byte, parts, threshold int) {
		if len(secret) == 0 || len(secret) > 10000 {
			return
		}

		// Normalize to valid ranges
		parts = (parts%253 + 2)          // 2-254
		threshold = (threshold%parts + 1) // 1 to parts
		if threshold < 2 {
			threshold = 2
		}
		if parts < threshold {
			parts = threshold
		}

		shares, err := shamir.Split(secret, parts, threshold)
		if err != nil {
			return
		}

		// Verify round-trip: threshold shares recover secret
		combined, err := shamir.Combine(shares[:threshold])
		if err != nil {
			t.Errorf("Combine failed on valid shares: %v", err)
			return
		}
		if !bytes.Equal(combined, secret) {
			t.Errorf("Round-trip mismatch want=%x got=%x", secret, combined)
		}
	})
}

// FuzzShamirSplitEdgeCases tests Split with edge-case parameters.
func FuzzShamirSplitEdgeCases(f *testing.F) {
	f.Add([]byte("edge"), 255, 2)
	f.Add([]byte{0xFF, 0x00, 0xAA}, 100, 50)

	f.Fuzz(func(t *testing.T, secret []byte, parts, threshold int) {
		if len(secret) == 0 || len(secret) > 10000 {
			return
		}
		// Split must never panic on any parameter combination
		_, _ = shamir.Split(secret, parts, threshold)
	})
}
