// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
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

	// Information that can be used to provide an encryption context for the
	// encrypted data. The length of the string representation of the associatedData
	// must be fewer than 4096 characters.
	AssociatedData map[string]string `mandatory:"false" json:"associatedData"`

	// Information that can be used to provide context for audit logging. It is a map that contains any addtional
	// data the users may have and will be added to the audit logs (if audit logging is enabled)
	LoggingContext map[string]string `mandatory:"false" json:"loggingContext"`
}

func (m GenerateKeyDetails) String() string {
	return common.PointerString(m)
}
