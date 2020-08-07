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

// ExplicitEncryptionOptions specifies options for configuring an explicit encryption context.
type ExplicitEncryptionOptions struct {
	KeyID      *primitive.Binary
	KeyAltName *string
	Algorithm  string
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
