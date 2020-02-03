// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// EncryptOptions represents options to explicitly encrypt a value.
type EncryptOptions struct {
	KeyID      *primitive.Binary
	KeyAltName *string
	Algorithm  string
}

// Encrypt creates a new EncryptOptions instance.
func Encrypt() *EncryptOptions {
	return &EncryptOptions{}
}

// SetKeyID specifies an _id of a data key. This should be a UUID (a primitive.Binary with subtype 4).
func (e *EncryptOptions) SetKeyID(keyID primitive.Binary) *EncryptOptions {
	e.KeyID = &keyID
	return e
}

// SetKeyAltName identifies a key vault document by 'keyAltName'.
func (e *EncryptOptions) SetKeyAltName(keyAltName string) *EncryptOptions {
	e.KeyAltName = &keyAltName
	return e
}

// SetAlgorithm specifies an algorithm to use for encryption. This should be AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic
// or AEAD_AES_256_CBC_HMAC_SHA_512-Random. This is required.
func (e *EncryptOptions) SetAlgorithm(algorithm string) *EncryptOptions {
	e.Algorithm = algorithm
	return e
}

// MergeEncryptOptions combines the argued EncryptOptions in a last-one wins fashion.
func MergeEncryptOptions(opts ...*EncryptOptions) *EncryptOptions {
	eo := Encrypt()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.KeyID != nil {
			eo.KeyID = opt.KeyID
		}
		if opt.KeyAltName != nil {
			eo.KeyAltName = opt.KeyAltName
		}
		if opt.Algorithm != "" {
			eo.Algorithm = opt.Algorithm
		}
	}

	return eo
}
