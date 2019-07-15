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

// CrossConnectGroup For use with Oracle Cloud Infrastructure FastConnect. A cross-connect group
// is a link aggregation group (LAG), which can contain one or more
// CrossConnect. Customers who are colocated with
// Oracle in a FastConnect location create and use cross-connect groups. For more
// information, see FastConnect Overview (https://docs.cloud.oracle.com/Content/Network/Concepts/fastconnect.htm).
// **Note:** If you're a provider who is setting up a physical connection to Oracle so customers
// can use FastConnect over the connection, be aware that your connection is modeled the
// same way as a colocated customer's (with `CrossConnect` and `CrossConnectGroup` objects, and so on).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type CrossConnectGroup struct {

	// The OCID of the compartment containing the cross-connect group.
	CompartmentId *string `mandatory:"false" json:"compartmentId"`

	// The display name of a user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// The cross-connect group's Oracle ID (OCID).
	Id *string `mandatory:"false" json:"id"`

	// The cross-connect group's current state.
	LifecycleState CrossConnectGroupLifecycleStateEnum `mandatory:"false" json:"lifecycleState,omitempty"`

	// A reference name or identifier for the physical fiber connection that this cross-connect
	// group uses.
	CustomerReferenceName *string `mandatory:"false" json:"customerReferenceName"`

	// The date and time the cross-connect group was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m CrossConnectGroup) String() string {
	return common.PointerString(m)
}

// CrossConnectGroupLifecycleStateEnum Enum with underlying type: string
type CrossConnectGroupLifecycleStateEnum string

// Set of constants representing the allowable values for CrossConnectGroupLifecycleStateEnum
const (
	CrossConnectGroupLifecycleStateProvisioning CrossConnectGroupLifecycleStateEnum = "PROVISIONING"
	CrossConnectGroupLifecycleStateProvisioned  CrossConnectGroupLifecycleStateEnum = "PROVISIONED"
	CrossConnectGroupLifecycleStateInactive     CrossConnectGroupLifecycleStateEnum = "INACTIVE"
	CrossConnectGroupLifecycleStateTerminating  CrossConnectGroupLifecycleStateEnum = "TERMINATING"
	CrossConnectGroupLifecycleStateTerminated   CrossConnectGroupLifecycleStateEnum = "TERMINATED"
)

var mappingCrossConnectGroupLifecycleState = map[string]CrossConnectGroupLifecycleStateEnum{
	"PROVISIONING": CrossConnectGroupLifecycleStateProvisioning,
	"PROVISIONED":  CrossConnectGroupLifecycleStateProvisioned,
	"INACTIVE":     CrossConnectGroupLifecycleStateInactive,
	"TERMINATING":  CrossConnectGroupLifecycleStateTerminating,
	"TERMINATED":   CrossConnectGroupLifecycleStateTerminated,
}

// GetCrossConnectGroupLifecycleStateEnumValues Enumerates the set of values for CrossConnectGroupLifecycleStateEnum
func GetCrossConnectGroupLifecycleStateEnumValues() []CrossConnectGroupLifecycleStateEnum {
	values := make([]CrossConnectGroupLifecycleStateEnum, 0)
	for _, v := range mappingCrossConnectGroupLifecycleState {
		values = append(values, v)
	}
	return values
}
