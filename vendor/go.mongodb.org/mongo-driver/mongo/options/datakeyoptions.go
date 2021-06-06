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
}

// DataKey creates a new DataKeyOptions instance.
func DataKey() *DataKeyOptions {
	return &DataKeyOptions{}
}

// SetMasterKey specifies a KMS-specific key used to encrypt the new data key.
//
// If being used with the AWS KMS provider, this option is required and must be a document with the following format:
// {region: string, key: string}.
//
// If being used with a local KMS provider, this option is not applicable and should not be specified.
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

// MergeDataKeyOptions combines the argued DataKeyOptions in a last-one wins fashion.
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
	}

	return dko
}
