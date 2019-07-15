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

// AuthenticationPolicy Authentication policy, currently set for the given compartment
type AuthenticationPolicy struct {

	// Password policy.
	PasswordPolicy *PasswordPolicy `mandatory:"false" json:"passwordPolicy"`

	// Compartment OCID.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`
}

func (m AuthenticationPolicy) String() string {
	return common.PointerString(m)
}
