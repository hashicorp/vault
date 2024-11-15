// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// These constants specify valid values for QueryType
// QueryType is used for Queryable Encryption.
const (
	QueryTypeEquality string = "equality"
)

// RangeOptions specifies index options for a Queryable Encryption field supporting "range" queries.
type RangeOptions struct {
	Min        *bson.RawValue
	Max        *bson.RawValue
	Sparsity   *int64
	TrimFactor *int32
	Precision  *int32
}

// EncryptOptions represents options to explicitly encrypt a value.
type EncryptOptions struct {
	KeyID            *primitive.Binary
	KeyAltName       *string
	Algorithm        string
	QueryType        string
	ContentionFactor *int64
	RangeOptions     *RangeOptions
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

// SetAlgorithm specifies an algorithm to use for encryption. This should be one of the following:
// - AEAD_AES_256_CBC_HMAC_SHA_512-Deterministic
// - AEAD_AES_256_CBC_HMAC_SHA_512-Random
// - Indexed
// - Unindexed
// - Range
// This is required.
// Indexed and Unindexed are used for Queryable Encryption.
func (e *EncryptOptions) SetAlgorithm(algorithm string) *EncryptOptions {
	e.Algorithm = algorithm
	return e
}

// SetQueryType specifies the intended query type. It is only valid to set if algorithm is "Indexed".
// This should be one of the following:
// - equality
// QueryType is used for Queryable Encryption.
func (e *EncryptOptions) SetQueryType(queryType string) *EncryptOptions {
	e.QueryType = queryType
	return e
}

// SetContentionFactor specifies the contention factor. It is only valid to set if algorithm is "Indexed".
// ContentionFactor is used for Queryable Encryption.
func (e *EncryptOptions) SetContentionFactor(contentionFactor int64) *EncryptOptions {
	e.ContentionFactor = &contentionFactor
	return e
}

// SetRangeOptions specifies the options to use for explicit encryption with range. It is only valid to set if algorithm is "Range".
func (e *EncryptOptions) SetRangeOptions(ro RangeOptions) *EncryptOptions {
	e.RangeOptions = &ro
	return e
}

// SetMin sets the range index minimum value.
func (ro *RangeOptions) SetMin(min bson.RawValue) *RangeOptions {
	ro.Min = &min
	return ro
}

// SetMax sets the range index maximum value.
func (ro *RangeOptions) SetMax(max bson.RawValue) *RangeOptions {
	ro.Max = &max
	return ro
}

// SetSparsity sets the range index sparsity.
func (ro *RangeOptions) SetSparsity(sparsity int64) *RangeOptions {
	ro.Sparsity = &sparsity
	return ro
}

// SetTrimFactor sets the range index trim factor.
func (ro *RangeOptions) SetTrimFactor(trimFactor int32) *RangeOptions {
	ro.TrimFactor = &trimFactor
	return ro
}

// SetPrecision sets the range index precision.
func (ro *RangeOptions) SetPrecision(precision int32) *RangeOptions {
	ro.Precision = &precision
	return ro
}

// MergeEncryptOptions combines the argued EncryptOptions in a last-one wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
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
		if opt.QueryType != "" {
			eo.QueryType = opt.QueryType
		}
		if opt.ContentionFactor != nil {
			eo.ContentionFactor = opt.ContentionFactor
		}
		if opt.RangeOptions != nil {
			eo.RangeOptions = opt.RangeOptions
		}
	}

	return eo
}
