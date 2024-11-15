// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package wrapping

import (
	"context"
)

type HmacComputer interface {
	// HmacKeyID is the ID of the key currently used for HMACing (if any)
	HmacKeyId(context.Context) (string, error)
}

type InitFinalizer interface {
	// Init allows performing any necessary setup calls before using a
	// Wrapper.
	Init(ctx context.Context, options ...Option) error

	// Finalize can be called when all usage of a Wrapper is done if any cleanup
	// or finalization is required.
	Finalize(ctx context.Context, options ...Option) error
}

// Wrapper is an an interface where supporting implementations allow for
// encrypting and decrypting data.
type Wrapper interface {
	// Type is the type of Wrapper
	Type(context.Context) (WrapperType, error)

	// KeyId is the ID of the key currently used for encryption
	KeyId(context.Context) (string, error)

	// SetConfig applies the given options to a wrapper and returns
	// configuration information. WithConfigMap will almost certainly be
	// required to be passed in to give wrapper-specific configuration
	// information to the wrapper. WithKeyId is also supported.
	SetConfig(ctx context.Context, options ...Option) (*WrapperConfig, error)

	// Encrypt encrypts the given byte slice and stores the resulting
	// information in the returned blob info. Which options are used depends on
	// the underlying wrapper. Supported options: WithAad.
	Encrypt(ctx context.Context, plaintext []byte, options ...Option) (*BlobInfo, error)
	// Decrypt decrypts the given byte slice and stores the resulting
	// information in the returned byte slice. Which options are used depends on
	// the underlying wrapper. Supported options: WithAad.
	Decrypt(ctx context.Context, ciphertext *BlobInfo, options ...Option) ([]byte, error)
}

// KeyExporter defines an optional interface for wrappers to implement that returns
// the "current" key bytes. This will be implementation-specific.
type KeyExporter interface {
	// KeyBytes returns the "current" key bytes
	KeyBytes(context.Context) ([]byte, error)
}
