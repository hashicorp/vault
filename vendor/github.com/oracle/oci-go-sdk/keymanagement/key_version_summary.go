// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// KeyVersionSummary The representation of KeyVersionSummary
type KeyVersionSummary struct {

	// The OCID of the compartment that contains this key version.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the key version.
	Id *string `mandatory:"true" json:"id"`

	// The OCID of the key associated with this key version.
	KeyId *string `mandatory:"true" json:"keyId"`

	// The date and time this key version was created, expressed in RFC 3339 (https://tools.ietf.org/html/rfc3339) timestamp format.
	// Example: `2018-04-03T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"true" json:"timeCreated"`

	// The OCID of the vault that contains this key version.
	VaultId *string `mandatory:"true" json:"vaultId"`
}

func (m KeyVersionSummary) String() string {
	return common.PointerString(m)
}
