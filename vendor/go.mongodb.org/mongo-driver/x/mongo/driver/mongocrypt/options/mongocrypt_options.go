// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package options

import (
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"
)

// MongoCryptOptions specifies options to configure a MongoCrypt instance.
type MongoCryptOptions struct {
	KmsProviders   bsoncore.Document
	LocalSchemaMap map[string]bsoncore.Document
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
