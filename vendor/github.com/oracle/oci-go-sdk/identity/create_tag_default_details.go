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

// CreateTagDefaultDetails The representation of CreateTagDefaultDetails
type CreateTagDefaultDetails struct {

	// The OCID of the compartment. The tag default will be applied to all new resources created in this compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the tag definition. The tag default will always assign a default value for this tag definition.
	TagDefinitionId *string `mandatory:"true" json:"tagDefinitionId"`

	// The default value for the tag definition. This will be applied to all new resources created in the compartment.
	Value *string `mandatory:"true" json:"value"`
}

func (m CreateTagDefaultDetails) String() string {
	return common.PointerString(m)
}
