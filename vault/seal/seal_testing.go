package seal

import (
	"context"

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

func NewToggleableTestSeal(opts *TestSealOpts) (*Access, *error) {
	if opts == nil {
		opts = new(TestSealOpts)
	}

	w := &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
	return &Access{
		Wrapper:        w,
		OverriddenType: opts.Name,
	}, &w.Error
}

type ToggleableWrapper struct {
	wrapping.Wrapper
	Error error
}

func (t ToggleableWrapper) Encrypt(ctx context.Context, bytes []byte, bytes2 []byte) (*wrapping.EncryptedBlobInfo, error) {
	if t.Error != nil {
		return nil, t.Error
	}
	return t.Wrapper.Encrypt(ctx, bytes, bytes2)
}

func (t ToggleableWrapper) Decrypt(ctx context.Context, info *wrapping.EncryptedBlobInfo, bytes []byte) ([]byte, error) {
	if t.Error != nil {
		return nil, t.Error
	}
	return t.Wrapper.Decrypt(ctx, info, bytes)
}

var _ wrapping.Wrapper = &ToggleableWrapper{}
