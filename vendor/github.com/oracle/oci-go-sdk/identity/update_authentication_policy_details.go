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

// UpdateAuthenticationPolicyDetails Update request for authentication policy, describes set of validation rules and their parameters to be updated.
type UpdateAuthenticationPolicyDetails struct {

	// Password policy.
	PasswordPolicy *PasswordPolicy `mandatory:"false" json:"passwordPolicy"`
}

func (m UpdateAuthenticationPolicyDetails) String() string {
	return common.PointerString(m)
}
