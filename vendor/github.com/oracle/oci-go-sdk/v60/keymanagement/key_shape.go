// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Vault Service Key Management API
//
// API for managing and performing operations with keys and vaults. (For the API for managing secrets, see the Vault Service
// Secret Management API. For the API for retrieving secrets, see the Vault Service Secret Retrieval API.)
//

package keymanagement

import (
	"fmt"
	"github.com/oracle/oci-go-sdk/v60/common"
	"strings"
)

// KeyShape The cryptographic properties of a key.
type KeyShape struct {

	// The algorithm used by a key's key versions to encrypt or decrypt.
	Algorithm KeyShapeAlgorithmEnum `mandatory:"true" json:"algorithm"`

	// The length of the key in bytes, expressed as an integer. Supported values include the following:
	//   - AES: 16, 24, or 32
	//   - RSA: 256, 384, or 512
	//   - ECDSA: 32, 48, or 66
	Length *int `mandatory:"true" json:"length"`

	// Supported curve IDs for ECDSA keys.
	CurveId KeyShapeCurveIdEnum `mandatory:"false" json:"curveId,omitempty"`
}

func (m KeyShape) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m KeyShape) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingKeyShapeAlgorithmEnum(string(m.Algorithm)); !ok && m.Algorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Algorithm: %s. Supported values are: %s.", m.Algorithm, strings.Join(GetKeyShapeAlgorithmEnumStringValues(), ",")))
	}

	if _, ok := GetMappingKeyShapeCurveIdEnum(string(m.CurveId)); !ok && m.CurveId != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for CurveId: %s. Supported values are: %s.", m.CurveId, strings.Join(GetKeyShapeCurveIdEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// KeyShapeAlgorithmEnum Enum with underlying type: string
type KeyShapeAlgorithmEnum string

// Set of constants representing the allowable values for KeyShapeAlgorithmEnum
const (
	KeyShapeAlgorithmAes   KeyShapeAlgorithmEnum = "AES"
	KeyShapeAlgorithmRsa   KeyShapeAlgorithmEnum = "RSA"
	KeyShapeAlgorithmEcdsa KeyShapeAlgorithmEnum = "ECDSA"
)

var mappingKeyShapeAlgorithmEnum = map[string]KeyShapeAlgorithmEnum{
	"AES":   KeyShapeAlgorithmAes,
	"RSA":   KeyShapeAlgorithmRsa,
	"ECDSA": KeyShapeAlgorithmEcdsa,
}

var mappingKeyShapeAlgorithmEnumLowerCase = map[string]KeyShapeAlgorithmEnum{
	"aes":   KeyShapeAlgorithmAes,
	"rsa":   KeyShapeAlgorithmRsa,
	"ecdsa": KeyShapeAlgorithmEcdsa,
}

// GetKeyShapeAlgorithmEnumValues Enumerates the set of values for KeyShapeAlgorithmEnum
func GetKeyShapeAlgorithmEnumValues() []KeyShapeAlgorithmEnum {
	values := make([]KeyShapeAlgorithmEnum, 0)
	for _, v := range mappingKeyShapeAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyShapeAlgorithmEnumStringValues Enumerates the set of values in String for KeyShapeAlgorithmEnum
func GetKeyShapeAlgorithmEnumStringValues() []string {
	return []string{
		"AES",
		"RSA",
		"ECDSA",
	}
}

// GetMappingKeyShapeAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyShapeAlgorithmEnum(val string) (KeyShapeAlgorithmEnum, bool) {
	enum, ok := mappingKeyShapeAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// KeyShapeCurveIdEnum Enum with underlying type: string
type KeyShapeCurveIdEnum string

// Set of constants representing the allowable values for KeyShapeCurveIdEnum
const (
	KeyShapeCurveIdP256 KeyShapeCurveIdEnum = "NIST_P256"
	KeyShapeCurveIdP384 KeyShapeCurveIdEnum = "NIST_P384"
	KeyShapeCurveIdP521 KeyShapeCurveIdEnum = "NIST_P521"
)

var mappingKeyShapeCurveIdEnum = map[string]KeyShapeCurveIdEnum{
	"NIST_P256": KeyShapeCurveIdP256,
	"NIST_P384": KeyShapeCurveIdP384,
	"NIST_P521": KeyShapeCurveIdP521,
}

var mappingKeyShapeCurveIdEnumLowerCase = map[string]KeyShapeCurveIdEnum{
	"nist_p256": KeyShapeCurveIdP256,
	"nist_p384": KeyShapeCurveIdP384,
	"nist_p521": KeyShapeCurveIdP521,
}

// GetKeyShapeCurveIdEnumValues Enumerates the set of values for KeyShapeCurveIdEnum
func GetKeyShapeCurveIdEnumValues() []KeyShapeCurveIdEnum {
	values := make([]KeyShapeCurveIdEnum, 0)
	for _, v := range mappingKeyShapeCurveIdEnum {
		values = append(values, v)
	}
	return values
}

// GetKeyShapeCurveIdEnumStringValues Enumerates the set of values in String for KeyShapeCurveIdEnum
func GetKeyShapeCurveIdEnumStringValues() []string {
	return []string{
		"NIST_P256",
		"NIST_P384",
		"NIST_P521",
	}
}

// GetMappingKeyShapeCurveIdEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingKeyShapeCurveIdEnum(val string) (KeyShapeCurveIdEnum, bool) {
	enum, ok := mappingKeyShapeCurveIdEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
