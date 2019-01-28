// +build !enterprise

package vault

import testing "github.com/mitchellh/go-testing-interface"

func NewTestSeal(testing.T, *TestSealOpts) Seal {
	return NewDefaultSeal()
}
