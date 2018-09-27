// +build !enterprise

package vault

import "github.com/mitchellh/go-testing-interface"

func NewTestSeal(testing.T, *TestSealOpts) Seal {
	return NewDefaultSeal()
}
