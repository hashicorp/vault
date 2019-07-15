// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements.
// For information about the Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
//

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/common"
)

// ChangeAutoScalingCompartmentDetails The configuration details for the move operation.
type ChangeAutoScalingCompartmentDetails struct {

	// The OCID (https://docs.cloud.oracle.com/iaas/Content/General/Concepts/identifiers.htm) of the compartment to move the autoscaling configuration to.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`
}

func (m ChangeAutoScalingCompartmentDetails) String() string {
	return common.PointerString(m)
}
