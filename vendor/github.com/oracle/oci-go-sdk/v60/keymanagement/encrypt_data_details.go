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

// EncryptDataDetails The representation of EncryptDataDetails
type EncryptDataDetails struct {

	// The OCID of the key to encrypt with.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The plaintext data to encrypt.
	Plaintext *string `mandatory:"true" json:"plaintext"`

	// Information that can be used to provide an encryption context for the
	// encrypted data. The length of the string representation of the associated data
	// must be fewer than 4096 characters.
	AssociatedData map[string]string `mandatory:"false" json:"associatedData"`

	// Information that provides context for audit logging. You can provide this additional
	// data as key-value pairs to include in the audit logs when audit logging is enabled.
	LoggingContext map[string]string `mandatory:"false" json:"loggingContext"`

	// The OCID of the key version used to encrypt the ciphertext.
	KeyVersionId *string `mandatory:"false" json:"keyVersionId"`

	// The encryption algorithm to use to encrypt and decrypt data with a customer-managed key.
	// `AES_256_GCM` indicates that the key is a symmetric key that uses the Advanced Encryption Standard (AES) algorithm and
	// that the mode of encryption is the Galois/Counter Mode (GCM). `RSA_OAEP_SHA_1` indicates that the
	// key is an asymmetric key that uses the RSA encryption algorithm and uses Optimal Asymmetric Encryption Padding (OAEP).
	// `RSA_OAEP_SHA_256` indicates that the key is an asymmetric key that uses the RSA encryption algorithm with a SHA-256 hash
	// and uses OAEP.
	EncryptionAlgorithm EncryptDataDetailsEncryptionAlgorithmEnum `mandatory:"false" json:"encryptionAlgorithm,omitempty"`
}

func (m EncryptDataDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m EncryptDataDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if _, ok := GetMappingEncryptDataDetailsEncryptionAlgorithmEnum(string(m.EncryptionAlgorithm)); !ok && m.EncryptionAlgorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for EncryptionAlgorithm: %s. Supported values are: %s.", m.EncryptionAlgorithm, strings.Join(GetEncryptDataDetailsEncryptionAlgorithmEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// EncryptDataDetailsEncryptionAlgorithmEnum Enum with underlying type: string
type EncryptDataDetailsEncryptionAlgorithmEnum string

// Set of constants representing the allowable values for EncryptDataDetailsEncryptionAlgorithmEnum
const (
	EncryptDataDetailsEncryptionAlgorithmAes256Gcm     EncryptDataDetailsEncryptionAlgorithmEnum = "AES_256_GCM"
	EncryptDataDetailsEncryptionAlgorithmRsaOaepSha1   EncryptDataDetailsEncryptionAlgorithmEnum = "RSA_OAEP_SHA_1"
	EncryptDataDetailsEncryptionAlgorithmRsaOaepSha256 EncryptDataDetailsEncryptionAlgorithmEnum = "RSA_OAEP_SHA_256"
)

var mappingEncryptDataDetailsEncryptionAlgorithmEnum = map[string]EncryptDataDetailsEncryptionAlgorithmEnum{
	"AES_256_GCM":      EncryptDataDetailsEncryptionAlgorithmAes256Gcm,
	"RSA_OAEP_SHA_1":   EncryptDataDetailsEncryptionAlgorithmRsaOaepSha1,
	"RSA_OAEP_SHA_256": EncryptDataDetailsEncryptionAlgorithmRsaOaepSha256,
}

var mappingEncryptDataDetailsEncryptionAlgorithmEnumLowerCase = map[string]EncryptDataDetailsEncryptionAlgorithmEnum{
	"aes_256_gcm":      EncryptDataDetailsEncryptionAlgorithmAes256Gcm,
	"rsa_oaep_sha_1":   EncryptDataDetailsEncryptionAlgorithmRsaOaepSha1,
	"rsa_oaep_sha_256": EncryptDataDetailsEncryptionAlgorithmRsaOaepSha256,
}

// GetEncryptDataDetailsEncryptionAlgorithmEnumValues Enumerates the set of values for EncryptDataDetailsEncryptionAlgorithmEnum
func GetEncryptDataDetailsEncryptionAlgorithmEnumValues() []EncryptDataDetailsEncryptionAlgorithmEnum {
	values := make([]EncryptDataDetailsEncryptionAlgorithmEnum, 0)
	for _, v := range mappingEncryptDataDetailsEncryptionAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetEncryptDataDetailsEncryptionAlgorithmEnumStringValues Enumerates the set of values in String for EncryptDataDetailsEncryptionAlgorithmEnum
func GetEncryptDataDetailsEncryptionAlgorithmEnumStringValues() []string {
	return []string{
		"AES_256_GCM",
		"RSA_OAEP_SHA_1",
		"RSA_OAEP_SHA_256",
	}
}

// GetMappingEncryptDataDetailsEncryptionAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingEncryptDataDetailsEncryptionAlgorithmEnum(val string) (EncryptDataDetailsEncryptionAlgorithmEnum, bool) {
	enum, ok := mappingEncryptDataDetailsEncryptionAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
