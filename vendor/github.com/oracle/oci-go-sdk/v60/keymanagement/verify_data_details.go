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

// VerifyDataDetails The representation of VerifyDataDetails
type VerifyDataDetails struct {

	// The OCID of the key used to sign the message.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The OCID of the key version used to sign the message.
	KeyVersionId *string `mandatory:"true" json:"keyVersionId"`

	// The base64-encoded binary data object denoting the cryptographic signature generated for the message.
	Signature *string `mandatory:"true" json:"signature"`

	// The base64-encoded binary data object denoting the message or message digest to sign. You can have a message up to 4096 bytes in size. To sign a larger message, provide the message digest.
	Message *string `mandatory:"true" json:"message"`

	// The algorithm to use to sign the message or message digest.
	// For RSA keys, supported signature schemes include PKCS #1 and RSASSA-PSS, along with
	// different hashing algorithms.
	// For ECDSA keys, ECDSA is the supported signature scheme with different hashing algorithms.
	// When you pass a message digest for signing, ensure that you specify the same hashing algorithm
	// as used when creating the message digest.
	SigningAlgorithm VerifyDataDetailsSigningAlgorithmEnum `mandatory:"true" json:"signingAlgorithm"`

	// Denotes whether the value of the message parameter is a raw message or a message digest.
	// The default value, `RAW`, indicates a message. To indicate a message digest, use `DIGEST`.
	MessageType VerifyDataDetailsMessageTypeEnum `mandatory:"false" json:"messageType,omitempty"`
}

func (m VerifyDataDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m VerifyDataDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}
	if _, ok := GetMappingVerifyDataDetailsSigningAlgorithmEnum(string(m.SigningAlgorithm)); !ok && m.SigningAlgorithm != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for SigningAlgorithm: %s. Supported values are: %s.", m.SigningAlgorithm, strings.Join(GetVerifyDataDetailsSigningAlgorithmEnumStringValues(), ",")))
	}

	if _, ok := GetMappingVerifyDataDetailsMessageTypeEnum(string(m.MessageType)); !ok && m.MessageType != "" {
		errMessage = append(errMessage, fmt.Sprintf("unsupported enum value for MessageType: %s. Supported values are: %s.", m.MessageType, strings.Join(GetVerifyDataDetailsMessageTypeEnumStringValues(), ",")))
	}
	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// VerifyDataDetailsMessageTypeEnum Enum with underlying type: string
type VerifyDataDetailsMessageTypeEnum string

// Set of constants representing the allowable values for VerifyDataDetailsMessageTypeEnum
const (
	VerifyDataDetailsMessageTypeRaw    VerifyDataDetailsMessageTypeEnum = "RAW"
	VerifyDataDetailsMessageTypeDigest VerifyDataDetailsMessageTypeEnum = "DIGEST"
)

var mappingVerifyDataDetailsMessageTypeEnum = map[string]VerifyDataDetailsMessageTypeEnum{
	"RAW":    VerifyDataDetailsMessageTypeRaw,
	"DIGEST": VerifyDataDetailsMessageTypeDigest,
}

var mappingVerifyDataDetailsMessageTypeEnumLowerCase = map[string]VerifyDataDetailsMessageTypeEnum{
	"raw":    VerifyDataDetailsMessageTypeRaw,
	"digest": VerifyDataDetailsMessageTypeDigest,
}

// GetVerifyDataDetailsMessageTypeEnumValues Enumerates the set of values for VerifyDataDetailsMessageTypeEnum
func GetVerifyDataDetailsMessageTypeEnumValues() []VerifyDataDetailsMessageTypeEnum {
	values := make([]VerifyDataDetailsMessageTypeEnum, 0)
	for _, v := range mappingVerifyDataDetailsMessageTypeEnum {
		values = append(values, v)
	}
	return values
}

// GetVerifyDataDetailsMessageTypeEnumStringValues Enumerates the set of values in String for VerifyDataDetailsMessageTypeEnum
func GetVerifyDataDetailsMessageTypeEnumStringValues() []string {
	return []string{
		"RAW",
		"DIGEST",
	}
}

// GetMappingVerifyDataDetailsMessageTypeEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVerifyDataDetailsMessageTypeEnum(val string) (VerifyDataDetailsMessageTypeEnum, bool) {
	enum, ok := mappingVerifyDataDetailsMessageTypeEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}

// VerifyDataDetailsSigningAlgorithmEnum Enum with underlying type: string
type VerifyDataDetailsSigningAlgorithmEnum string

