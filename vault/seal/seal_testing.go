package seal

import (
	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

type TestSealOpts struct {
	Logger     hclog.Logger
	StoredKeys StoredKeysSupport
	Secret     []byte
	Name       string
}

func NewTestSeal(opts *TestSealOpts) *Access {
	if opts == nil {
		opts = new(TestSealOpts)
	}

	return &Access{
		Wrapper:        wrapping.NewTestWrapper(opts.Secret),
		OverriddenType: opts.Name,
	}
}
