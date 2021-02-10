// Copyright (c) 2016, 2018, 2020, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Key Management Service API
//
// API for managing and performing operations with keys and vaults.
//

package keymanagement

import (
	"github.com/oracle/oci-go-sdk/common"
)

// DecryptedData The representation of DecryptedData
type DecryptedData struct {

	// The decrypted data, expressed as a base64-encoded value.
	Plaintext *string `mandatory:"true" json:"plaintext"`

	// Checksum of the decrypted data.
	PlaintextChecksum *string `mandatory:"true" json:"plaintextChecksum"`
}

func (m DecryptedData) String() string {
	return common.PointerString(m)
}
