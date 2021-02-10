// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// GenerateKeyDetails The representation of GenerateKeyDetails
type GenerateKeyDetails struct {

	// If true, the generated key is also returned unencrypted.
	IncludePlaintextKey *bool `mandatory:"true" json:"includePlaintextKey"`

	// The OCID of the master encryption key to encrypt the generated data encryption key with.
	KeyId *string `mandatory:"true" json:"keyId"`

	KeyShape *KeyShape `mandatory:"true" json:"keyShape"`

	// Information that can be used to provide an encryption context for the encrypted data.
	// The length of the string representation of the associated data must be fewer than 4096
	// characters.
	AssociatedData map[string]string `mandatory:"false" json:"associatedData"`

	// Information that provides context for audit logging. You can provide this additional
	// data by formatting it as key-value pairs to include in audit logs when audit logging is enabled.
	LoggingContext map[string]string `mandatory:"false" json:"loggingContext"`
}

func (m GenerateKeyDetails) String() string {
	return common.PointerString(m)
}
