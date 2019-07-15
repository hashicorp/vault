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

// AddUserToGroupDetails The representation of AddUserToGroupDetails
type AddUserToGroupDetails struct {

	// The OCID of the user.
	UserId *string `mandatory:"true" json:"userId"`

	// The OCID of the group.
	GroupId *string `mandatory:"true" json:"groupId"`
}

func (m AddUserToGroupDetails) String() string {
	return common.PointerString(m)
}
