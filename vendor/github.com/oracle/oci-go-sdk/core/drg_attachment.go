// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Core Services API
//
// API covering the Networking (https://docs.cloud.oracle.com/iaas/Content/Network/Concepts/overview.htm),
// Compute (https://docs.cloud.oracle.com/iaas/Content/Compute/Concepts/computeoverview.htm), and
// Block Volume (https://docs.cloud.oracle.com/iaas/Content/Block/Concepts/overview.htm) services. Use this API
// to manage resources such as virtual cloud networks (VCNs), compute instances, and
// block storage volumes.
//

package core

import (
	"github.com/oracle/oci-go-sdk/common"
)

// DrgAttachment A link between a DRG and VCN. For more information, see
// Overview of the Networking Service (https://docs.cloud.oracle.com/Content/Network/Concepts/overview.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type DrgAttachment struct {

	// The OCID of the compartment containing the DRG attachment.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The OCID of the DRG.
	DrgId *string `mandatory:"true" json:"drgId"`

	// The DRG attachment's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The DRG attachment's current state.
	LifecycleState DrgAttachmentLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// The OCID of the VCN.
	VcnId *string `mandatory:"true" json:"vcnId"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The OCID of the route table the DRG attachment is using. For information about why you
	// would associate a route table with a DRG attachment, see
	// Advanced Scenario: Transit Routing (https://docs.cloud.oracle.com/Content/Network/Tasks/transitrouting.htm).
	RouteTableId *string `mandatory:"false" json:"routeTableId"`

	// The date and time the DRG attachment was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m DrgAttachment) String() string {
	return common.PointerString(m)
}

// DrgAttachmentLifecycleStateEnum Enum with underlying type: string
type DrgAttachmentLifecycleStateEnum string

// Set of constants representing the allowable values for DrgAttachmentLifecycleStateEnum
const (
	DrgAttachmentLifecycleStateAttaching DrgAttachmentLifecycleStateEnum = "ATTACHING"
	DrgAttachmentLifecycleStateAttached  DrgAttachmentLifecycleStateEnum = "ATTACHED"
	DrgAttachmentLifecycleStateDetaching DrgAttachmentLifecycleStateEnum = "DETACHING"
	DrgAttachmentLifecycleStateDetached  DrgAttachmentLifecycleStateEnum = "DETACHED"
)

var mappingDrgAttachmentLifecycleState = map[string]DrgAttachmentLifecycleStateEnum{
	"ATTACHING": DrgAttachmentLifecycleStateAttaching,
	"ATTACHED":  DrgAttachmentLifecycleStateAttached,
	"DETACHING": DrgAttachmentLifecycleStateDetaching,
	"DETACHED":  DrgAttachmentLifecycleStateDetached,
}

// GetDrgAttachmentLifecycleStateEnumValues Enumerates the set of values for DrgAttachmentLifecycleStateEnum
func GetDrgAttachmentLifecycleStateEnumValues() []DrgAttachmentLifecycleStateEnum {
	values := make([]DrgAttachmentLifecycleStateEnum, 0)
	for _, v := range mappingDrgAttachmentLifecycleState {
		values = append(values, v)
	}
	return values
}
