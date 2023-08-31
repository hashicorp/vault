// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/vault/helper/testhelpers/corehelpers"
	"github.com/hashicorp/vault/vault/seal"
	testing "github.com/mitchellh/go-testing-interface"
)

func NewTestSeal(t testing.T, opts *seal.TestSealOpts) Seal {
	t.Helper()
	opts = seal.NewTestSealOpts(opts)
	logger := corehelpers.NewTestLogger(t).Named("sealAccess")

	switch opts.StoredKeys {
	case seal.StoredKeysSupportedShamirRoot:
		w := aeadwrapper.NewShamirWrapper()
		sealAccess, err := seal.NewAccessFromSealWrappers(logger, opts.Generation, true, []seal.SealWrapper{
			{
				Wrapper:  w,
				Priority: 1,
				Name:     "shamir",
			},
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
		w := aeadwrapper.NewShamirWrapper()
		sealAccess, err := seal.NewAccessFromSealWrappers(logger, opts.Generation, true, []seal.SealWrapper{
			{
				Wrapper:  w,
				Priority: 1,
				Name:     "shamir",
			},
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
