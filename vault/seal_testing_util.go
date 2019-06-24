// +build !enterprise

package vault

import (
	"github.com/hashicorp/go-hclog"
	shamirseal "github.com/hashicorp/vault/vault/seal/shamir"
	testing "github.com/mitchellh/go-testing-interface"
)

func NewTestSeal(t testing.T, opts *TestSealOpts) Seal {
	var logger hclog.Logger
	if opts != nil {
		logger = opts.Logger
	}
	return NewDefaultSeal(shamirseal.NewSeal(logger))
}
