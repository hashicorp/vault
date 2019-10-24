package vault

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/sdk/helper/logging"
	"github.com/hashicorp/vault/vault/seal"
	shamirseal "github.com/hashicorp/vault/vault/seal/shamir"
	testing "github.com/mitchellh/go-testing-interface"
)

func NewTestSeal(t testing.T, opts *TestSealOpts) Seal {
	t.Helper()
	if opts == nil {
		opts = &TestSealOpts{}
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(hclog.Debug)
	}

	switch opts.StoredKeys {
	case StoredKeysSupportedShamirMaster:
		newSeal := NewDefaultSeal(shamirseal.NewSeal(opts.Logger))
		// Need StoredShares set or this will look like a legacy shamir seal.
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    1,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	case StoredKeysNotSupported:
		newSeal := NewDefaultSeal(shamirseal.NewSeal(opts.Logger))
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    0,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	default:
		return NewAutoSeal(seal.NewTestSeal(opts.Secret))
	}
}
