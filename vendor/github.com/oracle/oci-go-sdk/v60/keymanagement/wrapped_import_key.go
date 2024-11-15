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

// WrappedImportKey The representation of WrappedImportKey
type WrappedImportKey struct {

	// The key material to import, wrapped by the vault's RSA public wrapping key and base64-encoded.
	KeyMaterial *string `mandatory:"true" json:"keyMaterial"`

	// The wrapping mechanism to use during key import.
	// `RSA_OAEP_AES_SHA256` invokes the RSA AES key wrap mechanism, which generates a temporary AES key. The temporary AES key is wrapped
	// by the vault's RSA public wrapping key, creating a wrapped temporary AES key. The temporary AES key is also used to wrap the private key material.
	// The wrapped temporary AES key and the wrapped exportable key material are concatenated, producing concatenated blob output that jointly represents them.
	// `RSA_OAEP_SHA256` means that the exportable key material is wrapped by the vault's RSA public wrapping key.
	WrappingAlgorithm WrappedImportKeyWrappingAlgorithmEnum `mandatory:"true" json:"wrappingAlgorithm"`
}

func (m WrappedImportKey) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m WrappedImportKey) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingWrappedImportKeyWrappingAlgorithmEnum(string(m.WrappingAlgorithm)); !ok && m.WrappingAlgorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for WrappingAlgorithm: %s. Supported values are: %s.", m.WrappingAlgorithm, strings.Join(GetWrappedImportKeyWrappingAlgorithmEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// WrappedImportKeyWrappingAlgorithmEnum Enum with underlying type: string
type WrappedImportKeyWrappingAlgorithmEnum string

// Set of constants representing the allowable values for WrappedImportKeyWrappingAlgorithmEnum
const (
	WrappedImportKeyWrappingAlgorithmSha256    WrappedImportKeyWrappingAlgorithmEnum = "RSA_OAEP_SHA256"
	WrappedImportKeyWrappingAlgorithmAesSha256 WrappedImportKeyWrappingAlgorithmEnum = "RSA_OAEP_AES_SHA256"
)

var mappingWrappedImportKeyWrappingAlgorithmEnum = map[string]WrappedImportKeyWrappingAlgorithmEnum{
	"RSA_OAEP_SHA256":     WrappedImportKeyWrappingAlgorithmSha256,
	"RSA_OAEP_AES_SHA256": WrappedImportKeyWrappingAlgorithmAesSha256,
}

var mappingWrappedImportKeyWrappingAlgorithmEnumLowerCase = map[string]WrappedImportKeyWrappingAlgorithmEnum{
	"rsa_oaep_sha256":     WrappedImportKeyWrappingAlgorithmSha256,
	"rsa_oaep_aes_sha256": WrappedImportKeyWrappingAlgorithmAesSha256,
}

// GetWrappedImportKeyWrappingAlgorithmEnumValues Enumerates the set of values for WrappedImportKeyWrappingAlgorithmEnum
func GetWrappedImportKeyWrappingAlgorithmEnumValues() []WrappedImportKeyWrappingAlgorithmEnum {
	values := make([]WrappedImportKeyWrappingAlgorithmEnum, 0)
	for _, v := range mappingWrappedImportKeyWrappingAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetWrappedImportKeyWrappingAlgorithmEnumStringValues Enumerates the set of values in String for WrappedImportKeyWrappingAlgorithmEnum
func GetWrappedImportKeyWrappingAlgorithmEnumStringValues() []string {
	return []string{
		"RSA_OAEP_SHA256",
		"RSA_OAEP_AES_SHA256",
	}
}

// GetMappingWrappedImportKeyWrappingAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingWrappedImportKeyWrappingAlgorithmEnum(val string) (WrappedImportKeyWrappingAlgorithmEnum, bool) {
	enum, ok := mappingWrappedImportKeyWrappingAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
