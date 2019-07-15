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

// UpdateTagDefaultDetails The representation of UpdateTagDefaultDetails
type UpdateTagDefaultDetails struct {

	// The default value for the tag definition. This will be applied to all resources created in the Compartment.
	Value *string `mandatory:"true" json:"value"`
}

func (m UpdateTagDefaultDetails) String() string {
	return common.PointerString(m)
}
