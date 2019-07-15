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

// UserCapabilities Properties indicating how the user is allowed to authenticate.
type UserCapabilities struct {

	// Indicates if the user can log in to the console.
	CanUseConsolePassword *bool `mandatory:"false" json:"canUseConsolePassword"`

	// Indicates if the user can use API keys.
	CanUseApiKeys *bool `mandatory:"false" json:"canUseApiKeys"`

	// Indicates if the user can use SWIFT passwords / auth tokens.
	CanUseAuthTokens *bool `mandatory:"false" json:"canUseAuthTokens"`

	// Indicates if the user can use SMTP passwords.
	CanUseSmtpCredentials *bool `mandatory:"false" json:"canUseSmtpCredentials"`

	// Indicates if the user can use SigV4 symmetric keys.
	CanUseCustomerSecretKeys *bool `mandatory:"false" json:"canUseCustomerSecretKeys"`
}

func (m UserCapabilities) String() string {
	return common.PointerString(m)
}
