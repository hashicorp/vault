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

// WrappedImportKey The representation of WrappedImportKey
type WrappedImportKey struct {

	// The wrapped/encrypted key material to import. It is encrypted using RSA wrapped key and Base64 encoded.
	KeyMaterial *string `mandatory:"true" json:"keyMaterial"`

	// The wrapping mechanism to be used during key import
	WrappingAlgorithm WrappedImportKeyWrappingAlgorithmEnum `mandatory:"true" json:"wrappingAlgorithm"`
}

func (m WrappedImportKey) String() string {
	return common.PointerString(m)
}

// WrappedImportKeyWrappingAlgorithmEnum Enum with underlying type: string
type WrappedImportKeyWrappingAlgorithmEnum string

// Set of constants representing the allowable values for WrappedImportKeyWrappingAlgorithmEnum
const (
	WrappedImportKeyWrappingAlgorithmRsaOaepSha256 WrappedImportKeyWrappingAlgorithmEnum = "RSA_OAEP_SHA256"
)

var mappingWrappedImportKeyWrappingAlgorithm = map[string]WrappedImportKeyWrappingAlgorithmEnum{
	"RSA_OAEP_SHA256": WrappedImportKeyWrappingAlgorithmRsaOaepSha256,
}

// GetWrappedImportKeyWrappingAlgorithmEnumValues Enumerates the set of values for WrappedImportKeyWrappingAlgorithmEnum
func GetWrappedImportKeyWrappingAlgorithmEnumValues() []WrappedImportKeyWrappingAlgorithmEnum {
	values := make([]WrappedImportKeyWrappingAlgorithmEnum, 0)
	for _, v := range mappingWrappedImportKeyWrappingAlgorithm {
		values = append(values, v)
	}
	return values
}