// Set of constants representing the allowable values for VerifyDataDetailsSigningAlgorithmEnum
const (
	VerifyDataDetailsSigningAlgorithmSha224RsaPkcsPss  VerifyDataDetailsSigningAlgorithmEnum = "SHA_224_RSA_PKCS_PSS"
	VerifyDataDetailsSigningAlgorithmSha256RsaPkcsPss  VerifyDataDetailsSigningAlgorithmEnum = "SHA_256_RSA_PKCS_PSS"
	VerifyDataDetailsSigningAlgorithmSha384RsaPkcsPss  VerifyDataDetailsSigningAlgorithmEnum = "SHA_384_RSA_PKCS_PSS"
	VerifyDataDetailsSigningAlgorithmSha512RsaPkcsPss  VerifyDataDetailsSigningAlgorithmEnum = "SHA_512_RSA_PKCS_PSS"
	VerifyDataDetailsSigningAlgorithmSha224RsaPkcs1V15 VerifyDataDetailsSigningAlgorithmEnum = "SHA_224_RSA_PKCS1_V1_5"
	VerifyDataDetailsSigningAlgorithmSha256RsaPkcs1V15 VerifyDataDetailsSigningAlgorithmEnum = "SHA_256_RSA_PKCS1_V1_5"
	VerifyDataDetailsSigningAlgorithmSha384RsaPkcs1V15 VerifyDataDetailsSigningAlgorithmEnum = "SHA_384_RSA_PKCS1_V1_5"
	VerifyDataDetailsSigningAlgorithmSha512RsaPkcs1V15 VerifyDataDetailsSigningAlgorithmEnum = "SHA_512_RSA_PKCS1_V1_5"
	VerifyDataDetailsSigningAlgorithmEcdsaSha256       VerifyDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_256"
	VerifyDataDetailsSigningAlgorithmEcdsaSha384       VerifyDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_384"
	VerifyDataDetailsSigningAlgorithmEcdsaSha512       VerifyDataDetailsSigningAlgorithmEnum = "ECDSA_SHA_512"
)

var mappingVerifyDataDetailsSigningAlgorithmEnum = map[string]VerifyDataDetailsSigningAlgorithmEnum{
	"SHA_224_RSA_PKCS_PSS":   VerifyDataDetailsSigningAlgorithmSha224RsaPkcsPss,
	"SHA_256_RSA_PKCS_PSS":   VerifyDataDetailsSigningAlgorithmSha256RsaPkcsPss,
	"SHA_384_RSA_PKCS_PSS":   VerifyDataDetailsSigningAlgorithmSha384RsaPkcsPss,
	"SHA_512_RSA_PKCS_PSS":   VerifyDataDetailsSigningAlgorithmSha512RsaPkcsPss,
	"SHA_224_RSA_PKCS1_V1_5": VerifyDataDetailsSigningAlgorithmSha224RsaPkcs1V15,
	"SHA_256_RSA_PKCS1_V1_5": VerifyDataDetailsSigningAlgorithmSha256RsaPkcs1V15,
	"SHA_384_RSA_PKCS1_V1_5": VerifyDataDetailsSigningAlgorithmSha384RsaPkcs1V15,
	"SHA_512_RSA_PKCS1_V1_5": VerifyDataDetailsSigningAlgorithmSha512RsaPkcs1V15,
	"ECDSA_SHA_256":          VerifyDataDetailsSigningAlgorithmEcdsaSha256,
	"ECDSA_SHA_384":          VerifyDataDetailsSigningAlgorithmEcdsaSha384,
	"ECDSA_SHA_512":          VerifyDataDetailsSigningAlgorithmEcdsaSha512,
}

var mappingVerifyDataDetailsSigningAlgorithmEnumLowerCase = map[string]VerifyDataDetailsSigningAlgorithmEnum{
	"sha_224_rsa_pkcs_pss":   VerifyDataDetailsSigningAlgorithmSha224RsaPkcsPss,
	"sha_256_rsa_pkcs_pss":   VerifyDataDetailsSigningAlgorithmSha256RsaPkcsPss,
	"sha_384_rsa_pkcs_pss":   VerifyDataDetailsSigningAlgorithmSha384RsaPkcsPss,
	"sha_512_rsa_pkcs_pss":   VerifyDataDetailsSigningAlgorithmSha512RsaPkcsPss,
	"sha_224_rsa_pkcs1_v1_5": VerifyDataDetailsSigningAlgorithmSha224RsaPkcs1V15,
	"sha_256_rsa_pkcs1_v1_5": VerifyDataDetailsSigningAlgorithmSha256RsaPkcs1V15,
	"sha_384_rsa_pkcs1_v1_5": VerifyDataDetailsSigningAlgorithmSha384RsaPkcs1V15,
	"sha_512_rsa_pkcs1_v1_5": VerifyDataDetailsSigningAlgorithmSha512RsaPkcs1V15,
	"ecdsa_sha_256":          VerifyDataDetailsSigningAlgorithmEcdsaSha256,
	"ecdsa_sha_384":          VerifyDataDetailsSigningAlgorithmEcdsaSha384,
	"ecdsa_sha_512":          VerifyDataDetailsSigningAlgorithmEcdsaSha512,
}

// GetVerifyDataDetailsSigningAlgorithmEnumValues Enumerates the set of values for VerifyDataDetailsSigningAlgorithmEnum
func GetVerifyDataDetailsSigningAlgorithmEnumValues() []VerifyDataDetailsSigningAlgorithmEnum {
	values := make([]VerifyDataDetailsSigningAlgorithmEnum, 0)
	for _, v := range mappingVerifyDataDetailsSigningAlgorithmEnum {
		values = append(values, v)
	}
	return values
}

// GetVerifyDataDetailsSigningAlgorithmEnumStringValues Enumerates the set of values in String for VerifyDataDetailsSigningAlgorithmEnum
func GetVerifyDataDetailsSigningAlgorithmEnumStringValues() []string {
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

// GetMappingVerifyDataDetailsSigningAlgorithmEnum performs case Insensitive comparison on enum value and return the desired enum
func GetMappingVerifyDataDetailsSigningAlgorithmEnum(val string) (VerifyDataDetailsSigningAlgorithmEnum, bool) {
	enum, ok := mappingVerifyDataDetailsSigningAlgorithmEnumLowerCase[strings.ToLower(val)]
	return enum, ok
}
