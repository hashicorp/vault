// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Resource Manager API
//
// API for the Resource Manager service. Use this API to install, configure, and manage resources via the "infrastructure-as-code" model. For more information, see Overview of Resource Manager (https://docs.cloud.oracle.com/iaas/Content/ResourceManager/Concepts/resourcemanager.htm).
//

package resourcemanager

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ChangeStackCompartmentDetails Defines the requirements and properties of changeStackCompartment operation.
type ChangeStackCompartmentDetails struct {

	// The OCID (https://docs.cloud.oracle.com/Content/General/Concepts/identifiers.htm) of the compartment
	// into which the Stack should be moved.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeStackCompartmentDetails) String() string {
	return common.PointerString(m)
}
