// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/vault/sdk/helper/logging"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

type TestSealOpts struct {
	Logger       hclog.Logger
	StoredKeys   StoredKeysSupport
	Secret       []byte
	Name         wrapping.WrapperType
	WrapperCount int
	Generation   uint64
}

func NewTestSealOpts(opts *TestSealOpts) *TestSealOpts {
	if opts == nil {
		opts = new(TestSealOpts)
	}
	if opts.WrapperCount == 0 {
		opts.WrapperCount = 1
	}
	if opts.Logger == nil {
		opts.Logger = logging.NewVaultLogger(hclog.Debug)
	}
	if opts.Generation == 0 {
		// we might at some point need to allow Generation == 0
		opts.Generation = 1
	}
	return opts
}

func NewTestSeal(opts *TestSealOpts) (Access, []*ToggleableWrapper) {
	opts = NewTestSealOpts(opts)
	wrappers := make([]*ToggleableWrapper, opts.WrapperCount)
	sealWrappers := make([]SealWrapper, opts.WrapperCount)
	for i := 0; i < opts.WrapperCount; i++ {
		wrappers[i] = &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
		sealWrappers[i] = SealWrapper{
			Wrapper:  wrappers[i],
			Priority: i + 1,
			Name:     fmt.Sprintf("%s-%d", opts.Name, i+1),
		}
	}

	sealAccess, err := NewAccessFromSealWrappers(nil, opts.Generation, true, sealWrappers)
	if err != nil {
		panic(err)
	}
	return sealAccess, wrappers
}

func NewToggleableTestSeal(opts *TestSealOpts) (Access, []func(error)) {
	opts = NewTestSealOpts(opts)

	wrappers := make([]*ToggleableWrapper, opts.WrapperCount)
	sealWrappers := make([]SealWrapper, opts.WrapperCount)
	funcs := make([]func(error), opts.WrapperCount)
	for i := 0; i < opts.WrapperCount; i++ {
		w := &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
		wrappers[i] = w
		sealWrappers[i] = SealWrapper{
			Wrapper:  wrappers[i],
			Priority: i + 1,
			Name:     fmt.Sprintf("%s-%d", opts.Name, i+1),
		}
		funcs[i] = w.SetError
	}

	sealAccess, err := NewAccessFromSealWrappers(nil, opts.Generation, true, sealWrappers)
	if err != nil {
		panic(err)
	}

	return sealAccess, funcs
}

type ToggleableWrapper struct {
	wrapping.Wrapper
	wrapperType  *wrapping.WrapperType
	error        error
	encryptError error
	l            sync.RWMutex
}

func (t *ToggleableWrapper) Encrypt(ctx context.Context, bytes []byte, opts ...wrapping.Option) (*wrapping.BlobInfo, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if t.error != nil {
		return nil, t.error
	}
	if t.encryptError != nil {
		return nil, t.encryptError
	}
	return t.Wrapper.Encrypt(ctx, bytes, opts...)
}

func (t ToggleableWrapper) Decrypt(ctx context.Context, info *wrapping.BlobInfo, opts ...wrapping.Option) ([]byte, error) {
	t.l.RLock()
	defer t.l.RUnlock()
	if t.error != nil {
		return nil, t.error
	}
	return t.Wrapper.Decrypt(ctx, info, opts...)
}

func (t *ToggleableWrapper) SetError(err error) {
	t.l.Lock()
	defer t.l.Unlock()
	t.error = err
}

// An error only occuring on encrypt
func (t *ToggleableWrapper) SetEncryptError(err error) {
	t.l.Lock()
	defer t.l.Unlock()
	t.encryptError = err
}

func (t *ToggleableWrapper) Type(ctx context.Context) (wrapping.WrapperType, error) {
	if t.wrapperType != nil {
		return *t.wrapperType, nil
	}
	return t.Wrapper.Type(ctx)
}

var _ wrapping.Wrapper = &ToggleableWrapper{}
