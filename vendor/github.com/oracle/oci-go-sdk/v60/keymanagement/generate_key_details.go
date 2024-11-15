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

// GenerateKeyDetails The representation of GenerateKeyDetails
type GenerateKeyDetails struct {

	// If true, the generated key is also returned unencrypted.
	IncludePlaintextKey *bool `mandatory:"true" json:"includePlaintextKey"`

	// The OCID of the master encryption key to encrypt the generated data encryption key with.
	KeyId *string `mandatory:"true" json:"keyId"`

	KeyShape *KeyShape `mandatory:"true" json:"keyShape"`

	// Information that can be used to provide an encryption context for the encrypted data.
	// The length of the string representation of the associated data must be fewer than 4096
	// characters.
	AssociatedData map[string]string `mandatory:"false" json:"associatedData"`

	// Information that provides context for audit logging. You can provide this additional
	// data by formatting it as key-value pairs to include in audit logs when audit logging is enabled.
	LoggingContext map[string]string `mandatory:"false" json:"loggingContext"`
}

func (m GenerateKeyDetails) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m GenerateKeyDetails) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}
