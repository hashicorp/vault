// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package seal

import (
	"context"
	"time"

	metrics "github.com/armon/go-metrics"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
)

type StoredKeysSupport int

const (
	// The 0 value of StoredKeysSupport is an invalid option
	StoredKeysInvalid StoredKeysSupport = iota
	StoredKeysNotSupported
	StoredKeysSupportedGeneric
	StoredKeysSupportedShamirRoot
)

func (s StoredKeysSupport) String() string {
	switch s {
	case StoredKeysNotSupported:
		return "Old-style Shamir"
	case StoredKeysSupportedGeneric:
		return "AutoUnseal"
	case StoredKeysSupportedShamirRoot:
		return "New-style Shamir"
	default:
		return "Invalid StoredKeys type"
	}
}

// Access is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access interface {
	wrapping.Wrapper
	wrapping.InitFinalizer

	GetWrapper() wrapping.Wrapper
}

type access struct {
	w wrapping.Wrapper
}

var _ Access = (*access)(nil)

func NewAccess(w wrapping.Wrapper) Access {
	return &access{
		w: w,
	}
}

func (a *access) KeyId(ctx context.Context) (string, error) {
	return a.w.KeyId(ctx)
}

func (a *access) SetConfig(ctx context.Context, options ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	return a.w.SetConfig(ctx, options...)
}

func (a *access) GetWrapper() wrapping.Wrapper {
	return a.w
}

func (a *access) Init(ctx context.Context, options ...wrapping.Option) error {
	if initWrapper, ok := a.w.(wrapping.InitFinalizer); ok {
		return initWrapper.Init(ctx, options...)
	}
	return nil
}

func (a *access) Type(ctx context.Context) (wrapping.WrapperType, error) {
	return a.w.Type(ctx)
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *access) Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (blob *wrapping.BlobInfo, err error) {
	wTyp, err := a.w.Type(ctx)
	if err != nil {
		return nil, err
	}

	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", wTyp.String(), "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", wTyp.String(), "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", wTyp.String(), "encrypt"}, 1)

	return a.w.Encrypt(ctx, plaintext, options...)
}

// Decrypt uses the underlying seal to decrypt the cryptotext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
func (a *access) Decrypt(ctx context.Context, data *wrapping.BlobInfo, options ...wrapping.Option) (pt []byte, err error) {
	wTyp, err := a.w.Type(ctx)
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", wTyp.String(), "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", wTyp.String(), "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", wTyp.String(), "decrypt"}, 1)

	return a.w.Decrypt(ctx, data, options...)
}

func (a *access) Finalize(ctx context.Context, options ...wrapping.Option) error {
	if finalizeWrapper, ok := a.w.(wrapping.InitFinalizer); ok {
		return finalizeWrapper.Finalize(ctx, options...)
	}
	return nil
}
