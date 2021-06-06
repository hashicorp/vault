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
	AwsProviderOpts   *AwsKmsProviderOptions
	LocalProviderOpts *LocalKmsProviderOptions
	LocalSchemaMap    map[string]bsoncore.Document
}

// MongoCrypt creates a new MongoCryptOptions instance.
func MongoCrypt() *MongoCryptOptions {
	return &MongoCryptOptions{}
}

// SetAwsProviderOptions specifies AWS KMS provider options.
func (mo *MongoCryptOptions) SetAwsProviderOptions(awsOpts *AwsKmsProviderOptions) *MongoCryptOptions {
	mo.AwsProviderOpts = awsOpts
	return mo
}

// SetLocalProviderOptions specifies local KMS provider options.
func (mo *MongoCryptOptions) SetLocalProviderOptions(localOpts *LocalKmsProviderOptions) *MongoCryptOptions {
	mo.LocalProviderOpts = localOpts
	return mo
}

// SetLocalSchemaMap specifies the local schema map.
func (mo *MongoCryptOptions) SetLocalSchemaMap(localSchemaMap map[string]bsoncore.Document) *MongoCryptOptions {
	mo.LocalSchemaMap = localSchemaMap
	return mo
}
