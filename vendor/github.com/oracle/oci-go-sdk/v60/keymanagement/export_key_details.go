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

// ExportKeyDetails The details of the key that you want to wrap and export.
type ExportKeyDetails struct {

	// The OCID of the master encryption key associated with the key version you want to export.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The encryption algorithm to use to encrypt exportable key material from a software-backed key. Specifying `RSA_OAEP_AES_SHA256`
	// invokes the RSA AES key wrap mechanism, which generates a temporary AES key. The temporary AES key is wrapped by the RSA public
	// wrapping key provided along with the request, creating a wrapped temporary AES key. The temporary AES key is also used to wrap
	// the exportable key material. The wrapped temporary AES key and the wrapped exportable key material are concatenated, producing
	// concatenated blob output that jointly represents them. Specifying `RSA_OAEP_SHA256` means that the software key is wrapped by
	// the RSA public wrapping key provided along with the request.
	Algorithm ExportKeyDetailsAlgorithmEnum `mandatory:"true" json:"algorithm"`

	// The PEM format of the 2048-bit, 3072-bit, or 4096-bit RSA wrapping key in your possession that you want to use to encrypt the key.
	PublicKey *string `mandatory:"true" json:"publicKey"`

	// The OCID of the specific key version to export. If not specified, the service exports the current key version.
	KeyVersionId *string `mandatory:"false" json:"keyVersionId"`

	// Information that provides context for audit logging. You can provide this additional
	// data as key-value pairs to include in the audit logs when audit logging is enabled.
	LoggingContext map[string]string `mandatory:"false" json:"loggingContext"`
}

func (m ExportKeyDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ExportKeyDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingExportKeyDetailsAlgorithmEnum(string(m.Algorithm)); !ok && m.Algorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for Algorithm: %s. Supported values are: %s.", m.Algorithm, strings.Join(GetExportKeyDetailsAlgorithmEnumStringValues(), ",")))
	}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// ExportKeyDetailsAlgorithmEnum Enum with underlying type: string
type ExportKeyDetailsAlgorithmEnum string

// Set of constants representing the allowable values for ExportKeyDetailsAlgorithmEnum
const (
	ExportKeyDetailsAlgorithmAesSha256 ExportKeyDetailsAlgorithmEnum = "RSA_OAEP_AES_SHA256"
	ExportKeyDetailsAlgorithmSha256    ExportKeyDetailsAlgorithmEnum = "RSA_OAEP_SHA256"
)

var mappingExportKeyDetailsAlgorithmEnum = map[string]ExportKeyDetailsAlgorithmEnum{
	"RSA_OAEP_AES_SHA256": ExportKeyDetailsAlgorithmAesSha256,
	"RSA_OAEP_SHA256":     ExportKeyDetailsAlgorithmSha256,
}

var mappingExportKeyDetailsAlgorithmEnumLowerCase = map[string]ExportKeyDetailsAlgorithmEnum{
	"rsa_oaep_aes_sha256": ExportKeyDetailsAlgorithmAesSha256,
	"rsa_oaep_sha256":     ExportKeyDetailsAlgorithmSha256,
}

// GetExportKeyDetailsAlgorithmEnumValues Enumerates the set of values for ExportKeyDetailsAlgorithmEnum
func GetExportKeyDetailsAlgorithmEnumValues() []ExportKeyDetailsAlgorithmEnum {
	values := make([]ExportKeyDetailsAlgorithmEnum, 0)
	for _, v := range mappingExportKeyDetailsAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetExportKeyDetailsAlgorithmEnumStringValues Enumerates the set of values in String for ExportKeyDetailsAlgorithmEnum
func GetExportKeyDetailsAlgorithmEnumStringValues() []string {
	return []string{
		"RSA_OAEP_AES_SHA256",
		"RSA_OAEP_SHA256",
	}
}

// GetMappingExportKeyDetailsAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingExportKeyDetailsAlgorithmEnum(val string) (ExportKeyDetailsAlgorithmEnum, bool) {
	enum, ok := mappingExportKeyDetailsAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
