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

// PasswordPolicy Password policy, currently set for the given compartment.
type PasswordPolicy struct {

	// Minimum password length required.
	MinimumPasswordLength *int `mandatory:"false" json:"minimumPasswordLength"`

	// At least one uppercase character required.
	IsUppercaseCharactersRequired *bool `mandatory:"false" json:"isUppercaseCharactersRequired"`

	// At least one lower case character required.
	IsLowercaseCharactersRequired *bool `mandatory:"false" json:"isLowercaseCharactersRequired"`

	// At least one numeric character required.
	IsNumericCharactersRequired *bool `mandatory:"false" json:"isNumericCharactersRequired"`

	// At least one special character required.
	IsSpecialCharactersRequired *bool `mandatory:"false" json:"isSpecialCharactersRequired"`

	// User name is allowed to be part of the password.
	IsUsernameContainmentAllowed *bool `mandatory:"false" json:"isUsernameContainmentAllowed"`
}

func (m PasswordPolicy) String() string {
	return common.PointerString(m)
}
