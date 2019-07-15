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

// CreateCustomerSecretKeyDetails The representation of CreateCustomerSecretKeyDetails
type CreateCustomerSecretKeyDetails struct {

	// The name you assign to the secret key during creation. Does not have to be unique, and it's changeable.
	DisplayName *string `mandatory:"true" json:"displayName"`
}

func (m CreateCustomerSecretKeyDetails) String() string {
	return common.PointerString(m)
}
