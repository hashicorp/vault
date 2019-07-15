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

// IdentityProviderGroupSummary A group created in an identity provider that can be mapped to a group in OCI
type IdentityProviderGroupSummary struct {

	// The OCID of the `IdentityProviderGroup`.
	Id *string `mandatory:"false" json:"id"`

	// The OCID of the `IdentityProvider` this group belongs to.
	IdentityProviderId *string `mandatory:"false" json:"identityProviderId"`

	// Display name of the group
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Identifier of the group in the identity provider
	ExternalIdentifier *string `mandatory:"false" json:"externalIdentifier"`

	// Date and time the `IdentityProviderGroup` was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`

	// Date and time the `IdentityProviderGroup` was last modified, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeModified *common.SDKTime `mandatory:"false" json:"timeModified"`
}

func (m IdentityProviderGroupSummary) String() string {
	return common.PointerString(m)
}
