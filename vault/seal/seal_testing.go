// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

type TestSealOpts struct {
	Logger       hclog.Logger
	StoredKeys   StoredKeysSupport
	Secret       []byte
	Name         wrapping.WrapperType
	WrapperCount int
}

func NewTestSeal(opts *TestSealOpts) (Access, []*ToggleableWrapper) {
	if opts == nil {
		opts = new(TestSealOpts)
	}
	if opts.WrapperCount == 0 {
		opts.WrapperCount = 1
	}

	wrappers := make([]*ToggleableWrapper, opts.WrapperCount)
	sealInfos := make([]SealInfo, opts.WrapperCount)
	for i := 0; i < opts.WrapperCount; i++ {
		wrappers[i] = &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
		sealInfos[i] = SealInfo{
			Wrapper:  wrappers[i],
			Priority: i + 1,
			Name:     fmt.Sprintf("%s-%d", opts.Name, i+1),
		}
	}

	sealAccess := NewAccess(sealInfos)
	return sealAccess, wrappers
}

func NewToggleableTestSeal(opts *TestSealOpts) (Access, []func(error)) {
	if opts == nil {
		opts = new(TestSealOpts)
	}
	if opts.WrapperCount == 0 {
		opts.WrapperCount = 1
	}

	wrappers := make([]*ToggleableWrapper, opts.WrapperCount)
	sealInfos := make([]SealInfo, opts.WrapperCount)
	funcs := make([]func(error), opts.WrapperCount)
	for i := 0; i < opts.WrapperCount; i++ {
		w := &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secret)}
		wrappers[i] = w
		sealInfos[i] = SealInfo{
			Wrapper:  wrappers[i],
			Priority: i + 1,
			Name:     fmt.Sprintf("%s-%d", opts.Name, i+1),
		}
		funcs[i] = w.SetError
	}

	sealAccess := NewAccess(sealInfos)

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
