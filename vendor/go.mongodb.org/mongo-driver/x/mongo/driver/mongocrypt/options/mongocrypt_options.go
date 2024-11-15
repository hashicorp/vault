// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"net/http"

	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MongoCryptOptions specifies options to configure a MongoCrypt instance.
type MongoCryptOptions struct {
	KmsProviders               bsoncore.Document
	LocalSchemaMap             map[string]bsoncore.Document
	BypassQueryAnalysis        bool
	EncryptedFieldsMap         map[string]bsoncore.Document
	CryptSharedLibDisabled     bool
	CryptSharedLibOverridePath string
	HTTPClient                 *http.Client
}

// MongoCrypt creates a new MongoCryptOptions instance.
func MongoCrypt() *MongoCryptOptions {
	return &MongoCryptOptions{}
}

// SetKmsProviders specifies the KMS providers map.
func (mo *MongoCryptOptions) SetKmsProviders(kmsProviders bsoncore.Document) *MongoCryptOptions {
	mo.KmsProviders = kmsProviders
	return mo
}

// SetLocalSchemaMap specifies the local schema map.
func (mo *MongoCryptOptions) SetLocalSchemaMap(localSchemaMap map[string]bsoncore.Document) *MongoCryptOptions {
	mo.LocalSchemaMap = localSchemaMap
	return mo
}

// SetBypassQueryAnalysis skips the NeedMongoMarkings state.
func (mo *MongoCryptOptions) SetBypassQueryAnalysis(bypassQueryAnalysis bool) *MongoCryptOptions {
	mo.BypassQueryAnalysis = bypassQueryAnalysis
	return mo
}

// SetEncryptedFieldsMap specifies the encrypted fields map.
func (mo *MongoCryptOptions) SetEncryptedFieldsMap(efcMap map[string]bsoncore.Document) *MongoCryptOptions {
	mo.EncryptedFieldsMap = efcMap
	return mo
}

// SetCryptSharedLibDisabled explicitly disables loading the crypt_shared library if set to true.
func (mo *MongoCryptOptions) SetCryptSharedLibDisabled(disabled bool) *MongoCryptOptions {
	mo.CryptSharedLibDisabled = disabled
	return mo
}

// SetCryptSharedLibOverridePath sets the override path to the crypt_shared library file. Setting
// an override path disables the default operating system dynamic library search path.
func (mo *MongoCryptOptions) SetCryptSharedLibOverridePath(path string) *MongoCryptOptions {
	mo.CryptSharedLibOverridePath = path
	return mo
}

// SetHTTPClient sets the http client.
func (mo *MongoCryptOptions) SetHTTPClient(httpClient *http.Client) *MongoCryptOptions {
	mo.HTTPClient = httpClient
	return mo
}
