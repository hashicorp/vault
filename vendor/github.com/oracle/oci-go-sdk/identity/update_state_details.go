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

// UpdateStateDetails The representation of UpdateStateDetails
type UpdateStateDetails struct {

	// Update state to blocked or unblocked. Only "false" is supported (for changing the state to unblocked).
	Blocked *bool `mandatory:"false" json:"blocked"`
}

func (m UpdateStateDetails) String() string {
	return common.PointerString(m)
}
