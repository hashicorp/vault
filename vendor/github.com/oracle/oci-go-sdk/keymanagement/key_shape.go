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

// KeyShape The cryptographic properties of a key.
type KeyShape struct {

	// The algorithm used by a key's key versions to encrypt or decrypt.
	Algorithm KeyShapeAlgorithmEnum `mandatory:"true" json:"algorithm"`

	// The length of the key, expressed as an integer. Values of 16, 24, or 32 are supported.
	Length *int `mandatory:"true" json:"length"`
}

func (m KeyShape) String() string {
	return common.PointerString(m)
}

// KeyShapeAlgorithmEnum Enum with underlying type: string
type KeyShapeAlgorithmEnum string

// Set of constants representing the allowable values for KeyShapeAlgorithmEnum
const (
	KeyShapeAlgorithmAes KeyShapeAlgorithmEnum = "AES"
	KeyShapeAlgorithmRsa KeyShapeAlgorithmEnum = "RSA"
)

var mappingKeyShapeAlgorithm = map[string]KeyShapeAlgorithmEnum{
	"AES": KeyShapeAlgorithmAes,
	"RSA": KeyShapeAlgorithmRsa,
}

// GetKeyShapeAlgorithmEnumValues Enumerates the set of values for KeyShapeAlgorithmEnum
func GetKeyShapeAlgorithmEnumValues() []KeyShapeAlgorithmEnum {
	values := make([]KeyShapeAlgorithmEnum, 0)
	for _, v := range mappingKeyShapeAlgorithm {
		values = append(values, v)
	}
	return values
}
