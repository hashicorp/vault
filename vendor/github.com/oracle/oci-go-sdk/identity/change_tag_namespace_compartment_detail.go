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

// ChangeTagNamespaceCompartmentDetail Details of the compartment the resource is being moved to.
type ChangeTagNamespaceCompartmentDetail struct {

	// The Oracle Cloud ID (OCID) of the destination compartment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeTagNamespaceCompartmentDetail) String() string {
	return common.PointerString(m)
}
