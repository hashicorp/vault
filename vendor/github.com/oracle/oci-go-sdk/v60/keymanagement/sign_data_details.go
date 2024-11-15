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

// SignDataDetails The representation of SignDataDetails
type SignDataDetails struct {

	// The base64-encoded binary data object denoting the message or message digest to sign. You can have a message up to 4096 bytes in size. To sign a larger message, provide the message digest.
	Message *string `mandatory:"true" json:"message"`

	// The OCID of the key used to sign the message.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The algorithm to use to sign the message or message digest.
	// For RSA keys, supported signature schemes include PKCS #1 and RSASSA-PSS, along with
	// different hashing algorithms.
	// For ECDSA keys, ECDSA is the supported signature scheme with different hashing algorithms.
	// When you pass a message digest for signing, ensure that you specify the same hashing algorithm
	// as used when creating the message digest.
	SigningAlgorithm SignDataDetailsSigningAlgorithmEnum `mandatory:"true" json:"signingAlgorithm"`

	// The OCID of the key version used to sign the message.
	KeyVersionId *string `mandatory:"false" json:"keyVersionId"`

	// Denotes whether the value of the message parameter is a raw message or a message digest.
	// The default value, `RAW`, indicates a message. To indicate a message digest, use `DIGEST`.
	MessageType SignDataDetailsMessageTypeEnum `mandatory:"false" json:"messageType,omitempty"`
}

func (m SignDataDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m SignDataDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingSignDataDetailsSigningAlgorithmEnum(string(m.SigningAlgorithm)); !ok && m.SigningAlgorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for SigningAlgorithm: %s. Supported values are: %s.", m.SigningAlgorithm, strings.Join(GetSignDataDetailsSigningAlgorithmEnumStringValues(), ",")))
	}

	if _, ok := GetMappingSignDataDetailsMessageTypeEnum(string(m.MessageType)); !ok && m.MessageType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for MessageType: %s. Supported values are: %s.", m.MessageType, strings.Join(GetSignDataDetailsMessageTypeEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// SignDataDetailsMessageTypeEnum Enum with underlying type: string
type SignDataDetailsMessageTypeEnum string

// Set of constants representing the allowable values for SignDataDetailsMessageTypeEnum
const (
	SignDataDetailsMessageTypeRaw    SignDataDetailsMessageTypeEnum = "RAW"
	SignDataDetailsMessageTypeDigest SignDataDetailsMessageTypeEnum = "DIGEST"
)

var mappingSignDataDetailsMessageTypeEnum = map[string]SignDataDetailsMessageTypeEnum{
	"RAW":    SignDataDetailsMessageTypeRaw,
	"DIGEST": SignDataDetailsMessageTypeDigest,
}

var mappingSignDataDetailsMessageTypeEnumLowerCase = map[string]SignDataDetailsMessageTypeEnum{
	"raw":    SignDataDetailsMessageTypeRaw,
	"digest": SignDataDetailsMessageTypeDigest,
}

