// Copyright (C) MongoDB, Inc. 2022-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

// RewrapManyDataKeyOptions represents all possible options used to decrypt and encrypt all matching data keys with a
// possibly new masterKey.
type RewrapManyDataKeyOptions struct {
	// Provider identifies the new KMS provider. If omitted, encrypting uses the current KMS provider.
	Provider *string

	// MasterKey identifies the new masterKey. If omitted, rewraps with the current masterKey.
	MasterKey interface{}
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
func (rmdko *RewrapManyDataKeyOptions) SetMasterKey(masterKey interface{}) *RewrapManyDataKeyOptions {
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
