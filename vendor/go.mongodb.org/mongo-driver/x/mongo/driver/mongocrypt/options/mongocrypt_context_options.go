// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// DataKeyOptions specifies options for creating a new data key.
type DataKeyOptions struct {
	KeyAltNames []string
	KeyMaterial []byte
	MasterKey   bsoncore.Document
}

// DataKey creates a new DataKeyOptions instance.
func DataKey() *DataKeyOptions {
	return &DataKeyOptions{}
}

// SetKeyAltNames specifies alternate key names.
func (dko *DataKeyOptions) SetKeyAltNames(names []string) *DataKeyOptions {
	dko.KeyAltNames = names
	return dko
}

// SetMasterKey specifies the master key.
func (dko *DataKeyOptions) SetMasterKey(key bsoncore.Document) *DataKeyOptions {
	dko.MasterKey = key
	return dko
}

// SetKeyMaterial specifies the key material.
func (dko *DataKeyOptions) SetKeyMaterial(keyMaterial []byte) *DataKeyOptions {
	dko.KeyMaterial = keyMaterial
	return dko
}

// QueryType describes the type of query the result of Encrypt is used for.
type QueryType int

// These constants specify valid values for QueryType
const (
	QueryTypeEquality QueryType = 1
)

// ExplicitEncryptionOptions specifies options for configuring an explicit encryption context.
type ExplicitEncryptionOptions struct {
	KeyID            *primitive.Binary
	KeyAltName       *string
	Algorithm        string
	QueryType        string
	ContentionFactor *int64
	RangeOptions     *ExplicitRangeOptions
}

// ExplicitRangeOptions specifies options for the range index.
type ExplicitRangeOptions struct {
	Min        *bsoncore.Value
	Max        *bsoncore.Value
	Sparsity   *int64
	TrimFactor *int32
	Precision  *int32
}

// ExplicitEncryption creates a new ExplicitEncryptionOptions instance.
func ExplicitEncryption() *ExplicitEncryptionOptions {
	return &ExplicitEncryptionOptions{}
}

// SetKeyID sets the key identifier.
func (eeo *ExplicitEncryptionOptions) SetKeyID(keyID primitive.Binary) *ExplicitEncryptionOptions {
	eeo.KeyID = &keyID
	return eeo
}

// SetKeyAltName sets the key alternative name.
func (eeo *ExplicitEncryptionOptions) SetKeyAltName(keyAltName string) *ExplicitEncryptionOptions {
	eeo.KeyAltName = &keyAltName
	return eeo
}

// SetAlgorithm specifies an encryption algorithm.
func (eeo *ExplicitEncryptionOptions) SetAlgorithm(algorithm string) *ExplicitEncryptionOptions {
	eeo.Algorithm = algorithm
	return eeo
}

// SetQueryType specifies the query type.
func (eeo *ExplicitEncryptionOptions) SetQueryType(queryType string) *ExplicitEncryptionOptions {
	eeo.QueryType = queryType
	return eeo
}

// SetContentionFactor specifies the contention factor.
func (eeo *ExplicitEncryptionOptions) SetContentionFactor(contentionFactor int64) *ExplicitEncryptionOptions {
	eeo.ContentionFactor = &contentionFactor
	return eeo
}

// SetRangeOptions specifies the range options.
func (eeo *ExplicitEncryptionOptions) SetRangeOptions(ro ExplicitRangeOptions) *ExplicitEncryptionOptions {
	eeo.RangeOptions = &ro
	return eeo
}

// RewrapManyDataKeyOptions represents all possible options used to decrypt and encrypt all matching data keys with a
// possibly new masterKey.
type RewrapManyDataKeyOptions struct {
	// Provider identifies the new KMS provider. If omitted, encrypting uses the current KMS provider.
	Provider *string

	// MasterKey identifies the new masterKey. If omitted, rewraps with the current masterKey.
	MasterKey bsoncore.Document
}

// RewrapManyDataKey creates a new RewrapManyDataKeyOptions instance.
func RewrapManyDataKey() *RewrapManyDataKeyOptions {
	return new(RewrapManyDataKeyOptions)
}

// SetProvider sets the value for the Provider field.
func (rmdko *RewrapManyDataKeyOptions) SetProvider(provider string) *RewrapManyDataKeyOptions {
	rmdko.Provider = &provider
	return rmdko
}

// SetMasterKey sets the value for the MasterKey field.
func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey bsoncore.Document) *RewrapManyDataKeyOptions {
	rmdko.MasterKey = masterKey
	return rmdko
}

// MergeRewrapManyDataKeyOptions combines the given RewrapManyDataKeyOptions instances into a single
// RewrapManyDataKeyOptions in a last one wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeRewrapManyDataKeyOptions(opts ...*RewrapManyDataKeyOptions) *RewrapManyDataKeyOptions {
	rmdkOpts := RewrapManyDataKey()
	for _, rmdko := range opts {
		if rmdko == nil {
			continue
		}
		if provider := rmdko.Provider; provider != nil {
			rmdkOpts.Provider = provider
		}
		if masterKey := rmdko.MasterKey; masterKey != nil {
			rmdkOpts.MasterKey = masterKey
		}
	}
	return rmdkOpts
}