// GetSignDataDetailsMessageTypeEnumValues Enumerates the set of values for SignDataDetailsMessageTypeEnum
func GetSignDataDetailsMessageTypeEnumValues() []SignDataDetailsMessageTypeEnum {
	values := make([]SignDataDetailsMessageTypeEnum, 0)
	for _, v := range mappingSignDataDetailsMessageTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetSignDataDetailsMessageTypeEnumStringValues Enumerates the set of values in String for SignDataDetailsMessageTypeEnum
func GetSignDataDetailsMessageTypeEnumStringValues() []string {
	return []string{
		"RAW",
		"DIGEST",
	}
}

// GetMappingSignDataDetailsMessageTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingSignDataDetailsMessageTypeEnum(val string) (SignDataDetailsMessageTypeEnum, bool) {
	enum, ok := mappingSignDataDetailsMessageTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// SignDataDetailsSigningAlgorithmEnum Enum with underlying type: string
type SignDataDetailsSigningAlgorithmEnum string

// Set of constants representing the allowable values for SignDataDetailsSigningAlgorithmEnum
const (
	SignDataDetailsSigningAlgorithmSha224RsaPkcsPss  SignDataDetailsSigningAlgorithmEnum = "SHA_224_RSA_PKCS_PSS"
	SignDataDetailsSigningAlgorithmSha256RsaPkcsPss  SignDataDetailsSigningAlgorithmEnum = "SHA_256_RSA_PKCS_PSS"
	SignDataDetailsSigningAlgorithmSha384RsaPkcsPss  SignDataDetailsSigningAlgorithmEnum = "SHA_384_RSA_PKCS_PSS"
	SignDataDetailsSigningAlgorithmSha512RsaPkcsPss  SignDataDetailsSigningAlgorithmEnum = "SHA_512_RSA_PKCS_PSS"
	SignDataDetailsSigningAlgorithmSha224RsaPkcs1V15 SignDataDetailsSigningAlgorithmEnum = "SHA_224_RSA_PKCS1_V1_5"
	SignDataDetailsSigningAlgorithmSha256RsaPkcs1V15 SignDataDetailsSigningAlgorithmEnum = "SHA_256_RSA_PKCS1_V1_5"
	SignDataDetailsSigningAlgorithmSha384RsaPkcs1V15 SignDataDetailsSigningAlgorithmEnum = "SHA_384_RSA_PKCS1_V1_5"
	SignDataDetailsSigningAlgorithmSha512RsaPkcs1V15 SignDataDetailsSigningAlgorithmEnum = "SHA_512_RSA_PKCS1_V1_5"
	SignDataDetailsSigningAlgorithmEcdsaSha256       SignDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_256"
	SignDataDetailsSigningAlgorithmEcdsaSha384       SignDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_384"
	SignDataDetailsSigningAlgorithmEcdsaSha512       SignDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_512"
)

var mappingSignDataDetailsSigningAlgorithmEnum = map[string]SignDataDetailsSigningAlgorithmEnum{
	"SHA_224_RSA_PKCS_PSS":   SignDataDetailsSigningAlgorithmSha224RsaPkcsPss,
	"SHA_256_RSA_PKCS_PSS":   SignDataDetailsSigningAlgorithmSha256RsaPkcsPss,
	"SHA_384_RSA_PKCS_PSS":   SignDataDetailsSigningAlgorithmSha384RsaPkcsPss,
	"SHA_512_RSA_PKCS_PSS":   SignDataDetailsSigningAlgorithmSha512RsaPkcsPss,
	"SHA_224_RSA_PKCS1_V1_5": SignDataDetailsSigningAlgorithmSha224RsaPkcs1V15,
	"SHA_256_RSA_PKCS1_V1_5": SignDataDetailsSigningAlgorithmSha256RsaPkcs1V15,
	"SHA_384_RSA_PKCS1_V1_5": SignDataDetailsSigningAlgorithmSha384RsaPkcs1V15,
	"SHA_512_RSA_PKCS1_V1_5": SignDataDetailsSigningAlgorithmSha512RsaPkcs1V15,
	"ECDSA_SHA_256":          SignDataDetailsSigningAlgorithmEcdsaSha256,
	"ECDSA_SHA_384":          SignDataDetailsSigningAlgorithmEcdsaSha384,
	"ECDSA_SHA_512":          SignDataDetailsSigningAlgorithmEcdsaSha512,
}

var mappingSignDataDetailsSigningAlgorithmEnumLowerCase = map[string]SignDataDetailsSigningAlgorithmEnum{
	"sha_224_rsa_pkcs_pss":   SignDataDetailsSigningAlgorithmSha224RsaPkcsPss,
	"sha_256_rsa_pkcs_pss":   SignDataDetailsSigningAlgorithmSha256RsaPkcsPss,
	"sha_384_rsa_pkcs_pss":   SignDataDetailsSigningAlgorithmSha384RsaPkcsPss,
	"sha_512_rsa_pkcs_pss":   SignDataDetailsSigningAlgorithmSha512RsaPkcsPss,
	"sha_224_rsa_pkcs1_v1_5": SignDataDetailsSigningAlgorithmSha224RsaPkcs1V15,
	"sha_256_rsa_pkcs1_v1_5": SignDataDetailsSigningAlgorithmSha256RsaPkcs1V15,
	"sha_384_rsa_pkcs1_v1_5": SignDataDetailsSigningAlgorithmSha384RsaPkcs1V15,
	"sha_512_rsa_pkcs1_v1_5": SignDataDetailsSigningAlgorithmSha512RsaPkcs1V15,
	"ecdsa_sha_256":          SignDataDetailsSigningAlgorithmEcdsaSha256,
	"ecdsa_sha_384":          SignDataDetailsSigningAlgorithmEcdsaSha384,
	"ecdsa_sha_512":          SignDataDetailsSigningAlgorithmEcdsaSha512,
}

// GetSignDataDetailsSigningAlgorithmEnumValues Enumerates the set of values for SignDataDetailsSigningAlgorithmEnum
func GetSignDataDetailsSigningAlgorithmEnumValues() []SignDataDetailsSigningAlgorithmEnum {
	values := make([]SignDataDetailsSigningAlgorithmEnum, 0)
	for _, v := range mappingSignDataDetailsSigningAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetSignDataDetailsSigningAlgorithmEnumStringValues Enumerates the set of values in String for SignDataDetailsSigningAlgorithmEnum
func GetSignDataDetailsSigningAlgorithmEnumStringValues() []string {
	return []string{
		"SHA_224_RSA_PKCS_PSS",
		"SHA_256_RSA_PKCS_PSS",
		"SHA_384_RSA_PKCS_PSS",
		"SHA_512_RSA_PKCS_PSS",
		"SHA_224_RSA_PKCS1_V1_5",
		"SHA_256_RSA_PKCS1_V1_5",
		"SHA_384_RSA_PKCS1_V1_5",
		"SHA_512_RSA_PKCS1_V1_5",
		"ECDSA_SHA_256",
		"ECDSA_SHA_384",
		"ECDSA_SHA_512",
	}
}

// GetMappingSignDataDetailsSigningAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingSignDataDetailsSigningAlgorithmEnum(val string) (SignDataDetailsSigningAlgorithmEnum, bool) {
	enum, ok := mappingSignDataDetailsSigningAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
