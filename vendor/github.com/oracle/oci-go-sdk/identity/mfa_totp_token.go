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

// MfaTotpToken Totp token for MFA
type MfaTotpToken struct {

	// The Totp token for MFA.
	TotpToken *string `mandatory:"false" json:"totpToken"`
}

func (m MfaTotpToken) String() string {
	return common.PointerString(m)
}
