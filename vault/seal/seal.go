package seal

import (
	"context"
	"time"

	metrics "github.com/armon/go-metrics"
	wrapping "github.com/hashicorp/go-kms-wrapping"
)

type StoredKeysSupport int

const (
	// The 0 value of StoredKeysSupport is an invalid option
	StoredKeysInvalid StoredKeysSupport = iota
	StoredKeysNotSupported
	StoredKeysSupportedGeneric
	StoredKeysSupportedShamirMaster
)

func (s StoredKeysSupport) String() string {
	switch s {
	case StoredKeysNotSupported:
		return "Old-style Shamir"
	case StoredKeysSupportedGeneric:
		return "AutoUnseal"
	case StoredKeysSupportedShamirMaster:
		return "New-style Shamir"
	default:
		return "Invalid StoredKeys type"
	}
}

// Access is the embedded implementation of autoSeal that contains logic
// specific to encrypting and decrypting data, or in this case keys.
type Access struct {
	wrapping.Wrapper
	OverriddenType string
}

func (a *Access) SetType(t string) {
	a.OverriddenType = t
}

func (a *Access) Type() string {
	if a.OverriddenType != "" {
		return a.OverriddenType
	}
	return a.Wrapper.Type()
}

// Encrypt uses the underlying seal to encrypt the plaintext and returns it.
func (a *Access) Encrypt(ctx context.Context, plaintext, aad []byte) (blob *wrapping.EncryptedBlobInfo, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "encrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", a.Wrapper.Type(), "encrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "encrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", a.Wrapper.Type(), "encrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "encrypt"}, 1)
	metrics.IncrCounter([]string{"seal", a.Wrapper.Type(), "encrypt"}, 1)

	return a.Wrapper.Encrypt(ctx, plaintext, aad)
}

// Decrypt uses the underlying seal to decrypt the cryptotext and returns it.
// Note that it is possible depending on the wrapper used that both pt and err
// are populated.
func (a *Access) Decrypt(ctx context.Context, data *wrapping.EncryptedBlobInfo, aad []byte) (pt []byte, err error) {
	defer func(now time.Time) {
		metrics.MeasureSince([]string{"seal", "decrypt", "time"}, now)
		metrics.MeasureSince([]string{"seal", a.Wrapper.Type(), "decrypt", "time"}, now)

		if err != nil {
			metrics.IncrCounter([]string{"seal", "decrypt", "error"}, 1)
			metrics.IncrCounter([]string{"seal", a.Wrapper.Type(), "decrypt", "error"}, 1)
		}
	}(time.Now())

	metrics.IncrCounter([]string{"seal", "decrypt"}, 1)
	metrics.IncrCounter([]string{"seal", a.Wrapper.Type(), "decrypt"}, 1)

	return a.Wrapper.Decrypt(ctx, data, aad)
}
