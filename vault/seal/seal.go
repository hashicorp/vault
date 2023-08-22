// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package seal

import (
	"context"
	"errors"
	"sort"
	"time"

	metrics "github.com/armon/go-metrics"
	wrapping "github.com/hashicorp/go-kms-wrapping/v2"
	"github.com/hashicorp/go-kms-wrapping/v2/aead"
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

type SealInfo struct {
	wrapping.Wrapper
	Priority int
	Name     string
}

// Access is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access interface {
	wrapping.Wrapper
	wrapping.InitFinalizer

	GetWrapper() wrapping.Wrapper
	SetShamirSealKey([]byte) error
	GetShamirKeyBytes(ctx context.Context) ([]byte, error)
	SealType(ctx context.Context) (SealType, error)
}

type access struct {
	wrappersByPriority []SealInfo
}

var _ Access = (*access)(nil)

func NewAccess(w []SealInfo) Access {
	a := &access{
		wrappersByPriority: w,
	}

	sort.Slice(a.wrappersByPriority, func(i int, j int) bool { return a.wrappersByPriority[i].Priority < a.wrappersByPriority[j].Priority })

	return a
}

func (a *access) KeyId(ctx context.Context) (string, error) {
	w := a.getDefaultWrapper()
	if w != nil {
		return w.KeyId(ctx)
	}
	return "", errors.New("no wrapper configured")
}

func (a *access) SetConfig(ctx context.Context, options ...wrapping.Option) (*wrapping.WrapperConfig, error) {
	w := a.getDefaultWrapper()
	if w != nil {
		return w.SetConfig(ctx, options...)
	}

	return nil, errors.New("no wrapper configured")
}

func (a *access) GetWrapper() wrapping.Wrapper {
	return a.getDefaultWrapper()
}

func (a *access) Init(ctx context.Context, options ...wrapping.Option) error {
	for _, w := range a.wrappersByPriority {
		if initWrapper, ok := w.Wrapper.(wrapping.InitFinalizer); ok {
			if err := initWrapper.Init(ctx, options...); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *access) getDefaultWrapper() wrapping.Wrapper {
	if len(a.wrappersByPriority) > 0 {
		return a.wrappersByPriority[0].Wrapper
	}
	return nil
}

func (a *access) Type(ctx context.Context) (wrapping.WrapperType, error) {
	return a.getDefaultWrapper().Type(ctx)
}

func (a *access) SealType(ctx context.Context) (SealType, error) {
	if len(a.wrappersByPriority) > 1 {
		return SealTypeMultiSeal, nil
	}

	wrapperType, err := a.getDefaultWrapper().Type(ctx)
	if err != nil {
		return "", err
	}

	return SealType(wrapperType), nil
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *access) Encrypt(ctx context.Context, plaintext []byte, options ...wrapping.Option) (*wrapping.BlobInfo, error) {
	var ciphertext *wrapping.BlobInfo
	for _, wrapper := range a.wrappersByPriority {
		wTyp, err := wrapper.Type(ctx)
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

		ciphertext, err = wrapper.Encrypt(ctx, plaintext, options...)
		if err != nil {
			// TODO (multiseal): logic for failures and setting the sentinel for retrying
			return nil, err
		}

		// TODO (multiseal): logic for new data structure with results from multiple seals
	}

	return ciphertext, nil
}

// Decrypt uses the underlying seal to decrypt the cryptotext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
func (a *access) Decrypt(ctx context.Context, data *wrapping.BlobInfo, options ...wrapping.Option) ([]byte, error) {
	var pt []byte
	var dErr error
	for _, wrapper := range a.wrappersByPriority {
		if func() bool {
			var wTyp wrapping.WrapperType
			wTyp, err := wrapper.Type(ctx)
			if err != nil {
				dErr = err
				return true
			}
			defer func(now time.Time) {
				metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
				metrics.MeasureSince([]string{"seal", wTyp.String(), "decrypt", "time"}, now)

				if dErr != nil {
					metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
					metrics.IncrCounter([]string{"seal", wTyp.String(), "decrypt", "error"}, 1)
				}
			}(time.Now())

			metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
			metrics.IncrCounter([]string{"seal", wTyp.String(), "decrypt"}, 1)

			pt, dErr = wrapper.Decrypt(ctx, data, options...)
			if err == nil {
				return true
			}
			return false
		}() {
			break
		}
		// TODO (multiseal): log an error?
	}

	return pt, dErr
}

func (a *access) Finalize(ctx context.Context, options ...wrapping.Option) error {
	var errs []error

	for _, w := range a.wrappersByPriority {
		if finalizeWrapper, ok := w.Wrapper.(wrapping.InitFinalizer); ok {
			if err := finalizeWrapper.Finalize(ctx, options...); err != nil {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (a *access) SetShamirSealKey(key []byte) error {
	if len(a.wrappersByPriority) == 0 {
		return errors.New("no wrapper configured")
	}

	wrapper := a.wrappersByPriority[0].Wrapper

	shamirWrapper, ok := wrapper.(*aead.ShamirWrapper)
	if !ok {
		return errors.New("seal is not a Shamir seal")
	}

	return shamirWrapper.SetAesGcmKeyBytes(key)
}

func (a *access) GetShamirKeyBytes(ctx context.Context) ([]byte, error) {
	if len(a.wrappersByPriority) == 0 {
		return nil, errors.New("no wrapper configured")
	}

	wrapper := a.wrappersByPriority[0].Wrapper

	shamirWrapper, ok := wrapper.(*aead.ShamirWrapper)
	if !ok {
		return nil, errors.New("seal is not a shamir seal")
	}

	return shamirWrapper.KeyBytes(ctx)
}

type SealType string

const (
	SealTypeMultiSeal         SealType = "multiseal"
	SealTypeAliCloudKms                = SealType(wrapping.WrapperTypeAliCloudKms)
	SealTypeAwsKms                     = SealType(wrapping.WrapperTypeAwsKms)
	SealTypeAzureKeyVault              = SealType(wrapping.WrapperTypeAzureKeyVault)
	SealTypeGcpCkms                    = SealType(wrapping.WrapperTypeGcpCkms)
	SealTypePkcs11                     = SealType(wrapping.WrapperTypePkcs11)
	SealTypeOciKms                     = SealType(wrapping.WrapperTypeOciKms)
	SealTypeShamir                     = SealType(wrapping.WrapperTypeShamir)
	SealTypeTransit                    = SealType(wrapping.WrapperTypeTransit)
	SealTypeHsmAutoDeprecated          = SealType(wrapping.WrapperTypeHsmAuto)
)

func (s SealType) String() string {
	return string(s)
}
