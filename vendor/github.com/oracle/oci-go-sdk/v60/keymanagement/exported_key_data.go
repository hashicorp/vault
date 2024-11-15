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

// ExportedKeyData The response to a request to export key material.
type ExportedKeyData struct {

	// The OCID of the key version.
	KeyVersionId *string `mandatory:"true" json:"keyVersionId"`

	// The OCID of the master encryption key associated with this key version.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The date and time this key version was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key version.
	VaultId *string `mandatory:"true" json:"vaultId"`

	// The base64-encoded exported key material, which is encrypted by using the public RSA wrapping key specified in the export request.
	EncryptedKey *string `mandatory:"true" json:"encryptedKey"`

	// The encryption algorithm to use to encrypt exportable key material from a key that persists on the server (as opposed to a key that
	// persists on a hardware security module and, therefore, cannot be exported). Specifying RSA_OAEP_AES_SHA256 invokes the RSA AES key
	// wrap mechanism, which generates a temporary AES key. The temporary AES key is wrapped by the RSA public wrapping key provided along
	// with the request, creating a wrapped temporary AES key. The temporary AES key is also used to wrap the exportable key material. The
	// wrapped temporary AES key and the wrapped exportable key material are concatenated, producing concatenated blob output that jointly
	// represents them. Specifying RSA_OAEP_SHA256 means that the exportable key material is wrapped by the RSA public wrapping key provided
	// along with the request.
	Algorithm ExportedKeyDataAlgorithmEnum `mandatory:"true" json:"algorithm"`
}

func (m ExportedKeyData) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ExportedKeyData) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingExportedKeyDataAlgorithmEnum(string(m.Algorithm)); !ok && m.Algorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Algorithm: %s. Supported values are: %s.", m.Algorithm, strings.Join(GetExportedKeyDataAlgorithmEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// ExportedKeyDataAlgorithmEnum Enum with underlying type: string
type ExportedKeyDataAlgorithmEnum string

// Set of constants representing the allowable values for ExportedKeyDataAlgorithmEnum
const (
	ExportedKeyDataAlgorithmAesSha256 ExportedKeyDataAlgorithmEnum = "RSA_OAEP_AES_SHA256"
	ExportedKeyDataAlgorithmSha256    ExportedKeyDataAlgorithmEnum = "RSA_OAEP_SHA256"
)

var mappingExportedKeyDataAlgorithmEnum = map[string]ExportedKeyDataAlgorithmEnum{
	"RSA_OAEP_AES_SHA256": ExportedKeyDataAlgorithmAesSha256,
	"RSA_OAEP_SHA256":     ExportedKeyDataAlgorithmSha256,
}

var mappingExportedKeyDataAlgorithmEnumLowerCase = map[string]ExportedKeyDataAlgorithmEnum{
	"rsa_oaep_aes_sha256": ExportedKeyDataAlgorithmAesSha256,
	"rsa_oaep_sha256":     ExportedKeyDataAlgorithmSha256,
}

// GetExportedKeyDataAlgorithmEnumValues Enumerates the set of values for ExportedKeyDataAlgorithmEnum
func GetExportedKeyDataAlgorithmEnumValues() []ExportedKeyDataAlgorithmEnum {
	values := make([]ExportedKeyDataAlgorithmEnum, 0)
	for _, v := range mappingExportedKeyDataAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetExportedKeyDataAlgorithmEnumStringValues Enumerates the set of values in String for ExportedKeyDataAlgorithmEnum
func GetExportedKeyDataAlgorithmEnumStringValues() []string {
	return []string{
		"RSA_OAEP_AES_SHA256",
		"RSA_OAEP_SHA256",
	}
}

// GetMappingExportedKeyDataAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingExportedKeyDataAlgorithmEnum(val string) (ExportedKeyDataAlgorithmEnum, bool) {
	enum, ok := mappingExportedKeyDataAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
