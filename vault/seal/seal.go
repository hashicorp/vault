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
type Access struct {
	wrapping.Wrapper
	WrapperType wrapping.WrapperType
}

func (a *Access) Init(ctx context.Context) error {
	if initWrapper, ok := a.Wrapper.(wrapping.InitFinalizer); ok {
		return initWrapper.Init(ctx)
	}
	return nil
}

func (a *Access) SetType(t wrapping.WrapperType) {
	a.WrapperType = t
}

func (a *Access) Type(ctx context.Context) (wrapping.WrapperType, error) {
	if a != nil && a.WrapperType != "" {
		return a.WrapperType, nil
	}
	return a.Wrapper.Type(ctx)
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *Access) Encrypt(ctx context.Context, plaintext, aad []byte) (blob *wrapping.BlobInfo, err error) {
	wTyp, err := a.Wrapper.Type(ctx)
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

	return a.Wrapper.Encrypt(ctx, plaintext, wrapping.WithAad(aad))
}

// Decrypt uses the underlying seal to decrypt the cryptotext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
func (a *Access) Decrypt(ctx context.Context, data *wrapping.BlobInfo, aad []byte) (pt []byte, err error) {
	wTyp, err := a.Wrapper.Type(ctx)
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

	return a.Wrapper.Decrypt(ctx, data, wrapping.WithAad(aad))
}

func (a *Access) Finalize(ctx context.Context) error {
	if finalizeWrapper, ok := a.Wrapper.(wrapping.InitFinalizer); ok {
		return finalizeWrapper.Finalize(ctx)
	}
	return nil
}
