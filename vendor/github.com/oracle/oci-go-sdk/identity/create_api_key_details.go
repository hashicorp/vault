// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Identity and Access Management Service API
//
// APIs for managing users, groups, compartments, and policies.
//

package identity

import (
	"github.com/oracle/oci-go-sdk/common"
)

// CreateApiKeyDetails The representation of CreateApiKeyDetails
type CreateApiKeyDetails struct {

	// The public key.  Must be an RSA key in PEM format.
	Key *string `mandatory:"true" json:"key"`
}

func (m CreateApiKeyDetails) String() string {
	return common.PointerString(m)
}
