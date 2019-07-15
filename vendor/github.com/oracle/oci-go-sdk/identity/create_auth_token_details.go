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

// CreateAuthTokenDetails The representation of CreateAuthTokenDetails
type CreateAuthTokenDetails struct {

	// The description you assign to the auth token during creation. Does not have to be unique, and it's changeable.
	Description *string `mandatory:"true" json:"description"`
}

func (m CreateAuthTokenDetails) String() string {
	return common.PointerString(m)
}
