// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"fmt"
	"sync"

	"github.com/hashicorp/go-hclog"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	UUID "github.com/hashicorp/go-uuid"
	"github.com/hashicorp/vault/sdk/helper/logging"
)

type TestSealOpts struct {
	Logger       hclog.Logger
	StoredKeys   StoredKeysSupport
	Secrets      [][]byte
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
	switch len(opts.Secrets) {
	case opts.WrapperCount:
		// all good, each wrapper has its own secret

	case 0:
		if opts.WrapperCount == 1 {
			// If there is only one wrapper, the default TestWrapper behaviour of reversing
			// the bytes slice is fine.
			opts.Secrets = [][]byte{nil}
		} else {
			// If there is more than one wrapper, each one needs a different secret
			for i := 0; i < opts.WrapperCount; i++ {
				uuid, err := UUID.GenerateUUID()
				if err != nil {
					panic(fmt.Sprintf("error generating secret: %v", err))
				}
				opts.Secrets = append(opts.Secrets, []byte(uuid))
			}
		}

	default:
		panic(fmt.Sprintf("wrong number of secrets %d vs %d wrappers", len(opts.Secrets), opts.WrapperCount))
	}
	return opts
}

func NewTestSeal(opts *TestSealOpts) (Access, []*ToggleableWrapper) {
	opts = NewTestSealOpts(opts)
	wrappers := make([]*ToggleableWrapper, opts.WrapperCount)
	sealWrappers := make([]*SealWrapper, opts.WrapperCount)
	ctx := context.Background()
	for i := 0; i < opts.WrapperCount; i++ {
		wrapperName := fmt.Sprintf("%s-%d", opts.Name, i+1)
		wrappers[i] = &ToggleableWrapper{Wrapper: wrapping.NewTestWrapper(opts.Secrets[i])}
		_, err := wrappers[i].Wrapper.SetConfig(context.Background(), wrapping.WithKeyId(wrapperName))
		if err != nil {
			panic(err)
		}
		wrapperType, err := wrappers[i].Type(ctx)
		if err != nil {
			panic(err)
		}
		sealWrappers[i] = NewSealWrapper(
			wrappers[i],
			i+1,
			wrapperName,
			wrapperType.String(),
			false,
			true,
		)
	}

	sealAccess, err := NewAccessFromSealWrappers(nil, opts.Generation, true, sealWrappers)
	if err != nil {
		panic(err)
	}
	return sealAccess, wrappers
}

type TestSealWrapperOpts struct {
	Logger       hclog.Logger
	Secret       []byte
	Name         wrapping.WrapperType
	WrapperCount int
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

func (t *ToggleableWrapper) Decrypt(ctx context.Context, info *wrapping.BlobInfo, opts ...wrapping.Option) ([]byte, error) {
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
