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

// Drg A dynamic routing gateway (DRG), which is a virtual router that provides a path for private
// network traffic between your VCN and your existing network. You use it with other Networking
// Service components to create an IPSec VPN or a connection that uses
// Oracle Cloud Infrastructure FastConnect. For more information, see
// Overview of the Networking Service (https://docs.cloud.oracle.com/Content/Network/Concepts/overview.htm).
// To use any of the API operations, you must be authorized in an IAM policy. If you're not authorized,
// talk to an administrator. If you're an administrator who needs to write policies to give users access, see
// Getting Started with Policies (https://docs.cloud.oracle.com/Content/Identity/Concepts/policygetstarted.htm).
// **Warning:** Oracle recommends that you avoid using any confidential information when you
// supply string values using the API.
type Drg struct {

	// The OCID of the compartment containing the DRG.
	CompartmentId *string `mandatory:"true" json:"compartmentId"`

	// The DRG's Oracle ID (OCID).
	Id *string `mandatory:"true" json:"id"`

	// The DRG's current state.
	LifecycleState DrgLifecycleStateEnum `mandatory:"true" json:"lifecycleState"`

	// Defined tags for this resource. Each key is predefined and scoped to a
	// namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Operations": {"CostCenter": "42"}}`
	DefinedTags map[string]map[string]interface{} `mandatory:"false" json:"definedTags"`

	// A user-friendly name. Does not have to be unique, and it's changeable.
	// Avoid entering confidential information.
	DisplayName *string `mandatory:"false" json:"displayName"`

	// Free-form tags for this resource. Each tag is a simple key-value pair with no
	// predefined name, type, or namespace. For more information, see Resource Tags (https://docs.cloud.oracle.com/Content/General/Concepts/resourcetags.htm).
	// Example: `{"Department": "Finance"}`
	FreeformTags map[string]string `mandatory:"false" json:"freeformTags"`

	// The date and time the DRG was created, in the format defined by RFC3339.
	// Example: `2016-08-25T21:10:29.600Z`
	TimeCreated *common.SDKTime `mandatory:"false" json:"timeCreated"`
}

func (m Drg) String() string {
	return common.PointerString(m)
}

// DrgLifecycleStateEnum Enum with underlying type: string
type DrgLifecycleStateEnum string

// Set of constants representing the allowable values for DrgLifecycleStateEnum
const (
	DrgLifecycleStateProvisioning DrgLifecycleStateEnum = "PROVISIONING"
	DrgLifecycleStateAvailable    DrgLifecycleStateEnum = "AVAILABLE"
	DrgLifecycleStateTerminating  DrgLifecycleStateEnum = "TERMINATING"
	DrgLifecycleStateTerminated   DrgLifecycleStateEnum = "TERMINATED"
)

var mappingDrgLifecycleState = map[string]DrgLifecycleStateEnum{
	"PROVISIONING": DrgLifecycleStateProvisioning,
	"AVAILABLE":    DrgLifecycleStateAvailable,
	"TERMINATING":  DrgLifecycleStateTerminating,
	"TERMINATED":   DrgLifecycleStateTerminated,
}

// GetDrgLifecycleStateEnumValues Enumerates the set of values for DrgLifecycleStateEnum
func GetDrgLifecycleStateEnumValues() []DrgLifecycleStateEnum {
	values := make([]DrgLifecycleStateEnum, 0)
	for _, v := range mappingDrgLifecycleState {
		values = append(values, v)
	}
	return values
}
