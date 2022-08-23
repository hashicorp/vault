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
		newSeal := NewDefaultSeal(&seal.Access{
			Wrapper: aeadwrapper.NewShamirWrapper(),
		})
		// Need StoredShares set or this will look like a legacy shamir seal.
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    1,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	case seal.StoredKeysNotSupported:
		newSeal := NewDefaultSeal(&seal.Access{
			Wrapper: aeadwrapper.NewShamirWrapper(),
		})
		newSeal.SetCachedBarrierConfig(&SealConfig{
			StoredShares:    0,
			SecretThreshold: 1,
			SecretShares:    1,
		})
		return newSeal
	default:
		seal, err := NewAutoSeal(seal.NewTestSeal(opts))
		if err != nil {
			t.Fatal(err)
		}
		return seal
	}
}
