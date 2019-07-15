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

// UpdateCustomerSecretKeyDetails The representation of UpdateCustomerSecretKeyDetails
type UpdateCustomerSecretKeyDetails struct {

	// The description you assign to the secret key. Does not have to be unique, and it's changeable.
	DisplayName *string `mandatory:"false" json:"displayName"`
}

func (m UpdateCustomerSecretKeyDetails) String() string {
	return common.PointerString(m)
}
