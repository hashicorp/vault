// +build !enterprise

package vault

import (
	shamirseal "github.com/hashicorp/vault/vault/seal/shamir"
	testing "github.com/mitchellh/go-testing-interface"
)

func NewTestSeal(t testing.T, opts *TestSealOpts) Seal {
	return NewDefaultSeal(shamirseal.NewSeal(opts.Logger))
}
