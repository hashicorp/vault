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

// UpdateIdpGroupMappingDetails The representation of UpdateIdpGroupMappingDetails
type UpdateIdpGroupMappingDetails struct {

	// The idp group name.
	IdpGroupName *string `mandatory:"false" json:"idpGroupName"`

	// The OCID of the group.
	GroupId *string `mandatory:"false" json:"groupId"`
}

func (m UpdateIdpGroupMappingDetails) String() string {
	return common.PointerString(m)
}
