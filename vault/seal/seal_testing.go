package seal

import (
	"context"
	"sync"

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

func NewToggleableTestSeal(opts *TestSealOpts) (*Access, func(error)) {
	if opts == nil {
		opts = new(TestSealOpts)
	}

	w := &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
	return &Access{
		Wrapper:        w,
		OverriddenType: opts.Name,
	}, w.SetError
}

type ToggleableWrapper struct {
	wrapping.Wrapper
	error error
	l     sync.RWMutex
}

func (t *ToggleableWrapper) Encrypt(ctx context.Context, bytes []byte, bytes2 []byte) (*wrapping.EncryptedBlobInfo, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if t.error != nil {
		return nil, t.error
	}
	return t.Wrapper.Encrypt(ctx, bytes, bytes2)
}

func (t ToggleableWrapper) Decrypt(ctx context.Context, info *wrapping.EncryptedBlobInfo, bytes []byte) ([]byte, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if t.error != nil {
		return nil, t.error
	}
	return t.Wrapper.Decrypt(ctx, info, bytes)
}

func (t *ToggleableWrapper) SetError(err error) {
	t.l.Lock()
	defer t.l.Unlock()
	t.error = err
}

var _ wrapping.Wrapper = &ToggleableWrapper{}
