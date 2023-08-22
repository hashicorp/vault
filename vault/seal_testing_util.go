// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package vault

import (
	"github.com/hashicorp/go-hclog"
	aeadwrapper "github.com/hashicorp/go-kms-wrapping/wrappers/aead/v2"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault/seal"
	testing "github.com/mitchellh/go-testing-interface"
)

func NewTestSeal(t testing.T, opts *seal.TestSealOpts) Seal {
	t.Helper()
	if opts == nil {
		opts = &seal.TestSealOpts{}
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(hclog.Debug)
	}

	switch opts.StoredKeys {
	case seal.StoredKeysSupportedShamirRoot:
		sealAccess := seal.NewAccess([]seal.SealInfo{
			{
				Wrapper:  aeadwrapper.NewShamirWrapper(),
				Priority: 1,
				Name:     "shamir",
			},
		})
		newSeal := NewDefaultSeal(sealAccess)
		// Need StoredShares set or this will look like a legacy shamir seal.
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    1,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	case seal.StoredKeysNotSupported:
		sealAccess := seal.NewAccess([]seal.SealInfo{
			{
				Wrapper:  aeadwrapper.NewShamirWrapper(),
				Priority: 1,
				Name:     "shamir",
			},
		})
		newSeal := NewDefaultSeal(sealAccess)
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    0,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	default:
		access, _ := seal.NewTestSeal(opts)
		seal, err := NewAutoSeal(access)
		if err != nil {
			t.Fatal(err)
		}
		return seal
	}
}
