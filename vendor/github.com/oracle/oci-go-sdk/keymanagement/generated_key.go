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

// GeneratedKey The representation of GeneratedKey
type GeneratedKey struct {

	// The encrypted data encryption key generated from a master encryption key.
	Ciphertext *string `mandatory:"true" json:"ciphertext"`

	// The plaintext data encryption key, a base64-encoded sequence of random bytes, which is
	// included if the GenerateDataEncryptionKey (https://docs.cloud.oracle.com/api/#/en/key/release/GeneratedKey/GenerateDataEncryptionKey)
	// request includes the `includePlaintextKey` parameter and sets its value to "true".
	Plaintext *string `mandatory:"false" json:"plaintext"`

	// The checksum of the plaintext data encryption key, which is included if the
	// GenerateDataEncryptionKey (https://docs.cloud.oracle.com/api/#/en/key/release/GeneratedKey/GenerateDataEncryptionKey)
	// request includes the `includePlaintextKey` parameter and sets its value to "true".
	PlaintextChecksum *string `mandatory:"false" json:"plaintextChecksum"`
}

func (m GeneratedKey) String() string {
	return common.PointerString(m)
}
