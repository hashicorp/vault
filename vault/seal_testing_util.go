// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	testing "testing"

	"github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/vault/seal"
)

// NewTestSeal creates a new seal for testing. If you want to use the same seal multiple times, such as for
// a cluster, use NewTestSealFunc instead.
func NewTestSeal(t testing.TB, opts *seal.TestSealOpts) Seal {
	t.Helper()
	opts = seal.NewTestSealOpts(opts)
	logger := corehelpers.NewTestLogger(t).Named("sealAccess")

	switch opts.StoredKeys {
	case seal.StoredKeysSupportedShamirRoot:
		sealAccess, err := seal.NewAccessFromSealWrappers(logger, opts.Generation, true, []*seal.SealWrapper{
			seal.NewSealWrapper(aead.NewShamirWrapper(), 1, "shamir", "shamir", false, true),
		})
		if err != nil {
			t.Fatal("error creating test seal", err)
		}
		newSeal := NewDefaultSeal(sealAccess)
		// Need StoredShares set or this will look like a legacy shamir seal.
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    1,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	case seal.StoredKeysNotSupported:
		sealAccess, err := seal.NewAccessFromSealWrappers(logger, opts.Generation, true, []*seal.SealWrapper{
			seal.NewSealWrapper(aead.NewShamirWrapper(), 1, "shamir", "shamir", false, true),
		})
		if err != nil {
			t.Fatal("error creating test seal", err)
		}
		newSeal := NewDefaultSeal(sealAccess)
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    0,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	default:
		access, _ := seal.NewTestSeal(opts)
		return NewAutoSeal(access)
	}
}

// NewTestSealFunc returns a function that creates seals. All such seals will have TestWrappers that
// share the same secret, thus making them equivalent.
func NewTestSealFunc(t testing.TB, opts *seal.TestSealOpts) func() Seal {
	testSeal := NewTestSeal(t, opts)

	return func() Seal {
		return cloneTestSeal(t, testSeal)
	}
}

// CloneTestSeal creates a new test seal that shares the same seal wrappers as `testSeal`.
func cloneTestSeal(t testing.TB, testSeal Seal) Seal {
	logger := corehelpers.NewTestLogger(t).Named("sealAccess")

	access, err := seal.NewAccessFromSealWrappers(logger, testSeal.GetAccess().Generation(), testSeal.GetAccess().GetSealGenerationInfo().IsRewrapped(), testSeal.GetAccess().GetAllSealWrappersByPriority())
	if err != nil {
		t.Fatalf("error cloning seal %v", err)
	}
	if testSeal.StoredKeysSupported() == seal.StoredKeysNotSupported {
		return NewDefaultSeal(access)
	}
	return NewAutoSeal(access)
}
