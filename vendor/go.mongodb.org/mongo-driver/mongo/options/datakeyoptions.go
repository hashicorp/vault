// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// DataKeyOptions represents all possible options used to create a new data key.
type DataKeyOptions struct {
	MasterKey   interface{}
	KeyAltNames []string

	// KeyMaterial is used to encrypt data. If omitted, keyMaterial is generated form a cryptographically secure random
	// source. "Key Material" is used interchangeably with "dataKey" and "Data Encryption Key" (DEK).
	KeyMaterial []byte
}

// DataKey creates a new DataKeyOptions instance.
func DataKey() *DataKeyOptions {
	return &DataKeyOptions{}
}

// SetMasterKey specifies a KMS-specific key used to encrypt the new data key.
//
// If being used with a local KMS provider, this option is not applicable and should not be specified.
//
// For the AWS, Azure, and GCP KMS providers, this option is required and must be a document. For each, the value of the
// "endpoint" or "keyVaultEndpoint" must be a host name with an optional port number (e.g. "foo.com" or "foo.com:443").
//
// When using AWS, the document must have the format:
//
//	{
//	  region: <string>,
//	  key: <string>,             // The Amazon Resource Name (ARN) to the AWS customer master key (CMK).
//	  endpoint: Optional<string> // An alternate host identifier to send KMS requests to.
//	}
//
// If unset, the "endpoint" defaults to "kms.<region>.amazonaws.com".
//
// When using Azure, the document must have the format:
//
//	{
//	  keyVaultEndpoint: <string>,  // A host identifier to send KMS requests to.
//	  keyName: <string>,
//	  keyVersion: Optional<string> // A specific version of the named key.
//	}
//
// If unset, "keyVersion" defaults to the key's primary version.
//
// When using GCP, the document must have the format:
//
//	{
//	  projectId: <string>,
//	  location: <string>,
//	  keyRing: <string>,
//	  keyName: <string>,
//	  keyVersion: Optional<string>, // A specific version of the named key.
//	  endpoint: Optional<string>    // An alternate host identifier to send KMS requests to.
//	}
//
// If unset, "keyVersion" defaults to the key's primary version and "endpoint" defaults to "cloudkms.googleapis.com".
func (dk *DataKeyOptions) SetMasterKey(masterKey interface{}) *DataKeyOptions {
	dk.MasterKey = masterKey
	return dk
}

// SetKeyAltNames specifies an optional list of string alternate names used to reference a key. If a key is created'
// with alternate names, encryption may refer to the key by a unique alternate name instead of by _id.
func (dk *DataKeyOptions) SetKeyAltNames(keyAltNames []string) *DataKeyOptions {
	dk.KeyAltNames = keyAltNames
	return dk
}

// SetKeyMaterial will set a custom keyMaterial to DataKeyOptions which can be used to encrypt data.
func (dk *DataKeyOptions) SetKeyMaterial(keyMaterial []byte) *DataKeyOptions {
	dk.KeyMaterial = keyMaterial
	return dk
}

// MergeDataKeyOptions combines the argued DataKeyOptions in a last-one wins fashion.
//
// Deprecated: Merging options structs will not be supported in Go Driver 2.0. Users should create a
// single options struct instead.
func MergeDataKeyOptions(opts ...*DataKeyOptions) *DataKeyOptions {
	dko := DataKey()
	for _, opt := range opts {
		if opt == nil {
			continue
		}

		if opt.MasterKey != nil {
			dko.MasterKey = opt.MasterKey
		}
		if opt.KeyAltNames != nil {
			dko.KeyAltNames = opt.KeyAltNames
		}
		if opt.KeyMaterial != nil {
			dko.KeyMaterial = opt.KeyMaterial
		}
	}

	return dko
}
