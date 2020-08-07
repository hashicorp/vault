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

// EncryptedData The representation of EncryptedData
type EncryptedData struct {

	// The encrypted data.
	Ciphertext *string `mandatory:"true" json:"ciphertext"`
}

func (m EncryptedData) String() string {
	return common.PointerString(m)
}
